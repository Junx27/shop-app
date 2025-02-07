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
