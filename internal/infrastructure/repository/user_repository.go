package repository

import (
	"context"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	DeleteByID(ctx context.Context, id string) error
}
