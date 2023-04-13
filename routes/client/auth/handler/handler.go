package authHandler

import (
	"health/models"
	"health/routes/client/auth"
	"health/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(auth.SignUpInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, types.BadResponse{
			Code: http.StatusBadRequest,
			Error: &types.Error{
				Message: err.Error(),
				Field:   "input data",
				Tag:     "auth",
			},
		})
		return
	}

	if err := auth.ValidateSignUpInput(inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})

		return
	}

	if err := h.useCase.SignUp(c.Request.Context(), inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SendVerifyCode(c *gin.Context) {
	inp := new(auth.SendVerifyCodeInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, types.BadResponse{
			Code: http.StatusBadRequest,
			Error: &types.Error{
				Message: err.Error(),
				Field:   "input data",
				Tag:     "auth",
			},
		})
		return
	}

	if err := auth.ValidateSendVerifyCodeInput(inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})

		return
	}

	if err := h.useCase.SendVerifyCode(c.Request.Context(), inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) CheckVerifyCode(c *gin.Context) {
	inp := new(auth.CheckVerifyCodeInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, types.BadResponse{
			Code: http.StatusBadRequest,
			Error: &types.Error{
				Message: err.Error(),
				Field:   "input data",
				Tag:     "auth",
			},
		})
		return
	}

	if err := auth.ValidateCheckVerifyCodeInput(inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})

		return
	}

	token, err := h.useCase.CheckVerifyCode(c.Request.Context(), inp)
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
			"token": token,
		},
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(auth.SignInInput)

	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, types.BadResponse{
			Code: http.StatusBadRequest,
			Error: &types.Error{
				Message: err.Error(),
				Field:   "input data",
				Tag:     "auth",
			},
		})
		return
	}

	if err := auth.ValidateSignInInput(inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})

		return
	}

	if err := h.useCase.SignIn(c.Request.Context(), inp); err != nil {
		c.JSON(http.StatusUnauthorized, types.BadResponse{
			Code:  http.StatusUnauthorized,
			Error: err,
		})

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetProfile(c *gin.Context) {
	inp := new(auth.GetProfileInput)

	// c токена вытаскиваем
	if user, exist := c.Get(auth.CtxUserKey); exist {
		inp.ID = user.(*models.User).ID
		inp.Email = user.(*models.User).Email
	}

	user, err := h.useCase.GetProfile(c.Request.Context(), inp)

	if err != nil {
		c.JSON(http.StatusUnauthorized, types.BadResponse{
			Code:  http.StatusUnauthorized,
			Error: err,
		})

		return
	}

	c.JSON(http.StatusOK, types.GoodResponse{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"user": user,
		},
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	inp := new(auth.UpdateProfileInput)

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
				Tag:     "auth",
			},
		})
		return
	}

	if err := auth.ValidateUpdateProfileInput(inp); err != nil {
		c.JSON(http.StatusNotAcceptable, types.BadResponse{
			Code:  http.StatusNotAcceptable,
			Error: err,
		})

		return
	}

	token, err := h.useCase.UpdateProfile(c.Request.Context(), inp)
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
			"token": token,
		},
	})
}