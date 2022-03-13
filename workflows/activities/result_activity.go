package activities

import (
	"context"

	"github.com/adinandradrs/cadence-sample/model"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func ResultActivity(ctx context.Context, req model.ApprovalRequest) (res bool, err error) {
	activity.GetLogger(ctx).Info("result activity finished with incoming payload", zap.Any("", req))
	//workflow end, may be attached with notification etc
	return req.IsApproved, nil
}
