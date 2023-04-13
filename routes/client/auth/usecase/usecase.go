package usecase

import (
	"context"
	"errors"
	"fmt"
	"health/models"
	"health/services/email"
	service_email "health/services/email"
	"math/rand"
	"strconv"
	"time"

	"health/routes/client/auth"
	"health/routes/client/role"
	"health/routes/client/userRole"
	"health/shared/types"
	"health/shared/utils"

	"github.com/dgrijalva/jwt-go/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

type UseCase struct {
	repo           auth.Repository
	roleRepo       role.Repository
	userRoleRepo   userRole.Repository
	mailer         *email.Mailer
	signingKey     []byte
	expireDuration time.Duration
}

func NewUseCase(
	repo auth.Repository,
	roleRepo role.Repository,
	userRoleRepo userRole.Repository,

	mailer *service_email.Mailer,
	signingKey []byte,
	tokenTTLHours time.Duration) *UseCase {
	return &UseCase{
		repo:         repo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,

		mailer:         mailer,
		signingKey:     signingKey,
		expireDuration: time.Hour * tokenTTLHours,
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePasswordHash(password1, password2 string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	if err != nil {
		return false, errors.New("invalid password")
	}
	return true, nil
}

func (a *UseCase) MakeClearUser(ctx context.Context, user *models.User) (*models.User, error) {
	utils.RemoveKeyFromStruct(user, "Password")
	utils.RemoveKeyFromStruct(user, "PasswordConfirm")

	userRoles, err := a.userRoleRepo.GetUserRoleByIDs(ctx, user.UserRoleIDs)
	if err != nil {
		return nil, err
	}

	roles, err := a.roleRepo.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	var resUserRoles []*models.UserRoleWithRole
	for _, userRole := range userRoles {
		for _, role := range roles {
			if userRole.RoleID == role.ID {
				resUserRoles = append(resUserRoles, &models.UserRoleWithRole{
					Name:      role.Name,
					Status:    userRole.Status,
					IsDefault: role.IsDefault,
				})

				break
			}
		}
	}

	user.UserRoles = resUserRoles

	return user, nil
}

func (a *UseCase) SignUp(ctx context.Context, inp *auth.SignUpInput) *types.Error {
	user, err := a.repo.GetUserByEmail(ctx, inp.Email)
	if err == nil {
		return &auth.ErrUserIsExist
	}

	hashPassword, err := HashPassword(inp.Password)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "sign-up",
			Tag:     "auth",
		}
	}

	hashPasswordConfirm, err := HashPassword(inp.PasswordConfirm)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "sign-up",
			Tag:     "auth",
		}
	}

	// Находим роль обычного юзера
	role, err := a.roleRepo.GetRoleByName(ctx, string(models.RoleNameUser))
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "sign-up",
			Tag:     "auth",
		}
	}

	userRoleId := primitive.NewObjectID()
	userId := primitive.NewObjectID()

	// Создаем юзер.роль
	userRole := models.UserRole{
		ID:     userRoleId.Hex(),
		UserID: userId.Hex(),
		RoleID: role.ID,
		Status: models.UserRoleStatusApproved,
	}
	err = a.userRoleRepo.CreateUserRole(ctx, &userRole)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "sign-up",
			Tag:     "auth",
		}
	}

	// подумать
	user = &models.User{
		ID:                   userId.Hex(),
		UserRoleIDs:          []string{userRoleId.Hex()},
		Email:                inp.Email,
		Password:             hashPassword,
		PasswordConfirm:      hashPasswordConfirm,
		Verified:             false,
		FinishedRegistration: false,
	}

	if err := a.repo.CreateUser(ctx, user); err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "sign-up",
			Tag:     "auth",
		}
	}

	return nil
}

type EmailContent struct {
	VerifyCode string
}

func (a *UseCase) SendVerifyCode(ctx context.Context, inp *auth.SendVerifyCodeInput) *types.Error {
	user, err := a.repo.GetUserByEmail(ctx, inp.Email)
	if err != nil {
		return &auth.ErrUserNotFound
	}

	isEqual, err := ComparePasswordHash(user.Password, inp.Password)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "send-verify-code",
			Tag:     "auth",
		}
	}
	if !isEqual {
		return &auth.ErrEmailOrPassword
	}

	minVal := 100000
	maxVal := 999999
	randomCode := rand.Intn(maxVal-minVal+1) + minVal

	user.VerifyCode = strconv.Itoa(randomCode)
	if err := a.repo.UpdateUser(ctx, user); err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "send-verify-code",
			Tag:     "auth",
		}
	}

	// Отправляем письмо
	emailMessage := service_email.Message{
		Subject:      "Your service: verify code",
		To:           []string{inp.Email},
		TemplateName: "VerifyCode",
		Content: EmailContent{
			VerifyCode: user.VerifyCode,
		},
	}
	a.mailer.Send(&emailMessage)

	return nil
}

func (a *UseCase) CheckVerifyCode(ctx context.Context, inp *auth.CheckVerifyCodeInput) (string, *types.Error) {
	user, err := a.repo.GetUserByEmail(ctx, inp.Email)
	if err != nil {
		return "", &auth.ErrUserNotFound
	}

	isEqual, err := ComparePasswordHash(user.Password, inp.Password)
	if err != nil {
		return "", &types.Error{
			Message: err.Error(),
			Field:   "check-verify-code",
			Tag:     "auth",
		}
	}
	if !isEqual {
		return "", &auth.ErrEmailOrPassword
	}

	if user.VerifyCode != inp.VerifyCode {
		return "", &auth.ErrVerifyCodeNotMatch
	}

	user.Verified = true
	user.VerifyCode = ""

	if err := a.repo.UpdateUser(ctx, user); err != nil {
		return "", &types.Error{
			Message: err.Error(),
			Field:   "sign-up",
			Tag:     "auth",
		}
	}

	return a.GetToken(ctx, user)
}

