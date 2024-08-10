package core

import (
	"context"
	"embed"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type LocalizerCtxKey string

const LocalizerCtxID LocalizerCtxKey = "LocalizerCtxID"

const (
	MsgSuccess                   = "S-0000"
	MsgInvalidEmailOrPassword    = "E-0000"
	MsgFailedToDecodeRequestBody = "E-0001"
	MsgInternalServerError       = "U-0000"
)

func GetLocalizer(r *http.Request) *i18n.Localizer {
	if localizer, ok := r.Context().Value(LocalizerCtxID).(*i18n.Localizer); ok {
		return localizer
	}

	return nil
}

func NewI18nModule(fs embed.FS) fx.Option {
	return fx.Module(
		"I18n Module",
		fx.Provide(func(appLifeCycle fx.Lifecycle) *i18n.Bundle {
			bundle := i18n.NewBundle(language.English)
			bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

			appLifeCycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					_, err := bundle.LoadMessageFileFS(fs, "resources/trans/locale.en.yaml")
					return err
				},
			})

			return bundle
		}),
	)
}
