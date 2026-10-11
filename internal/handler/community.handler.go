package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	"github.com/MGT06/EventHub_Backend.git/internal/service"
	"github.com/gin-gonic/gin"
)

type CommunityHandler struct {
	cs *service.CommunityService
}

func NewCommunityHandler(cs *service.CommunityService) *CommunityHandler {
	return &CommunityHandler{
		cs: cs,
	}
}

// Get Community
//
// @Summary			Get Detail Community
// @Description		Get Detail Community by Id Community
// @Tags			community
// @Produce			json
// @Router			/community/detail/{id}	[get]
// @Param			id			path	int		true	"community id"
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityHandler) GetCommunity(ctx *gin.Context) {
	communityId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	res, err := c.cs.GetCommunity(ctx.Request.Context(), communityId)
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
		Message: "Success Get Community",
	})

}

func (c *CommunityHandler) GetUpComingEventCommunity(ctx *gin.Context) {
	communityId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}
	res, err := c.cs.GetUpComingEventCommunity(ctx.Request.Context(), communityId)
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
		Message: "Success Get UpComing Event Communtiy",
	})
}

// Get Community by Search & Filter
//
// @Summary			Get Community by Search & Filter
// @Tags			community
// @Produce			json
// @Router			/community	[get]
// @Param			search		query		string	false	"key Search in query param"
// @Param			category	query		string	false	"key category in query param"
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityHandler) GetCommunityBySearchFilter(ctx *gin.Context) {
	search := ctx.Query("search")
	filter := ctx.Query("category")

	res, err := c.cs.GetCommunityBySearchFilter(ctx.Request.Context(), search, filter)
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
		Message: "Success Get Community",
	})

}

// Toggle Join Commmunity
//
// @Summary			Toggle Join Commmunity
// @Tags			community
// @Produce			json
// @Router			/community/join/{id}		[post]
// @Param			id			path	int		true	"community id"
// @Security 		BearerToken
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (e *CommunityHandler) ToggleJoinCommunity(ctx *gin.Context) {
	idCommunity, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	idUser, exist := ctx.Get("idUser")
	if !exist {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	isJoin, er := e.cs.ToggleJoinCommunity(ctx.Request.Context(), idUser.(int), idCommunity)
	if er != nil {
		log.Println(er)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	if isJoin {
		ctx.JSON(http.StatusOK, dto.Response{
			Success: true,
			Message: "Success Leave Event",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Success Join Event",
	})
}

// Get Community Members
//
// @Summary			Get Community Members
// @Tags			community
// @Produce			json
// @Router			/community/members/{id}		[get]
// @Param			id			path	int		true	"community id"
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityHandler) GetCommunityMembers(ctx *gin.Context) {
	communityId := ctx.Param("id")
	id, err := strconv.Atoi(communityId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	res, err := c.cs.GetCommunityMembers(ctx.Request.Context(), id)
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
		Message: "Success Get Community",
	})

}

// Get Popular Community
//
// @Summary			Get Popular Community
// @Tags			community
// @Produce			json
// @Router			/community/popular	[get]
// @Success			200		{object}	dto.Response
// @Failure			500		{object}	dto.ErrorResponse
func (c *CommunityHandler) GetPopularCommunity(ctx *gin.Context) {
	res, err := c.cs.GetPopularCommunity(ctx.Request.Context())
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
		Message: "Success Get Popular Community",
	})

}

func (c *CommunityHandler) GetCommunityDiscussions(ctx *gin.Context) {
	communityId, e := strconv.Atoi(ctx.Param("id"))
	if e != nil {
		log.Println(e)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "A system error has occurred",
		})
		return
	}

	res, err := c.cs.GetCommunityDiscussions(ctx.Request.Context(), communityId)
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
		Message: "Success Get Discussions Community",
	})
}
