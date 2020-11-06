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
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityPlatformOauthIdpConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityPlatformOauthIdpConfigCreate,
		Read:   resourceIdentityPlatformOauthIdpConfigRead,
		Update: resourceIdentityPlatformOauthIdpConfigUpdate,
		Delete: resourceIdentityPlatformOauthIdpConfigDelete,

		Importer: &schema.ResourceImporter{
			State: resourceIdentityPlatformOauthIdpConfigImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The client id of an OAuth client.`,
			},
			"issuer": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `For OIDC Idps, the issuer identifier.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the OauthIdpConfig. Must start with 'oidc.'.`,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The client secret of the OAuth client, to enable OIDC code flow.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Human friendly display name.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `If this config allows users to sign in with the provider.`,
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

func resourceIdentityPlatformOauthIdpConfigCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandIdentityPlatformOauthIdpConfigName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	displayNameProp, err := expandIdentityPlatformOauthIdpConfigDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	enabledProp, err := expandIdentityPlatformOauthIdpConfigEnabled(d.Get("enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enabled"); !isEmptyValue(reflect.ValueOf(enabledProp)) && (ok || !reflect.DeepEqual(v, enabledProp)) {
		obj["enabled"] = enabledProp
	}
	issuerProp, err := expandIdentityPlatformOauthIdpConfigIssuer(d.Get("issuer"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("issuer"); !isEmptyValue(reflect.ValueOf(issuerProp)) && (ok || !reflect.DeepEqual(v, issuerProp)) {
		obj["issuer"] = issuerProp
	}
	clientIdProp, err := expandIdentityPlatformOauthIdpConfigClientId(d.Get("client_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("client_id"); !isEmptyValue(reflect.ValueOf(clientIdProp)) && (ok || !reflect.DeepEqual(v, clientIdProp)) {
		obj["clientId"] = clientIdProp
	}
	clientSecretProp, err := expandIdentityPlatformOauthIdpConfigClientSecret(d.Get("client_secret"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("client_secret"); !isEmptyValue(reflect.ValueOf(clientSecretProp)) && (ok || !reflect.DeepEqual(v, clientSecretProp)) {
		obj["clientSecret"] = clientSecretProp
	}

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/oauthIdpConfigs?oauthIdpConfigId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new OauthIdpConfig: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthIdpConfig: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating OauthIdpConfig: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/oauthIdpConfigs/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating OauthIdpConfig %q: %#v", d.Id(), res)

	return resourceIdentityPlatformOauthIdpConfigRead(d, meta)
}

func resourceIdentityPlatformOauthIdpConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/oauthIdpConfigs/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthIdpConfig: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("IdentityPlatformOauthIdpConfig %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading OauthIdpConfig: %s", err)
	}

	if err := d.Set("name", flattenIdentityPlatformOauthIdpConfigName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthIdpConfig: %s", err)
	}
	if err := d.Set("display_name", flattenIdentityPlatformOauthIdpConfigDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthIdpConfig: %s", err)
	}
	if err := d.Set("enabled", flattenIdentityPlatformOauthIdpConfigEnabled(res["enabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthIdpConfig: %s", err)
	}
	if err := d.Set("issuer", flattenIdentityPlatformOauthIdpConfigIssuer(res["issuer"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthIdpConfig: %s", err)
	}
	if err := d.Set("client_id", flattenIdentityPlatformOauthIdpConfigClientId(res["clientId"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthIdpConfig: %s", err)
	}
	if err := d.Set("client_secret", flattenIdentityPlatformOauthIdpConfigClientSecret(res["clientSecret"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthIdpConfig: %s", err)
	}

	return nil
}

func resourceIdentityPlatformOauthIdpConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthIdpConfig: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	displayNameProp, err := expandIdentityPlatformOauthIdpConfigDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	enabledProp, err := expandIdentityPlatformOauthIdpConfigEnabled(d.Get("enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enabled"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, enabledProp)) {
		obj["enabled"] = enabledProp
	}
	issuerProp, err := expandIdentityPlatformOauthIdpConfigIssuer(d.Get("issuer"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("issuer"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, issuerProp)) {
		obj["issuer"] = issuerProp
	}
	clientIdProp, err := expandIdentityPlatformOauthIdpConfigClientId(d.Get("client_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("client_id"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, clientIdProp)) {
		obj["clientId"] = clientIdProp
	}
	clientSecretProp, err := expandIdentityPlatformOauthIdpConfigClientSecret(d.Get("client_secret"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("client_secret"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, clientSecretProp)) {
		obj["clientSecret"] = clientSecretProp
	}

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/oauthIdpConfigs/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating OauthIdpConfig %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("enabled") {
		updateMask = append(updateMask, "enabled")
	}

	if d.HasChange("issuer") {
		updateMask = append(updateMask, "issuer")
	}

	if d.HasChange("client_id") {
		updateMask = append(updateMask, "clientId")
	}

	if d.HasChange("client_secret") {
		updateMask = append(updateMask, "clientSecret")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating OauthIdpConfig %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating OauthIdpConfig %q: %#v", d.Id(), res)
	}

	return resourceIdentityPlatformOauthIdpConfigRead(d, meta)
}

func resourceIdentityPlatformOauthIdpConfigDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthIdpConfig: %s", err)
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/oauthIdpConfigs/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting OauthIdpConfig %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "OauthIdpConfig")
	}

	log.Printf("[DEBUG] Finished deleting OauthIdpConfig %q: %#v", d.Id(), res)
	return nil
}

func resourceIdentityPlatformOauthIdpConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/oauthIdpConfigs/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/oauthIdpConfigs/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenIdentityPlatformOauthIdpConfigName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return NameFromSelfLinkStateFunc(v)
}

func flattenIdentityPlatformOauthIdpConfigDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIdentityPlatformOauthIdpConfigEnabled(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIdentityPlatformOauthIdpConfigIssuer(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIdentityPlatformOauthIdpConfigClientId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIdentityPlatformOauthIdpConfigClientSecret(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandIdentityPlatformOauthIdpConfigName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformOauthIdpConfigDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformOauthIdpConfigEnabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformOauthIdpConfigIssuer(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformOauthIdpConfigClientId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformOauthIdpConfigClientSecret(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
