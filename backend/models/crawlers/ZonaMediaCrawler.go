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

type ZonaMediaCrawler struct {
	models.StaticWebCrawler
}

func (c ZonaMediaCrawler) RunForever() {
	log.Info("[ZONA MEDIA CRAWLER] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.ZonaMediaCrawler

	for {
		log.Debug("[ZONA MEDIA CRAWLER] Awake")

		err := c.FindAvailableTextStreams()
		if err != nil {
			log.Error("[ZONA MEDIA CRAWLER] Crawling failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("zona media").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("zona media").Inc()
		}

		time.Sleep(sleepTime)
	}
}

func (c ZonaMediaCrawler) FindAvailableTextStreams() error {
	fetchURL := "https://zona.media/onlines"

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

	doc.Find(".mz-materials li").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("a").Attr("href")
		if exists && strings.Contains(link, "/online/") {
			link = "https://zona.media" + link

			exists, err := dao.CheckTextStreamExistByURL(link)
			if err != nil {
				log.Error(err)
				return
			}

			if !exists {
				textStream := entities.TextStream{
					SourceID: 6,
					URL:      link,
					Name:     s.Find("header").Text(),
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

func (c ZonaMediaCrawler) StartTextStream(textStream entities.TextStream) {
	log.Info("[ZONA MEDIA CRAWLER] Starting text stream on " + textStream.URL)
	sleepTime := config.GetConfig().TextStreamUpdateInterval
	maxEmptyFetches := config.GetConfig().TextStreamMaxEmptyFetches

	emptyFetchesCounter := 0
	for {
		log.Debug("[ZONA MEDIA CRAWLER] Awake")

		fetched, err := c.fetchMessages(textStream)
		if err != nil {
			log.Error("[ZONA MEDIA CRAWLER] Crawler failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("zona media").Inc()
			emptyFetchesCounter++
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("zona media").Inc()
			if fetched > 0 {
				log.Debug("[ZONA MEDIA CRAWLER] New messages: " + strconv.Itoa(fetched))
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
			log.Error("[ZONA MEDIA CRAWLER] " + err.Error())
			time.Sleep(sleepTime)
		} else {
			finished = true
		}
	}

	log.Info("[ZONA MEDIA CRAWLER] Text stream finished: " + textStream.URL)
}

func (c ZonaMediaCrawler) fetchMessages(textStream entities.TextStream) (int, error) {
	log.Debug("[ZONA MEDIA CRAWLER] Fetching messages")
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

	selection := doc.Find(".mz-publish__text__item")
	for i, j := 0, len(selection.Nodes)-1; i < j; i, j = i+1, j-1 {
		selection.Nodes[i], selection.Nodes[j] = selection.Nodes[j], selection.Nodes[i]
	}

	selection.Each(func(i int, s *goquery.Selection) {
		t := s.Find(".mz-publish__time-header").Text()

		text := ""
		s.Find("body p").Each(func(i int, c *goquery.Selection) {
			text += c.Text() + "\n"
		})

		exists, err := dao.CheckMessageExist(text, textStream.ID)
		if err != nil {
			log.Error("[ZONA MEDIA CRAWLER] " + err.Error())
		} else {
			if !exists {
				newMessages++
				message := entities.Message{
					TextStreamID: textStream.ID,
					Body:         text,
					Time:         &t,
				}
				err := c.SaveMessage(message)
				if err != nil {
					log.Error("[ZONA MEDIA CRAWLER] " + err.Error())
				}
			}
		}
	})

	return newMessages, nil
}
