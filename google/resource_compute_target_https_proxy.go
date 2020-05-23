// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceComputeTargetHttpsProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeTargetHttpsProxyCreate,
		Read:   resourceComputeTargetHttpsProxyRead,
		Update: resourceComputeTargetHttpsProxyUpdate,
		Delete: resourceComputeTargetHttpsProxyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeTargetHttpsProxyImport,
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
				Description: `Name of the resource. Provided by the client when the resource is
created. The name must be 1-63 characters long, and comply with
RFC1035. Specifically, the name must be 1-63 characters long and match
the regular expression '[a-z]([-a-z0-9]*[a-z0-9])?' which means the
first character must be a lowercase letter, and all following
characters must be a dash, lowercase letter, or digit, except the last
character, which cannot be a dash.`,
			},
			"ssl_certificates": {
				Type:     schema.TypeList,
				Required: true,
				Description: `A list of SslCertificate resources that are used to authenticate
connections between users and the load balancer. At least one SSL
certificate must be specified.`,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
			},
			"url_map": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description: `A reference to the UrlMap resource that defines the mapping from URL
to the BackendService.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `An optional description of this resource.`,
			},
			"quic_override": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "ENABLE", "DISABLE", ""}, false),
				Description: `Specifies the QUIC override policy for this resource. This determines
whether the load balancer will attempt to negotiate QUIC with clients
or not. Can specify one of NONE, ENABLE, or DISABLE. If NONE is
specified, uses the QUIC policy with no user overrides, which is
equivalent to DISABLE. Default value: "NONE" Possible values: ["NONE", "ENABLE", "DISABLE"]`,
				Default: "NONE",
			},
			"ssl_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description: `A reference to the SslPolicy resource that will be associated with
