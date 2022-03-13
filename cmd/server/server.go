package main

import (
	"context"
	"os"

	"github.com/adinandradrs/cadence-sample/cmd"
	"github.com/adinandradrs/cadence-sample/workflows"
	"github.com/adinandradrs/cadence-sample/workflows/activities"
	"github.com/uber-go/tally"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

func main() {
	logger := cmd.Logger()
	cmd.Environment(logger)
	opt := worker.Options{
		Logger:                    logger,
		MetricsScope:              tally.NewTestScope(cmd.TaskList, map[string]string{}),
		BackgroundActivityContext: context.Background(),
	}
	w := worker.New(cmd.BuildCadenceClient(os.Getenv("CADENCE_HOST")), os.Getenv("CADENCE_DOMAIN"), cmd.TaskList, opt)
	w.RegisterWorkflow(workflows.RegisterFlow)
	w.RegisterActivity(activities.RegisterActivity)
	w.RegisterActivity(activities.DecisionActivity)
	w.RegisterActivity(activities.ResultActivity)
	if err := w.Start(); err != nil {
		panic("failed to start worker")
	}
	logger.Info("worfklow server - [worker is running]", zap.String("tasklist", cmd.TaskList))
	select {}
}
