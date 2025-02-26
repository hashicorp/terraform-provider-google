// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/discoveryengine/TargetSite.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package discoveryengine

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceDiscoveryEngineTargetSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceDiscoveryEngineTargetSiteCreate,
		Read:   resourceDiscoveryEngineTargetSiteRead,
		Delete: resourceDiscoveryEngineTargetSiteDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDiscoveryEngineTargetSiteImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"data_store_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The unique id of the data store.`,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The geographic location where the data store should reside. The value can
only be one of "global", "us" and "eu".`,
			},
			"provided_uri_pattern": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The user provided URI pattern from which the 'generated_uri_pattern' is
generated.`,
			},
			"exact_match": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Description: `If set to false, a uri_pattern is generated to include all pages whose
address contains the provided_uri_pattern. If set to true, an uri_pattern
is generated to try to be an exact match of the provided_uri_pattern or
just the specific page if the provided_uri_pattern is a specific one.
provided_uri_pattern is always normalized to generate the URI pattern to
be used by the search engine.`,
				Default: false,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateEnum([]string{"INCLUDE", "EXCLUDE", ""}),
				Description:  `The possible target site types. Possible values: ["INCLUDE", "EXCLUDE"]`,
			},
			"failure_reason": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Site search indexing failure reasons.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quota_failure": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Site verification state indicating the ownership and validity.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total_required_quota": {
										Type:     schema.TypeInt,
										Optional: true,
										Description: `This number is an estimation on how much total quota this project
needs to successfully complete indexing.`,
									},
								},
							},
						},
					},
				},
			},
			"generated_uri_pattern": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `This is system-generated based on the 'provided_uri_pattern'.`,
			},
			"indexing_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The indexing status.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The unique full resource name of the target site. Values are of the format
'projects/{project}/locations/{location}/collections/{collection_id}/dataStores/{data_store_id}/siteSearchEngine/targetSites/{target_site_id}'.
This field must be a UTF-8 encoded string with a length limit of 1024
characters.`,
			},
			"root_domain_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Root domain of the 'provided_uri_pattern'.`,
			},
			"site_verification_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Site ownership and validity verification status.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_verification_state": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidateEnum([]string{"VERIFIED", "UNVERIFIED", "EXEMPTED", ""}),
							Description:  `Site verification state indicating the ownership and validity. Possible values: ["VERIFIED", "UNVERIFIED", "EXEMPTED"]`,
						},
						"verify_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Latest site verification time.`,
						},
					},
				},
			},
			"target_site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique id of the target site.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The target site's last updated time.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceDiscoveryEngineTargetSiteCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	providedUriPatternProp, err := expandDiscoveryEngineTargetSiteProvidedUriPattern(d.Get("provided_uri_pattern"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("provided_uri_pattern"); !tpgresource.IsEmptyValue(reflect.ValueOf(providedUriPatternProp)) && (ok || !reflect.DeepEqual(v, providedUriPatternProp)) {
		obj["providedUriPattern"] = providedUriPatternProp
	}
	typeProp, err := expandDiscoveryEngineTargetSiteType(d.Get("type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("type"); !tpgresource.IsEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	exactMatchProp, err := expandDiscoveryEngineTargetSiteExactMatch(d.Get("exact_match"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("exact_match"); !tpgresource.IsEmptyValue(reflect.ValueOf(exactMatchProp)) && (ok || !reflect.DeepEqual(v, exactMatchProp)) {
		obj["exactMatch"] = exactMatchProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{DiscoveryEngineBasePath}}projects/{{project}}/locations/{{location}}/collections/default_collection/dataStores/{{data_store_id}}/siteSearchEngine/targetSites")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new TargetSite: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for TargetSite: %s", err)
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
		return fmt.Errorf("Error creating TargetSite: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = DiscoveryEngineOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Creating TargetSite", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")

		return fmt.Errorf("Error waiting to create TargetSite: %s", err)
	}

	if err := d.Set("name", flattenDiscoveryEngineTargetSiteName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = tpgresource.ReplaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating TargetSite %q: %#v", d.Id(), res)

	return resourceDiscoveryEngineTargetSiteRead(d, meta)
}

func resourceDiscoveryEngineTargetSiteRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{DiscoveryEngineBasePath}}{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for TargetSite: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("DiscoveryEngineTargetSite %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}

	if err := d.Set("name", flattenDiscoveryEngineTargetSiteName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("type", flattenDiscoveryEngineTargetSiteType(res["type"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("exact_match", flattenDiscoveryEngineTargetSiteExactMatch(res["exactMatch"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("generated_uri_pattern", flattenDiscoveryEngineTargetSiteGeneratedUriPattern(res["generatedUriPattern"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("root_domain_uri", flattenDiscoveryEngineTargetSiteRootDomainUri(res["rootDomainUri"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("site_verification_info", flattenDiscoveryEngineTargetSiteSiteVerificationInfo(res["siteVerificationInfo"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("indexing_status", flattenDiscoveryEngineTargetSiteIndexingStatus(res["indexingStatus"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("update_time", flattenDiscoveryEngineTargetSiteUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}
	if err := d.Set("failure_reason", flattenDiscoveryEngineTargetSiteFailureReason(res["failureReason"], d, config)); err != nil {
		return fmt.Errorf("Error reading TargetSite: %s", err)
	}

	return nil
}

func resourceDiscoveryEngineTargetSiteDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for TargetSite: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{DiscoveryEngineBasePath}}{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting TargetSite %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "TargetSite")
	}

	err = DiscoveryEngineOperationWaitTime(
		config, res, project, "Deleting TargetSite", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting TargetSite %q: %#v", d.Id(), res)
	return nil
}

func resourceDiscoveryEngineTargetSiteImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/collections/default_collection/dataStores/(?P<data_store_id>[^/]+)/siteSearchEngine/targetSites/(?P<target_site_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Set name based on the components
	if err := d.Set("name", "projects/{{project}}/locations/{{location}}/collections/default_collection/dataStores/{{data_store_id}}/siteSearchEngine/targetSites/{{target_site_id}}"); err != nil {
		return nil, fmt.Errorf("Error setting name: %s", err)
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenDiscoveryEngineTargetSiteName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteExactMatch(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteGeneratedUriPattern(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteRootDomainUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteSiteVerificationInfo(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["site_verification_state"] =
		flattenDiscoveryEngineTargetSiteSiteVerificationInfoSiteVerificationState(original["siteVerificationState"], d, config)
	transformed["verify_time"] =
		flattenDiscoveryEngineTargetSiteSiteVerificationInfoVerifyTime(original["verifyTime"], d, config)
	return []interface{}{transformed}
}
func flattenDiscoveryEngineTargetSiteSiteVerificationInfoSiteVerificationState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteSiteVerificationInfoVerifyTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteIndexingStatus(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDiscoveryEngineTargetSiteFailureReason(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["quota_failure"] =
		flattenDiscoveryEngineTargetSiteFailureReasonQuotaFailure(original["quotaFailure"], d, config)
	return []interface{}{transformed}
}
func flattenDiscoveryEngineTargetSiteFailureReasonQuotaFailure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["total_required_quota"] =
		flattenDiscoveryEngineTargetSiteFailureReasonQuotaFailureTotalRequiredQuota(original["totalRequiredQuota"], d, config)
	return []interface{}{transformed}
}
func flattenDiscoveryEngineTargetSiteFailureReasonQuotaFailureTotalRequiredQuota(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
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

func expandDiscoveryEngineTargetSiteProvidedUriPattern(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDiscoveryEngineTargetSiteType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDiscoveryEngineTargetSiteExactMatch(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
