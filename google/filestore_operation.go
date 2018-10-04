package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	file "google.golang.org/api/file/v1beta1"
)

type FilestoreOperationWaiter struct {
	Service *file.ProjectsLocationsService
	Op      *file.Operation
}

func (w *FilestoreOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		op, err := w.Service.Operations.Get(w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %v while polling for operation %s's 'done' status", op.Done, w.Op.Name)

		return op, fmt.Sprint(op.Done), nil
	}
}

func (w *FilestoreOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"false"},
		Target:  []string{"true"},
		Refresh: w.RefreshFunc(),
	}
}

func filestoreOperationWait(service *file.Service, op *file.Operation, project, activity string) error {
	return filestoreOperationWaitTime(service, op, project, activity, 4)
}

func filestoreOperationWaitTime(service *file.Service, op *file.Operation, project, activity string, timeoutMin int) error {
	if op.Done {
		if op.Error != nil {
			return fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
		}
		return nil
	}

	w := &FilestoreOperationWaiter{
		Service: service.Projects.Locations,
		Op:      op,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	op = opRaw.(*file.Operation)
	if op.Error != nil {
		return fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
	}

	return nil
}
