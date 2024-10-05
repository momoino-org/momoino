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

// profileHandler handles requests related to user profiles.
// It implements the core.HTTPRoute interface to define the HTTP route for retrieving user profile information.
type profileHandler struct {
	logger         *slog.Logger
	db             *gorm.DB
	userRepository UserRepository
}

// ProfileHandlerParams contains the dependencies needed to create a new profileHandler.
type ProfileHandlerParams struct {
	fx.In
	Logger         *slog.Logger
	DB             *gorm.DB
	UserRepository UserRepository
}

// ProfileResponse represents the structure of the response containing user profile information.
type ProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Locale    string    `json:"locale"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Ensure profileHandler implements the core.HTTPRoute interface.
var _ core.HTTPRoute = (*profileHandler)(nil)

// NewProfileHandler creates a new instance of profileHandler with the provided parameters.
//
// Params:
//   - params: Dependencies required by the handler.
//
// Returns:
//   - *profileHandler: A new instance of profileHandler.
func NewProfileHandler(params ProfileHandlerParams) *profileHandler {
	return &profileHandler{
		logger:         params.Logger,
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

// ServeHTTP handles incoming HTTP requests for retrieving user profile information.
//
// Params:
//   - w: The http.ResponseWriter to send the response.
//   - r: The incoming http.Request containing the request for user profile.
func (h *profileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseBuilder := core.NewResponseBuilder(r)

	authUser := core.MustGetAuthUserFromRequest(r)

	// Retrieve user details based on authenticated user's ID.
	user, err := h.userRepository.FindUserByID(ctx, h.db, authUser.GetID())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, responseBuilder.MessageID(core.MsgInvalidEmailOrPassword).Build())

			return
		}

		h.logger.ErrorContext(ctx, "Something went wrong when getting the user by email", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, responseBuilder.MessageID(core.MsgInternalServerError).Build())

		return
	}

	// Respond with the user's profile information.
	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseBuilder.Data(&ProfileResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Locale:    user.Locale,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}).Build())
}
