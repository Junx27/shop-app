package controller

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/helper"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	repositoryOrder entity.OrderRepository
	repositoryCart  entity.CartRepository
	orderService    entity.OrderService
}

func NewReportHandler(repositoryOrder entity.OrderRepository, repositoryCart entity.CartRepository, orderService entity.OrderService) *ReportHandler {
	return &ReportHandler{repositoryOrder: repositoryOrder, repositoryCart: repositoryCart, orderService: orderService}
}

func (h *ReportHandler) GetManyCart(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	carts, totalItems, err := h.repositoryCart.GetManyAdmin(ctx, pageInt, limitInt)
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
func (h *ReportHandler) GetManyOrder(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	orders, totalItems, err := h.repositoryOrder.GetManyAdmin(ctx, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limitInt)))

	if pageInt > totalPages {
		ctx.JSON(http.StatusNotFound, helper.FailedResponse("Data not found"))
		return
	}
	ctx.JSON(http.StatusOK, helper.PaginationResponse("Fetch data successfully", pageInt, limitInt, totalPages, totalItems, orders))
}

func (h *ReportHandler) SumaryReport(ctx *gin.Context) {
	orderReport, err := h.orderService.CalculateOrder(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	ctx.JSON(http.StatusOK, helper.SuccessResponse(("Fetch data successfully"), orderReport))
}
