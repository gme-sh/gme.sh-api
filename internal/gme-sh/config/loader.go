package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var (
	// ConfigPath is "config.toml" by default
	ConfigPath = "config.toml"
)

func init() {
	if val := os.Getenv("GME_CONFIG_PATH"); val != "" {
		ConfigPath = val
	}
}

// LoadConfig loads a config if the config exists, otherwise creates a default config
func LoadConfig() *Config {
	var cfg *Config

	// check if config file exists
	// if not, create a default config
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Println("└   Creating default config")
		if err := CreateDefault(); err != nil {
			log.Fatalln("Error creating config:", err)
			return nil
		}
	}

	// decode config from file "config.toml"
	if _, err := toml.DecodeFile(ConfigPath, &cfg); err != nil {
		log.Fatalln("Error decoding file:", err)
		return nil
	}

	log.Println("  ├ Dry-Redirect:", cfg.DryRedirect)
	log.Println("  ├ Web-Addr:", cfg.WebServer.Addr)
	log.Println("  └ Blocked-Hosts:", cfg.BlockedHosts)

	errs := FromEnv(cfg)
	for i, e := range errs {
		if e != nil {
			log.Fatalln("ERROR #", i, "loading config from env:", e)
			return nil
		}
	}

	return cfg
}
