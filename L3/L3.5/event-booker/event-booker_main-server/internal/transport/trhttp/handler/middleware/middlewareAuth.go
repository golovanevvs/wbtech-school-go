package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/zlog"
)

// IServiceForAuthHandler interface for auth handler
type IServiceForAuthHandler interface {
	ValidateToken(ctx context.Context, tokenString string) (userID int, userEmail string, err error)
}

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	lg *zlog.Zerolog
	sv IServiceForAuthHandler
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(parentLg *zlog.Zerolog, sv IServiceForAuthHandler) *AuthMiddleware {
	lg := parentLg.With().Str("component", "middleware-auth").Logger()
	return &AuthMiddleware{
		lg: &lg,
		sv: sv,
	}
}

// JWTMiddleware is the middleware function for JWT authentication
func (mw *AuthMiddleware) JWTMiddleware(c *gin.Context) {
	lg := mw.lg.With().Str("middleware", "JWTMiddleware").Logger()

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		lg.Warn().Msg("Authorization header is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		lg.Warn().Msg("Authorization header format is invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	userID, userEmail, err := mw.sv.ValidateToken(c.Request.Context(), tokenString)
	if err != nil {
		lg.Warn().Err(err).Msg("Failed to validate token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	// Set user ID in context for use in handlers
	c.Set("user_id", userID)
	c.Set("user_email", userEmail)

	c.Next()
}

// GetUserIDFromContext extracts user ID from gin context
func GetUserIDFromContext(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user_id not found in context")
	}

	userIDInt, ok := userID.(int)
	if !ok {
		return 0, errors.New("user_id is not of type int")
	}

	return userIDInt, nil
}
