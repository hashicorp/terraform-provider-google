package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

//These functions are used by both the `resource_container_node_pool` and `resource_container_cluster` for handling regional clusters

// Deprecated: For backward compatibility isZone is still working,
// but all new code should use IsZone in the tpgresource package instead.
func isZone(location string) bool {
	return tpgresource.IsZone(location)
}

// getLocation attempts to get values in this order (if they exist):
// - location argument in the resource config
// - region argument in the resource config
// - zone argument in the resource config
// - zone argument set in the provider config
//
// Deprecated: For backward compatibility getLocation is still working,
// but all new code should use GetLocation in the tpgresource package instead.
func getLocation(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetLocation(d, config)
}

// getZone reads the "zone" value from the given resource data and falls back
// to provider's value if not given.  If neither is provided, returns an error.
//
// Deprecated: For backward compatibility getZone is still working,
// but all new code should use GetZone in the tpgresource package instead.
func getZone(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetZone(d, config)
}
