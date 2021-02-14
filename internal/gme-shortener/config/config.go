package config

type DatabaseConfig struct {
	Backend string
	Mongo   *MongoConfig
	Redis   *RedisConfig
	BBolt   *BBoltConfig
}

// Config --> Config for Database implementations
type Config struct {
	Database *DatabaseConfig
}

// MongoConfig -> Config for MongoDB implementation
type MongoConfig struct {
	ApplyURI string
}

// RedisConfig -> Config for Redis implementation
type RedisConfig struct {
	Use      bool
	Addr     string
	Password string
	DB       int
}

// BBoltConfig -> Config for BBolt implementation
type BBoltConfig struct {
	Path string
}
