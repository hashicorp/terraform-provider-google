package google

import (
	"bytes"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

type ComputeOperationWaiter struct {
	Service *compute.Service
	Op      *compute.Operation
	Project string
}

func (w *ComputeOperationWaiter) State() string {
	return w.Op.Status
}

func (w *ComputeOperationWaiter) Error() error {
	if w.Op.Error != nil {
		return ComputeOperationError(*w.Op.Error)
	}
	return nil
}

func (w *ComputeOperationWaiter) SetOp(op interface{}) error {
	w.Op = op.(*compute.Operation)
	return nil
}

func (w *ComputeOperationWaiter) QueryOp() (interface{}, error) {
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
	return w.Op.Name
}

func (w *ComputeOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "RUNNING"}
}

func (w *ComputeOperationWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func computeOperationWait(client *compute.Service, op *compute.Operation, project, activity string) error {
	return computeOperationWaitTime(client, op, project, activity, 4)
}

func computeOperationWaitTime(client *compute.Service, op *compute.Operation, project, activity string, timeoutMinutes int) error {
	w := &ComputeOperationWaiter{
		Service: client,
		Op:      op,
		Project: project,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}

func computeBetaOperationWaitTime(client *compute.Service, op *computeBeta.Operation, project, activity string, timeoutMin int) error {
	opV1 := &compute.Operation{}
	err := Convert(op, opV1)
	if err != nil {
		return err
	}

	return computeOperationWaitTime(client, opV1, project, activity, timeoutMin)
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
