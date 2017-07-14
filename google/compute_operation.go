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

type ScopeType uint8

const (
	Global ScopeType = iota
	Region
	Zone
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

		// Find the scope type and scope name of an operation
		var scopeType ScopeType
		var scope string
		if w.Op.Zone != "" {
			scopeType = Zone
			scopeParts := strings.Split(w.Op.Zone, "/")
			scope = scopeParts[len(scopeParts)-1]
		} else if w.Op.Region != "" {
			scopeType = Region
			scopeParts := strings.Split(w.Op.Region, "/")
			scope = scopeParts[len(scopeParts)-1]
		} else {
			scopeType = Global
		}

		switch scopeType {
		case Global:
			op, err = w.Service.GlobalOperations.Get(w.Project, w.Op.Name).Do()
		case Region:
			op, err = w.Service.RegionOperations.Get(w.Project, scope, w.Op.Name).Do()
		case Zone:
			op, err = w.Service.ZoneOperations.Get(w.Project, scope, w.Op.Name).Do()
		default:
			return nil, "bad-type", fmt.Errorf("Invalid wait type: %#v", scopeType)
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

	return waitComputeOperationWaiter(w, timeoutMin, activity)
}

func computeOperationWaitGlobal(config *Config, op *compute.Operation, project, activity string) error {
	return computeOperationWaitGlobalTime(config, op, project, activity, 4)
}

func computeOperationWaitGlobalTime(config *Config, op *compute.Operation, project, activity string, timeoutMin int) error {
	w := &ComputeOperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Project: project,
	}

	return waitComputeOperationWaiter(w, timeoutMin, activity)
}

func computeOperationWaitRegion(config *Config, op *compute.Operation, project string, region, activity string) error {
	return computeOperationWaitRegionTime(config, op, project, region, 4, activity)
}

func computeOperationWaitRegionTime(config *Config, op *compute.Operation, project, region string, timeoutMin int, activity string) error {
	w := &ComputeOperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Project: project,
	}

	return waitComputeOperationWaiter(w, timeoutMin, activity)
}

func computeOperationWaitZone(config *Config, op *compute.Operation, project, zone, activity string) error {
	return computeOperationWaitZoneTime(config, op, project, zone, 4, activity)
}

func computeOperationWaitZoneTime(config *Config, op *compute.Operation, project, zone string, timeoutMin int, activity string) error {
	w := &ComputeOperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Project: project,
	}

	return waitComputeOperationWaiter(w, timeoutMin, activity)
}

func waitComputeOperationWaiter(w *ComputeOperationWaiter, timeoutMin int, activity string) error {
	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	op := opRaw.(*compute.Operation)
	if op.Error != nil {
		return ComputeOperationError(*op.Error)
	}

	return nil
}
