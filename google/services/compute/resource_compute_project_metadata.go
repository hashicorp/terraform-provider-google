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

func ResourceComputeProjectMetadata() *schema.Resource {
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
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A series of key value pairs.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceComputeProjectMetadataCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	projectID, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	md := &compute.Metadata{
		Items: expandComputeMetadata(d.Get("metadata").(map[string]interface{})),
	}

	err = resourceComputeProjectMetadataSet(projectID, userAgent, config, md, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
	}

	d.SetId(projectID)

	return resourceComputeProjectMetadataRead(d, meta)
}

func resourceComputeProjectMetadataRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	// At import time, we have no state to draw from. We'll wrongly pull the
	// provider default project if we use a normal GetProject, so we need to
	// rely on the `id` field being set to the project.
	// At any other time we can use GetProject, as state will have the correct
	// value; the project pulled from config / the provider / at import time.
	//
	// Note that if a user imports a project other than their provider project
	// and has left the project field unspecified, Terraform will not see a diff
	// but would create metadata for the provider project on a destroy/create.
	projectId := d.Id()

	project, err := config.NewComputeClient(userAgent).Projects.Get(projectId).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project metadata for project %q", projectId))
	}

	err = d.Set("metadata", FlattenMetadata(project.CommonInstanceMetadata))
	if err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}

	if err := d.Set("project", projectId); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

func resourceComputeProjectMetadataDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	projectID, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	md := &compute.Metadata{}
	err = resourceComputeProjectMetadataSet(projectID, userAgent, config, md, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
	}

	return resourceComputeProjectMetadataRead(d, meta)
}

func resourceComputeProjectMetadataSet(projectID, userAgent string, config *transport_tpg.Config, md *compute.Metadata, timeout time.Duration) error {
	createMD := func() error {
		log.Printf("[DEBUG] Loading project service: %s", projectID)
		project, err := config.NewComputeClient(userAgent).Projects.Get(projectID).Do()
		if err != nil {
			return fmt.Errorf("Error loading project '%s': %s", projectID, err)
		}

		md.Fingerprint = project.CommonInstanceMetadata.Fingerprint
		op, err := config.NewComputeClient(userAgent).Projects.SetCommonInstanceMetadata(projectID, md).Do()
		if err != nil {
			return fmt.Errorf("SetCommonInstanceMetadata failed: %s", err)
		}

		log.Printf("[DEBUG] SetCommonMetadata: %d (%s)", op.Id, op.SelfLink)
		return ComputeOperationWaitTime(config, op, project.Name, "SetCommonMetadata", userAgent, timeout)
	}

	err := transport_tpg.MetadataRetryWrapper(createMD)
	return err
}
