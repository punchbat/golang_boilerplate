package role

import (
	"health/shared/types"
)

var (
	ErrCantFindRole = types.Error{
		Message: "Cant find role",
		Field:   "id",
		Tag:     "role",
	}
	ErrUserIsUnauthorized = types.Error{
		Message: "User is unauthorized",
		Field:   "id",
		Tag:     "role",
	}
)
