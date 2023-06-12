// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"google.golang.org/api/compute/v1"
)

func findDiskByName(disks []*compute.AttachedDisk, id string) *compute.AttachedDisk {
	return tpgcompute.FindDiskByName(disks, id)
}
