package crawlers

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/db/dao"
	"EdgeNews/backend/models"
	"EdgeNews/backend/models/entities"
	"EdgeNews/backend/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var novayaGazetaWD selenium.WebDriver

type NovayaGazetaCrawler struct {
	models.DynamicWebCrawler
}

func (c NovayaGazetaCrawler) initWebDriver() (selenium.WebDriver, *selenium.Service, error) {
	seleniumPath := os.Getenv("selenium_path")
	geckoDriverPath := os.Getenv("gecko_driver_path")
	port := rand.Intn(10000) + 10000

	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(nil),                  // Output debug information to nowhere.
	}

	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		return nil, nil, err
	}

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return wd, service, err
	}

	log.Info("[DAEMON] Init standalone selenium on port: ", port)

	return wd, service, nil
}

func (c NovayaGazetaCrawler) RunForever() {
	log.Info("[NOVAYA GAZETA CRAWLER] Starting...")
	sleepTime := config.GetConfig().ScrappingIntervals.NovayaGazetaCrawler

	wd, s, err := c.initWebDriver()
	if err != nil {
		log.Fatal(err)
	}
	novayaGazetaWD = wd

	err = c.RecoverAfterRestart(7, c)
	if err != nil {
		log.Warn("[NOVAYA GAZETA CRAWLER] Recovering: " + err.Error())
	}

	err = wd.Quit()
	if err != nil {
		log.Error("[NOVAYA GAZETA CRAWLER] Quit browser: " + err.Error())
	}
	err = s.Stop()
	if err != nil {
		log.Error("[NOVAYA GAZETA CRAWLER] Stop selenium: " + err.Error())
	}

	for {
		log.Debug("[NOVAYA GAZETA CRAWLER] Awake")

		wd, s, err := c.initWebDriver()
		if err != nil {
			log.Fatal("[NOVAYA GAZETA CRAWLER] " + err.Error())
		}
		novayaGazetaWD = wd

		err = c.FindAvailableTextStreams()
		if err != nil {
			log.Error("[NOVAYA GAZETA CRAWLER] Crawling failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("novaya gazeta").Inc()
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("novaya gazeta").Inc()
		}

		err = wd.Quit()
		if err != nil {
			log.Error("[NOVAYA GAZETA CRAWLER] Quit browser: " + err.Error())
		}
		err = s.Stop()
		if err != nil {
			log.Error("[NOVAYA GAZETA CRAWLER] Stop selenium: " + err.Error())
		}

		time.Sleep(sleepTime)
	}
}

func (c NovayaGazetaCrawler) FindAvailableTextStreams() error {
	fetchURL := "https://novayagazeta.ru/features"

	err := novayaGazetaWD.Get(fetchURL)
	if err != nil {
		return err
	}

	articles, err := novayaGazetaWD.FindElements(selenium.ByTagName, "article")
	if err != nil {
		return err
	}

	var links []string
	var names []string

	for _, s := range articles {
		hrefs, _ := s.FindElements(selenium.ByTagName, "a")
		for _, a := range hrefs {
			href, _ := a.GetAttribute("href")

			if strings.Contains(href, "/articles/") {
				href = "https://novayagazeta.ru" + href
				h2, _ := s.FindElement(selenium.ByTagName, "h2")
				name, _ := h2.Text()

				log.Debug("[NOVAYA GAZETA CRAWLER] ", name, " -> ", href)
				links = append(links, href)
				names = append(names, name)
			}
		}
	}

	for i, href := range links {
		isOnline, err := c.checkArticleIsOnline(href)
		if err != nil {
			log.Error("[NOVAYA GAZETA CRAWLER] " + err.Error())
			continue
		}

		log.Debug("[NOVAYA GAZETA CRAWLER] Link: ", href, ", online: ", isOnline)

		if isOnline == true {
			exists, err := dao.CheckTextStreamExistByURL(href)
			if err != nil {
				log.Error(err)
				continue
			}

			if !exists {
				textStream := entities.TextStream{
					SourceID: 7,
					URL:      href,
					Name:     names[i],
					IsActive: true,
				}

				err := dao.AddTextStream(&textStream)
				if err != nil {
					log.Error("[NOVAYA GAZETA CRAWLER] " + err.Error())
				} else {
					go c.StartTextStream(textStream)
				}
			}
		}
	}

	return nil
}

func (c NovayaGazetaCrawler) StartTextStream(textStream entities.TextStream) {
	log.Info("[NOVAYA GAZETA CRAWLER] Starting text stream on " + textStream.URL)
	sleepTime := config.GetConfig().TextStreamUpdateInterval
	maxEmptyFetches := config.GetConfig().TextStreamMaxEmptyFetches

	emptyFetchesCounter := 0
	for {
		log.Debug("[NOVAYA GAZETA CRAWLER] Awake")

		fetched, err := c.fetchMessages(textStream)
		if err != nil {
			log.Error("[NOVAYA GAZETA CRAWLER] Crawler failed: " + err.Error())
			utils.GetMetrics().Scrapings.Failed.WithLabelValues("novaya gazeta").Inc()
			emptyFetchesCounter++
		} else {
			utils.GetMetrics().Scrapings.Done.WithLabelValues("novaya gazeta").Inc()
			if fetched > 0 {
				log.Debug("[NOVAYA GAZETA CRAWLER] New messages: " + strconv.Itoa(fetched))
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
			log.Error("[NOVAYA GAZETA CRAWLER] " + err.Error())
			time.Sleep(sleepTime)
		} else {
			finished = true
		}
	}

	log.Info("[NOVAYA GAZETA CRAWLER] Text stream finished: " + textStream.URL)
}

func (c NovayaGazetaCrawler) fetchMessages(textStream entities.TextStream) (int, error) {
	log.Debug("[NOVAYA GAZETA CRAWLER] Fetching messages")
	newMessages := 0

	wd, s, err := c.initWebDriver()
	if err != nil {
		return newMessages, err
	}

	defer s.Stop()
	defer wd.Quit()

	err = wd.Get(textStream.URL)
	if err != nil {
		return newMessages, err
	}

	articles, err := wd.FindElements(selenium.ByTagName, "article")
	if err != nil {
		return newMessages, err
	}

	for i, j := 0, len(articles)-1; i < j; i, j = i+1, j-1 {
		articles[i], articles[j] = articles[j], articles[i]
	}

	for _, s := range articles {
		tag, _ := s.FindElement(selenium.ByTagName, "time")
		t, _ := tag.Text()

		tag, _ = s.FindElement(selenium.ByTagName, "div")
		content, _ := tag.Text()

		tag, _ = s.FindElement(selenium.ByTagName, "h3")
		title, _ := tag.Text()

		exists, err := dao.CheckMessageExist(content, textStream.ID)
		if err != nil {
			log.Error("[NOVAYA GAZETA CRAWLER] " + err.Error())
		} else {
			if !exists {
				newMessages++
				message := entities.Message{
					TextStreamID: textStream.ID,
					Title:        &title,
					Body:         content,
					Time:         &t,
				}
				err := c.SaveMessage(message)
				if err != nil {
					log.Error("[NOVAYA GAZETA CRAWLER] " + err.Error())
				}
			}
		}
	}

	return newMessages, nil
}

func (c NovayaGazetaCrawler) checkArticleIsOnline(link string) (bool, error) {
	err := novayaGazetaWD.Get(link)
	if err != nil {
		return false, err
	}

	defer novayaGazetaWD.Back()

	links, err := novayaGazetaWD.FindElements(selenium.ByTagName, "a")
	if err != nil {
		return false, err
	}

	for _, l := range links {
		t, _ := l.Text()
		if t == "ОНЛАЙН" {
			return true, nil
		}
	}

	return false, nil
}
