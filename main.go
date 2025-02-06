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
	router.SetupAuthRouter(r, authService.(*service.AuthService))
	router.SetupUserRouter(r, db)
	router.SetupMenuRouter(r, db)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"App Name": "Shop App",
			"Author":   "Junx",
			"Version":  "1.0.0",
		})
	})
	r.Run(":8080")
}
