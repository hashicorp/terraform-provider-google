// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tags

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleTagsTagKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleTagsTagKeysRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceTagsTagKey().Schema),
				},
			},
		},
	}
}

func dataSourceGoogleTagsTagKeysRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	parent := d.Get("parent").(string)
	token := ""

	tagKeys := make([]map[string]interface{}, 0)

	for paginate := true; paginate; {
		resp, err := config.NewResourceManagerV3Client(userAgent).TagKeys.List().Parent(parent).PageSize(300).PageToken(token).Do()
		if err != nil {
			return fmt.Errorf("error reading tag key list: %s", err)
		}

		for _, tagKey := range resp.TagKeys {

			mappedData := map[string]interface{}{
				"name":            tagKey.Name,
				"namespaced_name": tagKey.NamespacedName,
				"short_name":      tagKey.ShortName,
				"parent":          tagKey.Parent,
				"create_time":     tagKey.CreateTime,
				"update_time":     tagKey.UpdateTime,
				"description":     tagKey.Description,
				"purpose":         tagKey.Purpose,
				"purpose_data":    tagKey.PurposeData,
			}
			tagKeys = append(tagKeys, mappedData)
		}
		token = resp.NextPageToken
		paginate = token != ""
	}

	d.SetId(parent)
	if err := d.Set("keys", tagKeys); err != nil {
		return fmt.Errorf("Error setting tag key name: %s", err)
	}

	return nil
}
