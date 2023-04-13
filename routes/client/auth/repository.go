package auth

import (
	"context"
	"health/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}
