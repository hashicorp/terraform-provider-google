// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrunv2

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/run/v2"
)

type RunAdminV2OperationWaiter struct {
	Config    *transport_tpg.Config
	UserAgent string
	Project   string
	tpgresource.CommonOperationWaiter
}

func (w *RunAdminV2OperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	url := fmt.Sprintf("%s%s", w.Config.CloudRunV2BasePath, w.CommonOperationWaiter.Op.Name)

	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    w.Config,
		Method:    "GET",
		Project:   w.Project,
		RawURL:    url,
		UserAgent: w.UserAgent,
	})
}

func createRunAdminV2Waiter(config *transport_tpg.Config, op *run.GoogleLongrunningOperation, project, activity, userAgent string) (*RunAdminV2OperationWaiter, error) {
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

func RunAdminV2OperationWaitTimeWithResponse(config *transport_tpg.Config, op *run.GoogleLongrunningOperation, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createRunAdminV2Waiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func RunAdminV2OperationWaitTime(config *transport_tpg.Config, op *run.GoogleLongrunningOperation, project, activity, userAgent string, timeout time.Duration) error {
	if op.Done {
		return nil
	}
	w, err := createRunAdminV2Waiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}
