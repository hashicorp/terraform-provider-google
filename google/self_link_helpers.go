package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Compare only the resource name of two self links/paths.
//
// Deprecated: For backward compatibility compareResourceNames is still working,
// but all new code should use CompareResourceNames in the tpgresource package instead.
func compareResourceNames(_, old, new string, _ *schema.ResourceData) bool {
	return tpgresource.CompareResourceNames("", old, new, nil)
}

// Compare only the relative path of two self links.
//
// Deprecated: For backward compatibility compareSelfLinkRelativePaths is still working,
// but all new code should use CompareSelfLinkRelativePaths in the tpgresource package instead.
func compareSelfLinkRelativePaths(_, old, new string, _ *schema.ResourceData) bool {
	return tpgresource.CompareSelfLinkRelativePaths("", old, new, nil)
}

// compareSelfLinkOrResourceName checks if two resources are the same resource
//
// Use this method when the field accepts either a name or a self_link referencing a resource.
// The value we store (i.e. `old` in this method), must be a self_link.
//
// Deprecated: For backward compatibility compareSelfLinkOrResourceName is still working,
// but all new code should use CompareSelfLinkOrResourceName in the tpgresource package instead.
func compareSelfLinkOrResourceName(_, old, new string, _ *schema.ResourceData) bool {
	return tpgresource.CompareSelfLinkOrResourceName("", old, new, nil)
}

// Hash the relative path of a self link.
//
// Deprecated: For backward compatibility selfLinkRelativePathHash is still working,
// but all new code should use SelfLinkRelativePathHash in the tpgresource package instead.
func selfLinkRelativePathHash(selfLink interface{}) int {
	return tpgresource.SelfLinkRelativePathHash(selfLink)
}

// Deprecated: For backward compatibility getRelativePath is still working,
// but all new code should use GetRelativePath in the tpgresource package instead.
func getRelativePath(selfLink string) (string, error) {
	return tpgresource.GetRelativePath(selfLink)
}

// Hash the name path of a self link.
//
// Deprecated: For backward compatibility selfLinkNameHash is still working,
// but all new code should use SelfLinkNameHash in the tpgresource package instead.
func selfLinkNameHash(selfLink interface{}) int {
	return tpgresource.SelfLinkNameHash(selfLink)
}

// Deprecated: For backward compatibility ConvertSelfLinkToV1 is still working,
// but all new code should use ConvertSelfLinkToV1 in the tpgresource package instead.
func ConvertSelfLinkToV1(link string) string {
	return tpgresource.ConvertSelfLinkToV1(link)
}

// Deprecated: For backward compatibility GetResourceNameFromSelfLink is still working,
// but all new code should use Hashcode in the GetResourceNameFromSelfLink package instead.
func GetResourceNameFromSelfLink(link string) string {
	return tpgresource.GetResourceNameFromSelfLink(link)
}

// Deprecated: For backward compatibility NameFromSelfLinkStateFunc is still working,
// but all new code should use NameFromSelfLinkStateFunc in the tpgresource package instead.
func NameFromSelfLinkStateFunc(v interface{}) string {
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

// Deprecated: For backward compatibility StoreResourceName is still working,
// but all new code should use StoreResourceName in the tpgresource package instead.
func StoreResourceName(resourceLink interface{}) string {
	return tpgresource.StoreResourceName(resourceLink)
}

// Deprecated: For backward compatibility GetZonalResourcePropertiesFromSelfLinkOrSchema is still working,
// but all new code should use GetZonalResourcePropertiesFromSelfLinkOrSchema in the tpgresource package instead.
func GetZonalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config) (string, string, string, error) {
	return tpgresource.GetZonalResourcePropertiesFromSelfLinkOrSchema(d, config)
}

// Deprecated: For backward compatibility GetRegionalResourcePropertiesFromSelfLinkOrSchema is still working,
// but all new code should use GetRegionalResourcePropertiesFromSelfLinkOrSchema in the tpgresource package instead.
func GetRegionalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config) (string, string, string, error) {
	return tpgresource.GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
}

// given a full locational (non-global) self link, returns the project + region/zone + name or an error
//
// Deprecated: For backward compatibility GetLocationalResourcePropertiesFromSelfLinkString is still working,
// but all new code should use GetLocationalResourcePropertiesFromSelfLinkString in the tpgresource package instead.
func GetLocationalResourcePropertiesFromSelfLinkString(selfLink string) (string, string, string, error) {
	return tpgresource.GetLocationalResourcePropertiesFromSelfLinkString(selfLink)
}

// This function supports selflinks that have regions and locations in their paths
//
// Deprecated: For backward compatibility GetRegionFromRegionalSelfLink is still working,
// but all new code should use GetRegionFromRegionalSelfLink in the tpgresource package instead.
func GetRegionFromRegionalSelfLink(selfLink string) string {
	return tpgresource.GetRegionFromRegionalSelfLink(selfLink)
}
