package google

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

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
			zoneURLParts := strings.Split(w.Op.Zone, "/")
			zone := zoneURLParts[len(zoneURLParts)-1]
			op, err = w.Service.ZoneOperations.Get(w.Project, zone, w.Op.Name).Do()
		} else if w.Op.Region != "" {
			regionURLParts := strings.Split(w.Op.Region, "/")
			region := regionURLParts[len(regionURLParts)-1]
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

func computeOperationWait(config *Config, op *compute.Operation, project, activity string) error {
	return computeOperationWaitTime(config, op, project, activity, 4)
}

func computeOperationWaitTime(config *Config, op *compute.Operation, project, activity string, timeoutMin int) error {
	w := &ComputeOperationWaiter{
		Service: config.clientCompute,
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
