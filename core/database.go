package core

import (
	"gorm.io/gorm"
)

func SetupDatabase(dbs []*gorm.DB) {
	for _, f := range dbMap {
		f(dbs)
	}
}
