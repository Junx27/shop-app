package database

import (
	"fmt"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	entities := []interface{}{
		&entity.User{},
	}
	for _, entity := range entities {
		if err := db.Migrator().DropTable(entity); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}

	return nil
}
