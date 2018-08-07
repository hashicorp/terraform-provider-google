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
			Delete: schema.DefaultTimeout(240 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"attached_disk": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attached_instance": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				// TODO play with this to determine if somebody can change between PERSISTANT vs SCRATCH.
				// And see if it's worth allowing a user to pass in
				Type:     schema.TypeString,
				Computed: true,
			},
			"self_link": {
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

	instance := d.Get("attached_instance").(string)
	zone, err := getZoneForAttachedDisk(d, config, instance)
	if err != nil {
		return err
	}
	instanceName := GetResourceNameFromSelfLink(instance)
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

	// TODO (chrisst) allow for override to timeouts
	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, project, 2, "disk to attach")
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

	attachedInstance, err := getAttachedInstance(d)
	if err != nil {
		return err
	}
	d.Set("attached_instance", attachedInstance)

	zone, err := getZoneForAttachedDisk(d, config, attachedInstance)
	if err != nil {
		return err
	}

	attachedDisk, err := getAttachedDisk(d)
	if err != nil {
		return err
	}
	d.Set("attached_disk", attachedDisk)

	instanceName := GetResourceNameFromSelfLink(attachedInstance)
	diskName := GetResourceNameFromSelfLink(attachedDisk)

	instance, err := config.clientCompute.Instances.Get(project, zone, instanceName).Do()
	if err != nil {
		return err
	}

	// Iterate through the instance's attached disk as this is the only way to
	// confirm the disk is actually attached
	ad := findDiskByName(instance.Disks, diskName)

	// Disk was not found attached to the referenced instance
	if ad == nil {
		d.SetId("")
		return nil
	}

	d.Set("type", ad.Type)

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

	ad := findDiskByName(instance.Disks, diskName)
	if ad == nil {
		return nil
	}

	op, err := config.clientCompute.Instances.DetachDisk(project, zone, instanceName, diskName).Do()
	if err != nil {
		return err
	}

	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, project,
		2, fmt.Sprintf("Detaching disk from %s/%s/%s", project, zone, instanceName))
	if waitErr != nil {
		return waitErr
	}

	return nil
}

func resourceAttachedDiskImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	// TODO (chrisst) make sure to add good examples to the docs
	// TODO (chrisst) using 'id' here is a problem. I either need to create a new computed variable (ew)
	// or rethink regex a bit
	err := parseImportId(
		[]string{"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/(?P<id>[^/]+)",
			"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<id>[^/]+)"}, d, config)
	if err != nil {
		return nil, err
	}

	// In all acceptable id formats the actual id will be the last in the path
	id := strings.Split(d.Id(), "/")
	if len(id) < 1 {
		return nil, fmt.Errorf("unable to parse resource id")
	}
	d.SetId(id[len(id)-1])

	return []*schema.ResourceData{d}, nil
}

// getZoneForAttachedDisk prioritizes the disk defined by the compute instance self link before standard logic
func getZoneForAttachedDisk(d *schema.ResourceData, c *Config, instance string) (string, error) {
	zone, err := GetZoneFromSelfLink(instance)
	if err != nil {
		return "", err
	}

	// If zone can't be inferred from the compute instance self link, fall back to project
	zone, err = getZone(d, c)
	if err != nil {
		return "", fmt.Errorf("%s to inherit from the attached instance use `self_link` instead of `name`", err)
	}

	return zone, nil
}

// getAttachedInstance uses fallback logic to look for the attached instance in multiple locations
// To enable importing this resource we need to handle the situation where only the ID is available and
// the attached instance hasn't been provided in the config and must be inferred from the id.
func getAttachedInstance(d *schema.ResourceData) (string, error) {
	attachedInstance := d.Get("attached_instance").(string)
	if attachedInstance != "" {
		return attachedInstance, nil
	}

	parts := strings.Split(d.Id(), ":")
	if len(parts) == 2 {
		return parts[0], nil
	}

	return "", fmt.Errorf("unable to determine the attached compute instance")
}

// getAttachedDisk uses fallback logic to look for the attached disk in multiple locations
// To enable importing this resource we need to handle the situation where only the ID is available and
// the attached disk hasn't been provided in the config and must be inferred from the id.
func getAttachedDisk(d *schema.ResourceData) (string, error) {
	attachedDisk := d.Get("attached_disk").(string)
	if attachedDisk != "" {
		return attachedDisk, nil
	}

	parts := strings.Split(d.Id(), ":")
	if len(parts) == 2 {
		return parts[1], nil
	}

	return "", fmt.Errorf("unable to determine the attached compute disk")
}

func findDiskByName(disks []*compute.AttachedDisk, id string) *compute.AttachedDisk {
	for _, disk := range disks {
		if compareSelfLinkOrResourceName("", disk.Source, id, nil) {
			return disk
		}
	}

	return nil
}
