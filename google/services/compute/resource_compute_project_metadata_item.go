// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/compute/v1"
)

type metadataPresentBehavior bool

const (
	failIfPresent    metadataPresentBehavior = true
	overwritePresent metadataPresentBehavior = false
)

func ResourceComputeProjectMetadataItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeProjectMetadataItemCreate,
		Read:   resourceComputeProjectMetadataItemRead,
		Update: resourceComputeProjectMetadataItemUpdate,
		Delete: resourceComputeProjectMetadataItemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The metadata key to set.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value to set for the given metadata key.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(7 * time.Minute),
			Update: schema.DefaultTimeout(7 * time.Minute),
			Delete: schema.DefaultTimeout(7 * time.Minute),
		},
		UseJSONNumber: true,
	}
}

func resourceComputeProjectMetadataItemCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	projectID, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	key := d.Get("key").(string)
	val := d.Get("value").(string)

	err = updateComputeCommonInstanceMetadata(config, projectID, key, userAgent, &val, d.Timeout(schema.TimeoutCreate), failIfPresent)
	if err != nil {
		return err
	}

	d.SetId(key)

	return nil
}

func resourceComputeProjectMetadataItemRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	projectID, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Loading project metadata: %s", projectID)
	project, err := config.NewComputeClient(userAgent).Projects.Get(projectID).Do()
	if err != nil {
		return fmt.Errorf("Error loading project '%s': %s", projectID, err)
	}

	md := FlattenMetadata(project.CommonInstanceMetadata)
	val, ok := md[d.Id()]
	if !ok {
		// Resource no longer exists
		d.SetId("")
		return nil
	}

	if err := d.Set("project", projectID); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("key", d.Id()); err != nil {
		return fmt.Errorf("Error setting key: %s", err)
	}
	if err := d.Set("value", val); err != nil {
		return fmt.Errorf("Error setting value: %s", err)
	}

	return nil
}

func resourceComputeProjectMetadataItemUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	projectID, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	if d.HasChange("value") {
		key := d.Get("key").(string)
		_, n := d.GetChange("value")
		new := n.(string)

		err = updateComputeCommonInstanceMetadata(config, projectID, key, userAgent, &new, d.Timeout(schema.TimeoutUpdate), overwritePresent)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceComputeProjectMetadataItemDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	projectID, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	key := d.Get("key").(string)

	err = updateComputeCommonInstanceMetadata(config, projectID, key, userAgent, nil, d.Timeout(schema.TimeoutDelete), overwritePresent)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func updateComputeCommonInstanceMetadata(config *transport_tpg.Config, projectID, key, userAgent string, afterVal *string, timeout time.Duration, failIfPresent metadataPresentBehavior) error {
	updateMD := func() error {
		lockName := fmt.Sprintf("projects/%s/commoninstancemetadata", projectID)
		transport_tpg.MutexStore.Lock(lockName)
		defer transport_tpg.MutexStore.Unlock(lockName)

		log.Printf("[DEBUG] Loading project metadata: %s", projectID)
		project, err := config.NewComputeClient(userAgent).Projects.Get(projectID).Do()
		if err != nil {
			return fmt.Errorf("Error loading project '%s': %s", projectID, err)
		}

		md := FlattenMetadata(project.CommonInstanceMetadata)

		val, ok := md[key]

		if !ok {
			if afterVal == nil {
				// Asked to set no value and we didn't find one - we're done
				return nil
			}
		} else {
			if failIfPresent {
				return fmt.Errorf("key %q already present in metadata for project %q. Use `terraform import` to manage it with Terraform", key, projectID)
			}
			if afterVal != nil && *afterVal == val {
				// Asked to set a value and it's already set - we're done.
				return nil
			}
		}

		if afterVal == nil {
			delete(md, key)
		} else {
			md[key] = *afterVal
		}

		// Attempt to write the new value now
		op, err := config.NewComputeClient(userAgent).Projects.SetCommonInstanceMetadata(
			projectID,
			&compute.Metadata{
				Fingerprint: project.CommonInstanceMetadata.Fingerprint,
				Items:       expandComputeMetadata(md),
			},
		).Do()

		if err != nil {
			return err
		}

		log.Printf("[DEBUG] SetCommonInstanceMetadata: %d (%s)", op.Id, op.SelfLink)

		return ComputeOperationWaitTime(config, op, project.Name, "SetCommonInstanceMetadata", userAgent, timeout)
	}

	return transport_tpg.MetadataRetryWrapper(updateMD)
}
