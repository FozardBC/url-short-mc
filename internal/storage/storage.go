package storage

import (
	"context"
	"errors"
)

var (
	ErrAliasNotFound      = errors.New("alias not found")
	ErrAliasAlreadyExists = errors.New("alias already exists")
)

type Storage interface {
	GetURL(alias string) (string, error)
	SaveURL(url string, alias string) error
	DeleteURL(alias string) error
	Ping(ctx context.Context, errShutDown chan error) error
	Close()
}
