package role

import (
	"fmt"
	"health/shared/types"

	"github.com/go-playground/validator"
)

type RoleInput struct {
	ID    string `json:"_id,omitempty"`
	Email string `json:"email"`

	RoleID string `json:"roleId,omitempty" validate:"required"`
}

func ValidateRoleInput(inp *RoleInput) *types.Error {
	validate := validator.New()
	err := validate.Struct(inp)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &types.Error{
					Message: fmt.Sprintf("%s is required", err.Field()),
					Field:   "id",
					Tag:     "role",
				}
			}
		}
	}

	return nil
}