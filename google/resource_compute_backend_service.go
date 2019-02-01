package google

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},

			"health_checks": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      selfLinkRelativePathHash,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
			},

			"iap": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oauth2_client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"oauth2_client_secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"oauth2_client_secret_sha256": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"backend": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      resourceGoogleComputeBackendServiceBackendHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: compareSelfLinkRelativePaths,
						},
						"balancing_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "UTILIZATION",
						},
						"capacity_scaler": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  1,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_rate": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_rate_per_instance": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"max_connections": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_connections_per_instance": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_utilization": {
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  0.8,
						},
					},
				},
			},

			"cdn_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_key_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include_host": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"include_protocol": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"include_query_string": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"query_string_blacklist": {
										Type:          schema.TypeSet,
										Optional:      true,
										Elem:          &schema.Schema{Type: schema.TypeString},
										ConflictsWith: []string{"cdn_policy.0.cache_key_policy.query_string_whitelist"},
									},
									"query_string_whitelist": {
										Type:          schema.TypeSet,
										Optional:      true,
										Elem:          &schema.Schema{Type: schema.TypeString},
										ConflictsWith: []string{"cdn_policy.0.cache_key_policy.query_string_blacklist"},
									},
								},
							},
						},
					},
				},
			},

			"custom_request_headers": {
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enable_cdn": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port_name": {
				Type:     schema.TypeString,
				Optional: true,
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

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "region has been removed as it was never used. For internal load balancing, use google_compute_region_backend_service",
			},

			"security_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"session_affinity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"affinity_cookie_ttl_sec": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"connection_draining_timeout_sec": {
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
	op, err := config.clientComputeBeta.BackendServices.Insert(
		project, service).Do()
	if err != nil {
		return fmt.Errorf("Error creating backend service: %s", err)
	}

	log.Printf("[DEBUG] Waiting for new backend service, operation: %#v", op)

	// Store the ID now
	d.SetId(service.Name)

	// Wait for the operation to complete
	waitErr := computeSharedOperationWait(config.clientCompute, op, project, "Creating Backend Service")
	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return waitErr
	}

	if v, ok := d.GetOk("security_policy"); ok {
		pol, err := ParseSecurityPolicyFieldValue(v.(string), d, config)
		op, err := config.clientComputeBeta.BackendServices.SetSecurityPolicy(
			project, service.Name, &computeBeta.SecurityPolicyReference{
				SecurityPolicy: pol.RelativeLink(),
			}).Do()
		if err != nil {
			return errwrap.Wrapf("Error setting Backend Service security policy: {{err}}", err)
		}
		waitErr := computeSharedOperationWait(config.clientCompute, op, project, "Adding Backend Service Security Policy")
		if waitErr != nil {
			return waitErr
		}
	}

	return resourceComputeBackendServiceRead(d, meta)
}

func resourceComputeBackendServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	service, err := config.clientComputeBeta.BackendServices.Get(project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Backend Service %q", d.Get("name").(string)))
	}

	d.Set("name", service.Name)
	d.Set("description", service.Description)
	d.Set("enable_cdn", service.EnableCDN)
	d.Set("port_name", service.PortName)
	d.Set("protocol", service.Protocol)
	d.Set("session_affinity", service.SessionAffinity)
	d.Set("affinity_cookie_ttl_sec", service.AffinityCookieTtlSec)
	d.Set("timeout_sec", service.TimeoutSec)
	d.Set("fingerprint", service.Fingerprint)
	d.Set("self_link", ConvertSelfLinkToV1(service.SelfLink))
	d.Set("backend", flattenBackends(service.Backends))
	d.Set("connection_draining_timeout_sec", service.ConnectionDraining.DrainingTimeoutSec)
	d.Set("iap", flattenIap(d, service.Iap))
	d.Set("project", project)
	guardedHealthChecks := make([]string, len(service.HealthChecks))
	for i, v := range service.HealthChecks {
		guardedHealthChecks[i] = ConvertSelfLinkToV1(v)
	}

	d.Set("health_checks", guardedHealthChecks)
	if err := d.Set("cdn_policy", flattenCdnPolicy(service.CdnPolicy)); err != nil {
		return err
	}
	d.Set("security_policy", service.SecurityPolicy)

	d.Set("custom_request_headers", nil)
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
	op, err := config.clientComputeBeta.BackendServices.Update(
		project, d.Id(), service).Do()
	if err != nil {
		return fmt.Errorf("Error updating backend service: %s", err)
	}

	err = computeSharedOperationWait(config.clientCompute, op, project, "Updating Backend Service")
	if err != nil {
		return err
	}

	if d.HasChange("security_policy") {
		pol, err := ParseSecurityPolicyFieldValue(d.Get("security_policy").(string), d, config)
		if err != nil {
			return err
		}
		op, err := config.clientComputeBeta.BackendServices.SetSecurityPolicy(
			project, service.Name, &computeBeta.SecurityPolicyReference{
				SecurityPolicy: pol.RelativeLink(),
			}).Do()
		if err != nil {
			return err
		}
		waitErr := computeSharedOperationWait(config.clientCompute, op, project, "Adding Backend Service Security Policy")
		if waitErr != nil {
			return waitErr
		}
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

func expandIap(configured []interface{}) *computeBeta.BackendServiceIAP {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &computeBeta.BackendServiceIAP{
		Enabled:            true,
		Oauth2ClientId:     data["oauth2_client_id"].(string),
		Oauth2ClientSecret: data["oauth2_client_secret"].(string),
		ForceSendFields:    []string{"Enabled", "Oauth2ClientId", "Oauth2ClientSecret"},
	}
}

func flattenIap(d *schema.ResourceData, iap *computeBeta.BackendServiceIAP) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if iap == nil || !iap.Enabled {
		return result
	}

	return append(result, map[string]interface{}{
		"oauth2_client_id":            iap.Oauth2ClientId,
		"oauth2_client_secret":        d.Get("iap.0.oauth2_client_secret"),
		"oauth2_client_secret_sha256": iap.Oauth2ClientSecretSha256,
	})
}

