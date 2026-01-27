package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	// User CRUD operations
	CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*UserResponse, error)
	GetUserByUsername(ctx context.Context, username string) (*UserResponse, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUsersCount(ctx context.Context) (int64, error)
}

type svc struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &svc{
		repo: repo,
	}
}

// CreateUser creates a new user with hashed password
func (s *svc) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	// Save to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByID retrieves a user by their ID
func (s *svc) GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByEmail retrieves a user by their email
func (s *svc) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByUsername retrieves a user by their username
func (s *svc) GetUserByUsername(ctx context.Context, username string) (*UserResponse, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// ListUsers retrieves all users with pagination
func (s *svc) ListUsers(ctx context.Context, limit, offset int) ([]*UserResponse, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		response := user.ToResponse()
		responses[i] = &response
	}

	return responses, nil
}

// UpdateUser updates an existing user
func (s *svc) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser soft deletes a user
func (s *svc) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetUsersCount returns the total number of active users
func (s *svc) GetUsersCount(ctx context.Context) (int64, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
