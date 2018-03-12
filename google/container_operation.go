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

func (w *ContainerOperationWaiter) Conf() *resource.StateChangeConf {
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

func containerOperationWait(config *Config, op *container.Operation, project, zone, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	w := &ContainerOperationWaiter{
		Service: config.clientContainer,
		Op:      op,
		Project: project,
		Zone:    zone,
	}

	state := w.Conf()
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = time.Duration(minTimeoutSeconds) * time.Second
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	return nil
}

func containerBetaOperationWait(config *Config, op *containerBeta.Operation, project, zone, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	opV1 := &container.Operation{}
	err := Convert(op, opV1)
	if err != nil {
		return err
	}

	return containerOperationWait(config, opV1, project, zone, activity, timeoutMinutes, minTimeoutSeconds)
}

func containerSharedOperationWait(config *Config, op interface{}, project, zone, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	if op == nil {
		panic("Attempted to wait on an Operation that was nil.")
	}

	switch op.(type) {
	case *container.Operation:
		return containerOperationWait(config, op.(*container.Operation), project, zone, activity, timeoutMinutes, minTimeoutSeconds)
	case *containerBeta.Operation:
		return containerBetaOperationWait(config, op.(*containerBeta.Operation), project, zone, activity, timeoutMinutes, minTimeoutSeconds)
	default:
		panic("Attempted to wait on an Operation of unknown type.")
	}
}
