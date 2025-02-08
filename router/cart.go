package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/middleware"
	"github.com/Junx27/shop-app/repository"
	"github.com/Junx27/shop-app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCartRouter(r *gin.Engine, db *gorm.DB, menuService *service.MenuService, cartService *service.CartService) {
	cartRepository := repository.NewCartRepository(db)
	cartHandler := controller.NewCartHandler(cartRepository, menuService, cartService)

	cartGroup := r.Group("/cart")
	cartGroup.Use(middleware.AuthProtected(db))
	{
		cartGroup.GET("", cartHandler.GetMany)
		cartGroup.GET("/:id", cartHandler.GetOne)
		cartGroup.POST("", cartHandler.CreateOne)
		cartGroup.PATCH("/increase/:id", cartHandler.Increase)
		cartGroup.PATCH("/decrease/:id", cartHandler.Decrease)
		cartGroup.GET("/total", cartHandler.CalculateTotalPrice)
		cartGroup.DELETE("/:id", cartHandler.DeleteOne)
	}
}
