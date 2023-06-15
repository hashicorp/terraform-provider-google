// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/cloudrun"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility PollCheckKnativeStatusFunc is still working,
// but all new code should use PollCheckKnativeStatusFunc in the cloudrun package instead.
func PollCheckKnativeStatusFunc(knativeRestResponse map[string]interface{}) func(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	return cloudrun.PollCheckKnativeStatusFunc(knativeRestResponse)
}
