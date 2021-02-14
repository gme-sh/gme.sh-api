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
		Database: &DatabaseConfig{
			Backend: "mongo",
			Mongo: &MongoConfig{
				ApplyURI: "mongodb://localhost:27017",
			},
			Redis: &RedisConfig{
				Use:      true,
				Addr:     "localhost",
				Password: "",
				DB:       0,
			},
			BBolt: &BBoltConfig{
				Path: "dbgoesbrr.rr",
			},
			Maria: &MariaConfig{
				Addr:        "localhost",
				User:        "root",
				Password:    "",
				DBName:      "stonks",
				TablePrefix: "stonks_",
			},
		},
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