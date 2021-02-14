package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/config"
	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/db"
	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/web"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
)

const (
	Banner = `
 ██████╗ ███╗   ███╗███████╗███████╗██╗  ██╗ ██████╗ ██████╗ ████████╗
██╔════╝ ████╗ ████║██╔════╝██╔════╝██║  ██║██╔═══██╗██╔══██╗╚══██╔══╝
██║  ███╗██╔████╔██║█████╗  ███████╗███████║██║   ██║██████╔╝   ██║   
██║   ██║██║╚██╔╝██║██╔══╝  ╚════██║██╔══██║██║   ██║██╔══██╗   ██║   
╚██████╔╝██║ ╚═╝ ██║███████╗███████║██║  ██║╚██████╔╝██║  ██║   ██║   
 ╚═════╝ ╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝`
	Version = "1.0.0-alpha" // semantic
)

func main() {
	fmt.Println(Banner)
	fmt.Println("Starting GMEshort", Version, "🚀")

	// load config
	log.Println("└ Loading config")

	var cfg *config.Config
	if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
		log.Fatalln("Error decoding file:", err)
		return
	}

	dbcfg := cfg.Database
	if s, err := json.Marshal(dbcfg); err != nil {
		log.Println("ERROR marshalling config:", err)
	} else {
		log.Println("config:", string(s))
	}

	// Update config from environment
	// Get mongo from environment
	if mdbs := os.Getenv("MONGODB_STRING"); mdbs != "" {
		dbcfg.Mongo.ApplyURI = mdbs
	}

	// Load database
	var database db.Database

	switch strings.ToLower(dbcfg.Backend) {
	case "mongo":
		log.Println("👉 Using MongoDB as backend")
		database = db.Must(db.NewMongoDatabase(dbcfg.Mongo.ApplyURI))
		break
	case "maria":
		log.Println("👉 Using MariaDB as backend")
		database = db.Must(db.NewMariaDB(*dbcfg.Maria))
		break
	case "bbolt":
		log.Println("👉 Using BBolt as backend")
		database = db.Must(db.NewBBoltDatabase(dbcfg.BBolt.Path))
		break
	default:
		log.Fatalln("🚨 Invalid database backend: '", dbcfg.Backend, "'")
		return
	}

	var redisClient *redis.Client = nil

	// Load redis
	if dbcfg.Redis.Use {
		log.Println("👉 Using redis")

		redisClient = db.NewRedisClient(*dbcfg.Redis)
		if res := redisClient.Set(context.TODO(), "heartbeat", 1, 0); res.Err() != nil {
			log.Fatalln("Error connecting to Redis:", res.Err())
			return
		}
	}

	// Create example data
	log.Println("└ Adding dummy data to database")
	link := short.ShortURL{
		ID:           "ddg",
		FullURL:      "https://duckduckgo.com/",
		CreationDate: time.Now(),
	}
	log.Println("Saving to database result:", database.SaveShortenedURL(link))

	server := web.NewWebServer(database, redisClient)
	go server.Start()

	log.Println("WebServer is (hopefully) up and running")
	log.Println("Press CTRL+C to exit gracefully")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// after CTRL+c
	log.Println("Shutting down webserver")
}
