package daemon

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/models/scrappers"
	log "github.com/sirupsen/logrus"
)

func StartAllScrappers() {
	log.Debug("[DAEMON] Starting scrappers...")
	c := config.GetConfig()

	if c.ScrappingEnabled.MeduzaAPI == true {
		medusa := scrappers.MeduzaScrapper{}
		go medusa.RunForever()
	}

	if c.ScrappingEnabled.NewsAPI == true {
		newsAPI := scrappers.NewsAPIScrapper{}
		go newsAPI.RunForever()
	}

	if c.ScrappingEnabled.NewscatcherAPI == true {
		newscatcherAPI := scrappers.NewscatcherAPIScrapper{}
		go newscatcherAPI.RunForever()
	}

	if c.ScrappingEnabled.OnoMediaAPI == true {
		onoMedia := scrappers.OnoMediaScrapper{}
		go onoMedia.RunForever()
	}
}
