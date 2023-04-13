package usecase

import (
	"context"
	"health/models"

	"health/routes/client/auth"
	"health/routes/client/role"
	"health/shared/types"
)

type UseCase struct {
	repoRole role.Repository
	repoUser auth.Repository
}

func NewUseCase(repoRole role.Repository, repoUser auth.Repository) *UseCase {
	return &UseCase{
		repoRole: repoRole,
		repoUser: repoUser,
	}
}

func (a *UseCase) GetRoles(ctx context.Context) ([]*models.Role, *types.Error) {
	roles, err := a.repoRole.GetRoles(ctx)

	if err != nil {
		return nil, &types.Error{
			Message: err.Error(),
			Field:   "roles",
			Tag:     "auth",
		}
	}

	return roles, nil
}

func (a *UseCase) GetRoleByID(ctx context.Context, id string) (*models.Role, *types.Error) {
	role, err := a.repoRole.GetRoleByID(ctx, id)

	if err != nil {
		return nil, &types.Error{
			Message: err.Error(),
			Field:   "id",
			Tag:     "auth-role",
		}
	}

	return role, nil
}