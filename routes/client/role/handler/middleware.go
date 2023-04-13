package roleHandler

import (
	"fmt"
	"health/models"
	"health/routes/client/auth"
	"health/routes/client/role"
	"health/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	roleName models.RoleName
	status   models.UserRoleStatus
}

func NewMiddlewareUser(usecase role.UseCase) gin.HandlerFunc {
	return (NewMiddleware(models.RoleNameUser, models.UserRoleStatusApproved)).HandleMiddleware
}

func NewMiddlewareSpecialist(usecase role.UseCase) gin.HandlerFunc {
	return (NewMiddleware(models.RoleNameSpecialist, models.UserRoleStatusApproved)).HandleMiddleware
}

func NewMiddlewareMinion(usecase role.UseCase) gin.HandlerFunc {
	return (NewMiddleware(models.RoleNameMinion, models.UserRoleStatusApproved)).HandleMiddleware
}

func NewMiddleware(roleName models.RoleName, status models.UserRoleStatus) *Middleware {
	return &Middleware{roleName: roleName, status: status}
}

func (m *Middleware) HandleMiddleware(c *gin.Context) {
	var userRoles []*models.UserRoleWithRole

	if user, exist := c.Get(auth.CtxUserKey); exist {
		userRoles = user.(*models.User).UserRoles
	}

	fmt.Println("HandleMiddleware", userRoles)

	for _, userRole := range userRoles {
		isRole, isApproved :=
			userRole.Name == m.roleName, userRole.Status == m.status

		if isRole && isApproved {
			c.Next()
			return
		}
	}

	c.JSON(http.StatusUnauthorized, types.BadResponse{
		Code:  http.StatusUnauthorized,
		Error: &role.ErrUserIsUnauthorized,
	})
	c.Abort()
}