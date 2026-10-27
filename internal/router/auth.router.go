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

func authRouter(router *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	authRoute := router.Group("/auth")

	nr := repo.NewNotificationRepo(db)
	nh := service.NewHub()
	ns := service.NewNotificationService(nr, nh)

	ar := repo.NewAuthRepo(db)
	as := service.NewAuthService(ar, rc, ns)
	ah := handler.NewAuthHandler(as)

	authRoute.POST("register", ah.Register)
	authRoute.POST("login", ah.Login)
	authRoute.POST("logout", middleware.CheckToken(rc), ah.Logout)
	authRoute.POST("forgot-password", ah.ForgotPassword)
	authRoute.POST("reset-password", ah.ResetPassword)
}