func (a *UseCase) SignIn(ctx context.Context, inp *auth.SignInInput) *types.Error {
	user, err := a.repo.GetUserByEmail(ctx, inp.Email)
	if err != nil {
		return &auth.ErrUserNotFound
	}

	isEqual, err := ComparePasswordHash(user.Password, inp.Password)
	if err != nil {
		return &types.Error{
			Message: err.Error(),
			Field:   "send-verify-code",
			Tag:     "auth",
		}
	}
	if !isEqual {
		return &auth.ErrEmailOrPassword
	}

	// Если юзер не verified, то он не может зайти
	if !user.Verified {
		return &auth.ErrUserIsUnauthorized
	}

	return nil
}

func (a *UseCase) GetToken(ctx context.Context, user *models.User) (string, *types.Error) {
	user, err := a.MakeClearUser(ctx, user)
	if err != nil {
		return "", &types.Error{
			Message: err.Error(),
			Field:   "create-token",
			Tag:     "auth",
		}
	}

	claims := AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	completeSignedToken, err := token.SignedString(a.signingKey)

	if err != nil {
		return "", &types.Error{
			Message: err.Error(),
			Field:   "create-token",
			Tag:     "auth",
		}
	}

	return completeSignedToken, nil
}

func (a *UseCase) ParseToken(ctx context.Context, accessToken string) (*models.User, *types.Error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return nil, &types.Error{
			Message: err.Error(),
			Field:   "parse-token",
			Tag:     "auth",
		}
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, &auth.ErrInvalidAccessToken
}

func (a *UseCase) GetProfile(ctx context.Context, inp *auth.GetProfileInput) (*models.User, *types.Error) {
	user, err := a.repo.GetUserById(ctx, inp.ID)

	if err != nil {
		return nil, &auth.ErrUserNotFound
	}

	user, err = a.MakeClearUser(ctx, user)
	if err != nil {
		return nil, &types.Error{
			Message: err.Error(),
			Field:   "get-profile",
			Tag:     "auth",
		}
	}

	return user, nil
}

func (a *UseCase) UpdateProfile(ctx context.Context, inp *auth.UpdateProfileInput) (string, *types.Error) {
	// * Получаем юзера
	user, err := a.repo.GetUserById(ctx, inp.ID)
	if err != nil {
		return "", &auth.ErrUserNotFound
	}

	// * Закончили регистрацию до конца
	user.FinishedRegistration = true

	user.IIN = inp.IIN
	user.Name = inp.Name
	user.Surname = inp.Surname
	user.Birthday = inp.Birthday
	user.Gender = inp.Gender
	user.Address = models.Address(inp.Address)

	// * Находим роль юзера
	roleUser, err := a.roleRepo.GetRoleByName(ctx, string(models.RoleNameUser))
	if err != nil {
		return "", &types.Error{
			Message: err.Error(),
			Field:   "update-profile",
			Tag:     "auth",
		}
	}

	var newUserRoleIDs []string
	// * Все текущие юзер-роли у юзера
	currUserRoles, err := a.userRoleRepo.GetUserRoleByIDs(ctx, user.UserRoleIDs)
	if err != nil {
		return "", &types.Error{
			Message: err.Error(),
			Field:   "user-role-ids",
			Tag:     "auth",
		}
	}

	for _, currUserRole := range currUserRoles {
		if currUserRole.RoleID == roleUser.ID {
			newUserRoleIDs = []string{currUserRole.ID}
		}
	}

	// * Находим все роли из запроса
	inpRoles, err := a.roleRepo.GetRoleByIDs(ctx, inp.RoleIDs)
	if err != nil {
		return "", &types.Error{
			Message: err.Error(),
			Field:   "role-ids",
			Tag:     "auth",
		}
	}

	for _, inpRole := range inpRoles {
		var userRoleID string

		for _, currUserRole := range currUserRoles {
			if currUserRole.RoleID == inpRole.ID {
				userRoleID = currUserRole.ID
			}

			if roleUser.ID == inpRole.ID {
				continue
			}
		}

		// * Сохраняем старую роль
		if userRoleID != "" {
			newUserRoleIDs = append(newUserRoleIDs, userRoleID)
		} else {
			// * Создаем новую роль
			userRoleId := primitive.NewObjectID()

			// Создаем юзер.роль
			userRole := models.UserRole{
				ID:     userRoleId.Hex(),
				UserID: user.ID,
				RoleID: inpRole.ID,
				// * У всех новых ролей, статус будет запрошенным
				Status:    models.UserRoleStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err = a.userRoleRepo.CreateUserRole(ctx, &userRole)
			if err != nil {
				return "", &types.Error{
					Message: err.Error(),
					Field:   "update-profile",
					Tag:     "auth",
				}
			}

			newUserRoleIDs = append(newUserRoleIDs, userRoleId.Hex())
		}
	}

	for _, currUserRole := range currUserRoles {
		var isDeleteUserRole bool = true

		for _, newUserRoleID := range newUserRoleIDs {
			if currUserRole.ID == newUserRoleID {
				isDeleteUserRole = false
			}
		}

		if isDeleteUserRole {
			err = a.userRoleRepo.DeleteUserRoleByID(ctx, currUserRole.ID)
			if err != nil {
				return "", &auth.ErrCantDeleteUserRole
			}
		}
	}

	// * Обновляем юзера
	user.UserRoleIDs = newUserRoleIDs
	if err = a.repo.UpdateUser(ctx, user); err != nil {
		return "", &auth.ErrCantUpdateUser
	}

	return a.GetToken(ctx, user)
}