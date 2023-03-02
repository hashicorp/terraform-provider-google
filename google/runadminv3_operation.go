package google

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/run/v2"
)

type RunAdminV2OperationWaiter struct {
	Config    *Config
	UserAgent string
	Project   string
	CommonOperationWaiter
}

func (w *RunAdminV2OperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	url := fmt.Sprintf("%s%s", w.Config.CloudRunV2BasePath, w.CommonOperationWaiter.Op.Name)

	return SendRequest(w.Config, "GET", w.Project, url, w.UserAgent, nil)
}

func createRunAdminV2Waiter(config *Config, op *run.GoogleLongrunningOperation, project, activity, userAgent string) (*RunAdminV2OperationWaiter, error) {
	w := &RunAdminV2OperationWaiter{
		Config:    config,
		UserAgent: userAgent,
		Project:   project,
	}
	if err := w.CommonOperationWaiter.SetOp(op); err != nil {
		return nil, err
	}
	return w, nil
}

func runAdminV2OperationWaitTimeWithResponse(config *Config, op *run.GoogleLongrunningOperation, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createRunAdminV2Waiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func runAdminV2OperationWaitTime(config *Config, op *run.GoogleLongrunningOperation, project, activity, userAgent string, timeout time.Duration) error {
	if op.Done {
		return nil
	}
	w, err := createRunAdminV2Waiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
