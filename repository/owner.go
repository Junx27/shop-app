package repository

import "context"

type HasUserID interface {
	GetUserID(id uint) (uint, error)
	GetManyByUser(ctx context.Context, userID uint, page, limit int) ([]interface{}, error)
}
