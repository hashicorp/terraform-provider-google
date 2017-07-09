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
		Pending: []string{"NOT_DONE"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
}

func (w *DataprocClusterOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := w.Service.Projects.Regions.Operations.Get(w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Progress of operation %q: is done -> %q", w.Op.Name, resp.Done)

		if resp.Done {
			return resp, "DONE", err
		}
		return resp, "NOT_DONE", err
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
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	return nil
}
