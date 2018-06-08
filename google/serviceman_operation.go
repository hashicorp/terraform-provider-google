package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/servicemanagement/v1"
)

type ServiceManagementOperationWaiter struct {
	Service *servicemanagement.APIService
	Op      *servicemanagement.Operation
}

func (w *ServiceManagementOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var op *servicemanagement.Operation
		var err error

		op, err = w.Service.Operations.Get(w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %v while polling for operation %s's 'done' status", op.Done, w.Op.Name)

		return op, fmt.Sprint(op.Done), nil
	}
}

func (w *ServiceManagementOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"false"},
		Target:  []string{"true"},
		Refresh: w.RefreshFunc(),
	}
}

func serviceManagementOperationWait(config *Config, op *servicemanagement.Operation, activity string) (googleapi.RawMessage, error) {
	return serviceManagementOperationWaitTime(config, op, activity, 10)
}

func serviceManagementOperationWaitTime(config *Config, op *servicemanagement.Operation, activity string, timeoutMin int) (googleapi.RawMessage, error) {
	if op.Done {
		if op.Error != nil {
			return nil, fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
		}
		return op.Response, nil
	}

	w := &ServiceManagementOperationWaiter{
		Service: config.clientServiceMan,
		Op:      op,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	op = opRaw.(*servicemanagement.Operation)
	if op.Error != nil {
		return nil, fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
	}

	return op.Response, nil
}
