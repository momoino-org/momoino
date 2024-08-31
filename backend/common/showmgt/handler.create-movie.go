package showmgt

import (
	"log/slog"
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	Kind             string    `json:"kind"`
	OriginalLanguage string    `json:"originalLanguage"`
	OriginalTitle    string    `json:"originalTitle"`
	OriginalOverview *string   `json:"originalOverview"`
	Keywords         *[]string `json:"keywords"`
	IsReleased       bool      `json:"isReleased"`
}

type CreateShowResponseBody struct {
	ID               uuid.UUID `json:"id"`
	Kind             string    `json:"kind"`
	OriginalLanguage string    `json:"originalLanguage"`
	OriginalTitle    string    `json:"originalTitle"`
	OriginalOverview *string   `json:"originalOverview"`
	Keywords         *[]string `json:"keywords"`
	IsReleased       bool      `json:"isReleased"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

var _ core.HTTPRoute = (*createMovieHandler)(nil)

func NewCreateMovieHandler(p CreateShowHandlerParams) *createMovieHandler {
	return &createMovieHandler{
		logger: p.Logger,
		db:     p.DB,
	}
}

func (h *createMovieHandler) Pattern() string {
	return "POST /api/v1/shows"
}

func (h *createMovieHandler) IsPrivateRoute() bool {
	return true
}

func (h *createMovieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	var requestBody CreateShowRequestBody
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		h.logger.ErrorContext(reqCtx, "Something went wrong when trying to decode request body", slog.Any("details", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	showModel := ShowModel{
		Kind:             requestBody.Kind,
		OriginalLanguage: requestBody.OriginalLanguage,
		OriginalTitle:    requestBody.OriginalTitle,
		OriginalOverview: requestBody.OriginalOverview,
		Keywords:         requestBody.Keywords,
		IsReleased:       requestBody.IsReleased,
	}

	if result := h.db.Create(&showModel); result.Error != nil {
		h.logger.ErrorContext(reqCtx, "Something went wrong when creating a show", slog.Any("details", result.Error))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgCannotCreateTheShow).Build())

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.MessageID(core.MsgSuccess).Data(&CreateShowResponseBody{
		ID:               showModel.ID,
		Kind:             showModel.Kind,
		OriginalLanguage: showModel.OriginalLanguage,
		OriginalTitle:    showModel.OriginalTitle,
		OriginalOverview: showModel.OriginalOverview,
		Keywords:         showModel.Keywords,
		IsReleased:       showModel.IsReleased,
		CreatedAt:        showModel.CreatedAt,
		UpdatedAt:        showModel.UpdatedAt,
	}).Build())
}
