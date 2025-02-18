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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/securitycenterv2/OrganizationSccBigQueryExport.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package securitycenterv2

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceSecurityCenterV2OrganizationSccBigQueryExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityCenterV2OrganizationSccBigQueryExportCreate,
		Read:   resourceSecurityCenterV2OrganizationSccBigQueryExportRead,
		Update: resourceSecurityCenterV2OrganizationSccBigQueryExportUpdate,
		Delete: resourceSecurityCenterV2OrganizationSccBigQueryExportDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSecurityCenterV2OrganizationSccBigQueryExportImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"big_query_export_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `This must be unique within the organization.`,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The organization whose Cloud Security Command Center the Big Query Export
Config lives in.`,
			},
			"dataset": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The dataset to write findings' updates to.
Its format is "projects/[projectId]/datasets/[bigquery_dataset_id]".
BigQuery Dataset unique ID must contain only letters (a-z, A-Z), numbers (0-9), or underscores (_).`,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
				Description:  `The description of the notification config (max of 1024 characters).`,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Expression that defines the filter to apply across create/update
events of findings. The
expression is a list of zero or more restrictions combined via
logical operators AND and OR. Parentheses are supported, and OR
has higher precedence than AND.

Restrictions have the form <field> <operator> <value> and may have
a - character in front of them to indicate negation. The fields
map to those defined in the corresponding resource.

The supported operators are:

* = for all value types.
* >, <, >=, <= for integer values.
* :, meaning substring matching, for strings.

The supported value types are:

* string literals in quotes.
* integer literals without quotes.
* boolean literals true and false without quotes.

See
[Filtering notifications](https://cloud.google.com/security-command-center/docs/how-to-api-filter-notifications)
for information on how to write a filter.`,
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `location Id is provided by organization. If not provided, Use global as default.`,
				Default:     "global",
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The resource name of this export, in the format
'organizations/{{organization}}/locations/{{location}}/bigQueryExports/{{big_query_export_id}}'.
This field is provided in responses, and is ignored when provided in create requests.`,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The time at which the BigQuery export was created. This field is set by the server and will be ignored if provided on export on creation.
A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits.
Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
			},
			"most_recent_editor": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Email address of the user who last edited the BigQuery export.
This field is set by the server and will be ignored if provided on export creation or update.`,
			},
			"principal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service account that needs permission to create table and upload data to the BigQuery dataset.`,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The most recent time at which the BigQuery export was updated. This field is set by the server and will be ignored if provided on export creation or update.
A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits.
Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSecurityCenterV2OrganizationSccBigQueryExportCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	datasetProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportDataset(d.Get("dataset"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("dataset"); !tpgresource.IsEmptyValue(reflect.ValueOf(datasetProp)) && (ok || !reflect.DeepEqual(v, datasetProp)) {
		obj["dataset"] = datasetProp
	}
	filterProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportFilter(d.Get("filter"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(filterProp)) && (ok || !reflect.DeepEqual(v, filterProp)) {
		obj["filter"] = filterProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterV2BasePath}}organizations/{{organization}}/locations/{{location}}/bigQueryExports?bigQueryExportId={{big_query_export_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new OrganizationSccBigQueryExport: %#v", obj)
	billingProject := ""

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
		return fmt.Errorf("Error creating OrganizationSccBigQueryExport: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/locations/{{location}}/bigQueryExports/{{big_query_export_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating OrganizationSccBigQueryExport %q: %#v", d.Id(), res)

	return resourceSecurityCenterV2OrganizationSccBigQueryExportRead(d, meta)
}

func resourceSecurityCenterV2OrganizationSccBigQueryExportRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterV2BasePath}}organizations/{{organization}}/locations/{{location}}/bigQueryExports/{{big_query_export_id}}")
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SecurityCenterV2OrganizationSccBigQueryExport %q", d.Id()))
	}

	if err := d.Set("name", flattenSecurityCenterV2OrganizationSccBigQueryExportName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}
	if err := d.Set("description", flattenSecurityCenterV2OrganizationSccBigQueryExportDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}
	if err := d.Set("dataset", flattenSecurityCenterV2OrganizationSccBigQueryExportDataset(res["dataset"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}
	if err := d.Set("create_time", flattenSecurityCenterV2OrganizationSccBigQueryExportCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}
	if err := d.Set("update_time", flattenSecurityCenterV2OrganizationSccBigQueryExportUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}
	if err := d.Set("most_recent_editor", flattenSecurityCenterV2OrganizationSccBigQueryExportMostRecentEditor(res["mostRecentEditor"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}
	if err := d.Set("principal", flattenSecurityCenterV2OrganizationSccBigQueryExportPrincipal(res["principal"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}
	if err := d.Set("filter", flattenSecurityCenterV2OrganizationSccBigQueryExportFilter(res["filter"], d, config)); err != nil {
		return fmt.Errorf("Error reading OrganizationSccBigQueryExport: %s", err)
	}

	return nil
}

func resourceSecurityCenterV2OrganizationSccBigQueryExportUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	nameProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	datasetProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportDataset(d.Get("dataset"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("dataset"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, datasetProp)) {
		obj["dataset"] = datasetProp
	}
	filterProp, err := expandSecurityCenterV2OrganizationSccBigQueryExportFilter(d.Get("filter"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("filter"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, filterProp)) {
		obj["filter"] = filterProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterV2BasePath}}organizations/{{organization}}/locations/{{location}}/bigQueryExports/{{big_query_export_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating OrganizationSccBigQueryExport %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("name") {
		updateMask = append(updateMask, "name")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("dataset") {
		updateMask = append(updateMask, "dataset")
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
			return fmt.Errorf("Error updating OrganizationSccBigQueryExport %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating OrganizationSccBigQueryExport %q: %#v", d.Id(), res)
		}

	}

	return resourceSecurityCenterV2OrganizationSccBigQueryExportRead(d, meta)
}

func resourceSecurityCenterV2OrganizationSccBigQueryExportDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{SecurityCenterV2BasePath}}organizations/{{organization}}/locations/{{location}}/bigQueryExports/{{big_query_export_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting OrganizationSccBigQueryExport %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "OrganizationSccBigQueryExport")
	}

	log.Printf("[DEBUG] Finished deleting OrganizationSccBigQueryExport %q: %#v", d.Id(), res)
	return nil
}

func resourceSecurityCenterV2OrganizationSccBigQueryExportImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^organizations/(?P<organization>[^/]+)/locations/(?P<location>[^/]+)/bigQueryExports/(?P<big_query_export_id>[^/]+)$",
		"^(?P<organization>[^/]+)/(?P<location>[^/]+)/(?P<big_query_export_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "organizations/{{organization}}/locations/{{location}}/bigQueryExports/{{big_query_export_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 6 {
		return nil, fmt.Errorf("unexpected format of ID (%q), expected organizations/{{organization}}/locations/{{location}}/bigQueryExports/{{big_query_export_id}}", d.Id())
	}

	if err := d.Set("organization", idParts[1]); err != nil {
		return nil, fmt.Errorf("error setting organization: %s", err)
	}

	if err := d.Set("big_query_export_id", idParts[5]); err != nil {
		return nil, fmt.Errorf("error setting big_query_export_id: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportDataset(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportMostRecentEditor(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportPrincipal(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityCenterV2OrganizationSccBigQueryExportFilter(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandSecurityCenterV2OrganizationSccBigQueryExportName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityCenterV2OrganizationSccBigQueryExportDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityCenterV2OrganizationSccBigQueryExportDataset(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityCenterV2OrganizationSccBigQueryExportFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
