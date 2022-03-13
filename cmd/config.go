package cmd

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/uber-go/tally"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() *zap.Logger {
	lc := zap.NewDevelopmentConfig()
	lc.DisableStacktrace = true
	lc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := lc.Build()
	return logger
}

func Environment(logger *zap.Logger) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			logger.Error("error while loading .env file", zap.Error(err))
		}
	} else {
		logger.Info("running service without configuration from .env")
	}
}

var (
	TaskList      = "sample-tasklist"
	ClientService = "cadence-client"
	OutboudKey    = "cadence-frontend"
)

func CadenceClient(h string, d string) client.Client {
	return client.NewClient(
		BuildCadenceClient(h),
		d,
		&client.Options{
			MetricsScope: tally.NewTestScope(TaskList, map[string]string{}),
			FeatureFlags: client.FeatureFlags{
				WorkflowExecutionAlreadyCompletedErrorEnabled: true,
			},
		},
	)
}

func BuildCadenceClient(h string) workflowserviceclient.Interface {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientService))
	if err != nil {
		panic("failed to setup transport channel")
	}
	d := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientService,
		Outbounds: yarpc.Outbounds{
			OutboudKey: {Unary: ch.NewSingleOutbound(h)},
		},
	})
	if err := d.Start(); err != nil {
		panic("failed to start dispatcher")
	}
	return workflowserviceclient.New(d.ClientConfig(OutboudKey))
}
