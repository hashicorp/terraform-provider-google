// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkeonprem

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

type gkeonpremOpError struct {
	*cloudresourcemanager.Status
}

func (e gkeonpremOpError) Error() string {
	var validationCheck map[string]interface{}

	for _, msg := range e.Details {
		detail := make(map[string]interface{})
		if err := json.Unmarshal(msg, &detail); err != nil {
			continue
		}

		if _, ok := detail["validationCheck"]; ok {
			delete(detail, "@type")
			validationCheck = detail
		}
	}

	if validationCheck != nil {
		bytes, err := json.MarshalIndent(validationCheck, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error code %v message: %s validation check: %s", e.Code, e.Message, validationCheck)
		}

		return fmt.Sprintf("Error code %v message: %s\n %s", e.Code, e.Message, bytes)
	}

	return fmt.Sprintf("Error code %v, message: %s", e.Code, e.Message)
}

type gkeonpremOperationWaiter struct {
	Config    *transport_tpg.Config
	UserAgent string
	Project   string
	Op        tpgresource.CommonOperation
}

func (w *gkeonpremOperationWaiter) State() string {
	if w == nil {
		return fmt.Sprintf("Operation is nil!")
	}

	return fmt.Sprintf("done: %v", w.Op.Done)
}

func (w *gkeonpremOperationWaiter) Error() error {
	if w != nil && w.Op.Error != nil {
		return &gkeonpremOpError{w.Op.Error}
	}
	return nil
}

func (w *gkeonpremOperationWaiter) IsRetryable(error) bool {
	return false
}

func (w *gkeonpremOperationWaiter) SetOp(op interface{}) error {
	if err := tpgresource.Convert(op, &w.Op); err != nil {
		return err
	}
	return nil
}

func (w *gkeonpremOperationWaiter) OpName() string {
	if w == nil {
		return "<nil>"
	}

	return w.Op.Name
}

func (w *gkeonpremOperationWaiter) PendingStates() []string {
	return []string{"done: false"}
}

func (w *gkeonpremOperationWaiter) TargetStates() []string {
	return []string{"done: true"}
}

func (w *gkeonpremOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	// Returns the proper get.
	url := fmt.Sprintf("%s%s", w.Config.GkeonpremBasePath, w.Op.Name)

	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    w.Config,
		Method:    "GET",
		Project:   w.Project,
		RawURL:    url,
		UserAgent: w.UserAgent,
	})
}

func creategkeonpremWaiter(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string) (*gkeonpremOperationWaiter, error) {
	w := &gkeonpremOperationWaiter{
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
func GkeonpremOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := creategkeonpremWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.Op.Response), response)
}

func GkeonpremOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := creategkeonpremWaiter(config, op, project, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}
