package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceProjectUsageBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectUsageBucketCreate,
		Read:   resourceProjectUsageBucketRead,
		Delete: resourceProjectUsageBucketDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectUsageBucketImportState,
		},

		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceProjectUsageBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	p, err := config.clientCompute.Projects.Get(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project data for project %s", project))
	}

	if p.UsageExportLocation == nil {
		log.Printf("[WARN] Removing usage export location resource %s because it's not enabled server-side.", project)
		d.SetId("")
	}

	d.Set("project", project)
	d.Set("prefix", p.UsageExportLocation.ReportNamePrefix)
	d.Set("bucket_name", p.UsageExportLocation.BucketName)
	return nil
}

func resourceProjectUsageBucketCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.Projects.SetUsageExportBucket(project, &compute.UsageExportLocation{
		ReportNamePrefix: d.Get("prefix").(string),
		BucketName:       d.Get("bucket_name").(string),
	}).Do()
	if err != nil {
		return err
	}
	d.SetId(project)
	err = computeOperationWait(config.clientCompute, op, project, "Setting usage export bucket.")
	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("project", project)

	return resourceProjectUsageBucketRead(d, meta)
}

func resourceProjectUsageBucketDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.Projects.SetUsageExportBucket(project, nil).Do()
	if err != nil {
		return err
	}

	err = computeOperationWait(config.clientCompute, op, project,
		"Setting usage export bucket to nil, automatically disabling usage export.")
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func resourceProjectUsageBucketImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	project := d.Id()
	d.Set("project", project)
	return []*schema.ResourceData{d}, nil
}
