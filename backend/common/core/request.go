package core

import (
	"github.com/gorilla/schema"
	"go.uber.org/fx"
)

func NewRequestModule() fx.Option {
	return fx.Module(
		"Request Module",
		fx.Provide(schema.NewDecoder),
	)
}
