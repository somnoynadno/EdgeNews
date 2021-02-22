package crawlers

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/db/dao"
	"EdgeNews/backend/models"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/utils"
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EchoMskCrawler struct {
	models.StaticWebCrawler
}

func (c EchoMskCrawler) RunForever() {
	log.Info("[ECHO MSK CRAWLER] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.EchoMskCrawler

	err := c.RecoverAfterRestart(4, c)
	if err != nil {
		log.Warn("[ECHO MSK CRAWLER] Recovering: " + err.Error())
	}

	for {
		log.Debug("[ECHO MSK CRAWLER] Awake")

		err := c.FindAvailableTextStreams()
		if err != nil {
			log.Error("[ECHO MSK CRAWLER] Crawling failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("echo msk").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("echo msk").Inc()
		}

		time.Sleep(sleepTime)
	}
}

func (c EchoMskCrawler) FindAvailableTextStreams() error {
	fetchURL := "https://echo.msk.ru/onlines/"

	resp, err := http.Get(fetchURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("status code: " + strconv.Itoa(resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".preview").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("a").Attr("href")
		if exists == true && strings.Contains(link, "/onlines/") {
			link = "https://echo.msk.ru" + link

			exists, err := dao.CheckTextStreamExistByURL(link)
			if err != nil {
				log.Error(err)
				return
			}

			if !exists {
				textStream := entities.TextStream{
					SourceID: 4,
					URL:      link,
					Name:     s.Find(".title").Text(),
					IsActive: true,
				}

				err := dao.AddTextStream(&textStream)
				if err != nil {
					log.Error(err)
				} else {
					go c.StartTextStream(textStream)
				}
			}
		}
	})

	return nil
}

func (c EchoMskCrawler) StartTextStream(textStream entities.TextStream) {
	log.Info("[ECHO MSK CRAWLER] Starting text stream on " + textStream.URL)
	sleepTime := config.GetConfig().TextStreamUpdateInterval
	maxEmptyFetches := config.GetConfig().TextStreamMaxEmptyFetches

	emptyFetchesCounter := 0
	for {
		log.Debug("[ECHO MSK CRAWLER] Awake")

		fetched, err := c.fetchMessages(textStream)
		if err != nil {
			log.Error("[ECHO MSK CRAWLER] Crawler failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("echo msk").Inc()
			emptyFetchesCounter++
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("echo msk").Inc()
			if fetched > 0 {
				log.Debug("[ECHO MSK CRAWLER] New messages: " + strconv.Itoa(fetched))
				emptyFetchesCounter = 0
				go dao.SetStreamUpdated(&textStream)
			} else {
				emptyFetchesCounter++
			}
		}

		if emptyFetchesCounter > maxEmptyFetches {
			break
		}

		time.Sleep(sleepTime)
	}

	finished := false
	for !finished {
		err := dao.FinishTextStream(&textStream)
		if err != nil {
			log.Error("[ECHO MSK CRAWLER] " + err.Error())
			time.Sleep(sleepTime)
		} else {
			finished = true
		}
	}

	log.Info("[ECHO MSK CRAWLER] Text stream finished: " + textStream.URL)
}

func (c EchoMskCrawler) fetchMessages(textStream entities.TextStream) (int, error) {
	log.Debug("[ECHO MSK CRAWLER] Fetching messages")
	newMessages := 0

	resp, err := http.Get(textStream.URL)
	if err != nil {
		return newMessages, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return newMessages, err
	}

	selection := doc.Find(".online__list-item")
	for i, j := 0, len(selection.Nodes)-1; i < j; i, j = i+1, j-1 {
		selection.Nodes[i], selection.Nodes[j] = selection.Nodes[j], selection.Nodes[i]
	}

	selection.Each(func(i int, s *goquery.Selection) {
		t := s.Find(".online__list-item-time").Text()
		content := s.Find(".online__list-item-content").Text()

		exists, err := dao.CheckMessageExist(content, textStream.ID)
		if err != nil {
			log.Error("[ECHO MSK CRAWLER] " + err.Error())
		} else {
			if !exists {
				newMessages++
				message := entities.Message{
					TextStreamID: textStream.ID,
					Body:         content,
					Time:         &t,
				}
				err := c.SaveMessage(message)
				if err != nil {
					log.Error("[ECHO MSK CRAWLER] " + err.Error())
				}
			}
		}
	})

	return newMessages, nil
}
