// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// SharedConfigForRegion returns a common config setup needed for the sweeper
// functions for a given region
//
// Deprecated: For backward compatibility SharedConfigForRegion is still working,
// but all new code should use SharedConfigForRegion in the sweeper package instead.
func SharedConfigForRegion(region string) (*transport_tpg.Config, error) {
	return sweeper.SharedConfigForRegion(region)
}

// Deprecated: For backward compatibility IsSweepableTestResource is still working,
// but all new code should use IsSweepableTestResource in the sweeper package instead.
func IsSweepableTestResource(resourceName string) bool {
	return sweeper.IsSweepableTestResource(resourceName)
}
