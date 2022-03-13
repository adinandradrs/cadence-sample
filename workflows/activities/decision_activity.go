package activities

import (
	"context"

	"github.com/adinandradrs/cadence-sample/model"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func DecisionActivity(ctx context.Context, req model.RegisterRequest) (out model.ApprovalRequest, err error) {
	req.Stage = model.STATUS_WAPPR
	info := activity.GetInfo(ctx)
	activity.GetLogger(ctx).Info("decision activity started with processed payload, state = pending",
		zap.Any("", out),
		zap.String("", info.WorkflowExecution.ID),
		zap.String("", info.ActivityID),
		zap.String("", info.WorkflowExecution.RunID))
	return out, activity.ErrResultPending
}
