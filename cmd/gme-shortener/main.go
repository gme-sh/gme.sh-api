package main

import (
	"encoding/json"
	"fmt"
	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/db/heartbeat"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

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

func Test(database db.Database) {
	log.Println("test:", database)
}

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

	config.FromEnv(dbcfg)

	// Update config from environment
	// Get mongo from environment
	if mdbs := os.Getenv("MONGODB_STRING"); mdbs != "" {
		dbcfg.Mongo.ApplyURI = mdbs
	}

	// Load persistentDB
	var persistentDB db.PersistentDatabase
	var tempDB db.TemporaryDatabase

	switch strings.ToLower(dbcfg.Backend) {
	case "mongo":
		log.Println("👉 Using MongoDB as backend")
		persistentDB = db.Must(db.NewMongoDatabase(dbcfg.Mongo.ApplyURI)).(db.PersistentDatabase)
		break
	case "maria":
		log.Println("👉 Using MariaDB as backend")
		persistentDB = db.Must(db.NewMariaDB(*dbcfg.Maria)).(db.PersistentDatabase)
		break
	case "bbolt":
		log.Println("👉 Using BBolt as backend")
		persistentDB = db.Must(db.NewBBoltDatabase(dbcfg.BBolt.Path)).(db.PersistentDatabase)
		break
	case "redis":
		log.Println("👉 Using Redis as backend")
		redisDB := db.Must(db.NewRedisDatabase(*dbcfg.Redis))

		persistentDB = redisDB.(db.PersistentDatabase)
		tempDB = redisDB.(db.TemporaryDatabase)
		break
	default:
		log.Fatalln("🚨 Invalid persistentDB backend: '", dbcfg.Backend, "'")
		return
	}

	var redisClient *redis.Client = nil

	// Load redis
	if dbcfg.Redis.Use {
		log.Println("👉 Using redis as temporary database")

		if tempDB == nil {
			tempDB = db.Must(db.NewRedisDatabase(*dbcfg.Redis)).(db.TemporaryDatabase)
		}
	}

	var hb chan bool
	if tempDB != nil {
		hb = heartbeat.CreateHeartbeatService(tempDB)
	} else {
		hb = make(chan bool, 1)
	}

	// Create example data
	log.Println("└ Adding dummy data to persistentDB")
	link := short.ShortURL{
		ID:           "ddg",
		FullURL:      "https://duckduckgo.com/",
		CreationDate: time.Now(),
	}
	log.Println("Saving to persistentDB result:", persistentDB.SaveShortenedURL(link))

	server := web.NewWebServer(persistentDB, redisClient)
	go server.Start()

	log.Println("WebServer is (hopefully) up and running")
	log.Println("Press CTRL+C to exit gracefully")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	hb <- true

	// after CTRL+c
	log.Println("Shutting down webserver")
}
