package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

type ResourceManagerOperationWaiter struct {
	Service *cloudresourcemanager.Service
	Op      *cloudresourcemanager.Operation
}

func (w *ResourceManagerOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		op, err := w.Service.Operations.Get(w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %v while polling for operation %s's 'done' status", op.Done, w.Op.Name)

		return op, fmt.Sprint(op.Done), nil
	}
}

func (w *ResourceManagerOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"false"},
		Target:  []string{"true"},
		Refresh: w.RefreshFunc(),
	}
}

func resourceManagerOperationWait(service *cloudresourcemanager.Service, op *cloudresourcemanager.Operation, activity string) error {
	return resourceManagerOperationWaitTime(service, op, activity, 4)
}

func resourceManagerOperationWaitTime(service *cloudresourcemanager.Service, op *cloudresourcemanager.Operation, activity string, timeoutMin int) error {
	if op.Done {
		if op.Error != nil {
			return fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
		}
		return nil
	}

	w := &ResourceManagerOperationWaiter{
		Service: service,
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

	op = opRaw.(*cloudresourcemanager.Operation)
	if op.Error != nil {
		return fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
	}

	return nil
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
