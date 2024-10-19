package core

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type AppConfig interface {
	// GetAppVersion retrieves the application's version from the configuration.
	// It returns the version as a string.
	GetAppVersion() string

	// GetCompatibleVersion retrieves the version of the database schema that the application is compatible with.
	// It returns the compatible version as a string.
	GetCompatibleVersion() string

	// GetMode retrieves the application's mode (development or production) from the configuration.
	// If the mode is not explicitly set, it defaults to development mode.
	GetMode() string

	// GetRevision retrieves the application's revision from the configuration.
	// It returns the revision as a string.
	GetRevision() string

	// IsDevelopment checks if the application is running in development mode.
	// It returns true if the mode is set to development, and false otherwise.
	IsDevelopment() bool

	// IsProduction checks if the application is running in production mode.
	// It returns true if the mode is set to production, and false otherwise.
	IsProduction() bool

	// IsTesting checks if the application is running in testing mode.
	// It returns true if the mode is set to testing, and false otherwise.
	IsTesting() bool

	// GetDatabaseConfig retrieves the database configuration from the application's configuration source.
	// It returns a pointer to a DatabaseConfig struct containing the database configuration details.
	GetDatabaseConfig() *DatabaseConfig

	// getCorsConfig retrieves the CORS (Cross-Origin Resource Sharing) configuration details from the
	// provided viper configuration. It retrieves the CORS configuration details from environment variables
	// using the Viper library. If any of the required environment variables are not set, default values are used.
	GetCorsConfig() *CorsConfig

	// GetSecretKey retrieves the secret key from "secret.key" file.
	GetSecretKey() []byte

	GetKeycloakProvider() *KeycloakProvider

	IsHTTPS() bool

	GetHost() string
}

type KeycloakProvider struct {
	Issuer       string
	WellKnownURL string
	JwksURI      string
	ClientID     string
	ClientSecret string
}

// DatabaseConfig is a struct that holds the database configuration details.
type DatabaseConfig struct {
	Host         string
	Port         uint16
	DatabaseName string
	Username     string
	Password     string
	MaxAttempts  uint
}

// CorsConfig holds the configuration settings for Cross-Origin Resource Sharing (CORS).
// Each field in this struct is populated from corresponding environment variables.
type CorsConfig struct {
	// AllowedOrigins: Specifies the list of allowed origins for CORS,
	// sourced from the environment variable "APP_CORS_ALLOWED_ORIGINS".
	//
	// Default value: ["*"]
	AllowedOrigins []string

	// AllowedMethods: Specifies the HTTP methods allowed for CORS,
	// sourced from the environment variable "APP_CORS_ALLOWED_METHODS".
	//
	// Default value: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
	AllowedMethods []string

	// AllowedHeaders: Specifies the headers allowed in CORS requests,
	// sourced from the environment variable "APP_CORS_ALLOWED_HEADERS".
	//
	// Default value: []
	AllowedHeaders []string

	// ExposedHeaders: Specifies the headers exposed to the browser in CORS responses,
	// sourced from the environment variable "APP_CORS_EXPOSED_HEADERS".
	//
	// Default value: []
	ExposedHeaders []string

	// AllowCredentials: Indicates whether credentials are allowed in CORS requests,
	// sourced from the environment variable "APP_CORS_ALLOW_CREDENTIALS".
	//
	// Default value: false
	AllowCredentials bool

	// MaxAge: Specifies the maximum age (in seconds) for which CORS preflight requests
	// can be cached, sourced from the environment variable "APP_CORS_MAX_AGE".
	// Default value: 0
	MaxAge int
}

// appConfig is a struct that holds the application's configuration.
type appConfig struct {
	appMode          string
	host             *url.URL
	databaseConfig   *DatabaseConfig
	corsConfig       *CorsConfig
	secretKey        []byte
	keycloakProvider *KeycloakProvider
}

