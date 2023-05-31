// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datastream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	datastream "google.golang.org/api/datastream/v1"
)

type DatastreamOperationWaiter struct {
	Config    *transport_tpg.Config
	UserAgent string
	Project   string
	Op        datastream.Operation
	tpgresource.CommonOperationWaiter
}

func (w *DatastreamOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	// Returns the proper get.
	url := fmt.Sprintf("%s%s", w.Config.DatastreamBasePath, w.Op.Name)

	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    w.Config,
		Method:    "GET",
		Project:   w.Project,
		RawURL:    url,
		UserAgent: w.UserAgent,
	})
}

func (w *DatastreamOperationWaiter) Error() error {
	if w != nil && w.Op.Error != nil {
		return &DatastreamOperationError{Op: w.Op}
	}
	return nil
}

func (w *DatastreamOperationWaiter) SetOp(op interface{}) error {
	w.CommonOperationWaiter.SetOp(op)
	if err := tpgresource.Convert(op, &w.Op); err != nil {
		return err
	}
	return nil
}

func createDatastreamWaiter(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string) (*DatastreamOperationWaiter, error) {
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
func DatastreamOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createDatastreamWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.Op.Response), response)
}

func DatastreamOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createDatastreamWaiter(config, op, project, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
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
