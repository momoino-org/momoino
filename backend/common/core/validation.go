package core

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func NewValidationModule() fx.Option {
	return fx.Module("Validation Module", fx.Provide(validator.New))
}
