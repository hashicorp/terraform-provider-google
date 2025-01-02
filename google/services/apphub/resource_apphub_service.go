// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
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

package apphub

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceApphubService() *schema.Resource {
	return &schema.Resource{
		Create: resourceApphubServiceCreate,
		Read:   resourceApphubServiceRead,
		Update: resourceApphubServiceUpdate,
		Delete: resourceApphubServiceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApphubServiceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Part of 'parent'. ApplicationId of parent AppHub Application. Example: projects/{HOST_PROJECT_ID}/locations/{LOCATION}/applications/{APPLICATION_ID}`,
			},
			"discovered_service": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.ProjectNumberDiffSuppress,
				Description:      `Immutable. The resource name of the original discovered service.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Part of 'parent'. Location of parent AppHub Application. Example: projects/{HOST_PROJECT_ID}/locations/{LOCATION}/applications/{APPLICATION_ID}`,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Description: `Part of 'parent'. ProjectId of parent AppHub Application. Example: projects/{HOST_PROJECT_ID}/locations/{LOCATION}/applications/{APPLICATION_ID}`

			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The Service identifier.`,
			},
			"attributes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Consumer provided attributes.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"business_owners": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Business team that ensures user needs are met and value is delivered`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Required. Email address of the contacts.`,
									},
									"display_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Contact's name.`,
									},
								},
							},
						},
						"criticality": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Criticality of the Application, Service, or Workload`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateEnum([]string{"MISSION_CRITICAL", "HIGH", "MEDIUM", "LOW"}),
										Description:  `Criticality type. Possible values: ["MISSION_CRITICAL", "HIGH", "MEDIUM", "LOW"]`,
									},
								},
							},
						},
						"developer_owners": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Developer team that owns development and coding.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Required. Email address of the contacts.`,
									},
									"display_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Contact's name.`,
									},
								},
							},
						},
						"environment": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Environment of the Application, Service, or Workload`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidateEnum([]string{"PRODUCTION", "STAGING", "TEST", "DEVELOPMENT"}),
										Description:  `Environment type. Possible values: ["PRODUCTION", "STAGING", "TEST", "DEVELOPMENT"]`,
									},
								},
							},
						},
						"operator_owners": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Operator team that ensures runtime and operations.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Required. Email address of the contacts.`,
									},
									"display_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Contact's name.`,
									},
								},
							},
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `User-defined description of a Service.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `User-defined name for the Service.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. Create time.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Identifier. The resource name of a Service. Format:
"projects/{host-project-id}/locations/{location}/applications/{application-id}/services/{service-id}"`,
			},
			"service_properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Properties of an underlying cloud resource that can comprise a Service.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gcp_project": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Output only. The service project identifier that the underlying cloud resource resides in.`,
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Output only. The location that the underlying resource resides in, for example, us-west1.`,
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Output only. The location that the underlying resource resides in if it is zonal, for example, us-west1-a).`,
						},
					},
				},
			},
			"service_reference": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Reference to an underlying networking resource that can comprise a Service.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `Output only. The underlying resource URI (For example, URI of Forwarding Rule, URL Map,
and Backend Service).`,
						},
					},
				},
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. Service state. Possible values: STATE_UNSPECIFIED CREATING ACTIVE DELETING DETACHED`,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. A universally unique identifier (UUID) for the 'Service' in the UUID4
format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. Update time.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceApphubServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandApphubServiceDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandApphubServiceDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	attributesProp, err := expandApphubServiceAttributes(d.Get("attributes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("attributes"); !tpgresource.IsEmptyValue(reflect.ValueOf(attributesProp)) && (ok || !reflect.DeepEqual(v, attributesProp)) {
		obj["attributes"] = attributesProp
	}
	discoveredServiceProp, err := expandApphubServiceDiscoveredService(d.Get("discovered_service"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("discovered_service"); !tpgresource.IsEmptyValue(reflect.ValueOf(discoveredServiceProp)) && (ok || !reflect.DeepEqual(v, discoveredServiceProp)) {
		obj["discoveredService"] = discoveredServiceProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApphubBasePath}}projects/{{project}}/locations/{{location}}/applications/{{application_id}}/services?serviceId={{service_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Service: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Service: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating Service: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/applications/{{application_id}}/services/{{service_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = ApphubOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Creating Service", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")

		return fmt.Errorf("Error waiting to create Service: %s", err)
	}

	if err := d.Set("name", flattenApphubServiceName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/applications/{{application_id}}/services/{{service_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Service %q: %#v", d.Id(), res)

	return resourceApphubServiceRead(d, meta)
}

func resourceApphubServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApphubBasePath}}projects/{{project}}/locations/{{location}}/applications/{{application_id}}/services/{{service_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Service: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ApphubService %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}

	if err := d.Set("name", flattenApphubServiceName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("display_name", flattenApphubServiceDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("description", flattenApphubServiceDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("service_reference", flattenApphubServiceServiceReference(res["serviceReference"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("service_properties", flattenApphubServiceServiceProperties(res["serviceProperties"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("attributes", flattenApphubServiceAttributes(res["attributes"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("discovered_service", flattenApphubServiceDiscoveredService(res["discoveredService"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("create_time", flattenApphubServiceCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("update_time", flattenApphubServiceUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("uid", flattenApphubServiceUid(res["uid"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}
	if err := d.Set("state", flattenApphubServiceState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading Service: %s", err)
	}

	return nil
}

func resourceApphubServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Service: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	displayNameProp, err := expandApphubServiceDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandApphubServiceDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	attributesProp, err := expandApphubServiceAttributes(d.Get("attributes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("attributes"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, attributesProp)) {
		obj["attributes"] = attributesProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ApphubBasePath}}projects/{{project}}/locations/{{location}}/applications/{{application_id}}/services/{{service_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Service %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("attributes") {
		updateMask = append(updateMask, "attributes")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// if updateMask is empty we are not updating anything so skip the post
	if len(updateMask) > 0 {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
			Headers:   headers,
		})

		if err != nil {
			return fmt.Errorf("Error updating Service %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Service %q: %#v", d.Id(), res)
		}

		err = ApphubOperationWaitTime(
			config, res, project, "Updating Service", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceApphubServiceRead(d, meta)
}

func resourceApphubServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Service: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{ApphubBasePath}}projects/{{project}}/locations/{{location}}/applications/{{application_id}}/services/{{service_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Service %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Service")
	}

	err = ApphubOperationWaitTime(
		config, res, project, "Deleting Service", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Service %q: %#v", d.Id(), res)
	return nil
}

func resourceApphubServiceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/applications/(?P<application_id>[^/]+)/services/(?P<service_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<application_id>[^/]+)/(?P<service_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<application_id>[^/]+)/(?P<service_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/applications/{{application_id}}/services/{{service_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenApphubServiceName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceServiceReference(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] =
		flattenApphubServiceServiceReferenceUri(original["uri"], d, config)
	return []interface{}{transformed}
}
func flattenApphubServiceServiceReferenceUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceServiceProperties(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["gcp_project"] =
		flattenApphubServiceServicePropertiesGcpProject(original["gcpProject"], d, config)
	transformed["location"] =
		flattenApphubServiceServicePropertiesLocation(original["location"], d, config)
	transformed["zone"] =
		flattenApphubServiceServicePropertiesZone(original["zone"], d, config)
	return []interface{}{transformed}
}
func flattenApphubServiceServicePropertiesGcpProject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceServicePropertiesLocation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceServicePropertiesZone(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["criticality"] =
		flattenApphubServiceAttributesCriticality(original["criticality"], d, config)
	transformed["environment"] =
		flattenApphubServiceAttributesEnvironment(original["environment"], d, config)
	transformed["developer_owners"] =
		flattenApphubServiceAttributesDeveloperOwners(original["developerOwners"], d, config)
	transformed["operator_owners"] =
		flattenApphubServiceAttributesOperatorOwners(original["operatorOwners"], d, config)
	transformed["business_owners"] =
		flattenApphubServiceAttributesBusinessOwners(original["businessOwners"], d, config)
	return []interface{}{transformed}
}
func flattenApphubServiceAttributesCriticality(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["type"] =
		flattenApphubServiceAttributesCriticalityType(original["type"], d, config)
	return []interface{}{transformed}
}
func flattenApphubServiceAttributesCriticalityType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributesEnvironment(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["type"] =
		flattenApphubServiceAttributesEnvironmentType(original["type"], d, config)
	return []interface{}{transformed}
}
func flattenApphubServiceAttributesEnvironmentType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributesDeveloperOwners(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"display_name": flattenApphubServiceAttributesDeveloperOwnersDisplayName(original["displayName"], d, config),
			"email":        flattenApphubServiceAttributesDeveloperOwnersEmail(original["email"], d, config),
		})
	}
	return transformed
}
func flattenApphubServiceAttributesDeveloperOwnersDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributesDeveloperOwnersEmail(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributesOperatorOwners(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"display_name": flattenApphubServiceAttributesOperatorOwnersDisplayName(original["displayName"], d, config),
			"email":        flattenApphubServiceAttributesOperatorOwnersEmail(original["email"], d, config),
		})
	}
	return transformed
}
func flattenApphubServiceAttributesOperatorOwnersDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributesOperatorOwnersEmail(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributesBusinessOwners(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"display_name": flattenApphubServiceAttributesBusinessOwnersDisplayName(original["displayName"], d, config),
			"email":        flattenApphubServiceAttributesBusinessOwnersEmail(original["email"], d, config),
		})
	}
	return transformed
}
func flattenApphubServiceAttributesBusinessOwnersDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceAttributesBusinessOwnersEmail(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceDiscoveredService(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceUid(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubServiceState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandApphubServiceDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCriticality, err := expandApphubServiceAttributesCriticality(original["criticality"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCriticality); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["criticality"] = transformedCriticality
	}

	transformedEnvironment, err := expandApphubServiceAttributesEnvironment(original["environment"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnvironment); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["environment"] = transformedEnvironment
	}

	transformedDeveloperOwners, err := expandApphubServiceAttributesDeveloperOwners(original["developer_owners"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDeveloperOwners); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["developerOwners"] = transformedDeveloperOwners
	}

	transformedOperatorOwners, err := expandApphubServiceAttributesOperatorOwners(original["operator_owners"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOperatorOwners); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["operatorOwners"] = transformedOperatorOwners
	}

	transformedBusinessOwners, err := expandApphubServiceAttributesBusinessOwners(original["business_owners"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedBusinessOwners); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["businessOwners"] = transformedBusinessOwners
	}

	return transformed, nil
}

func expandApphubServiceAttributesCriticality(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedType, err := expandApphubServiceAttributesCriticalityType(original["type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["type"] = transformedType
	}

	return transformed, nil
}

func expandApphubServiceAttributesCriticalityType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributesEnvironment(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedType, err := expandApphubServiceAttributesEnvironmentType(original["type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["type"] = transformedType
	}

	return transformed, nil
}

func expandApphubServiceAttributesEnvironmentType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributesDeveloperOwners(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedDisplayName, err := expandApphubServiceAttributesDeveloperOwnersDisplayName(original["display_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDisplayName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["displayName"] = transformedDisplayName
		}

		transformedEmail, err := expandApphubServiceAttributesDeveloperOwnersEmail(original["email"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEmail); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["email"] = transformedEmail
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandApphubServiceAttributesDeveloperOwnersDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributesDeveloperOwnersEmail(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributesOperatorOwners(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedDisplayName, err := expandApphubServiceAttributesOperatorOwnersDisplayName(original["display_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDisplayName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["displayName"] = transformedDisplayName
		}

		transformedEmail, err := expandApphubServiceAttributesOperatorOwnersEmail(original["email"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEmail); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["email"] = transformedEmail
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandApphubServiceAttributesOperatorOwnersDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributesOperatorOwnersEmail(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributesBusinessOwners(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedDisplayName, err := expandApphubServiceAttributesBusinessOwnersDisplayName(original["display_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedDisplayName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["displayName"] = transformedDisplayName
		}

		transformedEmail, err := expandApphubServiceAttributesBusinessOwnersEmail(original["email"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedEmail); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["email"] = transformedEmail
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandApphubServiceAttributesBusinessOwnersDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceAttributesBusinessOwnersEmail(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandApphubServiceDiscoveredService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
