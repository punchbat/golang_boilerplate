package routes

import (
	"net/http"

	authHandler "health/routes/client/auth/handler"
	roleHandler "health/routes/client/role/handler"
	userRoleHandler "health/routes/client/userRole/handler"
	"health/services/email"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRoutes(router *gin.Engine, db *mongo.Database, mailer *email.Mailer) {
	// Пингуем сервер
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// * AUTH
	authMiddleware := authHandler.RegisterHTTPEndpoints(router, db, mailer)

	// * API endpoints
	api := router.Group("/api")

	// * ROLE
	roleMiddlewareUser, roleMiddlewareSpecialist, roleMiddlewareMinion :=
		roleHandler.RegisterHTTPEndpoints(api, authMiddleware, db)

	// * USER.ROLE
	userRoleHandler.RegisterHTTPEndpoints(api, authMiddleware, db)

	// * CHECK ROLE MIDDLEWARES
	api.GET("/check-user", authMiddleware, roleMiddlewareUser, func(c *gin.Context) {
		c.String(http.StatusOK, "check-user")
	})

	api.GET("/check-specialist", authMiddleware, roleMiddlewareSpecialist, func(c *gin.Context) {
		c.String(http.StatusOK, "check-specialist")
	})

	api.GET("/check-minion", authMiddleware, roleMiddlewareMinion, func(c *gin.Context) {
		c.String(http.StatusOK, "check-minion")
	})

	api.GET("/check", authMiddleware, func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})
}