package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	strcase "github.com/stoewer/go-strcase"
	computeBeta "google.golang.org/api/compute/v0.beta"
	compute "google.golang.org/api/compute/v1"
)

func resourceComputeInstanceFromTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceFromTemplateCreate,
		Read:   resourceComputeInstanceRead,
		Update: resourceComputeInstanceUpdate,
		Delete: resourceComputeInstanceDelete,

		// Import doesn't really make sense, because you could just import
		// as a google_compute_instance.

		Timeouts: resourceComputeInstance().Timeouts,

		Schema:        computeInstanceFromTemplateSchema(),
		CustomizeDiff: resourceComputeInstance().CustomizeDiff,
	}
}

func computeInstanceFromTemplateSchema() map[string]*schema.Schema {
	s := resourceComputeInstance().Schema

	for _, field := range []string{"boot_disk", "machine_type", "network_interface"} {
		// The user can set these fields as an override, but doesn't need to -
		// the template values will be used if they're unset.
		s[field].Required = false
		s[field].Optional = true
	}

	// schema.SchemaConfigModeAttr allows these fields to be removed in Terraform 0.12.
	// Passing field_name = [] in this mode differentiates between an intentionally empty
	// block vs an ignored computed block.
	nic := s["network_interface"].Elem.(*schema.Resource)
	nic.Schema["alias_ip_range"].ConfigMode = schema.SchemaConfigModeAttr
	nic.Schema["access_config"].ConfigMode = schema.SchemaConfigModeAttr

	for _, field := range []string{"attached_disk", "guest_accelerator", "service_account", "scratch_disk"} {
		s[field].ConfigMode = schema.SchemaConfigModeAttr
	}

	// Remove deprecated/removed fields that are never d.Set. We can't
	// programmatically remove all of them, because some of them still have d.Set
	// calls.
	for _, field := range []string{"disk", "network"} {
		delete(s, field)
	}

	recurseOnSchema(s, func(field *schema.Schema) {
		// We don't want to accidentally use default values to override the instance
		// template, so remove defaults.
		field.Default = nil

		// Make non-required fields computed since they'll be set by the template.
		// Leave deprecated and removed fields alone because we don't set them.
		if !field.Required && !(field.Deprecated != "" || field.Removed != "") {
			field.Computed = true
		}
	})

	s["source_instance_template"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}

	return s
}

func recurseOnSchema(s map[string]*schema.Schema, f func(*schema.Schema)) {
	for _, field := range s {
		f(field)
		if e := field.Elem; e != nil {
			if r, ok := e.(*schema.Resource); ok {
				recurseOnSchema(r.Schema, f)
			}
		}
	}
}

func resourceComputeInstanceFromTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the zone
	z, err := getZone(d, config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Loading zone: %s", z)
	zone, err := config.clientCompute.Zones.Get(project, z).Do()
	if err != nil {
		return fmt.Errorf("Error loading zone '%s': %s", z, err)
	}

	instance, err := expandComputeInstance(project, d, config)
	if err != nil {
		return err
	}

	tpl, err := ParseInstanceTemplateFieldValue(d.Get("source_instance_template").(string), d, config)
	if err != nil {
		return err
	}

	it, err := config.clientComputeBeta.InstanceTemplates.Get(project, tpl.Name).Do()
	if err != nil {
		return err
	}

	instance.Disks, err = adjustInstanceFromTemplateDisks(d, config, it, zone, project)
	if err != nil {
		return err
	}

	// Force send all top-level fields that have been set in case they're overridden to zero values.
	// Initialize ForceSendFields to empty so we don't get things that the instance resource
	// always force-sends.
	instance.ForceSendFields = []string{}
	for f, s := range computeInstanceFromTemplateSchema() {
		// It seems that GetOkExists always returns true for sets.
		// TODO: confirm this and file issue against Terraform core.
		// In the meantime, don't force send sets.
		if s.Type == schema.TypeSet {
			continue
		}

		if _, exists := d.GetOkExists(f); exists {
			// Assume for now that all fields are exact snake_case versions of the API fields.
			// This won't necessarily always be true, but it serves as a good approximation and
			// can be adjusted later as we discover issues.
			instance.ForceSendFields = append(instance.ForceSendFields, strcase.UpperCamelCase(f))
		}
	}

	log.Printf("[INFO] Requesting instance creation")
	op, err := config.clientComputeBeta.Instances.Insert(project, zone.Name, instance).SourceInstanceTemplate(tpl.RelativeLink()).Do()
	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, z, instance.Name))

	// Wait for the operation to complete
	waitErr := computeOperationWaitTime(config, op, project,
		"instance to create", int(d.Timeout(schema.TimeoutCreate).Minutes()))
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	return resourceComputeInstanceRead(d, meta)
}

