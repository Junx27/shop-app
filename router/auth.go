package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/service"
	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(r *gin.Engine, authService *service.AuthService) {

	authHandler := controller.NewAuthHandler(authService)

	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)
	r.POST("/logout", authHandler.Logout)
}
