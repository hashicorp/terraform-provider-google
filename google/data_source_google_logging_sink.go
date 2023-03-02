package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleLoggingSink() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceLoggingSinkSchema())
	dsSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: `Required. An identifier for the resource in format: "projects/[PROJECT_ID]/sinks/[SINK_NAME]", "organizations/[ORGANIZATION_ID]/sinks/[SINK_NAME]", "billingAccounts/[BILLING_ACCOUNT_ID]/sinks/[SINK_NAME]", "folders/[FOLDER_ID]/sinks/[SINK_NAME]"`,
	}

	return &schema.Resource{
		Read:   dataSourceGoogleLoggingSinkRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleLoggingSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sinkId := d.Get("id").(string)

	sink, err := config.NewLoggingClient(userAgent).Sinks.Get(sinkId).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Logging Sink %s", d.Id()))
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	d.SetId(sinkId)

	return nil
}
