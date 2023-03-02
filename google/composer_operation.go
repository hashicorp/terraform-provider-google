package google

import (
	"fmt"
	"time"

	"google.golang.org/api/composer/v1"
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

func ComposerOperationWaitTime(config *Config, op *composer.Operation, project, activity, userAgent string, timeout time.Duration) error {
	w := &ComposerOperationWaiter{
		Service: config.NewComposerClient(userAgent).Projects.Locations,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
