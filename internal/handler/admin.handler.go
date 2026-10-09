package handler

import (
	"log"
	"net/http"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	"github.com/MGT06/EventHub_Backend.git/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	as *service.AdminService
}

func NewAdminHandler(as *service.AdminService) *AdminHandler {
	return &AdminHandler{
		as: as,
	}
}

// Get Data Dashboard Admin
//
// @Summary			Get Data Dashboard Admin
// @Tags			admin
// @Produce			json
// @Router			/admin	[get]
// @Security 		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (a *AdminHandler) GetDataDashboard(ctx *gin.Context) {
	res, err := a.as.GetDataDashboard(ctx.Request.Context())
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
		Message: "Success Get Data Dashboard",
	})
}

func (a *AdminHandler) GetUsers(ctx *gin.Context) {
	search := ctx.Query("search")

	res, err := a.as.GetUsers(ctx, search)
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
		Message: "Success Get Data Users",
	})
}

func (a *AdminHandler) GetEvents(ctx *gin.Context) {
	res, err := a.as.GetEvents(ctx)
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
		Message: "Success Get Data All Event",
	})
}