package service

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/transport/trhttp/handler/authHandler"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/transport/trhttp/handler/middleware"
)

// Service structure that combines all services
type Service struct {
	User *UserService
	Auth *AuthService
	Item *ItemService
}

// New creates a new Service structure
func New(
	cfg *Config,
	rp *repository.Repository,
	rs *pkgRetry.Retry,
) *Service {
	userService := NewUserService(rp.User())
	authService := NewAuthService(cfg, rp.User(), rp.RefreshToken())
	itemService := NewItemService(rp.Item(), rp.ItemHistory())

	return &Service{
		User: userService,
		Auth: authService,
		Item: itemService,
	}
}

// MiddlewareService returns the auth service for middleware
func (sv *Service) MiddlewareService() middleware.ISvForAuthHandler {
	return sv.Auth
}

// AuthService returns the auth service
func (sv *Service) AuthService() authHandler.ISvForAuthHandler {
	return sv.Auth
}

// ItemService returns the item service
func (sv *Service) ItemService() IItemService {
	return sv.Item
}
