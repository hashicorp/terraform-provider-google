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
	log.Printf("[DEBUG] hashing %v", m)

	if group, err := getRelativePath(m["group"].(string)); err != nil {
		log.Printf("[WARN] Error on retrieving relative path of instance group: %s", err)
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	} else {
		buf.WriteString(fmt.Sprintf("%s-", group))
	}

	if v, ok := m["balancing_mode"]; ok {
		if v == nil {
			v = ""
		}

		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["capacity_scaler"]; ok {
		if v == nil {
			v = 0.0
		}

		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["description"]; ok {
		if v == nil {
			v = ""
		}

		log.Printf("[DEBUG] writing description %s", v)
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["max_rate"]; ok {
		if v == nil {
			v = 0
		}

		buf.WriteString(fmt.Sprintf("%d-", int64(v.(int))))
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		if v == nil {
			v = 0.0
		}

		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["max_connections"]; ok {
		if v == nil {
			v = 0
		}

		switch v := v.(type) {
		case float64:
			// The Golang JSON library can't tell int values apart from floats,
			// because MM doesn't give fields strong types. Since another value
			// in this block was a real float, it assumed this was a float too.
			// It's not.
			// Note that math.Round in Go is from float64 -> float64, so it will
			// be a noop. int(floatVal) truncates extra parts, so if the float64
			// representation of an int falls below the real value we'll have
			// the wrong value. eg if 3 was represented as 2.999999, that would
			// convert to 2. So we add 0.5, ensuring that we'll truncate to the
			// correct value. This wouldn't remain true if we were far enough
			// from 0 that we were off by > 0.5, but no float conversion *could*
			// work correctly in that case. 53-bit floating types as the only
			// numeric type was not a good idea, thanks Javascript.
			var vInt int
			if v < 0 {
				vInt = int(v - 0.5)
			} else {
				vInt = int(v + 0.5)
			}

			log.Printf("[DEBUG] writing float value %f as integer value %v", v, vInt)
			buf.WriteString(fmt.Sprintf("%d-", vInt))
		default:
			buf.WriteString(fmt.Sprintf("%d-", int64(v.(int))))
		}
	}
	if v, ok := m["max_connections_per_instance"]; ok {
		if v == nil {
			v = 0
		}

		switch v := v.(type) {
		case float64:
			// The Golang JSON library can't tell int values apart from floats,
			// because MM doesn't give fields strong types. Since another value
			// in this block was a real float, it assumed this was a float too.
			// It's not.
			// Note that math.Round in Go is from float64 -> float64, so it will
			// be a noop. int(floatVal) truncates extra parts, so if the float64
			// representation of an int falls below the real value we'll have
			// the wrong value. eg if 3 was represented as 2.999999, that would
			// convert to 2. So we add 0.5, ensuring that we'll truncate to the
			// correct value. This wouldn't remain true if we were far enough
			// from 0 that we were off by > 0.5, but no float conversion *could*
			// work correctly in that case. 53-bit floating types as the only
			// numeric type was not a good idea, thanks Javascript.
			var vInt int
			if v < 0 {
				vInt = int(v - 0.5)
			} else {
				vInt = int(v + 0.5)
			}

			log.Printf("[DEBUG] writing float value %f as integer value %v", v, vInt)
			buf.WriteString(fmt.Sprintf("%d-", vInt))
		default:
			buf.WriteString(fmt.Sprintf("%d-", int64(v.(int))))
		}
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		if v == nil {
			v = 0.0
		}

		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}

	log.Printf("[DEBUG] computed hash value of %v from %v", hashcode.String(buf.String()), buf.String())
	return hashcode.String(buf.String())
}
