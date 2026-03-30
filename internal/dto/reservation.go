package dto

import "time"

type CreateReservationRequest struct {
	RoomID    int64     `json:"room_id" validate:"required,min=1"`
	StartTime time.Time `json:"start_time"` 
	EndTime   time.Time `json:"end_time"`
}

type CancelReservationRequest struct {
	Reason string `json:"reason" validate:"max=256"`
}

type ReservationResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RoomID    int64     `json:"room_id"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
}