func expandBackends(configured []interface{}) ([]*computeBeta.Backend, error) {
	backends := make([]*computeBeta.Backend, 0, len(configured))

	for _, raw := range configured {
		data := raw.(map[string]interface{})

		g, ok := data["group"]
		if !ok {
			return nil, errors.New("google_compute_backend_service.backend.group must be set")
		}

		b := computeBeta.Backend{
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
		if v, ok := data["max_connections"]; ok {
			b.MaxConnections = int64(v.(int))
			if b.MaxConnections == 0 {
				b.NullFields = append(b.NullFields, "MaxConnections")
			}
		}
		if v, ok := data["max_connections_per_instance"]; ok {
			b.MaxConnectionsPerInstance = int64(v.(int))
			if b.MaxConnectionsPerInstance == 0 {
				b.NullFields = append(b.NullFields, "MaxConnectionsPerInstance")
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

func flattenBackends(backends []*computeBeta.Backend) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(backends))

	for _, b := range backends {
		data := make(map[string]interface{})

		data["balancing_mode"] = b.BalancingMode
		data["capacity_scaler"] = b.CapacityScaler
		data["description"] = b.Description
		data["group"] = b.Group
		data["max_rate"] = b.MaxRate
		data["max_rate_per_instance"] = b.MaxRatePerInstance
		data["max_connections"] = b.MaxConnections
		data["max_connections_per_instance"] = b.MaxConnectionsPerInstance
		data["max_utilization"] = b.MaxUtilization
		result = append(result, data)
	}

	return result
}

func expandBackendService(d *schema.ResourceData) (*computeBeta.BackendService, error) {
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
	service := &computeBeta.BackendService{
		Name:         d.Get("name").(string),
		HealthChecks: healthChecks,
		Iap: &computeBeta.BackendServiceIAP{
			ForceSendFields: []string{"Enabled", "Oauth2ClientId", "Oauth2ClientSecret"},
		},
		CdnPolicy: &computeBeta.BackendServiceCdnPolicy{
			CacheKeyPolicy: &computeBeta.CacheKeyPolicy{
				ForceSendFields: []string{"IncludeProtocol", "IncludeHost", "IncludeQueryString", "QueryStringWhitelist", "QueryStringBlacklist"},
			},
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

	if v, ok := d.GetOk("affinity_cookie_ttl_sec"); ok {
		service.AffinityCookieTtlSec = int64(v.(int))
	}

	if v, ok := d.GetOk("timeout_sec"); ok {
		service.TimeoutSec = int64(v.(int))
	}

	if v, ok := d.GetOk("enable_cdn"); ok {
		service.EnableCDN = v.(bool)
	}

	connectionDrainingTimeoutSec := d.Get("connection_draining_timeout_sec")
	connectionDraining := &computeBeta.ConnectionDraining{
		DrainingTimeoutSec: int64(connectionDrainingTimeoutSec.(int)),
	}

	service.ConnectionDraining = connectionDraining

	if v, ok := d.GetOk("cdn_policy"); ok {
		c := expandCdnPolicy(v.([]interface{}))
		if c != nil {
			service.CdnPolicy = c
		}
	}

	return service, nil
}

func expandCdnPolicy(configured []interface{}) *computeBeta.BackendServiceCdnPolicy {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	ckp := data["cache_key_policy"].([]interface{})
	if len(ckp) == 0 {
		return nil
	}
	ckpData := ckp[0].(map[string]interface{})

	return &computeBeta.BackendServiceCdnPolicy{
		CacheKeyPolicy: &computeBeta.CacheKeyPolicy{
			IncludeHost:          ckpData["include_host"].(bool),
			IncludeProtocol:      ckpData["include_protocol"].(bool),
			IncludeQueryString:   ckpData["include_query_string"].(bool),
			QueryStringBlacklist: convertStringSet(ckpData["query_string_blacklist"].(*schema.Set)),
			QueryStringWhitelist: convertStringSet(ckpData["query_string_whitelist"].(*schema.Set)),
			ForceSendFields:      []string{"IncludeProtocol", "IncludeHost", "IncludeQueryString", "QueryStringWhitelist", "QueryStringBlacklist"},
		},
	}
}

func flattenCdnPolicy(pol *computeBeta.BackendServiceCdnPolicy) []map[string]interface{} {
	result := []map[string]interface{}{}
	if pol == nil || pol.CacheKeyPolicy == nil {
		return result
	}

	return append(result, map[string]interface{}{
		"cache_key_policy": []map[string]interface{}{
			{
				"include_host":           pol.CacheKeyPolicy.IncludeHost,
				"include_protocol":       pol.CacheKeyPolicy.IncludeProtocol,
				"include_query_string":   pol.CacheKeyPolicy.IncludeQueryString,
				"query_string_blacklist": schema.NewSet(schema.HashString, convertStringArrToInterface(pol.CacheKeyPolicy.QueryStringBlacklist)),
				"query_string_whitelist": schema.NewSet(schema.HashString, convertStringArrToInterface(pol.CacheKeyPolicy.QueryStringWhitelist)),
			},
		},
	})
}
