package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

// URLDeleter интерфейс для удаления URL
type URLDeleter interface {
	DeleteURL(alias string) error
}
