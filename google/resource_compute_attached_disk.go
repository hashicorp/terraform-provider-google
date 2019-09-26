package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	compute "google.golang.org/api/compute/v1"
)

func resourceComputeAttachedDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAttachedDiskCreate,
		Read:   resourceAttachedDiskRead,
		Delete: resourceAttachedDiskDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAttachedDiskImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Second),
			Delete: schema.DefaultTimeout(300 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"disk": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"instance": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"project": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
				Optional: true,
			},
			"zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
				Optional: true,
			},
			"device_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"mode": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "READ_WRITE",
				ValidateFunc: validation.StringInSlice([]string{"READ_ONLY", "READ_WRITE"}, false),
			},
		},
	}
}

func resourceAttachedDiskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zv, err := parseZonalFieldValue("instances", d.Get("instance").(string), "project", "zone", d, config, false)
	if err != nil {
		return err
	}

	disk := d.Get("disk").(string)
	diskName := GetResourceNameFromSelfLink(disk)
	diskSrc := fmt.Sprintf("projects/%s/zones/%s/disks/%s", zv.Project, zv.Zone, diskName)

	// Check if the disk is a regional disk
	if strings.Contains(disk, "regions") {
		rv, err := ParseRegionDiskFieldValue(disk, d, config)
		if err != nil {
			return err
		}
		diskSrc = rv.RelativeLink()
	}

	attachedDisk := compute.AttachedDisk{
		Source:     diskSrc,
		Mode:       d.Get("mode").(string),
		DeviceName: d.Get("device_name").(string),
	}

	op, err := config.clientCompute.Instances.AttachDisk(zv.Project, zv.Zone, zv.Name, &attachedDisk).Do()
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s", zv.Name, diskName))

	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, zv.Project,
		int(d.Timeout(schema.TimeoutCreate).Minutes()), "disk to attach")
	if waitErr != nil {
		d.SetId("")
		return waitErr
	}

	return resourceAttachedDiskRead(d, meta)
}

func resourceAttachedDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zv, err := parseZonalFieldValue("instances", d.Get("instance").(string), "project", "zone", d, config, false)
	if err != nil {
		return err
	}
	d.Set("project", zv.Project)
	d.Set("zone", zv.Zone)

	diskName := GetResourceNameFromSelfLink(d.Get("disk").(string))

	instance, err := config.clientCompute.Instances.Get(zv.Project, zv.Zone, zv.Name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("AttachedDisk %q", d.Id()))
	}

	// Iterate through the instance's attached disks as this is the only way to
	// confirm the disk is actually attached
	ad := findDiskByName(instance.Disks, diskName)
	if ad == nil {
		log.Printf("[WARN] Referenced disk wasn't found attached to this compute instance. Removing from state.")
		d.SetId("")
		return nil
	}

	d.Set("device_name", ad.DeviceName)
	d.Set("mode", ad.Mode)

	// Force the referenced resources to a self-link in state because it's more specific then name.
	instancePath, err := getRelativePath(instance.SelfLink)
	if err != nil {
		return err
	}
	d.Set("instance", instancePath)
	diskPath, err := getRelativePath(ad.Source)
	if err != nil {
		return err
	}
	d.Set("disk", diskPath)

	return nil
}

func resourceAttachedDiskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zv, err := parseZonalFieldValue("instances", d.Get("instance").(string), "project", "zone", d, config, false)
	if err != nil {
		return err
	}

	diskName := GetResourceNameFromSelfLink(d.Get("disk").(string))

	instance, err := config.clientCompute.Instances.Get(zv.Project, zv.Zone, zv.Name).Do()
	if err != nil {
		return err
	}

	// Confirm the disk is still attached before making the call to detach it. If the disk isn't listed as an attached
	// disk on the compute instance then return as though the delete call succeed since this is the desired state.
	ad := findDiskByName(instance.Disks, diskName)
	if ad == nil {
		return nil
	}

	op, err := config.clientCompute.Instances.DetachDisk(zv.Project, zv.Zone, zv.Name, ad.DeviceName).Do()
	if err != nil {
		return err
	}

	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, zv.Project,
		int(d.Timeout(schema.TimeoutDelete).Minutes()), fmt.Sprintf("Detaching disk from %s", zv.Name))
	if waitErr != nil {
		return waitErr
	}

	return nil
}

func resourceAttachedDiskImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	err := parseImportId(
		[]string{"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/[^/]+",
			"(?P<project>[^/]+)/(?P<zone>[^/]+)/[^/]+"}, d, config)
	if err != nil {
		return nil, err
	}

	// In all acceptable id formats the actual id will be the last in the path
	id := GetResourceNameFromSelfLink(d.Id())
	d.SetId(id)

	IDParts := strings.Split(d.Id(), ":")
	if len(IDParts) != 2 {
		return nil, fmt.Errorf("unable to determine attached disk id - id should be '{google_compute_instance.name}:{google_compute_disk.name}'")
	}
	d.Set("instance", IDParts[0])
	d.Set("disk", IDParts[1])

	return []*schema.ResourceData{d}, nil
}

func findDiskByName(disks []*compute.AttachedDisk, id string) *compute.AttachedDisk {
	for _, disk := range disks {
		if compareSelfLinkOrResourceName("", disk.Source, id, nil) {
			return disk
		}
	}

	return nil
}
