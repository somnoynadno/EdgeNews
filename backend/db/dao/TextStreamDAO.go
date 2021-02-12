package dao

import (
	"EdgeNews/backend/db"
	"EdgeNews/backend/models/entities"
	"github.com/jinzhu/gorm"
)

func CheckTextStreamExistByURL(url string) (bool, error) {
	var textStream entities.TextStream
	err := db.GetDB().Where("url = ?", url).First(&textStream).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else {
		return true, err
	}
}

func AddTextStream(textStream *entities.TextStream) error {
	return db.GetDB().Create(textStream).Error
}

func FinishTextStream(textStream *entities.TextStream) error {
	return db.GetDB().Model(textStream).Update("is_active", false).Error
}