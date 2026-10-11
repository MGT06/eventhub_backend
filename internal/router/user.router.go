package router

import (
	"path"

	"github.com/MGT06/EventHub_Backend.git/internal/handler"
	"github.com/MGT06/EventHub_Backend.git/internal/middleware"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
	"github.com/MGT06/EventHub_Backend.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func UserRouter(router *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	userRoute := router.Group("/user")

	router.Static("img", path.Join("public", "img"))

	ur := repo.NewUserRepo(db)
	us := service.NewUserService(ur, rc)
	uh := handler.NewUserHandler(us)

	userRoute.GET("", middleware.CheckToken(rc), uh.GetProfileUser)
	userRoute.PATCH("edit", middleware.CheckToken(rc), uh.EditProfileUser)
	userRoute.POST("change-password", middleware.CheckToken(rc), uh.ChangePassword)
	userRoute.GET("headerinfo", middleware.CheckToken(rc), uh.GetUserHeaderInformation)
}
