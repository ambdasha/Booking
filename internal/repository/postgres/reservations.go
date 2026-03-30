package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"booking/internal/domain"
	"booking/internal/errs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepo struct {
	db *pgxpool.Pool
}

func NewReservationRepo(db *pgxpool.Pool) *ReservationRepo {
	return &ReservationRepo{db: db}
}

func (r *ReservationRepo) Create(ctx context.Context, res domain.Reservation) (domain.Reservation, error) {
	const q = `
	INSERT INTO reservations (user_id, room_id, start_time, end_time, status)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at;
	`
	err := r.db.QueryRow(ctx, q, res.UserID, res.RoomID, res.StartTime, res.EndTime, res.Status).
		Scan(&res.ID, &res.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23P01": // exclusion_violation (EXCLUDE); попытались создать пересекающуюся бронь
				return domain.Reservation{}, errs.ErrConflict
			case "23503": // foreign_key_violation; например room_id не существует
				return domain.Reservation{}, errs.ErrNotFound
			}
		}
		return domain.Reservation{}, fmt.Errorf("create reservation: %w", err)
	}

	return res, nil
}

func (r *ReservationRepo) GetByID(ctx context.Context, id int64) (domain.Reservation, error) {
	const q = `
	SELECT id, user_id, room_id, start_time, end_time, status, cancelled_at, cancellation_reason, created_at
	FROM reservations
	WHERE id = $1;
	`
	var res domain.Reservation
	err := r.db.QueryRow(ctx, q, id).Scan(
		&res.ID, &res.UserID, &res.RoomID, &res.StartTime, &res.EndTime, &res.Status,
		&res.CancelledAt, &res.CancellationReason, &res.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Reservation{}, errs.ErrNotFound
		}
		return domain.Reservation{}, fmt.Errorf("get reservation: %w", err)
	}
	return res, nil
}

func (r *ReservationRepo) Cancel(ctx context.Context, id int64, reason *string) error {
	const q = `
	UPDATE reservations
	SET status = 'cancelled',
		cancelled_at = $2,
		cancellation_reason = $3
	WHERE id = $1;
	`
	now := time.Now().UTC()
	ct, err := r.db.Exec(ctx, q, id, now, reason)
	if err != nil {
		return fmt.Errorf("cancel reservation: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (r *ReservationRepo) ListByUser(ctx context.Context, userID int64, status string) ([]domain.Reservation, error) {
	const q = `
	SELECT id, user_id, room_id, start_time, end_time, status, cancelled_at, cancellation_reason, created_at
	FROM reservations
	WHERE user_id = $1
	AND ($2 = '' OR status = $2)
	ORDER BY start_time DESC;
	`
	rows, err := r.db.Query(ctx, q, userID, status)
	if err != nil {
		return nil, fmt.Errorf("list reservations: %w", err)
	}
	defer rows.Close()

	var out []domain.Reservation
	for rows.Next() {
		var res domain.Reservation
		if err := rows.Scan(
			&res.ID, &res.UserID, &res.RoomID, &res.StartTime, &res.EndTime, &res.Status,
			&res.CancelledAt, &res.CancellationReason, &res.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan reservation: %w", err)
		}
		out = append(out, res)
	}
	return out, nil
}