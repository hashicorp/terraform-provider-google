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

func DataSourceGoogleTagsTagValue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleTagsTagValueRead,

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

func dataSourceGoogleTagsTagValueRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	var tagValueMatch *resourceManagerV3.TagValue
	parent := d.Get("parent").(string)
	shortName := d.Get("short_name").(string)
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewResourceManagerV3Client(userAgent).TagValues.List().Parent(parent).PageSize(300).PageToken(token).Do()
		if err != nil {
			return fmt.Errorf("error reading tag value list: %s", err)
		}

		for _, tagValue := range resp.TagValues {
			if tagValue.ShortName == shortName {
				if tagValueMatch != nil {
					return errors.New("more than one matching tag value found")
				}
				tagValueMatch = tagValue
			}
		}
		token = resp.NextPageToken
		paginate = token != ""
	}

	if tagValueMatch == nil {
		return fmt.Errorf("tag value with short_name %s not found under parent %s", shortName, parent)
	}

	d.SetId(tagValueMatch.Name)
	nameParts := strings.Split(tagValueMatch.Name, "/")
	if err := d.Set("name", nameParts[1]); err != nil {
		return fmt.Errorf("Error setting tag value name: %s", err)
	}
	if err := d.Set("namespaced_name", tagValueMatch.NamespacedName); err != nil {
		return fmt.Errorf("Error setting tag value namespaced_name: %s", err)
	}
	if err := d.Set("create_time", tagValueMatch.CreateTime); err != nil {
		return fmt.Errorf("Error setting tag value create_time: %s", err)
	}
	if err := d.Set("update_time", tagValueMatch.UpdateTime); err != nil {
		return fmt.Errorf("Error setting tag value update_time: %s", err)
	}
	if err := d.Set("description", tagValueMatch.Description); err != nil {
		return fmt.Errorf("Error setting tag value description: %s", err)
	}

	return nil
}
