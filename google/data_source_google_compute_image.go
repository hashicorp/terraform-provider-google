package google

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	var image *compute.Image
	if v, ok := d.GetOk("name"); ok {
		log.Printf("[DEBUG] Fetching image %s", v.(string))
		image, err = config.clientCompute.Images.Get(project, v.(string)).Do()
		log.Printf("[DEBUG] Fetched image %s", v.(string))
	} else if v, ok := d.GetOk("family"); ok {
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

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	if err := d.Set("name", image.Name); err != nil {
		return fmt.Errorf("Error reading name: %s", err)
	}
	if err := d.Set("family", image.Family); err != nil {
		return fmt.Errorf("Error reading family: %s", err)
	}
	if err := d.Set("archive_size_bytes", image.ArchiveSizeBytes); err != nil {
		return fmt.Errorf("Error reading archive_size_bytes: %s", err)
	}
	if err := d.Set("creation_timestamp", image.CreationTimestamp); err != nil {
		return fmt.Errorf("Error reading creation_timestamp: %s", err)
	}
	if err := d.Set("description", image.Description); err != nil {
		return fmt.Errorf("Error reading description: %s", err)
	}
	if err := d.Set("disk_size_gb", image.DiskSizeGb); err != nil {
		return fmt.Errorf("Error reading disk_size_gb: %s", err)
	}
	if err := d.Set("image_id", strconv.FormatUint(image.Id, 10)); err != nil {
		return fmt.Errorf("Error reading image_id: %s", err)
	}
	if err := d.Set("image_encryption_key_sha256", ieks256); err != nil {
		return fmt.Errorf("Error reading image_encryption_key_sha256: %s", err)
	}
	if err := d.Set("label_fingerprint", image.LabelFingerprint); err != nil {
		return fmt.Errorf("Error reading label_fingerprint: %s", err)
	}
	if err := d.Set("labels", image.Labels); err != nil {
		return fmt.Errorf("Error reading labels: %s", err)
	}
	if err := d.Set("licenses", image.Licenses); err != nil {
		return fmt.Errorf("Error reading licenses: %s", err)
	}
	if err := d.Set("self_link", image.SelfLink); err != nil {
		return fmt.Errorf("Error reading self_link: %s", err)
	}
	if err := d.Set("source_disk", image.SourceDisk); err != nil {
		return fmt.Errorf("Error reading source_disk: %s", err)
	}
	if err := d.Set("source_disk_encryption_key_sha256", sdeks256); err != nil {
		return fmt.Errorf("Error reading source_disk_encryption_key_sha256: %s", err)
	}
	if err := d.Set("source_disk_id", image.SourceDiskId); err != nil {
		return fmt.Errorf("Error reading source_disk_id: %s", err)
	}
	if err := d.Set("source_image_id", image.SourceImageId); err != nil {
		return fmt.Errorf("Error reading source_image_id: %s", err)
	}
	if err := d.Set("status", image.Status); err != nil {
		return fmt.Errorf("Error reading status: %s", err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/global/images/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}
