package core

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
)

// NewHTTPClient creates and returns a new instance of resty.Client.
// This function is used to create a new HTTP client for making requests to external APIs.
func NewHTTPClient() *resty.Client {
	return resty.New()
}

func NewHTTPClientModule() fx.Option {
	return fx.Module(
		"Http Client Module",
		fx.Provide(NewHTTPClient),
	)
}
