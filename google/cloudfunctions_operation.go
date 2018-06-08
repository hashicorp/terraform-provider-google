package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/cloudfunctions/v1"
)

type CloudFunctionsOperationWaiter struct {
	Service *cloudfunctions.Service
	Op      *cloudfunctions.Operation
}

func (w *CloudFunctionsOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		op, err := w.Service.Operations.Get(w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		status := "PENDING"
		if op.Done == true {
			status = "DONE"
		}

		log.Printf("[DEBUG] Got %q when asking for operation %q", status, w.Op.Name)
		return op, status, nil
	}
}

func (w *CloudFunctionsOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
}

func cloudFunctionsOperationWait(client *cloudfunctions.Service,
	op *cloudfunctions.Operation, activity string) error {
	return cloudFunctionsOperationWaitTime(client, op, activity, 4)
}

func cloudFunctionsOperationWaitTime(client *cloudfunctions.Service, op *cloudfunctions.Operation,
	activity string, timeoutMin int) error {
	if op.Done {
		if op.Error != nil {
			return fmt.Errorf(op.Error.Message)
		}
		return nil
	}

	w := &CloudFunctionsOperationWaiter{
		Service: client,
		Op:      op,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	resultOp := opRaw.(*cloudfunctions.Operation)
	if resultOp.Error != nil {
		return fmt.Errorf(resultOp.Error.Message)
	}

	return nil
}
