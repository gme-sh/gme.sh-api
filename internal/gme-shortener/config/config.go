package config

// Config --> Config for Database implementations
type Config struct {
	Mongo *MongoConfig
	Redis *RedisConfig
	BBolt *BBoltConfig
}

// MongoConfig -> Config for MongoDB implementation
type MongoConfig struct {
	ApplyURI string `json:"apply_uri"`
}

// RedisConfig -> Config for Redis implementation
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// BBoltConfig -> Config for BBolt implementation
type BBoltConfig struct {
	path string
}
