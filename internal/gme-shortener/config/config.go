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
<<<<<<< HEAD
=======
<<<<<<< Updated upstream
=======
>>>>>>> main

// BBoltConfig -> Config for BBolt implementation
type BBoltConfig struct {
	Path string
}
<<<<<<< HEAD
=======

// MariaConfig -> Config for Maria Imlementation
type MariaConfig struct {
	user     string
	password string
	dbname   string
}
>>>>>>> Stashed changes
>>>>>>> main
