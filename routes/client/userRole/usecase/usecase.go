package usecase

import (
	"context"
	"health/models"

	"health/routes/client/auth"
	"health/routes/client/role"
	"health/routes/client/userRole"
	"health/shared/types"
	"health/shared/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UseCase struct {
	repo     userRole.Repository
	roleRepo role.Repository
	userRepo auth.Repository
}

func NewUseCase(repo userRole.Repository, roleRepo role.Repository, userRepo auth.Repository) *UseCase {
	return &UseCase{
		repo:     repo,
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

func (a *UseCase) AddRole(ctx context.Context, inp *userRole.RoleInput) *types.Error {
	user, err := a.userRepo.GetUserById(ctx, inp.ID)
	if err != nil {
		return &auth.ErrUserNotFound
	}

	roleEntity, err := a.roleRepo.GetRoleByID(ctx, inp.RoleID)
	if err != nil {
		return &userRole.ErrCantFindRole
	}

	// * Существует ли уже роль?
	userRoleEntities, err := a.repo.GetUserRoleByIDs(ctx, user.UserRoleIDs)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "user-role-ids",
			Tag:     "user-role",
		}
	}
	for _, userRoleEntity := range userRoleEntities {
		if userRoleEntity.RoleID == roleEntity.ID {
			return &userRole.ErrRoleIsExist
		}
	}

	userRoleID := primitive.NewObjectID()

	// Создаем юзер.роль
	userRole := models.UserRole{
		ID:     userRoleID.Hex(),
		UserID: user.ID,
		RoleID: roleEntity.ID,
		Status: models.UserRoleStatusPending,
	}
	err = a.repo.CreateUserRole(ctx, &userRole)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "create-user-role",
			Tag:     "user-role",
		}
	}

	user.UserRoleIDs = append(user.UserRoleIDs, userRoleID.Hex())
	if err := a.userRepo.UpdateUser(ctx, user); err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "update-user",
			Tag:     "user-role",
		}
	}

	return nil
}

func (a *UseCase) RemoveRole(ctx context.Context, inp *userRole.RoleInput) *types.Error {
	user, err := a.userRepo.GetUserById(ctx, inp.ID)
	if err != nil {
		return &auth.ErrUserNotFound
	}

	roleEntity, err := a.roleRepo.GetRoleByID(ctx, inp.RoleID)
	if err != nil {
		return &userRole.ErrCantFindRole
	}

	userRoleEntities, err := a.repo.GetUserRoleByIDs(ctx, user.UserRoleIDs)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "remove-role",
			Tag:     "user-role",
		}
	}

	// * Существует ли уже роль?
	isExist := false

	for _, userRoleEntity := range userRoleEntities {
		if userRoleEntity.RoleID == roleEntity.ID {
			isExist = true

			for index, userRoleID := range user.UserRoleIDs {
				if userRoleEntity.ID == userRoleID {
					user.UserRoleIDs = utils.RemoveElementByIndex(user.UserRoleIDs, index)
				}
			}
		}
	}

	if !isExist {
		return &userRole.ErrRoleIsNotExist
	}

	if err := a.userRepo.UpdateUser(ctx, user); err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "remove-role",
			Tag:     "user-role",
		}
	}

	return nil
}