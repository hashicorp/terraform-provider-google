package google

import (
	"fmt"

	"google.golang.org/api/container/v1beta1"
)

type ContainerOperationWaiter struct {
	Service  *container.Service
	Op       *container.Operation
	Project  string
	Location string
}

func (w *ContainerOperationWaiter) State() string {
	return w.Op.Status
}

func (w *ContainerOperationWaiter) Error() error {
	if w.Op.StatusMessage != "" {
		return fmt.Errorf(w.Op.StatusMessage)
	}
	return nil
}

func (w *ContainerOperationWaiter) SetOp(op interface{}) error {
	w.Op = op.(*container.Operation)
	return nil
}

func (w *ContainerOperationWaiter) QueryOp() (interface{}, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/operations/%s",
		w.Project, w.Location, w.Op.Name)
	return w.Service.Projects.Locations.Operations.Get(name).Do()
}

func (w *ContainerOperationWaiter) OpName() string {
	return w.Op.Name
}

func (w *ContainerOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "RUNNING"}
}

func (w *ContainerOperationWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func containerOperationWait(config *Config, op *container.Operation, project, location, activity string, timeoutMinutes int) error {
	w := &ContainerOperationWaiter{
		Service:  config.clientContainerBeta,
		Op:       op,
		Project:  project,
		Location: location,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
