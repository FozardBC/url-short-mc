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
	URL(alias string) (string, error)
	SaveURL(url string, alias string) error
	RemoveAlias(alias string) error
	Ping(ctx context.Context, errShutDown chan error) error
	Close()
}
