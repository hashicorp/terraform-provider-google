package google

import (
	"fmt"

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

func serviceManagementOperationWait(config *Config, op *servicemanagement.Operation, activity string) (googleapi.RawMessage, error) {
	return serviceManagementOperationWaitTime(config, op, activity, 10)
}

func serviceManagementOperationWaitTime(config *Config, op *servicemanagement.Operation, activity string, timeoutMinutes int) (googleapi.RawMessage, error) {
	w := &ServiceManagementOperationWaiter{
		Service: config.clientServiceMan,
	}

	if err := w.SetOp(op); err != nil {
		return nil, err
	}

	if err := OperationWait(w, activity, timeoutMinutes, config.PollInterval); err != nil {
		return nil, err
	}
	return w.Op.Response, nil
}
