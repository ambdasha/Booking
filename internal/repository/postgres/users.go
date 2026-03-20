package postgres

import (
	"context"
	"errors"
	"fmt"

	"booking/internal/domain"
	"booking/internal/errs"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u domain.User) (domain.User, error) {
	const q = `
		INSERT INTO users (email, password_hash, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at;
		`

	err := r.db.QueryRow(ctx, q, u.Email, u.PasswordHash, u.Name, u.Role).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError

		//23505 - SQLSTATE код Postgres для unique_violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, errs.ErrConflict
		}
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	return u, nil

}


func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const q = `
		SELECT id, email, password_hash, name, role, created_at
		FROM users
		WHERE email = $1;
		`

	var u domain.User

	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		return domain.User{}, errs.ErrNotFound
	}
	return u, nil
	
}