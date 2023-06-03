// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type SqlAdminOperationWaiter struct {
	Service *sqladmin.Service
	Op      *sqladmin.Operation
	Project string
}

func (w *SqlAdminOperationWaiter) State() string {
	if w == nil {
		return "Operation Waiter is nil!"
	}

	if w.Op == nil {
		return "Operation is nil!"
	}

	return w.Op.Status
}

func (w *SqlAdminOperationWaiter) Error() error {
	if w != nil && w.Op != nil && w.Op.Error != nil {
		return SqlAdminOperationError(*w.Op.Error)
	}
	return nil
}

func (w *SqlAdminOperationWaiter) IsRetryable(error) bool {
	return false
}

func (w *SqlAdminOperationWaiter) SetOp(op interface{}) error {
	if op == nil {
		// Starting as a log statement, this may be a useful error in the future
		log.Printf("[DEBUG] attempted to set nil op")
	}

	sqlOp, ok := op.(*sqladmin.Operation)
	w.Op = sqlOp
	if !ok {
		return fmt.Errorf("Unable to set operation. Bad type!")
	}

	return nil
}

func (w *SqlAdminOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, waiter is unset or nil.")
	}

	if w.Op == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}

	if w.Service == nil {
		return nil, fmt.Errorf("Cannot query operation, service is nil.")
	}

	var op interface{}
	var err error
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			op, err = w.Service.Operations.Get(w.Project, w.Op.Name).Do()
			return err
		},
		Timeout: transport_tpg.DefaultRequestTimeout,
	})

	return op, err
}

func (w *SqlAdminOperationWaiter) OpName() string {
	if w == nil {
		return "<nil waiter>"
	}

	if w.Op == nil {
		return "<nil op>"
	}

	return w.Op.Name
}

func (w *SqlAdminOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "RUNNING"}
}

func (w *SqlAdminOperationWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func SqlAdminOperationWaitTime(config *transport_tpg.Config, res interface{}, project, activity, userAgent string, timeout time.Duration) error {
	op := &sqladmin.Operation{}
	err := tpgresource.Convert(res, op)
	if err != nil {
		return err
	}

	w := &SqlAdminOperationWaiter{
		Service: config.NewSqlAdminClient(userAgent),
		Op:      op,
		Project: project,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}

// SqlAdminOperationError wraps sqladmin.OperationError and implements the
// error interface so it can be returned.
type SqlAdminOperationError sqladmin.OperationErrors

func (e SqlAdminOperationError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}

	return buf.String()
}

// Retry if Cloud SQL operation returns a 429 with a specific message for
// concurrent operations.
func IsSqlInternalError(err error) (bool, string) {
	if gerr, ok := err.(*SqlAdminOperationError); ok {
		// SqlAdminOperationError is a non-interface type so we need to cast it through
		// a layer of interface{}.  :)
		var ierr interface{}
		ierr = gerr
		if serr, ok := ierr.(*sqladmin.OperationErrors); ok && serr.Errors[0].Code == "INTERNAL_ERROR" {
			return true, "Received an internal error, which is sometimes retryable for some SQL resources.  Optimistically retrying."
		}

	}
	return false, ""
}
