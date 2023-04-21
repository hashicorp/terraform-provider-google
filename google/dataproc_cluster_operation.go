package google

import (
	"fmt"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"time"

	"google.golang.org/api/dataproc/v1"
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

func dataprocClusterOperationWait(config *transport_tpg.Config, op *dataproc.Operation, activity, userAgent string, timeout time.Duration) error {
	w := &DataprocClusterOperationWaiter{
		Service: config.NewDataprocClient(userAgent),
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
