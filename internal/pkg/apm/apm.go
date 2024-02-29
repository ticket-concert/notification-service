package apm

import (
	"notification-service/configs"
	"os"

	"go.elastic.co/apm"
	"go.elastic.co/apm/transport"
)

func InitConnection(apmUrl, apmScretToken string) {
	os.Setenv("ELASTIC_APM_SERVER_URL", apmUrl)
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", apmScretToken)

	_, _ = transport.InitDefault()
}

func GetTracer() *apm.Tracer {
	tracer, _ := apm.NewTracer(configs.GetConfig().ServiceName, configs.GetConfig().ServiceVersion)
	return tracer
}
