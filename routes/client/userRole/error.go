package userRole

import (
	"health/shared/types"
)

var (
	ErrCantFindRole = types.Error{
		Message: "Can`t find role",
		Field:   "role_id",
		Tag:     "user-role",
	}
	ErrRoleIsExist = types.Error{
		Message: "This role already exists for the user",
		Field:   "role_id",
		Tag:     "user-role",
	}
	ErrRoleIsNotExist = types.Error{
		Message: "This role does not exist for the user",
		Field:   "role_id",
		Tag:     "user-role",
	}
)
