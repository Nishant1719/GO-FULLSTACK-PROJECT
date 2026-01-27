package users

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for user data operations
type Repository interface {
	// Create creates a new user in the database
	Create(ctx context.Context, user *User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*User, error)

	// List retrieves all users with optional pagination
	List(ctx context.Context, limit, offset int) ([]*User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *User) error

	// Delete deletes a user by their ID (soft delete by setting is_active to false)
	Delete(ctx context.Context, id uuid.UUID) error

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}
