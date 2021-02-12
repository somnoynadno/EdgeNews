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

type newscatcherAPINews struct {
	Title         string
	Link          string
	PublishedDate string `json:"published_date"`
	Rights        string
	Topic         string
	Summary       string
	Author        *string
}

type newscatcherAPIResponse struct {
	Status   string
	Articles []newscatcherAPINews
}

type NewscatcherAPIScrapper struct {
	models.ScrapperAPI
}

func (s NewscatcherAPIScrapper) RunOnce() error {
	fetchURL := "https://newscatcher.p.rapidapi.com/v1/latest_headlines?lang=ru"

	req, _ := http.NewRequest("GET", fetchURL, nil)
	req.Header.Add("x-rapidapi-key", os.Getenv("newscatcher_api_token"))
	req.Header.Add("x-rapidapi-host", "newscatcher.p.rapidapi.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result newscatcherAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	for _, a := range result.Articles {
		news := entities.News{
			Title:    a.Title,
			Date:     &a.PublishedDate,
			URL:      &a.Link,
			Rights:   &a.Rights,
			Body:     &a.Summary,
			Tag:      &a.Topic,
			Author:   a.Author,
			SourceID: 3,
		}

		err := s.SaveNews(news)
		if err != nil {
			log.Error("[NEWSCATCHER API] " + err.Error())
		}
	}

	return nil
}

func (s NewscatcherAPIScrapper) RunForever() {
	log.Info("[NEWSCATCHER API] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.NewscatcherAPI

	for {
		log.Debug("[NEWSCATCHER API] Awake")

		err := s.RunOnce()
		if err != nil {
			log.Error("[NEWSCATCHER API] Scrapper failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("newscatcher api").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("newscatcher api").Inc()
		}

		time.Sleep(sleepTime * time.Second)
	}
}
