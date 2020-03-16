package google

import (
	"google.golang.org/api/servicenetworking/v1"
)

type ServiceNetworkingOperationWaiter struct {
	Service *servicenetworking.APIService
	CommonOperationWaiter
}

func (w *ServiceNetworkingOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func serviceNetworkingOperationWait(config *Config, op *servicenetworking.Operation, activity string) error {
	return serviceNetworkingOperationWaitTime(config, op, activity, 10)
}

func serviceNetworkingOperationWaitTime(config *Config, op *servicenetworking.Operation, activity string, timeoutMinutes int) error {
	w := &ServiceNetworkingOperationWaiter{
		Service: config.clientServiceNetworking,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes, config.PollInterval)
}
