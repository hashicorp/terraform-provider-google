// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func resourceComputeInstanceTemplateMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Compute Instance Template State v0; migrating to v1")
		return migrateComputeInstanceTemplateStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateComputeInstanceTemplateStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	// automatic_restart is stored in two places. The top-level automatic_restart value is deprecated, so let's delete
	// it from the state map for now. For paranoia's sake, we compare it to the value stored in scheduling as well.
	ar := is.Attributes["automatic_restart"]
	delete(is.Attributes, "automatic_restart")

	schedulingCount, ok := is.Attributes["scheduling.#"]
	if ok && schedulingCount != "0" && schedulingCount != "1" {
		return nil, fmt.Errorf("Found multiple scheduling blocks when there should only be one")
	}

	if !ok || schedulingCount == "0" {
		// Either scheduling is missing or empty; go ahead and add
		is.Attributes["scheduling.#"] = "1"
		is.Attributes["scheduling.0.automatic_restart"] = ar
	}

	schedAr := is.Attributes["scheduling.0.automatic_restart"]
	if ar != schedAr {
		// Here we could try to choose one value over the other, but in reality they should never be out of sync; error
		// for now
		return nil, fmt.Errorf("Found differing values for automatic_restart in state, unsure how to proceed. automatic_restart = %#v, scheduling.0.automatic_restart = %#v", ar, schedAr)
	}

	// We also nuke "on_host_maintenance" as it's been deprecated as well. Here we don't check the current value though
	// as the authoritative value has always been maintained in the scheduling block.
	delete(is.Attributes, "on_host_maintenance")

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
