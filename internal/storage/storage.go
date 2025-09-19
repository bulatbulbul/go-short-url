package storage

import "errors"

var (
	ErrURLFound  = errors.New("url не найден")
	ErrURLExists = errors.New("url существует")
)
