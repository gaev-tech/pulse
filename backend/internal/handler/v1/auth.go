package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gaevivan/pulse/internal/handler/middleware"
	userusecase "github.com/gaevivan/pulse/internal/usecase/user"
	"github.com/gaevivan/pulse/pkg/apierror"
	"github.com/gaevivan/pulse/pkg/validator"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	usecase *userusecase.UseCase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(usecase *userusecase.UseCase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

// SendMagicLink godoc
// @Summary     Отправить magic-link
// @Tags        auth
// @Accept      json
// @Param       body body object{email=string} true "Email"
// @Success     200
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/magic-link [post]
func (handler *AuthHandler) SendMagicLink(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := validator.Decode(request, &body); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}
	if err := validator.Email(body.Email); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}

	if err := handler.usecase.SendMagicLink(request.Context(), body.Email); err != nil {
		switch {
		case errors.Is(err, userusecase.ErrEmailUnavailable):
			apierror.EmailUnavailable(writer)
		case errors.Is(err, userusecase.ErrDatabaseUnavailable):
			apierror.DatabaseUnavailable(writer)
		case errors.Is(err, userusecase.ErrTokenGenerationFailed):
			apierror.TokenGenerationFailed(writer)
		default:
			apierror.Internal(writer)
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
}

// VerifyMagicLink godoc
// @Summary     Верифицировать magic-link токен
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object{token=string} true "Токен"
// @Success     200 {object} object{access_token=string,refresh_token=string,user=object{id=string,email=string,username=string}}
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Failure     401 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/magic-link/verify [post]
func (handler *AuthHandler) VerifyMagicLink(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Token string `json:"token"`
	}
	if err := validator.Decode(request, &body); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}
	if err := validator.Required(body.Token, "token"); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}

	result, err := handler.usecase.VerifyMagicLink(request.Context(), body.Token)
	if err != nil {
		switch {
		case errors.Is(err, userusecase.ErrInvalidToken), errors.Is(err, userusecase.ErrTokenUsed):
			apierror.Unauthorized(writer)
		case errors.Is(err, userusecase.ErrDatabaseUnavailable):
			apierror.DatabaseUnavailable(writer)
		case errors.Is(err, userusecase.ErrTokenGenerationFailed):
			apierror.TokenGenerationFailed(writer)
		case errors.Is(err, userusecase.ErrUsernameConflict):
			apierror.UsernameConflict(writer)
		default:
			apierror.Internal(writer)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]any{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"user": map[string]string{
			"id":       result.User.ID,
			"email":    result.User.Email,
			"username": result.User.Username,
		},
	})
}

// Refresh godoc
// @Summary     Обновить токены
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object{refresh_token=string} true "Refresh token"
// @Success     200 {object} object{access_token=string,refresh_token=string}
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Failure     401 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/refresh [post]
func (handler *AuthHandler) Refresh(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := validator.Decode(request, &body); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}
	if err := validator.Required(body.RefreshToken, "refresh_token"); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}

	result, err := handler.usecase.Refresh(request.Context(), body.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, userusecase.ErrInvalidToken):
			apierror.Unauthorized(writer)
		case errors.Is(err, userusecase.ErrUserNotFound):
			apierror.UserNotFound(writer)
		case errors.Is(err, userusecase.ErrDatabaseUnavailable):
			apierror.DatabaseUnavailable(writer)
		case errors.Is(err, userusecase.ErrTokenGenerationFailed):
			apierror.TokenGenerationFailed(writer)
		default:
			apierror.Internal(writer)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]string{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

// Logout godoc
// @Summary     Выйти из системы
// @Tags        auth
// @Accept      json
// @Param       body body object{refresh_token=string} true "Refresh token"
// @Success     204
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/logout [post]
func (handler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := validator.Decode(request, &body); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}
	if err := validator.Required(body.RefreshToken, "refresh_token"); err != nil {
		apierror.BadRequest(writer, err.Error())
		return
	}

	// Silently succeed even for invalid/revoked tokens.
	_ = handler.usecase.Logout(request.Context(), body.RefreshToken)
	writer.WriteHeader(http.StatusNoContent)
}

// Me godoc
// @Summary     Текущий пользователь
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} object{id=string,email=string,username=string}
// @Failure     401 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/me [get]
func (handler *AuthHandler) Me(writer http.ResponseWriter, request *http.Request) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok {
		apierror.Unauthorized(writer)
		return
	}

	currentUser, err := handler.usecase.GetUserByID(request.Context(), userID)
	if err != nil {
		apierror.DatabaseUnavailable(writer)
		return
	}
	if currentUser == nil {
		apierror.UserNotFound(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]string{
		"id":       currentUser.ID,
		"email":    currentUser.Email,
		"username": currentUser.Username,
	})
}
