package controller

import (
	"fmt"
	"math"
	"net/http"

	"strconv"

	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	repository  entity.CartRepository
	menuService entity.MenuService
	cartService entity.CartService
}

func NewCartHandler(repository entity.CartRepository, menuService entity.MenuService, cartService entity.CartService) *CartHandler {
	return &CartHandler{repository: repository, menuService: menuService, cartService: cartService}
}

func (h *CartHandler) GetMany(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromCookie(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	carts, totalItems, err := h.repository.GetMany(ctx, userID, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limitInt)))

	if pageInt > totalPages {
		ctx.JSON(http.StatusNotFound, helper.FailedResponse("Data not found"))
		return
	}

	ctx.JSON(http.StatusOK, helper.PaginationResponse("Fetch data successfully", pageInt, limitInt, totalPages, totalItems, carts))
}

func (h *CartHandler) CalculateTotalPrice(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, helper.SuccessResponse("Total price calculated successfully", calculation))
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
	userID, err := helper.GetUserIDFromCookie(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	existingCart, err := h.repository.FindByUserAndMenuAndStatus(ctx, userID, cart.MenuID, "pending")
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to check existing cart"))
		return
	}
	if existingCart != nil {
		existingCart.Quantity += cart.Quantity
		subTotal, err := h.menuService.CalculateSubTotal(ctx, cart.MenuID, existingCart.Quantity)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Quantity is unavailable"))
			return
		}
		existingCart.Subtotal = subTotal
		err = h.menuService.DecreaseMenu(ctx, cart.MenuID, cart.Quantity)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Please check quantity"))
			return
		}
		updatedCart, err := h.repository.UpdateOne(ctx, existingCart)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to update cart"))
			return
		}

		ctx.JSON(http.StatusOK, helper.SuccessResponse("Cart updated successfully", updatedCart))
		return
	}
	subTotal, err := h.menuService.CalculateSubTotal(ctx, cart.MenuID, cart.Quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to calculate sub total"))
		return
	}
	err = h.menuService.DecreaseMenu(ctx, cart.MenuID, cart.Quantity)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Please check quantity"))
		return
	}
	cart.UserID = userID
	cart.Subtotal = subTotal
	createCart, err := h.repository.CreateOne(ctx, cart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to create data"))
		return
	}
	ctx.JSON(http.StatusCreated, helper.SuccessResponse("Create data successfully", createCart))
}

func (h *CartHandler) Increase(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Invalid ID"))
		return
	}
	qty := 1
	err = h.cartService.IncreaseCart(ctx, uint(id), qty)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to increase quantity"))
		return
	}
	ctx.JSON(http.StatusOK, helper.SuccessResponse("Increase quantity successfully", nil))
}

func (h *CartHandler) Decrease(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.FailedResponse("Invalid ID"))
		return
	}
	qty := 1
	err = h.cartService.DecreaseCart(ctx, uint(id), qty)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to decrease quantity"))
		return
	}
	ctx.JSON(http.StatusOK, helper.SuccessResponse("Decrease quantity successfully", nil))
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
