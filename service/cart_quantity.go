package service

import (
	"context"

	"github.com/Junx27/shop-app/entity"
)

type QuantityService struct {
	repository entity.CartRepository
}

func NewQuantityService(repository entity.CartRepository) entity.CartService {
	return &QuantityService{repository: repository}
}

func (s *QuantityService) IncreaseCart(ctx context.Context, id uint, qty int) error {
	err := s.repository.UpdateQuantity(ctx, id, "increase", qty)
	if err != nil {
		return err
	}
	return nil
}

func (s *QuantityService) DecreaseCart(ctx context.Context, id uint, qty int) error {
	err := s.repository.UpdateQuantity(ctx, id, "decrease", qty)
	if err != nil {
		return err
	}
	return nil
}

func (s *QuantityService) CalculatePrice(ctx context.Context, userID uint, status string) (*entity.CalculateCart, error) {
	carts, err := s.repository.GetManyByUserAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	var totalPrice int
	var totalQuantity int
	for _, cart := range carts {
		totalPrice += cart.Subtotal
		totalQuantity += cart.Quantity
	}

	return &entity.CalculateCart{
		TotalItems:    len(carts),
		TotalQuantity: totalQuantity,
		TotalPrice:    float64(totalPrice),
	}, nil
}
