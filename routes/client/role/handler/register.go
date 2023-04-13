package roleHandler

import (
	authRepository "health/routes/client/auth/repository"
	"health/routes/client/role/repository"
	"health/routes/client/role/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, db *mongo.Database) (gin.HandlerFunc, gin.HandlerFunc, gin.HandlerFunc) {
	// Создаем repository, все взаимодействия с db в ней
	repo := repository.NewRepository(db)
	authRepository := authRepository.NewRepository(db)

	// Создаем usecase, вся бизнес-логика в нем
	uc := usecase.NewUseCase(
		repo,
		authRepository,
	)

	// Create the middleware instance
	mUser := NewMiddlewareUser(uc)
	mSpecialist := NewMiddlewareSpecialist(uc)
	mMinion := NewMiddlewareMinion(uc)

	// Create the handler
	h := NewHandler(uc)

	// Create the endpoints
	endpoints := router.Group("/role/v1")
	{
		endpoints.GET("/list", h.GetRoles)
	}

	return mUser, mSpecialist, mMinion
}
