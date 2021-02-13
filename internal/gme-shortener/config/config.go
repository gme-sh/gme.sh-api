package config

type Config struct {
	Mongo *MongoConfig
	Redis *RedisConfig
}

type MongoConfig struct {
	ApplyURI string `json:"apply_uri"`
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
