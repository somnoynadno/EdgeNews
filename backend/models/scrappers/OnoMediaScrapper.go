package scrappers

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/models"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type onoMediaSingleArticle struct {
	Title      string
	Link       string
	SourceType string `json:"source_type"`
	CreatedAt  string
}

type onoMediaOutput struct {
	Articles [][]onoMediaSingleArticle
}

type onoMediaResponse struct {
	Outputs []onoMediaOutput
}

type OnoMediaScrapper struct {
	models.ScrapperAPI
}

func (s OnoMediaScrapper) RunOnce() error {
	fetchURL := "https://api.onomedia.today/api/period_outputs?skip=0&count=100&amount=3&hours=3&noise=false"

	resp, err := http.Get(fetchURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result onoMediaResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	for _, o := range result.Outputs {
		for _, row := range o.Articles {
			for _, sa := range row {
				u, err := url.Parse(sa.Link)
				if err != nil {
					log.Error("[ONO MEDIA API] " + err.Error())
					continue
				}

				news := entities.News{
					Title:    sa.Title,
					Date:     &sa.CreatedAt,
					URL:      &sa.Link,
					Rights:   &u.Host,
					Tag:      &sa.SourceType,
					SourceID: 6,
				}

				err = s.SaveNews(news)
				if err != nil {
					log.Error("[ONO MEDIA API] " + err.Error())
				}
			}
		}
	}

	return nil
}

func (s OnoMediaScrapper) RunForever() {
	log.Info("[ONO MEDIA API] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.OnoMediaAPI

	for {
		log.Debug("[ONO MEDIA API] Awake")

		err := s.RunOnce()
		if err != nil {
			log.Error("[ONO MEDIA API] Scrapper failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("ono media api").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("ono media api").Inc()
		}

		time.Sleep(sleepTime)
	}
}
