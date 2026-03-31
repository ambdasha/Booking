package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"booking/internal/domain"
	"booking/internal/errs"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlockRepo struct {
	db *pgxpool.Pool
}

func NewBlockRepo(db *pgxpool.Pool) *BlockRepo {
	return &BlockRepo{db: db}
}

func (r *BlockRepo) Create(ctx context.Context, b domain.RoomBlock) (domain.RoomBlock, error) {
	const q = `
	INSERT INTO room_blocks (room_id, start_time, end_time, reason, created_by)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at;
	`
	err := r.db.QueryRow(ctx, q, b.RoomID, b.StartTime, b.EndTime, b.Reason, b.CreatedBy).
		Scan(&b.ID, &b.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P01" { // exclusion_violation; если новый блок пересекается с уже существующим блоком той же комнаты.
			return domain.RoomBlock{}, errs.ErrConflict
		}
		return domain.RoomBlock{}, fmt.Errorf("create block: %w", err)
	}
	return b, nil
}

//получить блокировки комнаты в диапазоне времени
func (r *BlockRepo) ListByRoomAndRange(ctx context.Context, roomID int64, from, to time.Time) ([]domain.RoomBlock, error) {
	const q = `
	SELECT id, room_id, start_time, end_time, reason, created_by, created_at
	FROM room_blocks
	WHERE room_id = $1 AND start_time < $3 AND end_time > $2
	ORDER BY start_time;
	`
	rows, err := r.db.Query(ctx, q, roomID, from, to)
	if err != nil {
		return nil, fmt.Errorf("list blocks: %w", err)
	}
	defer rows.Close()

	var out []domain.RoomBlock
	
	for rows.Next() {
		var b domain.RoomBlock
		if err := rows.Scan(&b.ID, &b.RoomID, &b.StartTime, &b.EndTime, &b.Reason, &b.CreatedBy, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan block: %w", err)
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *BlockRepo) Delete(ctx context.Context, blockID int64) error {
	const q = `DELETE FROM room_blocks WHERE id = $1;`
	ct, err := r.db.Exec(ctx, q, blockID)
	if err != nil {
		return fmt.Errorf("delete block: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ErrNotFound
	}
	return nil
}