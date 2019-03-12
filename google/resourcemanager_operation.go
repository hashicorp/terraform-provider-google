package google

import (
	"fmt"

	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

type ResourceManagerOperationWaiter struct {
	Service *cloudresourcemanager.Service
	CommonOperationWaiter
}

func (w *ResourceManagerOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func resourceManagerOperationWaitTime(service *cloudresourcemanager.Service, op *cloudresourcemanager.Operation, activity string, timeoutMin int) error {
	w := &ResourceManagerOperationWaiter{
		Service: service,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMin)
}

func resourceManagerOperationWait(service *cloudresourcemanager.Service, op *cloudresourcemanager.Operation, activity string) error {
	return resourceManagerOperationWaitTime(service, op, activity, 4)
}

func resourceManagerV2Beta1OperationWait(service *cloudresourcemanager.Service, op *resourceManagerV2Beta1.Operation, activity string) error {
	return resourceManagerV2Beta1OperationWaitTime(service, op, activity, 4)
}

func resourceManagerV2Beta1OperationWaitTime(service *cloudresourcemanager.Service, op *resourceManagerV2Beta1.Operation, activity string, timeoutMin int) error {
	opV1 := &cloudresourcemanager.Operation{}
	err := Convert(op, opV1)
	if err != nil {
		return err
	}

	return resourceManagerOperationWaitTime(service, opV1, activity, timeoutMin)
}
