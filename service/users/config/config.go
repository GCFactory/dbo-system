package config

import (
	platformConfig "github.com/GCFactory/dbo-system/platform/config"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

// App config struct
type Config struct {
	App        map[interface{}]interface{} `yaml:"app"`
	Env        string                      `yaml:"env"`
	Version    string                      `yaml:"version"`
	HTTPServer HTTPServerConfig            `yaml:"http-server,omitempty" mapstructure:"http-server"`

	Logger  platformConfig.Logger `yaml:"logger"`
	Jaeger  Jaeger                `yaml:"jaeger"`
	Metrics Metrics               `yaml:"metrics"`
	Docs    Docs

	KafkaConsumer KafkaConsumer `yaml:"kafkaConsumer"`
	KafkaProducer KafkaProducer `yaml:"kafkaConsumer"`

	Postgres platformConfig.PostgresConfig `yaml:"postgres,omitempty"`
	AWS      AWS                           `yaml:"aws,omitempty"`

	Cookie  Cookie  `yaml:"cookie,omitempty"`
	Session Session `yaml:"session,omitempty"`
}

// Swagger configuration
type Docs struct {
	Enable bool
	Title  string
	Prefix string
}

// HTTP Server config struct
type HTTPServerConfig struct {
	Port              string
	PprofPort         string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

// AWS S3
type AWS struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	UseSSL        bool
	MinioEndpoint string `yaml:"minio-endpoint,omitempty"`
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AllowEmptyEnv(false)
	//v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
