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

type AuthHandler struct {
	usecase *userusecase.UseCase
}

func NewAuthHandler(usecase *userusecase.UseCase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

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
		apierror.Internal(writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

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
		if errors.Is(err, userusecase.ErrInvalidToken) || errors.Is(err, userusecase.ErrTokenUsed) {
			apierror.Unauthorized(writer)
			return
		}
		apierror.Internal(writer)
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
		if errors.Is(err, userusecase.ErrInvalidToken) {
			apierror.Unauthorized(writer)
			return
		}
		apierror.Internal(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]string{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

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

	// Silently succeed even for invalid/revoked tokens
	_ = handler.usecase.Logout(request.Context(), body.RefreshToken)
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *AuthHandler) Me(writer http.ResponseWriter, request *http.Request) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok {
		apierror.Unauthorized(writer)
		return
	}

	currentUser, err := handler.usecase.GetUserByID(request.Context(), userID)
	if err != nil || currentUser == nil {
		apierror.Internal(writer)
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
