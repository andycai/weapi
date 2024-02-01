package database

import (
	"fmt"

	"gorm.io/gorm"
)

func DecorateLike(tx *gorm.DB, field, value string) *gorm.DB {
	if value != "" {
		return tx.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%s%%", value))
	}

	return tx
}

func DecorateEqualInt(tx *gorm.DB, field string, value int) *gorm.DB {
	if value > 0 {
		return tx.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return tx
}
