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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/chronicle/RuleDeployment.yaml
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

func chronicleRuleDeploymentNilRunFrequencyDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	if new == "" {
		return true
	}
	return false
}

func ResourceChronicleRuleDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceChronicleRuleDeploymentCreate,
		Read:   resourceChronicleRuleDeploymentRead,
		Update: resourceChronicleRuleDeploymentUpdate,
		Delete: resourceChronicleRuleDeploymentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceChronicleRuleDeploymentImport,
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
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The Rule ID of the rule.`,
			},
			"alerting": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Whether detections resulting from this deployment should be considered
alerts.`,
			},
			"archived": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `The archive state of the rule deployment.
Cannot be set to true unless enabled is set to false i.e.
archiving requires a two-step process: first, disable the rule by
setting 'enabled' to false, then set 'archive' to true.
If set to true, alerting will automatically be set to false.
If currently set to true, enabled, alerting, and run_frequency cannot be
updated.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the rule is currently deployed continuously against incoming data.`,
			},
			"run_frequency": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: chronicleRuleDeploymentNilRunFrequencyDiffSuppressFunc,
				Description: `The run frequency of the rule deployment.
Possible values:
LIVE
HOURLY
DAILY`,
			},
			"archive_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The timestamp when the rule deployment archive state was last set to true. If the rule deployment's current archive state is not set to true, the field will be empty.`,
			},
			"consumer_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Description: `Output only. The names of the associated/chained consumer rules. Rules are considered
consumers of this rule if their rule text explicitly filters on this rule's ruleid.
Format:
projects/{project}/locations/{location}/instances/{instance}/rules/{rule}`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"execution_state": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The execution state of the rule deployment.
Possible values:
DEFAULT
LIMITED
PAUSED`,
			},
			"last_alert_status_change_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The timestamp when the rule deployment alert state was lastly changed. This is filled regardless of the current alert state.E.g. if the current alert status is false, this timestamp will be the timestamp when the alert status was changed to false.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The resource name of the rule deployment.
Note that RuleDeployment is a child of the overall Rule, not any individual
revision, so the resource ID segment for the Rule resource must not
reference a specific revision.
Format:
projects/{project}/locations/{location}/instances/{instance}/rules/{rule}/deployment`,
			},
			"producer_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Description: `Output only. The names of the associated/chained producer rules. Rules are considered
