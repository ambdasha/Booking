package dto

import "time"


type Interval struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type AvailabilityResponse struct {
	RoomID int64     `json:"room_id"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`

	Reservations []Interval `json:"reservations"` //занятые активными бронями
	Blocks       []Interval `json:"blocks"` //занятые админскими блокировками
}