// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func isVcrEnabled() bool {
	return acctest.IsVcrEnabled()
}

// VcrTest is a wrapper for resource.Test to swap out providers for VCR providers and handle VCR specific things
// Can be called when VCR is not enabled, and it will behave as normal
func VcrTest(t *testing.T, c resource.TestCase) {
	acctest.VcrTest(t, c)
}
