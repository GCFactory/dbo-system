package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

// App config struct
type Config struct {
	Logger  Logger  `yaml:"logger"`
	Jaeger  Jaeger  `yaml:"jaeger"`
	Metrics Metrics `yaml:"metrics"`

	Env        string                      `yaml:"env"`
	App        map[interface{}]interface{} `yaml:"app"`
	HTTPServer HTTPServerConfig            `yaml:"http-server,omitempty"`

	Postgres PostgresConfig `yaml:"postgres,omitempty"`
	Redis    RedisConfig    `yaml:"redis,omitempty"`
	MongoDB  MongoDB        `yaml:"mongodb,omitempty"`
	AWS      AWS            `yaml:"aws,omitempty"`

	Cookie  Cookie  `yaml:"cookie,omitempty"`
	Session Session `yaml:"session,omitempty"`
}

// HTTP Server config struct
type HTTPServerConfig struct {
	AppVersion        string
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
			return nil, errors.New("config file not found")
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
