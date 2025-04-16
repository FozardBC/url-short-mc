package hashmap

import (
	"context"
	"microservice_t/internal/storage"
)

type HashmapStorage struct {
	data map[string]string
}

func New() *HashmapStorage {
	return &HashmapStorage{
		data: map[string]string{"vk": "https://vk.com"},
		data: map[string]string{"vk": "https://vk.com"},
	}
}

func (m *HashmapStorage) GetURL(alias string) (string, error) {
func (m *HashmapStorage) GetURL(alias string) (string, error) {

	url, ok := m.data[alias]

	if !ok {
		return "", storage.ErrAliasNotFound
	}

	return url, nil
}

func (m *HashmapStorage) SaveURL(url string, alias string) error {
	_, ok := m.data[alias]
	if ok {
		return storage.ErrAliasAlreadyExists
	}

	m.data[alias] = url

	return nil
}

func (m *HashmapStorage) DeleteURL(alias string) error {
func (m *HashmapStorage) DeleteURL(alias string) error {
	_, ok := m.data[alias]
	if !ok {
	if !ok {
		return storage.ErrAliasNotFound
	}

	delete(m.data, alias)

	return nil
}

func (m *HashmapStorage) Ping(ctx context.Context, errShutDown chan error) error {
	return nil
}

func (m *HashmapStorage) Close() {

}
