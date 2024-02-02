package weapi

import (
	"gorm.io/gorm"
)

func AutoMultiMigrate(dbs []*gorm.DB) error {
	for _, db := range dbs {
		if err := MakeMigrates(db, models); err != nil {
			return err
		}
	}
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	return MakeMigrates(db, models)
}

func MakeMigrates(db *gorm.DB, insts []any) error {
	for _, v := range insts {
		if err := db.AutoMigrate(v); err != nil {
			return err
		}
	}
	return nil
}
