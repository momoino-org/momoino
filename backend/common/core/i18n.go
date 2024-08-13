package core

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type LocalizerCtxKey string

type I18nBundleParams struct {
	fx.In

	AppLifeCycle fx.Lifecycle
	LocaleFS     fs.FS
}

const localizerCtxID LocalizerCtxKey = "LocalizerCtxID"

const (
	MsgSuccess                              = "S-0000"
	MsgInvalidEmailOrPassword               = "E-0000"
	MsgFailedToDecodeRequestBody            = "E-0001"
	MsgNeedToLogin                          = "E-0002"
	MsgCannotProcessYourLogin               = "E-0003"
	MsgPasswordChangeErrorDueToUserNotFound = "E-0004"
	MsgValidationFailed                     = "E-0005"
	MsgInternalServerError                  = "U-0000"
)

func GetLocalizer(r *http.Request) *i18n.Localizer {
	if localizer, ok := r.Context().Value(localizerCtxID).(*i18n.Localizer); ok {
		return localizer
	}

	return nil
}

func WithLocalizer(r *http.Request, localizer *i18n.Localizer) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), localizerCtxID, localizer))
}

func NewI18nBundle(params I18nBundleParams) *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	params.AppLifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := bundle.LoadMessageFileFS(params.LocaleFS, "resources/trans/locale.en.yaml")
			return err
		},
	})

	return bundle
}

func NewI18nModule(fs fs.FS) fx.Option {
	return fx.Module(
		"I18n Module",
		fx.Provide(func(appLifeCycle fx.Lifecycle) *i18n.Bundle {
			return NewI18nBundle(I18nBundleParams{
				AppLifeCycle: appLifeCycle,
				LocaleFS:     fs,
			})
		}),
	)
}
