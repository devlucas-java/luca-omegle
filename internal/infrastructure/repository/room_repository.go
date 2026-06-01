package repository

import (
	"context"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
)

type RoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	FindByID(ctx context.Context, roomID string) (*entity.Room, error)
	DeleteByID(ctx context.Context, roomID string) error
}
