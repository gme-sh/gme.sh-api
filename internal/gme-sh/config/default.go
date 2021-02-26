package config

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

// CreateDefault -> create default config
func CreateDefault() (err error) {
	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	err = e.Encode(Config{
		DryRedirect: false,
		BlockedHosts: &BlockedHosts{
			Hosts: []string{"gme.sh"},
		},
		Backends: &BackendConfig{
			PersistentBackend: "Mongo",
			StatsBackend:      "Redis",
			PubSubBackend:     "Redis",
			CacheBackend:      "Shared",
		},
		Database: &DatabaseConfig{
			Mongo: &MongoConfig{
				ApplyURI:           "mongodb://localhost:27017",
				Database:           "stonksdb",
				ShortURLCollection: "stonks-url-collection",
			},
			Redis: &RedisConfig{
				Addr:     "localhost",
				Password: "",
				DB:       0,
			},
			BBolt: &BBoltConfig{
				Path:                  "dbgoesbrr.rr",
				FileMode:              0666,
				ShortedURLsBucketName: "stonks-url-bucket",
			},
			Maria: &MariaConfig{
				Addr:        "localhost",
				User:        "root",
				Password:    "",
				DBName:      "stonks",
				TablePrefix: "stonks_",
			},
		},
		WebServer: nil,
	})
	if err != nil {
		log.Fatalln("Error encoding default config:", err)
		return
	}

	if err = ioutil.WriteFile("config.toml", buf.Bytes(), 0666); err != nil {
		log.Fatalln("Error saving default config:", err)
		return
	}
	return
}
