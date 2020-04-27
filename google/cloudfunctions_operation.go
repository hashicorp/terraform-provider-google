package google

import (
	"fmt"
	"time"

	"google.golang.org/api/cloudfunctions/v1"
)

type CloudFunctionsOperationWaiter struct {
	Service *cloudfunctions.Service
	CommonOperationWaiter
}

func (w *CloudFunctionsOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func cloudFunctionsOperationWait(config *Config, op *cloudfunctions.Operation, activity string, timeout time.Duration) error {
	w := &CloudFunctionsOperationWaiter{
		Service: config.clientCloudFunctions,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
