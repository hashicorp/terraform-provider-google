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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/storagecontrol/OrganizationIntelligenceConfig.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package storagecontrol

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceStorageControlOrganizationIntelligenceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageControlOrganizationIntelligenceConfigCreate,
		Read:   resourceStorageControlOrganizationIntelligenceConfigRead,
		Update: resourceStorageControlOrganizationIntelligenceConfigUpdate,
		Delete: resourceStorageControlOrganizationIntelligenceConfigDelete,

		Importer: &schema.ResourceImporter{
			State: resourceStorageControlOrganizationIntelligenceConfigImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Identifier of the GCP Organization. For GCP org, this field should be organization number.`,
			},
			"edition_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `Edition configuration of the Storage Intelligence resource. Valid values are INHERIT, DISABLED, TRIAL and STANDARD.`,
			},
			"filter": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: intelligenceFilterDiffSuppress,
				Description:      `Filter over location and bucket using include or exclude semantics. Resources that match the include or exclude filter are exclusively included or excluded from the Storage Intelligence plan.`,
				MaxItems:         1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"excluded_cloud_storage_buckets": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: intelligenceFilterExcludedCloudStorageBucketsDiffSuppress,
							Description:      `Buckets to exclude from the Storage Intelligence plan.`,
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_id_regexes": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `List of bucket id regexes to exclude in the storage intelligence plan.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
							ConflictsWith: []string{"filter.0.included_cloud_storage_buckets"},
							AtLeastOneOf:  []string{"filter.0.included_cloud_storage_buckets", "filter.0.excluded_cloud_storage_buckets", "filter.0.included_cloud_storage_locations", "filter.0.excluded_cloud_storage_locations"},
						},
						"excluded_cloud_storage_locations": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: intelligenceFilterExcludedCloudStorageLocationsDiffSuppress,
							Description:      `Locations to exclude from the Storage Intelligence plan.`,
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"locations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `List of locations.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
							ConflictsWith: []string{"filter.0.included_cloud_storage_locations"},
							AtLeastOneOf:  []string{"filter.0.included_cloud_storage_buckets", "filter.0.excluded_cloud_storage_buckets", "filter.0.included_cloud_storage_locations", "filter.0.excluded_cloud_storage_locations"},
						},
						"included_cloud_storage_buckets": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: intelligenceFilterincludedCloudStorageBucketsDiffSuppress,
							Description:      `Buckets to include in the Storage Intelligence plan.`,
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_id_regexes": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `List of bucket id regexes to exclude in the storage intelligence plan.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
							ConflictsWith: []string{"filter.0.excluded_cloud_storage_buckets"},
							AtLeastOneOf:  []string{"filter.0.included_cloud_storage_buckets", "filter.0.excluded_cloud_storage_buckets", "filter.0.included_cloud_storage_locations", "filter.0.excluded_cloud_storage_locations"},
						},
						"included_cloud_storage_locations": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: intelligenceFilterincludedCloudStorageLocationsDiffSuppress,
							Description:      `Locations to include in the Storage Intelligence plan.`,
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"locations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `List of locations.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
							ConflictsWith: []string{"filter.0.excluded_cloud_storage_locations"},
							AtLeastOneOf:  []string{"filter.0.included_cloud_storage_buckets", "filter.0.excluded_cloud_storage_buckets", "filter.0.included_cloud_storage_locations", "filter.0.excluded_cloud_storage_locations"},
						},
					},
				},
			},
			"effective_intelligence_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The Intelligence config that is effective for the resource.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effective_edition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The 'StorageIntelligence' edition that is applicable for the resource.`,
						},
						"intelligence_config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Intelligence config resource that is applied for the target resource.`,
						},
					},
				},
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time at which the Storage Intelligence Config resource is last updated.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceStorageControlOrganizationIntelligenceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	editionConfigProp, err := expandStorageControlOrganizationIntelligenceConfigEditionConfig(d.Get("edition_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("edition_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(editionConfigProp)) && (ok || !reflect.DeepEqual(v, editionConfigProp)) {
		obj["editionConfig"] = editionConfigProp
	}
	filterProp, err := expandStorageControlOrganizationIntelligenceConfigFilter(d.Get("filter"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(filterProp)) && (ok || !reflect.DeepEqual(v, filterProp)) {
		obj["filter"] = filterProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageControlBasePath}}organizations/{{name}}/locations/global/intelligenceConfig")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new OrganizationIntelligenceConfig: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	updateMask := []string{"filter"}

	if d.HasChange("edition_config") {
		updateMask = append(updateMask, "editionConfig")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}
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
		return fmt.Errorf("Error creating OrganizationIntelligenceConfig: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{name}}/locations/global/intelligenceConfig")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating OrganizationIntelligenceConfig %q: %#v", d.Id(), res)

	return resourceStorageControlOrganizationIntelligenceConfigRead(d, meta)
}

func resourceStorageControlOrganizationIntelligenceConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageControlBasePath}}organizations/{{name}}/locations/global/intelligenceConfig")
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("StorageControlOrganizationIntelligenceConfig %q", d.Id()))
	}

	if err := d.Set("edition_config", flattenStorageControlOrganizationIntelligenceConfigEditionConfig(res["editionConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationIntelligenceConfig: %s", err)
	}
	if err := d.Set("update_time", flattenStorageControlOrganizationIntelligenceConfigUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationIntelligenceConfig: %s", err)
	}
	if err := d.Set("filter", flattenStorageControlOrganizationIntelligenceConfigFilter(res["filter"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationIntelligenceConfig: %s", err)
	}
	if err := d.Set("effective_intelligence_config", flattenStorageControlOrganizationIntelligenceConfigEffectiveIntelligenceConfig(res["effectiveIntelligenceConfig"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationIntelligenceConfig: %s", err)
	}

	return nil
}

func resourceStorageControlOrganizationIntelligenceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	editionConfigProp, err := expandStorageControlOrganizationIntelligenceConfigEditionConfig(d.Get("edition_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("edition_config"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, editionConfigProp)) {
		obj["editionConfig"] = editionConfigProp
	}
	filterProp, err := expandStorageControlOrganizationIntelligenceConfigFilter(d.Get("filter"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, filterProp)) {
		obj["filter"] = filterProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageControlBasePath}}organizations/{{name}}/locations/global/intelligenceConfig")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating OrganizationIntelligenceConfig %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("edition_config") {
		updateMask = append(updateMask, "editionConfig")
	}

	if d.HasChange("filter") {
		updateMask = append(updateMask, "filter")
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
			return fmt.Errorf("Error updating OrganizationIntelligenceConfig %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating OrganizationIntelligenceConfig %q: %#v", d.Id(), res)
		}

	}

	return resourceStorageControlOrganizationIntelligenceConfigRead(d, meta)
}

func resourceStorageControlOrganizationIntelligenceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] StorageControl OrganizationIntelligenceConfig resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceStorageControlOrganizationIntelligenceConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^organizations/(?P<name>[^/]+)/locations/global/intelligenceConfig$",
		"^(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{name}}/locations/global/intelligenceConfig")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenStorageControlOrganizationIntelligenceConfigEditionConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageControlOrganizationIntelligenceConfigUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageControlOrganizationIntelligenceConfigFilter(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["excluded_cloud_storage_buckets"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBuckets(original["excludedCloudStorageBuckets"], d, config)
	transformed["included_cloud_storage_buckets"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBuckets(original["includedCloudStorageBuckets"], d, config)
	transformed["excluded_cloud_storage_locations"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocations(original["excludedCloudStorageLocations"], d, config)
	transformed["included_cloud_storage_locations"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocations(original["includedCloudStorageLocations"], d, config)
	return []interface{}{transformed}
}
func flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBuckets(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket_id_regexes"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBucketsBucketIdRegexes(original["bucketIdRegexes"], d, config)
	return []interface{}{transformed}
}
func flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBucketsBucketIdRegexes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBuckets(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["bucket_id_regexes"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBucketsBucketIdRegexes(original["bucketIdRegexes"], d, config)
	return []interface{}{transformed}
}
func flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBucketsBucketIdRegexes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocations(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["locations"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocationsLocations(original["locations"], d, config)
	return []interface{}{transformed}
}
func flattenStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocationsLocations(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocations(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["locations"] =
		flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocationsLocations(original["locations"], d, config)
	return []interface{}{transformed}
}
func flattenStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocationsLocations(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageControlOrganizationIntelligenceConfigEffectiveIntelligenceConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["intelligence_config"] =
		flattenStorageControlOrganizationIntelligenceConfigEffectiveIntelligenceConfigIntelligenceConfig(original["intelligenceConfig"], d, config)
	transformed["effective_edition"] =
		flattenStorageControlOrganizationIntelligenceConfigEffectiveIntelligenceConfigEffectiveEdition(original["effectiveEdition"], d, config)
	return []interface{}{transformed}
}
func flattenStorageControlOrganizationIntelligenceConfigEffectiveIntelligenceConfigIntelligenceConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageControlOrganizationIntelligenceConfigEffectiveIntelligenceConfigEffectiveEdition(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandStorageControlOrganizationIntelligenceConfigEditionConfig(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedExcludedCloudStorageBuckets, err := expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBuckets(original["excluded_cloud_storage_buckets"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExcludedCloudStorageBuckets); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["excludedCloudStorageBuckets"] = transformedExcludedCloudStorageBuckets
	}

	transformedIncludedCloudStorageBuckets, err := expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBuckets(original["included_cloud_storage_buckets"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIncludedCloudStorageBuckets); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["includedCloudStorageBuckets"] = transformedIncludedCloudStorageBuckets
	}

	transformedExcludedCloudStorageLocations, err := expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocations(original["excluded_cloud_storage_locations"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExcludedCloudStorageLocations); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["excludedCloudStorageLocations"] = transformedExcludedCloudStorageLocations
	}

	transformedIncludedCloudStorageLocations, err := expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocations(original["included_cloud_storage_locations"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedIncludedCloudStorageLocations); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["includedCloudStorageLocations"] = transformedIncludedCloudStorageLocations
	}

	return transformed, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBuckets(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucketIdRegexes, err := expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBucketsBucketIdRegexes(original["bucket_id_regexes"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["bucketIdRegexes"] = transformedBucketIdRegexes
	}

	return transformed, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageBucketsBucketIdRegexes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBuckets(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedBucketIdRegexes, err := expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBucketsBucketIdRegexes(original["bucket_id_regexes"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["bucketIdRegexes"] = transformedBucketIdRegexes
	}

	return transformed, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageBucketsBucketIdRegexes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocations(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedLocations, err := expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocationsLocations(original["locations"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["locations"] = transformedLocations
	}

	return transformed, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterExcludedCloudStorageLocationsLocations(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocations(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedLocations, err := expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocationsLocations(original["locations"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["locations"] = transformedLocations
	}

	return transformed, nil
}

func expandStorageControlOrganizationIntelligenceConfigFilterIncludedCloudStorageLocationsLocations(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
