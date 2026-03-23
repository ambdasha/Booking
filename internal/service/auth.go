package service

import (
	"context"
	"time"

	"booking/internal/auth"
	"booking/internal/domain"
	"booking/internal/dto"
	"booking/internal/errs"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type AuthService struct {
	users     UserRepository
	jwtSecret string
}

func NewAuthService(users UserRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret}
}


//последущие методы Register и Login отвечают за бизнес-логику, то есть хеширование паролей и их проверка, выдача jwt и обработка ошибок


func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (domain.User, error) {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return domain.User{}, err
	}

	u := domain.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Role:         "user",
	}
	return s.users.Create(ctx, u)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.TokenResponse, error) {
	u, err := s.users.GetByEmail(ctx, req.Email)

	if err != nil {
		return dto.TokenResponse{}, errs.ErrInvalidCredentials
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		return dto.TokenResponse{}, errs.ErrInvalidCredentials
	}

	token, err := auth.SignJWT(s.jwtSecret, u.ID, u.Role, 2*time.Hour)

	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{AccessToken: token, TokenType: "Bearer"}, nil
}