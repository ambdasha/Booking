package errs

import "errors"

//набор сигнальных ошибок, чтобы разные слои проекта понимали друг друга без строковых сравнений и без утечек деталей
var (
	ErrNotFound           = errors.New("not_found")
	ErrConflict           = errors.New("conflict")
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrValidation         = errors.New("validation_error")
)