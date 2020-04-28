package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project": {
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

	err = resourceComputeProjectMetadataSet(projectID, config, md, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
	}

	d.SetId(projectID)

	return resourceComputeProjectMetadataRead(d, meta)
}

func resourceComputeProjectMetadataRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// At import time, we have no state to draw from. We'll wrongly pull the
	// provider default project if we use a normal getProject, so we need to
	// rely on the `id` field being set to the project.
	// At any other time we can use getProject, as state will have the correct
	// value; the project pulled from config / the provider / at import time.
	//
	// Note that if a user imports a project other than their provider project
	// and has left the project field unspecified, Terraform will not see a diff
	// but would create metadata for the provider project on a destroy/create.
	projectId := d.Id()

	project, err := config.clientCompute.Projects.Get(projectId).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project metadata for project %q", projectId))
	}

	err = d.Set("metadata", flattenMetadata(project.CommonInstanceMetadata))
	if err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}

	d.Set("project", projectId)

	return nil
}

func resourceComputeProjectMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectID, err := getProject(d, config)
	if err != nil {
		return err
	}

	md := &compute.Metadata{}
	err = resourceComputeProjectMetadataSet(projectID, config, md, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
	}

	return resourceComputeProjectMetadataRead(d, meta)
}

func resourceComputeProjectMetadataSet(projectID string, config *Config, md *compute.Metadata, timeout time.Duration) error {
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
		return computeOperationWaitTime(config, op, project.Name, "SetCommonMetadata", timeout)
	}

	err := MetadataRetryWrapper(createMD)
	return err
}
