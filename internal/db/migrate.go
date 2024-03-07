package db

import (
	"github.com/SicParv1sMagna/AtomHackMarsService/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Document{}, &model.File{})
	if err != nil {
		panic("cant migrate db")
	}

	return nil
}
