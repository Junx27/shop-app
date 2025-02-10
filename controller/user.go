package controller

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/helper"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository entity.UserReopository
}

func NewUserHandler(repository entity.UserReopository) *UserHandler {
	return &UserHandler{repository: repository}
}

func (h *UserHandler) GetMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	users, totalItems, err := h.repository.GetMany(ctx, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limitInt)))

	if pageInt > totalPages {
		ctx.JSON(http.StatusNotFound, helper.FailedResponse("Data not found"))
		return
	}
	ctx.JSON(http.StatusOK, helper.PaginationResponse("Fetch data successfully", pageInt, limitInt, totalPages, totalItems, users))
}
