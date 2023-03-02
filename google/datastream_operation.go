package google

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	datastream "google.golang.org/api/datastream/v1"
)

type DatastreamOperationWaiter struct {
	Config    *Config
	UserAgent string
	Project   string
	Op        datastream.Operation
	CommonOperationWaiter
}

func (w *DatastreamOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	// Returns the proper get.
	url := fmt.Sprintf("%s%s", w.Config.DatastreamBasePath, w.Op.Name)

	return SendRequest(w.Config, "GET", w.Project, url, w.UserAgent, nil)
}

func (w *DatastreamOperationWaiter) Error() error {
	if w != nil && w.Op.Error != nil {
		return &DatastreamOperationError{Op: w.Op}
	}
	return nil
}

func (w *DatastreamOperationWaiter) SetOp(op interface{}) error {
	w.CommonOperationWaiter.SetOp(op)
	if err := Convert(op, &w.Op); err != nil {
		return err
	}
	return nil
}

func createDatastreamWaiter(config *Config, op map[string]interface{}, project, activity, userAgent string) (*DatastreamOperationWaiter, error) {
	w := &DatastreamOperationWaiter{
		Config:    config,
		UserAgent: userAgent,
		Project:   project,
	}
	if err := w.SetOp(op); err != nil {
		return nil, err
	}
	return w, nil
}

// nolint: deadcode,unused
func DatastreamOperationWaitTimeWithResponse(config *Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createDatastreamWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.Op.Response), response)
}

func DatastreamOperationWaitTime(config *Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
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

// DatastreamOperationError wraps datastream.Status and implements the
// error interface so it can be returned.
type DatastreamOperationError struct {
	Op datastream.Operation
}

func (e DatastreamOperationError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Op.Error.Details {
		buf.Write(err)
		buf.WriteString("\n")
	}
	if validations := e.extractFailedValidationResult(); validations != nil {
		buf.Write(validations)
		buf.WriteString("\n")
	}

	return buf.String()
}

// extractFailedValidationResult extracts the internal failed validations
// if there are any.
func (e DatastreamOperationError) extractFailedValidationResult() []byte {
	var metadata datastream.OperationMetadata
	data, err := e.Op.Metadata.MarshalJSON()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil
	}
	if metadata.ValidationResult == nil {
		return nil
	}
	var res []byte
	for _, v := range metadata.ValidationResult.Validations {
		if v.State == "FAILED" {
			data, err := v.MarshalJSON()
			if err != nil {
				return nil
			}
			res = append(res, data...)
			res = append(res, []byte("\n")...)
		}
	}
	return res
}
