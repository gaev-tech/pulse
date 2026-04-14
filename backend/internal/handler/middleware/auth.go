package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gaevivan/pulse/internal/domain/user"
	"github.com/gaevivan/pulse/internal/infrastructure/jwt"
	"github.com/gaevivan/pulse/pkg/apierror"
)

type contextKey string

// UserIDKey is the context key for the authenticated user ID.
const UserIDKey contextKey = "user_id"

// UserIDFromContext returns the authenticated user ID from ctx.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok && userID != ""
}

// Auth is HTTP middleware for JWT and PAT authentication.
type Auth struct {
	jwt  *jwt.Manager
	pats user.PATRepository
}

// NewAuth creates a new Auth middleware.
func NewAuth(jwtManager *jwt.Manager, pats user.PATRepository) *Auth {
	return &Auth{jwt: jwtManager, pats: pats}
}

// Required enforces authentication, returning 401 if no valid credential is present.
func (auth *Auth) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userID, err := auth.resolveUserID(request)
		if err != nil || userID == "" {
			apierror.Unauthorized(writer)
			return
		}
		ctx := context.WithValue(request.Context(), UserIDKey, userID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (auth *Auth) resolveUserID(request *http.Request) (string, error) {
	if patHeader := request.Header.Get("X-API-Token"); patHeader != "" {
		hash := sha256Hash(patHeader)
		return auth.pats.GetUserIDByHash(request.Context(), hash)
	}

	authHeader := request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return auth.jwt.ParseAccessToken(token)
}

func sha256Hash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}