const (
	// TestingMode: Indicates that the application is running in testing mode.
	TestingMode = "testing"

	// DevelopmentMode: Indicates that the application is running in development mode.
	DevelopmentMode = "development"

	// ProductionMode: Indicates that the application is running in production mode.
	ProductionMode = "production"
)

var _ AppConfig = (*appConfig)(nil)

var (
	// AppVersion is the current version of the migration application.
	AppVersion string

	// CompatibleVersion is the version of the database schema that the migration application is compatible with.
	CompatibleVersion string

	// AppRevision is the revision of the application.
	AppRevision string
)

func (appCfg *appConfig) GetAppVersion() string {
	return AppVersion
}

func (appCfg *appConfig) GetCompatibleVersion() string {
	return CompatibleVersion
}

func (appCfg *appConfig) GetMode() string {
	return appCfg.appMode
}

func (appCfg *appConfig) IsDevelopment() bool {
	return appCfg.GetMode() == DevelopmentMode
}

func (appCfg *appConfig) IsProduction() bool {
	return appCfg.GetMode() == ProductionMode
}

func (appCfg *appConfig) IsTesting() bool {
	return appCfg.GetMode() == TestingMode
}

func (appCfg *appConfig) GetRevision() string {
	return AppRevision
}

func (appCfg *appConfig) GetDatabaseConfig() *DatabaseConfig {
	return appCfg.databaseConfig
}

func (appCfg *appConfig) GetCorsConfig() *CorsConfig {
	return appCfg.corsConfig
}

func (appCfg *appConfig) GetSecretKey() []byte {
	return appCfg.secretKey
}

func (appCfg *appConfig) IsHTTPS() bool {
	return appCfg.host.Scheme == "https"
}

func (appCfg *appConfig) GetHost() string {
	return appCfg.host.Host
}

func (appCfg *appConfig) GetKeycloakProvider() *KeycloakProvider {
	return appCfg.keycloakProvider
}

// initAppMode retrieves the application mode from the provided viper configuration.
// If the mode is not explicitly set in the configuration, it defaults to development mode.
func initAppMode(v *viper.Viper) string {
	appMode := v.GetString("mode")

	if slices.Contains([]string{TestingMode, DevelopmentMode, ProductionMode}, appMode) {
		return appMode
	}

	return DevelopmentMode
}

// initDatabaseConfig retrieves the database configuration details from the provided viper configuration.
// It returns a pointer to a DatabaseConfig struct containing the database configuration details.
func initDatabaseConfig(v *viper.Viper) *DatabaseConfig {
	v.SetDefault("database_host", "localhost")
	v.SetDefault("database_port", "5432")
	v.SetDefault("database_name", "momoiro-wano")
	v.SetDefault("database_username", "root")
	v.SetDefault("database_password", "Keep!t5ecret")
	v.SetDefault("database_max_attempts", "3")

	databaseCfg := &DatabaseConfig{
		Host:         v.GetString("database_host"),
		Port:         v.GetUint16("database_port"),
		DatabaseName: v.GetString("database_name"),
		Username:     v.GetString("database_username"),
		Password:     v.GetString("database_password"),
		MaxAttempts:  v.GetUint("database_max_attempts"),
	}

	return databaseCfg
}

// initCorsConfig retrieves the CORS (Cross-Origin Resource Sharing) configuration details from the
// provided viper configuration. It retrieves the CORS configuration details from environment variables
// using the Viper library. If any of the required environment variables are not set, default values are used.
func initCorsConfig(v *viper.Viper) *CorsConfig {
	// Set default values for CORS configuration settings if not explicitly set in environment variables
	v.SetDefault("cors_allowed_origins", "*")
	v.SetDefault("cors_allowed_methods", []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	})
	v.SetDefault("cors_allowed_headers", "")
	v.SetDefault("cors_exposed_headers", "")
	v.SetDefault("cors_allow_credentials", "false")
	v.SetDefault("cors_max_age", "0")

	// Retrieve the CORS configuration details from environment variables using the Viper library
	return &CorsConfig{
		AllowedOrigins:   v.GetStringSlice("cors_allowed_origins"),
		AllowedMethods:   v.GetStringSlice("cors_allowed_methods"),
		AllowedHeaders:   v.GetStringSlice("cors_allowed_headers"),
		ExposedHeaders:   v.GetStringSlice("cors_exposed_headers"),
		AllowCredentials: v.GetBool("cors_allow_credentials"),
		MaxAge:           v.GetInt("cors_max_age"),
	}
}

