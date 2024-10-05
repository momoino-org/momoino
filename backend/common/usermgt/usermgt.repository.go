package usermgt

import (
	"context"
	"wano-island/common/core"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// UserRepository defines the interface for interacting with user data in the database.
type UserRepository interface {
	// FindUserByID retrieves a user by their unique ID from the database.
	//
	// Params:
	//  - ctx: The context for the operation.
	//  - db: The GORM database connection.
	//  - userID: The unique identifier of the user to retrieve.
	//
	// Returns:
	//  - A pointer to the UserModel if found, nil otherwise.
	//  - An error if any occurred during the operation.
	FindUserByID(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*UserModel, error)

	// FindUserByUsername retrieves a user by their username from the database.
	//
	// Params:
	//  - ctx: The context for the operation.
	//  - db: The GORM database connection.
	//  - username: The username of the user to retrieve.
	//
	// Returns:
	//  - A pointer to the UserModel if found, nil otherwise.
	//  - An error if any occurred during the operation.
	FindUserByUsername(ctx context.Context, db *gorm.DB, username string) (*UserModel, error)

	// FindUserByEmail retrieves a user by their email from the database.
	//
	// Params:
	//  - ctx: The context for the operation.
	//  - db: The GORM database connection.
	//  - email: The email of the user to retrieve.
	//
	// Returns:
	//  - A pointer to the UserModel if found, nil otherwise.
	//  - An error if any occurred during the operation.
	FindUserByEmail(ctx context.Context, db *gorm.DB, email string) (*UserModel, error)

	// ChangePassword updates the password of a user in the database.
	//
	// Params:
	//  - ctx: The context for the operation.
	//  - db: The GORM database connection.
	//  - userID: The unique identifier of the user whose password needs to be updated.
	//  - password: The new password for the user.
	//
	// Returns:
	//  - A pointer to the updated UserModel.
	//  - An error if any occurred during the operation.
	ChangePassword(ctx context.Context, db *gorm.DB, userID string, password string) (*UserModel, error)

	// FirstOrCreateUser creates a new user in the database if one does not already exist.
	// If a user with the same username or email already exists, it retrieves the existing user.
	//
	// Params:
	//  - ctx: The context for the operation.
	//  - db: The GORM database connection.
	//  - params: The parameters required to create or find a user.
	//
	// Returns:
	//  - A pointer to the UserModel.
	//  - An error if any occurred during the operation.
	FirstOrCreateUser(ctx context.Context, db *gorm.DB, params CreateUserParams) (*UserModel, error)
}

// LinkedProvider represents a relationship between a user and an OAuth2 provider.
type LinkedProvider struct {
	// Provider is the OAuth2 provider associated with the user.
	Provider OAuth2ProviderModel

	// OpenID is the unique identifier provided by the OAuth2 provider for the user.
	OpenID string
}

// CreateUserParams defines the parameters required for creating a new user.
// This includes user credentials, personal information, and linked providers.
type CreateUserParams struct {
	// The unique username for the user.
	Username string

	// The user's email address.
	Email string

	// Indicates whether the user's email has been verified.
	VerifiedEmail bool

	// A pointer to the user's password.
	Password *string

	// The user's first name.
	FirstName string

	// The user's last name.
	LastName string

	// The user's locale or preferred language setting.
	Locale string

	// The user or system entity responsible for creating this user.
	CreatedBy core.PrincipalUser

	// A list of linked third-party authentication providers for the user.
	LinkedProviders []LinkedProvider
}

// userRepository is a concrete implementation of the UserRepository interface.
// It provides methods for interacting with the user data in the database.
type userRepository struct{}

// UserRepositoryParams defines the dependencies required for constructing a user repository.
// This struct uses the fx.In tag to indicate that it should be populated by the Fx dependency injection system.
type UserRepositoryParams struct {
	fx.In
}

// _ ensures that *userRepository implements the UserRepository interface.
// This is a compile-time check that the concrete type `userRepository` satisfies the `UserRepository` interface.
// If `userRepository` does not implement all the methods in the `UserRepository` interface,
// the code will fail to compile.
var _ UserRepository = (*userRepository)(nil)

// NewUserRepository creates a new instance of the UserRepository interface.
// It initializes and returns a new userRepository struct.
//
// Params:
//   - params: An instance of UserRepositoryParams, which is used to inject dependencies into the userRepository.
//
// Returns:
//   - A pointer to a new userRepository instance.
func NewUserRepository(params UserRepositoryParams) *userRepository {
	return &userRepository{}
}

func (u userRepository) FindUserByID(ctx context.Context, db *gorm.DB, userID uuid.UUID) (*UserModel, error) {
	var user UserModel

	if result := db.WithContext(ctx).First(&user, "id = ?", userID); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u userRepository) FindUserByUsername(ctx context.Context, db *gorm.DB, username string) (*UserModel, error) {
	var model UserModel

	if result := db.WithContext(ctx).Where("username = ?", username).First(&model); result.Error != nil {
		return nil, result.Error
	}

	return &model, nil
}

func (u userRepository) FindUserByEmail(ctx context.Context, db *gorm.DB, email string) (*UserModel, error) {
	var user UserModel

	if result := db.WithContext(ctx).Where("email", email).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u userRepository) ChangePassword(
	ctx context.Context,
	db *gorm.DB,
	userID string,
	password string,
) (*UserModel, error) {
	user := &UserModel{
		Model: core.Model{
			ID: uuid.MustParse(userID),
		},
	}

	if result := db.Model(&user).Update("password", password); result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u userRepository) FirstOrCreateUser(
	ctx context.Context,
	db *gorm.DB,
	params CreateUserParams,
) (*UserModel, error) {
	var userModel UserModel

	result := db.WithContext(ctx).
		Where("username = ?", params.Username).
		Attrs(UserModel{
			Username:      params.Username,
			Email:         params.Email,
			VerifiedEmail: params.VerifiedEmail,
			FirstName:     params.FirstName,
			LastName:      params.LastName,
			Password:      params.Password,
			Locale:        params.Locale,
			HasCreatedByColumn: core.HasCreatedByColumn{
				CreatedBy: params.CreatedBy.GetUsername(),
			},
			LinkedProviders: lo.Map(params.LinkedProviders, func(linkedProvider LinkedProvider, _ int) OAuth2UserModel {
				return OAuth2UserModel{
					ProviderID: linkedProvider.Provider.ID,
					OpenID:     linkedProvider.OpenID,
				}
			}),
		}).
		FirstOrCreate(&userModel)

	if result.Error != nil {
		return nil, result.Error
	}

	return &userModel, nil
}
