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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDialogflowAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceDialogflowAgentCreate,
		Read:   resourceDialogflowAgentRead,
		Update: resourceDialogflowAgentUpdate,
		Delete: resourceDialogflowAgentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDialogflowAgentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"default_language_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The default language of the agent as a language tag. [See Language Support](https://cloud.google.com/dialogflow/docs/reference/language) 
for a list of the currently supported language codes. This field cannot be updated after creation.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of this agent.`,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
				Description: `The time zone of this agent from the [time zone database](https://www.iana.org/time-zones), e.g., America/New_York,
Europe/Paris.`,
			},
			"api_version": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"API_VERSION_V1", "API_VERSION_V2", "API_VERSION_V2_BETA_1", ""}, false),
				Description: `API version displayed in Dialogflow console. If not specified, V2 API is assumed. Clients are free to query
different service endpoints for different API versions. However, bots connectors and webhook calls will follow 
the specified API version.
* API_VERSION_V1: Legacy V1 API.
* API_VERSION_V2: V2 API.
* API_VERSION_V2_BETA_1: V2beta1 API.`,
			},
			"avatar_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The URI of the agent's avatar, which are used throughout the Dialogflow console. When an image URL is entered
into this field, the Dialogflow will save the image in the backend. The address of the backend image returned
from the API will be shown in the [avatarUriBackend] field.`,
			},
			"classification_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
				Description: `To filter out false positive results and still get variety in matched natural language inputs for your agent,
you can tune the machine learning classification threshold. If the returned score value is less than the threshold
value, then a fallback intent will be triggered or, if there are no fallback intents defined, no intent will be 
triggered. The score values range from 0.0 (completely uncertain) to 1.0 (completely certain). If set to 0.0, the 
default of 0.3 is used.`,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 500),
				Description:  `The description of this agent. The maximum length is 500 characters. If exceeded, the request is rejected.`,
			},
			"enable_logging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Determines whether this agent should log conversation queries.`,
			},
			"match_mode": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"MATCH_MODE_HYBRID", "MATCH_MODE_ML_ONLY", ""}, false),
				Description: `Determines how intents are detected from user queries.
* MATCH_MODE_HYBRID: Best for agents with a small number of examples in intents and/or wide use of templates
syntax and composite entities.
* MATCH_MODE_ML_ONLY: Can be used for agents with a large number of examples in intents, especially the ones
using @sys.any or very large developer entities.`,
			},
			"supported_language_codes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of all languages supported by this agent (except for the defaultLanguageCode).`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tier": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"TIER_STANDARD", "TIER_ENTERPRISE", "TIER_ENTERPRISE_PLUS", ""}, false),
				Description: `The agent tier. If not specified, TIER_STANDARD is assumed.
* TIER_STANDARD: Standard tier.
* TIER_ENTERPRISE: Enterprise tier (Essentials).
* TIER_ENTERPRISE_PLUS: Enterprise tier (Plus).
NOTE: This field seems to have eventual consistency in the API. Updating this field to a new value, or even 
creating a new agent with a tier that is different from a previous agent in the same project will take some
time to propagate. The provider will wait for the API to show consistency, which can lead to longer apply times.`,
			},
			"avatar_uri_backend": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The URI of the agent's avatar as returned from the API. Output only. To provide an image URL for the agent avatar,
