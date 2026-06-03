package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/repository"
	"github.com/redis/go-redis/v9"
)

type RoomCh struct {
	client *redis.Client
	expire time.Duration
	key    string
}

func NewRoomCh(client *redis.Client, expire time.Duration) repository.RoomRepository {
	return &RoomCh{client: client, expire: expire}
}

func (r *RoomCh) Set(ctx context.Context, room *entity.Room) error {
	data, err := json.Marshal(room)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, toRedisID(room.ID), data, r.expire).Err()
	if err != nil {
		return err
	}
	return r.client.SAdd(ctx, r.key, room.ID).Err()
}

func (r *RoomCh) FindByID(ctx context.Context, roomID string) (*entity.Room, error) {
	data, err := r.client.Get(ctx, toRedisID(roomID)).Bytes()
	if err != nil {
		return nil, err
	}
	var room entity.Room
	if err := json.Unmarshal(data, &room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomCh) DeleteByID(ctx context.Context, roomID string) error {
	err := r.client.Del(ctx, toRedisID(roomID)).Err()
	if err != nil {
		return err
	}

	return r.client.SRem(ctx, r.key, roomID).Err()

}

func (r *RoomCh) GetIDs(ctx context.Context) []string {
	result := r.client.SMembers(ctx, r.key)
	return result.Val()
}

func toRedisID(id string) string {
	return fmt.Sprintf("room:%s", id)
}
