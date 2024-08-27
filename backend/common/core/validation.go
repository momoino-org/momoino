package core

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"go.uber.org/fx"
	"golang.org/x/text/language"
)

func NewValidator(uni *ut.UniversalTranslator) *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	trans, _ := uni.GetTranslator(language.English.String())

	_ = en_translations.RegisterDefaultTranslations(v, trans)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		//nolint:mnd // No need to fix
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

func TranslateValidationErrors(
	r *http.Request,
	uni *ut.UniversalTranslator,
	v error,
) validator.ValidationErrorsTranslations {
	var validationErrs validator.ValidationErrors

	if errors.As(v, &validationErrs) {
		authUser := MustGetAuthUserFromRequest(r)
		trans, _ := uni.GetTranslator(authUser.Locale)
		translations := validator.ValidationErrorsTranslations{}

		for _, v := range validationErrs {
			translations[v.Field()] = v.Translate(trans)
		}

		return translations
	}

	panic("Only support translating validator.ValidationErrors")
}

func NewValidationModule() fx.Option {
	return fx.Module("Validation Module", fx.Provide(NewValidator))
}
