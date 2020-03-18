package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleComputeInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeInstance().Schema)

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "name", "self_link", "project", "zone")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, zone, name, err := GetZonalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	instance, err := config.clientComputeBeta.Instances.Get(project, zone, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Instance %s", name))
	}

	md := flattenMetadataBeta(instance.Metadata)
	if err = d.Set("metadata", md); err != nil {
		return fmt.Errorf("error setting metadata: %s", err)
	}

	d.Set("can_ip_forward", instance.CanIpForward)
	d.Set("machine_type", GetResourceNameFromSelfLink(instance.MachineType))

	// Set the networks
	// Use the first external IP found for the default connection info.
	networkInterfaces, _, internalIP, externalIP, err := flattenNetworkInterfaces(d, config, instance.NetworkInterfaces)
	if err != nil {
		return err
	}
	if err := d.Set("network_interface", networkInterfaces); err != nil {
		return err
	}

	// Fall back on internal ip if there is no external ip.  This makes sense in the situation where
	// terraform is being used on a cloud instance and can therefore access the instances it creates
	// via their internal ips.
	sshIP := externalIP
	if sshIP == "" {
		sshIP = internalIP
	}

	// Initialize the connection info
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": sshIP,
	})

	// Set the metadata fingerprint if there is one.
	if instance.Metadata != nil {
		d.Set("metadata_fingerprint", instance.Metadata.Fingerprint)
	}

	// Set the tags fingerprint if there is one.
	if instance.Tags != nil {
		d.Set("tags_fingerprint", instance.Tags.Fingerprint)
		d.Set("tags", convertStringArrToInterface(instance.Tags.Items))
	}

	if err := d.Set("labels", instance.Labels); err != nil {
		return err
	}

	if instance.LabelFingerprint != "" {
		d.Set("label_fingerprint", instance.LabelFingerprint)
	}

	attachedDisks := []map[string]interface{}{}
	scratchDisks := []map[string]interface{}{}
	for _, disk := range instance.Disks {
		if disk.Boot {
			err = d.Set("boot_disk", flattenBootDisk(d, disk, config))
			if err != nil {
				return err
			}
		} else if disk.Type == "SCRATCH" {
			scratchDisks = append(scratchDisks, flattenScratchDisk(disk))
		} else {
			di := map[string]interface{}{
				"source":      ConvertSelfLinkToV1(disk.Source),
				"device_name": disk.DeviceName,
				"mode":        disk.Mode,
			}
			if key := disk.DiskEncryptionKey; key != nil {
				di["disk_encryption_key_sha256"] = key.Sha256
				di["kms_key_self_link"] = key.KmsKeyName
			}
			attachedDisks = append(attachedDisks, di)
		}
	}
	// Remove nils from map in case there were disks in the config that were not present on read;
	// i.e. a disk was detached out of band
	ads := []map[string]interface{}{}
	for _, d := range attachedDisks {
		if d != nil {
			ads = append(ads, d)
		}
	}

	err = d.Set("service_account", flattenServiceAccounts(instance.ServiceAccounts))
	if err != nil {
		return err
	}

	err = d.Set("scheduling", flattenScheduling(instance.Scheduling))
	if err != nil {
		return err
	}

	err = d.Set("guest_accelerator", flattenGuestAccelerators(instance.GuestAccelerators))
	if err != nil {
		return err
	}

	err = d.Set("scratch_disk", scratchDisks)
	if err != nil {
		return err
	}

	err = d.Set("shielded_instance_config", flattenShieldedVmConfig(instance.ShieldedVmConfig))
	if err != nil {
		return err
	}

	err = d.Set("enable_display", flattenEnableDisplay(instance.DisplayDevice))
	if err != nil {
		return err
	}

	d.Set("attached_disk", ads)
	d.Set("cpu_platform", instance.CpuPlatform)
	d.Set("min_cpu_platform", instance.MinCpuPlatform)
	d.Set("deletion_protection", instance.DeletionProtection)
	d.Set("self_link", ConvertSelfLinkToV1(instance.SelfLink))
	d.Set("instance_id", fmt.Sprintf("%d", instance.Id))
	d.Set("project", project)
	d.Set("zone", GetResourceNameFromSelfLink(instance.Zone))
	d.Set("current_status", instance.Status)
	d.Set("name", instance.Name)
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, instance.Zone, instance.Name))
	return nil
}