// getSecretKey retrieves and validates the secret key from the provided
// Viper configuration instance. The secret key is expected to be a string
// that is trimmed of any leading or trailing whitespace and must meet the
// required length.
//
// If the length of the secret key does not match the expected AES key length
// (AESSecretKeyLength), the function returns an error.
//
// Parameters:
//   - v: A pointer to a Viper instance from which the secret key is retrieved.
//
// Returns:
//   - A byte slice containing the secret key if valid.
//   - An error if the secret key is of invalid length.
func getSecretKey(v *viper.Viper) ([]byte, error) {
	secretValue := strings.TrimSpace(v.GetString("secret_key"))

	if len(secretValue) != AESSecretKeyLength {
		return nil, errors.New("secret key length is invalid")
	}

	return []byte(secretValue), nil
}

func initHost(v *viper.Viper) (*url.URL, error) {
	u, err := url.Parse(v.GetString("host"))
	if err != nil {
		return nil, fmt.Errorf("cannot parse APP_HOST: %w", err)
	}

	return u, nil
}

func initKeycloakProider(v *viper.Viper, httpClient *resty.Client) (*KeycloakProvider, error) {
	type WellKnownResponse struct {
		Issuer  string `json:"issuer"`
		JwksURI string `json:"jwks_uri"`
	}

	wellKnownURL := v.GetString("keycloak_well_known_url")
	wellKnownResponse := &WellKnownResponse{}

	resp, err := httpClient.R().EnableTrace().SetResult(wellKnownResponse).Get(wellKnownURL)
	if err != nil || !resp.IsSuccess() {
		return nil, fmt.Errorf("failed to retrieve well-known response: (status = %v) %w", resp.Status(), err)
	}

	return &KeycloakProvider{
		Issuer:       wellKnownResponse.Issuer,
		WellKnownURL: wellKnownURL,
		JwksURI:      wellKnownResponse.JwksURI,
		ClientID:     v.GetString("keycloak_client_id"),
		ClientSecret: v.GetString("keycloak_client_secret"),
	}, nil
}

// NewAppConfig initializes and returns a new instance of the application's configuration.
// The configuration is loaded from environment variables and files using the Viper library.
// The function retrieves the application's mode, database connection details, and JWT configuration.
func NewAppConfig(httpClient *resty.Client) (*appConfig, error) {
	// Initialize a new Viper instance
	viperInstance := viper.New()

	// Set the environment variable prefix for Viper
	viperInstance.SetEnvPrefix("app")

	// Enable automatic environment variable loading
	viperInstance.AutomaticEnv()

	host, err := initHost(viperInstance)
	if err != nil {
		return nil, err
	}

	secretKey, err := getSecretKey(viperInstance)
	if err != nil {
		return nil, err
	}

	keycloakProvider, err := initKeycloakProider(viperInstance, httpClient)
	if err != nil {
		return nil, err
	}

	// Create a new appConfig instance with the retrieved configuration details
	return &appConfig{
		secretKey:        secretKey,
		host:             host,
		appMode:          initAppMode(viperInstance),
		databaseConfig:   initDatabaseConfig(viperInstance),
		corsConfig:       initCorsConfig(viperInstance),
		keycloakProvider: keycloakProvider,
	}, nil
}

// NewConfigModule is an Fx option that provides an instance of the application's configuration.
func NewConfigModule() fx.Option {
	return fx.Module(
		"Config Module",
		fx.Provide(
			fx.Annotate(NewAppConfig, fx.As(new(AppConfig))),
		),
	)
}
