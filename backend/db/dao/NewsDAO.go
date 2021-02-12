package dao

import (
	"EdgeNews/backend/db"
	"EdgeNews/backend/models/entities"
	"github.com/jinzhu/gorm"
)

func CheckTitleExist(title string) (bool, error) {
	var news entities.News
	err := db.GetDB().Where("title = ?", title).First(&news).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else {
		return true, err
	}
}

func AddNews(news *entities.News) error {
	return db.GetDB().Create(news).Error
}