package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/middleware"
	"github.com/Junx27/shop-app/repository"
	"github.com/Junx27/shop-app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCartRouter(r *gin.Engine, db *gorm.DB, menuService *service.CalculateService) {
	cartRepository := repository.NewCartRepository(db)
	cartHandler := controller.NewCartHandler(cartRepository, menuService)

	cartGroup := r.Group("/carts")
	cartGroup.Use(middleware.AuthProtected(db))
	{
		cartGroup.GET("/:id", cartHandler.GetOne)
		cartGroup.POST("", cartHandler.CreateOne)
		cartGroup.DELETE("/:id", cartHandler.DeleteOne)
	}
}
