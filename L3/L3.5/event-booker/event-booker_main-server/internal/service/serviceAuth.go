package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// IUserRpForAuth interface for user repository (needed for authentication)
type IUserRpForAuth interface {
	GetByEmail(email string) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	GetByID(id int) (*model.User, error)
	Update(user *model.User) error
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
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *Config, ur IUserRpForAuth, rtr IRefreshTokenRp) *AuthService {
	return &AuthService{userRp: ur, refreshTokenRp: rtr, cfg: cfg}
}

// Register registers a new user
func (sv *AuthService) Register(ctx context.Context, email, password, name string) (*model.User, error) {
	existingUser, err := sv.userRp.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		Name:         name,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := sv.userRp.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// GetUserByID returns a user by ID
func (sv *AuthService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return sv.userRp.GetByID(id)
}

// UpdateUser updates a user
func (sv *AuthService) UpdateUser(ctx context.Context, user *model.User) error {
	return sv.userRp.Update(user)
}

// Login authenticates user and returns access and refresh tokens
func (sv *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := sv.userRp.GetByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid email or password")
	}

	accessToken, err := sv.generateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := sv.generateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateAccessToken generates a new access token
func (sv *AuthService) generateAccessToken(user *model.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sv.cfg.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "event-booker",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sv.cfg.JWTSecret))
}

// generateRefreshToken generates a new refresh token
func (sv *AuthService) generateRefreshToken(user *model.User) (string, error) {
	refreshTokenValue := fmt.Sprintf("refresh_%d_%d", user.ID, time.Now().Unix())

	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenValue,
		ExpiresAt: time.Now().Add(sv.cfg.RefreshTokenExpiry),
		CreatedAt: time.Now(),
	}

	_, err := sv.refreshTokenRp.Create(refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return refreshTokenValue, nil
}

// RefreshTokens refreshes access and refresh tokens
func (sv *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := sv.refreshTokenRp.GetByToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = sv.refreshTokenRp.DeleteByToken(refreshToken)
		return "", "", fmt.Errorf("refresh token has expired")
	}

	user, err := sv.userRp.GetByID(token.UserID)
	if err != nil {
		return "", "", fmt.Errorf("user not found")
	}

	newAccessToken, err := sv.generateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	err = sv.refreshTokenRp.DeleteByToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	newRefreshToken, err := sv.generateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

// ValidateToken validates an access token
func (sv *AuthService) ValidateToken(ctx context.Context, tokenString string) (int, string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sv.cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	return claims.UserID, claims.Email, nil
}
