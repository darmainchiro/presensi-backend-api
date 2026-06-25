package service

import (
	"context"
	"errors"
	"backend-api/repository"
)

type UserService interface {
	ValidateUserExists(ctx context.Context, id int64) error
	GetUserDetails(ctx context.Context, id int64) (*repository.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) ValidateUserExists(ctx context.Context, id int64) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("terjadi kesalahan saat mengecek data pengguna")
	}
	if user == nil {
		return errors.New("pengguna tidak terdaftar dalam sistem")
	}
	return nil
}

func (s *userService) GetUserDetails(ctx context.Context, id int64) (*repository.User, error){
	return s.repo.GetByID(ctx, id)
}