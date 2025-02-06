package controller

import (
	"net/http"

	"github.com/Junx27/shop-app/entity"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository entity.UserReopository
}

func NewUserHandler(repository entity.UserReopository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (h *UserHandler) GetMany(ctx *gin.Context) {
	users, err := h.repository.GetMany()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}
