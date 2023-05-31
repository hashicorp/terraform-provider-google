// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicemanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/servicemanagement/v1"
)

type ServiceManagementOperationWaiter struct {
	Service *servicemanagement.APIService
	tpgresource.CommonOperationWaiter
}

func (w *ServiceManagementOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func ServiceManagementOperationWaitTime(config *transport_tpg.Config, op *servicemanagement.Operation, activity, userAgent string, timeout time.Duration) (googleapi.RawMessage, error) {
	w := &ServiceManagementOperationWaiter{
		Service: config.NewServiceManClient(userAgent),
	}

	if err := w.SetOp(op); err != nil {
		return nil, err
	}

	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return nil, err
	}
	return w.Op.Response, nil
}
