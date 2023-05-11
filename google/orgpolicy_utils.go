package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// OrgPolicyPolicy has a custom import method because the parent field needs to allow an additional forward slash
// to represent the type of parent (e.g. projects/{project_id}).
func resourceOrgPolicyPolicyCustomImport(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^(?P<parent>[^/]+/?[^/]*)/policies/(?P<name>[^/]+)",
		"^(?P<parent>[^/]+/?[^/]*)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsRecursive(d, config, "{{parent}}/policies/{{name}}", false, 0)
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}
