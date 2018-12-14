package google

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func resourceComputeRegionBackendService() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionBackendServiceCreate,
		Read:   resourceComputeRegionBackendServiceRead,
		Update: resourceComputeRegionBackendServiceUpdate,
		Delete: resourceComputeRegionBackendServiceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},

			"health_checks": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
			},

			"backend": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				Set:      resourceGoogleComputeRegionBackendServiceBackendHash,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"session_affinity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"connection_draining_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceComputeRegionBackendServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hc := d.Get("health_checks").(*schema.Set).List()
	healthChecks := make([]string, 0, len(hc))
	for _, v := range hc {
		healthChecks = append(healthChecks, v.(string))
	}

	service := computeBeta.BackendService{
		Name:                d.Get("name").(string),
		HealthChecks:        healthChecks,
		LoadBalancingScheme: "INTERNAL",
	}

	var err error
	if v, ok := d.GetOk("backend"); ok {
		service.Backends, err = expandBackends(v.(*schema.Set).List())
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("description"); ok {
		service.Description = v.(string)
	}

	if v, ok := d.GetOk("protocol"); ok {
		service.Protocol = v.(string)
	}

	if v, ok := d.GetOk("session_affinity"); ok {
		service.SessionAffinity = v.(string)
	}

	if v, ok := d.GetOk("timeout_sec"); ok {
		service.TimeoutSec = int64(v.(int))
	}

	if v, ok := d.GetOk("connection_draining_timeout_sec"); ok {
		connectionDraining := &computeBeta.ConnectionDraining{
			DrainingTimeoutSec: int64(v.(int)),
		}

		service.ConnectionDraining = connectionDraining
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Region Backend Service: %#v", service)

	op, err := config.clientComputeBeta.RegionBackendServices.Insert(
		project, region, &service).Do()
	if err != nil {
		return fmt.Errorf("Error creating backend service: %s", err)
	}

	log.Printf("[DEBUG] Waiting for new backend service, operation: %#v", op)

	d.SetId(service.Name)

	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating Region Backend Service")
	if err != nil {
		return err
	}

	return resourceComputeRegionBackendServiceRead(d, meta)
}

func resourceComputeRegionBackendServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	service, err := config.clientCompute.RegionBackendServices.Get(
		project, region, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Region Backend Service %q", d.Get("name").(string)))
	}

	d.Set("description", service.Description)
	d.Set("protocol", service.Protocol)
	d.Set("session_affinity", service.SessionAffinity)
	d.Set("timeout_sec", service.TimeoutSec)
	d.Set("connection_draining_timeout_sec", service.ConnectionDraining.DrainingTimeoutSec)
	d.Set("fingerprint", service.Fingerprint)
	d.Set("self_link", service.SelfLink)
	err = d.Set("backend", flattenRegionBackends(service.Backends))
	if err != nil {
		return err
	}
	d.Set("health_checks", service.HealthChecks)
	d.Set("project", project)
	d.Set("region", region)

	return nil
}

func resourceComputeRegionBackendServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	hc := d.Get("health_checks").(*schema.Set).List()
	healthChecks := make([]string, 0, len(hc))
	for _, v := range hc {
		healthChecks = append(healthChecks, v.(string))
	}

	service := computeBeta.BackendService{
		Name:                d.Get("name").(string),
		Fingerprint:         d.Get("fingerprint").(string),
		HealthChecks:        healthChecks,
		LoadBalancingScheme: "INTERNAL",
	}

	// Optional things
	if v, ok := d.GetOk("backend"); ok {
		service.Backends, err = expandBackends(v.(*schema.Set).List())
		if err != nil {
			return err
		}
	}
	if v, ok := d.GetOk("description"); ok {
		service.Description = v.(string)
	}
	if v, ok := d.GetOk("protocol"); ok {
		service.Protocol = v.(string)
	}
	if v, ok := d.GetOk("session_affinity"); ok {
		service.SessionAffinity = v.(string)
	}
	if v, ok := d.GetOk("timeout_sec"); ok {
		service.TimeoutSec = int64(v.(int))
	}

	if d.HasChange("connection_draining_timeout_sec") {
		connectionDraining := &computeBeta.ConnectionDraining{
			DrainingTimeoutSec: int64(d.Get("connection_draining_timeout_sec").(int)),
		}

		service.ConnectionDraining = connectionDraining
	}

	log.Printf("[DEBUG] Updating existing Backend Service %q: %#v", d.Id(), service)
	op, err := config.clientComputeBeta.RegionBackendServices.Update(
		project, region, d.Id(), &service).Do()
	if err != nil {
		return fmt.Errorf("Error updating backend service: %s", err)
	}

	d.SetId(service.Name)

	err = computeSharedOperationWait(config.clientCompute, op, project, "Updating Backend Service")
	if err != nil {
		return err
	}

	return resourceComputeRegionBackendServiceRead(d, meta)
}

func resourceComputeRegionBackendServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting backend service %s", d.Id())
	op, err := config.clientCompute.RegionBackendServices.Delete(
		project, region, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting backend service: %s", err)
	}

	err = computeOperationWait(config.clientCompute, op, project, "Deleting Backend Service")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceGoogleComputeRegionBackendServiceBackendHash(v interface{}) int {
	if v == nil {
		return 0
	}

	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if group, err := getRelativePath(m["group"].(string)); err != nil {
		log.Printf("[WARN] Error on retrieving relative path of instance group: %s", err)
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	} else {
		buf.WriteString(fmt.Sprintf("%s-", group))
	}

	if v, ok := m["description"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func flattenRegionBackends(backends []*compute.Backend) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(backends))

	for _, b := range backends {
		data := make(map[string]interface{})

		data["description"] = b.Description
		data["group"] = b.Group
		result = append(result, data)
	}

	return result
}
