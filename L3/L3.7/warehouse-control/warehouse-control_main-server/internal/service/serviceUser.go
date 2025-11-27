package service

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// IUserRp interface for user repository
type IUserRp interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id int) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
}

// UserService service for working with users
type UserService struct {
	rp IUserRp
}

// NewUserService creates a new UserService
func NewUserService(rp IUserRp) *UserService {
	return &UserService{rp: rp}
}

// Create creates a new user
func (sv *UserService) Create(ctx context.Context, user *model.User, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = string(hashedPassword)

	createdUser, err := sv.rp.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// GetByID returns a user by ID
func (sv *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	return sv.rp.GetByID(id)
}

// GetByUsername returns a user by username
func (sv *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return sv.rp.GetByUsername(username)
}

// Update updates a user
func (sv *UserService) Update(ctx context.Context, user *model.User) error {
	return sv.rp.Update(user)
}

// Delete deletes a user by ID
func (sv *UserService) Delete(ctx context.Context, id int) error {
	return sv.rp.Delete(id)
}

// DeleteUser deletes a user by ID (alias for Delete for consistency)
func (sv *UserService) DeleteUser(ctx context.Context, id int) error {
	return sv.rp.Delete(id)
}
