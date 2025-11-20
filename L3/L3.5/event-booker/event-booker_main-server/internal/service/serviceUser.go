package service

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo interface for user repository
type UserRepo interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
}

// UserService service for working with users
type UserService struct {
	rp UserRepo
}

// NewUserService creates a new UserService
func NewUserService(rp UserRepo) *UserService {
	return &UserService{rp: rp}
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, user *model.User, password string) (*model.User, error) {
	// Hash the password
	_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Save the user
	createdUser, err := s.rp.Create(user)
	if err != nil {
		return nil, err
	}

	// Here you can add logic to save the hashed password separately
	// or add a password field to the user model

	return createdUser, nil
}

// GetByID returns a user by ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	return s.rp.GetByID(id)
}

// GetByEmail returns a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.rp.GetByEmail(email)
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, user *model.User) error {
	return s.rp.Update(user)
}

// Delete deletes a user by ID
func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.rp.Delete(id)
}
