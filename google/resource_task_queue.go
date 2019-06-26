package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/cloudtasks/v2"
)

func resourceTaskQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceTaskQueueCreateUpdate,
		Read:   resourceTaskQueueRead,
		Update: resourceTaskQueueCreateUpdate,
		Delete: resourceTaskQueueDelete,
		Importer: &schema.ResourceImporter{
			State: resourceTaskQueueImportState,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"region"},
			},
			"region": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"location"},
			},
			"rate_limits": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_burst_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_dispatches_per_second": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
						"max_concurrent_dispatches": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"retry": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_attempts": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_doublings": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"max_backoff": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateRegexp(`^[0-9]*s$`),
						},
						"min_backoff": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateRegexp(`^[0-9]*\.?[0-9]+s$`),
						},
					},
				},
			},
		},
	}
}

func resourceTaskQueueRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location := d.Get("location")
	name := d.Get("name")

	url := fmt.Sprintf("projects/%s/locations/%s/queues/%s", project, location, name)

	resp, err := config.clientCloudTasks.Projects.Locations.Queues.Get(url).Do()
	if err != nil {
		return fmt.Errorf("Error reading task queue: %s %s", location, err)
	}

	d.Set("retry", flattenRetryConfig(resp.RetryConfig))
	d.Set("rate_limits", flattenRateLimits(resp.RateLimits))

	return nil
}

func resourceTaskQueueCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location := d.Get("location")
	name := d.Get("name")

	url := fmt.Sprintf("projects/%s/locations/%s/queues/%s", project, location, name)
	id := fmt.Sprintf("%s/%s/%s", project, location, name)

	rb := &cloudtasks.Queue{
		Name:        url,
		RateLimits:  expandRateLimits(d.Get("rate_limits").([]interface{})),
		RetryConfig: expandRetryConfig(d.Get("retry").([]interface{})),
	}

	log.Printf("[DEBUG] Updating Task Queue: %#v", name)

	resp, err := config.clientCloudTasks.Projects.Locations.Queues.Patch(url, rb).Do()
	if err != nil {
		return fmt.Errorf("Error updating task queue %s: %s", name, err)
	}

	d.SetId(id)

	log.Printf("[DEBUG] Finished updating new Task Queue: %#v\n%#v\n", name, resp)

	return resourceTaskQueueRead(d, meta)
}

func resourceTaskQueueDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location := d.Get("location")
	name := d.Get("name")

	url := fmt.Sprintf("projects/%s/locations/%s/queues/%s", project, location, name)

	log.Printf("[DEBUG] Deleting Task Queue: %#v", name)

	resp, err := config.clientCloudTasks.Projects.Locations.Queues.Delete(url).Do()
	if err != nil {
		return fmt.Errorf("Error deleting task queue: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Task Queue: %#v\n%#v\n", name, resp)

	return nil
}

func resourceTaskQueueImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import id %q. Expecting {project}/{location}/{queue}", d.Id())
	}

	d.Set("location", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}

func flattenRateLimits(in *cloudtasks.RateLimits) []interface{} {
	m := make(map[string]interface{})

	if in != nil {
		m["max_dispatches_per_second"] = in.MaxDispatchesPerSecond
		m["max_concurrent_dispatches"] = in.MaxConcurrentDispatches
		m["max_burst_size"] = in.MaxBurstSize
	}

	return []interface{}{m}
}

func expandRateLimits(configured interface{}) *cloudtasks.RateLimits {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	rateLimits := l[0].(map[string]interface{})

	return &cloudtasks.RateLimits{
		MaxBurstSize:            int64(rateLimits["max_burst_size"].(int)),
		MaxConcurrentDispatches: int64(rateLimits["max_concurrent_dispatches"].(int)),
		MaxDispatchesPerSecond:  rateLimits["max_dispatches_per_second"].(float64),
	}
}

func flattenRetryConfig(in *cloudtasks.RetryConfig) []interface{} {
	m := make(map[string]interface{})

	if in != nil {
		m["max_attempts"] = in.MaxAttempts
		m["max_doublings"] = in.MaxDoublings
		m["max_backoff"] = in.MaxBackoff
		m["min_backoff"] = in.MinBackoff
	}

	return []interface{}{m}
}

func expandRetryConfig(configured interface{}) *cloudtasks.RetryConfig {
	l := configured.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	retry := l[0].(map[string]interface{})

	return &cloudtasks.RetryConfig{
		MaxAttempts:  int64(retry["max_attempts"].(int)),
		MaxDoublings: int64(retry["max_doublings"].(int)),
		MaxBackoff:   retry["max_backoff"].(string),
		MinBackoff:   retry["min_backoff"].(string),
	}
}
