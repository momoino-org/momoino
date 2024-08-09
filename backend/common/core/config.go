package core

import (
	"slices"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	// developmentMode: Indicates that the application is running in development mode.
	developmentMode = "development"

	// productionMode: Indicates that the application is running in production mode.
	productionMode = "production"
)

var (
	// AppVersion is the current version of the migration application.
	AppVersion string

	// CompatibleVersion is the version of the database schema that the migration application is compatible with.
	CompatibleVersion string

	// AppMode is the mode in which the application is running (ex., "development", "production").
	AppMode string

	// AppRevision is the revision of the application.
	AppRevision string
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

	// GetDatabaseConfig retrieves the database configuration from the application's configuration source.
	// It returns a pointer to a DatabaseConfig struct containing the database configuration details.
	GetDatabaseConfig() *DatabaseConfig
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

// AppConfig is a struct that holds the application's configuration.
type appConfig struct {
	configSource *viper.Viper
}

func (appCfg *appConfig) GetAppVersion() string {
	return AppVersion
}

func (appCfg *appConfig) GetCompatibleVersion() string {
	return CompatibleVersion
}

func (appCfg *appConfig) GetMode() string {
	if slices.Contains([]string{developmentMode, productionMode}, AppMode) {
		return AppMode
	}

	return developmentMode
}

func (appCfg *appConfig) IsDevelopment() bool {
	return AppMode == developmentMode
}

func (appCfg *appConfig) IsProduction() bool {
	return AppMode == productionMode
}

func (appCfg *appConfig) GetRevision() string {
	return AppRevision
}

func (appCfg *appConfig) GetDatabaseConfig() *DatabaseConfig {
	databaseCfg := &DatabaseConfig{
		Host:         appCfg.configSource.GetString("database.host"),
		Port:         appCfg.configSource.GetUint16("database.port"),
		DatabaseName: appCfg.configSource.GetString("database.name"),
		Username:     appCfg.configSource.GetString("database.username"),
		Password:     appCfg.configSource.GetString("database.password"),
		MaxAttempts:  appCfg.configSource.GetUint("database.max-attempts"),
	}

	return databaseCfg
}

// NewConfigModule returns an fx.Option that provides a new AppConfig instance.
// This option can be used to initialize the application's configuration module.
func NewConfigModule() fx.Option {
	return fx.Module(
		"Config Module",
		fx.Provide(func() AppConfig {
			viperInstance := viper.New()
			viperInstance.SetEnvPrefix("app")
			viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viperInstance.AutomaticEnv()

			return &appConfig{
				configSource: viperInstance,
			}
		}),
	)
}
