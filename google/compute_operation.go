package google

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/compute/v1"
)

type ComputeOperationWaiter struct {
	Service *compute.Service
	Op      *compute.Operation
	Context context.Context
	Project string
}

func (w *ComputeOperationWaiter) State() string {
	if w == nil || w.Op == nil {
		return "<nil>"
	}

	return w.Op.Status
}

func (w *ComputeOperationWaiter) Error() error {
	if w != nil && w.Op != nil && w.Op.Error != nil {
		return ComputeOperationError(*w.Op.Error)
	}
	return nil
}

func (w *ComputeOperationWaiter) IsRetryable(err error) bool {
	if oe, ok := err.(ComputeOperationError); ok {
		for _, e := range oe.Errors {
			if e.Code == "RESOURCE_NOT_READY" {
				return true
			}
		}
	}
	return false
}

func (w *ComputeOperationWaiter) SetOp(op interface{}) error {
	var ok bool
	w.Op, ok = op.(*compute.Operation)
	if !ok {
		return fmt.Errorf("Unable to set operation. Bad type!")
	}
	return nil
}

func (w *ComputeOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil || w.Op == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	if w.Context != nil {
		select {
		case <-w.Context.Done():
			log.Println("[WARN] request has been cancelled early")
			return w.Op, errors.New("unable to finish polling, context has been cancelled")
		default:
			// default must be here to keep the previous case from blocking
		}
	}
	if w.Op.Zone != "" {
		zone := GetResourceNameFromSelfLink(w.Op.Zone)
		return w.Service.ZoneOperations.Get(w.Project, zone, w.Op.Name).Do()
	} else if w.Op.Region != "" {
		region := GetResourceNameFromSelfLink(w.Op.Region)
		return w.Service.RegionOperations.Get(w.Project, region, w.Op.Name).Do()
	}
	return w.Service.GlobalOperations.Get(w.Project, w.Op.Name).Do()
}

func (w *ComputeOperationWaiter) OpName() string {
	if w == nil || w.Op == nil {
		return "<nil> Compute Op"
	}

	return w.Op.Name
}

func (w *ComputeOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "RUNNING"}
}

func (w *ComputeOperationWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func computeOperationWaitTime(config *Config, res interface{}, project, activity, userAgent string, timeout time.Duration) error {
	op := &compute.Operation{}
	err := Convert(res, op)
	if err != nil {
		return err
	}

	w := &ComputeOperationWaiter{
		Service: config.NewComputeClient(userAgent),
		Context: config.context,
		Op:      op,
		Project: project,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
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
