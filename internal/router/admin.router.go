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

func AdminRouter(router *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	adminRoute := router.Group("/admin")

	ar := repo.NewAdminRepo(db)
	as := service.NewAdminService(ar)
	ah := handler.NewAdminHandler(as)

	adminRoute.GET("",middleware.CheckToken(rc), middleware.AdminAccess, ah.GetDataDashboard)
	adminRoute.GET("users", middleware.CheckToken(rc), middleware.AdminAccess, ah.GetUsers)
	adminRoute.GET("events", middleware.CheckToken(rc), middleware.AdminAccess, ah.GetEvents)
	adminRoute.GET("communities", middleware.CheckToken(rc), middleware.AdminAccess, ah.GetCommunities)
}