package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceComputeBackendServiceMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Compute Backend Service State v0; migrating to v1")
		is, err := migrateBackendServiceStateV0toV1(is)
		if err != nil {
			return is, err
		}
		return is, nil
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateBackendServiceStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	oldHashToValue := map[string]map[string]interface{}{}
	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "backend.") || k == "backend.#" {
			continue
		}

		// Key is now of the form backend.%d.%s
		kParts := strings.Split(k, ".")

		// Sanity check: two parts should be there and <N> should be a number
		badFormat := false
		if len(kParts) != 3 {
			badFormat = true
		} else if _, err := strconv.Atoi(kParts[1]); err != nil {
			badFormat = true
		}

		if badFormat {
			return is, fmt.Errorf("migration error: found backend key in unexpected format: %s", k)
		}

		if oldHashToValue[kParts[1]] == nil {
			oldHashToValue[kParts[1]] = map[string]interface{}{}
		}
		oldHashToValue[kParts[1]][kParts[2]] = v
	}

	oldHashToNewHash := map[string]int{}
	for k, v := range oldHashToValue {
		oldHashToNewHash[k] = resourceGoogleComputeBackendServiceBackendHash(v)
	}

	values := map[string]string{}
	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "backend.") {
			continue
		}

		if k == "backend.#" {
			continue
		}

		// Key is now of the form backend.%d.%s
		kParts := strings.Split(k, ".")
		newKey := fmt.Sprintf("%s.%d.%s", kParts[0], oldHashToNewHash[kParts[1]], kParts[2])
		values[newKey] = v
		delete(is.Attributes, k)
	}

	for k, v := range values {
		is.Attributes[k] = v
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
