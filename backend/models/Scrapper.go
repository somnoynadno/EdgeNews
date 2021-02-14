package models

import (
	"EdgeNews/backend/db/dao"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/server"
	"EdgeNews/backend/utils"
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
	return errors.New("should implement")
}

func (s ScrapperAPI) RunForever() {
	log.Panic("should implement")
}

func (s ScrapperAPI) SaveNews(news entities.News) error {
	exist, err := dao.CheckNewsExistByTitle(news.Title)
	if err != nil {
		return err
	}

	if !exist {
		err = dao.AddNews(&news)
		if err != nil {
			return err
		}

		b, _ := json.Marshal(news)
		utils.GetMetrics().WS.BroadcastMessages.WithLabelValues("api scrapper").Inc()
		go server.GetNewsHub().SendMessage(b)
	}

	return nil
}
