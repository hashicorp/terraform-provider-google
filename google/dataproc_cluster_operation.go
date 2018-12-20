package google

import (
	"google.golang.org/api/dataproc/v1"
)

type DataprocClusterOperationWaiter struct {
	Service *dataproc.Service
	CommonOperationWaiter
}

func (w *DataprocClusterOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Projects.Regions.Operations.Get(w.Op.Name).Do()
}

func dataprocClusterOperationWait(config *Config, op *dataproc.Operation, activity string, timeoutMinutes int) error {
	w := &DataprocClusterOperationWaiter{
		Service: config.clientDataproc,
	}
	err := w.SetOp(op)
	if err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
