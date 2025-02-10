package router

import (
	"github.com/Junx27/shop-app/controller"
	"github.com/Junx27/shop-app/middleware"
	"github.com/Junx27/shop-app/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRouter(r *gin.Engine, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userHandler := controller.NewUserHandler(userRepository)

	userGroup := r.Group("/users")
	userGroup.Use(middleware.AuthProtected(db))
	{
		userGroup.GET("/", middleware.RoleRequired("admin"), userHandler.GetMany)
	}
}
