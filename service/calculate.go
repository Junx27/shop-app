package service

import (
	"context"

	"github.com/Junx27/shop-app/entity"
)

type CalculateService struct {
	repository entity.MenuRepository
}

func NewCalculateService(repository entity.MenuRepository) entity.MenuService {
	return &CalculateService{repository: repository}
}

func (s *CalculateService) CalculateSubTotal(ctx context.Context, id uint, qty int) (int, error) {
	subTotal := 0
	menu, err := s.repository.GetOne(ctx, id)
	if err != nil {
		return subTotal, err
	}
	return menu.Price * qty, nil
}
