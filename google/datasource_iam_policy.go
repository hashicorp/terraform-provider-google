package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
)

// Deprecated: For backward compatibility DataSourceIamPolicy is still working,
// but all new code should use DataSourceIamPolicy in the tpgiamresource package instead.
func DataSourceIamPolicy(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc, options ...func(*tpgiamresource.IamSettings)) *schema.Resource {
	return tpgiamresource.DataSourceIamPolicy(parentSpecificSchema, newUpdaterFunc, options...)
}
