package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/RomanV1/go-rest-api/internal/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service interface {
	GetUserByID(ctx context.Context, uuid uuid.UUID) (User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]User, error)
	CreateUser(ctx context.Context, dto CreateUserDto) (User, error)
	UpdateUser(ctx context.Context, uuid uuid.UUID, dto UpdateUserDto) (User, error)
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
}

type service struct {
	repo   Repository
	logger *logrus.Logger
}

func NewService(repo Repository, logger *logrus.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s service) GetUserByID(ctx context.Context, uuid uuid.UUID) (User, error) {
	user, err := s.repo.GetOne(ctx, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, sql.ErrNoRows
		}
		return User{}, err
	}

	return user, nil
}

func (s service) GetAllUsers(ctx context.Context, limit, offset int) ([]User, error) {
	users, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s service) CreateUser(ctx context.Context, dto CreateUserDto) (User, error) {
	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return User{}, err
	}

	user, err := s.repo.Create(ctx, CreateUserDto{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hash,
	})
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s service) UpdateUser(ctx context.Context, uuid uuid.UUID, dto UpdateUserDto) (User, error) {
	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return User{}, err
	}

	user, err := s.repo.Update(ctx, uuid, UpdateUserDto{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hash,
	})
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s service) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	return s.repo.Delete(ctx, uuid)
}
