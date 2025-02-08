package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/middleware"
	"github.com/Junx27/shop-app/repository"
	"github.com/Junx27/shop-app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRouter(r *gin.Engine, db *gorm.DB, cartService *service.CartService) {
	OrderRepository := repository.NewOrderRepository(db)
	orderHandler := controller.NewOrderHandler(OrderRepository, cartService)

	orderGroup := r.Group("/orders")
	orderGroup.Use(middleware.AuthProtected(db))
	{
		orderGroup.GET("", orderHandler.GetMany)
		orderGroup.POST("", orderHandler.CreateOne)
		orderGroup.PATCH("/payment/:id", orderHandler.UpdatePayment)
	}
}