the TargetHttpsProxy resource. If not set, the TargetHttpsProxy
resource will not have any SSL policy configured.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Creation timestamp in RFC3339 text format.`,
			},
			"proxy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The unique identifier for the resource.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeTargetHttpsProxyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	descriptionProp, err := expandComputeTargetHttpsProxyDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	nameProp, err := expandComputeTargetHttpsProxyName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	quicOverrideProp, err := expandComputeTargetHttpsProxyQuicOverride(d.Get("quic_override"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("quic_override"); !isEmptyValue(reflect.ValueOf(quicOverrideProp)) && (ok || !reflect.DeepEqual(v, quicOverrideProp)) {
		obj["quicOverride"] = quicOverrideProp
	}
	sslCertificatesProp, err := expandComputeTargetHttpsProxySslCertificates(d.Get("ssl_certificates"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ssl_certificates"); !isEmptyValue(reflect.ValueOf(sslCertificatesProp)) && (ok || !reflect.DeepEqual(v, sslCertificatesProp)) {
		obj["sslCertificates"] = sslCertificatesProp
	}
	sslPolicyProp, err := expandComputeTargetHttpsProxySslPolicy(d.Get("ssl_policy"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ssl_policy"); !isEmptyValue(reflect.ValueOf(sslPolicyProp)) && (ok || !reflect.DeepEqual(v, sslPolicyProp)) {
		obj["sslPolicy"] = sslPolicyProp
	}
	urlMapProp, err := expandComputeTargetHttpsProxyUrlMap(d.Get("url_map"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("url_map"); !isEmptyValue(reflect.ValueOf(urlMapProp)) && (ok || !reflect.DeepEqual(v, urlMapProp)) {
		obj["urlMap"] = urlMapProp
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/targetHttpsProxies")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new TargetHttpsProxy: %#v", obj)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating TargetHttpsProxy: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/global/targetHttpsProxies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(
		config, res, project, "Creating TargetHttpsProxy",
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create TargetHttpsProxy: %s", err)
	}

	log.Printf("[DEBUG] Finished creating TargetHttpsProxy %q: %#v", d.Id(), res)

	return resourceComputeTargetHttpsProxyRead(d, meta)
}

func resourceComputeTargetHttpsProxyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/targetHttpsProxies/{{name}}")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeTargetHttpsProxy %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}

	if err := d.Set("creation_timestamp", flattenComputeTargetHttpsProxyCreationTimestamp(res["creationTimestamp"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("description", flattenComputeTargetHttpsProxyDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("proxy_id", flattenComputeTargetHttpsProxyProxyId(res["id"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("name", flattenComputeTargetHttpsProxyName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("quic_override", flattenComputeTargetHttpsProxyQuicOverride(res["quicOverride"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("ssl_certificates", flattenComputeTargetHttpsProxySslCertificates(res["sslCertificates"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("ssl_policy", flattenComputeTargetHttpsProxySslPolicy(res["sslPolicy"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("url_map", flattenComputeTargetHttpsProxyUrlMap(res["urlMap"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading TargetHttpsProxy: %s", err)
	}

	return nil
}

func resourceComputeTargetHttpsProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("ssl_certificates") {
		obj := make(map[string]interface{})

		sslCertificatesProp, err := expandComputeTargetHttpsProxySslCertificates(d.Get("ssl_certificates"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("ssl_certificates"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, sslCertificatesProp)) {
			obj["sslCertificates"] = sslCertificatesProp
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/targetHttpsProxies/{{name}}/setSslCertificates")
		if err != nil {
			return err
		}
		res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating TargetHttpsProxy %q: %s", d.Id(), err)
		}

		err = computeOperationWaitTime(
			config, res, project, "Updating TargetHttpsProxy",
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		d.SetPartial("ssl_certificates")
	}
	if d.HasChange("ssl_policy") {
		obj := make(map[string]interface{})

		sslPolicyProp, err := expandComputeTargetHttpsProxySslPolicy(d.Get("ssl_policy"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("ssl_policy"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, sslPolicyProp)) {
			obj["sslPolicy"] = sslPolicyProp
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/targetHttpsProxies/{{name}}/setSslPolicy")
		if err != nil {
			return err
		}
		res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating TargetHttpsProxy %q: %s", d.Id(), err)
		}

		err = computeOperationWaitTime(
			config, res, project, "Updating TargetHttpsProxy",
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		d.SetPartial("ssl_policy")
	}
	if d.HasChange("url_map") {
		obj := make(map[string]interface{})

		urlMapProp, err := expandComputeTargetHttpsProxyUrlMap(d.Get("url_map"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("url_map"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, urlMapProp)) {
			obj["urlMap"] = urlMapProp
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/targetHttpsProxies/{{name}}/setUrlMap")
		if err != nil {
			return err
		}
		res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating TargetHttpsProxy %q: %s", d.Id(), err)
		}

		err = computeOperationWaitTime(
			config, res, project, "Updating TargetHttpsProxy",
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		d.SetPartial("url_map")
	}
	if d.HasChange("quic_override") {
		obj := make(map[string]interface{})

		quicOverrideProp, err := expandComputeTargetHttpsProxyQuicOverride(d.Get("quic_override"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("quic_override"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, quicOverrideProp)) {
			obj["quicOverride"] = quicOverrideProp
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/targetHttpsProxies/{{name}}/setQuicOverride")
		if err != nil {
			return err
		}
		res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating TargetHttpsProxy %q: %s", d.Id(), err)
		}

		err = computeOperationWaitTime(
			config, res, project, "Updating TargetHttpsProxy",
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}

		d.SetPartial("quic_override")
	}

	d.Partial(false)

	return resourceComputeTargetHttpsProxyRead(d, meta)
}

func resourceComputeTargetHttpsProxyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/targetHttpsProxies/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting TargetHttpsProxy %q", d.Id())

	res, err := sendRequestWithTimeout(config, "DELETE", project, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "TargetHttpsProxy")
	}

	err = computeOperationWaitTime(
		config, res, project, "Deleting TargetHttpsProxy",
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting TargetHttpsProxy %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeTargetHttpsProxyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/global/targetHttpsProxies/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/global/targetHttpsProxies/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeTargetHttpsProxyCreationTimestamp(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeTargetHttpsProxyDescription(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeTargetHttpsProxyProxyId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeTargetHttpsProxyName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenComputeTargetHttpsProxyQuicOverride(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil || isEmptyValue(reflect.ValueOf(v)) {
		return "NONE"
	}

	return v
}

func flattenComputeTargetHttpsProxySslCertificates(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return convertAndMapStringArr(v.([]interface{}), ConvertSelfLinkToV1)
}

func flattenComputeTargetHttpsProxySslPolicy(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeTargetHttpsProxyUrlMap(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func expandComputeTargetHttpsProxyDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetHttpsProxyName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetHttpsProxyQuicOverride(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeTargetHttpsProxySslCertificates(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			return nil, fmt.Errorf("Invalid value for ssl_certificates: nil")
		}
		f, err := parseGlobalFieldValue("sslCertificates", raw.(string), "project", d, config, true)
		if err != nil {
			return nil, fmt.Errorf("Invalid value for ssl_certificates: %s", err)
		}
		req = append(req, f.RelativeLink())
	}
	return req, nil
}

func expandComputeTargetHttpsProxySslPolicy(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("sslPolicies", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for ssl_policy: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeTargetHttpsProxyUrlMap(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("urlMaps", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for url_map: %s", err)
	}
	return f.RelativeLink(), nil
}
