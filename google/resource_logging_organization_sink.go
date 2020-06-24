package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLoggingOrganizationSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingOrganizationSinkCreate,
		Read:   resourceLoggingOrganizationSinkRead,
		Delete: resourceLoggingOrganizationSinkDelete,
		Update: resourceLoggingOrganizationSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("org_id"),
		},
	}
	schm.Schema["org_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: `The numeric ID of the organization to be exported to the sink.`,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "organizations/", "", 1)
		},
	}
	schm.Schema["include_children"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: `Whether or not to include children organizations in the sink export. If true, logs associated with child projects are also exported; otherwise only logs relating to the provided organization are included.`,
	}

	return schm
}

func resourceLoggingOrganizationSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	org := d.Get("org_id").(string)
	id, sink := expandResourceLoggingSink(d, "organizations", org)
	sink.IncludeChildren = d.Get("include_children").(bool)

	// Must use a unique writer, since all destinations are in projects.
	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err := config.clientLogging.Organizations.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingOrganizationSinkRead(d, meta)
}

func resourceLoggingOrganizationSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink, err := config.clientLogging.Organizations.Sinks.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Organization Logging Sink %s", d.Get("name").(string)))
	}

	flattenResourceLoggingSink(d, sink)
	d.Set("include_children", sink.IncludeChildren)

	return nil
}

func resourceLoggingOrganizationSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)
	// It seems the API might actually accept an update for include_children; this is not in the list of updatable
	// properties though and might break in the future. Always include the value to prevent it changing.
	sink.IncludeChildren = d.Get("include_children").(bool)
	sink.ForceSendFields = append(sink.ForceSendFields, "IncludeChildren")

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err := config.clientLogging.Organizations.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	return resourceLoggingOrganizationSinkRead(d, meta)
}

func resourceLoggingOrganizationSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientLogging.Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}
