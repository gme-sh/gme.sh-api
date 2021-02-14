package config

type DatabaseConfig struct {
	Backend string
	Mongo   *MongoConfig
	Redis   *RedisConfig
	BBolt   *BBoltConfig
	Maria   *MariaConfig
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

// MariaConfig -> Config for Maria Imlementation
type MariaConfig struct {
	Addr        string
	User        string
	Password    string
	DBName      string
	TablePrefix string
}
