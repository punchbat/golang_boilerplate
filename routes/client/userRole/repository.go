package userRole

import (
	"context"
	"health/models"
)

type Repository interface {
	CreateUserRole(ctx context.Context, userRole *models.UserRole) error
	GetUserRoleByID(ctx context.Context, id string) (*models.UserRole, error)
	GetUserRoleByIDs(ctx context.Context, ids []string) ([]*models.UserRole, error)
	DeleteUserRoleByID(ctx context.Context, id string) error
}
