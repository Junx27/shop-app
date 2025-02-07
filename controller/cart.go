package controller

import (
	"net/http"

	"strconv"

	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/helper"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	repository entity.CartRepository
	service    entity.MenuService
}

func NewCartHandler(repository entity.CartRepository, service entity.MenuService) *CartHandler {
	return &CartHandler{repository: repository, service: service}
}

func (h *CartHandler) GetOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	cart, err := h.repository.GetOne(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	ctx.JSON(http.StatusOK, helper.SuccessResponse("Fetch data successfully", cart))
}

func (h *CartHandler) CreateOne(ctx *gin.Context) {
	cart := &entity.Cart{}
	if err := ctx.ShouldBindJSON(&cart); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Invalid request data"))
		return
	}
	subTotal, err := h.service.CalculateSubTotal(ctx, cart.MenuID, cart.Quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to calculate sub total"))
		return
	}
	cart.Subtotal = subTotal

	createCart, err := h.repository.CreateOne(ctx, cart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to create data"))
		return
	}

	ctx.JSON(http.StatusCreated, helper.SuccessResponse("Create data successfully", createCart))
}

func (h *CartHandler) DeleteOne(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Invalid ID"))
		return
	}
	err = h.repository.DeleteOne(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to delete data"))
		return
	}
	ctx.JSON(http.StatusOK, helper.SuccessResponse("Delete data successfully", nil))
}
