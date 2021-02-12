package dao

import (
	"EdgeNews/backend/db"
	"EdgeNews/backend/models/entities"
)

func GetAllSources() ([]entities.Source, error) {
	var sources []entities.Source
	err := db.GetDB().Preload("ScrapperType").Find(&sources).Error
	return sources, err
}
