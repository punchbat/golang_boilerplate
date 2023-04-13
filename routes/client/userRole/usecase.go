package userRole

import (
	"context"
	"health/shared/types"
)

type UseCase interface {
	AddRole(ctx context.Context, inp *RoleInput) *types.Error
	RemoveRole(ctx context.Context, inp *RoleInput) *types.Error
}
