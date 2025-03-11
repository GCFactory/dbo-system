package config

type RegistrationOrch struct {
	Host             string `yaml:"Host"`
	Port             string `yaml:"Port"`
	Retry            int    `yaml:"Retry"`
	TimeWaitRetry    int    `yaml:"TimeWaitRetry"`
	TimeWaitResponse int    `yaml:"TimeWaitResponse"`
}
