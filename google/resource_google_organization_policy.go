// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func canonicalOrgPolicyConstraint(constraint string) string {
	return resourcemanager.CanonicalOrgPolicyConstraint(constraint)
}
