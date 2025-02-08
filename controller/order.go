package controller

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/helper"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	repository  entity.OrderRepository
	cartService entity.CartService
}

func NewOrderHandler(repository entity.OrderRepository, cartService entity.CartService) *OrderHandler {
	return &OrderHandler{repository: repository, cartService: cartService}
}

func (h *OrderHandler) GetMany(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	menus, totalItems, err := h.repository.GetMany(ctx, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limitInt)))

	if pageInt > totalPages {
		pageInt = totalPages
		menus, _, err = h.repository.GetMany(ctx, pageInt, limitInt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
			return
		}
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Page not found"))
		return
	}
	response := helper.PaginationResponse(menus, pageInt, limitInt, totalPages, totalItems)
	ctx.JSON(http.StatusOK, helper.SuccessResponse(("Fetch data successfully"), response))
}

func (h *OrderHandler) CreateOne(ctx *gin.Context) {
	order := &entity.Order{}
	userID, err := helper.GetUserIDFromCookie(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	calculation, err := h.cartService.CalculatePrice(ctx, uint(userID), "pending")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse(fmt.Sprintf("Failed to calculate total price: %v", err)))
		return
	}
	order.UserID = userID
	order.Total = int(calculation.TotalPrice)
	createOrder, err := h.repository.CreateOne(ctx, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to create data"))
		return
	}

	err = h.cartService.UpdateOrderIDInPendingCarts(ctx, uint(userID), createOrder.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse(fmt.Sprintf("Failed to update order_id on pending carts: %v", err)))
		return
	}
	ctx.JSON(http.StatusCreated, helper.SuccessResponse("Create data successfully", createOrder))
}
