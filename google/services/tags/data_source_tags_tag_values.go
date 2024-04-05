// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tags

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleTagsTagValues() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleTagsTagValuesRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceTagsTagValue().Schema),
				},
			},
		},
	}
}

func dataSourceGoogleTagsTagValuesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	parent := d.Get("parent").(string)
	token := ""

	tagValues := make([]map[string]interface{}, 0)

	for paginate := true; paginate; {
		resp, err := config.NewResourceManagerV3Client(userAgent).TagValues.List().Parent(parent).PageSize(300).PageToken(token).Do()
		if err != nil {
			return fmt.Errorf("error reading tag value list: %s", err)
		}

		for _, tagValue := range resp.TagValues {
			mappedData := map[string]interface{}{
				"name":            tagValue.Name,
				"namespaced_name": tagValue.NamespacedName,
				"short_name":      tagValue.ShortName,
				"parent":          tagValue.Parent,
				"create_time":     tagValue.CreateTime,
				"update_time":     tagValue.UpdateTime,
				"description":     tagValue.Description,
			}

			tagValues = append(tagValues, mappedData)
		}
		token = resp.NextPageToken
		paginate = token != ""
	}

	d.SetId(parent)

	if err := d.Set("values", tagValues); err != nil {
		return fmt.Errorf("Error setting tag values: %s", err)
	}

	return nil
}
