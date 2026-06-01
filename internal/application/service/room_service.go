package service

import (
	"context"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/repository"
)

type RoomService struct {
	RoomRepository repository.RoomRepository
	UserRepository repository.UserRepository
}

func NewRoomService(repository repository.RoomRepository) *RoomService {
	return &RoomService{RoomRepository: repository}
}

func (r *RoomService) JoinRoom(ctx context.Context, user *entity.User) error {

	return nil
}

func (r *RoomService) LeaveRoom(ctx context.Context, user *entity.User) error {

	return nil
}

func (r *RoomService) CreateRoom(ctx context.Context, room *entity.Room) error {
	return r.RoomRepository.Create(ctx, room)
}
