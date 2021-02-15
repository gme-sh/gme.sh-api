package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// FromEnv -> Please correct this function
// taken from r/programminghorror
func FromEnv(config *DatabaseConfig) {
	// MongoDB
	if mdbap := os.Getenv("MONGODB_APPLYURI"); mdbap != "" {
		config.Mongo.ApplyURI = mdbap
	}
	if val := os.Getenv("MONGO_DATABASE"); val != "" {
		config.Mongo.Database = val
	}
	if val := os.Getenv("MONGO_SHORT_URL_COLLECTION"); val != "" {
		config.Mongo.ShortURLCollection = val
	}

	// Redis
	if rduse := os.Getenv("REDIS_USE"); rduse != "" {
		config.Redis.Use = strings.ToLower(rduse) == "true"
	}
	if rdaddr := os.Getenv("REDIS_ADDR"); rdaddr != "" {
		config.Redis.Addr = rdaddr
	}
	if rdpw := os.Getenv("REDIS_PASS"); rdpw != "" {
		config.Redis.Password = rdpw
	}
	if rddb := os.Getenv("REDIS_DB"); rddb != "" {
		atoi, err := strconv.Atoi(rddb)
		if err == nil {
			config.Redis.DB = atoi
		} else {
			log.Fatalln("🚨 REDIS_DB: Invalid number (int):", rddb, "(error):", err)
		}
	}

	// BBolt
	if bbp := os.Getenv("BBOLT_PATH"); bbp != "" {
		config.BBolt.Path = bbp
	}
	if val := os.Getenv("BBOLT_FILE_MODE"); val != "" {
		i, err := strconv.Atoi(val)
		if err == nil {
			config.BBolt.FileMode = os.FileMode(i)
		} else {
			log.Fatalln("🚨 BBOLT_FILE_MODE: Invalid number (int):", val, "(error):", err)
		}
	}
	if val := os.Getenv("BBOLT_SHORTED_URL_BUCKET_NAME"); val != "" {
		config.BBolt.ShortedURLsBucketName = val
	}

	// MariaDB
	if val := os.Getenv("MARIADB_ADDR"); val != "" {
		config.Maria.Addr = val
	}
	if val := os.Getenv("MARIADB_USER"); val != "" {
		config.Maria.User = val
	}
	if val := os.Getenv("MARIADB_PASS"); val != "" {
		config.Maria.Password = val
	}
	if val := os.Getenv("MARIADB_DB"); val != "" {
		config.Maria.DBName = val
	}
	if val := os.Getenv("MARIADB_TABLE_PREFIX"); val != "" {
		config.Maria.TablePrefix = val
	}
}
