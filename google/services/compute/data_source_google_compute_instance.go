// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeInstance() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInstance().Schema)

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name", "self_link", "project", "zone")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, zone, name, err := tpgresource.GetZonalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	instance, err := config.NewComputeClient(userAgent).Instances.Get(project, zone, name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Instance %s", name))
	}

	md := flattenMetadataBeta(instance.Metadata)
	if err = d.Set("metadata", md); err != nil {
		return fmt.Errorf("error setting metadata: %s", err)
	}

	if err := d.Set("can_ip_forward", instance.CanIpForward); err != nil {
		return fmt.Errorf("Error setting can_ip_forward: %s", err)
	}
	if err := d.Set("machine_type", tpgresource.GetResourceNameFromSelfLink(instance.MachineType)); err != nil {
		return fmt.Errorf("Error setting machine_type: %s", err)
	}

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
		if err := d.Set("metadata_fingerprint", instance.Metadata.Fingerprint); err != nil {
			return fmt.Errorf("Error setting metadata_fingerprint: %s", err)
		}
	}

	// Set the tags fingerprint if there is one.
	if instance.Tags != nil {
		if err := d.Set("tags_fingerprint", instance.Tags.Fingerprint); err != nil {
			return fmt.Errorf("Error setting tags_fingerprint: %s", err)
		}
		if err := d.Set("tags", tpgresource.ConvertStringArrToInterface(instance.Tags.Items)); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	}

	if err := d.Set("labels", instance.Labels); err != nil {
		return err
	}

	if instance.LabelFingerprint != "" {
		if err := d.Set("label_fingerprint", instance.LabelFingerprint); err != nil {
			return fmt.Errorf("Error setting label_fingerprint: %s", err)
		}
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
				"source":      tpgresource.ConvertSelfLinkToV1(disk.Source),
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

	err = d.Set("shielded_instance_config", flattenShieldedVmConfig(instance.ShieldedInstanceConfig))
	if err != nil {
		return err
	}

	err = d.Set("enable_display", flattenEnableDisplay(instance.DisplayDevice))
	if err != nil {
		return err
	}

	if err := d.Set("attached_disk", ads); err != nil {
		return fmt.Errorf("Error setting attached_disk: %s", err)
	}
	if err := d.Set("cpu_platform", instance.CpuPlatform); err != nil {
		return fmt.Errorf("Error setting cpu_platform: %s", err)
	}
	if err := d.Set("min_cpu_platform", instance.MinCpuPlatform); err != nil {
		return fmt.Errorf("Error setting min_cpu_platform: %s", err)
	}
	if err := d.Set("deletion_protection", instance.DeletionProtection); err != nil {
		return fmt.Errorf("Error setting deletion_protection: %s", err)
	}
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(instance.SelfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("instance_id", fmt.Sprintf("%d", instance.Id)); err != nil {
		return fmt.Errorf("Error setting instance_id: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", tpgresource.GetResourceNameFromSelfLink(instance.Zone)); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("current_status", instance.Status); err != nil {
		return fmt.Errorf("Error setting current_status: %s", err)
	}
	if err := d.Set("name", instance.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, tpgresource.GetResourceNameFromSelfLink(instance.Zone), instance.Name))
	return nil
}
