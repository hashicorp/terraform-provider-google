package google

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sqladmin/v1beta4"
)

type SqlAdminOperationWaiter struct {
	Service *sqladmin.Service
	Op      *sqladmin.Operation
	Project string
}

func (w *SqlAdminOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] self_link: %s", w.Op.SelfLink)
		op, err := w.Service.Operations.Get(w.Project, w.Op.Name).Do()

		if e, ok := err.(*googleapi.Error); ok && (e.Code == 429 || e.Code == 503) {
			return w.Op, "PENDING", nil
		} else if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %q when asking for operation %q", op.Status, w.Op.Name)

		return op, op.Status, nil
	}
}

func (w *SqlAdminOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
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

func sqladminOperationWait(config *Config, op *sqladmin.Operation, project, activity string) error {
	return sqladminOperationWaitTime(config, op, project, activity, 10)
}

func sqladminOperationWaitTime(config *Config, op *sqladmin.Operation, project, activity string, timeoutMinutes int) error {
	if op.Status == "DONE" {
		if op.Error != nil {
			return SqlAdminOperationError(*op.Error)
		}
		return nil
	}

	w := &SqlAdminOperationWaiter{
		Service: config.clientSqlAdmin,
		Op:      op,
		Project: project,
	}

	state := w.Conf()
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = 2 * time.Second
	state.Delay = 5 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s (op %s): %s", activity, op.Name, err)
	}

	op = opRaw.(*sqladmin.Operation)
	if op.Error != nil {
		return SqlAdminOperationError(*op.Error)
	}

	return nil
}
