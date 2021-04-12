package config

import (
	"os"
)

// Config -> Config for Database implementations
type Config struct {
	DryRedirect             bool          `env:"DRY_REDIRECT"`
	BlockedHosts            *BlockedHosts `env:"BLOCKED_HOSTS"`
	ExpirationCheckInterval duration      `env:"EXPIRATION_CHECK_INTERVAL"`
	ExpirationDryRun        bool          `env:"EXPIRATION_DRY_RUN"`
	Backends                *BackendConfig
	Database                *DatabaseConfig
	WebServer               *WebServerConfig
}

type DummyConfig struct {
	DryRedirect             bool          `env:"DRY_REDIRECT"`
	BlockedHosts            *BlockedHosts `env:"BLOCKED_HOSTS"`
	ExpirationCheckInterval string        `env:"EXPIRATION_CHECK_INTERVAL"`
	ExpirationDryRun        bool          `env:"EXPIRATION_DRY_RUN"`
	Backends                *BackendConfig
	Database                *DatabaseConfig
	WebServer               *WebServerConfig
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
	Addr       string `env:"WEB_ADDR"`
	DefaultURL string `env:"DEFAULT_URL"`
}

// MongoConfig -> Config for MongoDB implementation
type MongoConfig struct {
	ApplyURI           string `env:"MDB_APPLY_URI"`
	Database           string `env:"MDB_DATABASE"`
	ShortURLCollection string `env:"MDB_COLLECTION_SHORT_URLS"`
	MetaCollection     string `env:"MDB_COLLECTION_META"`
	TplCollection      string `env:"MDB_COLLECTION_TPL"`
	PoolCollection     string `env:"MDB_POOL_COLLECTION"`
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
	MetaBucketName        string      `env:"BBOLT_BUCKET_META"`
	TplBucketName         string      `env:"BBOLT_BUCKET_TPL"`
	PoolBucketName        string      `env:"BBOLT_BUCKET_POOL"`
}

// MariaConfig -> Config for Maria Imlementation
type MariaConfig struct {
	Addr        string `env:"MARIA_ADDR"`
	User        string `env:"MARIA_USER"`
	Password    string `env:"MARIA_PASS"`
	DBName      string `env:"MARIA_DATABASE"`
	TablePrefix string `env:"MARIA_TABLE_PREFIX"`
}
