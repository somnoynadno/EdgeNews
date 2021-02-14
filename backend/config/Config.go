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
	MeduzaAPI       time.Duration
	NewsAPI         time.Duration
	NewscatcherAPI  time.Duration
	OnoMediaAPI     time.Duration
	EchoMskCrawler  time.Duration
	KulichkiCrawler time.Duration
}

type ScrappingEnabled struct {
	MeduzaAPI       bool
	NewsAPI         bool
	NewscatcherAPI  bool
	OnoMediaAPI     bool
	EchoMskCrawler  bool
	KulichkiCrawler bool
}

var config *Config

func init() {
	defaultConfig := Config{
		ScrappingIntervals: ScrappingIntervals{
			MeduzaAPI:       1 * 60,
			NewsAPI:         3 * 60,
			NewscatcherAPI:  6 * 60,
			OnoMediaAPI:     1 * 60,
			EchoMskCrawler:  3 * 60,
			KulichkiCrawler: 1 * 60,
		},
		ScrappingEnabled: ScrappingEnabled{
			MeduzaAPI:       true,
			NewsAPI:         true,
			NewscatcherAPI:  true,
			OnoMediaAPI:     true,
			EchoMskCrawler:  true,
			KulichkiCrawler: true,
		},
		TextStreamUpdateInterval:  1 * 60,
		TextStreamMaxEmptyFetches: 12 * 60,
	}

	config = &defaultConfig
}

func GetConfig() *Config {
	return config
}
