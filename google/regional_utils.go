package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//These functions are used by both the `resource_container_node_pool` and `resource_container_cluster` for handling regional clusters

func isZone(location string) bool {
	return len(strings.Split(location, "-")) == 3
}

func getLocation(d *schema.ResourceData, config *Config) (string, error) {
	if v, ok := d.GetOk("location"); ok {
		return v.(string), nil
	} else if v, isRegionalCluster := d.GetOk("region"); isRegionalCluster {
		return v.(string), nil
	} else {
		// If region is not explicitly set, use "zone" (or fall back to the provider-level zone).
		// For now, to avoid confusion, we require region to be set in the config to create a regional
		// cluster rather than falling back to the provider-level region.
		return getZone(d, config)
	}
}

// getZone reads the "zone" value from the given resource data and falls back
// to provider's value if not given.  If neither is provided, returns an error.
func getZone(d TerraformResourceData, config *Config) (string, error) {
	res, ok := d.GetOk("zone")
	if !ok {
		if config.Zone != "" {
			return config.Zone, nil
		}
		return "", fmt.Errorf("Cannot determine zone: set in this resource, or set provider-level zone.")
	}
	return GetResourceNameFromSelfLink(res.(string)), nil
}
