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

package google

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceLoggingLinkedDataset() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingLinkedDatasetCreate,
		Read:   resourceLoggingLinkedDatasetRead,
		Delete: resourceLoggingLinkedDatasetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceLoggingLinkedDatasetImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `The bucket to which the linked dataset is attached.`,
			},
			"link_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The id of the linked dataset.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Describes this link. The maximum length of the description is 8000 characters.`,
			},
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The location of the linked dataset.`,
			},
			"parent": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The parent of the linked dataset.`,
			},
			"bigquery_dataset": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Description: `The information of a BigQuery Dataset. When a link is created, a BigQuery dataset is created along
with it, in the same project as the LogBucket it's linked to. This dataset will also have BigQuery
Views corresponding to the LogViews in the bucket.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dataset_id": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `Output only. The full resource name of the BigQuery dataset. The DATASET_ID will match the ID
of the link, so the link must match the naming restrictions of BigQuery datasets
(alphanumeric characters and underscores only). The dataset will have a resource path of
"bigquery.googleapis.com/projects/[PROJECT_ID]/datasets/[DATASET_ID]"`,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. The creation timestamp of the link. A timestamp in RFC3339 UTC "Zulu" format,
with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z"
and "2014-10-02T15:01:23.045123456Z".`,
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The linked dataset lifecycle state.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name of the linked dataset. The name can have up to 100 characters. A valid link id
(at the end of the link name) must only have alphanumeric characters and underscores within it.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceLoggingLinkedDatasetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandLoggingLinkedDatasetDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	obj, err = resourceLoggingLinkedDatasetEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{LoggingBasePath}}{{parent}}/locations/{{location}}/buckets/{{bucket}}/links?linkId={{link_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new LinkedDataset: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating LinkedDataset: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/locations/{{location}}/buckets/{{bucket}}/links/{{link_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = LoggingOperationWaitTimeWithResponse(
		config, res, &opRes, "Creating LinkedDataset", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")

		return fmt.Errorf("Error waiting to create LinkedDataset: %s", err)
	}

	if err := d.Set("name", flattenLoggingLinkedDatasetName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = tpgresource.ReplaceVars(d, config, "{{parent}}/locations/{{location}}/buckets/{{bucket}}/links/{{link_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating LinkedDataset %q: %#v", d.Id(), res)

	return resourceLoggingLinkedDatasetRead(d, meta)
}

func resourceLoggingLinkedDatasetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{LoggingBasePath}}{{parent}}/locations/{{location}}/buckets/{{bucket}}/links/{{link_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("LoggingLinkedDataset %q", d.Id()))
	}

	if err := d.Set("name", flattenLoggingLinkedDatasetName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading LinkedDataset: %s", err)
	}
	if err := d.Set("description", flattenLoggingLinkedDatasetDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading LinkedDataset: %s", err)
	}
	if err := d.Set("create_time", flattenLoggingLinkedDatasetCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading LinkedDataset: %s", err)
	}
	if err := d.Set("lifecycle_state", flattenLoggingLinkedDatasetLifecycleState(res["lifecycleState"], d, config)); err != nil {
		return fmt.Errorf("Error reading LinkedDataset: %s", err)
	}
	if err := d.Set("bigquery_dataset", flattenLoggingLinkedDatasetBigqueryDataset(res["bigqueryDataset"], d, config)); err != nil {
		return fmt.Errorf("Error reading LinkedDataset: %s", err)
	}

	return nil
}

func resourceLoggingLinkedDatasetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{LoggingBasePath}}{{parent}}/locations/{{location}}/buckets/{{bucket}}/links/{{link_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting LinkedDataset %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "LinkedDataset")
	}

	err = LoggingOperationWaitTime(
		config, res, "Deleting LinkedDataset", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting LinkedDataset %q: %#v", d.Id(), res)
	return nil
}

func resourceLoggingLinkedDatasetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"(?P<parent>.+)/locations/(?P<location>[^/]+)/buckets/(?P<bucket>[^/]+)/links/(?P<link_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/locations/{{location}}/buckets/{{bucket}}/links/{{link_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenLoggingLinkedDatasetName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingLinkedDatasetDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingLinkedDatasetCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingLinkedDatasetLifecycleState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingLinkedDatasetBigqueryDataset(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["dataset_id"] =
		flattenLoggingLinkedDatasetBigqueryDatasetDatasetId(original["datasetId"], d, config)
	return []interface{}{transformed}
}
func flattenLoggingLinkedDatasetBigqueryDatasetDatasetId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandLoggingLinkedDatasetDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func resourceLoggingLinkedDatasetEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// Extract any empty fields from the bucket field.
	parent := d.Get("parent").(string)
	bucket := d.Get("bucket").(string)
	parent, err := ExtractFieldByPattern("parent", parent, bucket, "((projects|folders|organizations|billingAccounts)/[a-z0-9A-Z-]*)/locations/.*")
	if err != nil {
		return nil, fmt.Errorf("error extracting parent field: %s", err)
	}
	location := d.Get("location").(string)
	location, err = ExtractFieldByPattern("location", location, bucket, "[a-zA-Z]*/[a-z0-9A-Z-]*/locations/([a-z0-9-]*)/buckets/.*")
	if err != nil {
		return nil, fmt.Errorf("error extracting location field: %s", err)
	}
	// Set parent to the extracted value.
	d.Set("parent", parent)
	// Set all the other fields to their short forms before forming url and setting ID.
	bucket = tpgresource.GetResourceNameFromSelfLink(bucket)
	name := d.Get("name").(string)
	name = tpgresource.GetResourceNameFromSelfLink(name)
	d.Set("location", location)
	d.Set("bucket", bucket)
	d.Set("name", name)
	return obj, nil
}
