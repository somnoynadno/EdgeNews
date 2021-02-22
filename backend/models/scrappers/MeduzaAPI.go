package scrappers

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/models"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type meduzaNewsSource struct {
	Name string
}

type meduzaNewsTag struct {
	Name string
}

type meduzaNews struct {
	URL         string
	Title       string
	SecondTitle string
	PubDate     string `json:"pub_date"`
	Source      meduzaNewsSource
	Tag         meduzaNewsTag
}

type meduzaResponse struct {
	Documents  map[string]meduzaNews
	Collection []string
}

type MeduzaScrapper struct {
	models.ScrapperAPI
}

func (s MeduzaScrapper) RunOnce() error {
	fetchURL := "https://meduza.io/api/v3/search?chrono=news&locale=ru&page=0&per_page=24"

	resp, err := http.Get(fetchURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result meduzaResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	for _, c := range result.Collection {
		mn := result.Documents[c]
		mn.URL = "https://meduza.io/" + mn.URL

		news := entities.News{
			Title:       mn.Title,
			Date:        &mn.PubDate,
			URL:         &mn.URL,
			Rights:      &mn.Source.Name,
			Tag:         &mn.Tag.Name,
			Description: &mn.SecondTitle,
			SourceID:    1,
		}

		err := s.SaveNews(news)
		if err != nil {
			log.Error("[MEDUZA API] " + err.Error())
		}
	}

	return nil
}

func (s MeduzaScrapper) RunForever() {
	log.Info("[MEDUZA API] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.MeduzaAPI

	for {
		log.Debug("[MEDUZA API] Awake")

		err := s.RunOnce()
		if err != nil {
			log.Error("[MEDUZA API] Scrapper failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("meduza api").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("meduza api").Inc()
		}

		time.Sleep(sleepTime)
	}
}
