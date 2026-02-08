package errors

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")  // дубликат, например email
	ErrBadInput       = errors.New("bad input") // валидация
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidAuth    = errors.New("invalid auth") // неверный логин/пароль или сессия
	ErrSessionExpired = errors.New("session expired")
)
