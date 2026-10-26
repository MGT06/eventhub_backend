package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	errorTemplate "github.com/MGT06/EventHub_Backend.git/internal/error"
	"github.com/MGT06/EventHub_Backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserHandler struct {
	us *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		us: us,
	}
}

// Get User Profile
//
// @Summary			Get User Profile
// @Tags			user
// @Produce			json
// @Router			/user	[get]
// @Security 		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserHandler) GetProfileUser(ctx *gin.Context) {
	id, exists := ctx.Get("idUser")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	res, err := u.us.GetProfileUser(ctx.Request.Context(), id.(int))

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
		Message: "Success Get Profile",
	})

}

// Edit User
//
// @Summary			Update user info
// @Tags			user
// @Accept			mpfd
// @Produce			json
// @Router			/user/edit	[patch]
// @Security 		BearerToken
// @Param			name			formData	string	true	"update name user"
// @Param			bio				formData	string	false	"update bio user"
// @Param			user_location	formData	string	false	"update location user"
// @Param			position		formData	string	false	"update position user"
// @Param			avatar_url		formData	file	false	"update avatar user"
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserHandler) EditProfileUser(ctx *gin.Context) {
	id, exists := ctx.Get("idUser")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	var body dto.EditUserProfile
	if err := ctx.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	var filepath *string

	if body.Avatar != nil {
		ext := path.Ext(body.Avatar.Filename)
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
		default:
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "Unsupported image type",
			})
			return
		}

		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), body.Name, path.Ext(body.Avatar.Filename))
		savePath := path.Join("public", "img", "persons", filename)

		if err := ctx.SaveUploadedFile(body.Avatar, savePath); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Message: "A system error has occurred",
			})
			return
		}

		path := path.Join("img", "persons", filename)
		filepath = &path
	}

	res, err := u.us.EditProfileUser(ctx, body, id.(int), filepath)

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
		Data: res,
		Message: "Success Edit Profile",
	})

}

// Change Password
//
// @Summary			Change Password
// @Description		Change Password User
// @Tags			user
// @Accept			json
// @Produce			json
// @Router			/user/change-password	[post]
// @Param			body	body	dto.ChangePassword	true	"change password"
// @Security 		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrorResponse
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserHandler) ChangePassword(ctx *gin.Context) {
	idUser, exist := ctx.Get("idUser")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	var body dto.ChangePassword
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	if err := u.us.ChangePassword(ctx, idUser.(int), body); err != nil {
		if errors.Is(err, errorTemplate.ErrInvalidInputs) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "Please fill in all required fields",
			})
			return
		}
		if errors.Is(err, errorTemplate.ErrEmailPasswordIncorrect) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: errorTemplate.ErrEmailPasswordIncorrect.Error(),
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
		Message: "Change Password Success",
	})
}

// Get User Header Information
//
// @Summary			Get User Header Information
// @Tags			user
// @Produce			json
// @Router			/user/headerinfo	[get]
// @Security 		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrorResponse
// @Failure			500		{object}	dto.ErrorResponse
func (u *UserHandler) GetUserHeaderInformation(ctx *gin.Context) {
	id, exists := ctx.Get("idUser")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	res, err := u.us.GetUserHeaderInformation(ctx.Request.Context(), id.(int))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    res,
		Message: "Success Get Information Header User",
	})
}
