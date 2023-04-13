package role

import (
	"context"
	"health/models"
	"health/shared/types"
)

type UseCase interface {
	GetRoles(ctx context.Context) ([]*models.Role, *types.Error)
	GetRoleByID(ctx context.Context, id string) (*models.Role, *types.Error)
}
