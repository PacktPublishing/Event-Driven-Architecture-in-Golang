package config

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/stackus/dotenv"

	"eda-in-golang/internal/rpc"
	"eda-in-golang/internal/web"
)

type (
	PGConfig struct {
		Conn string `required:"true"`
	}

	NatsConfig struct {
		URL    string `required:"true"`
		Stream string `default:"mallbots"`
	}

	OtelConfig struct {
		ServiceName      string `envconfig:"SERVICE_NAME" default:"mallbots"`
		ExporterEndpoint string `envconfig:"EXPORTER_OTLP_ENDPOINT" default:"http://collector:4317"`
	}
	AppConfig struct {
		Environment     string
		LogLevel        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG              PGConfig
		Nats            NatsConfig
		Rpc             rpc.RpcConfig
		Web             web.WebConfig
		Otel            OtelConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg AppConfig, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}

	err = envconfig.Process("", &cfg)

	return
}
