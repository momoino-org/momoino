package usermgt

import (
	"errors"
	"net/http"
	"wano-island/common/core"

	"log/slog"

	"github.com/go-chi/render"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// changePasswordHandler is responsible for handling password change requests.
// It implements the core.HTTPRoute interface to define the HTTP route for changing passwords.
type changePasswordHandler struct {
	db                  *gorm.DB
	logger              *slog.Logger
	validator           *validator.Validate
	universalTranslator *ut.UniversalTranslator
	userService         UserService
	userRepository      UserRepository
}

// ChangePasswordHandlerParams contains the dependencies needed to create a new changePasswordHandler.
type ChangePasswordHandlerParams struct {
	fx.In
	DB                  *gorm.DB
	Logger              *slog.Logger
	Validator           *validator.Validate
	UniversalTranslator *ut.UniversalTranslator
	UserService         UserService
	UserRepository      UserRepository
}

// ChangePasswordRequest represents the payload for a password change request.
type ChangePasswordRequest struct {
	// Current password of the user
	CurrentPassword string `json:"currentPassword" validate:"required"`

	// New password to set
	NewPassword string `json:"newPassword" validate:"required"`

	// Confirmation of the new password
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required,eqfield=newPassword"`
}

// Ensure changePasswordHandler implements the core.HTTPRoute interface.
var _ core.HTTPRoute = (*changePasswordHandler)(nil)

// NewChangePasswordHandler creates a new instance of changePasswordHandler with the provided parameters.
//
// Params:
//   - params: Dependencies required by the handler.
//
// Returns:
//   - *changePasswordHandler: A new instance of changePasswordHandler.
func NewChangePasswordHandler(params ChangePasswordHandlerParams) *changePasswordHandler {
	return &changePasswordHandler{
		db:                  params.DB,
		logger:              params.Logger,
		validator:           params.Validator,
		universalTranslator: params.UniversalTranslator,
		userService:         params.UserService,
		userRepository:      params.UserRepository,
	}
}

func (h *changePasswordHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern:   "POST /api/v1/profile/change-password",
		IsPrivate: true,
	}
}

// ServeHTTP handles incoming HTTP requests for changing passwords.
//
// Params:
//   - w: The http.ResponseWriter to send the response.
//   - r: The incoming http.Request containing the password change details.
func (h *changePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request body.
	var request ChangePasswordRequest
	if err := render.DecodeJSON(r.Body, &request); err != nil {
		h.logger.ErrorContext(r.Context(), "Cannot decode the request body", core.DetailsLogAttr(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	authUser := core.MustGetAuthUserFromRequest(r)

	// Validate the request fields.
	if err := h.validator.Struct(request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, core.NewResponseBuilder(r).
			MessageID(core.MsgValidationFailed).
			Data(core.TranslateValidationErrors(r, h.universalTranslator, err)).
			Build())

		return
	}

	// Find the user by ID.
	user, err := h.userRepository.FindUserByID(r.Context(), h.db, authUser.GetID())

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.ErrorContext(r.Context(), "The user does not exist in the database", core.DetailsLogAttr(err))
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgPasswordChangeErrorDueToUserNotFound).Build())

			return
		}

		h.logger.ErrorContext(r.Context(), "Something went wrong when trying to get the user by id", core.DetailsLogAttr(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	// Compare the current password with the stored hashed password.
	if err = h.userService.ComparePassword(
		r.Context(),
		[]byte(request.CurrentPassword),
		[]byte(*user.Password),
	); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInvalidCurrentPassword).Build())

			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	// Hash the new password.
	hashedPassword, err := h.userService.HashPassword(r.Context(), request.NewPassword)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	// Update the user's password in the repository.
	if _, err := h.userRepository.ChangePassword(
		r.Context(),
		h.db,
		user.ID.String(),
		string(hashedPassword),
	); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"Something went wrong when trying to change the user password",
			core.DetailsLogAttr(err),
		)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	// Respond with success.
	render.Status(r, http.StatusOK)
	render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgSuccess).Build())
}
