package tpgiamresource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var IamPolicyBaseDataSourceSchema = map[string]*schema.Schema{
	"policy_data": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func DataSourceIamPolicy(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc NewResourceIamUpdaterFunc, options ...func(*IamSettings)) *schema.Resource {
	settings := &IamSettings{}
	for _, o := range options {
		o(settings)
	}

	return &schema.Resource{
		Read: DatasourceIamPolicyRead(newUpdaterFunc),
		// if non-empty, this will be used to send a deprecation message when the
		// datasource is used.
		DeprecationMessage: settings.DeprecationMessage,
		Schema:             tpgresource.MergeSchemas(IamPolicyBaseDataSourceSchema, parentSpecificSchema),
		UseJSONNumber:      true,
	}
}

func DatasourceIamPolicyRead(newUpdaterFunc NewResourceIamUpdaterFunc) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		policy, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Resource %q with IAM Policy", updater.DescribeResource()))
		}

		if err := d.Set("etag", policy.Etag); err != nil {
			return fmt.Errorf("Error setting etag: %s", err)
		}
		if err := d.Set("policy_data", marshalIamPolicy(policy)); err != nil {
			return fmt.Errorf("Error setting policy_data: %s", err)
		}
		d.SetId(updater.GetResourceId())

		return nil
	}
}
