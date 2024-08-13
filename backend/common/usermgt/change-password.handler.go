package usermgt

import (
	"errors"
	"net/http"
	"wano-island/common/core"

	"log/slog"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type changePasswordHandler struct {
	db             *gorm.DB
	validator      *validator.Validate
	logger         core.Logger
	userService    UserService
	userRepository UserRepository
}

type ChangePasswordHandlerParams struct {
	fx.In
	DB             *gorm.DB
	Validator      *validator.Validate
	Logger         core.Logger
	UserService    UserService
	UserRepository UserRepository
}

type ChangePasswordRequest struct {
	CurrentPassword    string `json:"currentPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required,eqfield=NewPassword"`
}

var _ core.HTTPRoute = (*changePasswordHandler)(nil)

func NewChangePasswordHandler(params ChangePasswordHandlerParams) *changePasswordHandler {
	return &changePasswordHandler{
		db:             params.DB,
		validator:      params.Validator,
		logger:         params.Logger,
		userService:    params.UserService,
		userRepository: params.UserRepository,
	}
}

func (h *changePasswordHandler) Pattern() string {
	return "POST /api/v1/profile/change-password"
}

func (h *changePasswordHandler) IsPrivateRoute() bool {
	return true
}

func (h *changePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request ChangePasswordRequest

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgFailedToDecodeRequestBody).Build())

		return
	}

	authUser, ok := core.GetAuthUserFromRequest(r)
	if !ok {
		h.logger.ErrorContext(r.Context(), "Cannot get auth user from the request context")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	if err := h.validator.Struct(request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgValidationFailed).Build())

		return
	}

	user, err := h.userRepository.FindUserByID(r.Context(), h.db, authUser.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.ErrorContext(r.Context(), "The user does not exist in the database", slog.Any("details", err))
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgPasswordChangeErrorDueToUserNotFound).Build())

			return
		}

		h.logger.ErrorContext(r.Context(), "Something went wrong when trying to get the user by id", slog.Any("details", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	if err = h.userService.ComparePassword(
		r.Context(),
		[]byte(request.CurrentPassword),
		[]byte(user.Password),
	); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInvalidEmailOrPassword).Build())

			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	hashedPassword, err := h.userService.HashPassword(r.Context(), request.NewPassword)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	if _, err := h.userRepository.ChangePassword(
		r.Context(),
		h.db,
		user.ID.String(),
		string(*hashedPassword),
	); err != nil {
		h.logger.ErrorContext(
			r.Context(),
			"Something went wrong when trying to change the user password",
			slog.Any("details", err),
		)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgInternalServerError).Build())

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, core.NewResponseBuilder(r).MessageID(core.MsgSuccess).Build())
}
