package google

import (
	"encoding/json"
	"fmt"
	"time"
)

type VertexAIOperationWaiter struct {
	Config    *Config
	UserAgent string
	Project   string
	CommonOperationWaiter
}

func (w *VertexAIOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}

	region := GetRegionFromRegionalSelfLink(w.CommonOperationWaiter.Op.Name)

	// Returns the proper get.
	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/%s", region, w.CommonOperationWaiter.Op.Name)

	return SendRequest(w.Config, "GET", w.Project, url, w.UserAgent, nil)
}

func createVertexAIWaiter(config *Config, op map[string]interface{}, project, activity, userAgent string) (*VertexAIOperationWaiter, error) {
	w := &VertexAIOperationWaiter{
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
func VertexAIOperationWaitTimeWithResponse(config *Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createVertexAIWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func VertexAIOperationWaitTime(config *Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createVertexAIWaiter(config, op, project, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
