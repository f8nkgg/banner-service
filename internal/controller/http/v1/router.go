package v1

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, bannerController *BannerController) {
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	authenticated := server.Group("/")
	authenticated.Use(authenticate)
	authenticated.POST("/banner", bannerController.createBanner)
	authenticated.GET("/user_banner", bannerController.getBanner)
	authenticated.GET("/banner", bannerController.getBanners)
	authenticated.DELETE("/banner/:id", bannerController.deleteBanner)
	authenticated.PATCH("/banner/:id", bannerController.updateBanner)
	authenticated.GET("/banner/history/:id", bannerController.getBannersHistoryByID)
}
