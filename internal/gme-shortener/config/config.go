package config

type Config struct {
	Mongo MongoConfig
}

type MongoConfig struct {
	ApplyURI string `json:"apply_uri"`
}
