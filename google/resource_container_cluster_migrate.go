package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
	case 1:
		log.Println("[INFO] Found Container Cluster State v1; migrating to v2")
		return migrateClusterStateV1toV2(is)
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

func migrateClusterStateV1toV2(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	newScopes := []string{}

	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "node_config.0.oauth_scopes") {
			continue
		}

		if k == "node_config.0.oauth_scopes.#" {
			continue
		}

		// Key is now of the form node_config.0.oauth_scopes.%d
		kParts := strings.Split(k, ".")

		// Sanity check: two parts should be there and <N> should be a number
		badFormat := false
		if len(kParts) != 4 {
			badFormat = true
		} else if _, err := strconv.Atoi(kParts[3]); err != nil {
			badFormat = true
		}

		if badFormat {
			return is, fmt.Errorf("migration error: found node_config.0.oauth_scopes key in unexpected format: %s", k)
		}

		newScopes = append(newScopes, v)
		delete(is.Attributes, k)
	}

	for _, v := range newScopes {
		hash := schema.HashString(canonicalizeServiceScope(v))
		newKey := fmt.Sprintf("node_config.0.oauth_scopes.%d", hash)
		is.Attributes[newKey] = v
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
