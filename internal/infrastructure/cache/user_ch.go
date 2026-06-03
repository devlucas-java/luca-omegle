package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/repository"
	"github.com/redis/go-redis/v9"
)

type UserCH struct {
	client *redis.Client
	expire time.Duration
}

func NewUserCH(c *redis.Client, expire time.Duration) repository.UserRepository {
	return &UserCH{client: c, expire: expire}
}

var _ repository.UserRepository = (*UserCH)(nil)

func (u *UserCH) Set(ctx context.Context, user *entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return u.client.Set(ctx, userKey(user.ID), data, u.expire).Err()
}

func (u *UserCH) FindByID(ctx context.Context, id string) (*entity.User, error) {
	data, err := u.client.Get(ctx, userKey(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var user entity.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserCH) ExistsByID(ctx context.Context, id string) (bool, error) {
	n, err := u.client.Exists(ctx, userKey(id)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (u *UserCH) DeleteByID(ctx context.Context, id string) error {
	return u.client.Del(ctx, userKey(id)).Err()
}

func userKey(id string) string { return "user:" + id }
