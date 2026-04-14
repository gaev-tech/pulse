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
	magicLinkTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

var (
	ErrInvalidToken  = errors.New("invalid or expired token")
	ErrTokenUsed     = errors.New("token already used")
)

type UseCase struct {
	users         user.Repository
	magicLinks    user.MagicLinkRepository
	refreshTokens user.RefreshTokenRepository
	jwt           *jwt.Manager
	email         email.Sender
	frontendURL   string
}

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

func (useCase *UseCase) SendMagicLink(ctx context.Context, emailAddr string) error {
	rawToken, err := generateToken()
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	tokenHash := hashToken(rawToken)
	expiresAt := time.Now().Add(magicLinkTTL)

	if err := useCase.magicLinks.Create(ctx, emailAddr, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("create magic link: %w", err)
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", useCase.frontendURL, rawToken)
	if err := useCase.email.SendMagicLink(ctx, emailAddr, link); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

type VerifyResult struct {
	AccessToken  string
	RefreshToken string
	User         *user.User
}

func (useCase *UseCase) VerifyMagicLink(ctx context.Context, rawToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawToken)

	magicToken, err := useCase.magicLinks.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("get magic link: %w", err)
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

	existingUser, err := useCase.users.GetByEmail(ctx, magicToken.Email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	var currentUser *user.User
	if existingUser != nil {
		currentUser = existingUser
	} else {
		username, err := useCase.generateUsername(ctx, magicToken.Email)
		if err != nil {
			return nil, fmt.Errorf("generate username: %w", err)
		}
		currentUser, err = useCase.users.Create(ctx, magicToken.Email, username)
		if err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	if err := useCase.magicLinks.MarkUsed(ctx, magicToken.ID); err != nil {
		return nil, fmt.Errorf("mark used: %w", err)
	}

	return useCase.issueTokenPair(ctx, currentUser)
}

func (useCase *UseCase) Refresh(ctx context.Context, rawRefreshToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawRefreshToken)

	storedToken, err := useCase.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return nil, ErrInvalidToken
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	if err := useCase.refreshTokens.Revoke(ctx, storedToken.ID); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	currentUser, err := useCase.users.GetByID(ctx, storedToken.UserID)
	if err != nil || currentUser == nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return useCase.issueTokenPair(ctx, currentUser)
}

func (useCase *UseCase) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := hashToken(rawRefreshToken)

	storedToken, err := useCase.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("get refresh token: %w", err)
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return ErrInvalidToken
	}

	return useCase.refreshTokens.Revoke(ctx, storedToken.ID)
}

func (useCase *UseCase) issueTokenPair(ctx context.Context, currentUser *user.User) (*VerifyResult, error) {
	accessToken, err := useCase.jwt.GenerateAccessToken(currentUser.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	rawRefreshToken, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshTokenHash := hashToken(rawRefreshToken)
	expiresAt := time.Now().Add(refreshTokenTTL)

	if _, err := useCase.refreshTokens.Create(ctx, currentUser.ID, refreshTokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &VerifyResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		User:         currentUser,
	}, nil
}

func (useCase *UseCase) generateUsername(ctx context.Context, emailAddr string) (string, error) {
	prefix := strings.Split(emailAddr, "@")[0]
	prefix = strings.ToLower(prefix)

	exists, err := useCase.users.ExistsByUsername(ctx, prefix)
	if err != nil {
		return "", err
	}
	if !exists {
		return prefix, nil
	}

	for range 10 {
		suffix, err := generateToken()
		if err != nil {
			return "", err
		}
		candidate := prefix + "_" + suffix[:4]
		exists, err := useCase.users.ExistsByUsername(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not generate unique username")
}

func (useCase *UseCase) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return useCase.users.GetByID(ctx, id)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
