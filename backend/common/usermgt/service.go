package usermgt

import (
	"context"
	"errors"
	"log/slog"
	"wano-island/common/core"

	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ComparePassword(ctx context.Context, password []byte, hasedPassword []byte) error
	HashPassword(ctx context.Context, password string) (*[]byte, error)
}

type userService struct {
	logger core.Logger
}

type UserServiceParams struct {
	fx.In
	Logger core.Logger
}

func NewUserService(params UserServiceParams) *userService {
	return &userService{
		logger: params.Logger,
	}
}

func (s *userService) ComparePassword(ctx context.Context, password []byte, hasedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hasedPassword, password)

	if err == nil {
		return nil
	}

	if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		s.logger.ErrorContext(ctx, "Something went wrong when comparing password", slog.Any("details", err))
		return err
	}

	return err
}

func (s *userService) HashPassword(ctx context.Context, password string) (*[]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		s.logger.ErrorContext(ctx, "Something went wrong when hashing password", slog.Any("details", err))
		return nil, err
	}

	return &hashedPassword, nil
}
