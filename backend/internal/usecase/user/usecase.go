package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gaevivan/pulse/internal/domain/user"
	"github.com/gaevivan/pulse/internal/infrastructure/email"
	"github.com/gaevivan/pulse/internal/infrastructure/jwt"
)

const (
	magicLinkTTL        = 15 * time.Minute
	refreshTokenTTL     = 30 * 24 * time.Hour
	maxUsernameAttempts = 10
)

var (
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrTokenUsed             = errors.New("token already used")
	ErrEmailUnavailable      = errors.New("email service unavailable")
	ErrDatabaseUnavailable   = errors.New("database unavailable")
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrUsernameConflict      = errors.New("could not generate unique username")
	ErrUserNotFound          = errors.New("user not found")
)

// UseCase implements auth business logic.
type UseCase struct {
	users         user.Repository
	magicLinks    user.MagicLinkRepository
	refreshTokens user.RefreshTokenRepository
	jwt           *jwt.Manager
	email         email.Sender
	frontendURL   string
}

// New creates a new UseCase.
func New(
	users user.Repository,
	magicLinks user.MagicLinkRepository,
	refreshTokens user.RefreshTokenRepository,
	jwtManager *jwt.Manager,
	emailSender email.Sender,
	frontendURL string,
) *UseCase {
	return &UseCase{
		users:         users,
		magicLinks:    magicLinks,
		refreshTokens: refreshTokens,
		jwt:           jwtManager,
		email:         emailSender,
		frontendURL:   frontendURL,
	}
}

// SendMagicLink generates a magic link and emails it to the given address.
func (usecase *UseCase) SendMagicLink(ctx context.Context, emailAddr string) error {
	rawToken, err := generateToken()
	if err != nil {
		return errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate token: %w", err))
	}

	tokenHash := hashToken(rawToken)
	expiresAt := time.Now().Add(magicLinkTTL)

	if err := usecase.magicLinks.Create(ctx, emailAddr, tokenHash, expiresAt); err != nil {
		return errors.Join(ErrDatabaseUnavailable, fmt.Errorf("create magic link: %w", err))
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", usecase.frontendURL, rawToken)
	if err := usecase.email.SendMagicLink(ctx, emailAddr, link); err != nil {
		return errors.Join(ErrEmailUnavailable, fmt.Errorf("send email: %w", err))
	}

	return nil
}

// VerifyResult holds the tokens and user returned after successful magic-link verification.
type VerifyResult struct {
	AccessToken  string
	RefreshToken string
	User         *user.User
}

// VerifyMagicLink validates the raw token and returns a token pair and the user.
func (usecase *UseCase) VerifyMagicLink(ctx context.Context, rawToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawToken)

	magicToken, err := usecase.magicLinks.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get magic link: %w", err))
	}
	if magicToken == nil {
		return nil, ErrInvalidToken
	}
	if magicToken.UsedAt != nil {
		return nil, ErrTokenUsed
	}
	if time.Now().After(magicToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	existingUser, err := usecase.users.GetByEmail(ctx, magicToken.Email)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get user: %w", err))
	}

	currentUser := existingUser
	if currentUser == nil {
		username, err := usecase.generateUsername(ctx, magicToken.Email)
		if err != nil {
			return nil, fmt.Errorf("generate username: %w", err)
		}
		currentUser, err = usecase.users.Create(ctx, magicToken.Email, username)
		if err != nil {
			return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("create user: %w", err))
		}
	}

	if err := usecase.magicLinks.MarkUsed(ctx, magicToken.ID); err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("mark used: %w", err))
	}

	return usecase.issueTokenPair(ctx, currentUser)
}

// Refresh rotates the refresh token and returns a new token pair.
func (usecase *UseCase) Refresh(ctx context.Context, rawRefreshToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawRefreshToken)

	storedToken, err := usecase.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get refresh token: %w", err))
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return nil, ErrInvalidToken
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	if err := usecase.refreshTokens.Revoke(ctx, storedToken.ID); err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("revoke old token: %w", err))
	}

	currentUser, err := usecase.users.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get user: %w", err))
	}
	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	return usecase.issueTokenPair(ctx, currentUser)
}

// Logout revokes the given refresh token. Silently succeeds for unknown tokens.
func (usecase *UseCase) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := hashToken(rawRefreshToken)

	storedToken, err := usecase.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("get refresh token: %w", err)
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return ErrInvalidToken
	}

	return usecase.refreshTokens.Revoke(ctx, storedToken.ID)
}

// GetUserByID returns the user with the given ID, or nil if not found.
func (usecase *UseCase) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return usecase.users.GetByID(ctx, id)
}

func (usecase *UseCase) issueTokenPair(ctx context.Context, currentUser *user.User) (*VerifyResult, error) {
	accessToken, err := usecase.jwt.GenerateAccessToken(currentUser.ID)
	if err != nil {
		return nil, errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate access token: %w", err))
	}

	rawRefreshToken, err := generateToken()
	if err != nil {
		return nil, errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate refresh token: %w", err))
	}

	refreshTokenHash := hashToken(rawRefreshToken)
	expiresAt := time.Now().Add(refreshTokenTTL)

	if _, err := usecase.refreshTokens.Create(ctx, currentUser.ID, refreshTokenHash, expiresAt); err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("store refresh token: %w", err))
	}

	return &VerifyResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		User:         currentUser,
	}, nil
}

func (usecase *UseCase) generateUsername(ctx context.Context, emailAddr string) (string, error) {
	prefix := strings.ToLower(strings.Split(emailAddr, "@")[0])

	exists, err := usecase.users.ExistsByUsername(ctx, prefix)
	if err != nil {
		return "", errors.Join(ErrDatabaseUnavailable, fmt.Errorf("check username: %w", err))
	}
	if !exists {
		return prefix, nil
	}

	for range maxUsernameAttempts {
		suffix, err := generateToken()
		if err != nil {
			return "", errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate suffix: %w", err))
		}
		candidate := prefix + "_" + suffix[:4]
		exists, err := usecase.users.ExistsByUsername(ctx, candidate)
		if err != nil {
			return "", errors.Join(ErrDatabaseUnavailable, fmt.Errorf("check username: %w", err))
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", ErrUsernameConflict
}

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
