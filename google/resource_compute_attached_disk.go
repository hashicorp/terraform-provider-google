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
				Description:      `name or self_link of the disk that will be attached.`,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"instance": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      `name or self_link of the compute instance that the disk will be attached to. If the self_link is provided then zone and project are extracted from the self link. If only the name is used then zone and project must be defined as properties on the resource or provider.`,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"project": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: `The project that the referenced compute instance is a part of. If instance is referenced by its self_link the project defined in the link will take precedence.`,
			},
			"zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: `The zone that the referenced compute instance is located within. If instance is referenced by its self_link the zone defined in the link will take precedence.`,
			},
			"device_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `Specifies a unique device name of your choice that is reflected into the /dev/disk/by-id/google-* tree of a Linux operating system running within the instance. This name can be used to reference the device for mounting, resizing, and so on, from within the instance. If not specified, the server chooses a default device name to apply to this disk, in the form persistent-disks-x, where x is a number assigned by Google Compute Engine.`,
			},
			"mode": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "READ_WRITE",
				Description:  `The mode in which to attach this disk, either READ_WRITE or READ_ONLY. If not specified, the default is to attach the disk in READ_WRITE mode.`,
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

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s/%s", zv.Project, zv.Zone, zv.Name, diskName))

	waitErr := computeOperationWaitTime(config, op, zv.Project,
		"disk to attach", d.Timeout(schema.TimeoutCreate))
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

	waitErr := computeOperationWaitTime(config, op, zv.Project,
		fmt.Sprintf("Detaching disk from %s", zv.Name), d.Timeout(schema.TimeoutDelete))
	if waitErr != nil {
		return waitErr
	}

	return nil
}

func resourceAttachedDiskImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	err := parseImportId(
		[]string{"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/(?P<instance>[^/]+)/(?P<disk>[^/]+)",
			"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<instance>[^/]+)/(?P<disk>[^/]+)"}, d, config)
	if err != nil {
		return nil, err
	}

	id, err := replaceVars(d, config, "projects/{{project}}/zones/{{zone}}/instances/{{instance}}/{{disk}}")
	if err != nil {
		return nil, err
	}
	d.SetId(id)

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
