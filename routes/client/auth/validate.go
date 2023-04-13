package auth

import (
	"fmt"
	"health/models"
	"health/shared/types"
	"strconv"
	"time"

	"github.com/go-playground/validator"
)

type SignUpInput struct {
	Email           string `json:"email"            validate:"required,email"`
	Password        string `json:"password"         validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=@!?"`
	PasswordConfirm string `json:"passwordConfirm"  validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=@!?"`
}

func ValidateSignUpInput(inp *SignUpInput) *types.Error {
	validate := validator.New()
	err := validate.Struct(inp)

	if inp.Password != inp.PasswordConfirm {
		return &ErrPasswordNotEqual
	}

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &types.Error{
					Message: fmt.Sprintf("%s is required", err.Field()),
					Field:   "email",
					Tag:     "auth",
				}
			case "email":
				return &types.Error{
					Message: fmt.Sprintf("%s is not a valid email", inp.Email),
					Field:   "email",
					Tag:     "auth",
				}
			case "min":
				return &types.Error{
					Message: fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			case "containsany":
				return &types.Error{
					Message: fmt.Sprintf("%s should contain at least one %s character", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			}
		}
	}

	return nil
}

type SendVerifyCodeInput struct {
	Email    string `json:"email"        validate:"required,email"`
	Password string `json:"password"     validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789"`
}

func ValidateSendVerifyCodeInput(inp *SendVerifyCodeInput) *types.Error {
	validate := validator.New()
	err := validate.Struct(inp)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &types.Error{
					Message: fmt.Sprintf("%s is required", err.Field()),
					Field:   "email",
					Tag:     "auth",
				}
			case "min":
				return &types.Error{
					Message: fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			case "containsany":
				return &types.Error{
					Message: fmt.Sprintf("%s should contain at least one %s character", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			}
		}
	}

	return nil
}

type CheckVerifyCodeInput struct {
	Email      string `json:"email"        validate:"required,email"`
	Password   string `json:"password"     validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789"`
	VerifyCode string `json:"verifyCode"   validate:"required,len=6"`
}

func ValidateCheckVerifyCodeInput(inp *CheckVerifyCodeInput) *types.Error {
	validate := validator.New()
	err := validate.Struct(inp)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &types.Error{
					Message: fmt.Sprintf("%s is required", err.Field()),
					Field:   "email",
					Tag:     "auth",
				}
			case "min":
				return &types.Error{
					Message: fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			case "containsany":
				return &types.Error{
					Message: fmt.Sprintf("%s should contain at least one %s character", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			case "len":
				return &types.Error{
					Message: fmt.Sprintf("%s must be %s characters", err.Field(), err.Param()),
					Field:   "verify-code",
					Tag:     "auth",
				}
			}
		}
	}

	return nil
}

type SignInInput struct {
	Email    string `json:"email"        validate:"required,email"`
	Password string `json:"password"     validate:"required,min=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789"`
}

func ValidateSignInInput(inp *SignInInput) *types.Error {
	validate := validator.New()
	err := validate.Struct(inp)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &types.Error{
					Message: fmt.Sprintf("%s is required", err.Field()),
					Field:   "email",
					Tag:     "auth",
				}
			case "email":
				return &types.Error{
					Message: fmt.Sprintf("%s is not a valid email", inp.Email),
					Field:   "email",
					Tag:     "auth",
				}
			case "min":
				return &types.Error{
					Message: fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			case "containsany":
				return &types.Error{
					Message: fmt.Sprintf("%s should contain at least one %s character", err.Field(), err.Param()),
					Field:   "password",
					Tag:     "auth",
				}
			}
		}
	}

	return nil
}

type Address struct {
	Country     string `json:"country" validate:"required"`
	City        string `json:"city" validate:"required"`
	Street      string
	HouseNumber string
}

type UpdateProfileInput struct {
	ID    string `json:"_id,omitempty"`
	Email string `json:"email"`

	IIN      int           `json:"IIN"             validate:"required,number,IIN_custom_validation"`
	Name     string        `json:"name"`
	Surname  string        `json:"surname"`
	Birthday time.Time     `json:"birthday"        validate:"required,birthday_custom_validation"`
	Gender   models.Gender `json:"gender"          validate:"required,gender_custom_validation"`
	Address  Address       `json:"address"         validate:"required,dive"`

	RoleIDs []string `json:"roleIds,omitempty" validate:"required"`
}

func ValidateUpdateProfileInput(inp *UpdateProfileInput) *types.Error {
	validate := validator.New()

	err := validate.RegisterValidation("IIN_custom_validation", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().(int)

		return len(strconv.Itoa(value)) == 12
	})
	if err != nil {
		return &types.Error{
			Message: fmt.Sprintf("IIN must be 12 characters"),
			Field:   "IIN",
			Tag:     "auth",
		}
	}

	err = validate.RegisterValidation("birthday_custom_validation", func(fl validator.FieldLevel) bool {
		age := time.Since(fl.Field().Interface().(time.Time)).Hours() / 24 / 365.25

		return age > 0 && age < 130
	})
	if err != nil {
		return &types.Error{
			Message: fmt.Sprintf("Age must be between 0 and 130"),
			Field:   "Age",
			Tag:     "auth",
		}
	}

	err = validate.RegisterValidation("gender_custom_validation", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().(models.Gender)

		return value.String() != "unknown"
	})
	if err != nil {
		return &types.Error{
			Message: fmt.Sprintf("Gender is unknown"),
			Field:   "gender",
			Tag:     "auth",
		}
	}

	err = validate.Struct(inp)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &types.Error{
					Message: fmt.Sprintf("%s is required", err.Field()),
					Field:   "email",
					Tag:     "auth",
				}
			}
		}
	}

	// validate the address
	err = validate.Struct(inp.Address)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return &types.Error{
					Message: fmt.Sprintf("%s is required", err.Field()),
					Field:   fmt.Sprintf("address.%s", err.Field()),
					Tag:     "auth",
				}
			}
		}
	}

	return nil
}

type GetProfileInput struct {
	ID    string `json:"_id,omitempty"`
	Email string `json:"email"`
}