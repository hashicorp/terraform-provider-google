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

func ResourceProjectUsageBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectUsageBucketCreate,
		Read:   resourceProjectUsageBucketRead,
		Delete: resourceProjectUsageBucketDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectUsageBucketImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The bucket to store reports in.`,
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `A prefix for the reports, for instance, the project name.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The project to set the export bucket on. If it is not provided, the provider project is used.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceProjectUsageBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	p, err := config.NewComputeClient(userAgent).Projects.Get(project).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project data for project %s", project))
	}

	if p.UsageExportLocation == nil {
		log.Printf("[WARN] Removing usage export location resource %s because it's not enabled server-side.", project)
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("prefix", p.UsageExportLocation.ReportNamePrefix); err != nil {
		return fmt.Errorf("Error setting prefix: %s", err)
	}
	if err := d.Set("bucket_name", p.UsageExportLocation.BucketName); err != nil {
		return fmt.Errorf("Error setting bucket_name: %s", err)
	}
	return nil
}

func resourceProjectUsageBucketCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.NewComputeClient(userAgent).Projects.SetUsageExportBucket(project, &compute.UsageExportLocation{
		ReportNamePrefix: d.Get("prefix").(string),
		BucketName:       d.Get("bucket_name").(string),
	}).Do()
	if err != nil {
		return err
	}
	d.SetId(project)
	err = ComputeOperationWaitTime(config, op, project, "Setting usage export bucket.", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		d.SetId("")
		return err
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return resourceProjectUsageBucketRead(d, meta)
}

func resourceProjectUsageBucketDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.NewComputeClient(userAgent).Projects.SetUsageExportBucket(project, nil).Do()
	if err != nil {
		return err
	}

	err = ComputeOperationWaitTime(config, op, project,
		"Setting usage export bucket to nil, automatically disabling usage export.", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceProjectUsageBucketImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	project := d.Id()
	if err := d.Set("project", project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
