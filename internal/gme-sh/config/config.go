package config

import "os"

// Config -> Config for Database implementations
type Config struct {
	DryRedirect  bool          `env:"DRY_REDIRECT"`
	BlockedHosts *BlockedHosts `env:"BLOCKED_HOSTS"`
	Backends     *BackendConfig
	Database     *DatabaseConfig
	WebServer    *WebServerConfig
}

type BackendConfig struct {
	PersistentBackend string `env:"PERSISTENT_BACKEND"`
	StatsBackend      string `env:"STATS_BACKEND"`
	PubSubBackend     string `env:"PUBSUB_BACKEND"`
	CacheBackend      string `env:"CACHE_BACKEND"`
}

// DatabaseConfig -> Config for PersistentDatabase implementations
type DatabaseConfig struct {
	Mongo *MongoConfig
	Redis *RedisConfig
	BBolt *BBoltConfig
	Maria *MariaConfig
}

// WebServerConfig -> Config for web.WebServer
type WebServerConfig struct {
	Addr string `env:"WEB_ADDR"`
}

// MongoConfig -> Config for MongoDB implementation
type MongoConfig struct {
	ApplyURI           string `env:"MDB_APPLY_URI"`
	Database           string `env:"MDB_DATABASE"`
	ShortURLCollection string `env:"MDB_COLLECTION_SHORT_URLS"`
}

// RedisConfig -> Config for Redis implementation
type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASS"`
	DB       int    `env:"REDIS_DATABASE"`
}

// BBoltConfig -> Config for BBolt implementation
type BBoltConfig struct {
	Path                  string      `env:"BBOLT_PATH"`
	FileMode              os.FileMode `env:"BBOLT_FILE_MODE"`
	ShortedURLsBucketName string      `env:"BBOLT_BUCKET_SHORT_URLS"`
}

// MariaConfig -> Config for Maria Imlementation
type MariaConfig struct {
	Addr        string `env:"MARIA_ADDR"`
	User        string `env:"MARIA_USER"`
	Password    string `env:"MARIA_PASS"`
	DBName      string `env:"MARIA_DATABASE"`
	TablePrefix string `env:"MARIA_TABLE_PREFIX"`
}
