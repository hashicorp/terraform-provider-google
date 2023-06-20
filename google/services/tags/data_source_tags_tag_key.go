// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tags

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	resourceManagerV3 "google.golang.org/api/cloudresourcemanager/v3"
)

func DataSourceGoogleTagsTagKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleTagsTagKeyRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			"short_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespaced_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleTagsTagKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	var tagKeyMatch *resourceManagerV3.TagKey
	parent := d.Get("parent").(string)
	shortName := d.Get("short_name").(string)
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewResourceManagerV3Client(userAgent).TagKeys.List().Parent(parent).PageSize(300).PageToken(token).Do()
		if err != nil {
			return fmt.Errorf("error reading tag key list: %s", err)
		}

		for _, tagKey := range resp.TagKeys {
			if tagKey.ShortName == shortName {
				if tagKeyMatch != nil {
					return errors.New("more than one matching tag key found")
				}
				tagKeyMatch = tagKey
			}
		}
		token = resp.NextPageToken
		paginate = token != ""
	}

	if tagKeyMatch == nil {
		return fmt.Errorf("tag key with short_name %s not found under parent %s", shortName, parent)
	}

	d.SetId(tagKeyMatch.Name)
	nameParts := strings.Split(tagKeyMatch.Name, "/")
	if err := d.Set("name", nameParts[1]); err != nil {
		return fmt.Errorf("Error setting tag key name: %s", err)
	}
	if err := d.Set("namespaced_name", tagKeyMatch.NamespacedName); err != nil {
		return fmt.Errorf("Error setting tag key namespaced_name: %s", err)
	}
	if err := d.Set("create_time", tagKeyMatch.CreateTime); err != nil {
		return fmt.Errorf("Error setting tag key create_time: %s", err)
	}
	if err := d.Set("update_time", tagKeyMatch.UpdateTime); err != nil {
		return fmt.Errorf("Error setting tag key update_time: %s", err)
	}
	if err := d.Set("description", tagKeyMatch.Description); err != nil {
		return fmt.Errorf("Error setting tag key description: %s", err)
	}

	return nil
}