producers for this rule if this rule explicitly filters on their ruleid.
Format:
projects/{project}/locations/{location}/instances/{instance}/rules/{rule}`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceChronicleRuleDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule}}/deployment")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RuleDeployment: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ChronicleRuleDeployment %q", d.Id()))
	}

	obj := make(map[string]interface{})
	updateMask := []string{}

	enabledProp, err := expandChronicleRuleDeploymentEnabled(d.Get("enabled"), d, config)
	if err != nil {
		return err
	}
	alertingProp, err := expandChronicleRuleDeploymentAlerting(d.Get("alerting"), d, config)
	if err != nil {
		return err
	}
	archivedProp, err := expandChronicleRuleDeploymentArchived(d.Get("archived"), d, config)
	if err != nil {
		return err
	}
	runFrequencyProp, err := expandChronicleRuleDeploymentRunFrequency(d.Get("run_frequency"), d, config)
	if err != nil {
		return err
	}

	if res != nil {
		enabledValue, enabledExists := res["enabled"]
		if enabledExists {
			enabled := flattenChronicleRuleDeploymentEnabled(enabledValue, d, config)
			if !reflect.DeepEqual(enabledProp, enabled) {
				obj["enabled"] = enabledProp
				updateMask = append(updateMask, "enabled")
			}
		} else {
			// Handle the case where "enabled" is missing from the API response
			if v, ok := d.GetOkExists("enabled"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
				obj["enabled"] = enabledProp
				updateMask = append(updateMask, "enabled")
			}
		}

		alertingValue, alertingExists := res["alerting"]
		if alertingExists {
			alerting := flattenChronicleRuleDeploymentAlerting(alertingValue, d, config)
			if !reflect.DeepEqual(alertingProp, alerting) {
				obj["alerting"] = alertingProp
				updateMask = append(updateMask, "alerting")
			}
		} else {
			// Handle the case where "alerting" is missing from the API response
			if v, ok := d.GetOkExists("alerting"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
				obj["alerting"] = alertingProp
				updateMask = append(updateMask, "alerting")
			}
		}

		archivedValue, archivedExists := res["archived"]
		if archivedExists {
			archived := flattenChronicleRuleDeploymentArchived(archivedValue, d, config)
			if !reflect.DeepEqual(archivedProp, archived) {
				obj["archived"] = archivedProp
				updateMask = append(updateMask, "archived")
			}
		} else {
			// Handle the case where "archived" is missing from the API response
			if v, ok := d.GetOkExists("archived"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
				obj["archived"] = archivedProp
				updateMask = append(updateMask, "archived")
			}
		}

		runFrequencyValue, runFrequencyExists := res["runFrequency"]
		if runFrequencyExists {
			runFrequency := flattenChronicleRuleDeploymentRunFrequency(runFrequencyValue, d, config)
			_, ok := d.GetOkExists("run_frequency")
			if !reflect.DeepEqual(runFrequencyProp, runFrequency) && ok {
				obj["runFrequency"] = runFrequencyProp
				updateMask = append(updateMask, "runFrequency")
			}
		} else {
			// Handle the case where "run_frequency" is missing from the API response
			if v, ok := d.GetOkExists("run_frequency"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
				obj["runFrequency"] = runFrequencyProp
				updateMask = append(updateMask, "runFrequency")
			}
		}
	} else {
		if v, ok := d.GetOkExists("enabled"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
			obj["enabled"] = enabledProp
			updateMask = append(updateMask, "enabled")
		}
		if v, ok := d.GetOkExists("alerting"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
			obj["alerting"] = alertingProp
			updateMask = append(updateMask, "alerting")
		}
		if v, ok := d.GetOkExists("archived"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
			obj["archived"] = archivedProp
			updateMask = append(updateMask, "archived")
		}
		if v, ok := d.GetOkExists("run_frequency"); ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
			obj["runFrequency"] = runFrequencyProp
			updateMask = append(updateMask, "runFrequency")
		}
	}

	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating RuleDeployment %q: %#v", d.Id(), obj)

	res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
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
		return fmt.Errorf("Error updating RuleDeployment %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating RuleDeployment %q: %#v", d.Id(), res)
	}

	if err := d.Set("name", flattenChronicleRuleDeploymentName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule}}/deployment")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating RuleDeployment %q: %#v", d.Id(), res)

	return resourceChronicleRuleDeploymentRead(d, meta)
}

func resourceChronicleRuleDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule}}/deployment")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RuleDeployment: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ChronicleRuleDeployment %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}

	if err := d.Set("name", flattenChronicleRuleDeploymentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("enabled", flattenChronicleRuleDeploymentEnabled(res["enabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("alerting", flattenChronicleRuleDeploymentAlerting(res["alerting"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("archived", flattenChronicleRuleDeploymentArchived(res["archived"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("archive_time", flattenChronicleRuleDeploymentArchiveTime(res["archiveTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("run_frequency", flattenChronicleRuleDeploymentRunFrequency(res["runFrequency"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("execution_state", flattenChronicleRuleDeploymentExecutionState(res["executionState"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("producer_rules", flattenChronicleRuleDeploymentProducerRules(res["producerRules"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("consumer_rules", flattenChronicleRuleDeploymentConsumerRules(res["consumerRules"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}
	if err := d.Set("last_alert_status_change_time", flattenChronicleRuleDeploymentLastAlertStatusChangeTime(res["lastAlertStatusChangeTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading RuleDeployment: %s", err)
	}

	return nil
}

func resourceChronicleRuleDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for RuleDeployment: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	enabledProp, err := expandChronicleRuleDeploymentEnabled(d.Get("enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, enabledProp)) {
		obj["enabled"] = enabledProp
	}
	alertingProp, err := expandChronicleRuleDeploymentAlerting(d.Get("alerting"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("alerting"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, alertingProp)) {
		obj["alerting"] = alertingProp
	}
	archivedProp, err := expandChronicleRuleDeploymentArchived(d.Get("archived"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("archived"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, archivedProp)) {
		obj["archived"] = archivedProp
	}
	runFrequencyProp, err := expandChronicleRuleDeploymentRunFrequency(d.Get("run_frequency"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("run_frequency"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, runFrequencyProp)) {
		obj["runFrequency"] = runFrequencyProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule}}/deployment")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating RuleDeployment %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("enabled") {
		updateMask = append(updateMask, "enabled")
	}

	if d.HasChange("alerting") {
		updateMask = append(updateMask, "alerting")
	}

	if d.HasChange("archived") {
		updateMask = append(updateMask, "archived")
	}

	if d.HasChange("run_frequency") {
		updateMask = append(updateMask, "runFrequency")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}
	// removeRunFrequencyFromUpdateMask removes 'runFrequency' from the updateMask in a URL.
	removeRunFrequencyFromUpdateMask := func(url string) string {
		// Remove "runFrequency" and handle commas.
		url = strings.ReplaceAll(url, "%2CrunFrequency", "")
		url = strings.ReplaceAll(url, "runFrequency%2C", "")
		url = strings.ReplaceAll(url, "runFrequency", "")

		// Remove extra commas.
		url = strings.ReplaceAll(url, "%2C%2C", "%2C")

		//Remove trailing commas.
		url = strings.TrimSuffix(url, "%2C")

		return url
	}

	// Remove "runFrequency" and handle commas if run_frequency not configured by user
	if _, ok := d.GetOk("run_frequency"); !ok {
		url = removeRunFrequencyFromUpdateMask(url)
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
			return fmt.Errorf("Error updating RuleDeployment %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating RuleDeployment %q: %#v", d.Id(), res)
		}

	}

	return resourceChronicleRuleDeploymentRead(d, meta)
}

func resourceChronicleRuleDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] Chronicle RuleDeployment resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceChronicleRuleDeploymentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/instances/(?P<instance>[^/]+)/rules/(?P<rule>[^/]+)/deployment$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<instance>[^/]+)/(?P<rule>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<instance>[^/]+)/(?P<rule>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/rules/{{rule}}/deployment")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenChronicleRuleDeploymentName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentEnabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentAlerting(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentArchived(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentArchiveTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentRunFrequency(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentExecutionState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentProducerRules(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentConsumerRules(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleRuleDeploymentLastAlertStatusChangeTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandChronicleRuleDeploymentEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleRuleDeploymentAlerting(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleRuleDeploymentArchived(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleRuleDeploymentRunFrequency(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
