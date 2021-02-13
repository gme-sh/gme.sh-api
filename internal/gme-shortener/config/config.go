package config

type Config struct {
	Mongo *MongoConfig
}

type MongoConfig struct {
	ApplyURI string `json:"apply_uri"`
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
