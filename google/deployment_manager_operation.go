package google

import (
	"bytes"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
)

type DeploymentManagerOperationWaiter struct {
	Config       *Config
	Project      string
	OperationUrl string
	ComputeOperationWaiter
}

func (w *DeploymentManagerOperationWaiter) IsRetryable(error) bool {
	return false
}

func (w *DeploymentManagerOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil || w.Op == nil || w.Op.SelfLink == "" {
		return nil, fmt.Errorf("cannot query unset/nil operation")
	}
	resp, err := sendRequest(w.Config, "GET", w.Project, w.Op.SelfLink, nil)
	if err != nil {
		return nil, err
	}
	op := &compute.Operation{}
	if err := Convert(resp, op); err != nil {
		return nil, fmt.Errorf("could not convert response to operation: %v", err)
	}
	return op, nil
}

func deploymentManagerOperationWaitTime(config *Config, resp interface{}, project, activity string, timeout time.Duration) error {
	op := &compute.Operation{}
	err := Convert(resp, op)
	if err != nil {
		return err
	}

	w := &DeploymentManagerOperationWaiter{
		Config:       config,
		OperationUrl: op.SelfLink,
		ComputeOperationWaiter: ComputeOperationWaiter{
			Project: project,
		},
	}
	if err := w.SetOp(op); err != nil {
		return err
	}

	return OperationWait(w, activity, timeout, config.PollInterval)
}

func (w *DeploymentManagerOperationWaiter) Error() error {
	if w != nil && w.Op != nil && w.Op.Error != nil {
		return DeploymentManagerOperationError{
			HTTPStatusCode: w.Op.HttpErrorStatusCode,
			HTTPMessage:    w.Op.HttpErrorMessage,
			OperationError: *w.Op.Error,
		}
	}
	return nil
}

// DeploymentManagerOperationError wraps information from the compute.Operation
// in an implementation of Error.
type DeploymentManagerOperationError struct {
	HTTPStatusCode int64
	HTTPMessage    string
	compute.OperationError
}

func (e DeploymentManagerOperationError) Error() string {
	var buf bytes.Buffer
	buf.WriteString("Deployment Manager returned errors for this operation, likely due to invalid configuration.")
	buf.WriteString(fmt.Sprintf("Operation failed with HTTP error %d: %s.", e.HTTPStatusCode, e.HTTPMessage))
	buf.WriteString("Errors returned: \n")
	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}
	return buf.String()
}
