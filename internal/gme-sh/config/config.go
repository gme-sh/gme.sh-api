package config

import "os"

// DatabaseConfig -> Config for PersistentDatabase implementations
type DatabaseConfig struct {
	Backend           string
	EnableSharedCache bool
	Mongo             *MongoConfig
	Redis             *RedisConfig
	BBolt             *BBoltConfig
	Maria             *MariaConfig
}

// Config --> Config for Database implementations
type Config struct {
	Database *DatabaseConfig
}

// MongoConfig -> Config for MongoDB implementation
type MongoConfig struct {
	ApplyURI           string
	Database           string
	ShortURLCollection string
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
	Path                  string
	FileMode              os.FileMode
	ShortedURLsBucketName string
}

// MariaConfig -> Config for Maria Imlementation
type MariaConfig struct {
	Addr        string
	User        string
	Password    string
	DBName      string
	TablePrefix string
}
