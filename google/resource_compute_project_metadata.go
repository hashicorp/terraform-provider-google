package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeProjectMetadata() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeProjectMetadataCreateOrUpdate,
		Read:   resourceComputeProjectMetadataRead,
		Update: resourceComputeProjectMetadataCreateOrUpdate,
		Delete: resourceComputeProjectMetadataDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeProjectMetadataCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	md := &compute.Metadata{
		Items: expandComputeMetadata(d.Get("metadata").(map[string]interface{})),
	}

	err = resourceComputeProjectMetadataSet(projectID, config, md)
	if err != nil {
		return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
	}

	return resourceComputeProjectMetadataRead(d, meta)
}

func resourceComputeProjectMetadataRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.Id() == "" {
		projectID, err := getProject(d, config)
		if err != nil {
			return err
		}
		d.SetId(projectID)
	}

	// Load project service
	log.Printf("[DEBUG] Loading project service: %s", d.Id())
	project, err := config.clientCompute.Projects.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project metadata for project %q", d.Id()))
	}

	err = d.Set("metadata", flattenMetadata(project.CommonInstanceMetadata))
	if err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}

	d.Set("project", d.Id())
	d.SetId(d.Id())
	return nil
}

func resourceComputeProjectMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	md := &compute.Metadata{}
	err = resourceComputeProjectMetadataSet(projectID, config, md)
	if err != nil {
		return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
	}

	return resourceComputeProjectMetadataRead(d, meta)
}

func resourceComputeProjectMetadataSet(projectID string, config *Config, md *compute.Metadata) error {
	createMD := func() error {
		log.Printf("[DEBUG] Loading project service: %s", projectID)
		project, err := config.clientCompute.Projects.Get(projectID).Do()
		if err != nil {
			return fmt.Errorf("Error loading project '%s': %s", projectID, err)
		}

		md.Fingerprint = project.CommonInstanceMetadata.Fingerprint
		op, err := config.clientCompute.Projects.SetCommonInstanceMetadata(projectID, md).Do()
		if err != nil {
			return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
		}

		log.Printf("[DEBUG] SetCommonMetadata: %d (%s)", op.Id, op.SelfLink)
		return computeOperationWait(config.clientCompute, op, project.Name, "SetCommonMetadata")
	}

	err := MetadataRetryWrapper(createMD)
	return err
}
