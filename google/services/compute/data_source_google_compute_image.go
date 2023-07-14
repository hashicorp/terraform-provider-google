// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeImageRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "family", "filter"},
			},
			"family": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "family", "filter"},
			},
			"filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"name", "family", "filter"},
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_encryption_key_sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"licenses": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"source_disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_disk_encryption_key_sha256": {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceGoogleComputeImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	var image *compute.Image
	if v, ok := d.GetOk("name"); ok {
		log.Printf("[DEBUG] Fetching image %s", v.(string))
		image, err = config.NewComputeClient(userAgent).Images.Get(project, v.(string)).Do()
		log.Printf("[DEBUG] Fetched image %s", v.(string))
	} else if v, ok := d.GetOk("family"); ok {
		log.Printf("[DEBUG] Fetching latest non-deprecated image from family %s", v.(string))
		image, err = config.NewComputeClient(userAgent).Images.GetFromFamily(project, v.(string)).Do()
		log.Printf("[DEBUG] Fetched latest non-deprecated image from family %s", v.(string))
	} else if v, ok := d.GetOk("filter"); ok {
		images, err := config.NewComputeClient(userAgent).Images.List(project).Filter(v.(string)).Do()
		if err != nil {
			return fmt.Errorf("error retrieving list of images: %s", err)
		}

		if len(images.Items) == 1 {
			for _, im := range images.Items {
				image = im
			}
		} else if mr, ok := d.GetOk("most_recent"); len(images.Items) >= 1 && ok && mr.(bool) {
			most_recent := time.UnixMicro(0)
			for _, im := range images.Items {
				parsedTS, err := time.Parse(time.RFC3339, im.CreationTimestamp)
				if err != nil {
					return fmt.Errorf("error parsing creation timestamp: %w", err)
				}

				if parsedTS.After(most_recent) {
					most_recent = parsedTS
					image = im
				}
			}
		} else {
			return fmt.Errorf("your filter has returned more than one image or no image. Please refine your filter to return exactly one image")
		}
	} else {
		return fmt.Errorf("one of name, family or filters must be set")
	}

	if err != nil {
		return fmt.Errorf("error retrieving image information: %s", err)
	}

	var ieks256, sdeks256 string

	if image.SourceDiskEncryptionKey != nil {
		sdeks256 = image.SourceDiskEncryptionKey.Sha256
	}

	if image.ImageEncryptionKey != nil {
		ieks256 = image.ImageEncryptionKey.Sha256
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("name", image.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("family", image.Family); err != nil {
		return fmt.Errorf("Error setting family: %s", err)
	}
	if err := d.Set("archive_size_bytes", image.ArchiveSizeBytes); err != nil {
		return fmt.Errorf("Error setting archive_size_bytes: %s", err)
	}
	if err := d.Set("creation_timestamp", image.CreationTimestamp); err != nil {
		return fmt.Errorf("Error setting creation_timestamp: %s", err)
	}
	if err := d.Set("description", image.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("disk_size_gb", image.DiskSizeGb); err != nil {
		return fmt.Errorf("Error setting disk_size_gb: %s", err)
	}
	if err := d.Set("image_id", strconv.FormatUint(image.Id, 10)); err != nil {
		return fmt.Errorf("Error setting image_id: %s", err)
	}
	if err := d.Set("image_encryption_key_sha256", ieks256); err != nil {
		return fmt.Errorf("Error setting image_encryption_key_sha256: %s", err)
	}
	if err := d.Set("label_fingerprint", image.LabelFingerprint); err != nil {
		return fmt.Errorf("Error setting label_fingerprint: %s", err)
	}
	if err := d.Set("labels", image.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}
	if err := d.Set("licenses", image.Licenses); err != nil {
		return fmt.Errorf("Error setting licenses: %s", err)
	}
	if err := d.Set("self_link", image.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("source_disk", image.SourceDisk); err != nil {
		return fmt.Errorf("Error setting source_disk: %s", err)
	}
	if err := d.Set("source_disk_encryption_key_sha256", sdeks256); err != nil {
		return fmt.Errorf("Error setting source_disk_encryption_key_sha256: %s", err)
	}
	if err := d.Set("source_disk_id", image.SourceDiskId); err != nil {
		return fmt.Errorf("Error setting source_disk_id: %s", err)
	}
	if err := d.Set("source_image_id", image.SourceImageId); err != nil {
		return fmt.Errorf("Error setting source_image_id: %s", err)
	}
	if err := d.Set("status", image.Status); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/global/images/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}
