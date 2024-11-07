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

package accesscontextmanager

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"slices"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func AccessContextManagerServicePerimeterIngressPolicyEgressToResourcesDiffSupressFunc(_, _, _ string, d *schema.ResourceData) bool {
	old, new := d.GetChange("egress_to.0.resources")

	oldResources, err := tpgresource.InterfaceSliceToStringSlice(old)
	if err != nil {
		log.Printf("[ERROR] Failed to convert config value: %s", err)
		return false
	}

	newResources, err := tpgresource.InterfaceSliceToStringSlice(new)
	if err != nil {
		log.Printf("[ERROR] Failed to convert config value: %s", err)
		return false
	}

	sort.Strings(oldResources)
	sort.Strings(newResources)

	return slices.Equal(oldResources, newResources)
}

func AccessContextManagerServicePerimeterIngressPolicyIngressToResourcesDiffSupressFunc(_, _, _ string, d *schema.ResourceData) bool {
	old, new := d.GetChange("ingress_to.0.resources")

	oldResources, err := tpgresource.InterfaceSliceToStringSlice(old)
	if err != nil {
		log.Printf("[ERROR] Failed to convert config value: %s", err)
		return false
	}

	newResources, err := tpgresource.InterfaceSliceToStringSlice(new)
	if err != nil {
		log.Printf("[ERROR] Failed to convert config value: %s", err)
		return false
	}

	sort.Strings(oldResources)
	sort.Strings(newResources)

	return slices.Equal(oldResources, newResources)
}

func AccessContextManagerServicePerimeterIngressPolicyIdentityTypeDiffSupressFunc(_, old, new string, _ *schema.ResourceData) bool {
	if old == "" && new == "IDENTITY_TYPE_UNSPECIFIED" {
		return true
	}

	return old == new
}

func ResourceAccessContextManagerServicePerimeterIngressPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessContextManagerServicePerimeterIngressPolicyCreate,
		Read:   resourceAccessContextManagerServicePerimeterIngressPolicyRead,
		Delete: resourceAccessContextManagerServicePerimeterIngressPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"perimeter": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The name of the Service Perimeter to add this resource to.`,
			},
			"ingress_from": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Description: `Defines the conditions on the source of a request causing this 'IngressPolicy'
to apply.`,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identities": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Description: `A list of identities that are allowed access through this ingress policy.
Should be in the format of email address. The email address should represent
individual user or service account only.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"identity_type": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							ValidateFunc:     verify.ValidateEnum([]string{"ANY_IDENTITY", "ANY_USER_ACCOUNT", "ANY_SERVICE_ACCOUNT", ""}),
							DiffSuppressFunc: AccessContextManagerServicePerimeterIdentityTypeDiffSupressFunc,
							Description: `Specifies the type of identities that are allowed access from outside the
perimeter. If left unspecified, then members of 'identities' field will be
allowed access. Possible values: ["ANY_IDENTITY", "ANY_USER_ACCOUNT", "ANY_SERVICE_ACCOUNT"]`,
						},
						"sources": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: `Sources that this 'IngressPolicy' authorizes access from.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_level": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Description: `An 'AccessLevel' resource name that allow resources within the
'ServicePerimeters' to be accessed from the internet. 'AccessLevels' listed
must be in the same policy as this 'ServicePerimeter'. Referencing a nonexistent
'AccessLevel' will cause an error. If no 'AccessLevel' names are listed,
resources within the perimeter can only be accessed via Google Cloud calls
with request origins within the perimeter.
Example 'accessPolicies/MY_POLICY/accessLevels/MY_LEVEL.'
If * is specified, then all IngressSources will be allowed.`,
									},
									"resource": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Description: `A Google Cloud resource that is allowed to ingress the perimeter.
Requests from these resources will be allowed to access perimeter data.
Currently only projects and VPCs are allowed.
Project format: 'projects/{projectNumber}'
VPC network format:
'//compute.googleapis.com/projects/{PROJECT_ID}/global/networks/{NAME}'.
The project may be in any Google Cloud organization, not just the
organization that the perimeter is defined in. '*' is not allowed, the case
of allowing all Google Cloud resources only is not supported.`,
									},
								},
							},
						},
					},
				},
			},
			"ingress_to": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Description: `Defines the conditions on the 'ApiOperation' and request destination that cause
this 'IngressPolicy' to apply.`,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operations": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Description: `A list of 'ApiOperations' the sources specified in corresponding 'IngressFrom'
are allowed to perform in this 'ServicePerimeter'.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method_selectors": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Description: `API methods or permissions to allow. Method or permission must belong to
the service specified by serviceName field. A single 'MethodSelector' entry
with '*' specified for the method field will allow all methods AND
permissions for the service specified in 'serviceName'.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"method": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
													Description: `Value for method should be a valid method name for the corresponding
serviceName in 'ApiOperation'. If '*' used as value for 'method', then
ALL methods and permissions are allowed.`,
												},
												"permission": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
													Description: `Value for permission should be a valid Cloud IAM permission for the
corresponding 'serviceName' in 'ApiOperation'.`,
												},
											},
										},
									},
									"service_name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Description: `The name of the API whose methods or permissions the 'IngressPolicy' or
