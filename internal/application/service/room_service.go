package service

import (
	"context"
	"fmt"

	"github.com/devlucas-java/luca-omegle/internal/domain/entity"
	"github.com/devlucas-java/luca-omegle/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type RoomService struct {
	roomRepo repository.RoomRepository
	userRepo repository.UserRepository
	log      *logger.Logger
}

func NewRoomService(
	roomRepo repository.RoomRepository,
	userRepo repository.UserRepository,
	log *logger.Logger,
) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
		userRepo: userRepo,
		log:      log,
	}
}

func (s *RoomService) JoinRandom(ctx context.Context, userIDs []string) (*entity.Room, error) {
	if len(userIDs) != entity.MaxUsersPerRoom {
		return nil, fmt.Errorf("expected %d users, got %d", entity.MaxUsersPerRoom, len(userIDs))
	}

	users := make([]*entity.User, 0, len(userIDs))
	for _, id := range userIDs {
		user, err := s.findUser(ctx, id)
		if err != nil {
			return nil, err
		}

		if user.RoomID != "" {
			if err = s.LeaveRoom(ctx, user); err != nil {
				return nil, err
			}
		}

		users = append(users, user)
	}

	room, err := s.CreateRoom(ctx, users)
	if err != nil {
		return nil, err
	}

	s.log.Info(fmt.Sprintf("room %s created for users %v", room.ID, userIDs))
	return room, nil
}

func (s *RoomService) LeaveRoom(ctx context.Context, user *entity.User) error {
	room, err := s.roomRepo.FindByID(ctx, user.RoomID)
	if err != nil {
		return err
	}

	room.RemoveUser(user.ID)

	if room.IsEmpty() {
		if err = s.roomRepo.DeleteByID(ctx, room.ID); err != nil {
			return err
		}
	} else {
		if err = s.roomRepo.Set(ctx, room); err != nil {
			return err
		}
	}

	user.ClearRoom()
	if err = s.userRepo.Set(ctx, user); err != nil {
		return err
	}

	s.log.Info(fmt.Sprintf("user %s left room %s", user.ID, room.ID))
	return nil
}

func (s *RoomService) CreateRoom(ctx context.Context, users []*entity.User) (*entity.Room, error) {
	room := entity.NewRoom()

	for _, user := range users {
		if room.IsFull() {
			return nil, fmt.Errorf("room capacity exceeded (max %d)", entity.MaxUsersPerRoom)
		}

		room.AddUser(user.ID)
		user.AssignRoom(room.ID)

		if err := s.userRepo.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	if err := s.roomRepo.Set(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) findUser(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
