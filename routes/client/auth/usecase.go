package auth

import (
	"context"
	"health/models"
	"health/shared/types"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, inp *SignUpInput) *types.Error
	SendVerifyCode(ctx context.Context, inp *SendVerifyCodeInput) *types.Error
	CheckVerifyCode(ctx context.Context, inp *CheckVerifyCodeInput) (string, *types.Error)

	SignIn(ctx context.Context, inp *SignInInput) *types.Error
	ParseToken(ctx context.Context, accessToken string) (*models.User, *types.Error)
	UpdateProfile(ctx context.Context, inp *UpdateProfileInput) (string, *types.Error)

	GetProfile(ctx context.Context, inp *GetProfileInput) (*models.User, *types.Error)
}