'EgressPolicy' want to allow. A single 'ApiOperation' with 'serviceName'
field set to '*' will allow all methods AND permissions for all services.`,
									},
								},
							},
						},
						"resources": {
							Type:             schema.TypeList,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: AccessContextManagerServicePerimeterIngressPolicyIngressToResourcesDiffSupressFunc,
							Description: `A list of resources, currently only projects in the form
'projects/<projectnumber>', protected by this 'ServicePerimeter'
that are allowed to be accessed by sources defined in the
corresponding 'IngressFrom'. A request matches if it contains
a resource in this list. If '*' is specified for resources,
then this 'IngressTo' rule will authorize access to all
resources inside the perimeter, provided that the request
also matches the 'operations' field.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func resourceAccessContextManagerServicePerimeterIngressPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	ingressFromProp, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(d.Get("ingress_from"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ingress_from"); !tpgresource.IsEmptyValue(reflect.ValueOf(ingressFromProp)) && (ok || !reflect.DeepEqual(v, ingressFromProp)) {
		obj["ingressFrom"] = ingressFromProp
	}
	ingressToProp, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(d.Get("ingress_to"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ingress_to"); !tpgresource.IsEmptyValue(reflect.ValueOf(ingressToProp)) && (ok || !reflect.DeepEqual(v, ingressToProp)) {
		obj["ingressTo"] = ingressToProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "{{perimeter}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{AccessContextManagerBasePath}}{{perimeter}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new ServicePerimeterIngressPolicy: %#v", obj)

	obj, err = resourceAccessContextManagerServicePerimeterIngressPolicyPatchCreateEncoder(d, meta, obj)
	if err != nil {
		return err
	}
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": "status.ingressPolicies"})
	if err != nil {
		return err
	}
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating ServicePerimeterIngressPolicy: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{perimeter}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = AccessContextManagerOperationWaitTimeWithResponse(
		config, res, &opRes, "Creating ServicePerimeterIngressPolicy", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")

		return fmt.Errorf("Error waiting to create ServicePerimeterIngressPolicy: %s", err)
	}

	if _, ok := opRes["status"]; ok {
		opRes, err = flattenNestedAccessContextManagerServicePerimeterIngressPolicy(d, meta, opRes)
		if err != nil {
			return fmt.Errorf("Error getting nested object from operation response: %s", err)
		}
		if opRes == nil {
			// Object isn't there any more - remove it from the state.
			return fmt.Errorf("Error decoding response from operation, could not find nested object")
		}
	}
	if err := d.Set("ingress_from", flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(opRes["ingressFrom"], d, config)); err != nil {
		return err
	}
	if err := d.Set("ingress_to", flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(opRes["ingressTo"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = tpgresource.ReplaceVars(d, config, "{{perimeter}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating ServicePerimeterIngressPolicy %q: %#v", d.Id(), res)

	return resourceAccessContextManagerServicePerimeterIngressPolicyRead(d, meta)
}

func resourceAccessContextManagerServicePerimeterIngressPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AccessContextManagerBasePath}}{{perimeter}}")
	if err != nil {
		return err
	}

	billingProject := ""

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("AccessContextManagerServicePerimeterIngressPolicy %q", d.Id()))
	}

	res, err = flattenNestedAccessContextManagerServicePerimeterIngressPolicy(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Object isn't there any more - remove it from the state.
		log.Printf("[DEBUG] Removing AccessContextManagerServicePerimeterIngressPolicy because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	if err := d.Set("ingress_from", flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(res["ingressFrom"], d, config)); err != nil {
		return fmt.Errorf("Error reading ServicePerimeterIngressPolicy: %s", err)
	}
	if err := d.Set("ingress_to", flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(res["ingressTo"], d, config)); err != nil {
		return fmt.Errorf("Error reading ServicePerimeterIngressPolicy: %s", err)
	}

	return nil
}

