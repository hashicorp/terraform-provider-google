package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	datastream "google.golang.org/api/datastream/v1"
	"time"
)

type DatastreamOperationWaiter struct {
	Config    *Config
	UserAgent string
	Project   string
	CommonOperationWaiter
}

func (w *DatastreamOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	// Returns the proper get.
	url := fmt.Sprintf("%s%s", w.Config.DatastreamBasePath, w.CommonOperationWaiter.Op.Name)

	return sendRequest(w.Config, "GET", w.Project, url, w.UserAgent, nil)
}

func (w *DatastreamOperationWaiter) Error() error {
	if w != nil && w.Op.Error != nil {
		return DatastreamError(*w.Op.Error)
	}
	return nil
}

func createDatastreamWaiter(config *Config, op map[string]interface{}, project, activity, userAgent string) (*DatastreamOperationWaiter, error) {
	w := &DatastreamOperationWaiter{
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
func datastreamOperationWaitTimeWithResponse(config *Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createDatastreamWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func datastreamOperationWaitTime(config *Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createDatastreamWaiter(config, op, project, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}

// DatastreamError wraps datastream.Status and implements the
// error interface so it can be returned.
type DatastreamError datastream.Status

func (e DatastreamError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Details {
		buf.Write(err)
		buf.WriteString("\n")
	}

	return buf.String()
}
