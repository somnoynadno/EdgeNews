package dao

import (
	"EdgeNews/backend/db"
	"EdgeNews/backend/models/entities"
	"github.com/jinzhu/gorm"
)

func CheckMessageExistByBody(body string) (bool, error) {
	var message entities.Message
	err := db.GetDB().Where("body = ?", body).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else {
		return true, err
	}
}

func AddMessage(message *entities.Message) error {
	return db.GetDB().Create(message).Error
}
