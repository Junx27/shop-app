package service

import (
	"context"

	"github.com/Junx27/shop-app/entity"
)

type OrderService struct {
	repository entity.OrderRepository
}

func NewOrderService(repository entity.OrderRepository) entity.OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) CalculateOrder(ctx context.Context) (*entity.OrderReport, error) {
	orders, err := s.repository.GetManyByStatus(ctx, "pending")
	if err != nil {
		return nil, err
	}
	ordersPaid, err := s.repository.GetManyByStatus(ctx, "success")
	if err != nil {
		return nil, err
	}

	var amount int
	for _, order := range orders {
		amount += order.Total
	}

	var paidAmount int
	for _, order := range ordersPaid {
		paidAmount += order.Total
	}
	totalSales := amount + paidAmount
	return &entity.OrderReport{
		Amount:     float64(amount),
		PaidAmount: float64(paidAmount),
		TotalSales: float64(totalSales),
	}, nil
}
