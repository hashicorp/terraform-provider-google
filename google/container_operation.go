package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/container/v1"
	containerBeta "google.golang.org/api/container/v1beta1"
)

type ContainerOperationWaiter struct {
	Service *container.Service
	Op      *container.Operation
	Project string
	Zone    string
}

type ContainerBetaOperationWaiter struct {
	Service  *containerBeta.Service
	Op       *containerBeta.Operation
	Project  string
	Location string
}

func (w *ContainerOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
}

func (w *ContainerBetaOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
}

func (w *ContainerOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := w.Service.Projects.Zones.Operations.Get(
			w.Project, w.Zone, w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		if resp.StatusMessage != "" {
			return resp, resp.Status, fmt.Errorf(resp.StatusMessage)
		}

		log.Printf("[DEBUG] Progress of operation %q: %q", w.Op.Name, resp.Status)

		return resp, resp.Status, err
	}
}

func (w *ContainerBetaOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		name := fmt.Sprintf("projects/%s/locations/%s/operations/%s",
			w.Project, w.Location, w.Op.Name)
		resp, err := w.Service.Projects.Locations.Operations.Get(name).Do()

		if err != nil {
			return nil, "", err
		}

		if resp.StatusMessage != "" {
			return resp, resp.Status, fmt.Errorf(resp.StatusMessage)
		}

		log.Printf("[DEBUG] Progress of operation %q: %q", w.Op.Name, resp.Status)

		return resp, resp.Status, err
	}
}

func containerOperationWait(config *Config, op *container.Operation, project, zone, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	if op.Status == "DONE" {
		if op.StatusMessage != "" {
			return fmt.Errorf(op.StatusMessage)
		}
		return nil
	}

	w := &ContainerOperationWaiter{
		Service: config.clientContainer,
		Op:      op,
		Project: project,
		Zone:    zone,
	}

	state := w.Conf()
	return waitForState(state, activity, timeoutMinutes, minTimeoutSeconds)
}

func containerBetaOperationWait(config *Config, op *containerBeta.Operation, project, location, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	if op.Status == "DONE" {
		if op.StatusMessage != "" {
			return fmt.Errorf(op.StatusMessage)
		}
		return nil
	}

	w := &ContainerBetaOperationWaiter{
		Service:  config.clientContainerBeta,
		Op:       op,
		Project:  project,
		Location: location,
	}

	state := w.Conf()
	return waitForState(state, activity, timeoutMinutes, minTimeoutSeconds)
}

func waitForState(state *resource.StateChangeConf, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = time.Duration(minTimeoutSeconds) * time.Second
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}
	return nil
}

func containerSharedOperationWait(config *Config, op interface{}, project, location, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	if op == nil {
		panic("Attempted to wait on an Operation that was nil.")
	}

	switch op.(type) {
	case *container.Operation:
		return containerOperationWait(config, op.(*container.Operation), project, location, activity, timeoutMinutes, minTimeoutSeconds)
	case *containerBeta.Operation:
		return containerBetaOperationWait(config, op.(*containerBeta.Operation), project, location, activity, timeoutMinutes, minTimeoutSeconds)
	default:
		panic("Attempted to wait on an Operation of unknown type.")
	}
}
