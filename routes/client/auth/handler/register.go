package authHandler

import (
	"health/routes/client/auth/repository"
	"health/routes/client/auth/usecase"
	roleRepository "health/routes/client/role/repository"
	userRoleRepository "health/routes/client/userRole/repository"
	"health/services/email"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterHTTPEndpoints(router *gin.Engine, db *mongo.Database, mailer *email.Mailer) gin.HandlerFunc {
	// Создаем repository, все взаимодействия с db в ней
	repo := repository.NewRepository(db)
	roleRepository := roleRepository.NewRepository(db)
	userRoleRepository := userRoleRepository.NewRepository(db)

	// Создаем usecase, вся бизнес-логика в нем
	uc := usecase.NewUseCase(
		repo,
		roleRepository,
		userRoleRepository,

		mailer,
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetDuration("auth.token_ttl"),
	)

	// Create the middleware instance
	m := NewMiddleware(uc)

	// Create the handler
	h := NewHandler(uc)

	// Create the endpoints
	endpoints := router.Group("/auth/v1")
	{
		endpoints.POST("/sign-up", h.SignUp)
		endpoints.POST("/send-verify-code", h.SendVerifyCode)
		endpoints.POST("/check-verify-code", h.CheckVerifyCode)
		endpoints.POST("/sign-in", h.SignIn)

		// * проверяем на наличие аутентификации
		endpoints.POST("/update-profile", m, h.UpdateProfile)
		endpoints.GET("/get-profile", m, h.GetProfile)
	}

	return m
}
