package models

import (
	"EdgeNews/backend/db/dao"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/server"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
)

type AbstractScrapper interface {
	RunForever()
	RunOnce() error
	SaveNews(news entities.News) error
}

type ScrapperAPI struct {
	AbstractScrapper
}

func (s ScrapperAPI) RunOnce() error {
	return errors.New("not implemented")
}

func (s ScrapperAPI) RunForever() {
	log.Error("not implemented")
}

func (s ScrapperAPI) SaveNews(news entities.News) error {
	exist, err := dao.CheckTitleExist(news.Title)
	if err != nil {
		return err
	}

	if !exist {
		err = dao.AddNews(&news)
		if err != nil {
			return err
		}

		b, _ := json.Marshal(news)
		go server.GetNewsHub().SendMessage(b)
	}

	return nil
}
