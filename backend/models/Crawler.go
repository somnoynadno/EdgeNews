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

type AbstractCrawler interface {
	RunForever()
	FindAvailableTextStreams() error
	StartTextStream(textStream entities.TextStream)
	SaveMessage(message entities.Message) error
}

type StaticWebCrawler struct {
	AbstractCrawler
}

func (c StaticWebCrawler) RunForever() {
	log.Panic("should implement")
}

func (c StaticWebCrawler) FindAvailableTextStreams() error {
	return errors.New("should implement")
}

func (c StaticWebCrawler) StartTextStream(textStream entities.TextStream) {
	log.Panic("should implement")
}

func (c StaticWebCrawler) SaveMessage(message entities.Message) error {
	exist, err := dao.CheckMessageExistByBody(message.Body)
	if err != nil {
		return err
	}

	if !exist {
		err = dao.AddMessage(&message)
		if err != nil {
			return err
		}

		b, _ := json.Marshal(message)
		utils.GetMetrics().WS.BroadcastMessages.WithLabelValues("static crawler").Inc()
		go server.GetTextStreamHub().SendMessage(b)
	}

	return nil
}