// Instances have disks spread across multiple schema properties. This function
// ensures that overriding one of these properties does not override the others.
func adjustInstanceFromTemplateDisks(d *schema.ResourceData, config *Config, it *computeBeta.InstanceTemplate, zone *compute.Zone, project string) ([]*computeBeta.AttachedDisk, error) {
	disks := []*computeBeta.AttachedDisk{}
	if _, hasBootDisk := d.GetOk("boot_disk"); hasBootDisk {
		bootDisk, err := expandBootDisk(d, config, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, bootDisk)
	} else {
		// boot disk was not overridden, so use the one from the instance template
		for _, disk := range it.Properties.Disks {
			if disk.Boot {
				if disk.Source != "" {
					// Instances need a URL for the disk, but instance templates only have the name
					disk.Source = fmt.Sprintf("projects/%s/zones/%s/disks/%s", project, zone.Name, disk.Source)
				}
				if disk.InitializeParams != nil {
					if dt := disk.InitializeParams.DiskType; dt != "" {
						// Instances need a URL for the disk type, but instance templates
						// only have the name (since they're global).
						disk.InitializeParams.DiskType = fmt.Sprintf("zones/%s/diskTypes/%s", zone.Name, dt)
					}
				}
				disks = append(disks, disk)
				break
			}
		}
	}

	if _, hasScratchDisk := d.GetOk("scratch_disk"); hasScratchDisk {
		scratchDisks, err := expandScratchDisks(d, config, project)
		if err != nil {
			return nil, err
		}
		disks = append(disks, scratchDisks...)
	} else {
		// scratch disks were not overridden, so use the ones from the instance template
		for _, disk := range it.Properties.Disks {
			if disk.Type == "SCRATCH" {
				if disk.InitializeParams != nil {
					if dt := disk.InitializeParams.DiskType; dt != "" {
						// Instances need a URL for the disk type, but instance templates
						// only have the name (since they're global).
						disk.InitializeParams.DiskType = fmt.Sprintf("zones/%s/diskTypes/%s", zone.Name, dt)
					}
				}
				disks = append(disks, disk)
			}
		}
	}

	attachedDisksCount := d.Get("attached_disk.#").(int)
	if attachedDisksCount > 0 {
		for i := 0; i < attachedDisksCount; i++ {
			diskConfig := d.Get(fmt.Sprintf("attached_disk.%d", i)).(map[string]interface{})
			disk, err := expandAttachedDisk(diskConfig, d, config)
			if err != nil {
				return nil, err
			}

			disks = append(disks, disk)
		}
	} else {
		// attached disks were not overridden, so use the ones from the instance template
		for _, disk := range it.Properties.Disks {
			if !disk.Boot && disk.Type != "SCRATCH" {
				if s := disk.Source; s != "" {
					// Instances need a URL for the disk source, but instance templates
					// only have the name (since they're global).
					disk.Source = fmt.Sprintf("zones/%s/disks/%s", zone.Name, s)
				}
				if disk.InitializeParams != nil {
					if dt := disk.InitializeParams.DiskType; dt != "" {
						// Instances need a URL for the disk type, but instance templates
						// only have the name (since they're global).
						disk.InitializeParams.DiskType = fmt.Sprintf("zones/%s/diskTypes/%s", zone.Name, dt)
					}
				}
				disks = append(disks, disk)
			}
		}
	}

	return disks, nil
}
