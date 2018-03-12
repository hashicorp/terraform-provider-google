package google

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strings"
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

type computeAddressId struct {
	Project string
	Region  string
	Name    string
}

func (s computeAddressId) canonicalId() string {
	return fmt.Sprintf(computeAddressIdTemplate, s.Project, s.Region, s.Name)
}

func parseComputeAddressId(id string, config *Config) (*computeAddressId, error) {
	var parts []string
	if computeAddressLinkRegex.MatchString(id) {
		parts = computeAddressLinkRegex.FindStringSubmatch(id)

		return &computeAddressId{
			Project: parts[1],
			Region:  parts[2],
			Name:    parts[3],
		}, nil
	} else {
		parts = strings.Split(id, "/")
	}

	if len(parts) == 3 {
		return &computeAddressId{
			Project: parts[0],
			Region:  parts[1],
			Name:    parts[2],
		}, nil
	} else if len(parts) == 2 {
		// Project is optional.
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{region}/{name}` id format.")
		}

		return &computeAddressId{
			Project: config.Project,
			Region:  parts[0],
			Name:    parts[1],
		}, nil
	} else if len(parts) == 1 {
		// Project and region is optional
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{name}` id format.")
		}
		if config.Region == "" {
			return nil, fmt.Errorf("The default region for the provider must be set when using the `{name}` id format.")
		}

		return &computeAddressId{
			Project: config.Project,
			Region:  config.Region,
			Name:    parts[0],
		}, nil
	}

	return nil, fmt.Errorf("Invalid compute address id. Expecting resource link, `{project}/{region}/{name}`, `{region}/{name}` or `{name}` format.")
}
