package crawlers

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/db/dao"
	"EdgeNews/backend/models"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/utils"
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type KulichkiCrawler struct {
	models.StaticWebCrawler
}

func (c KulichkiCrawler) RunForever() {
	log.Info("[KULICHKI CRAWLER] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.KulichkiCrawler

	for {
		log.Debug("[KULICHKI CRAWLER] Awake")

		err := c.FindAvailableTextStreams()
		if err != nil {
			log.Error("[KULICHKI CRAWLER] Crawling failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("kulichki").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("kulichki").Inc()
		}

		time.Sleep(sleepTime * time.Second)
	}
}

func (c KulichkiCrawler) FindAvailableTextStreams() error {
	fetchURL := "https://football.kulichki.net/trans/"

	resp, err := http.Get(fetchURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("status code: " + strconv.Itoa(resp.StatusCode))
	}

	// fucking encoding
	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(resp.Body)

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	doc.Find("center table tr").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a")

		link, exist := a.Attr("href")
		if exist && strings.Contains(link, "trans") {
			link = "https://football.kulichki.net" + link

			exists, err := dao.CheckTextStreamExistByURL(link)
			if err != nil {
				log.Error(err)
				return
			}

			if !exists {
				textStream := entities.TextStream{
					SourceID: 5,
					URL:      link,
					Name:     strings.Split(a.Text(), " - ")[0],
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

func (c KulichkiCrawler) StartTextStream(textStream entities.TextStream) {
	log.Info("[KULICHKI CRAWLER] Starting text stream on " + textStream.URL)
	sleepTime := config.GetConfig().TextStreamUpdateInterval
	maxEmptyFetches := config.GetConfig().TextStreamMaxEmptyFetches

	emptyFetchesCounter := 0
	for {
		log.Debug("[KULICHKI CRAWLER] Awake")

		fetched, err := c.fetchMessages(textStream)
		if err != nil {
			log.Error("[KULICHKI CRAWLER] Crawler failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("kulichki").Inc()
			emptyFetchesCounter++
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("kulichki").Inc()
			if fetched > 0 {
				log.Debug("[KULICHKI CRAWLER] New messages: " + strconv.Itoa(fetched))
				emptyFetchesCounter = 0
				go dao.SetStreamUpdated(&textStream)
			} else {
				emptyFetchesCounter++
			}
		}

		if emptyFetchesCounter > maxEmptyFetches {
			break
		}

		time.Sleep(sleepTime * time.Second)
	}

	finished := false
	for !finished {
		err := dao.FinishTextStream(&textStream)
		if err != nil {
			log.Error("[KULICHKI CRAWLER] " + err.Error())
			time.Sleep(sleepTime * time.Second)
		} else {
			finished = true
		}
	}

	log.Info("[KULICHKI CRAWLER] Text stream finished: " + textStream.URL)
}

func (c KulichkiCrawler) fetchMessages(textStream entities.TextStream) (int, error) {
	log.Debug("[KULICHKI CRAWLER] Fetching messages")
	newMessages := 0

	resp, err := http.Get(textStream.URL)
	if err != nil {
		return newMessages, err
	}
	defer resp.Body.Close()

	// fucking encoding
	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(resp.Body)

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return newMessages, err
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return newMessages, err
	}

	selection := doc.Find("table tr")
	for i, j := 0, len(selection.Nodes)-1; i < j; i, j = i+1, j-1 {
		selection.Nodes[i], selection.Nodes[j] = selection.Nodes[j], selection.Nodes[i]
	}

	selection.Each(func(i int, s *goquery.Selection) {
		t := strings.TrimSpace(s.Find("td font").Text())
		t = strings.Split(t, "\n")[0]

		if t != "" && len(t) < 5 {
			content := strings.TrimSpace(s.Find("td span").Text())

			exists, err := dao.CheckMessageExistByBody(content)
			if err != nil {
				log.Error("[KULICHKI CRAWLER] " + err.Error())
			} else {
				if !exists && len(content) < 2400 {
					newMessages++
					message := entities.Message{
						TextStreamID: textStream.ID,
						Body:         content,
						Time:         &t,
					}
					err := c.SaveMessage(message)
					if err != nil {
						log.Error("[KULICHKI CRAWLER] " + err.Error())
					}
				}
			}
		}
	})

	return newMessages, nil
}
