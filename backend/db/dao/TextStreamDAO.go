package dao

import (
	"EdgeNews/backend/db"
	"EdgeNews/backend/models/entities"
	"github.com/jinzhu/gorm"
	"time"
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

func SetStreamUpdated(textStream *entities.TextStream) error {
	return db.GetDB().Model(textStream).Update("last_stream_update", time.Now()).Error
}

func FinishTextStream(textStream *entities.TextStream) error {
	return db.GetDB().Model(textStream).Update("is_active", false).Error
}

func GetActiveTextStreams() ([]entities.TextStream, error) {
	var textStreams []entities.TextStream
	err := db.GetDB().Preload("Source").Preload("Source.ScrapperType").
		Where("is_active = true").Find(&textStreams).Error
	return textStreams, err
}

func GetActiveTextStreamsBySourceID(sourceID uint) ([]entities.TextStream, error) {
	var textStreams []entities.TextStream
	err := db.GetDB().Preload("Source").Preload("Source.ScrapperType").
		Where("source_id = ?", sourceID).Where("is_active = true").Find(&textStreams).Error
	return textStreams, err
}