package google

import (
	"google.golang.org/api/spanner/v1"
)

type SpannerDatabaseOperationWaiter struct {
	Service *spanner.Service
	CommonOperationWaiter
}

func (w *SpannerDatabaseOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Projects.Instances.Databases.Operations.Get(w.Op.Name).Do()
}

func spannerDatabaseOperationWait(config *Config, op *spanner.Operation, activity string, timeoutMinutes int) error {
	w := &SpannerDatabaseOperationWaiter{
		Service: config.clientSpanner,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
