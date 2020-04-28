package google

import (
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/servicemanagement/v1"
)

type ServiceManagementOperationWaiter struct {
	Service *servicemanagement.APIService
	CommonOperationWaiter
}

func (w *ServiceManagementOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func serviceManagementOperationWaitTime(config *Config, op *servicemanagement.Operation, activity string, timeout time.Duration) (googleapi.RawMessage, error) {
	w := &ServiceManagementOperationWaiter{
		Service: config.clientServiceMan,
	}

	if err := w.SetOp(op); err != nil {
		return nil, err
	}

	if err := OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return nil, err
	}
	return w.Op.Response, nil
}
