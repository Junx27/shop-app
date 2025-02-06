package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/middleware"
	"github.com/Junx27/shop-app/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupMenuRouter(r *gin.Engine, db *gorm.DB) {
	menuRepository := repository.NewMenuRepository(db)
	menuHandler := controller.NewEventHandler(menuRepository)

	menuGroup := r.Group("/menus")
	menuGroup.Use(middleware.AuthProtected(db))
	{
		menuGroup.GET("", menuHandler.GetMany)
		menuGroup.Static("/image", "./uploads")
		menuGroup.GET("/download/:id", menuHandler.DownloadImage)
		menuGroup.GET("/:id", menuHandler.GetOne)
		menuGroup.POST("", middleware.RoleRequired("admin"), menuHandler.CreateOne)
		menuGroup.PUT("/:id", middleware.RoleRequired("admin"), menuHandler.UpdateOne)
		menuGroup.DELETE("/:id", middleware.RoleRequired("admin"), menuHandler.DeleteOne)
	}
}
