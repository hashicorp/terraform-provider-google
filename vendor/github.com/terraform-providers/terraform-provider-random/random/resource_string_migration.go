package random

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
)

func resourceRandomStringMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found random string state v0; migrating to v1")
		return migrateStringStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func redactAttributes(is *terraform.InstanceState) map[string]string {
	redactedAttributes := make(map[string]string)
	for k, v := range is.Attributes {
		redactedAttributes[k] = v
		if k == "id" || k == "result" {
			redactedAttributes[k] = "<sensitive>"
		}
	}
	return redactedAttributes
}

func migrateStringStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Random String Attributes before Migration: %#v", redactAttributes(is))

	keys := []string{"min_numeric", "min_upper", "min_lower", "min_special"}
	for _, k := range keys {
		if v := is.Attributes[k]; v == "" {
			is.Attributes[k] = "0"
		}
	}

	log.Printf("[DEBUG] Random String Attributes after State Migration: %#v", redactAttributes(is))

	return is, nil
}
