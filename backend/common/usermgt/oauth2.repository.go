package usermgt

import (
	"context"
	"database/sql"
	"net/http"
	"wano-island/common/core"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// OAuth2ProviderRepository defines the contract for managing OAuth2 providers within the system.
// It provides methods for retrieving and creating OAuth2 provider records in the database.
type OAuth2ProviderRepository interface {
	// GetMany retrieves multiple OAuth2 provider models based on the incoming HTTP request.
	//
	// This function is expected to process the request, potentially applying filtering,
	// sorting, or pagination based on the request parameters.
	//
	// Parameters:
	//   - r (*http.Request): The HTTP request containing any relevant query parameters
	//     for filtering, sorting, or pagination.
	//
	// Returns:
	//   - []OAuth2ProviderModel: A slice containing the retrieved OAuth2 provider models.
	//   - *int64: A pointer to the total count of available records matching the request criteria
	//     (useful for pagination purposes).
	//   - error: An error object if the retrieval fails; otherwise, nil.
	GetMany(r *http.Request) ([]OAuth2ProviderModel, *int64, error)

	// Get retrieves an OAuth2 provider by its name.
	//
	// Params:
	//   - ctx: The context for the operation, which may include cancellation or deadlines.
	//   - name: The name of the OAuth2 provider to retrieve.
	//
	// Returns:
	//   - A pointer to the OAuth2ProviderModel if found.
	//   - An error if the provider is not found or if the operation fails.
	Get(ctx context.Context, name string) (*OAuth2ProviderModel, error)

	// Create inserts a new OAuth2 provider into the system.
	//
	// Params:
	//   - ctx: The context for the operation.
	//   - params: The parameters needed to create a new OAuth2 provider.
	//
	// Returns:
	//   - A pointer to the newly created OAuth2ProviderModel.
	//   - An error if any issues occur during the creation process.
	Create(ctx context.Context, params CreateOAuth2ProviderParams) (*OAuth2ProviderModel, error)
}

// CreateOAuth2ProviderParams holds the input parameters required to create an OAuth2 provider.
type CreateOAuth2ProviderParams struct {
	// The name of the OAuth2 provider (e.g., Google, Facebook).
	Provider string

	// The client ID associated with the OAuth2 provider.
	ClientID string

	// The client secret associated with the OAuth2 provider.
	ClientSecret string

	// The URL to redirect users to after they authenticate.
	RedirectURL string

	// The OAuth2 scopes required by the provider for access (e.g., profile, email).
	Scopes []string

	// A flag indicating whether the OAuth2 provider is currently active or not.
	IsEnabled bool

	// The user or system entity responsible for creating the provider.
	CreatedBy core.PrincipalUser
}

// oauth2ProviderRepository provides the concrete implementation of the OAuth2ProviderRepository interface.
// It interacts with the database to manage OAuth2 provider records.
type oauth2ProviderRepository struct {
	db              *gorm.DB
	aesGCMEncryptor core.Encryptor
}

// OAuth2ProviderRepositoryParams defines the dependencies required to instantiate an OAuth2 provider repository.
//
// Fields:
//   - DB: The GORM database connection used to interact with the database.
//   - AESGCMEncryptor: The encryptor used to securely store sensitive information, such as client secrets.
type OAuth2ProviderRepositoryParams struct {
	fx.In
	DB              *gorm.DB
	AESGCMEncryptor core.Encryptor `name:"aes-gcm"`
}

// Ensure that *oauth2ProviderRepository implements the OAuth2ProviderRepository interface.
// This compile-time check enforces the contract between the interface and its implementation.
var _ (OAuth2ProviderRepository) = (*oauth2ProviderRepository)(nil)

// NewOAuth2ProviderRepository constructs and returns a new instance of oauth2ProviderRepository.
//
// Params:
//   - p: The dependencies required to initialize the repository, injected via the Fx framework.
//
// Returns:
//   - A new oauth2ProviderRepository instance configured with the provided database and encryptor.
func NewOAuth2ProviderRepository(params OAuth2ProviderRepositoryParams) *oauth2ProviderRepository {
	return &oauth2ProviderRepository{
		db:              params.DB,
		aesGCMEncryptor: params.AESGCMEncryptor,
	}
}

func (r *oauth2ProviderRepository) GetMany(request *http.Request) ([]OAuth2ProviderModel, *int64, error) {
	var (
		totalRows      int64
		oauth2Provider []OAuth2ProviderModel
	)

	err := r.db.WithContext(request.Context()).Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(&OAuth2ProviderModel{}).
			Count(&totalRows); result.Error != nil {
			return result.Error
		}

		if result := tx.Omit("ClientID", "ClientSecret", "RedirectURL", "Scopes").
			Scopes(core.Paginate(request)).
			Find(&oauth2Provider); result.Error != nil {
			return result.Error
		}

		return nil
	}, &sql.TxOptions{ReadOnly: true})

	if err != nil {
		return nil, nil, err
	}

	return oauth2Provider, &totalRows, nil
}

func (r *oauth2ProviderRepository) Get(ctx context.Context, name string) (*OAuth2ProviderModel, error) {
	var oauth2Provider OAuth2ProviderModel

	result := r.db.WithContext(ctx).
		Where(&OAuth2ProviderModel{Provider: name}).
		First(&oauth2Provider)

	if result.Error != nil {
		return nil, result.Error
	}

	return &oauth2Provider, nil
}

func (r *oauth2ProviderRepository) Create(
	ctx context.Context,
	params CreateOAuth2ProviderParams,
) (*OAuth2ProviderModel, error) {
	oauth2Provider := OAuth2ProviderModel{
		Provider:     params.Provider,
		ClientID:     params.ClientID,
		ClientSecret: params.ClientSecret,
		RedirectURL:  params.RedirectURL,
		Scopes:       params.Scopes,
		IsEnabled:    params.IsEnabled,
		HasCreatedByColumn: core.HasCreatedByColumn{
			CreatedBy: params.CreatedBy.GetUsername(),
		},
	}

	if result := r.db.WithContext(ctx).Create(&oauth2Provider); result.Error != nil {
		return nil, result.Error
	}

	return &oauth2Provider, nil
}
