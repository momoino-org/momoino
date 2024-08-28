package core

import (
	"slices"

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

	// GetJWTConfig retrieves the JWT configuration from the application's configuration source.
	// It returns a pointer to a JWTConfig struct containing the JWT configuration details.
	GetJWTConfig() *JWTConfig
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

// JWTConfig holds the configuration details for JSON Web Tokens (JWT).
type JWTConfig struct {
	PublicKey             []byte
	PrivateKey            []byte
	AccessTokenExpiresIn  int64
	RefreshTokenExpiresIn int64
}

// AppConfig is a struct that holds the application's configuration.
type appConfig struct {
	configSource *viper.Viper
}

const (
	// TestingMode: Indicates that the application is running in testing mode.
	TestingMode = "testing"

	// DevelopmentMode: Indicates that the application is running in development mode.
	DevelopmentMode = "development"

	// ProductionMode: Indicates that the application is running in production mode.
	ProductionMode = "production"
)

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
	appMode := appCfg.configSource.GetString("mode")

	if slices.Contains([]string{TestingMode, DevelopmentMode, ProductionMode}, appMode) {
		return appMode
	}

	return DevelopmentMode
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
	databaseCfg := &DatabaseConfig{
		Host:         appCfg.configSource.GetString("database_host"),
		Port:         appCfg.configSource.GetUint16("database_port"),
		DatabaseName: appCfg.configSource.GetString("database_name"),
		Username:     appCfg.configSource.GetString("database_username"),
		Password:     appCfg.configSource.GetString("database_password"),
		MaxAttempts:  appCfg.configSource.GetUint("database_max_attempts"),
	}

	return databaseCfg
}

func (appCfg *appConfig) GetJWTConfig() *JWTConfig {
	return &JWTConfig{
		PublicKey:             []byte(appCfg.configSource.GetString("jwt_public_key")),
		PrivateKey:            []byte(appCfg.configSource.GetString("jwt_private_key")),
		AccessTokenExpiresIn:  appCfg.configSource.GetInt64("jwt_access_token_expires_in"),
		RefreshTokenExpiresIn: appCfg.configSource.GetInt64("jwt_refresh_token_expires_in"),
	}
}

// NewConfigModule returns an fx.Option that provides a new AppConfig instance.
// This option can be used to initialize the application's configuration module.
func NewConfigModule() fx.Option {
	return fx.Module(
		"Config Module",
		fx.Provide(func() AppConfig {
			viperInstance := viper.New()
			viperInstance.SetEnvPrefix("app")
			// viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viperInstance.AutomaticEnv()

			return &appConfig{
				configSource: viperInstance,
			}
		}),
	)
}
