package authHandler

import (
	"health/routes/client/auth"
	"health/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	usecase auth.UseCase
}

func NewMiddleware(usecase auth.UseCase) gin.HandlerFunc {
	return (&Middleware{
		usecase: usecase,
	}).Handle
}

func (m *Middleware) Handle(c *gin.Context) {
	tokenFromHeader := c.GetHeader("Authorization")

	if tokenFromHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.usecase.ParseToken(c.Request.Context(), tokenFromHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &types.BadResponse{
			Code:  http.StatusUnauthorized,
			Error: err,
		})

		return
	}

	c.Set(auth.CtxUserKey, user)
}
