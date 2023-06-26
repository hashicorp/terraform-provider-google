// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sweeper

import (
	"encoding/hex"
	"hash/crc32"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func AddTestSweepers(name string, sweeper func(region string) error) {
	_, filename, _, _ := runtime.Caller(0)
	hash := crc32.NewIEEE()
	hash.Write([]byte(filename))
	hashedFilename := hex.EncodeToString(hash.Sum(nil))
	uniqueName := name + "_" + hashedFilename

	resource.AddTestSweepers(uniqueName, &resource.Sweeper{
		Name: name,
		F:    sweeper,
	})
}
