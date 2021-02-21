package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/config"
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/db"
	"github.com/full-stack-gods/gme.sh-api/internal/gme-sh/web"
)

const (
	// This Banner is displayed when the API is started
	Banner = `
                                         /$$                               /$$
                                        | $$                              | $$
 ██████╗ ███╗   ███╗███████╗   /$$$$$$$ | $$$$$$$    /$$$$$$   /$$$$$$   /$$$$$$
██╔════╝ ████╗ ████║██╔════╝  /$$_____/ | $$__  $$  /$$__  $$ /$$__  $$ |_  $$_/
██║  ███╗██╔████╔██║█████╗   |  $$$$$$  | $$  \ $$ | $$  \ $$ | $$  \__/   | $$
██║   ██║██║╚██╔╝██║██╔══╝    \____  $$ | $$  | $$ | $$  | $$ | $$         | $$ /$$
╚██████╔╝██║ ╚═╝ ██║███████╗  /$$$$$$$/ | $$  | $$ |  $$$$$$/ | $$         |  $$$$/
 ╚═════╝ ╚═╝     ╚═╝╚══════╝ |_______/  |__/  |__/  \______/  |__/          \____/`

	// Version of the backend
	Version = "1.0.0-alpha" // semantic
)

func main() {
	fmt.Println(Banner)
	fmt.Println("Starting $GMEshort", Version, "🚀")
	fmt.Println()

	//// Config
	log.Println("└ Loading config")
	cfg := config.LoadConfig()
	if cfg == nil {
		return
	}
	dbcfg := cfg.Database
	////

	//// Database
	// persistentDB is used to store short urls (persistent, obviously)
	var persistentDB db.PersistentDatabase
	// tempDB is used to store temporary information for short urls (eg. stats, caching)
	var tempDB db.TemporaryDatabase

	if strings.ToLower(dbcfg.Backend) == "redis" {
		log.Println("👉 Using Redis as backend")
		redisDB := db.Must(db.NewRedisDatabase(dbcfg.Redis))

		persistentDB = redisDB.(db.PersistentDatabase)
		tempDB = redisDB.(db.TemporaryDatabase)
	}
	if dbcfg.Redis.Use {
		log.Println("👉 Using redis as temporary database")

		if tempDB == nil {
			tempDB = db.Must(db.NewRedisDatabase(dbcfg.Redis)).(db.TemporaryDatabase)
		}
	}

	var cache db.DBCache

	if dbcfg.EnableSharedCache {
		if tempDB == nil {
			log.Fatalln("Cannot enable shared cache when no temporary database is set! (e. g. Redis)")
			return
		}
		cache = db.NewSharedCache(tempDB)

		// subscribe to shared cache
		// e. g. Redis Pub-Sub
		go func() {
			log.Println("SCACHE :: Subscribing to redis channels ...")
			if err := cache.(*db.SharedCache).Subscribe(); err != nil {
				log.Println("SCACHE :: Error:", err)
			}
		}()
	} else {
		cache = db.NewLocalCache()
	}

	switch strings.ToLower(dbcfg.Backend) {
	case "mongo":
		log.Println("👉 Using MongoDB as backend")
		persistentDB = db.Must(db.NewMongoDatabase(dbcfg.Mongo, cache)).(db.PersistentDatabase)
		break
	case "maria":
		log.Println("👉 Using MariaDB as backend")
		persistentDB = db.Must(db.NewMariaDB(dbcfg.Maria, cache)).(db.PersistentDatabase)
		break
	case "bbolt":
		log.Println("👉 Using BBolt as backend")
		persistentDB = db.Must(db.NewBBoltDatabase(dbcfg.BBolt, cache)).(db.PersistentDatabase)
		break
	case "redis":
		break
	default:
		log.Fatalln("🚨 Invalid persistentDB backend: '", dbcfg.Backend, "'")
		return
	}

	var hb chan bool
	if tempDB != nil {
		hb = db.CreateHeartbeatService(tempDB)
	} else {
		hb = make(chan bool, 1)
	}
	////

	//// Web-Server
	server := web.NewWebServer(persistentDB, tempDB, cfg)
	go server.Start()
	////

	log.Println("WebServer is (hopefully) up and running")
	log.Println("Press CTRL+C to exit gracefully")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	hb <- true

	// after CTRL+c
	log.Println("Shutting down webserver")
}
