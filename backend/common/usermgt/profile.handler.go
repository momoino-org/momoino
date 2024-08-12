package usermgt

import (
	"errors"
	"log/slog"
	"net/http"
	"time"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type profileHandler struct {
	logger         core.Logger
	db             *gorm.DB
	userRepository UserRepository
}

type ProfileHandlerParams struct {
	fx.In
	Logger         core.Logger
	DB             *gorm.DB
	UserRepository UserRepository
}

type ProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var _ core.HTTPRoute = (*profileHandler)(nil)

func NewProfileHandler(params ProfileHandlerParams) *profileHandler {
	return &profileHandler{
		db:             params.DB,
		userRepository: params.UserRepository,
	}
}

func (h *profileHandler) Pattern() string {
	return "GET /api/v1/profile"
}

func (h *profileHandler) IsPrivateRoute() bool {
	return true
}

func (h *profileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := core.GetAuthUser(r)
	responseBuilder := core.NewResponseBuilder(r)

	user, err := h.userRepository.FindUserByID(ctx, h.db, authUser.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidEmailOrPassword).Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when getting the user by email", slog.Any("details", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(&ProfileResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}).Build())
}
