package google

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeBackendService() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeBackendServiceCreate,
		Read:   resourceComputeBackendServiceRead,
		Update: resourceComputeBackendServiceUpdate,
		Delete: resourceComputeBackendServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},

			"health_checks": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
			},

			"iap": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oauth2_client_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"oauth2_client_secret": &schema.Schema{
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == fmt.Sprintf("%x", sha256.Sum256([]byte(new))) {
									return true
								}
								return false
							},
						},
					},
				},
			},

			"backend": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Set:      resourceGoogleComputeBackendServiceBackendHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},
						"balancing_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "UTILIZATION",
						},
						"capacity_scaler": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  1,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_rate": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_rate_per_instance": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"max_utilization": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0.8,
						},
					},
				},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"enable_cdn": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"port_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "region has been removed as it was never used. For internal load balancing, use google_compute_region_backend_service",
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"session_affinity": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"timeout_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"connection_draining_timeout_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
		},
	}
}

func resourceComputeBackendServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service, err := expandBackendService(d)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Backend Service: %#v", service)
	op, err := config.clientCompute.BackendServices.Insert(
		project, service).Do()
	if err != nil {
		return fmt.Errorf("Error creating backend service: %s", err)
	}

	log.Printf("[DEBUG] Waiting for new backend service, operation: %#v", op)

	// Store the ID now
	d.SetId(service.Name)

	// Wait for the operation to complete
	waitErr := computeOperationWait(config.clientCompute, op, project, "Creating Backend Service")
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	return resourceComputeBackendServiceRead(d, meta)
}

func resourceComputeBackendServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	service, err := config.clientCompute.BackendServices.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Backend Service %q", d.Get("name").(string)))
	}

	d.Set("name", service.Name)
	d.Set("description", service.Description)
	d.Set("enable_cdn", service.EnableCDN)
	d.Set("port_name", service.PortName)
	d.Set("protocol", service.Protocol)
	d.Set("session_affinity", service.SessionAffinity)
	d.Set("timeout_sec", service.TimeoutSec)
	d.Set("fingerprint", service.Fingerprint)
	d.Set("self_link", service.SelfLink)
	d.Set("backend", flattenBackends(service.Backends))
	d.Set("connection_draining_timeout_sec", service.ConnectionDraining.DrainingTimeoutSec)
	d.Set("iap", flattenIap(service.Iap))
	d.Set("project", project)
	d.Set("health_checks", service.HealthChecks)

	return nil
}

func resourceComputeBackendServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service, err := expandBackendService(d)
	if err != nil {
		return err
	}
	service.Fingerprint = d.Get("fingerprint").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating existing Backend Service %q: %#v", d.Id(), service)
	op, err := config.clientCompute.BackendServices.Update(
		project, d.Id(), service).Do()
	if err != nil {
		return fmt.Errorf("Error updating backend service: %s", err)
	}

	d.SetId(service.Name)

	err = computeOperationWait(config.clientCompute, op, project, "Updating Backend Service")
	if err != nil {
		return err
	}

	return resourceComputeBackendServiceRead(d, meta)
}

func resourceComputeBackendServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting backend service %s", d.Id())
	op, err := config.clientCompute.BackendServices.Delete(
		project, d.Id()).Do()
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

func expandIap(configured []interface{}) *compute.BackendServiceIAP {
	data := configured[0].(map[string]interface{})
	iap := &compute.BackendServiceIAP{
		Enabled:            true,
		Oauth2ClientId:     data["oauth2_client_id"].(string),
		Oauth2ClientSecret: data["oauth2_client_secret"].(string),
		ForceSendFields:    []string{"Enabled", "Oauth2ClientId", "Oauth2ClientSecret"},
	}

	return iap
}

func flattenIap(iap *compute.BackendServiceIAP) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if iap == nil || !iap.Enabled {
		return result
	}

	result = append(result, map[string]interface{}{
		"oauth2_client_id":     iap.Oauth2ClientId,
		"oauth2_client_secret": iap.Oauth2ClientSecretSha256,
	})

	return result
}

