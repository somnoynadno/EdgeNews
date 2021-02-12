package main

import (
	"EdgeNews/backend/config"
	"EdgeNews/backend/daemon"
	"EdgeNews/backend/db"
	"EdgeNews/backend/server"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)

	if os.Getenv("ENV") == "PRODUCTION" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	_ = config.GetConfig()

	con := db.GetDB()
	err := con.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	go daemon.StartAllScrappers()
	go daemon.StartAllCrawlers()

	server.InitRouter()
}