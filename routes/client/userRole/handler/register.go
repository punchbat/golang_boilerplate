package userRoleHandler

import (
	authRepository "health/routes/client/auth/repository"
	roleRepository "health/routes/client/role/repository"
	"health/routes/client/userRole/repository"
	"health/routes/client/userRole/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, db *mongo.Database) {
	// Создаем repository, все взаимодействия с db в ней
	repo := repository.NewRepository(db)
	roleRepository := roleRepository.NewRepository(db)
	authRepository := authRepository.NewRepository(db)

	// Создаем usecase, вся бизнес-логика в нем
	uc := usecase.NewUseCase(
		repo,
		roleRepository,
		authRepository,
	)

	// Create the handler
	h := NewHandler(uc)

	// Create the endpoints
	endpoints := router.Group("/user-role/v1")
	{
		endpoints.POST("/add", authMiddleware, h.AddRole)
		endpoints.POST("/remove", authMiddleware, h.RemoveRole)
	}
}
