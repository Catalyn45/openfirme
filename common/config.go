package common

import (
	"encoding/json"
	"os"
)

type WebsiteGeneratorConfig struct {
	Enabled bool
	TemplatesDirectory string
	OutputDirectory string
}

type DosareJuridiceClientConfig struct {
	RequstTimeoutInSeconds int
}

type CacheConfig struct {
	Enabled bool

	DefaultCacheTimeInMinutes int
	CleanupCacheIntervalTimeInMinutes int

	DosareJuridiceCacheTimeInMinutes int

	MinPagedCacheTimeInMinutes int
	MaxPagedCacheTimeInMinutes int

	MinSearchCacheTimeInMinutes int
	MaxSearchCacheTimeInMinutes int

	CacheSaveEnabled bool
	CacheSaveFilePath string
	CacheSaveIntervalInMinutes int
}

type DBConfig struct {
	DefaultTimeoutInSeconds int
	SearchFirmeTimeoutInSeconds int

	LogQueries bool
}

type ServerConfig struct {
	Host string
	Port int

	PublicDirectory string
}

type Config struct {
	DataDirectory string
	DBFilePath string

	ServerConfig ServerConfig
	DBConfig DBConfig
	CacheConfig CacheConfig
	DosareJuridiceClientConfig DosareJuridiceClientConfig
	WebsiteGeneratorConfig WebsiteGeneratorConfig
}

const configOverridePath = "./config.json"

var config *Config
func init() {
	config = &Config {
		DataDirectory: "./data",
		DBFilePath: "./foo.db",

		ServerConfig: ServerConfig {
			Host: "127.0.0.1",
			Port: 8080,
			PublicDirectory: "./public",
		},

		WebsiteGeneratorConfig: WebsiteGeneratorConfig {
			Enabled: false,
			TemplatesDirectory: "./templates",
			OutputDirectory: "./public",
		},

		DBConfig: DBConfig {
			DefaultTimeoutInSeconds: 15,
			SearchFirmeTimeoutInSeconds: 10,
			LogQueries: false,
		},

		CacheConfig : CacheConfig {
			Enabled: true,

			DefaultCacheTimeInMinutes: 20,
			CleanupCacheIntervalTimeInMinutes: 10,

			DosareJuridiceCacheTimeInMinutes: 30,

			MinPagedCacheTimeInMinutes: 20,
			MaxPagedCacheTimeInMinutes: 40,

			MinSearchCacheTimeInMinutes: 5,
			MaxSearchCacheTimeInMinutes: 10,

			CacheSaveEnabled: true,
			CacheSaveFilePath: "./cache.gob",
			CacheSaveIntervalInMinutes: 5,
		},

		DosareJuridiceClientConfig: DosareJuridiceClientConfig {
			RequstTimeoutInSeconds: 20,
		},
	}

	// check if there is json file to override configs
	_, err := os.Stat(configOverridePath)
	if err != nil {
		return
	}

	fileContent, err := os.ReadFile(configOverridePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileContent, config)
	if err != nil {
		panic(err)
	}
}