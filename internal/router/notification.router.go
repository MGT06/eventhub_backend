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

func notificationRouter(router *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	notificationRoute := router.Group("/notification")

	hs := service.NewHub()
	
	nr := repo.NewNotificationRepo(db)
	ns := service.NewNotificationService(nr, hs)
	nh := handler.NewNotificationHandler(ns)
	
	sh := handler.NewStreamHandler(hs, ns)

	notificationRoute.GET("", middleware.CheckToken(rc), nh.GetMyNotification)
	notificationRoute.GET("stream", middleware.CheckToken(rc), sh.Stream)
}
