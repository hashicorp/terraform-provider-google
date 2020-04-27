package google

import (
	"fmt"
	"time"

	composer "google.golang.org/api/composer/v1beta1"
)

type ComposerOperationWaiter struct {
	Service *composer.ProjectsLocationsService
	CommonOperationWaiter
}

func (w *ComposerOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func composerOperationWaitTime(config *Config, op *composer.Operation, project, activity string, timeout time.Duration) error {
	w := &ComposerOperationWaiter{
		Service: config.clientComposer.Projects.Locations,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
