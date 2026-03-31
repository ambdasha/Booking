package domain

import "time"


type RoomBlock struct {
	ID        int64
	RoomID    int64
	StartTime time.Time
	EndTime   time.Time
	Reason    string
	CreatedBy *int64
	CreatedAt time.Time
}