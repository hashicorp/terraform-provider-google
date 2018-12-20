package google

import (
	"bytes"

	"google.golang.org/api/sqladmin/v1beta4"
)

type SqlAdminOperationWaiter struct {
	Service *sqladmin.Service
	Op      *sqladmin.Operation
	Project string
}

func (w *SqlAdminOperationWaiter) State() string {
	return w.Op.Status
}

func (w *SqlAdminOperationWaiter) Error() error {
	if w.Op.Error != nil {
		return SqlAdminOperationError(*w.Op.Error)
	}
	return nil
}

func (w *SqlAdminOperationWaiter) SetOp(op interface{}) error {
	w.Op = op.(*sqladmin.Operation)
	return nil
}

func (w *SqlAdminOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Operations.Get(w.Project, w.Op.Name).Do()
}

func (w *SqlAdminOperationWaiter) OpName() string {
	return w.Op.Name
}

func (w *SqlAdminOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "RUNNING"}
}

func (w *SqlAdminOperationWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func sqladminOperationWait(config *Config, op *sqladmin.Operation, project, activity string) error {
	return sqladminOperationWaitTime(config, op, project, activity, 10)
}

func sqladminOperationWaitTime(config *Config, op *sqladmin.Operation, project, activity string, timeoutMinutes int) error {
	w := &SqlAdminOperationWaiter{
		Service: config.clientSqlAdmin,
		Op:      op,
		Project: project,
	}
	err := w.SetOp(op)
	if err != nil {
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
