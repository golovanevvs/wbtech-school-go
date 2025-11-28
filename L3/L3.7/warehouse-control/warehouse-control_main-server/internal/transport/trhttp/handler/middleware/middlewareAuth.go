package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/zlog"
)

// ISvForAuthHandler interface for auth handler
type ISvForAuthHandler interface {
	ValidateToken(ctx context.Context, tokenString string) (userID int, userRole string, err error)
	// RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	lg *zlog.Zerolog
	sv ISvForAuthHandler
	// publicHost         string
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(parentLg *zlog.Zerolog, sv ISvForAuthHandler, accessTokenExp, refreshTokenExp time.Duration) *AuthMiddleware {
	lg := parentLg.With().Str("component", "middleware-auth").Logger()
	return &AuthMiddleware{
		lg: &lg,
		sv: sv,
		// publicHost:         publicHost,
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp,
	}
}

// JWTMiddleware is the middleware function for JWT authentication
func (mw *AuthMiddleware) JWTMiddleware(c *gin.Context) {
	lg := mw.lg.With().Str("middleware", "JWTMiddleware").Logger()

	tokenString, err := c.Cookie("access_token")
	if err != nil || tokenString == "" {
		lg.Warn().Msg("Access token cookie is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})
		c.Abort()
		return
	}

	lg.Debug().Msgf("%s Validating token...", pkgConst.Starting)

	userID, userRole, err := mw.sv.ValidateToken(c.Request.Context(), tokenString)
	if err != nil {
		// lg.Warn().Err(err).Msg("Failed to validate token, attempting refresh")

		lg.Warn().Msgf("%s Access token cookie is missing", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})

		// refreshToken, err := c.Cookie("refresh_token")
		// if err != nil || refreshToken == "" {
		// 	lg.Warn().Msg("Refresh token cookie is missing")
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})
		// 	c.Abort()
		// 	return
		// }

		// newAccessToken, newRefreshToken, err := mw.sv.RefreshTokens(c.Request.Context(), refreshToken)
		// if err != nil {
		// 	lg.Warn().Err(err).Msg("Failed to refresh tokens")
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})
		// 	c.Abort()
		// 	return
		// }

		// c.SetCookie("access_token", newAccessToken, int(mw.accessTokenExp.Seconds()), "/", extractDomain(mw.webHost), true, true)
		// c.SetCookie("refresh_token", newRefreshToken, int(mw.refreshTokenExp.Seconds()), "/", extractDomain(mw.webHost), true, true)

		// userID, userRole, err = mw.sv.ValidateToken(c.Request.Context(), newAccessToken)
		// if err != nil {
		// 	lg.Warn().Err(err).Msg("Failed to validate refreshed token")
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": pkgErrors.ErrUnauthorized.Error()})
		// 	c.Abort()
		// 	return
		// }

		// lg.Debug().Msg("Token refreshed and validated successfully")
	}

	c.Set("user_id", userID)
	c.Set("user_role", userRole)

	lg.Debug().Int("user_id", userID).Str("user_role", userRole).Msgf("%s Token validated successfully", pkgConst.Finished)

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

// extractDomain extracts the domain from a URL
// func extractDomain(url string) string {
// 	if strings.HasPrefix(url, "https://") {
// 		return strings.TrimPrefix(url, "https://")
// 	}
// 	if strings.HasPrefix(url, "http://") {
// 		return strings.TrimPrefix(url, "http://")
// 	}
// 	return url
// }
