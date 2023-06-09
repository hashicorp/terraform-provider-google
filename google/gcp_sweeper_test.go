// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	_ "github.com/hashicorp/terraform-provider-google/google/services/firebase"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
