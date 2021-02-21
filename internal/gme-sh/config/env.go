package config

import (
	"github.com/qiangxue/go-env"
	"log"
)

// FromEnv -> Please correct this function
// taken from r/programminghorror
func FromEnv(cfg *Config) (err []error) {
	loader := env.New("GME_", log.Printf)
	err = append(err, loader.Load(cfg))

	// Database
	err = append(err, loader.Load(cfg.Database))
	err = append(err, loader.Load(cfg.Database.Mongo))
	err = append(err, loader.Load(cfg.Database.Redis))
	err = append(err, loader.Load(cfg.Database.BBolt))

	// Web Server
	err = append(err, loader.Load(cfg.WebServer))
	return
}
