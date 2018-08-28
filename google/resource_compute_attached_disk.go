package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
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
			"attached_disk": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"attached_instance": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAttachedDiskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	attachedInstance := d.Get("attached_instance").(string)
	zone, err := getZoneForAttachedDisk(d, config, attachedInstance)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(attachedInstance)
	diskName := GetResourceNameFromSelfLink(d.Get("attached_disk").(string))

	attachedDisk := compute.AttachedDisk{
		DeviceName: diskName,
		Source:     fmt.Sprintf("/projects/%s/zones/%s/disks/%s", project, zone, diskName),
	}

	op, err := config.clientCompute.Instances.AttachDisk(project, zone, instanceName, &attachedDisk).Do()
	if err != nil {
		return err
	}

	// TODO (chrisst) change format of the internal id to include project/zone
	d.SetId(fmt.Sprintf("%s:%s", instanceName, diskName))

	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, project,
		int(d.Timeout(schema.TimeoutCreate).Minutes()), "disk to attach")
	if waitErr != nil {
		d.SetId("")
		return waitErr
	}

	return resourceAttachedDiskRead(d, meta)
}

func resourceAttachedDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	d.Set("project", project)

	attachedInstance := d.Get("attached_instance").(string)
	zone, err := getZoneForAttachedDisk(d, config, attachedInstance)
	if err != nil {
		return err
	}
	d.Set("zone", zone)

	instanceName := GetResourceNameFromSelfLink(attachedInstance)
	diskName := GetResourceNameFromSelfLink(d.Get("attached_disk").(string))

	instance, err := config.clientCompute.Instances.Get(project, zone, instanceName).Do()
	if err != nil {
		return err
	}

	// Iterate through the instance's attached disk as this is the only way to
	// confirm the disk is actually attached
	ad := findDiskByName(instance.Disks, diskName)
	if ad == nil {
		// Disk was not found attached to the referenced instance
		d.SetId("")
		return nil
	}

	// Force the referenced resources to a self-link in state because it's more specific then name.
	d.Set("attached_instance", instance.SelfLink)
	d.Set("attached_disk", ad.Source)

	return nil
}

func resourceAttachedDiskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	attachedInstance := d.Get("attached_instance").(string)
	zone, err := getZoneForAttachedDisk(d, config, attachedInstance)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(attachedInstance)
	diskName := GetResourceNameFromSelfLink(d.Get("attached_disk").(string))

	instance, err := config.clientCompute.Instances.Get(project, zone, instanceName).Do()
	if err != nil {
		return err
	}

	// Confirm the disk is still attached before making the call to detach it. If the disk isn't listed as an attached
	// disk on the compute instance then return as though the delete call succeed since this is the desired state.
	ad := findDiskByName(instance.Disks, diskName)
	if ad == nil {
		return nil
	}

	op, err := config.clientCompute.Instances.DetachDisk(project, zone, instanceName, diskName).Do()
	if err != nil {
		return err
	}

	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, project,
		int(d.Timeout(schema.TimeoutDelete).Minutes()), fmt.Sprintf("Detaching disk from %s", attachedInstance))
	if waitErr != nil {
		return waitErr
	}

	return nil
}

func resourceAttachedDiskImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	// TODO (chrisst) make sure to add good examples to the docs
	err := parseImportId(
		[]string{"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/[^/]+",
			"(?P<project>[^/]+)/(?P<zone>[^/]+)/[^/]+"}, d, config)
	if err != nil {
		return nil, err
	}

	// In all acceptable id formats the actual id will be the last in the path
	id := strings.Split(d.Id(), "/")
	if len(id) < 1 {
		return nil, fmt.Errorf("unable to parse resource id")
	}
	d.SetId(id[len(id)-1])

	IDParts := strings.Split(d.Id(), ":")
	if len(IDParts) != 2 {
		return nil, fmt.Errorf("unable to determine attached disk id - id should be 'google_compute_instance.name:google_compute_disk.name'")
	}
	d.Set("attached_instance", IDParts[0])
	d.Set("attached_disk", IDParts[1])

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

// getZoneForAttachedDisk prioritizes the disk defined by the compute instance self link before standard logic
func getZoneForAttachedDisk(d *schema.ResourceData, c *Config, instance string) (string, error) {
	zone, err := GetZoneFromSelfLink(instance)
	if err == nil {
		return zone, nil
	}

	// If zone can't be inferred from the compute instance self link, fall back to project
	zone, err = getZone(d, c)
	if err != nil {
		return "", fmt.Errorf("%s to inherit from the attached instance use `self_link` instead of `name`", err)
	}

	return zone, nil
}
