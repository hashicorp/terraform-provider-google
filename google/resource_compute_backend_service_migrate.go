package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
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

func resourceGoogleComputeBackendServiceBackendHash(v interface{}) int {
	if v == nil {
		return 0
	}

	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if group, err := getRelativePath(m["group"].(string)); err != nil {
		log.Printf("[WARN] Error on retrieving relative path of instance group: %s", err)
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	} else {
		buf.WriteString(fmt.Sprintf("%s-", group))
	}

	if v, ok := m["balancing_mode"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["capacity_scaler"]; ok {
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["description"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["max_rate"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", int64(v.(int))))
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["max_connections"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", int64(v.(int))))
	}
	if v, ok := m["max_connections_per_instance"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", int64(v.(int))))
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}

	return hashcode.String(buf.String())
}
