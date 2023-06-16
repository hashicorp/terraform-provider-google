// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func resourceContainerClusterMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Container Cluster State v0; migrating to v1")
		return migrateClusterStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateClusterStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	newZones := []string{}

	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "additional_zones.") {
			continue
		}

		if k == "additional_zones.#" {
			continue
		}

		// Key is now of the form additional_zones.%d
		kParts := strings.Split(k, ".")

		// Sanity check: two parts should be there and <N> should be a number
		badFormat := false
		if len(kParts) != 2 {
			badFormat = true
		} else if _, err := strconv.Atoi(kParts[1]); err != nil {
			badFormat = true
		}

		if badFormat {
			return is, fmt.Errorf("migration error: found additional_zones key in unexpected format: %s", k)
		}

		newZones = append(newZones, v)
		delete(is.Attributes, k)
	}

	for _, v := range newZones {
		hash := schema.HashString(v)
		newKey := fmt.Sprintf("additional_zones.%d", hash)
		is.Attributes[newKey] = v
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
