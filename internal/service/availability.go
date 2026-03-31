package service

import (
	"context"
	"time"

	"booking/internal/domain"
	"booking/internal/dto"
	"booking/internal/errs"
)

type ReservationLister interface {
	ListByRoomAndRange(ctx context.Context, roomID int64, from, to time.Time) ([]domain.Reservation, error)
}

type BlockLister interface {
	ListByRoomAndRange(ctx context.Context, roomID int64, from, to time.Time) ([]domain.RoomBlock, error)
}

type AvailabilityService struct {
	resRepo ReservationLister
	blkRepo BlockLister
}

func NewAvailabilityService(resRepo ReservationLister, blkRepo BlockLister) *AvailabilityService {
	return &AvailabilityService{resRepo: resRepo, blkRepo: blkRepo}
}

func (s *AvailabilityService) Get(ctx context.Context, roomID int64, from, to time.Time) (dto.AvailabilityResponse, error) {
	if roomID <= 0 || from.IsZero() || to.IsZero() || !from.Before(to) {
		return dto.AvailabilityResponse{}, errs.ErrValidation
	}

	from = from.UTC()
	to = to.UTC()

	reservations, err := s.resRepo.ListByRoomAndRange(ctx, roomID, from, to)
	if err != nil {
		return dto.AvailabilityResponse{}, err
	}

	blocks, err := s.blkRepo.ListByRoomAndRange(ctx, roomID, from, to)
	if err != nil {
		return dto.AvailabilityResponse{}, err
	}

	out := dto.AvailabilityResponse{
		RoomID:       roomID,
		From:         from,
		To:           to,
		Reservations: make([]dto.Interval, 0, len(reservations)),
		Blocks:       make([]dto.Interval, 0, len(blocks)),
	}

	for _, r := range reservations {
		out.Reservations = append(out.Reservations, dto.Interval{StartTime: r.StartTime, EndTime: r.EndTime})
	}
	for _, b := range blocks {
		out.Blocks = append(out.Blocks, dto.Interval{StartTime: b.StartTime, EndTime: b.EndTime})
	}

	return out, nil
}