package config

import (
	"time"
)

type Config struct {
	ScrappingIntervals ScrappingIntervals
	ScrappingEnabled   ScrappingEnabled
}

type ScrappingIntervals struct {
	Meduza         time.Duration
	NewsAPI        time.Duration
	NewscatcherAPI time.Duration
}

type ScrappingEnabled struct {
	Meduza         bool
	NewsAPI        bool
	NewscatcherAPI bool
}

var config *Config

func init() {
	defaultConfig := Config{
		ScrappingIntervals: ScrappingIntervals{
			Meduza:         1 * 60,
			NewsAPI:        4 * 60,
			NewscatcherAPI: 5 * 60,
		},
		ScrappingEnabled: ScrappingEnabled{
			Meduza:         true,
			NewsAPI:        false,
			NewscatcherAPI: false,
		},
	}

	config = &defaultConfig
}

func GetConfig() *Config {
	return config
}
