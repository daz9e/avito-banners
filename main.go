package main

import (
	"avito-banners/config"
	"avito-banners/db"
	"avito-banners/handlers"
	"avito-banners/middleware"
	"avito-banners/tools"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func startRouter() {
	r := gin.Default()

	r.GET("/user_banner", middleware.AuthMiddleware(), handlers.GetUserBanner)
	r.GET("/banner", middleware.AdminAuthMiddleware(), handlers.GetBanners)
	r.POST("/banner", middleware.AdminAuthMiddleware(), handlers.CreateBanner)
	r.PATCH("/banner/:id", middleware.AdminAuthMiddleware(), handlers.UpdateBanner)
	r.DELETE("/banner/:id", middleware.AdminAuthMiddleware(), handlers.DeleteBanner)

	r.GET("/banners/:banner_id/versions", middleware.AdminAuthMiddleware(), handlers.GetBannerVersions)
	r.POST("/banners/:version_id/restore", middleware.AdminAuthMiddleware(), handlers.RestoreBannerVersion)

	log.Fatal(r.Run(":8080"))
}

func main() {
	config.Database = db.SetupDatabase()
	tools.SetupRedis()
	startRouter()
}
