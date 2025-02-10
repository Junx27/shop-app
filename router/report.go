package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/middleware"
	"github.com/Junx27/shop-app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupReportRouter(r *gin.Engine, db *gorm.DB, orderService *service.OrderService, orderRepo entity.OrderRepository, cartRepo entity.CartRepository) {
	reportHandler := controller.NewReportHandler(orderRepo, cartRepo, orderService)
	reportGroup := r.Group("/reports")
	reportGroup.Use(middleware.AuthProtected(db), middleware.RoleRequired("admin"))
	{
		reportGroup.GET("/carts", reportHandler.GetManyCart)
		reportGroup.GET("/orders", reportHandler.GetManyOrder)
		reportGroup.GET("", reportHandler.SumaryReport)
	}
}
