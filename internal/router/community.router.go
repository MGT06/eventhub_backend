package router

import (
	"github.com/MGT06/EventHub_Backend.git/internal/handler"
	"github.com/MGT06/EventHub_Backend.git/internal/middleware"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
	"github.com/MGT06/EventHub_Backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func CommunityRouter(ctx *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	communityRoute := ctx.Group("/community")

	cr := repo.NewCommunityRepo(db)
	cs := service.NewCommunityService(cr)
	ch := handler.NewCommunityHandler(cs)

	communityRoute.GET("", ch.GetCommunityBySearchFilter)
	communityRoute.GET("detail/:id", ch.GetCommunity)
	communityRoute.POST("join/:id", middleware.CheckToken(rc), ch.ToggleJoinCommunity)
	communityRoute.GET("members/:id", ch.GetCommunityMembers)
	communityRoute.GET("popular", ch.GetPopularCommunity)
	communityRoute.GET("event/:id", ch.GetUpComingEventCommunity)
}
