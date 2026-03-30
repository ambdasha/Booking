package domain

import "time"

//константы для статусов бронирования
const (
	ReservationPending   = "pending" //бронь создана, но не подтверждена
	ReservationConfirmed = "confirmed" // бронь активна
	ReservationCancelled = "cancelled" //бронь отменили
	ReservationExpired   = "expired" // бронь истекла
)

type Reservation struct {
	ID                 int64
	UserID             int64
	RoomID             int64
	StartTime          time.Time
	EndTime            time.Time
	Status             string
	CancelledAt        *time.Time
	CancellationReason *string
	CreatedAt          time.Time
}
