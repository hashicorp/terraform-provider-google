// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package deploymentmanager

import (
	"bytes"
	"fmt"
	"time"

	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/compute/v1"
)

type DeploymentManagerOperationWaiter struct {
	Config       *transport_tpg.Config
	UserAgent    string
	Project      string
	OperationUrl string
	tpgcompute.ComputeOperationWaiter
}

func (w *DeploymentManagerOperationWaiter) IsRetryable(error) bool {
	return false
}

func (w *DeploymentManagerOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil || w.Op == nil || w.Op.SelfLink == "" {
		return nil, fmt.Errorf("cannot query unset/nil operation")
	}

	resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    w.Config,
		Method:    "GET",
		Project:   w.Project,
		RawURL:    w.Op.SelfLink,
		UserAgent: w.UserAgent,
	})
	if err != nil {
		return nil, err
	}
	op := &compute.Operation{}
	if err := tpgresource.Convert(resp, op); err != nil {
		return nil, fmt.Errorf("could not convert response to operation: %v", err)
	}
	return op, nil
}

func DeploymentManagerOperationWaitTime(config *transport_tpg.Config, resp interface{}, project, activity, userAgent string, timeout time.Duration) error {
	op := &compute.Operation{}
	err := tpgresource.Convert(resp, op)
	if err != nil {
		return err
	}

	w := &DeploymentManagerOperationWaiter{
		Config:       config,
		UserAgent:    userAgent,
		OperationUrl: op.SelfLink,
		ComputeOperationWaiter: tpgcompute.ComputeOperationWaiter{
			Project: project,
		},
	}
	if err := w.SetOp(op); err != nil {
		return err
	}

	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
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
