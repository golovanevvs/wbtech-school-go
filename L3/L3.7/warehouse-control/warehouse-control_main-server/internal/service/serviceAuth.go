package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/zlog"
	"golang.org/x/crypto/bcrypt"
)

// IUserRpForAuth interface for user repository (needed for authentication)
type IUserRpForAuth interface {
	GetByUsername(username string) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	GetByID(id int) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
}

// IRefreshTokenRp interface for refresh token repository
type IRefreshTokenRp interface {
	Create(token *model.RefreshToken) (*model.RefreshToken, error)
	GetByToken(token string) (*model.RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID int) error
	DeleteExpired() error
}

// AuthService service for working with authentication
type AuthService struct {
	userRp         IUserRpForAuth
	refreshTokenRp IRefreshTokenRp
	cfg            *Config
}

// Claims represents JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *Config, ur IUserRpForAuth, rtr IRefreshTokenRp) *AuthService {
	return &AuthService{userRp: ur, refreshTokenRp: rtr, cfg: cfg}
}

// Register registers a new user
func (sv *AuthService) Register(ctx context.Context, username, password, name, role string) (*model.User, error) {
	existingUser, err := sv.userRp.GetByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		UserName:     username,
		Name:         name,
		PasswordHash: string(hashedPassword),
		UserRole:     role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := sv.userRp.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// DeleteUser deletes a user by ID
func (sv *AuthService) DeleteUser(ctx context.Context, id int) error {
	return sv.userRp.Delete(id)
}

// DeleteRefreshTokensByUserID deletes all refresh tokens for a user
func (sv *AuthService) DeleteRefreshTokensByUserID(ctx context.Context, userID int) error {
	lg := zlog.Logger.With().Str("service", "DeleteRefreshTokensByUserID").Logger()

	err := sv.refreshTokenRp.DeleteByUserID(userID)
	if err != nil {
		lg.Warn().Err(err).Int("userID", userID).Msgf("%s failed to delete refresh tokens", pkgConst.Warn)
		return fmt.Errorf("failed to delete refresh tokens: %w", err)
	}

	lg.Debug().Int("userID", userID).Msgf("%s all refresh tokens deleted successfully", pkgConst.OpSuccess)
	return nil
}

// Login authenticates user and returns access and refresh tokens
func (sv *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	lg := zlog.Logger.With().Str("service", "Login").Logger()

	user, err := sv.userRp.GetByUsername(username)
	if err != nil {
		return "", "", fmt.Errorf("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid username or password")
	}

	err = sv.refreshTokenRp.DeleteByUserID(user.ID)
	if err != nil {
		lg.Warn().Err(err).Int("userID", user.ID).Msgf("%s failed to delete old refresh tokens during login", pkgConst.Warn)
	} else {
		lg.Debug().Int("userID", user.ID).Msgf("%s old refresh tokens deleted successfully during login", pkgConst.OpSuccess)
	}

	accessToken, err := sv.generateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := sv.generateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	lg.Debug().Int("userID", user.ID).Str("username", username).Msgf("%s user logged in successfully", pkgConst.OpSuccess)

	return accessToken, refreshToken, nil
}

// generateAccessToken generates a new access token
func (sv *AuthService) generateAccessToken(user *model.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		UserRole: user.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sv.cfg.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "warehouse-control",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sv.cfg.JWTSecret))
}

// generateRefreshToken generates a new refresh token
func (sv *AuthService) generateRefreshToken(user *model.User) (string, error) {
	lg := zlog.Logger.With().Str("service", "generateRefreshToken").Logger()

	refreshTokenValue := fmt.Sprintf("refresh_%d_%d", user.ID, time.Now().Unix())

	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenValue,
		ExpiresAt: time.Now().Add(sv.cfg.RefreshTokenExpiry),
		CreatedAt: time.Now(),
	}

	_, err := sv.refreshTokenRp.Create(refreshToken)
	if err != nil {
		lg.Warn().Err(err).Int("userID", user.ID).Str("refreshTokenValue", refreshTokenValue).Msgf("%s failed to save refresh token to database", pkgConst.Warn)
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	lg.Debug().Int("userID", user.ID).Str("refreshTokenValue", refreshTokenValue).Msgf("%s refresh token saved to database successfully", pkgConst.OpSuccess)
	return refreshTokenValue, nil
}

// RefreshTokens refreshes access and refresh tokens
func (sv *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	lg := zlog.Logger.With().Str("service", "RefreshTokens").Logger()

	token, err := sv.refreshTokenRp.GetByToken(refreshToken)
	if err != nil {
		lg.Warn().Err(err).Str("refreshToken", refreshToken).Msgf("%s refresh token not found in database", pkgConst.Warn)
		return "", "", fmt.Errorf("invalid refresh token")
	}

	if token.ExpiresAt.Before(time.Now()) {
		lg.Warn().Str("refreshToken", refreshToken).Time("expiresAt", token.ExpiresAt).Msgf("%s refresh token expired", pkgConst.Warn)
		_ = sv.refreshTokenRp.DeleteByToken(refreshToken)
		return "", "", fmt.Errorf("refresh token has expired")
	}

	lg.Debug().Int("userID", token.UserID).Msgf("%s refresh token found, generating new tokens", pkgConst.Info)

	user, err := sv.userRp.GetByID(token.UserID)
	if err != nil {
		lg.Warn().Err(err).Int("userID", token.UserID).Msgf("%s user not found for refresh token", pkgConst.Warn)
		return "", "", fmt.Errorf("user not found")
	}

	newAccessToken, err := sv.generateAccessToken(user)
	if err != nil {
		lg.Warn().Err(err).Int("userID", token.UserID).Msgf("%s failed to generate new access token", pkgConst.Warn)
		return "", "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	err = sv.refreshTokenRp.DeleteByUserID(token.UserID)
	if err != nil {
		lg.Warn().Err(err).Int("userID", token.UserID).Msgf("%s failed to delete old refresh tokens", pkgConst.Warn)
		return "", "", fmt.Errorf("failed to delete old refresh tokens: %w", err)
	}

	lg.Debug().Int("userID", token.UserID).Msgf("%s all old refresh tokens deleted, generating new refresh token", pkgConst.Info)

	newRefreshToken, err := sv.generateRefreshToken(user)
	if err != nil {
		lg.Warn().Err(err).Int("userID", token.UserID).Msgf("%s failed to generate new refresh token", pkgConst.Warn)
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	lg.Debug().Int("userID", token.UserID).Msgf("%s tokens refreshed successfully", pkgConst.OpSuccess)
	return newAccessToken, newRefreshToken, nil
}

// ValidateToken validates an access token
func (sv *AuthService) ValidateToken(ctx context.Context, tokenString string) (int, string, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sv.cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, "", "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, "", "", errors.New("invalid token")
	}

	user, err := sv.userRp.GetByID(claims.UserID)
	if err != nil {
		return 0, "", "", fmt.Errorf("user not found: %w", err)
	}

	return claims.UserID, claims.UserRole, user.Name, nil
}

// GetUserByID returns a user by ID
func (sv *AuthService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return sv.userRp.GetByID(id)
}

// GetUserByUsername returns a user by username
func (sv *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return sv.userRp.GetByUsername(username)
}

// UpdateUser updates a user
func (sv *AuthService) UpdateUser(ctx context.Context, user *model.User) error {
	return sv.userRp.Update(user)
}
