package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

var instancesSelfLinkPattern = regexp.MustCompile(fmt.Sprintf(zonalLinkBasePattern, "instances"))

func resourceComputeTargetPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeTargetPoolCreate,
		Read:   resourceComputeTargetPoolRead,
		Delete: resourceComputeTargetPoolDelete,
		Update: resourceComputeTargetPoolUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceTargetPoolStateImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backup_pool": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"failover_ratio": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},

			"health_checks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
			},

			"instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(v interface{}) string {
						return canonicalizeInstanceRef(v.(string))
					},
				},
				Set: func(v interface{}) int {
					return schema.HashString(canonicalizeInstanceRef(v.(string)))
				},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"session_affinity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "NONE",
			},
		},
	}
}

func canonicalizeInstanceRef(instanceRef string) string {
	// instances can also be specified in the config as a URL or <zone>/<project>
	parts := instancesSelfLinkPattern.FindStringSubmatch(instanceRef)
	// parts[0] = full match
	// parts[1] = project
	// parts[2] = zone
	// parts[3] = instance name

	if len(parts) < 4 {
		return instanceRef
	}

	return fmt.Sprintf("%s/%s", parts[2], parts[3])
	// return fmt.Sprintf("%s/%s/%s", parts[1], parts[2], parts[3])
}

// Healthchecks need to exist before being referred to from the target pool.
func convertHealthChecks(healthChecks []interface{}, d *schema.ResourceData, config *Config) ([]string, error) {
	if len(healthChecks) == 0 {
		return []string{}, nil
	}

	hc, err := ParseHttpHealthCheckFieldValue(healthChecks[0].(string), d, config)
	if err != nil {
		return nil, err
	}

	return []string{hc.RelativeLink()}, nil
}

// Instances do not need to exist yet, so we simply generate URLs.
// Instances can be full URLS or zone/name
func convertInstancesToUrls(d *schema.ResourceData, config *Config, project string, names *schema.Set) ([]string, error) {
	urls := make([]string, len(names.List()))
	for i, nameI := range names.List() {
		name := nameI.(string)
		// assume that any URI will start with https://
		if strings.HasPrefix(name, "https://") {
			urls[i] = name
		} else {
			splitName := strings.Split(name, "/")
			if len(splitName) != 2 {
				return nil, fmt.Errorf("Invalid instance name, require URL or zone/name: %s", name)
			} else {
				url, err := replaceVars(d, config, fmt.Sprintf(
					"{{ComputeBasePath}}projects/%s/zones/%s/instances/%s",
					project, splitName[0], splitName[1]))
				if err != nil {
					return nil, err
				}
				urls[i] = url
			}
		}
	}
	return urls, nil
}

func resourceComputeTargetPoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	hchkUrls, err := convertHealthChecks(d.Get("health_checks").([]interface{}), d, config)
	if err != nil {
		return err
	}

	instanceUrls, err := convertInstancesToUrls(d, config, project, d.Get("instances").(*schema.Set))
	if err != nil {
		return err
	}

	// Build the parameter
	tpool := &compute.TargetPool{
		BackupPool:      d.Get("backup_pool").(string),
		Description:     d.Get("description").(string),
		HealthChecks:    hchkUrls,
		Instances:       instanceUrls,
		Name:            d.Get("name").(string),
		SessionAffinity: d.Get("session_affinity").(string),
	}
	if d.Get("failover_ratio") != nil {
		tpool.FailoverRatio = d.Get("failover_ratio").(float64)
	}
	log.Printf("[DEBUG] TargetPool insert request: %#v", tpool)
	op, err := config.clientCompute.TargetPools.Insert(
		project, region, tpool).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 && strings.Contains(gerr.Message, "httpHealthChecks") {
			return fmt.Errorf("Health check %s is not a valid HTTP health check", d.Get("health_checks").([]interface{})[0])
		}
		return fmt.Errorf("Error creating TargetPool: %s", err)
	}

	// It probably maybe worked, so store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/targetPools/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(config, op, project, "Creating Target Pool", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	return resourceComputeTargetPoolRead(d, meta)
}

func resourceComputeTargetPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	d.Partial(true)

	if d.HasChange("health_checks") {

		from_, to_ := d.GetChange("health_checks")
		fromUrls, err := convertHealthChecks(from_.([]interface{}), d, config)
		if err != nil {
			return err
		}
		toUrls, err := convertHealthChecks(to_.([]interface{}), d, config)
		if err != nil {
			return err
		}
		add, remove := calcAddRemove(fromUrls, toUrls)

		removeReq := &compute.TargetPoolsRemoveHealthCheckRequest{
			HealthChecks: make([]*compute.HealthCheckReference, len(remove)),
		}
		for i, v := range remove {
			removeReq.HealthChecks[i] = &compute.HealthCheckReference{HealthCheck: v}
		}
		op, err := config.clientCompute.TargetPools.RemoveHealthCheck(
			project, region, name, removeReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating health_check: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Updating Target Pool", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		addReq := &compute.TargetPoolsAddHealthCheckRequest{
			HealthChecks: make([]*compute.HealthCheckReference, len(add)),
		}
		for i, v := range add {
			addReq.HealthChecks[i] = &compute.HealthCheckReference{HealthCheck: v}
		}
		op, err = config.clientCompute.TargetPools.AddHealthCheck(
			project, region, name, addReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating health_check: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Updating Target Pool", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		d.SetPartial("health_checks")
	}

	if d.HasChange("instances") {

		old_, new_ := d.GetChange("instances")
		old := old_.(*schema.Set)
		new := new_.(*schema.Set)

		addUrls, err := convertInstancesToUrls(d, config, project, new.Difference(old))
		if err != nil {
			return err
		}
		removeUrls, err := convertInstancesToUrls(d, config, project, old.Difference(new))
		if err != nil {
			return err
		}

		addReq := &compute.TargetPoolsAddInstanceRequest{
			Instances: make([]*compute.InstanceReference, len(addUrls)),
		}
		for i, v := range addUrls {
			addReq.Instances[i] = &compute.InstanceReference{Instance: v}
		}
		op, err := config.clientCompute.TargetPools.AddInstance(
			project, region, name, addReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating instances: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Updating Target Pool", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		removeReq := &compute.TargetPoolsRemoveInstanceRequest{
			Instances: make([]*compute.InstanceReference, len(removeUrls)),
		}
		for i, v := range removeUrls {
			removeReq.Instances[i] = &compute.InstanceReference{Instance: v}
		}
		op, err = config.clientCompute.TargetPools.RemoveInstance(
			project, region, name, removeReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating instances: %s", err)
		}
		err = computeOperationWaitTime(config, op, project, "Updating Target Pool", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		d.SetPartial("instances")
	}

	if d.HasChange("backup_pool") {
		bpool_name := d.Get("backup_pool").(string)
		tref := &compute.TargetReference{
			Target: bpool_name,
		}
		op, err := config.clientCompute.TargetPools.SetBackup(
			project, region, name, tref).Do()
		if err != nil {
			return fmt.Errorf("Error updating backup_pool: %s", err)
		}

		err = computeOperationWaitTime(config, op, project, "Updating Target Pool", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		d.SetPartial("backup_pool")
	}

	d.Partial(false)

	return resourceComputeTargetPoolRead(d, meta)
}

func convertInstancesFromUrls(urls []string) []string {
	result := make([]string, 0, len(urls))
	for _, url := range urls {
		urlArray := strings.Split(url, "/")
		instance := fmt.Sprintf("%s/%s", urlArray[len(urlArray)-3], urlArray[len(urlArray)-1])
		result = append(result, instance)
	}
	return result
}

func resourceComputeTargetPoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	tpool, err := config.clientCompute.TargetPools.Get(
		project, region, d.Get("name").(string)).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Target Pool %q", d.Get("name").(string)))
	}

	d.Set("self_link", tpool.SelfLink)
	d.Set("backup_pool", tpool.BackupPool)
	d.Set("description", tpool.Description)
	d.Set("failover_ratio", tpool.FailoverRatio)
	d.Set("health_checks", tpool.HealthChecks)
	if tpool.Instances != nil {
		d.Set("instances", convertInstancesFromUrls(tpool.Instances))
	} else {
		d.Set("instances", nil)
	}
	d.Set("name", tpool.Name)
	d.Set("region", GetResourceNameFromSelfLink(tpool.Region))
	d.Set("session_affinity", tpool.SessionAffinity)
	d.Set("project", project)
	return nil
}

func resourceComputeTargetPoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the TargetPool
	op, err := config.clientCompute.TargetPools.Delete(
		project, region, d.Get("name").(string)).Do()
	if err != nil {
		return fmt.Errorf("Error deleting TargetPool: %s", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting Target Pool", d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceTargetPoolStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/targetPools/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/regions/{{region}}/targetPools/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
