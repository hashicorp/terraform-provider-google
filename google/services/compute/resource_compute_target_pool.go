// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"

	"google.golang.org/api/compute/v1"
)

var instancesSelfLinkPattern = regexp.MustCompile(fmt.Sprintf(tpgresource.ZonalLinkBasePattern, "instances"))

func ResourceComputeTargetPool() *schema.Resource {
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `A unique name for the resource, required by GCE. Changing this forces a new resource to be created.`,
			},

			"backup_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: `URL to the backup target pool. Must also set failover_ratio.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Textual description field.`,
			},

			"failover_ratio": {
				Type:        schema.TypeFloat,
				Optional:    true,
				ForceNew:    true,
				Description: `Ratio (0 to 1) of failed nodes before using the backup pool (which must also be set).`,
			},

			"health_checks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				},
				Description: `List of zero or one health check name or self_link. Only legacy google_compute_http_health_check is supported.`,
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
				Description: `List of instances in the pool. They can be given as URLs, or in the form of "zone/name". Note that the instances need not exist at the time of target pool creation, so there is no need to use the Terraform interpolators to create a dependency on the instances from the target pool.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Where the target pool resides. Defaults to project region.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},

			"session_affinity": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "NONE",
				Description: `How to distribute load. Options are "NONE" (no affinity). "CLIENT_IP" (hash of the source/dest addresses / ports), and "CLIENT_IP_PROTO" also includes the protocol (default "NONE").`,
			},
		},
		UseJSONNumber: true,
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
func convertHealthChecks(healthChecks []interface{}, d *schema.ResourceData, config *transport_tpg.Config) ([]string, error) {
	if len(healthChecks) == 0 {
		return []string{}, nil
	}

	hc, err := tpgresource.ParseHttpHealthCheckFieldValue(healthChecks[0].(string), d, config)
	if err != nil {
		return nil, err
	}

	return []string{hc.RelativeLink()}, nil
}

// Instances do not need to exist yet, so we simply generate URLs.
// Instances can be full URLS or zone/name
func convertInstancesToUrls(d *schema.ResourceData, config *transport_tpg.Config, project string, names *schema.Set) ([]string, error) {
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
				url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf(
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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
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
	op, err := config.NewComputeClient(userAgent).TargetPools.Insert(
		project, region, tpool).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 && strings.Contains(gerr.Message, "httpHealthChecks") {
			return fmt.Errorf("Health check %s is not a valid HTTP health check", d.Get("health_checks").([]interface{})[0])
		}
		return fmt.Errorf("Error creating TargetPool: %s", err)
	}

	// It probably maybe worked, so store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/targetPools/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = ComputeOperationWaitTime(config, op, project, "Creating Target Pool", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	return resourceComputeTargetPoolRead(d, meta)
}

func resourceComputeTargetPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
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
		add, remove := tpgresource.CalcAddRemove(fromUrls, toUrls)

		removeReq := &compute.TargetPoolsRemoveHealthCheckRequest{
			HealthChecks: make([]*compute.HealthCheckReference, len(remove)),
		}
		for i, v := range remove {
			removeReq.HealthChecks[i] = &compute.HealthCheckReference{HealthCheck: v}
		}
		op, err := config.NewComputeClient(userAgent).TargetPools.RemoveHealthCheck(
			project, region, name, removeReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating health_check: %s", err)
		}

		err = ComputeOperationWaitTime(config, op, project, "Updating Target Pool", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		addReq := &compute.TargetPoolsAddHealthCheckRequest{
			HealthChecks: make([]*compute.HealthCheckReference, len(add)),
		}
		for i, v := range add {
			addReq.HealthChecks[i] = &compute.HealthCheckReference{HealthCheck: v}
		}
		op, err = config.NewComputeClient(userAgent).TargetPools.AddHealthCheck(
			project, region, name, addReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating health_check: %s", err)
		}

		err = ComputeOperationWaitTime(config, op, project, "Updating Target Pool", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
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
		op, err := config.NewComputeClient(userAgent).TargetPools.AddInstance(
			project, region, name, addReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating instances: %s", err)
		}

		err = ComputeOperationWaitTime(config, op, project, "Updating Target Pool", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		removeReq := &compute.TargetPoolsRemoveInstanceRequest{
			Instances: make([]*compute.InstanceReference, len(removeUrls)),
		}
		for i, v := range removeUrls {
			removeReq.Instances[i] = &compute.InstanceReference{Instance: v}
		}
		op, err = config.NewComputeClient(userAgent).TargetPools.RemoveInstance(
			project, region, name, removeReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating instances: %s", err)
		}
		err = ComputeOperationWaitTime(config, op, project, "Updating Target Pool", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	if d.HasChange("backup_pool") {
		bpool_name := d.Get("backup_pool").(string)
		tref := &compute.TargetReference{
			Target: bpool_name,
		}
		op, err := config.NewComputeClient(userAgent).TargetPools.SetBackup(
			project, region, name, tref).Do()
		if err != nil {
			return fmt.Errorf("Error updating backup_pool: %s", err)
		}

		err = ComputeOperationWaitTime(config, op, project, "Updating Target Pool", userAgent, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	tpool, err := config.NewComputeClient(userAgent).TargetPools.Get(
		project, region, d.Get("name").(string)).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Target Pool %q", d.Get("name").(string)))
	}

	if err := d.Set("self_link", tpool.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("backup_pool", tpool.BackupPool); err != nil {
		return fmt.Errorf("Error setting backup_pool: %s", err)
	}
	if err := d.Set("description", tpool.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("failover_ratio", tpool.FailoverRatio); err != nil {
		return fmt.Errorf("Error setting failover_ratio: %s", err)
	}
	if err := d.Set("health_checks", tpool.HealthChecks); err != nil {
		return fmt.Errorf("Error setting health_checks: %s", err)
	}
	if tpool.Instances != nil {
		if err := d.Set("instances", convertInstancesFromUrls(tpool.Instances)); err != nil {
			return fmt.Errorf("Error setting instances: %s", err)
		}
	} else {
		if err := d.Set("instances", nil); err != nil {
			return fmt.Errorf("Error setting instances: %s", err)
		}
	}
	if err := d.Set("name", tpool.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("region", tpgresource.GetResourceNameFromSelfLink(tpool.Region)); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("session_affinity", tpool.SessionAffinity); err != nil {
		return fmt.Errorf("Error setting session_affinity: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	return nil
}

func resourceComputeTargetPoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	// Delete the TargetPool
	op, err := config.NewComputeClient(userAgent).TargetPools.Delete(
		project, region, d.Get("name").(string)).Do()
	if err != nil {
		return fmt.Errorf("Error deleting TargetPool: %s", err)
	}

	err = ComputeOperationWaitTime(config, op, project, "Deleting Target Pool", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceTargetPoolStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/targetPools/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/regions/{{region}}/targetPools/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
