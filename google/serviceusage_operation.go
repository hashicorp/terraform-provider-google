package google

import (
	"google.golang.org/api/serviceusage/v1beta1"
)

type ServiceUsageOperationWaiter struct {
	Service *serviceusage.APIService
	CommonOperationWaiter
}

func (w *ServiceUsageOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func serviceUsageOperationWait(config *Config, op *serviceusage.Operation, activity string) error {
	return serviceUsageOperationWaitTime(config, op, activity, 10)
}

func serviceUsageOperationWaitTime(config *Config, op *serviceusage.Operation, activity string, timeoutMinutes int) error {
	w := &ServiceUsageOperationWaiter{
		Service: config.clientServiceUsage,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
