// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"strings"
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestUnitMtls_urlSwitching(t *testing.T) {
	t.Parallel()
	for key, bp := range transport_tpg.DefaultBasePaths {
		url := getMtlsEndpoint(bp)
		if !strings.Contains(url, ".mtls.") {
			t.Errorf("%s: mtls conversion unsuccessful preconv - %s postconv - %s", key, bp, url)
		}
	}
}
