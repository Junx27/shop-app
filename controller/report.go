package controller

import (
	"net/http"

	"github.com/Junx27/shop-app/entity"
	"github.com/Junx27/shop-app/helper"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	repository   entity.OrderRepository
	orderService entity.OrderService
}

func NewReportHandler(repository entity.OrderRepository, orderService entity.OrderService) *ReportHandler {
	return &ReportHandler{repository: repository, orderService: orderService}
}

func (h *ReportHandler) SumaryReport(ctx *gin.Context) {
	orderReport, err := h.orderService.CalculateOrder(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.FailedResponse("Failed to fetch data"))
		return
	}
	ctx.JSON(http.StatusOK, helper.SuccessResponse(("Fetch data successfully"), orderReport))
}
