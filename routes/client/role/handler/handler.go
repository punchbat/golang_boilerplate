package roleHandler

import (
	"health/routes/client/role"
	"health/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase role.UseCase
}

func NewHandler(useCase role.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetRoles(c *gin.Context) {
	roles, err := h.useCase.GetRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})
		return
	}

	c.JSON(http.StatusOK, types.GoodResponse{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"roles": roles,
		},
	})
}
