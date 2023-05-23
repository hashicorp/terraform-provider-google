package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Deprecated: For backward compatibility ResourceIamBinding is still working,
// but all new code should use ResourceIamBinding in the tpgiamresource package instead.
func ResourceIamBinding(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc, resourceIdParser tpgiamresource.ResourceIdParserFunc, options ...func(*tpgiamresource.IamSettings)) *schema.Resource {
	return tpgiamresource.ResourceIamBinding(parentSpecificSchema, newUpdaterFunc, resourceIdParser, options...)
}

// Deprecated: For backward compatibility expandIamCondition is still working,
// but all new code should use ExpandIamCondition in the tpgiamresource package instead.
func expandIamCondition(v interface{}) *cloudresourcemanager.Expr {
	return tpgiamresource.ExpandIamCondition(v)
}

// Deprecated: For backward compatibility flattenIamCondition is still working,
// but all new code should use FlattenIamCondition in the tpgiamresource package instead.
func flattenIamCondition(condition *cloudresourcemanager.Expr) []map[string]interface{} {
	return tpgiamresource.FlattenIamCondition(condition)
}
