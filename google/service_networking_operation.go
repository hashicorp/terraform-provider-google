package google

import (
	"time"

	"google.golang.org/api/servicenetworking/v1"
)

type ServiceNetworkingOperationWaiter struct {
	Service             *servicenetworking.APIService
	Project             string
	UserProjectOverride bool
	CommonOperationWaiter
}

func (w *ServiceNetworkingOperationWaiter) QueryOp() (interface{}, error) {
	opGetCall := w.Service.Operations.Get(w.Op.Name)
	if w.UserProjectOverride {
		opGetCall.Header().Add("X-Goog-User-Project", w.Project)
	}
	return opGetCall.Do()
}

func ServiceNetworkingOperationWaitTime(config *Config, op *servicenetworking.Operation, activity, userAgent, project string, timeout time.Duration) error {
	w := &ServiceNetworkingOperationWaiter{
		Service:             config.NewServiceNetworkingClient(userAgent),
		Project:             project,
		UserProjectOverride: config.UserProjectOverride,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
