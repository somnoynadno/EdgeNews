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
			MeduzaAPI:      5 * 60,
			NewsAPI:        10 * 60,
			NewscatcherAPI: 15 * 60,
			EchoMskCrawler: 3 * 60,
		},
		ScrappingEnabled: ScrappingEnabled{
			MeduzaAPI:      true,
			NewsAPI:        true,
			NewscatcherAPI: true,
			EchoMskCrawler: true,
		},
		TextStreamUpdateInterval:  1*60,
		TextStreamMaxEmptyFetches: 12*60*60,
	}

	config = &defaultConfig
}

func GetConfig() *Config {
	return config
}
