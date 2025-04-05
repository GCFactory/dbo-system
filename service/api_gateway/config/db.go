package config

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// Redis config
type RedisConfig struct {
	RedisAddr     string `yaml:"Host"`
	RedisPassword string `yaml:"RedisPassword"`
	MaxRetries    int    `yaml:"MaxRetries"`
	User          string `yaml:"User"`
	DB            int    `yaml:"DbId"`
	DialTimeout   int    `yaml:"DialTimeout"` // таймаут для установки новых соединений
	Timeout       int    `yaml:"Timeout"`     // таймаут для чтения и записи
}

// MongoDB config
type MongoDB struct {
	MongoURI string
}
