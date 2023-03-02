package google

import (
	"encoding/json"
	"fmt"
	"time"
)

type ContainerAttachedOperationWaiter struct {
	Config    *Config
	UserAgent string
	Project   string
	CommonOperationWaiter
}

func (w *ContainerAttachedOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}

	region := GetRegionFromRegionalSelfLink(w.CommonOperationWaiter.Op.Name)

	// Returns the proper get.
	url := fmt.Sprintf("https://%s-gkemulticloud.googleapis.com/v1/%s", region, w.CommonOperationWaiter.Op.Name)

	return SendRequest(w.Config, "GET", w.Project, url, w.UserAgent, nil)
}

func createContainerAttachedWaiter(config *Config, op map[string]interface{}, project, activity, userAgent string) (*ContainerAttachedOperationWaiter, error) {
	w := &ContainerAttachedOperationWaiter{
		Config:    config,
		UserAgent: userAgent,
		Project:   project,
	}
	if err := w.CommonOperationWaiter.SetOp(op); err != nil {
		return nil, err
	}
	return w, nil
}

// nolint: deadcode,unused
func ContainerAttachedOperationWaitTimeWithResponse(config *Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createContainerAttachedWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func ContainerAttachedOperationWaitTime(config *Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createContainerAttachedWaiter(config, op, project, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
