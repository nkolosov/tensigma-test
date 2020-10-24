package config

import (
	"os"
	"time"
)

type (
	GRPCConfig struct {
		Host string
		Port string
	}

	MongoConfig struct {
		Connection     string
		ConnectTimeout time.Duration
		SocketTimeout  time.Duration
		MinPoolSize    uint64
		MaxPoolSize    uint64
	}

	PipelineConfig struct {
		DownloaderWorkersCount  int
		DownloaderTempDirectory string

		ExporterWorkersCount int
	}

	Config struct {
		GRPC     GRPCConfig
		Mongo    MongoConfig
		Pipeline PipelineConfig
	}
)

func MustConfigure() *Config {
	config := &Config{}
	config.GRPC = withGRPC()
	config.Mongo = withMongoDB()
	config.Pipeline = withPipeline()

	return config
}

func withGRPC() GRPCConfig {
	grpc := GRPCConfig{
		Host: "127.0.0.1",
		Port: "9090",
	}

	if len(os.Getenv("GRPC_HOST")) > 0 {
		grpc.Host = os.Getenv("GRPC_HOST")
	}

	if len(os.Getenv("GRPC_PORT")) > 0 {
		grpc.Port = os.Getenv("GRPC_PORT")
	}

	return grpc
}

func withMongoDB() MongoConfig {
	return MongoConfig{
		Connection:     os.Getenv("MONGO_CONNECTION_STRING"),
		ConnectTimeout: 5 * time.Second,
		SocketTimeout:  5 * time.Second,
		MinPoolSize:    10,
		MaxPoolSize:    100,
	}
}

func withPipeline() PipelineConfig {
	return PipelineConfig{
		DownloaderWorkersCount:  1,
		DownloaderTempDirectory: "/tmp",
		ExporterWorkersCount:    1,
	}
}
