package authHandler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// ISvForAuthHandler interface for auth handler service
type ISvForAuthHandler interface {
	Register(ctx context.Context, username, password, name, role string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	ValidateToken(ctx context.Context, tokenString string) (int, string, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
}

// AuthHandler handles authentication requests
type AuthHandler struct {
	lg              *zlog.Zerolog
	sv              ISvForAuthHandler
	publicHost      string
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

// New creates a new AuthHandler
func New(parentLg *zlog.Zerolog, sv ISvForAuthHandler, publicHost string, accessTokenExp, refreshTokenExp time.Duration) *AuthHandler {
	lg := parentLg.With().Str("component", "handler-authHandler").Logger()
	return &AuthHandler{
		lg:              &lg,
		sv:              sv,
		publicHost:      publicHost,
		accessTokenExp:  accessTokenExp,
		refreshTokenExp: refreshTokenExp,
	}
}

// RegisterPublicRoutes registers the public auth routes (no authentication required)
func (hd *AuthHandler) RegisterPublicRoutes(rt *ginext.RouterGroup) {
	auth := rt.Group("")
	{
		auth.POST("/register", hd.Register)
		auth.POST("/login", hd.Login)
		auth.POST("/refresh", hd.Refresh)
		auth.POST("/logout", hd.Logout)
	}
}

// RegisterProtectedRoutes registers the protected auth routes (authentication required)
func (hd *AuthHandler) RegisterProtectedRoutes(rt *ginext.RouterGroup) {
	auth := rt.Group("/auth")
	{
		auth.GET("/me", hd.GetCurrentUser)
		auth.PUT("/update", hd.UpdateUser)
		auth.DELETE("/delete", hd.DeleteUser)
	}
}

// Register handles user registration
func (hd *AuthHandler) Register(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Register").Logger()

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Int("status", http.StatusBadRequest).Msgf("%s invalid content-type", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrContentTypeAJ.Error()})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	user, err := hd.sv.Register(c.Request.Context(), req.Username, req.Password, req.Name, req.Role)
	if err != nil {
		lg.Warn().Err(err).Str("username", req.Username).Msgf("%s failed to register user", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := hd.sv.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		lg.Warn().Err(err).Str("username", req.Username).Msgf("%s failed to generate tokens after registration", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to generate tokens"})
		return
	}

	lg.Debug().Int("user_id", user.ID).Str("username", user.UserName).Msgf("%s user registered successfully", pkgConst.OpSuccess)

	c.SetCookie("access_token", accessToken, int(hd.accessTokenExp.Seconds()), "/", extractDomain(hd.publicHost), true, true)
	c.SetCookie("refresh_token", refreshToken, int(hd.refreshTokenExp.Seconds()), "/", extractDomain(hd.publicHost), true, true)

	c.JSON(http.StatusCreated, ginext.H{
		"message": "Registration successful",
		"user":    user,
	})
}

// Login handles user login
func (hd *AuthHandler) Login(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Login").Logger()

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Int("status", http.StatusBadRequest).Msgf("%s invalid content-type", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrContentTypeAJ.Error()})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := hd.sv.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		lg.Warn().Err(err).Str("username", req.Username).Msgf("%s failed to login user", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Str("username", req.Username).Msgf("%s user logged in successfully", pkgConst.OpSuccess)

	c.SetCookie("access_token", accessToken, int(hd.accessTokenExp.Seconds()), "/", extractDomain(hd.publicHost), true, true)
	c.SetCookie("refresh_token", refreshToken, int(hd.refreshTokenExp.Seconds()), "/", extractDomain(hd.publicHost), true, true)

	c.JSON(http.StatusOK, ginext.H{
		"message": "Login successful",
	})
}

// Refresh handles token refresh
func (hd *AuthHandler) Refresh(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Refresh").Logger()

	var refreshToken string

	// Try to get refresh token from cookie first
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		// If not found in cookie, try to get from request body
		if strings.Contains(c.ContentType(), "application/json") {
			var req struct {
				RefreshToken string `json:"refreshToken" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
				c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
				return
			}
			refreshToken = req.RefreshToken
		} else {
			lg.Warn().Str("content-type", c.ContentType()).Msgf("%s No refresh token found in cookie or request body", pkgConst.Warn)
			c.JSON(http.StatusBadRequest, ginext.H{"error": "Refresh token is required"})
			return
		}
	}

	newAccessToken, newRefreshToken, err := hd.sv.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		lg.Warn().Err(err).Msgf("%s failed to refresh tokens", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Msgf("%s tokens refreshed successfully", pkgConst.OpSuccess)

	c.SetCookie("access_token", newAccessToken, int(hd.accessTokenExp.Seconds()), "/", extractDomain(hd.publicHost), true, true)
	c.SetCookie("refresh_token", newRefreshToken, int(hd.refreshTokenExp.Seconds()), "/", extractDomain(hd.publicHost), true, true)

	c.JSON(http.StatusOK, ginext.H{
		"message": "Tokens refreshed successfully",
	})
}

// GetCurrentUser handles getting current user
func (hd *AuthHandler) GetCurrentUser(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "GetCurrentUser").Logger()

	userID, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msgf("%s User ID not found in context", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		lg.Warn().Msgf("%s User ID is not of type int", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	user, err := hd.sv.GetUserByID(c.Request.Context(), userIDInt)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s user not found in database, clearing invalid cookies", pkgConst.Warn)

		c.SetCookie("access_token", "", -1, "/", extractDomain(hd.publicHost), true, true)
		c.SetCookie("refresh_token", "", -1, "/", extractDomain(hd.publicHost), true, true)

		c.JSON(http.StatusUnauthorized, ginext.H{"error": "User not found"})
		return
	}

	lg.Debug().Int("user_id", user.ID).Msgf("%s current user retrieved successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating user
func (hd *AuthHandler) UpdateUser(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "UpdateUser").Logger()

	userID, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msgf("%s User ID not found in context", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		lg.Warn().Msgf("%s User ID is not of type int", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	_, err := hd.sv.GetUserByID(c.Request.Context(), userIDInt)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s user not found in database, clearing invalid cookies", pkgConst.Warn)

		c.SetCookie("access_token", "", -1, "/", extractDomain(hd.publicHost), true, true)
		c.SetCookie("refresh_token", "", -1, "/", extractDomain(hd.publicHost), true, true)

		c.JSON(http.StatusUnauthorized, ginext.H{"error": "User not found"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		UserRole string `json:"user_role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	user := &model.User{
		ID:        userIDInt,
		Name:      req.Name,
		UserRole:  req.UserRole,
		UpdatedAt: time.Now(),
	}

	err = hd.sv.UpdateUser(c.Request.Context(), user)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s failed to update user", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to update user"})
		return
	}

	updatedUser, err := hd.sv.GetUserByID(c.Request.Context(), userIDInt)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s failed to get updated user", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to get updated user"})
		return
	}

	lg.Debug().Int("user_id", userIDInt).Msgf("%s user updated successfully", pkgConst.OpSuccess)

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles user deletion
func (hd *AuthHandler) DeleteUser(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "DeleteUser").Logger()

	userID, exists := c.Get("user_id")
	if !exists {
		lg.Warn().Msgf("%s User ID not found in context", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, ginext.H{"error": pkgErrors.ErrUnauthorized.Error()})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		lg.Warn().Msgf("%s User ID is not of type int", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Internal server error"})
		return
	}

	_, err := hd.sv.GetUserByID(c.Request.Context(), userIDInt)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s user not found in database, clearing invalid cookies", pkgConst.Warn)

		c.SetCookie("access_token", "", -1, "/", extractDomain(hd.publicHost), true, true)
		c.SetCookie("refresh_token", "", -1, "/", extractDomain(hd.publicHost), true, true)

		c.JSON(http.StatusUnauthorized, ginext.H{"error": "User not found"})
		return
	}

	err = hd.sv.DeleteUser(c.Request.Context(), userIDInt)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s failed to delete user", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to delete user"})
		return
	}

	lg.Debug().Int("user_id", userIDInt).Msgf("%s user deleted successfully", pkgConst.OpSuccess)

	c.SetCookie("access_token", "", -1, "/", extractDomain(hd.publicHost), true, true)
	c.SetCookie("refresh_token", "", -1, "/", extractDomain(hd.publicHost), true, true)

	c.JSON(http.StatusOK, ginext.H{
		"message": "User deleted successfully",
	})
}

// Logout handles user logout
func (hd *AuthHandler) Logout(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Logout").Logger()

	lg.Debug().Msg("User logout initiated")

	c.SetCookie("access_token", "", -1, "/", extractDomain(hd.publicHost), true, true)
	c.SetCookie("refresh_token", "", -1, "/", extractDomain(hd.publicHost), true, true)

	lg.Debug().Msgf("%s user logged out successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, ginext.H{
		"message": "Logout successful",
	})
}

// extractDomain extracts the domain from a URL
func extractDomain(url string) string {
	if strings.HasPrefix(url, "https://") {
		return strings.TrimPrefix(url, "https://")
	}
	if strings.HasPrefix(url, "http://") {
		return strings.TrimPrefix(url, "http://")
	}
	return url
}
