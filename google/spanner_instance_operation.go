package google

import (
	"fmt"

	"google.golang.org/api/spanner/v1"
)

type SpannerInstanceOperationWaiter struct {
	Service *spanner.Service
	CommonOperationWaiter
}

func (w *SpannerInstanceOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Projects.Instances.Operations.Get(w.Op.Name).Do()
}

func spannerOperationWaitTime(spanner *spanner.Service, op *spanner.Operation, _ string, activity string, timeoutMinutes int) error {
	if op.Name == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w := &SpannerInstanceOperationWaiter{
		Service: spanner,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
