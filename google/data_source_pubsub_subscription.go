package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGooglePubsubSubscription() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourcePubsubSubscription().Schema)
	addRequiredFieldsToSchema(dsSchema, "name")
	addOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGooglePubsubSubscriptionRead,
		Schema: dsSchema,
	}
}

func dataSourceGooglePubsubSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := replaceVars(d, config, "projects/{{project}}/subscriptions/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourcePubsubSubscriptionRead(d, meta)
}
