// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package appengine

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/appengine/v1"
)

var (
	appEngineOperationIdRegexp = regexp.MustCompile(fmt.Sprintf("apps/%s/operations/(.*)", verify.ProjectRegex))
)

type AppEngineOperationWaiter struct {
	Service *appengine.APIService
	AppId   string
	tpgresource.CommonOperationWaiter
}

func (w *AppEngineOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	matches := appEngineOperationIdRegexp.FindStringSubmatch(w.Op.Name)
	if len(matches) != 2 {
		return nil, fmt.Errorf("Expected %d results of parsing operation name, got %d from %s", 2, len(matches), w.Op.Name)
	}
	return w.Service.Apps.Operations.Get(w.AppId, matches[1]).Do()
}

func AppEngineOperationWaitTimeWithResponse(config *transport_tpg.Config, res interface{}, response *map[string]interface{}, appId, activity, userAgent string, timeout time.Duration) error {
	op := &appengine.Operation{}
	err := tpgresource.Convert(res, op)
	if err != nil {
		return err
	}

	w := &AppEngineOperationWaiter{
		Service: config.NewAppEngineClient(userAgent),
		AppId:   appId,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func AppEngineOperationWaitTime(config *transport_tpg.Config, res interface{}, appId, activity, userAgent string, timeout time.Duration) error {
	op := &appengine.Operation{}
	err := tpgresource.Convert(res, op)
	if err != nil {
		return err
	}

	w := &AppEngineOperationWaiter{
		Service: config.NewAppEngineClient(userAgent),
		AppId:   appId,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}
