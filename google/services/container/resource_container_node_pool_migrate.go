// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func resourceContainerNodePoolMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Container Node Pool State v0; migrating to v1")
		return migrateNodePoolStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateNodePoolStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)
	log.Printf("[DEBUG] ID before migration: %s", is.ID)

	is.ID = fmt.Sprintf("%s/%s/%s", is.Attributes["zone"], is.Attributes["cluster"], is.Attributes["name"])

	log.Printf("[DEBUG] ID after migration: %s", is.ID)
	return is, nil
}
