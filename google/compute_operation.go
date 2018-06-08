package google

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

type ComputeOperationWaiter struct {
	Service *compute.Service
	Op      *compute.Operation
	Project string
}

func (w *ComputeOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var op *compute.Operation
		var err error

		if w.Op.Zone != "" {
			zone := GetResourceNameFromSelfLink(w.Op.Zone)
			op, err = w.Service.ZoneOperations.Get(w.Project, zone, w.Op.Name).Do()
		} else if w.Op.Region != "" {
			region := GetResourceNameFromSelfLink(w.Op.Region)
			op, err = w.Service.RegionOperations.Get(w.Project, region, w.Op.Name).Do()
		} else {
			op, err = w.Service.GlobalOperations.Get(w.Project, w.Op.Name).Do()
		}
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %q when asking for operation %q", op.Status, w.Op.Name)
		return op, op.Status, nil
	}
}

func (w *ComputeOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
}

// ComputeOperationError wraps compute.OperationError and implements the
// error interface so it can be returned.
type ComputeOperationError compute.OperationError

func (e ComputeOperationError) Error() string {
	var buf bytes.Buffer
	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}

	return buf.String()
}

func computeOperationWait(client *compute.Service, op *compute.Operation, project, activity string) error {
	return computeOperationWaitTime(client, op, project, activity, 4)
}

func computeOperationWaitTime(client *compute.Service, op *compute.Operation, project, activity string, timeoutMin int) error {
	if op.Status == "DONE" {
		if op.Error != nil {
			return ComputeOperationError(*op.Error)
		}
		return nil
	}

	w := &ComputeOperationWaiter{
		Service: client,
		Op:      op,
		Project: project,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	resultOp := opRaw.(*compute.Operation)
	if resultOp.Error != nil {
		return ComputeOperationError(*resultOp.Error)
	}

	return nil
}

func computeBetaOperationWaitTime(client *compute.Service, op *computeBeta.Operation, project, activity string, timeoutMin int) error {
	opV1 := &compute.Operation{}
	err := Convert(op, opV1)
	if err != nil {
		return err
	}

	return computeOperationWaitTime(client, opV1, project, activity, timeoutMin)
}
