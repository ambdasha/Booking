//правила создания комнат,
package service

import (
	"context"

	"booking/internal/domain"
	"booking/internal/dto"
)

type RoomRepository interface {
	ListActive(ctx context.Context) ([]domain.Room, error)
	GetByID(ctx context.Context, id int64) (domain.Room, error)
	Create(ctx context.Context, room domain.Room) (domain.Room, error)
	Update(ctx context.Context, room domain.Room) (domain.Room, error)
	Deactivate(ctx context.Context, id int64) error
}


type RoomService struct {
	repo RoomRepository
}

func NewRoomService(repo RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

//получение списка комнат
func (s *RoomService) List(ctx context.Context) ([]domain.Room, error) {
	return s.repo.ListActive(ctx)
}

//получение комнаты по id
func (s *RoomService) Get(ctx context.Context, id int64) (domain.Room, error) {
	return s.repo.GetByID(ctx, id)
}

//создание комнаты; преобразует из dto в domain room и сохраняет в бд
func (s *RoomService) Create(ctx context.Context, req dto.CreateRoomRequest) (domain.Room, error) {
	room := domain.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Location:    req.Location,
	}
	return s.repo.Create(ctx, room)
}

//обновление комнаты
func (s *RoomService) Update(ctx context.Context, id int64, req dto.UpdateRoomRequest) (domain.Room, error) {
	room := domain.Room{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Location:    req.Location,
	}
	return s.repo.Update(ctx, room)
}

//дизактивирует комнату, при необходимости комнату модно булет вернуть
func (s *RoomService) Deactivate(ctx context.Context, id int64) error {
	return s.repo.Deactivate(ctx, id)
}