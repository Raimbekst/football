package domain

import (
	"errors"
)

var (
	ErrUserAlreadyExist  = errors.New("пользователь уже существует")
	ErrUserDoesNotExist  = errors.New("неверный номер телефона или пароль")
	ErrUserNotRegistered = errors.New("пользователь не зарегистрирован")
	ErrPasswordNotMatch  = errors.New("пароли не совпадают")
	ErrInvalidSecretCode = errors.New("неверный секретный код")
	ErrInvalidPassword   = errors.New("неверный старый пароль")
	ErrNotFound          = errors.New("не найдено")
	ErrUserCommented     = errors.New("пользователь уже коммент оставил")
)
