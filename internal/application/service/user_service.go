package service

import (
	"context"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type UserService struct {
	Repository repository.UserRepository
	log        *logger.Logger
}

func NewUserService(
	repository repository.UserRepository,
	log *logger.Logger,
) *UserService {
	return &UserService{
		Repository: repository,
		log:        log,
	}
}

func (u *UserService) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	return user, u.Repository.Set(ctx, user)
}

func (u *UserService) Delete(ctx context.Context, user *entity.User) error {
	return u.Repository.DeleteByID(ctx, user.ID)
}
