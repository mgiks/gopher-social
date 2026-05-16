package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	store "github.com/mgiks/gopher-social/internal/store/db"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExprTime = time.Minute * 5

func (s UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s UserStore) Set(ctx context.Context, user store.User) error {
	if user.ID == 0 {
		return fmt.Errorf("user ID not provided")
	}

	cacheKey := fmt.Sprintf("user-%d", user.ID)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, jsonData, UserExprTime).Err()
}
