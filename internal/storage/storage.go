package storage

import "errors"

var (
	ErrUserAlreadyExists = errors.New("пользователь уже существует")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrServiceNotFound   = errors.New("сервис не найден")
)
