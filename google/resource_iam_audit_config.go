package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

// Deprecated: For backward compatibility ResourceIamAuditConfig is still working,
// but all new code should use ResourceIamAuditConfig in the tpgiamresource package instead.
func ResourceIamAuditConfig(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc, resourceIdParser tpgiamresource.ResourceIdParserFunc, options ...func(*tpgiamresource.IamSettings)) *schema.Resource {
	return tpgiamresource.ResourceIamAuditConfig(parentSpecificSchema, newUpdaterFunc, resourceIdParser, options...)
}
