package showmgt

import (
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/lib/pq"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type createMovieHandler struct {
	logger *slog.Logger
	db     *gorm.DB
}

type CreateShowHandlerParams struct {
	fx.In

	Logger *slog.Logger
	DB     *gorm.DB
}

type CreateShowRequestBody struct {
	Kind             string   `json:"kind"`
	OriginalLanguage string   `json:"originalLanguage"`
	OriginalTitle    string   `json:"originalTitle"`
	OriginalOverview *string  `json:"originalOverview"`
	Keywords         []string `json:"keywords"`
	IsReleased       bool     `json:"isReleased"`
}

var _ core.HTTPRoute = (*createMovieHandler)(nil)

func NewCreateMovieHandler(p CreateShowHandlerParams) *createMovieHandler {
	return &createMovieHandler{
		logger: p.Logger,
		db:     p.DB,
	}
}

func (h *createMovieHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern:        "POST /api/v1/shows",
		IsPrivate: true,
	}
}

func (h *createMovieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	var requestBody CreateShowRequestBody
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		h.logger.ErrorContext(reqCtx, "Something went wrong when trying to decode request body", core.DetailsLogAttr(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	showModel := ShowModel{
		Kind:             requestBody.Kind,
		OriginalLanguage: requestBody.OriginalLanguage,
		OriginalTitle:    requestBody.OriginalTitle,
		OriginalOverview: requestBody.OriginalOverview,
		Keywords:         pq.StringArray(requestBody.Keywords),
		IsReleased:       requestBody.IsReleased,
	}

	if result := h.db.Create(&showModel); result.Error != nil {
		h.logger.ErrorContext(reqCtx, "Something went wrong when creating a show", core.DetailsLogAttr(result.Error))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgCannotCreateTheShow).Build())

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, responseBuilder.MessageID(core.MsgSuccess).Data(ToShowDTO(&showModel)).Build())
}
