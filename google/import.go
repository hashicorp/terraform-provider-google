package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Parse an import id extracting field values using the given list of regexes.
// They are applied in order. The first in the list is tried first.
//
// e.g:
// - projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/subnetworks/(?P<name>[^/]+) (applied first)
// - (?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+),
// - (?P<name>[^/]+) (applied last)
//
// Deprecated: For backward compatibility ParseImportId is still working,
// but all new code should use ParseImportId in the tpgresource package instead.
func ParseImportId(idRegexes []string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) error {
	return tpgresource.ParseImportId(idRegexes, d, config)
}

// Parse an import id extracting field values using the given list of regexes.
// They are applied in order. The first in the list is tried first.
// This does not mutate any of the parameters, returning a map of matches
// Similar to ParseImportId in import.go, but less import specific
//
// e.g:
// - projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/subnetworks/(?P<name>[^/]+) (applied first)
// - (?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+),
// - (?P<name>[^/]+) (applied last)
//
// Deprecated: For backward compatibility getImportIdQualifiers is still working,
// but all new code should use GetImportIdQualifiers in the tpgresource package instead.
func getImportIdQualifiers(idRegexes []string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, id string) (map[string]string, error) {
	return tpgresource.GetImportIdQualifiers(idRegexes, d, config, id)
}
