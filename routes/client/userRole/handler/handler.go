package userRoleHandler

import (
	"health/models"
	"health/routes/client/auth"
	"health/routes/client/userRole"
	"health/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase userRole.UseCase
}

func NewHandler(useCase userRole.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) AddRole(c *gin.Context) {
	inp := new(userRole.RoleInput)

	// c токена вытаскиваем
	if user, exist := c.Get(auth.CtxUserKey); exist {
		inp.ID = user.(*models.User).ID
		inp.Email = user.(*models.User).Email
	}

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, types.BadResponse{
			Code: http.StatusBadRequest,
			Error: &types.Error{
				Message: err.Error(),
				Field:   "input data",
				Tag:     "user-role",
			},
		})
		return
	}

	if err := userRole.ValidateRoleInput(inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})

		return
	}

	if err := h.useCase.AddRole(c.Request.Context(), inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) RemoveRole(c *gin.Context) {
	inp := new(userRole.RoleInput)

	// c токена вытаскиваем
	if user, exist := c.Get(auth.CtxUserKey); exist {
		inp.ID = user.(*models.User).ID
		inp.Email = user.(*models.User).Email
	}

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, types.BadResponse{
			Code: http.StatusBadRequest,
			Error: &types.Error{
				Message: err.Error(),
				Field:   "input data",
				Tag:     "user-role",
			},
		})
		return
	}

	if err := userRole.ValidateRoleInput(inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})

		return
	}

	if err := h.useCase.RemoveRole(c.Request.Context(), inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})
		return
	}

	c.Status(http.StatusOK)
}