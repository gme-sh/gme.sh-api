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
<<<<<<< Updated upstream
=======

// BBoltConfig -> Config for BBolt implementation
type BBoltConfig struct {
	path string
}

// MariaConfig -> Config for Maria Imlementation
type MariaConfig struct {
	user     string
	password string
	dbname   string
}
>>>>>>> Stashed changes
