package authHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgErrors"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// ISvForAuthHandler interface for auth handler service
type ISvForAuthHandler interface {
	Register(ctx context.Context, email, password, name string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	ValidateToken(ctx context.Context, tokenString string) (int, string, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
}

// AuthHandler handles authentication requests
type AuthHandler struct {
	lg *zlog.Zerolog
	sv ISvForAuthHandler
}

// New creates a new AuthHandler
func New(parentLg *zlog.Zerolog, sv ISvForAuthHandler) *AuthHandler {
	lg := parentLg.With().Str("component", "handler-authHandler").Logger()
	return &AuthHandler{
		lg: &lg,
		sv: sv,
	}
}

// RegisterPublicRoutes registers the public auth routes (no authentication required)
func (hd *AuthHandler) RegisterPublicRoutes(rt *ginext.RouterGroup) {
	auth := rt.Group("")
	{
		auth.POST("/register", hd.Register)
		auth.POST("/login", hd.Login)
		auth.POST("/refresh", hd.Refresh)
	}
}

// RegisterProtectedRoutes registers the protected auth routes (authentication required)
func (hd *AuthHandler) RegisterProtectedRoutes(rt *ginext.RouterGroup) {
	auth := rt.Group("/auth")
	{
		auth.GET("/me", hd.GetCurrentUser)
		auth.PUT("/update", hd.UpdateUser)
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
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	user, err := hd.sv.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		lg.Warn().Err(err).Str("email", req.Email).Msgf("%s failed to register user", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Int("user_id", user.ID).Str("email", user.Email).Msgf("%s user registered successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusCreated, ginext.H{
		"user": user,
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
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := hd.sv.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		lg.Warn().Err(err).Str("email", req.Email).Msgf("%s failed to login user", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Str("email", req.Email).Msgf("%s user logged in successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, ginext.H{
		"token":        accessToken,
		"refreshToken": refreshToken,
	})
}

// Refresh handles token refresh
func (hd *AuthHandler) Refresh(c *gin.Context) {
	lg := hd.lg.With().Str("handler", "Refresh").Logger()

	if !strings.Contains(c.ContentType(), "application/json") {
		lg.Warn().Str("content-type", c.ContentType()).Int("status", http.StatusBadRequest).Msgf("%s invalid content-type", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": pkgErrors.ErrContentTypeAJ.Error()})
		return
	}

	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	newAccessToken, newRefreshToken, err := hd.sv.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		lg.Warn().Err(err).Msgf("%s failed to refresh tokens", pkgConst.Warn)
		c.JSON(http.StatusUnauthorized, ginext.H{"error": err.Error()})
		return
	}

	lg.Debug().Msgf("%s tokens refreshed successfully", pkgConst.OpSuccess)
	c.JSON(http.StatusOK, ginext.H{
		"token":        newAccessToken,
		"refreshToken": newRefreshToken,
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
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s failed to get user", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
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

	var req struct {
		Name                  string  `json:"name"`
		TelegramUsername      *string `json:"telegramUsername"`
		TelegramNotifications *bool   `json:"telegramNotifications"`
		EmailNotifications    *bool   `json:"emailNotifications"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Warn().Err(err).Msgf("%s error bind json", pkgConst.Warn)
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	user := &model.User{
		ID:               userIDInt,
		Name:             req.Name,
		TelegramUsername: req.TelegramUsername,
	}

	err := hd.sv.UpdateUser(c.Request.Context(), user)
	if err != nil {
		lg.Warn().Err(err).Int("user_id", userIDInt).Msgf("%s failed to update user", pkgConst.Warn)
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
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
