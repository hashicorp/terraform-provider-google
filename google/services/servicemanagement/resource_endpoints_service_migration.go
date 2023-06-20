// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicemanagement

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func migrateEndpointsService(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		if is.Attributes["protoc_output"] == "" {
			log.Println("[DEBUG] Nothing to migrate to V1.")
			return is, nil
		}
		is.Attributes["protoc_output_base64"] = base64.StdEncoding.EncodeToString([]byte(is.Attributes["protoc_output"]))
		is.Attributes["protoc_output"] = ""
		return is, nil
	default:
		return nil, fmt.Errorf("Unexpected schema version: %d", v)
	}
}
