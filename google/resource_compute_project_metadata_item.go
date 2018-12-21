package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeProjectMetadataItem() *schema.Resource {
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceComputeProjectMetadataItemCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	key := d.Get("key").(string)
	val := d.Get("value").(string)

	err = updateComputeCommonInstanceMetadata(config, projectID, key, &val, int(d.Timeout(schema.TimeoutCreate).Minutes()))
	if err != nil {
		return err
	}

	d.SetId(key)

	return nil
}

func resourceComputeProjectMetadataItemRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Loading project metadata: %s", projectID)
	project, err := config.clientCompute.Projects.Get(projectID).Do()
	if err != nil {
		return fmt.Errorf("Error loading project '%s': %s", projectID, err)
	}

	md := flattenMetadata(project.CommonInstanceMetadata)
	val, ok := md[d.Id()]
	if !ok {
		// Resource no longer exists
		d.SetId("")
		return nil
	}

	d.Set("project", projectID)
	d.Set("key", d.Id())
	d.Set("value", val)

	return nil
}

func resourceComputeProjectMetadataItemUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	if d.HasChange("value") {
		key := d.Get("key").(string)
		_, n := d.GetChange("value")
		new := n.(string)

		err = updateComputeCommonInstanceMetadata(config, projectID, key, &new, int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceComputeProjectMetadataItemDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	key := d.Get("key").(string)

	err = updateComputeCommonInstanceMetadata(config, projectID, key, nil, int(d.Timeout(schema.TimeoutDelete).Minutes()))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func updateComputeCommonInstanceMetadata(config *Config, projectID string, key string, afterVal *string, timeout int) error {
	updateMD := func() error {
		log.Printf("[DEBUG] Loading project metadata: %s", projectID)
		project, err := config.clientCompute.Projects.Get(projectID).Do()
		if err != nil {
			return fmt.Errorf("Error loading project '%s': %s", projectID, err)
		}

		md := flattenMetadata(project.CommonInstanceMetadata)

		val, ok := md[key]

		if !ok {
			if afterVal == nil {
				// Asked to set no value and we didn't find one - we're done
				return nil
			}
		} else {
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
		op, err := config.clientCompute.Projects.SetCommonInstanceMetadata(
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

		return computeOperationWaitTime(config.clientCompute, op, project.Name, "SetCommonInstanceMetadata", timeout)
	}

	return MetadataRetryWrapper(updateMD)
}
