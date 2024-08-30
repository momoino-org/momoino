package core

import (
	"context"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

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
	MsgInvalidCurrentPassword               = "E-0006"
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
			files, err := glob(params.LocaleFS, "resources/trans/locale.*.yaml")
			if err != nil {
				return err
			}

			for _, file := range files {
				if _, err := bundle.LoadMessageFileFS(params.LocaleFS, file); err != nil {
					return err
				}
			}

			return nil
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

// glob is a function that performs a glob-style pattern matching on a given file system (fs.FS)
// and returns a list of matching file paths. It uses fs.WalkDir to traverse the file system
// and filepath.Match to check if each file path matches the provided pattern.
func glob(f fs.FS, pattern string) ([]string, error) {
	var matches []string

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && !strings.HasSuffix(pattern, "/**") {
			return nil
		}

		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return err
		}

		if matched {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}
