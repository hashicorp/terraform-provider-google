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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/chronicle/Rule.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package chronicle

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
)

func ResourceChronicleRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceChronicleRuleCreate,
		Read:   resourceChronicleRuleRead,
		Update: resourceChronicleRuleUpdate,
		Delete: resourceChronicleRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceChronicleRuleImport,
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
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The unique identifier for the Chronicle instance, which is the same as the customer ID.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The location of the resource. This is the geographical region where the Chronicle instance resides, such as "us" or "europe-west2".`,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Description: `The etag for this rule.
If this is provided on update, the request will succeed if and only if it
matches the server-computed value, and will fail with an ABORTED error
otherwise.
Populated in BASIC view and FULL view.`,
			},
			"scope": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.ProjectNumberDiffSuppress,
				Description: `Resource name of the DataAccessScope bound to this rule.
Populated in BASIC view and FULL view.
If reference lists are used in the rule, validations will be performed
against this scope to ensure that the reference lists are compatible with
both the user's and the rule's scopes.
The scope should be in the format:
"projects/{project}/locations/{location}/instances/{instance}/dataAccessScopes/{scope}".`,
			},
			"text": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The YARA-L content of the rule.
Populated in FULL view.`,
			},
			"allowed_run_frequencies": {
				Type:     schema.TypeList,
				Computed: true,
				Description: `Output only. The run frequencies that are allowed for the rule.
Populated in BASIC view and FULL view.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"author": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. The author of the rule. Extracted from the meta section of text.
Populated in BASIC view and FULL view.`,
			},
			"compilation_diagnostics": {
				Type:     schema.TypeList,
				Computed: true,
				Description: `Output only. A list of a rule's corresponding compilation diagnostic messages
such as compilation errors and compilation warnings.
Populated in FULL view.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"position": {
							Type:     schema.TypeList,
							Optional: true,
							Description: `CompilationPosition represents the location of a compilation diagnostic in
rule text.`,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_column": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Output only. End column number, beginning at 1.`,
									},
									"end_line": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Output only. End line number, beginning at 1.`,
									},
									"start_column": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Output only. Start column number, beginning at 1.`,
									},
									"start_line": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Output only. Start line number, beginning at 1.`,
									},
								},
							},
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Output only. The diagnostic message.`,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `Output only. The severity of a rule's compilation diagnostic.
Possible values:
SEVERITY_UNSPECIFIED
WARNING
ERROR`,
						},
						"uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Output only. Link to documentation that describes a diagnostic in more detail.`,
						},
					},
				},
			},
			"compilation_state": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. The current compilation state of the rule.
Populated in FULL view.
Possible values:
COMPILATION_STATE_UNSPECIFIED
SUCCEEDED
FAILED`,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. The timestamp of when the rule was created.
Populated in FULL view.`,
			},
			"data_tables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Output only. Resource names of the data tables used in this rule.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. Display name of the rule.
Populated in BASIC view and FULL view.`,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `Output only. Additional metadata specified in the meta section of text.
Populated in FULL view.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Full resource name for the rule. This unique identifier is generated using values provided for the URL parameters.
Format:
projects/{project}/locations/{location}/instances/{instance}/rules/{rule}`,
			},
			"near_real_time_live_rule_eligible": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `Output only. Indicate the rule can run in near real time live rule.
If this is true, the rule uses the near real time live rule when the run
frequency is set to LIVE.`,
			},
			"reference_lists": {
				Type:     schema.TypeList,
				Computed: true,
				Description: `Output only. Resource names of the reference lists used in this rule.
Populated in FULL view.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"revision_create_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. The timestamp of when the rule revision was created.
Populated in FULL, REVISION_METADATA_ONLY views.`,
			},
			"revision_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. The revision ID of the rule.
A new revision is created whenever the rule text is changed in any way.
Format: v_{10 digits}_{9 digits}
Populated in REVISION_METADATA_ONLY view and FULL view.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `Rule Id is the ID of the Rule.`,
			},
			"severity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Severity represents the severity level of the rule.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `The display name of the severity level. Extracted from the meta section of
