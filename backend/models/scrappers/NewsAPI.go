package scrappers

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/models"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type newsAPINewsSource struct {
	Name string
}

type newsAPINews struct {
	Title       string
	Description string
	URL         string
	PublishedAt string
	Author      *string
	Source      newsAPINewsSource
}

type newsAPIResponse struct {
	Status       string
	TotalResults int
	Articles     []newsAPINews
}

type NewsAPIScrapper struct {
	models.ScrapperAPI
}

func (s NewsAPIScrapper) RunOnce() error {
	fetchURL := "http://newsapi.org/v2/top-headlines?country=ru&apiKey=" + os.Getenv("news_api_token")

	resp, err := http.Get(fetchURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result newsAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	for _, a := range result.Articles {
		news := entities.News{
			Title:       a.Title,
			Date:        &a.PublishedAt,
			URL:         &a.URL,
			Rights:      &a.Source.Name,
			Description: &a.Description,
			Author:      a.Author,
			SourceID:    2,
		}

		err := s.SaveNews(news)
		if err != nil {
			log.Error("[NEWS API] " + err.Error())
		}
	}

	return nil
}

func (s NewsAPIScrapper) RunForever() {
	log.Info("[NEWS API] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.NewsAPI

	for {
		log.Debug("[NEWS API] Awake")

		err := s.RunOnce()
		if err != nil {
			log.Error("[NEWS API] Scrapper failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("news api").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("news api").Inc()
		}

		time.Sleep(sleepTime * time.Second)
	}
}
