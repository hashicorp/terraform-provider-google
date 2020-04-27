package google

import (
	"time"

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
	return serviceNetworkingOperationWaitTime(config, op, activity, 10*time.Minute)
}

func serviceNetworkingOperationWaitTime(config *Config, op *servicenetworking.Operation, activity string, timeout time.Duration) error {
	w := &ServiceNetworkingOperationWaiter{
		Service: config.clientServiceNetworking,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
