package service

import (
	"context"
	"time"

	"booking/internal/domain"
	"booking/internal/dto"
	"booking/internal/errs"
)

type ReservationRepository interface {
	Create(ctx context.Context, res domain.Reservation) (domain.Reservation, error)
	GetByID(ctx context.Context, id int64) (domain.Reservation, error)
	Cancel(ctx context.Context, id int64, reason *string) error
	ListByUser(ctx context.Context, userID int64, status string) ([]domain.Reservation, error)
}

type RoomGetter interface {
	GetByID(ctx context.Context, id int64) (domain.Room, error)
}

type ReservationService struct {
	repo  ReservationRepository
	rooms RoomGetter
}

func NewReservationService(repo ReservationRepository, rooms RoomGetter) *ReservationService {
	return &ReservationService{repo: repo, rooms: rooms}
}

func (s *ReservationService) Create(ctx context.Context, userID int64, req dto.CreateReservationRequest) (domain.Reservation, error) {
	if req.RoomID <= 0 || req.StartTime.IsZero() || req.EndTime.IsZero() {
		return domain.Reservation{}, errs.ErrValidation
	}

	start := req.StartTime.UTC()
	end := req.EndTime.UTC()

	if !start.Before(end) {
		return domain.Reservation{}, errs.ErrValidation
	}

	// пример ограничений (можешь менять)
	if end.Sub(start) < 15*time.Minute || end.Sub(start) > 8*time.Hour {
		return domain.Reservation{}, errs.ErrValidation
	}

	// нельзя бронировать в прошлом (маленькая “погрешность”)
	if start.Before(time.Now().UTC().Add(-1 * time.Minute)) {
		return domain.Reservation{}, errs.ErrValidation
	}

	room, err := s.rooms.GetByID(ctx, req.RoomID)
	if err != nil {
		return domain.Reservation{}, err
	}
	if !room.IsActive {
		return domain.Reservation{}, errs.ErrForbidden
	}

	res := domain.Reservation{
		UserID:    userID,
		RoomID:    req.RoomID,
		StartTime: start,
		EndTime:   end,
		Status:    domain.ReservationConfirmed,
	}

	created, err := s.repo.Create(ctx, res)
	if err != nil {
		return domain.Reservation{}, err
	}
	return created, nil
}

func (s *ReservationService) Get(ctx context.Context, actorID int64, actorRole string, id int64) (domain.Reservation, error) {
	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Reservation{}, err
	}

	if actorRole != "admin" && res.UserID != actorID {
		return domain.Reservation{}, errs.ErrForbidden
	}
	return res, nil
}

func (s *ReservationService) ListMy(ctx context.Context, userID int64, status string) ([]domain.Reservation, error) {
	return s.repo.ListByUser(ctx, userID, status)
}

func (s *ReservationService) Cancel(ctx context.Context, actorID int64, actorRole string, id int64, reason string) error {
	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if actorRole != "admin" && res.UserID != actorID {
		return errs.ErrForbidden
	}

	// нельзя отменять уже отменённое/истёкшее
	if res.Status == domain.ReservationCancelled || res.Status == domain.ReservationExpired {
		return errs.ErrConflict
	}

	// нельзя отменять, если уже началось
	if !time.Now().UTC().Before(res.StartTime) {
		return errs.ErrConflict
	}

	var reasonPtr *string
	if reason != "" {
		reasonPtr = &reason
	}

	return s.repo.Cancel(ctx, id, reasonPtr)
}