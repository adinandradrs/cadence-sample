package activities

import (
	"context"

	"github.com/adinandradrs/cadence-sample/model"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func RegisterActivity(ctx context.Context, req model.RegisterRequest) (out model.RegisterRequest, err error) {
	//put some input validation if you wish
	out = req
	activity.GetLogger(ctx).Info("register activity started with incoming payload", zap.Any("", out))
	return out, nil
}
