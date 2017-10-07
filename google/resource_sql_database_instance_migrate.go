package google

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceSqlDatabaseInstanceMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found SQL Database Instance State v0; migrating to v1")
		is, err := migrateSqlDatabaseInstanceStateV0toV1(is)
		if err != nil {
			return is, err
		}
		return is, nil
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateSqlDatabaseInstanceStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)
	idx := 0
	networkCount := 0
	newNetworks := make(map[string]string)
	keys := make([]string, len(is.Attributes))
	for k, _ := range is.Attributes {
		keys[idx] = k
		idx++

	}
	sort.Strings(keys)
	for _, k := range keys {
		if !strings.HasPrefix(k, "settings.0.ip_configuration.0.") {
			continue
		}

		if k == "settings.0.ip_configuration.0.authorized_networks.#" {
			continue
		}

		// We have a key that looks like "settings.0.ip_configuration.0.authorized_networks.<listn>" and we know it's not
		// settings.0.ip_configuration.0.authorized_networks.# because we deleted it above, so it must be
		// settings.0.ip_configuration.0.authorized_networks.<listn>
		// All that's left is to convert the list to a hash.
		kParts := strings.Split(k, ".")

		// Sanity check: all seven parts should be there and <listn> should be a number
		badFormat := false
		if len(kParts) != 7 {
			badFormat = true
		} else if _, err := strconv.Atoi(kParts[5]); err != nil {
			badFormat = true
		}

		if badFormat {
			return is, fmt.Errorf(
				"migration error: found network key in unexpected format: %s", k)
		}

		// Get the values for all items in the set that make up the hash.
		vTime := is.Attributes[fmt.Sprintf("settings.0.ip_configuration.0.authorized_networks.%d.expiration_time", kParts[5])]
		vName := is.Attributes[fmt.Sprintf("settings.0.ip_configuration.0.authorized_networks.%d.name", kParts[5])]
		vValue := is.Attributes[fmt.Sprintf("settings.0.ip_configuration.0.authorized_networks.%d.value", kParts[5])]

		// Generate the hash based on the expected values using the actual hash function.
		networkHash := resourceSqlDatabaseInstanceAuthNetworkHash(struct {
			expiration_time string
			name            string
			value           string
		}{
			expiration_time: vTime,
			name:            vName,
			value:           vValue,
		})

		newK := fmt.Sprintf("settings.0.ip_configuration.0.authorized_networks.%s.%s", networkHash, kParts[6])
		networkCount++
		newNetworks[newK] = is.Attributes[k]
		delete(is.Attributes, k)
	}

	for k, v := range newNetworks {
		is.Attributes[k] = v
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
