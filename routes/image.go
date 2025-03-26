package routes

import (
	"r2-gallery/controllers"
	"r2-gallery/middleware"

	"github.com/gin-gonic/gin"
)

func SetupImageRoutes(router *gin.RouterGroup) {
	imageGroup := router.Group("/images")
	imageGroup.Use(middleware.AuthMiddleware())
	{
		imageGroup.POST("/upload", controllers.UploadImage)
		imageGroup.GET("/", controllers.ListImages)
		imageGroup.DELETE("/:id", controllers.DeleteImage)
	}
}
