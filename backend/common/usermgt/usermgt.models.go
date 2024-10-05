package usermgt

import (
	"wano-island/common/core"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// UserModel represents a user in the system.
type UserModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasCreatedByColumn
	core.HasUpdatedAtColumn

	// The unique username of the user.
	Username string `gorm:"type:string;size:256;not null;unique"`

	// The user's unique email address.
	Email string `gorm:"type:string;size:256;not null;unique"`

	// A flag indicating whether the user's email has been verified.
	VerifiedEmail bool `gorm:"type:bool;not null"`

	// The hashed password of the user. This field is nullable to support OAuth2 users.
	Password *string `gorm:"type:string;size:256;"`

	// The user's first name
	FirstName string `gorm:"type:string;size:64"`

	// The user's last name
	LastName string `gorm:"type:string;size:64"`

	// The user's preferred language or locale (e.g., "en" for English). Defaults to "en" if not specified.
	Locale string `gorm:"type:string;not null;default:en"`

	// A list of OAuth2 providers linked to this user.
	LinkedProviders []OAuth2UserModel `gorm:"foreignKey:LocalID"`
}

// OAuth2ProviderModel represents an OAuth2 provider configuration in the system.
// It stores the necessary credentials and settings for authenticating with an external OAuth2 provider.
type OAuth2ProviderModel struct {
	core.Model
	core.HasCreatedAtColumn
	core.HasCreatedByColumn

	// The name of the OAuth2 provider (e.g., "google", "facebook").
	Provider string `gorm:"type:string;size:256;not null;unique"`

	// The client ID used to authenticate with the provider, encrypted for security.
	ClientID string `gorm:"type:text;not null;serializer:encryption"`

	// The client secret used to authenticate with the provider, encrypted for security.
	ClientSecret string `gorm:"type:text;not null;serializer:encryption"`

	// The URL where the provider will redirect after authentication.
	RedirectURL string `gorm:"type:string;size:256;not null"`

	// The permissions requested from the provider (e.g., "email", "profile").
	Scopes pq.StringArray `gorm:"type:text[];not null"`

	// A flag indicating whether the OAuth2 provider is enabled.
	IsEnabled bool `gorm:"type:boolean;not null"`

	// A list of users associated with this provider.
	Users []OAuth2UserModel `gorm:"foreignKey:ProviderID"`
}

// OAuth2UserModel represents the association between a user in the system and an OAuth2 provider.
// It links the user's internal ID with the external provider's unique identifier for the user (OpenID).
type OAuth2UserModel struct {
	// The unique ID of the OAuth2 provider.
	ProviderID uuid.UUID `gorm:"primaryKey;not null;type:uuid"`

	// The unique ID of the local user in the system.
	LocalID uuid.UUID `gorm:"primaryKey;not null;type:uuid"`

	// The unique identifier for the user in the external OAuth2 provider.
	OpenID string `gorm:"primaryKey;type:string;size:256;not null"`
}

func (UserModel) TableName() string {
	return "public.users"
}

func (OAuth2ProviderModel) TableName() string {
	return "public.oauth2_providers"
}

func (OAuth2UserModel) TableName() string {
	return "public.oauth2_users"
}
