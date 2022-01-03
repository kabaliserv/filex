package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		Storage Storage
		Server  Server
	}

	Storage struct {
		S3Bucket   string `envconfig:"FILEX_STORAGE_S3_BUCKET"`
		S3EndPoint string `envconfig:"FILEX_STORAGE_S3_ENDPOINT"`
	}

	Server struct {
		Addr  string `envconfig:"-"`
		Host  string `envconfig:"FILEX_SERVER_HOST" default:"localhost:3000"`
		Port  string `envconfig:"FILEX_SERVER_PORT" default:":3000"`
		Proto string `envconfig:"FILEX_SERVER_PROTO" default:"http"`
	}
)

func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)

	return cfg, err
}
