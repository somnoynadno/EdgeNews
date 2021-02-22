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
	RecoverAfterRestart(sourceID uint, crawler AbstractCrawler) error
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
	exist, err := dao.CheckMessageExist(message.Body, message.TextStreamID)
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

func (c StaticWebCrawler) RecoverAfterRestart(sourceID uint, crawler AbstractCrawler) error {
	ts, err := dao.GetActiveTextStreamsBySourceID(sourceID)
	if err != nil {
		return err
	}
	for _, t := range ts {
		go crawler.StartTextStream(t)
	}

	return nil
}

type DynamicWebCrawler struct {
	AbstractCrawler
}

func (c DynamicWebCrawler) RunForever() {
	log.Panic("should implement")
}

func (c DynamicWebCrawler) FindAvailableTextStreams() error {
	return errors.New("should implement")
}

func (c DynamicWebCrawler) StartTextStream(textStream entities.TextStream) {
	log.Panic("should implement")
}

func (c DynamicWebCrawler) SaveMessage(message entities.Message) error {
	exist, err := dao.CheckMessageExist(message.Body, message.TextStreamID)
	if err != nil {
		return err
	}

	if !exist {
		err = dao.AddMessage(&message)
		if err != nil {
			return err
		}

		b, _ := json.Marshal(message)
		utils.GetMetrics().WS.BroadcastMessages.WithLabelValues("dynamic crawler").Inc()
		go server.GetTextStreamHub().SendMessage(b)
	}

	return nil
}

func (c DynamicWebCrawler) RecoverAfterRestart(sourceID uint, crawler AbstractCrawler) error {
	ts, err := dao.GetActiveTextStreamsBySourceID(sourceID)
	if err != nil {
		return err
	}
	for _, t := range ts {
		go crawler.StartTextStream(t)
	}

	return nil
}
