package cache

import (
	"context"
	"time"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/repository"
	"github.com/redis/go-redis/v9"
)

type RoomCh struct {
	client *redis.Client
	expire time.Duration
}

func NewRoomCh(client *redis.Client, expire time.Duration) repository.RoomRepository {
	return &RoomCh{client: client, expire: expire}
}

func (r *RoomCh) Create(ctx context.Context, room *entity.Room) error {
	return r.client.Set(ctx, room.ID, room, r.expire).Err()
}
func (r *RoomCh) FindByID(ctx context.Context, roomID string) (*entity.Room, error) {
	var room entity.Room
	err := r.client.Get(ctx, roomID).Scan(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}
func (r *RoomCh) DeleteByID(ctx context.Context, roomID string) error {
	return r.client.Del(ctx, roomID).Err()
}