the rule text.`,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Possible values:
RULE_TYPE_UNSPECIFIED
SINGLE_EVENT
MULTI_EVENT`,
			},
			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Policy to determine if the rule should be deleted forcefully.
If deletion_policy = "FORCE", any retrohunts and any detections associated with the rule
will also be deleted. If deletion_policy = "DEFAULT", the call will only succeed if the
rule has no associated retrohunts, including completed retrohunts, and no
associated detections. Regardless of this field's value, the rule
deployment associated with this rule will also be deleted.
Possible values: DEFAULT, FORCE`,
				Default: "DEFAULT",
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

func resourceChronicleRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	textProp, err := expandChronicleRuleText(d.Get("text"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("text"); !tpgresource.IsEmptyValue(reflect.ValueOf(textProp)) && (ok || !reflect.DeepEqual(v, textProp)) {
		obj["text"] = textProp
	}
	scopeProp, err := expandChronicleRuleScope(d.Get("scope"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("scope"); !tpgresource.IsEmptyValue(reflect.ValueOf(scopeProp)) && (ok || !reflect.DeepEqual(v, scopeProp)) {
		obj["scope"] = scopeProp
	}
	etagProp, err := expandChronicleRuleEtag(d.Get("etag"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("etag"); !tpgresource.IsEmptyValue(reflect.ValueOf(etagProp)) && (ok || !reflect.DeepEqual(v, etagProp)) {
		obj["etag"] = etagProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Rule: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Rule: %s", err)
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
		return fmt.Errorf("Error creating Rule: %s", err)
	}
	// Set computed resource properties from create API response so that they're available on the subsequent Read
	// call.
	// Setting `name` field so that `id_from_name` flattener will work properly.
	if err := d.Set("name", flattenChronicleRuleName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}
	if err := d.Set("rule_id", flattenChronicleRuleRuleId(res["ruleId"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "rule_id": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Rule id is set by API and required to GET the connection
	// it is set by reading the "name" field rather than a field in the response
	if err := d.Set("rule_id", flattenChronicleRuleRuleId("", d, config)); err != nil {
		return fmt.Errorf("Error reading Rule ID: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Rule %q: %#v", d.Id(), res)

	return resourceChronicleRuleRead(d, meta)
}

func resourceChronicleRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Rule: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ChronicleRule %q", d.Id()))
	}

	// Explicitly set virtual fields to default values if unset
	if _, ok := d.GetOkExists("deletion_policy"); !ok {
		if err := d.Set("deletion_policy", "DEFAULT"); err != nil {
			return fmt.Errorf("Error setting deletion_policy: %s", err)
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}

	if err := d.Set("name", flattenChronicleRuleName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("rule_id", flattenChronicleRuleRuleId(res["ruleId"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("text", flattenChronicleRuleText(res["text"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("metadata", flattenChronicleRuleMetadata(res["metadata"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("scope", flattenChronicleRuleScope(res["scope"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("near_real_time_live_rule_eligible", flattenChronicleRuleNearRealTimeLiveRuleEligible(res["nearRealTimeLiveRuleEligible"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("revision_id", flattenChronicleRuleRevisionId(res["revisionId"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("severity", flattenChronicleRuleSeverity(res["severity"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("revision_create_time", flattenChronicleRuleRevisionCreateTime(res["revisionCreateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("compilation_state", flattenChronicleRuleCompilationState(res["compilationState"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("type", flattenChronicleRuleType(res["type"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("reference_lists", flattenChronicleRuleReferenceLists(res["referenceLists"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("display_name", flattenChronicleRuleDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("create_time", flattenChronicleRuleCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("author", flattenChronicleRuleAuthor(res["author"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("allowed_run_frequencies", flattenChronicleRuleAllowedRunFrequencies(res["allowedRunFrequencies"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("etag", flattenChronicleRuleEtag(res["etag"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("compilation_diagnostics", flattenChronicleRuleCompilationDiagnostics(res["compilationDiagnostics"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}
	if err := d.Set("data_tables", flattenChronicleRuleDataTables(res["dataTables"], d, config)); err != nil {
		return fmt.Errorf("Error reading Rule: %s", err)
	}

	return nil
}

func resourceChronicleRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Rule: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	textProp, err := expandChronicleRuleText(d.Get("text"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("text"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, textProp)) {
		obj["text"] = textProp
	}
	scopeProp, err := expandChronicleRuleScope(d.Get("scope"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("scope"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, scopeProp)) {
		obj["scope"] = scopeProp
	}
	etagProp, err := expandChronicleRuleEtag(d.Get("etag"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("etag"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, etagProp)) {
		obj["etag"] = etagProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Rule %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("text") {
		updateMask = append(updateMask, "text")
	}

	if d.HasChange("scope") {
		updateMask = append(updateMask, "scope")
	}

	if d.HasChange("etag") {
		updateMask = append(updateMask, "etag")
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
			return fmt.Errorf("Error updating Rule %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Rule %q: %#v", d.Id(), res)
		}

	}

	return resourceChronicleRuleRead(d, meta)
}

func resourceChronicleRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Rule: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	// Forcefully delete any retrohunts and any detections associated with the rule.
	if deletionPolicy := d.Get("deletion_policy"); deletionPolicy == "FORCE" {
		url = url + "?force=true"
	}

	log.Printf("[DEBUG] Deleting Rule %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "Rule")
	}

	log.Printf("[DEBUG] Finished deleting Rule %q: %#v", d.Id(), res)
	return nil
}

func resourceChronicleRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/instances/(?P<instance>[^/]+)/rules/(?P<rule_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<instance>[^/]+)/(?P<rule_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<instance>[^/]+)/(?P<rule_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Explicitly set virtual fields to default values on import
	if err := d.Set("deletion_policy", "DEFAULT"); err != nil {
		return nil, fmt.Errorf("Error setting deletion_policy: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenChronicleRuleName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleRuleId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	parts := strings.Split(d.Get("name").(string), "/")
	return parts[len(parts)-1]
}

func flattenChronicleRuleText(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleMetadata(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleScope(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleNearRealTimeLiveRuleEligible(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleRevisionId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleSeverity(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["display_name"] =
		flattenChronicleRuleSeverityDisplayName(original["displayName"], d, config)
	return []interface{}{transformed}
}
func flattenChronicleRuleSeverityDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleRevisionCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleCompilationState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleReferenceLists(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleAuthor(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleAllowedRunFrequencies(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleEtag(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleCompilationDiagnostics(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"message":  flattenChronicleRuleCompilationDiagnosticsMessage(original["message"], d, config),
			"position": flattenChronicleRuleCompilationDiagnosticsPosition(original["position"], d, config),
			"severity": flattenChronicleRuleCompilationDiagnosticsSeverity(original["severity"], d, config),
			"uri":      flattenChronicleRuleCompilationDiagnosticsUri(original["uri"], d, config),
		})
	}
	return transformed
}
func flattenChronicleRuleCompilationDiagnosticsMessage(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleCompilationDiagnosticsPosition(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["start_line"] =
		flattenChronicleRuleCompilationDiagnosticsPositionStartLine(original["startLine"], d, config)
	transformed["start_column"] =
		flattenChronicleRuleCompilationDiagnosticsPositionStartColumn(original["startColumn"], d, config)
	transformed["end_line"] =
		flattenChronicleRuleCompilationDiagnosticsPositionEndLine(original["endLine"], d, config)
	transformed["end_column"] =
		flattenChronicleRuleCompilationDiagnosticsPositionEndColumn(original["endColumn"], d, config)
	return []interface{}{transformed}
}
func flattenChronicleRuleCompilationDiagnosticsPositionStartLine(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenChronicleRuleCompilationDiagnosticsPositionStartColumn(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenChronicleRuleCompilationDiagnosticsPositionEndLine(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenChronicleRuleCompilationDiagnosticsPositionEndColumn(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenChronicleRuleCompilationDiagnosticsSeverity(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleCompilationDiagnosticsUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDataTables(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandChronicleRuleText(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleRuleScope(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleRuleEtag(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