func resourceAccessContextManagerServicePerimeterIngressPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	lockName, err := tpgresource.ReplaceVars(d, config, "{{perimeter}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{AccessContextManagerBasePath}}{{perimeter}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	obj, err = resourceAccessContextManagerServicePerimeterIngressPolicyPatchDeleteEncoder(d, meta, obj)
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "ServicePerimeterIngressPolicy")
	}
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": "status.ingressPolicies"})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting ServicePerimeterIngressPolicy %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "ServicePerimeterIngressPolicy")
	}

	err = AccessContextManagerOperationWaitTime(
		config, res, "Deleting ServicePerimeterIngressPolicy", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting ServicePerimeterIngressPolicy %q: %#v", d.Id(), res)
	return nil
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["identity_type"] =
		flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentityType(original["identityType"], d, config)
	transformed["identities"] =
		flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentities(original["identities"], d, config)
	transformed["sources"] =
		flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSources(original["sources"], d, config)
	return []interface{}{transformed}
}
func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentityType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentities(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSources(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"access_level": flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesAccessLevel(original["accessLevel"], d, config),
			"resource":     flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesResource(original["resource"], d, config),
		})
	}
	return transformed
}
func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesAccessLevel(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesResource(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["resources"] =
		flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToResources(original["resources"], d, config)
	transformed["operations"] =
		flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperations(original["operations"], d, config)
	return []interface{}{transformed}
}
func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToResources(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	rawConfigValue := d.Get("ingress_to.0.resources")

	// Convert config value to []string
	configValue, err := tpgresource.InterfaceSliceToStringSlice(rawConfigValue)
	if err != nil {
		log.Printf("[ERROR] Failed to convert config value: %s", err)
		return v
	}

	// Convert v to []string
	apiStringValue, err := tpgresource.InterfaceSliceToStringSlice(v)
	if err != nil {
		log.Printf("[ERROR] Failed to convert API value: %s", err)
		return v
	}

	sortedStrings, err := tpgresource.SortStringsByConfigOrder(configValue, apiStringValue)
	if err != nil {
		log.Printf("[ERROR] Could not sort API response value: %s", err)
		return v
	}

	return sortedStrings
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperations(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"service_name":     flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsServiceName(original["serviceName"], d, config),
			"method_selectors": flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectors(original["methodSelectors"], d, config),
		})
	}
	return transformed
}
func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsServiceName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectors(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"method":     flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsMethod(original["method"], d, config),
			"permission": flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsPermission(original["permission"], d, config),
		})
	}
	return transformed
}
func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsMethod(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsPermission(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedIdentityType, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentityType(original["identity_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIdentityType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["identityType"] = transformedIdentityType
	}

	transformedIdentities, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentities(original["identities"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIdentities); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["identities"] = transformedIdentities
	}

	transformedSources, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSources(original["sources"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSources); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["sources"] = transformedSources
	}

	return transformed, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentityType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromIdentities(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSources(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedAccessLevel, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesAccessLevel(original["access_level"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedAccessLevel); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["accessLevel"] = transformedAccessLevel
		}

		transformedResource, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesResource(original["resource"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedResource); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["resource"] = transformedResource
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesAccessLevel(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFromSourcesResource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedResources, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToResources(original["resources"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResources); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["resources"] = transformedResources
	}

	transformedOperations, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperations(original["operations"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOperations); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["operations"] = transformedOperations
	}

	return transformed, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToResources(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperations(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedServiceName, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsServiceName(original["service_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedServiceName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["serviceName"] = transformedServiceName
		}

		transformedMethodSelectors, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectors(original["method_selectors"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedMethodSelectors); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["methodSelectors"] = transformedMethodSelectors
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsServiceName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectors(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedMethod, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsMethod(original["method"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedMethod); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["method"] = transformedMethod
		}

		transformedPermission, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsPermission(original["permission"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPermission); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["permission"] = transformedPermission
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsMethod(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressToOperationsMethodSelectorsPermission(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func flattenNestedAccessContextManagerServicePerimeterIngressPolicy(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	var v interface{}
	var ok bool

	v, ok = res["status"]
	if !ok || v == nil {
		return nil, nil
	}
	res = v.(map[string]interface{})

	v, ok = res["ingressPolicies"]
	if !ok || v == nil {
		return nil, nil
	}

	switch v.(type) {
	case []interface{}:
		break
	case map[string]interface{}:
		// Construct list out of single nested resource
		v = []interface{}{v}
	default:
		return nil, fmt.Errorf("expected list or map for value status.ingressPolicies. Actual value: %v", v)
	}

	_, item, err := resourceAccessContextManagerServicePerimeterIngressPolicyFindNestedObjectInList(d, meta, v.([]interface{}))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func resourceAccessContextManagerServicePerimeterIngressPolicyFindNestedObjectInList(d *schema.ResourceData, meta interface{}, items []interface{}) (index int, item map[string]interface{}, err error) {
	expectedIngressFrom, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(d.Get("ingress_from"), d, meta.(*transport_tpg.Config))
	if err != nil {
		return -1, nil, err
	}
	expectedFlattenedIngressFrom := flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(expectedIngressFrom, d, meta.(*transport_tpg.Config))
	expectedIngressTo, err := expandNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(d.Get("ingress_to"), d, meta.(*transport_tpg.Config))
	if err != nil {
		return -1, nil, err
	}
	expectedFlattenedIngressTo := flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(expectedIngressTo, d, meta.(*transport_tpg.Config))

	// Search list for this resource.
	for idx, itemRaw := range items {
		if itemRaw == nil {
			continue
		}
		item := itemRaw.(map[string]interface{})

		itemIngressFrom := flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressFrom(item["ingressFrom"], d, meta.(*transport_tpg.Config))
		// IsEmptyValue check so that if one is nil and the other is "", that's considered a match
		if !(tpgresource.IsEmptyValue(reflect.ValueOf(itemIngressFrom)) && tpgresource.IsEmptyValue(reflect.ValueOf(expectedFlattenedIngressFrom))) && !reflect.DeepEqual(itemIngressFrom, expectedFlattenedIngressFrom) {
			log.Printf("[DEBUG] Skipping item with ingressFrom= %#v, looking for %#v)", itemIngressFrom, expectedFlattenedIngressFrom)
			continue
		}
		itemIngressTo := flattenNestedAccessContextManagerServicePerimeterIngressPolicyIngressTo(item["ingressTo"], d, meta.(*transport_tpg.Config))
		// IsEmptyValue check so that if one is nil and the other is "", that's considered a match
		if !(tpgresource.IsEmptyValue(reflect.ValueOf(itemIngressTo)) && tpgresource.IsEmptyValue(reflect.ValueOf(expectedFlattenedIngressTo))) && !reflect.DeepEqual(itemIngressTo, expectedFlattenedIngressTo) {
			log.Printf("[DEBUG] Skipping item with ingressTo= %#v, looking for %#v)", itemIngressTo, expectedFlattenedIngressTo)
			continue
		}
		log.Printf("[DEBUG] Found item for resource %q: %#v)", d.Id(), item)
		return idx, item, nil
	}
	return -1, nil, nil
}

// PatchCreateEncoder handles creating request data to PATCH parent resource
// with list including new object.
func resourceAccessContextManagerServicePerimeterIngressPolicyPatchCreateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	currItems, err := resourceAccessContextManagerServicePerimeterIngressPolicyListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	_, found, err := resourceAccessContextManagerServicePerimeterIngressPolicyFindNestedObjectInList(d, meta, currItems)
	if err != nil {
		return nil, err
	}

	// Return error if item already created.
	if found != nil {
		return nil, fmt.Errorf("Unable to create ServicePerimeterIngressPolicy, existing object already found: %+v", found)
	}

	// Return list with the resource to create appended
	res := map[string]interface{}{
		"ingressPolicies": append(currItems, obj),
	}
	wrapped := map[string]interface{}{
		"status": res,
	}
	res = wrapped

	return res, nil
}

// PatchDeleteEncoder handles creating request data to PATCH parent resource
// with list excluding object to delete.
func resourceAccessContextManagerServicePerimeterIngressPolicyPatchDeleteEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	currItems, err := resourceAccessContextManagerServicePerimeterIngressPolicyListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	idx, item, err := resourceAccessContextManagerServicePerimeterIngressPolicyFindNestedObjectInList(d, meta, currItems)
	if err != nil {
		return nil, err
	}
	if item == nil {
		// Spoof 404 error for proper handling by Delete (i.e. no-op)
		return nil, tpgresource.Fake404("nested", "AccessContextManagerServicePerimeterIngressPolicy")
	}

	updatedItems := append(currItems[:idx], currItems[idx+1:]...)
	res := map[string]interface{}{
		"ingressPolicies": updatedItems,
	}
	wrapped := map[string]interface{}{
		"status": res,
	}
	res = wrapped

	return res, nil
}

// ListForPatch handles making API request to get parent resource and
// extracting list of objects.
func resourceAccessContextManagerServicePerimeterIngressPolicyListForPatch(d *schema.ResourceData, meta interface{}) ([]interface{}, error) {
	config := meta.(*transport_tpg.Config)
	url, err := tpgresource.ReplaceVars(d, config, "{{AccessContextManagerBasePath}}{{perimeter}}")
	if err != nil {
		return nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, err
	}

	var v interface{}
	var ok bool
	if v, ok = res["status"]; ok && v != nil {
		res = v.(map[string]interface{})
	} else {
		return nil, nil
	}

	v, ok = res["ingressPolicies"]
	if ok && v != nil {
		ls, lsOk := v.([]interface{})
		if !lsOk {
			return nil, fmt.Errorf(`expected list for nested field "ingressPolicies"`)
		}
		return ls, nil
	}
	return nil, nil
}
