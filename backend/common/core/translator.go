package core

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/fx"
)

func NewUniversalTranslator() *ut.UniversalTranslator {
	en := en.New()

	return ut.New(en, en)
}

func NewTranslationModule() fx.Option {
	return fx.Module(
		"Translation Module",
		fx.Provide(NewUniversalTranslator),
	)
}
