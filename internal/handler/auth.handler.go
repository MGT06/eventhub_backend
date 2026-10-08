package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/MGT06/EventHub_Backend.git/docs"
	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	errorTemplate "github.com/MGT06/EventHub_Backend.git/internal/error"
	"github.com/MGT06/EventHub_Backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthHandler struct {
	as *service.AuthService
}

func NewAuthHandler(as *service.AuthService) *AuthHandler {
	return &AuthHandler{
		as: as,
	}
}

// Register
//
// @Summary			Register
// @Description		Register Account
// @Tags			auth
// @Accept			json
// @Produce			json
// @Router			/auth/register	[post]
// @Param			newAccount	body	dto.Register	true	"Register"
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrorResponse
// @Failure			409		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (a *AuthHandler) Register(ctx *gin.Context) {
	var newAccount dto.Register
	if err := ctx.ShouldBindWith(&newAccount, binding.JSON); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	if err := a.as.Register(ctx.Request.Context(), newAccount); err != nil {
		log.Println(err)
		if errors.Is(err, errorTemplate.ErrInvalidInputs) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: errorTemplate.ErrInvalidInputs.Error(),
			})
			return
		}
		if errors.Is(err, errorTemplate.ErrEmailExist) {
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{
				Success: false,
				Message: errorTemplate.ErrEmailExist.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Account created successfully",
	})
}

// Login
//
// @Summary			Login
// @Description		Login using email and password
// @Tags			auth
// @Accept			json
// @Produce			json
// @Router			/auth/login	[post]
// @Param			account	body	dto.Login	true	"Login"
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrorResponse
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (a *AuthHandler) Login(ctx *gin.Context) {
	var account dto.Login
	if err := ctx.ShouldBindWith(&account, binding.JSON); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Invalid email format",
		})
		return
	}

	token, err := a.as.Login(ctx.Request.Context(), account)
	if err != nil {
		log.Println(err)
		if errors.Is(err, errorTemplate.ErrInvalidInputs) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "Please fill in all required fields",
			})
			return
		}
		if errors.Is(err, errorTemplate.ErrEmailPasswordIncorrect) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"token": token,
		},
		Message: "Login Success",
	})
}


// Logout
//
// @Summary			Logut
// @Description		User Logout
// @Tags			auth
// @Accept			json
// @Produce			json
// @Router			/auth/logout	[post]
// @Security 		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrorResponse
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (a *AuthHandler) Logout(ctx *gin.Context) {
	idUser, exist := ctx.Get("idUser")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	jwtID, exist := ctx.Get("jti")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	exp, exist := ctx.Get("exp")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	expTime, ok := exp.(time.Time)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	ttl := time.Until(expTime)

	if err := a.as.Logout(ctx, idUser.(int), jwtID.(string), ttl); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, dto.Response{
		Success: true,
		Message: "logout success",
	})
}

// Forgot Password
//
// @Summary			Forgot Password
// @Tags			auth
// @Accept			json
// @Produce			json
// @Router			/auth/forgot-password	[post]
// @Param			body	body	dto.RequestForgotPassword	true	"forgot password"
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (a *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var body dto.RequestForgotPassword
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	res, err := a.as.ForgotPassword(ctx, body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    res,
		Message: "Email exists, please reset your password",
	})
}

// Reset Password
//
// @Summary			Reset Password
// @Tags			auth
// @Accept			json
// @Produce			json
// @Router			/auth/reset-password	[post]
// @Param			body	body	dto.ResetPassword	true	"reset password"
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (a *AuthHandler) ResetPassword(ctx *gin.Context) {
	var body dto.ResetPassword
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	if err := a.as.ResetPassword(ctx, body); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Success Reset Password",
	})
}
