package role

import (
	"context"
	"health/models"
)

type Repository interface {
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	GetRoleByID(ctx context.Context, id string) (*models.Role, error)
	GetRoleByIDs(ctx context.Context, id []string) ([]*models.Role, error)
	GetRoles(ctx context.Context) ([]*models.Role, error)
}
