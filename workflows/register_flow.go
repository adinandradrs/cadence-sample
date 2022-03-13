package workflows

import (
	"strconv"
	"time"

	"github.com/adinandradrs/cadence-sample/cmd"
	"github.com/adinandradrs/cadence-sample/model"
	"github.com/adinandradrs/cadence-sample/workflows/activities"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func RegisterFlow(ctx workflow.Context, inp model.RegisterRequest) (result string, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("flow register", zap.Any("register_payload", inp))
	var out model.RegisterRequest
	var out2 model.ApprovalRequest

	opt := workflow.ActivityOptions{
		TaskList:               cmd.TaskList,
		ScheduleToStartTimeout: 10 * time.Minute,
		ScheduleToCloseTimeout: 10 * time.Minute,
		StartToCloseTimeout:    10 * time.Minute,
		HeartbeatTimeout:       time.Second * 10,
	}
	ctx1 := workflow.WithActivityOptions(ctx, opt)
	err = workflow.ExecuteActivity(ctx1, activities.RegisterActivity, inp).Get(ctx1, &out)
	if err != nil {
		return "", err
	}

	opt = workflow.ActivityOptions{
		TaskList:               cmd.TaskList,
		ScheduleToStartTimeout: 10 * time.Minute,
		StartToCloseTimeout:    10 * time.Minute,
	}
	ctx2 := workflow.WithActivityOptions(ctx, opt)
	err = workflow.ExecuteActivity(ctx2, activities.DecisionActivity, inp).Get(ctx2, &out2)
	if err != nil {
		return "", err
	}
	logger.Info("continue after manual decision....", zap.Any("", out2))

	opt = workflow.ActivityOptions{
		TaskList:               cmd.TaskList,
		ScheduleToStartTimeout: 10 * time.Minute,
		StartToCloseTimeout:    10 * time.Minute,
	}
	ctx3 := workflow.WithActivityOptions(ctx, opt)
	final := false
	err = workflow.ExecuteActivity(ctx3, activities.ResultActivity, out2).Get(ctx3, nil)
	if err != nil {
		logger.Error("COMPLETED with error", zap.Error(err))
		return "", err
	} else {
		logger.Info("COMPLETED with result", zap.Bool("", final))
		return "COMPLETED with result = " + strconv.FormatBool(final), nil
	}

}
