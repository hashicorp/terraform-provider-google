package google

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func resourceComputeAddressMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Container Node Pool State v0; migrating to v1")
		return migrateComputeAddressV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateComputeAddressV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)
	log.Printf("[DEBUG] ID before migration: %s", is.ID)

	config := meta.(*Config)

	project, err := getProjectFromInstanceState(is, config)
	if err != nil {
		return is, err
	}

	region, err := getRegionFromInstanceState(is, config)
	if err != nil {
		return is, err
	}

	is.ID = computeAddressId{
		Project: project,
		Region:  region,
		Name:    is.Attributes["name"],
	}.canonicalId()

	log.Printf("[DEBUG] ID after migration: %s", is.ID)
	return is, nil
}
