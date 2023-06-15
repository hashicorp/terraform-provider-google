// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"google.golang.org/api/compute/v1"
)

func ResourceComputeAttachedDisk() *schema.Resource {
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
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
			},
			"instance": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      `name or self_link of the compute instance that the disk will be attached to. If the self_link is provided then zone and project are extracted from the self link. If only the name is used then zone and project must be defined as properties on the resource or provider.`,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
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
		UseJSONNumber: true,
	}
}

func resourceAttachedDiskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	zv, err := tpgresource.ParseZonalFieldValue("instances", d.Get("instance").(string), "project", "zone", d, config, false)
	if err != nil {
		return err
	}

	disk := d.Get("disk").(string)
	diskName := tpgresource.GetResourceNameFromSelfLink(disk)
	diskSrc := fmt.Sprintf("projects/%s/zones/%s/disks/%s", zv.Project, zv.Zone, diskName)

	// Check if the disk is a regional disk
	if strings.Contains(disk, "regions") {
		rv, err := tpgresource.ParseRegionDiskFieldValue(disk, d, config)
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

	op, err := config.NewComputeClient(userAgent).Instances.AttachDisk(zv.Project, zv.Zone, zv.Name, &attachedDisk).Do()
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s/%s", zv.Project, zv.Zone, zv.Name, diskName))

	waitErr := ComputeOperationWaitTime(config, op, zv.Project,
		"disk to attach", userAgent, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		d.SetId("")
		return waitErr
	}

	return resourceAttachedDiskRead(d, meta)
}

func resourceAttachedDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	zv, err := tpgresource.ParseZonalFieldValue("instances", d.Get("instance").(string), "project", "zone", d, config, false)
	if err != nil {
		return err
	}
	if err := d.Set("project", zv.Project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", zv.Zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}

	diskName := tpgresource.GetResourceNameFromSelfLink(d.Get("disk").(string))

	instance, err := config.NewComputeClient(userAgent).Instances.Get(zv.Project, zv.Zone, zv.Name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("AttachedDisk %q", d.Id()))
	}

	// Iterate through the instance's attached disks as this is the only way to
	// confirm the disk is actually attached
	ad := FindDiskByName(instance.Disks, diskName)
	if ad == nil {
		log.Printf("[WARN] Referenced disk wasn't found attached to this compute instance. Removing from state.")
		d.SetId("")
		return nil
	}

	if err := d.Set("device_name", ad.DeviceName); err != nil {
		return fmt.Errorf("Error setting device_name: %s", err)
	}
	if err := d.Set("mode", ad.Mode); err != nil {
		return fmt.Errorf("Error setting mode: %s", err)
	}

	// Force the referenced resources to a self-link in state because it's more specific then name.
	instancePath, err := tpgresource.GetRelativePath(instance.SelfLink)
	if err != nil {
		return err
	}
	if err := d.Set("instance", instancePath); err != nil {
		return fmt.Errorf("Error setting instance: %s", err)
	}
	diskPath, err := tpgresource.GetRelativePath(ad.Source)
	if err != nil {
		return err
	}
	if err := d.Set("disk", diskPath); err != nil {
		return fmt.Errorf("Error setting disk: %s", err)
	}

	return nil
}

func resourceAttachedDiskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	zv, err := tpgresource.ParseZonalFieldValue("instances", d.Get("instance").(string), "project", "zone", d, config, false)
	if err != nil {
		return err
	}

	diskName := tpgresource.GetResourceNameFromSelfLink(d.Get("disk").(string))

	instance, err := config.NewComputeClient(userAgent).Instances.Get(zv.Project, zv.Zone, zv.Name).Do()
	if err != nil {
		return err
	}

	// Confirm the disk is still attached before making the call to detach it. If the disk isn't listed as an attached
	// disk on the compute instance then return as though the delete call succeed since this is the desired state.
	ad := FindDiskByName(instance.Disks, diskName)
	if ad == nil {
		return nil
	}

	op, err := config.NewComputeClient(userAgent).Instances.DetachDisk(zv.Project, zv.Zone, zv.Name, ad.DeviceName).Do()
	if err != nil {
		return err
	}

	waitErr := ComputeOperationWaitTime(config, op, zv.Project,
		fmt.Sprintf("Detaching disk from %s", zv.Name), userAgent, d.Timeout(schema.TimeoutDelete))
	if waitErr != nil {
		return waitErr
	}

	return nil
}

func resourceAttachedDiskImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	err := tpgresource.ParseImportId(
		[]string{"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/(?P<instance>[^/]+)/(?P<disk>[^/]+)",
			"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<instance>[^/]+)/(?P<disk>[^/]+)"}, d, config)
	if err != nil {
		return nil, err
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/zones/{{zone}}/instances/{{instance}}/{{disk}}")
	if err != nil {
		return nil, err
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func FindDiskByName(disks []*compute.AttachedDisk, id string) *compute.AttachedDisk {
	for _, disk := range disks {
		if tpgresource.CompareSelfLinkOrResourceName("", disk.Source, id, nil) {
			return disk
		}
	}

	return nil
}
