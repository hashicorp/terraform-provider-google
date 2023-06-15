// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// readDiskType finds the disk type with the given name.
func readDiskType(c *transport_tpg.Config, d tpgresource.TerraformResourceData, name string) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseZonalFieldValue("diskTypes", name, "project", "zone", d, c, false)
}

// readRegionDiskType finds the disk type with the given name.
func readRegionDiskType(c *transport_tpg.Config, d tpgresource.TerraformResourceData, name string) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseRegionalFieldValue("diskTypes", name, "project", "region", "zone", d, c, false)
}
