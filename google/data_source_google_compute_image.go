package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	compute "google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeImageRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"family"},
			},
			"family": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
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
		},
	}
}

func dataSourceGoogleComputeImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	params := []string{project}
	var image *compute.Image
	if v, ok := d.GetOk("name"); ok {
		params = append(params, v.(string))
		log.Printf("[DEBUG] Fetching image %s", v.(string))
		image, err = config.clientCompute.Images.Get(project, v.(string)).Do()
		log.Printf("[DEBUG] Fetched image %s", v.(string))
	} else if v, ok := d.GetOk("family"); ok {
		params = append(params, "family", v.(string))
		log.Printf("[DEBUG] Fetching latest non-deprecated image from family %s", v.(string))
		image, err = config.clientCompute.Images.GetFromFamily(project, v.(string)).Do()
		log.Printf("[DEBUG] Fetched latest non-deprecated image from family %s", v.(string))
	} else {
		return fmt.Errorf("one of name or family must be set")
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

	d.Set("project", project)
	d.Set("name", image.Name)
	d.Set("family", image.Family)
	d.Set("archive_size_bytes", image.ArchiveSizeBytes)
	d.Set("creation_timestamp", image.CreationTimestamp)
	d.Set("description", image.Description)
	d.Set("disk_size_gb", image.DiskSizeGb)
	d.Set("image_id", strconv.FormatUint(image.Id, 10))
	d.Set("image_encryption_key_sha256", ieks256)
	d.Set("label_fingerprint", image.LabelFingerprint)
	d.Set("labels", image.Labels)
	d.Set("licenses", image.Licenses)
	d.Set("self_link", image.SelfLink)
	d.Set("source_disk", image.SourceDisk)
	d.Set("source_disk_encryption_key_sha256", sdeks256)
	d.Set("source_disk_id", image.SourceDiskId)
	d.Set("source_image_id", image.SourceImageId)
	d.Set("status", image.Status)

	d.SetId(strings.Join(params, "/"))

	return nil
}