the [avatarUri] field can be used.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDialogflowAgentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	displayNameProp, err := expandDialogflowAgentDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	defaultLanguageCodeProp, err := expandDialogflowAgentDefaultLanguageCode(d.Get("default_language_code"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("default_language_code"); !isEmptyValue(reflect.ValueOf(defaultLanguageCodeProp)) && (ok || !reflect.DeepEqual(v, defaultLanguageCodeProp)) {
		obj["defaultLanguageCode"] = defaultLanguageCodeProp
	}
	supportedLanguageCodesProp, err := expandDialogflowAgentSupportedLanguageCodes(d.Get("supported_language_codes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("supported_language_codes"); !isEmptyValue(reflect.ValueOf(supportedLanguageCodesProp)) && (ok || !reflect.DeepEqual(v, supportedLanguageCodesProp)) {
		obj["supportedLanguageCodes"] = supportedLanguageCodesProp
	}
	timeZoneProp, err := expandDialogflowAgentTimeZone(d.Get("time_zone"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("time_zone"); !isEmptyValue(reflect.ValueOf(timeZoneProp)) && (ok || !reflect.DeepEqual(v, timeZoneProp)) {
		obj["timeZone"] = timeZoneProp
	}
	descriptionProp, err := expandDialogflowAgentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	avatarUriProp, err := expandDialogflowAgentAvatarUri(d.Get("avatar_uri"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("avatar_uri"); !isEmptyValue(reflect.ValueOf(avatarUriProp)) && (ok || !reflect.DeepEqual(v, avatarUriProp)) {
		obj["avatarUri"] = avatarUriProp
	}
	enableLoggingProp, err := expandDialogflowAgentEnableLogging(d.Get("enable_logging"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_logging"); !isEmptyValue(reflect.ValueOf(enableLoggingProp)) && (ok || !reflect.DeepEqual(v, enableLoggingProp)) {
		obj["enableLogging"] = enableLoggingProp
	}
	matchModeProp, err := expandDialogflowAgentMatchMode(d.Get("match_mode"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("match_mode"); !isEmptyValue(reflect.ValueOf(matchModeProp)) && (ok || !reflect.DeepEqual(v, matchModeProp)) {
		obj["matchMode"] = matchModeProp
	}
	classificationThresholdProp, err := expandDialogflowAgentClassificationThreshold(d.Get("classification_threshold"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("classification_threshold"); !isEmptyValue(reflect.ValueOf(classificationThresholdProp)) && (ok || !reflect.DeepEqual(v, classificationThresholdProp)) {
		obj["classificationThreshold"] = classificationThresholdProp
	}
	apiVersionProp, err := expandDialogflowAgentApiVersion(d.Get("api_version"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("api_version"); !isEmptyValue(reflect.ValueOf(apiVersionProp)) && (ok || !reflect.DeepEqual(v, apiVersionProp)) {
		obj["apiVersion"] = apiVersionProp
	}
	tierProp, err := expandDialogflowAgentTier(d.Get("tier"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("tier"); !isEmptyValue(reflect.ValueOf(tierProp)) && (ok || !reflect.DeepEqual(v, tierProp)) {
		obj["tier"] = tierProp
	}

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}projects/{{project}}/agent")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Agent: %#v", obj)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Agent: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{project}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Agent %q: %#v", d.Id(), res)

	// The tier field is eventually consistent, we need to test for consistency before moving on.
	// Otherwise, the user will see diffs on the field.
	if d.HasChange("tier") {
		old, new := d.GetChange("tier")
		readUrl, err := replaceVars(d, config, "{{DialogflowBasePath}}projects/{{project}}/agent")
		if err != nil {
			return fmt.Errorf("Error preparing read URL: %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{old.(string)},
			Target:  []string{new.(string)},
			Refresh: func() (interface{}, string, error) {
				res, err := sendRequest(config, "GET", project, readUrl, nil)
				if err != nil {
					return 0, "", err
				}
				return res, res["tier"].(string), nil
			},
			Timeout:                   40 * time.Minute,
			MinTimeout:                10 * time.Second,
			ContinuousTargetOccurence: 10,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Timed out waiting for agent tier to return correct value.  Waiting for %s, got %s.", new, old)
		}
	}

	return resourceDialogflowAgentRead(d, meta)
}

func resourceDialogflowAgentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}projects/{{project}}/agent")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("DialogflowAgent %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}

	if err := d.Set("display_name", flattenDialogflowAgentDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("default_language_code", flattenDialogflowAgentDefaultLanguageCode(res["defaultLanguageCode"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("supported_language_codes", flattenDialogflowAgentSupportedLanguageCodes(res["supportedLanguageCodes"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("time_zone", flattenDialogflowAgentTimeZone(res["timeZone"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("description", flattenDialogflowAgentDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("avatar_uri_backend", flattenDialogflowAgentAvatarUriBackend(res["avatarUri"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("enable_logging", flattenDialogflowAgentEnableLogging(res["enableLogging"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("match_mode", flattenDialogflowAgentMatchMode(res["matchMode"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("classification_threshold", flattenDialogflowAgentClassificationThreshold(res["classificationThreshold"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("api_version", flattenDialogflowAgentApiVersion(res["apiVersion"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}
	if err := d.Set("tier", flattenDialogflowAgentTier(res["tier"], d, config)); err != nil {
		return fmt.Errorf("Error reading Agent: %s", err)
	}

	return nil
}

func resourceDialogflowAgentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandDialogflowAgentDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	defaultLanguageCodeProp, err := expandDialogflowAgentDefaultLanguageCode(d.Get("default_language_code"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("default_language_code"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, defaultLanguageCodeProp)) {
		obj["defaultLanguageCode"] = defaultLanguageCodeProp
	}
	supportedLanguageCodesProp, err := expandDialogflowAgentSupportedLanguageCodes(d.Get("supported_language_codes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("supported_language_codes"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, supportedLanguageCodesProp)) {
		obj["supportedLanguageCodes"] = supportedLanguageCodesProp
	}
	timeZoneProp, err := expandDialogflowAgentTimeZone(d.Get("time_zone"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("time_zone"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, timeZoneProp)) {
		obj["timeZone"] = timeZoneProp
	}
	descriptionProp, err := expandDialogflowAgentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	avatarUriProp, err := expandDialogflowAgentAvatarUri(d.Get("avatar_uri"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("avatar_uri"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, avatarUriProp)) {
		obj["avatarUri"] = avatarUriProp
	}
	enableLoggingProp, err := expandDialogflowAgentEnableLogging(d.Get("enable_logging"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_logging"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, enableLoggingProp)) {
		obj["enableLogging"] = enableLoggingProp
	}
	matchModeProp, err := expandDialogflowAgentMatchMode(d.Get("match_mode"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("match_mode"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, matchModeProp)) {
		obj["matchMode"] = matchModeProp
	}
	classificationThresholdProp, err := expandDialogflowAgentClassificationThreshold(d.Get("classification_threshold"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("classification_threshold"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, classificationThresholdProp)) {
		obj["classificationThreshold"] = classificationThresholdProp
	}
	apiVersionProp, err := expandDialogflowAgentApiVersion(d.Get("api_version"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("api_version"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, apiVersionProp)) {
		obj["apiVersion"] = apiVersionProp
	}
	tierProp, err := expandDialogflowAgentTier(d.Get("tier"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("tier"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, tierProp)) {
		obj["tier"] = tierProp
	}

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}projects/{{project}}/agent")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Agent %q: %#v", d.Id(), obj)
	_, err = sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating Agent %q: %s", d.Id(), err)
	}

	// The tier field is eventually consistent, we need to test for consistency before moving on.
	// Otherwise, the user will see diffs on the field.
	if d.HasChange("tier") {
		old, new := d.GetChange("tier")
		readUrl, err := replaceVars(d, config, "{{DialogflowBasePath}}projects/{{project}}/agent")
		if err != nil {
			return fmt.Errorf("Error preparing read URL: %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{old.(string)},
			Target:  []string{new.(string)},
			Refresh: func() (interface{}, string, error) {
				res, err := sendRequest(config, "GET", project, readUrl, nil)
				if err != nil {
					return 0, "", err
				}
				return res, res["tier"].(string), nil
			},
			Timeout:                   40 * time.Minute,
			MinTimeout:                10 * time.Second,
			ContinuousTargetOccurence: 10,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Timed out waiting for agent tier to return correct value.  Waiting for %s, got %s.", new, old)
		}
	}
	return resourceDialogflowAgentRead(d, meta)
}

func resourceDialogflowAgentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}projects/{{project}}/agent")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Agent %q", d.Id())

	res, err := sendRequestWithTimeout(config, "DELETE", project, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Agent")
	}

	log.Printf("[DEBUG] Finished deleting Agent %q: %#v", d.Id(), res)
	return nil
}

func resourceDialogflowAgentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"(?P<project>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{project}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenDialogflowAgentDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentDefaultLanguageCode(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentSupportedLanguageCodes(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentTimeZone(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentDescription(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentAvatarUriBackend(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentEnableLogging(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentMatchMode(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentClassificationThreshold(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentApiVersion(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowAgentTier(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandDialogflowAgentDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentDefaultLanguageCode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentSupportedLanguageCodes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentTimeZone(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentAvatarUri(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentEnableLogging(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentMatchMode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentClassificationThreshold(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentApiVersion(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowAgentTier(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
