package config

import (
	"time"
)

type Config struct {
	ScrappingIntervals        ScrappingIntervals
	ScrappingEnabled          ScrappingEnabled
	TextStreamUpdateInterval  time.Duration
	TextStreamMaxEmptyFetches int
}

type ScrappingIntervals struct {
	MeduzaAPI        time.Duration
	NewsAPI          time.Duration
	NewscatcherAPI   time.Duration
	OnoMediaAPI      time.Duration
	EchoMskCrawler   time.Duration
	KulichkiCrawler  time.Duration
	ZonaMediaCrawler time.Duration
}

type ScrappingEnabled struct {
	MeduzaAPI        bool
	NewsAPI          bool
	NewscatcherAPI   bool
	OnoMediaAPI      bool
	EchoMskCrawler   bool
	KulichkiCrawler  bool
	ZonaMediaCrawler bool
}

var config *Config

func init() {
	defaultConfig := Config{
		ScrappingIntervals: ScrappingIntervals{
			MeduzaAPI:        1 * time.Minute,
			NewsAPI:          3 * time.Minute,
			NewscatcherAPI:   6 * time.Minute,
			OnoMediaAPI:      1 * time.Minute,
			EchoMskCrawler:   3 * time.Minute,
			KulichkiCrawler:  3 * time.Minute,
			ZonaMediaCrawler: 3 * time.Minute,
		},
		ScrappingEnabled: ScrappingEnabled{
			MeduzaAPI:        false,
			NewsAPI:          false,
			NewscatcherAPI:   false,
			OnoMediaAPI:      false,
			EchoMskCrawler:   true,
			KulichkiCrawler:  false,
			ZonaMediaCrawler: true,
		},
		TextStreamUpdateInterval:  1 * time.Minute,
		TextStreamMaxEmptyFetches: 12 * 60,
	}

	config = &defaultConfig
}

func GetConfig() *Config {
	return config
}
