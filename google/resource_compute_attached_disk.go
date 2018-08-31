package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"auto_delete": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"device_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"interface": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "SCSI",
				ValidateFunc: validation.StringInSlice([]string{"NVME", "SCSI"}, false),
			},
			"mode": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "READ_WRITE",
				ValidateFunc: validation.StringInSlice([]string{"READ_ONLY", "READ_WRITE"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "PERSISTENT",
				ValidateFunc: validation.StringInSlice([]string{"SCRATCH", "PERSISTENT"}, false),
			},
		},
	}
}

func resourceAttachedDiskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProjectForAttachedDisk(d, config)
	if err != nil {
		return err
	}
	zone, err := getZoneForAttachedDisk(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance").(string))
	diskName := GetResourceNameFromSelfLink(d.Get("disk").(string))
	autoDelete := d.Get("auto_delete").(bool)
	diskInterface := d.Get("interface").(string)
	mode := d.Get("mode").(string)
	diskType := d.Get("type").(string)

	attachedDisk := compute.AttachedDisk{
		Source:     fmt.Sprintf("/projects/%s/zones/%s/disks/%s", project, zone, diskName),
		AutoDelete: autoDelete,
		Interface:  diskInterface,
		Mode:       mode,
		Type:       diskType,
	}

	deviceName := d.Get("device_name").(string)
	if deviceName != "" {
		attachedDisk.DeviceName = deviceName
	}

	op, err := config.clientCompute.Instances.AttachDisk(project, zone, instanceName, &attachedDisk).Do()
	if err != nil {
		return err
	}

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

	project, err := getProjectForAttachedDisk(d, config)
	if err != nil {
		return err
	}
	d.Set("project", project)

	zone, err := getZoneForAttachedDisk(d, config)
	if err != nil {
		return err
	}
	d.Set("zone", zone)

	instanceName := GetResourceNameFromSelfLink(d.Get("instance").(string))
	diskName := GetResourceNameFromSelfLink(d.Get("disk").(string))

	instance, err := config.clientCompute.Instances.Get(project, zone, instanceName).Do()
	if err != nil {
		return err
	}

	// Iterate through the instance's attached disk as this is the only way to
	// confirm the disk is actually attached
	ad := findDiskByName(instance.Disks, diskName)
	if ad == nil {
		log.Printf("[WARN] Refereecned disk wasn't found attached to this compute instance. Unsetting resource id.")
		d.SetId("")
		return nil
	}

	d.Set("device_name", ad.DeviceName)
	d.Set("auto_delete", ad.AutoDelete)
	d.Set("interface", ad.Interface)
	d.Set("mode", ad.Mode)
	d.Set("type", ad.Type)

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
	log.Printf("ZOMG destroying")
	config := meta.(*Config)

	project, err := getProjectForAttachedDisk(d, config)
	if err != nil {
		return err
	}
	zone, err := getZoneForAttachedDisk(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance").(string))
	diskName := GetResourceNameFromSelfLink(d.Get("disk").(string))

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

	log.Printf("ZOMG calling detach")
	op, err := config.clientCompute.Instances.DetachDisk(project, zone, instanceName, ad.DeviceName).Do()
	if err != nil {
		return err
	}

	waitErr := computeSharedOperationWaitTime(config.clientCompute, op, project,
		int(d.Timeout(schema.TimeoutDelete).Minutes()), fmt.Sprintf("Detaching disk from %s", instanceName))
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
		return nil, fmt.Errorf("unable to determine attached disk id - id should be 'google_compute_instance.name:google_compute_disk.name'")
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

// getZoneForAttachedDisk prioritizes the zone defined by the compute instance self link before standard logic
func getZoneForAttachedDisk(d *schema.ResourceData, c *Config) (string, error) {
	attachedInstance := d.Get("instance").(string)

	zone, err := GetPathVariableFromSelfLink(attachedInstance, "zone")
	if err == nil {
		return zone, nil
	}

	// If zone can't be inferred from the compute instance self link, fall back to project
	zone, err = getZone(d, c)
	if err != nil {
		return "", fmt.Errorf("%s to inherit from the attached instance compute use `self_link` instead of `name`", err)
	}

	return zone, nil
}

// getProjectForAttachedDisk prioritizes the project defined by the compute instance self link before standard logic
func getProjectForAttachedDisk(d *schema.ResourceData, c *Config) (string, error) {
	attachedInstance := d.Get("instance").(string)

	project, err := GetPathVariableFromSelfLink(attachedInstance, "project")
	if err == nil {
		return project, nil
	}

	// If project can't be inferred from the compute instance self link, fall back to project
	project, err = getProject(d, c)
	if err != nil {
		return "", fmt.Errorf("%s to inherit from the attached compute instance use `self_link` instead of `name`", err)
	}

	return project, nil
}
