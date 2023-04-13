package auth

import (
	"health/shared/types"
)

var (
	ErrEmailOrPassword = types.Error{
		Message: "Email or Password is wrong",
		Field:   "email/password",
		Tag:     "auth",
	}
	ErrPasswordNotEqual = types.Error{
		Message: "Password not equal",
		Field:   "password",
		Tag:     "auth",
	}
	ErrUserIsUnauthorized = types.Error{
		Message: "User is unauthorized",
		Field:   "token",
		Tag:     "auth",
	}
	ErrUserNotFound = types.Error{
		Message: "User not found",
		Field:   "sign-in",
		Tag:     "auth",
	}
	ErrUserIsExist = types.Error{
		Message: "User is exist",
		Field:   "sign-up",
		Tag:     "auth",
	}
	ErrVerifyCodeNotMatch = types.Error{
		Message: "Verify code is invalid",
		Field:   "sign-up",
		Tag:     "auth",
	}
	ErrInvalidAccessToken = types.Error{
		Message: "Invalid access token",
		Field:   "parse token",
		Tag:     "auth",
	}
	ErrRoleIsExist = types.Error{
		Message: "This role already exists for the user",
		Field:   "role-ids",
		Tag:     "auth",
	}
	ErrCantFindUserRole = types.Error{
		Message: "Can`t find role",
		Field:   "role_id",
		Tag:     "user-role",
	}
	ErrCantDeleteUserRole = types.Error{
		Message: "Can`t delete role",
		Field:   "role_id",
		Tag:     "user-role",
	}
	ErrCantUpdateUser = types.Error{
		Message: "Can`t update user",
		Field:   "user_id",
		Tag:     "user",
	}
)
