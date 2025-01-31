package utils

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	} else if configPath == "local" {
		return "./config/local"
	} else {
		return configPath
	}
}
