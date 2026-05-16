package cache

import (
	"context"

	store "github.com/mgiks/gopher-social/internal/store/db"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, store.User) error
	}
}

func NewCacheStore(rdb *redis.Client) Store {
	return Store{
		Users: UserStore{rdb: rdb},
	}
}
