package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

var organizationLoggingSinkSchema = map[string]*schema.Schema{
	"org_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The numeric ID of the organization to be exported to the sink. Note that either [ORG_ID] or "organizations/[ORG_ID]" is accepted.`,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "organizations/", "", 1)
		},
	},
	"include_children": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: `Whether or not to include children organizations in the sink export. If true, logs associated with child projects are also exported; otherwise only logs relating to the provided organization are included.`,
	},
}

func organizationLoggingSinkID(d *schema.ResourceData, config *Config) (string, error) {
	organization := d.Get("org_id").(string)
	sinkName := d.Get("name").(string)

	if !strings.HasPrefix(organization, "organization") {
		organization = "organizations/" + organization
	}

	id := fmt.Sprintf("%s/sinks/%s", organization, sinkName)
	return id, nil
}

func organizationLoggingSinksPath(d *schema.ResourceData, config *Config) (string, error) {
	organization := d.Get("org_id").(string)

	if !strings.HasPrefix(organization, "organization") {
		organization = "organizations/" + organization
	}

	id := fmt.Sprintf("%s/sinks", organization)
	return id, nil
}

func resourceLoggingOrganizationSink() *schema.Resource {
	organizationLoggingSinkRead := resourceLoggingSinkRead(flattenOrganizationLoggingSink)
	organizationLoggingSinkCreate := resourceLoggingSinkCreate(organizationLoggingSinksPath, expandOrganizationLoggingSinkForCreate, organizationLoggingSinkRead)
	organizationLoggingSinkUpdate := resourceLoggingSinkUpdate(expandOrganizationLoggingSinkForUpdate, organizationLoggingSinkRead)

	return &schema.Resource{
		Create: resourceLoggingSinkAcquireOrCreate(organizationLoggingSinkID, organizationLoggingSinkCreate, organizationLoggingSinkUpdate),
		Read:   organizationLoggingSinkRead,
		Update: organizationLoggingSinkUpdate,
		Delete: resourceLoggingSinkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("org_id"),
		},
		Schema:        resourceLoggingSinkSchema(organizationLoggingSinkSchema),
		UseJSONNumber: true,
	}
}

func expandOrganizationLoggingSinkForCreate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, uniqueWriterIdentity bool) {
	obj = expandLoggingSink(d)

	obj["includeChildren"] = d.Get("include_children").(bool)
	uniqueWriterIdentity = true
	return
}

func flattenOrganizationLoggingSink(d *schema.ResourceData, res map[string]interface{}, config *Config) error {
	if err := flattenLoggingSinkBase(d, res); err != nil {
		return err
	}

	if err := d.Set("include_children", res["includeChildren"]); err != nil {
		return fmt.Errorf("Error setting include_children: %s", err)
	}

	return nil
}

func expandOrganizationLoggingSinkForUpdate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, updateMask string, uniqueWriterIdentity bool) {
	obj, updateFields := expandResourceLoggingSinkForUpdateBase(d)

	obj["includeChildren"] = d.Get("include_children").(bool)
	updateFields = append(updateFields, "include_children")

	updateMask = strings.Join(updateFields, ",")
	uniqueWriterIdentity = true
	return
}
