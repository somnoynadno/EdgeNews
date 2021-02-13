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

func GetMessagesByTextStreamID(textStreamID int) ([]entities.Message, error) {
	var messages []entities.Message
	err := db.GetDB().Where("text_stream_id = ?", textStreamID).Find(&messages).Error
	return messages, err
}
