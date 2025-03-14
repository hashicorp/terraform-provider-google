// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeImagesRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"archive_size_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"labels": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"source_disk": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeImagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for image: %s", err)
	}

	filter := d.Get("filter").(string)

	images := make([]map[string]interface{}, 0)

	imageList, err := config.NewComputeClient(userAgent).Images.List(project).Filter(filter).Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Images : %s", project), fmt.Sprintf("Images : %s", project))
	}

	for _, image := range imageList.Items {
		images = append(images, map[string]interface{}{
			"name":               image.Name,
			"family":             image.Family,
			"self_link":          image.SelfLink,
			"archive_size_bytes": image.ArchiveSizeBytes,
			"creation_timestamp": image.CreationTimestamp,
			"description":        image.Description,
			"disk_size_gb":       image.DiskSizeGb,
			"image_id":           image.Id,
			"labels":             image.Labels,
			"source_disk":        image.SourceDisk,
			"source_disk_id":     image.SourceDiskId,
			"source_image_id":    image.SourceImageId,
		})
	}

	if err := d.Set("images", images); err != nil {
		return fmt.Errorf("Error retrieving images: %s", err)
	}

	d.SetId(fmt.Sprintf(
		"projects/%s/global/images",
		project,
	))

	return nil
}
