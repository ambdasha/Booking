package dto

import "time"

//структуры для HTTP входа/выхода

type CreateBlockRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Reason    string    `json:"reason" validate:"max=256"`
}

type BlockResponse struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}