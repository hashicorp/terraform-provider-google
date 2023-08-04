// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/compute/v1"
)

type ComputeOperationWaiter struct {
	Service *compute.Service
	Op      *compute.Operation
	Context context.Context
	Project string
}

func (w *ComputeOperationWaiter) State() string {
	if w == nil || w.Op == nil {
		return "<nil>"
	}

	return w.Op.Status
}

func (w *ComputeOperationWaiter) Error() error {
	if w != nil && w.Op != nil && w.Op.Error != nil {
		return ComputeOperationError(*w.Op.Error)
	}
	return nil
}

func (w *ComputeOperationWaiter) IsRetryable(err error) bool {
	if oe, ok := err.(ComputeOperationError); ok {
		for _, e := range oe.Errors {
			if e.Code == "RESOURCE_NOT_READY" {
				return true
			}
		}
	}
	return false
}

func (w *ComputeOperationWaiter) SetOp(op interface{}) error {
	var ok bool
	w.Op, ok = op.(*compute.Operation)
	if !ok {
		return fmt.Errorf("Unable to set operation. Bad type!")
	}
	return nil
}

func (w *ComputeOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil || w.Op == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	if w.Context != nil {
		select {
		case <-w.Context.Done():
			log.Println("[WARN] request has been cancelled early")
			return w.Op, errors.New("unable to finish polling, context has been cancelled")
		default:
			// default must be here to keep the previous case from blocking
		}
	}
	if w.Op.Zone != "" {
		zone := tpgresource.GetResourceNameFromSelfLink(w.Op.Zone)
		return w.Service.ZoneOperations.Get(w.Project, zone, w.Op.Name).Do()
	} else if w.Op.Region != "" {
		region := tpgresource.GetResourceNameFromSelfLink(w.Op.Region)
		return w.Service.RegionOperations.Get(w.Project, region, w.Op.Name).Do()
	}
	return w.Service.GlobalOperations.Get(w.Project, w.Op.Name).Do()
}

func (w *ComputeOperationWaiter) OpName() string {
	if w == nil || w.Op == nil {
		return "<nil> Compute Op"
	}

	return w.Op.Name
}

func (w *ComputeOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "RUNNING"}
}

func (w *ComputeOperationWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func ComputeOperationWaitTime(config *transport_tpg.Config, res interface{}, project, activity, userAgent string, timeout time.Duration) error {
	op := &compute.Operation{}
	err := tpgresource.Convert(res, op)
	if err != nil {
		return err
	}

	w := &ComputeOperationWaiter{
		Service: config.NewComputeClient(userAgent),
		Context: config.Context,
		Op:      op,
		Project: project,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}

// ComputeOperationError wraps compute.OperationError and implements the
// error interface so it can be returned.
type ComputeOperationError compute.OperationError

func (e ComputeOperationError) Error() string {
	buf := bytes.NewBuffer(nil)
	for _, err := range e.Errors {
		writeOperationError(buf, err)
	}

	return buf.String()
}

const errMsgSep = "\n\n"

func writeOperationError(w io.StringWriter, opError *compute.OperationErrorErrors) {
	w.WriteString(opError.Message + "\n")

	var lm *compute.LocalizedMessage
	var link *compute.HelpLink

	for _, ed := range opError.ErrorDetails {
		if opError.Code == "QUOTA_EXCEEDED" && ed.QuotaInfo != nil {
			w.WriteString("\tmetric name = " + ed.QuotaInfo.MetricName + "\n")
			w.WriteString("\tlimit name = " + ed.QuotaInfo.LimitName + "\n")
			if ed.QuotaInfo.Limit != 0 {
				w.WriteString("\tlimit = " + fmt.Sprint(ed.QuotaInfo.Limit) + "\n")
			}
			if ed.QuotaInfo.FutureLimit != 0 {
				w.WriteString("\tfuture limit = " + fmt.Sprint(ed.QuotaInfo.FutureLimit) + "\n")
				w.WriteString("\trollout status = in progress\n")
			}
			if ed.QuotaInfo.Dimensions != nil {
				w.WriteString("\tdimensions = " + fmt.Sprint(ed.QuotaInfo.Dimensions) + "\n")
			}
			break
		}
		if lm == nil && ed.LocalizedMessage != nil {
			lm = ed.LocalizedMessage
		}

		if link == nil && ed.Help != nil && len(ed.Help.Links) > 0 {
			link = ed.Help.Links[0]
		}

		if lm != nil && link != nil {
			break
		}
	}

	if lm != nil && lm.Message != "" {
		w.WriteString(errMsgSep)
		w.WriteString(lm.Message + "\n")
	}

	if link != nil {
		w.WriteString(errMsgSep)

		if link.Description != "" {
			w.WriteString(link.Description + "\n")
		}

		if link.Url != "" {
			w.WriteString(link.Url + "\n")
		}
	}
}
