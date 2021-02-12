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
	MeduzaAPI      time.Duration
	NewsAPI        time.Duration
	NewscatcherAPI time.Duration
	EchoMskCrawler time.Duration
}

type ScrappingEnabled struct {
	MeduzaAPI      bool
	NewsAPI        bool
	NewscatcherAPI bool
	EchoMskCrawler bool
}

var config *Config

func init() {
	defaultConfig := Config{
		ScrappingIntervals: ScrappingIntervals{
			MeduzaAPI:      1 * 60,
			NewsAPI:        4 * 60,
			NewscatcherAPI: 5 * 60,
			EchoMskCrawler: 1 * 60,
		},
		ScrappingEnabled: ScrappingEnabled{
			MeduzaAPI:      true,
			NewsAPI:        false,
			NewscatcherAPI: false,
			EchoMskCrawler: true,
		},
		TextStreamUpdateInterval:  1*60,
		TextStreamMaxEmptyFetches: 3,
	}

	config = &defaultConfig
}

func GetConfig() *Config {
	return config
}
