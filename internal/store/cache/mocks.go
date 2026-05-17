package cache

import (
	"context"

	store "github.com/mgiks/gopher-social/internal/store/db"
)

func NewMockStore() Store {
	return Store{
		Users: MockUserStore{},
	}
}

type MockUserStore struct{}

func (s MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	return nil, nil
}

func (s MockUserStore) Set(ctx context.Context, user store.User) error {
	return nil
}
