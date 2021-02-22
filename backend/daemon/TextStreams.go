package daemon

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/models/crawlers"
	log "github.com/sirupsen/logrus"
)

func StartAllCrawlers() {
	log.Debug("[DAEMON] Starting crawlers...")
	c := config.GetConfig()

	if c.ScrappingEnabled.EchoMskCrawler == true {
		echoMskCrawler := crawlers.EchoMskCrawler{}
		go echoMskCrawler.RunForever()
	}

	if c.ScrappingEnabled.ZonaMediaCrawler == true {
		zonaMediaCrawler := crawlers.ZonaMediaCrawler{}
		go zonaMediaCrawler.RunForever()
	}

	if c.ScrappingEnabled.NovayaGazetaCrawler == true {
		novayaGazetaCrawler := crawlers.NovayaGazetaCrawler{}
		go novayaGazetaCrawler.RunForever()
	}
}
