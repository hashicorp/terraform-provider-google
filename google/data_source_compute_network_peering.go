package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const regexGCEName = "^(?:[a-z](?:[-a-z0-9]{0,61}[a-z0-9])?)$"

func DataSourceComputeNetworkPeering() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(ResourceComputeNetworkPeering().Schema)
	addRequiredFieldsToSchema(dsSchema, "name", "network")

	dsSchema["name"].ValidateFunc = validateRegexp(regexGCEName)
	dsSchema["network"].ValidateFunc = validateRegexp(peerNetworkLinkRegex)
	return &schema.Resource{
		Read:   dataSourceComputeNetworkPeeringRead,
		Schema: dsSchema,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(4 * time.Minute),
		},
	}
}

func dataSourceComputeNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	networkFieldValue, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", networkFieldValue.Name, d.Get("name").(string)))

	return resourceComputeNetworkPeeringRead(d, meta)
}
