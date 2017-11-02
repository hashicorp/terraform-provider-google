package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/dataproc/v1"
)

type DataprocClusterOperationWaiter struct {
	Service *dataproc.Service
	Op      *dataproc.Operation
}

func (w *DataprocClusterOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"false"},
		Target:  []string{"true"},
		Refresh: w.RefreshFunc(),
	}
}

func (w *DataprocClusterOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		op, err := w.Service.Projects.Regions.Operations.Get(w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %v while polling for operation %s's 'done' status", op.Done, w.Op.Name)

		return op, fmt.Sprint(op.Done), nil
	}
}

func dataprocClusterOperationWait(config *Config, op *dataproc.Operation, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	w := &DataprocClusterOperationWaiter{
		Service: config.clientDataproc,
		Op:      op,
	}

	state := w.Conf()
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = time.Duration(minTimeoutSeconds) * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	op = opRaw.(*dataproc.Operation)
	if op.Error != nil {
		return fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
	}

	return nil
}
