package google

import (
	"bytes"
	"fmt"
	"log"
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
	Service   *compute.Service
	Op        *compute.Operation
	Project   string
	Scope     string
	ScopeType ScopeType
}

func (w *ComputeOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var op *compute.Operation
		var err error

		switch w.ScopeType {
		case Global:
			op, err = w.Service.GlobalOperations.Get(w.Project, w.Op.Name).Do()
		case Region:
			op, err = w.Service.RegionOperations.Get(w.Project, w.Scope, w.Op.Name).Do()
		case Zone:
			op, err = w.Service.ZoneOperations.Get(w.Project, w.Scope, w.Op.Name).Do()
		default:
			return nil, "bad-type", fmt.Errorf("Invalid wait type: %#v", w.ScopeType)
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

func computeOperationWaitGlobal(config *Config, op *compute.Operation, project, activity string) error {
	return computeOperationWaitGlobalTime(config, op, project, activity, 4)
}

func computeOperationWaitGlobalTime(config *Config, op *compute.Operation, project, activity string, timeoutMin int) error {
	w := &ComputeOperationWaiter{
		Service:   config.clientCompute,
		Op:        op,
		Project:   project,
		ScopeType: Global,
	}

	return waitComputeOperationWaiter(w, timeoutMin, activity)
}

func computeOperationWaitRegion(config *Config, op *compute.Operation, project string, region, activity string) error {
	return computeOperationWaitRegionTime(config, op, project, region, 4, activity)
}

func computeOperationWaitRegionTime(config *Config, op *compute.Operation, project, region string, timeoutMin int, activity string) error {
	w := &ComputeOperationWaiter{
		Service:   config.clientCompute,
		Op:        op,
		Project:   project,
		ScopeType: Region,
		Scope:     region,
	}

	return waitComputeOperationWaiter(w, timeoutMin, activity)
}

func computeOperationWaitZone(config *Config, op *compute.Operation, project, zone, activity string) error {
	return computeOperationWaitZoneTime(config, op, project, zone, 4, activity)
}

func computeOperationWaitZoneTime(config *Config, op *compute.Operation, project, zone string, timeoutMin int, activity string) error {
	w := &ComputeOperationWaiter{
		Service:   config.clientCompute,
		Op:        op,
		Project:   project,
		ScopeType: Zone,
		Scope:     zone,
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
