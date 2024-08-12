package usermgt

import (
	"context"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByID(ctx context.Context, db *gorm.DB, userID string) (*UserModel, error)
	FindUserByEmail(ctx context.Context, db *gorm.DB, email string) (*UserModel, error)
}

type userRepository struct{}

type UserRepositoryParams struct {
	fx.In
}

var _ UserRepository = (*userRepository)(nil)

func NewUserRepository(params UserRepositoryParams) *userRepository {
	return &userRepository{}
}

func (u userRepository) FindUserByID(ctx context.Context, db *gorm.DB, userID string) (*UserModel, error) {
	var user UserModel

	if result := db.WithContext(ctx).First(&user, "id = ?", userID); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u userRepository) FindUserByEmail(ctx context.Context, db *gorm.DB, email string) (*UserModel, error) {
	var user UserModel

	if result := db.WithContext(ctx).Where("email", email).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
