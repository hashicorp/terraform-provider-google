package google

import (
	composer "google.golang.org/api/composer/v1beta1"
)

type ComposerOperationWaiter struct {
	Service *composer.ProjectsLocationsService
	CommonOperationWaiter
}

func (w *ComposerOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func composerOperationWaitTime(service *composer.Service, op *composer.Operation, project, activity string, timeoutMinutes int) error {
	w := &ComposerOperationWaiter{
		Service: service.Projects.Locations,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