func expandBackends(configured []interface{}) ([]*compute.Backend, error) {
	backends := make([]*compute.Backend, 0, len(configured))

	for _, raw := range configured {
		data := raw.(map[string]interface{})

		g, ok := data["group"]
		if !ok {
			return nil, errors.New("google_compute_backend_service.backend.group must be set")
		}

		b := compute.Backend{
			Group: g.(string),
		}

		if v, ok := data["balancing_mode"]; ok {
			b.BalancingMode = v.(string)
		}
		if v, ok := data["capacity_scaler"]; ok {
			b.CapacityScaler = v.(float64)
			b.ForceSendFields = append(b.ForceSendFields, "CapacityScaler")
		}
		if v, ok := data["description"]; ok {
			b.Description = v.(string)
		}
		if v, ok := data["max_rate"]; ok {
			b.MaxRate = int64(v.(int))
			if b.MaxRate == 0 {
				b.NullFields = append(b.NullFields, "MaxRate")
			}
		}
		if v, ok := data["max_rate_per_instance"]; ok {
			b.MaxRatePerInstance = v.(float64)
			if b.MaxRatePerInstance == 0 {
				b.NullFields = append(b.NullFields, "MaxRatePerInstance")
			}
		}
		if v, ok := data["max_utilization"]; ok {
			b.MaxUtilization = v.(float64)
			b.ForceSendFields = append(b.ForceSendFields, "MaxUtilization")
		}

		backends = append(backends, &b)
	}

	return backends, nil
}

func flattenBackends(backends []*compute.Backend) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(backends))

	for _, b := range backends {
		data := make(map[string]interface{})

		data["balancing_mode"] = b.BalancingMode
		data["capacity_scaler"] = b.CapacityScaler
		data["description"] = b.Description
		data["group"] = b.Group
		data["max_rate"] = b.MaxRate
		data["max_rate_per_instance"] = b.MaxRatePerInstance
		data["max_utilization"] = b.MaxUtilization
		result = append(result, data)
	}

	return result
}

func expandBackendService(d *schema.ResourceData) (*compute.BackendService, error) {
	hc := d.Get("health_checks").(*schema.Set).List()
	healthChecks := make([]string, 0, len(hc))
	for _, v := range hc {
		healthChecks = append(healthChecks, v.(string))
	}

	// The IAP service is enabled and disabled by adding or removing
	// the IAP configuration block (and providing the client id
	// and secret). We are force sending the three required API fields
	// to enable/disable IAP at all times here, and relying on Golang's
	// type defaults to enable or disable IAP in the existence or absence
	// of the block, instead of checking if the block exists, zeroing out
	// fields, etc.
	service := &compute.BackendService{
		Name:         d.Get("name").(string),
		HealthChecks: healthChecks,
		Iap: &compute.BackendServiceIAP{
			ForceSendFields: []string{"Enabled", "Oauth2ClientId", "Oauth2ClientSecret"},
		},
	}

	if v, ok := d.GetOk("iap"); ok {
		service.Iap = expandIap(v.([]interface{}))
	}

	var err error
	if v, ok := d.GetOk("backend"); ok {
		service.Backends, err = expandBackends(v.(*schema.Set).List())
		if err != nil {
			return nil, err
		}
	}

	if v, ok := d.GetOk("description"); ok {
		service.Description = v.(string)
	}

	if v, ok := d.GetOk("port_name"); ok {
		service.PortName = v.(string)
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

	if v, ok := d.GetOk("enable_cdn"); ok {
		service.EnableCDN = v.(bool)
	}

	connectionDrainingTimeoutSec := d.Get("connection_draining_timeout_sec")
	connectionDraining := &compute.ConnectionDraining{
		DrainingTimeoutSec: int64(connectionDrainingTimeoutSec.(int)),
	}

	service.ConnectionDraining = connectionDraining

	return service, nil
}

func resourceGoogleComputeBackendServiceBackendHash(v interface{}) int {
	if v == nil {
		return 0
	}

	var buf bytes.Buffer
	m := v.(map[string]interface{})

	group, _ := getRelativePath(m["group"].(string))
	buf.WriteString(fmt.Sprintf("%s-", group))

	if v, ok := m["balancing_mode"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["capacity_scaler"]; ok {
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["description"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["max_rate"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", int64(v.(int))))
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}

	return hashcode.String(buf.String())
}
