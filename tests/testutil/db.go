package testutil
//помогает чистить БД и менять роль пользователя
import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ResetDB(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
	TRUNCATE TABLE
	room_blocks,
	reservations,
	rooms,
	users
	RESTART IDENTITY CASCADE;
	`)
	return err
}

//находит пользователя по email и меняет ему роль на admin
func PromoteAdmin(ctx context.Context, db *pgxpool.Pool, email string) error {
	_, err := db.Exec(ctx, `UPDATE users SET role='admin' WHERE email=$1`, email)
	return err
}