package service
//правила и валидность создания/удаления блокировок

import (
	"context"
	"time"

	"booking/internal/domain"
	"booking/internal/dto"
	"booking/internal/errs"
)
//связь между service и repository

type BlockRepository interface {
	Create(ctx context.Context, b domain.RoomBlock) (domain.RoomBlock, error)
	Delete(ctx context.Context, blockID int64) error
}


type BlockService struct {
	repo  BlockRepository
	rooms RoomGetter 
}

func NewBlockService(repo BlockRepository, rooms RoomGetter) *BlockService {
	return &BlockService{repo: repo, rooms: rooms}
}

func (s *BlockService) Create(ctx context.Context, adminID int64, roomID int64, req dto.CreateBlockRequest) (domain.RoomBlock, error) {
	if roomID <= 0 || req.StartTime.IsZero() || req.EndTime.IsZero() {  //невалидный ID
		return domain.RoomBlock{}, errs.ErrValidation
	}

	//нормализация времени в UTC
	start := req.StartTime.UTC()
	end := req.EndTime.UTC()
	
	if !start.Before(end) { //если start >= end интервал неправильный
		return domain.RoomBlock{}, errs.ErrValidation
	}
	if end.Sub(start) < 15*time.Minute || end.Sub(start) > 24*time.Hour { //ограничения на длительность блокировки: 15 минут минимум, 24 часа максимум
		return domain.RoomBlock{}, errs.ErrValidation
	}

	room, err := s.rooms.GetByID(ctx, roomID)

	if err != nil { //если комната не найдена ошибка
		return domain.RoomBlock{}, err 
	}
	if !room.IsActive { // запрет доступа
		return domain.RoomBlock{}, errs.ErrForbidden
	}

	createdBy := adminID
	//сбор домена, который в дальнейшем уйдет в бд слой
	b := domain.RoomBlock{
		RoomID:    roomID,
		StartTime: start,
		EndTime:   end,
		Reason:    req.Reason,
		CreatedBy: &createdBy,
	}
	return s.repo.Create(ctx, b)
}


//сервис не даёт даже пытаться удалять "0" или отрицательные ID
func (s *BlockService) Delete(ctx context.Context, blockID int64) error {
	if blockID <= 0 {
		return errs.ErrValidation
	}
	return s.repo.Delete(ctx, blockID)
}