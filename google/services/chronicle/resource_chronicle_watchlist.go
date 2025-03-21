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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/chronicle/Watchlist.yaml
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

func ResourceChronicleWatchlist() *schema.Resource {
	return &schema.Resource{
		Create: resourceChronicleWatchlistCreate,
		Read:   resourceChronicleWatchlistRead,
		Update: resourceChronicleWatchlistUpdate,
		Delete: resourceChronicleWatchlistDelete,

		Importer: &schema.ResourceImporter{
			State: resourceChronicleWatchlistImport,
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
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Required. Display name of the watchlist.
Note that it must be at least one character and less than 63 characters
(https://google.aip.dev/148).`,
			},
			"entity_population_mechanism": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Mechanism to populate entities in the watchlist.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"manual": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Entities are added manually.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
					},
				},
			},
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Optional. Description of the watchlist.`,
			},
			"multiplying_factor": {
				Type:     schema.TypeFloat,
				Optional: true,
				Description: `Optional. Weight applied to the risk score for entities
in this watchlist.
The default is 1.0 if it is not specified.`,
			},
			"watchlist_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Description: `Optional. The ID to use for the watchlist,
which will become the final component of the watchlist's resource name.
This value should be 4-63 characters, and valid characters
are /a-z-/.`,
			},
			"watchlist_user_preferences": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: `A collection of user preferences for watchlist UI configuration.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pinned": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Optional. Whether the watchlist is pinned on the dashboard.`,
						},
					},
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. Time the watchlist was created.`,
			},
			"entity_count": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Count of different types of entities in the watchlist.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Output only. Count of asset type entities in the watchlist.`,
						},
						"user": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Output only. Count of user type entities in the watchlist.`,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Identifier. Resource name of the watchlist. This unique identifier is generated using values provided for the URL parameters.
Format:
projects/{project}/locations/{location}/instances/{instance}/watchlists/{watchlist}`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. Time the watchlist was last updated.`,
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

func resourceChronicleWatchlistCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	multiplyingFactorProp, err := expandChronicleWatchlistMultiplyingFactor(d.Get("multiplying_factor"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("multiplying_factor"); !tpgresource.IsEmptyValue(reflect.ValueOf(multiplyingFactorProp)) && (ok || !reflect.DeepEqual(v, multiplyingFactorProp)) {
		obj["multiplyingFactor"] = multiplyingFactorProp
	}
	displayNameProp, err := expandChronicleWatchlistDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandChronicleWatchlistDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	entityPopulationMechanismProp, err := expandChronicleWatchlistEntityPopulationMechanism(d.Get("entity_population_mechanism"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("entity_population_mechanism"); !tpgresource.IsEmptyValue(reflect.ValueOf(entityPopulationMechanismProp)) && (ok || !reflect.DeepEqual(v, entityPopulationMechanismProp)) {
		obj["entityPopulationMechanism"] = entityPopulationMechanismProp
	}
	watchlistUserPreferencesProp, err := expandChronicleWatchlistWatchlistUserPreferences(d.Get("watchlist_user_preferences"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("watchlist_user_preferences"); !tpgresource.IsEmptyValue(reflect.ValueOf(watchlistUserPreferencesProp)) && (ok || !reflect.DeepEqual(v, watchlistUserPreferencesProp)) {
		obj["watchlistUserPreferences"] = watchlistUserPreferencesProp
	}
	watchlistIdProp, err := expandChronicleWatchlistWatchlistId(d.Get("watchlist_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("watchlist_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(watchlistIdProp)) && (ok || !reflect.DeepEqual(v, watchlistIdProp)) {
		obj["watchlistId"] = watchlistIdProp
	}

	obj, err = resourceChronicleWatchlistEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/watchlists?watchlistId={{watchlist_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Watchlist: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Watchlist: %s", err)
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
		return fmt.Errorf("Error creating Watchlist: %s", err)
	}
	if err := d.Set("name", flattenChronicleWatchlistName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/watchlists/{{watchlist_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	if tpgresource.IsEmptyValue(reflect.ValueOf(d.Get("watchlist_id"))) {
		// watchlist id is set by API when unset and required to GET the connection
		// it is set by reading the "name" field rather than a field in the response
		if err := d.Set("watchlist_id", flattenChronicleWatchlistWatchlistId("", d, config)); err != nil {
			return fmt.Errorf("Error reading Watchlist ID: %s", err)
		}
	}

	log.Printf("[DEBUG] Finished creating Watchlist %q: %#v", d.Id(), res)

	return resourceChronicleWatchlistRead(d, meta)
}

func resourceChronicleWatchlistRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/watchlists/{{watchlist_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Watchlist: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ChronicleWatchlist %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}

	if err := d.Set("name", flattenChronicleWatchlistName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("multiplying_factor", flattenChronicleWatchlistMultiplyingFactor(res["multiplyingFactor"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("create_time", flattenChronicleWatchlistCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("update_time", flattenChronicleWatchlistUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("display_name", flattenChronicleWatchlistDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("description", flattenChronicleWatchlistDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("entity_population_mechanism", flattenChronicleWatchlistEntityPopulationMechanism(res["entityPopulationMechanism"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("entity_count", flattenChronicleWatchlistEntityCount(res["entityCount"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("watchlist_user_preferences", flattenChronicleWatchlistWatchlistUserPreferences(res["watchlistUserPreferences"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}
	if err := d.Set("watchlist_id", flattenChronicleWatchlistWatchlistId(res["watchlistId"], d, config)); err != nil {
		return fmt.Errorf("Error reading Watchlist: %s", err)
	}

	return nil
}

func resourceChronicleWatchlistUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Watchlist: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	multiplyingFactorProp, err := expandChronicleWatchlistMultiplyingFactor(d.Get("multiplying_factor"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("multiplying_factor"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, multiplyingFactorProp)) {
		obj["multiplyingFactor"] = multiplyingFactorProp
	}
	displayNameProp, err := expandChronicleWatchlistDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandChronicleWatchlistDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	entityPopulationMechanismProp, err := expandChronicleWatchlistEntityPopulationMechanism(d.Get("entity_population_mechanism"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("entity_population_mechanism"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, entityPopulationMechanismProp)) {
		obj["entityPopulationMechanism"] = entityPopulationMechanismProp
	}
	watchlistUserPreferencesProp, err := expandChronicleWatchlistWatchlistUserPreferences(d.Get("watchlist_user_preferences"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("watchlist_user_preferences"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, watchlistUserPreferencesProp)) {
		obj["watchlistUserPreferences"] = watchlistUserPreferencesProp
	}

	obj, err = resourceChronicleWatchlistEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/watchlists/{{watchlist_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Watchlist %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("multiplying_factor") {
		updateMask = append(updateMask, "multiplyingFactor")
	}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("entity_population_mechanism") {
		updateMask = append(updateMask, "entityPopulationMechanism")
	}

	if d.HasChange("watchlist_user_preferences") {
		updateMask = append(updateMask, "watchlistUserPreferences")
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
			return fmt.Errorf("Error updating Watchlist %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating Watchlist %q: %#v", d.Id(), res)
		}

	}

	return resourceChronicleWatchlistRead(d, meta)
}

func resourceChronicleWatchlistDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Watchlist: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{ChronicleBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/watchlists/{{watchlist_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Watchlist %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "Watchlist")
	}

	log.Printf("[DEBUG] Finished deleting Watchlist %q: %#v", d.Id(), res)
	return nil
}

func resourceChronicleWatchlistImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/instances/(?P<instance>[^/]+)/watchlists/(?P<watchlist_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<instance>[^/]+)/(?P<watchlist_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<instance>[^/]+)/(?P<watchlist_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/watchlists/{{watchlist_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenChronicleWatchlistName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleWatchlistMultiplyingFactor(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleWatchlistCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleWatchlistUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleWatchlistDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleWatchlistDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleWatchlistEntityPopulationMechanism(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["manual"] =
		flattenChronicleWatchlistEntityPopulationMechanismManual(original["manual"], d, config)
	return []interface{}{transformed}
}
func flattenChronicleWatchlistEntityPopulationMechanismManual(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	transformed := make(map[string]interface{})
	return []interface{}{transformed}
}

func flattenChronicleWatchlistEntityCount(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["user"] =
		flattenChronicleWatchlistEntityCountUser(original["user"], d, config)
	transformed["asset"] =
		flattenChronicleWatchlistEntityCountAsset(original["asset"], d, config)
	return []interface{}{transformed}
}
func flattenChronicleWatchlistEntityCountUser(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenChronicleWatchlistEntityCountAsset(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenChronicleWatchlistWatchlistUserPreferences(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["pinned"] =
		flattenChronicleWatchlistWatchlistUserPreferencesPinned(original["pinned"], d, config)
	return []interface{}{transformed}
}
func flattenChronicleWatchlistWatchlistUserPreferencesPinned(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenChronicleWatchlistWatchlistId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	parts := strings.Split(d.Get("name").(string), "/")
	return parts[len(parts)-1]
}

func expandChronicleWatchlistMultiplyingFactor(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleWatchlistDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleWatchlistDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleWatchlistEntityPopulationMechanism(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedManual, err := expandChronicleWatchlistEntityPopulationMechanismManual(original["manual"], d, config)
	if err != nil {
		return nil, err
	} else {
		transformed["manual"] = transformedManual
	}

	return transformed, nil
}

func expandChronicleWatchlistEntityPopulationMechanismManual(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 {
		return nil, nil
	}

	if l[0] == nil {
		transformed := make(map[string]interface{})
		return transformed, nil
	}
	transformed := make(map[string]interface{})

	return transformed, nil
}

func expandChronicleWatchlistWatchlistUserPreferences(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPinned, err := expandChronicleWatchlistWatchlistUserPreferencesPinned(original["pinned"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPinned); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["pinned"] = transformedPinned
	}

	return transformed, nil
}

func expandChronicleWatchlistWatchlistUserPreferencesPinned(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandChronicleWatchlistWatchlistId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func resourceChronicleWatchlistEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// watchlist_id is needed to qualify the URL but cannot be sent in the body
	delete(obj, "watchlistId")
	return obj, nil
}
