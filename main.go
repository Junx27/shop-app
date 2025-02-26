package main

import (
	"github.com/Junx27/shop-app/config"
	"github.com/Junx27/shop-app/database"
	"github.com/Junx27/shop-app/repository"
	"github.com/Junx27/shop-app/router"
	"github.com/Junx27/shop-app/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewEnvConfig()
	db := database.Init(cfg, database.DBMigrator)
	r := gin.Default()
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	menuRepository := repository.NewMenuRepository(db)
	menuService := service.NewCalculateService(menuRepository)
	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewQuantityService(cartRepository)
	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository)
	router.SetupAuthRouter(r, authService.(*service.AuthService))
	router.SetupUserRouter(r, db)
	router.SetupMenuRouter(r, db)
	router.SetupCartRouter(r, db, menuService.(*service.MenuService), cartService.(*service.CartService))
	router.SetupOrderRouter(r, db, cartService.(*service.CartService))
	router.SetupReportRouter(r, db, orderService.(*service.OrderService), orderRepository, cartRepository)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"App Name": "Shop App",
			"Author":   "Junx",
			"Version":  "1.0.0",
		})
	})
	r.Run(":8080")
}
