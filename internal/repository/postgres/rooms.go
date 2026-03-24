//функции, которые умеют читать b писать комнаты в Postgres

package postgres

import (
	"context"
	"fmt"

	"booking/internal/domain"
	"booking/internal/errs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepo struct {
	db *pgxpool.Pool
}

func NewRoomRepo(db *pgxpool.Pool) *RoomRepo {
	return &RoomRepo{db: db}
}

func (r *RoomRepo) ListActive(ctx context.Context) ([]domain.Room, error) {
	const q = `
		SELECT id, name, description, capacity, location, is_active, created_at
		FROM rooms
		WHERE is_active = true
		ORDER BY id;
		`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list rooms: %w", err)
	}
	defer rows.Close()

	var res []domain.Room
	for rows.Next() {
		var room domain.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.Location, &room.IsActive, &room.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan room: %w", err)
		}
		res = append(res, room)
	}
	return res, nil
}

func (r *RoomRepo) GetByID(ctx context.Context, id int64) (domain.Room, error) {
	const q = `
		SELECT id, name, description, capacity, location, is_active, created_at
		FROM rooms
		WHERE id = $1;
		`
	var room domain.Room
	err := r.db.QueryRow(ctx, q, id).Scan(
		&room.ID, &room.Name, &room.Description, &room.Capacity, &room.Location, &room.IsActive, &room.CreatedAt,
	)
	if err != nil {
		return domain.Room{}, errs.ErrNotFound
	}
	return room, nil
}

func (r *RoomRepo) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	const q = `
		INSERT INTO rooms (name, description, capacity, location, is_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id, created_at;
		`
	err := r.db.QueryRow(ctx, q, room.Name, room.Description, room.Capacity, room.Location).
		Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		return domain.Room{}, fmt.Errorf("create room: %w", err)
	}
	room.IsActive = true
	return room, nil
}

func (r *RoomRepo) Update(ctx context.Context, room domain.Room) (domain.Room, error) {
	const q = `
		UPDATE rooms
		SET name = $2, description = $3, capacity = $4, location = $5
		WHERE id = $1
		RETURNING id, name, description, capacity, location, is_active, created_at;
		`
	var out domain.Room
	err := r.db.QueryRow(ctx, q, room.ID, room.Name, room.Description, room.Capacity, room.Location).
		Scan(&out.ID, &out.Name, &out.Description, &out.Capacity, &out.Location, &out.IsActive, &out.CreatedAt)
	if err != nil {
		return domain.Room{}, errs.ErrNotFound
	}
	return out, nil
}

func (r *RoomRepo) Deactivate(ctx context.Context, id int64) error {
	const q = `
		UPDATE rooms SET is_active = false WHERE id = $1;
		`
	ct, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("deactivate room: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}