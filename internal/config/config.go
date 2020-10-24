package config

import (
	"os"
	"strconv"
)

type (
	HTTPConfig struct {
		Port int
	}

	StorageConfig struct {
		Host string
		Port int
	}

	Config struct {
		HTTP    HTTPConfig
		Storage StorageConfig
	}
)

func MustConfigure() *Config {
	var httpPort int
	var storageHost string
	var storagePort int

	if len(os.Getenv("HTTP_PORT")) > 0 {
		httpPort, _ = strconv.Atoi(os.Getenv("HTTP_PORT"))
	}

	if len(os.Getenv("STORAGE_HOST")) > 0 {
		storageHost = os.Getenv("STORAGE_HOST")
	}

	if len(os.Getenv("STORAGE_PORT")) > 0 {
		storagePort, _ = strconv.Atoi(os.Getenv("STORAGE_PORT"))
	}

	return &Config{
		HTTPConfig{Port: httpPort},
		StorageConfig{
			Host: storageHost,
			Port: storagePort,
		},
	}
}
