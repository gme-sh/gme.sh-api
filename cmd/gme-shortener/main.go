package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/full-stack-gods/GMEshortener/internal/gme-shortener/config"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/db"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/short"
	"github.com/full-stack-gods/GMEshortener/pkg/gme-shortener/web"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Get mongo from environment
	if mdbs := os.Getenv("MONGODB_STRING"); mdbs != "" {
		cfg.Mongo.ApplyURI = mdbs
	}

	// connect to database
	log.Println("└ Connecting to database")
	database, err := db.NewMongoDatabase(cfg.Mongo.ApplyURI)
	if err != nil {
		log.Fatalln("Error connecting:", err)
		return
	}

	// Create example data
	log.Println("└ Adding dummy data to database")
	link := short.ShortURL{
		ID:           "ddg",
		FullURL:      "https://duckduckgo.com/",
		CreationDate: time.Now(),
	}
	log.Println("Saving to database result:", database.SaveShortenedURL(link))

	server := web.NewWebServer(database)
	go server.Start()

	log.Println("WebServer is (hopefully) up and running")
	log.Println("Press CTRL+C to exit gracefully")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// after CTRL+c
	log.Println("Shutting down webserver")
}
