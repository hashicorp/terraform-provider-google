// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/container/v1"
)

type ContainerOperationWaiter struct {
	Service             *container.Service
	Context             context.Context
	Op                  *container.Operation
	Project             string
	Location            string
	UserProjectOverride bool
}

func (w *ContainerOperationWaiter) State() string {
	if w == nil || w.Op == nil {
		return "<nil>"
	}
	return w.Op.Status
}

func (w *ContainerOperationWaiter) Error() error {
	if w == nil || w.Op == nil {
		return nil
	}

	// Error gets called during operation polling to see if there is an error.
	// Since container's operation doesn't have an "error" field, we must wait
	// until it's done and check the status message
	for _, pending := range w.PendingStates() {
		if w.Op.Status == pending {
			return nil
		}
	}

	if w.Op.StatusMessage != "" {
		return fmt.Errorf(w.Op.StatusMessage)
	}

	return nil
}

func (w *ContainerOperationWaiter) IsRetryable(error) bool {
	return false
}

func (w *ContainerOperationWaiter) SetOp(op interface{}) error {
	var ok bool
	w.Op, ok = op.(*container.Operation)
	if !ok {
		return fmt.Errorf("Unable to set operation. Bad type!")
	}
	return nil
}

func (w *ContainerOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil || w.Op == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	name := fmt.Sprintf("projects/%s/locations/%s/operations/%s",
		w.Project, w.Location, w.Op.Name)

	var op *container.Operation
	select {
	case <-w.Context.Done():
		log.Println("[WARN] request has been cancelled early")
		return op, errors.New("unable to finish polling, context has been cancelled")
	default:
		// default must be here to keep the previous case from blocking
	}
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (opErr error) {
			opGetCall := w.Service.Projects.Locations.Operations.Get(name)
			if w.UserProjectOverride {
				opGetCall.Header().Add("X-Goog-User-Project", w.Project)
			}
			op, opErr = opGetCall.Do()
			return opErr
		},
		Timeout: transport_tpg.DefaultRequestTimeout,
	})

	return op, err
}

func (w *ContainerOperationWaiter) OpName() string {
	if w == nil || w.Op == nil {
		return "<nil>"
	}
	return w.Op.Name
}

func (w *ContainerOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "RUNNING"}
}

func (w *ContainerOperationWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func ContainerOperationWait(config *transport_tpg.Config, op *container.Operation, project, location, activity, userAgent string, timeout time.Duration) error {
	w := &ContainerOperationWaiter{
		Service:             config.NewContainerClient(userAgent),
		Context:             config.Context,
		Op:                  op,
		Project:             project,
		Location:            location,
		UserProjectOverride: config.UserProjectOverride,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}

	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}
