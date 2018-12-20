package google

import (
	"google.golang.org/api/spanner/v1"
)

type SpannerInstanceOperationWaiter struct {
	Service *spanner.Service
	CommonOperationWaiter
}

func (w *SpannerInstanceOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Projects.Instances.Operations.Get(w.Op.Name).Do()
}

func spannerInstanceOperationWait(config *Config, op *spanner.Operation, activity string, timeoutMinutes int) error {
	w := &SpannerInstanceOperationWaiter{
		Service: config.clientSpanner,
	}
	err := w.SetOp(op)
	if err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
