package usermgt

import (
	"log/slog"
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

type getOAuth2ProvidersHandler struct {
	logger                   *slog.Logger
	oauth2ProviderRepository OAuth2ProviderRepository
}

type GetOAuth2ProvidersHandlerParams struct {
	fx.In
	Logger                   *slog.Logger
	OAuth2ProviderRepository OAuth2ProviderRepository
}

type OAuth2ProviderDTO struct {
	ID        string    `json:"id"`
	Provider  string    `json:"name"`
	IsEnabled bool      `json:"isEnabled"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

var _ core.HTTPRoute = (*getOAuth2ProvidersHandler)(nil)

func NewGetOAuth2ProvidersHandler(params GetOAuth2ProvidersHandlerParams) *getOAuth2ProvidersHandler {
	return &getOAuth2ProvidersHandler{
		logger:                   params.Logger,
		oauth2ProviderRepository: params.OAuth2ProviderRepository,
	}
}

func (h *getOAuth2ProvidersHandler) Pattern() string {
	return "GET /api/v1/providers"
}

func (h *getOAuth2ProvidersHandler) IsPrivateRoute() bool {
	return true
}

func (h *getOAuth2ProvidersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseBuilder := core.NewResponseBuilder(r)

	providers, totalRows, err := h.oauth2ProviderRepository.GetMany(r)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot get providers", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	providerDTOs := lo.Map(providers, func(provider OAuth2ProviderModel, _ int) OAuth2ProviderDTO {
		return OAuth2ProviderDTO{
			ID:        provider.ID.String(),
			Provider:  provider.Provider,
			IsEnabled: provider.IsEnabled,
			CreatedAt: provider.CreatedAt,
			CreatedBy: provider.CreatedBy,
		}
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(providerDTOs).Pagination(*totalRows).Build())
}
