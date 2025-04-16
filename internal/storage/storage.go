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
	GetURL(ctx context.Context, alias string) (string, error)
	SaveURL(ctx context.Context, url string, alias string) error
	DeleteURL(ctx context.Context, alias string) error
	Ping(ctx context.Context, errShutDown chan error)
	Close()
}
