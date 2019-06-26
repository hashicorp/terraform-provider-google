package google

import (
	"fmt"

	"google.golang.org/api/dataproc/v1beta2"
)

type DataprocClusterOperationWaiter struct {
	Service *dataproc.Service
	CommonOperationWaiter
}

func (w *DataprocClusterOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Projects.Regions.Operations.Get(w.Op.Name).Do()
}

func dataprocClusterOperationWait(config *Config, op *dataproc.Operation, activity string, timeoutMinutes int) error {
	w := &DataprocClusterOperationWaiter{
		Service: config.clientDataprocBeta,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
