package showmgt

import (
	"log/slog"
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type getShowsHandler struct {
	logger *slog.Logger
	db     *gorm.DB
}

type GetShowsHandlerParams struct {
	fx.In

	Logger *slog.Logger
	DB     *gorm.DB
}

var _ core.HTTPRoute = (*getShowsHandler)(nil)

func NewGetShowsHandler(p GetShowsHandlerParams) *getShowsHandler {
	return &getShowsHandler{
		logger: p.Logger,
		db:     p.DB,
	}
}

func (h *getShowsHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern:        "GET /api/v1/shows",
		IsPrivate: true,
	}
}

func (h *getShowsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	var showModels []ShowModel

	pageSize := core.GetPageSize(r)
	offset := core.GetOffset(r)

	var totalRows int64
	if result := h.db.Model(&ShowModel{}).Count(&totalRows); result.Error != nil {
		h.logger.ErrorContext(reqCtx, "Something went wrong when getting total rows", core.DetailsLogAttr(result.Error))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	if result := h.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&showModels); result.Error != nil {
		h.logger.ErrorContext(reqCtx, "Something went wrong when getting shows", core.DetailsLogAttr(result.Error))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	showDTOs := lo.Map(showModels, func(showModel ShowModel, _ int) *ShowDTO {
		return ToShowDTO(&showModel)
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(showDTOs).Pagination(totalRows).Build())
}
