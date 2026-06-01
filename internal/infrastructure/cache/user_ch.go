package cache

import (
	"context"
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
	return &UserCH{
		client: c,
		expire: expire,
	}
}

func (u *UserCH) Create(ctx context.Context, user *entity.User) error {
	return u.client.Set(ctx, user.ID, user, u.expire).Err()
}
func (u *UserCH) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user *entity.User
	err := u.client.Get(ctx, id).Scan(&user)
	return user, err
}
func (u *UserCH) ExistsByID(ctx context.Context, id string) (bool, error) {
	v, err := u.client.Get(ctx, id).Result()
	if err != nil {
		return false, err
	}
	return v != "", nil
}
func (u *UserCH) DeleteByID(ctx context.Context, id string) error {
	err := u.client.Del(ctx, id).Err()
	if err != nil {
		return err
	}
	return nil
}
