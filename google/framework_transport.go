// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func sendFrameworkRequest(p *fwtransport.FrameworkProviderConfig, method, project, rawurl, userAgent string, body map[string]interface{}, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, diag.Diagnostics) {
	return fwtransport.SendFrameworkRequest(p, method, project, rawurl, userAgent, body, errorRetryPredicates...)
}

func sendFrameworkRequestWithTimeout(p *fwtransport.FrameworkProviderConfig, method, project, rawurl, userAgent string, body map[string]interface{}, timeout time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, diag.Diagnostics) {
	return fwtransport.SendFrameworkRequestWithTimeout(p, method, project, rawurl, userAgent, body, timeout, errorRetryPredicates...)
}
