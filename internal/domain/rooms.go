package domain

import "time"

//то с чем работает бизнес-логика
type Room struct {
	ID          int64
	Name        string
	Description string
	Capacity    int
	Location    string
	IsActive    bool
	CreatedAt   time.Time
}