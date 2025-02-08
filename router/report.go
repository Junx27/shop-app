package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/middleware"
	"github.com/Junx27/shop-app/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupReportRouter(r *gin.Engine, db *gorm.DB, orderService *service.OrderService, orderRepo entity.OrderRepository) {
	reportHandler := controller.NewReportHandler(orderRepo, orderService)
	reportGroup := r.Group("/reports")
	reportGroup.Use(middleware.AuthProtected(db))
	{
		reportGroup.GET("", middleware.RoleRequired("admin"), reportHandler.SumaryReport)
	}
}
