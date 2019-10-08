package google

import (
	"bytes"
	"fmt"
	"log"

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
	err = retryTimeDuration(
		func() error {
			op, err = w.Service.Operations.Get(w.Project, w.Op.Name).Do()
			return err
		},

		DefaultRequestTimeout,
	)

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

func sqlAdminOperationWait(service *sqladmin.Service, op *sqladmin.Operation, project, activity string) error {
	return sqlAdminOperationWaitTime(service, op, project, activity, 10)
}

func sqlAdminOperationWaitTime(service *sqladmin.Service, op *sqladmin.Operation, project, activity string, timeoutMinutes int) error {
	w := &SqlAdminOperationWaiter{
		Service: service,
		Op:      op,
		Project: project,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
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
