package cache

import (
	"AwesomeProject/internal/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const UserExpDate = time.Hour

type UserStore struct {
	rbd *redis.Client
}

func (s *UserStore) Get(ctx context.Context, id int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", id)
	data, err := s.rbd.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
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
	if user.ID != id {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rbd.Set(ctx, cacheKey, string(data), UserExpDate).Err()
}
