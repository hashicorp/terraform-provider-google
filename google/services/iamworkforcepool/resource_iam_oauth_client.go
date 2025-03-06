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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/iamworkforcepool/OauthClient.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package iamworkforcepool

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

func ResourceIAMWorkforcePoolOauthClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceIAMWorkforcePoolOauthClientCreate,
		Read:   resourceIAMWorkforcePoolOauthClientRead,
		Update: resourceIAMWorkforcePoolOauthClientUpdate,
		Delete: resourceIAMWorkforcePoolOauthClientDelete,

		Importer: &schema.ResourceImporter{
			State: resourceIAMWorkforcePoolOauthClientImport,
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
			"allowed_grant_types": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Required. The list of OAuth grant types is allowed for the OauthClient.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_redirect_uris": {
				Type:     schema.TypeList,
				Required: true,
				Description: `Required. The list of redirect uris that is allowed to redirect back
when authorization process is completed.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_scopes": {
				Type:     schema.TypeList,
				Required: true,
				Description: `Required. The list of scopes that the OauthClient is allowed to request during
OAuth flows.

The following scopes are supported:

* 'https://www.googleapis.com/auth/cloud-platform': See, edit, configure,
and delete your Google Cloud data and see the email address for your Google
Account.
* 'openid': The OAuth client can associate you with your personal
information on Google Cloud.
* 'email': The OAuth client can read a federated identity's email address.
* 'groups': The OAuth client can read a federated identity's groups.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Resource ID segment making up resource 'name'. It identifies the resource within its parent collection as described in https://google.aip.dev/122.`,
			},
			"oauth_client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Required. The ID to use for the OauthClient, which becomes the final component of
the resource name. This value should be a string of 6 to 63 lowercase
letters, digits, or hyphens. It must start with a letter, and cannot have a
trailing hyphen. The prefix 'gcp-' is reserved for use by Google, and may
not be specified.`,
			},
			"client_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `Immutable. The type of OauthClient. Either public or private.
For private clients, the client secret can be managed using the dedicated
OauthClientCredential resource.
Possible values:
CLIENT_TYPE_UNSPECIFIED
PUBLIC_CLIENT
CONFIDENTIAL_CLIENT`,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `A user-specified description of the OauthClient.

Cannot exceed 256 characters.`,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Whether the OauthClient is disabled. You cannot use a disabled OAuth
client.`,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `A user-specified display name of the OauthClient.

Cannot exceed 32 characters.`,
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. The system-generated OauthClient id.`,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Time after which the OauthClient will be permanently purged and cannot
be recovered.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Immutable. Identifier. The resource name of the OauthClient.

Format:'projects/{project}/locations/{location}/oauthClients/{oauth_client}'.`,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The state of the OauthClient.
Possible values:
STATE_UNSPECIFIED
ACTIVE
DELETED`,
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

func resourceIAMWorkforcePoolOauthClientCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	allowedScopesProp, err := expandIAMWorkforcePoolOauthClientAllowedScopes(d.Get("allowed_scopes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allowed_scopes"); !tpgresource.IsEmptyValue(reflect.ValueOf(allowedScopesProp)) && (ok || !reflect.DeepEqual(v, allowedScopesProp)) {
		obj["allowedScopes"] = allowedScopesProp
	}
	disabledProp, err := expandIAMWorkforcePoolOauthClientDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	displayNameProp, err := expandIAMWorkforcePoolOauthClientDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandIAMWorkforcePoolOauthClientDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	allowedGrantTypesProp, err := expandIAMWorkforcePoolOauthClientAllowedGrantTypes(d.Get("allowed_grant_types"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allowed_grant_types"); !tpgresource.IsEmptyValue(reflect.ValueOf(allowedGrantTypesProp)) && (ok || !reflect.DeepEqual(v, allowedGrantTypesProp)) {
		obj["allowedGrantTypes"] = allowedGrantTypesProp
	}
	clientTypeProp, err := expandIAMWorkforcePoolOauthClientClientType(d.Get("client_type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("client_type"); !tpgresource.IsEmptyValue(reflect.ValueOf(clientTypeProp)) && (ok || !reflect.DeepEqual(v, clientTypeProp)) {
		obj["clientType"] = clientTypeProp
	}
	allowedRedirectUrisProp, err := expandIAMWorkforcePoolOauthClientAllowedRedirectUris(d.Get("allowed_redirect_uris"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allowed_redirect_uris"); !tpgresource.IsEmptyValue(reflect.ValueOf(allowedRedirectUrisProp)) && (ok || !reflect.DeepEqual(v, allowedRedirectUrisProp)) {
		obj["allowedRedirectUris"] = allowedRedirectUrisProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{IAMWorkforcePoolBasePath}}projects/{{project}}/locations/{{location}}/oauthClients?oauthClientId={{oauth_client_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new OauthClient: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthClient: %s", err)
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
		return fmt.Errorf("Error creating OauthClient: %s", err)
	}
	if err := d.Set("name", flattenIAMWorkforcePoolOauthClientName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/oauthClients/{{oauth_client_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// This is useful if the resource in question doesn't have a perfectly consistent API
	// That is, the Operation for Create might return before the Get operation shows the
	// completed state of the resource.
	time.Sleep(5 * time.Second)

	log.Printf("[DEBUG] Finished creating OauthClient %q: %#v", d.Id(), res)

	return resourceIAMWorkforcePoolOauthClientRead(d, meta)
}

func resourceIAMWorkforcePoolOauthClientRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{IAMWorkforcePoolBasePath}}projects/{{project}}/locations/{{location}}/oauthClients/{{oauth_client_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthClient: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("IAMWorkforcePoolOauthClient %q", d.Id()))
	}

	res, err = resourceIAMWorkforcePoolOauthClientDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing IAMWorkforcePoolOauthClient because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}

	if err := d.Set("allowed_scopes", flattenIAMWorkforcePoolOauthClientAllowedScopes(res["allowedScopes"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("name", flattenIAMWorkforcePoolOauthClientName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("state", flattenIAMWorkforcePoolOauthClientState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("disabled", flattenIAMWorkforcePoolOauthClientDisabled(res["disabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("client_id", flattenIAMWorkforcePoolOauthClientClientId(res["clientId"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("display_name", flattenIAMWorkforcePoolOauthClientDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("description", flattenIAMWorkforcePoolOauthClientDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("allowed_grant_types", flattenIAMWorkforcePoolOauthClientAllowedGrantTypes(res["allowedGrantTypes"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("expire_time", flattenIAMWorkforcePoolOauthClientExpireTime(res["expireTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("client_type", flattenIAMWorkforcePoolOauthClientClientType(res["clientType"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}
	if err := d.Set("allowed_redirect_uris", flattenIAMWorkforcePoolOauthClientAllowedRedirectUris(res["allowedRedirectUris"], d, config)); err != nil {
		return fmt.Errorf("Error reading OauthClient: %s", err)
	}

	return nil
}

func resourceIAMWorkforcePoolOauthClientUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthClient: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	allowedScopesProp, err := expandIAMWorkforcePoolOauthClientAllowedScopes(d.Get("allowed_scopes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allowed_scopes"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, allowedScopesProp)) {
		obj["allowedScopes"] = allowedScopesProp
	}
	disabledProp, err := expandIAMWorkforcePoolOauthClientDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}
	displayNameProp, err := expandIAMWorkforcePoolOauthClientDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	descriptionProp, err := expandIAMWorkforcePoolOauthClientDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	allowedGrantTypesProp, err := expandIAMWorkforcePoolOauthClientAllowedGrantTypes(d.Get("allowed_grant_types"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allowed_grant_types"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, allowedGrantTypesProp)) {
		obj["allowedGrantTypes"] = allowedGrantTypesProp
	}
	allowedRedirectUrisProp, err := expandIAMWorkforcePoolOauthClientAllowedRedirectUris(d.Get("allowed_redirect_uris"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allowed_redirect_uris"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, allowedRedirectUrisProp)) {
		obj["allowedRedirectUris"] = allowedRedirectUrisProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{IAMWorkforcePoolBasePath}}projects/{{project}}/locations/{{location}}/oauthClients/{{oauth_client_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating OauthClient %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("allowed_scopes") {
		updateMask = append(updateMask, "allowedScopes")
	}

	if d.HasChange("disabled") {
		updateMask = append(updateMask, "disabled")
	}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("allowed_grant_types") {
		updateMask = append(updateMask, "allowedGrantTypes")
	}

	if d.HasChange("allowed_redirect_uris") {
		updateMask = append(updateMask, "allowedRedirectUris")
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
			return fmt.Errorf("Error updating OauthClient %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating OauthClient %q: %#v", d.Id(), res)
		}

	}

	// This is useful if the resource in question doesn't have a perfectly consistent API
	// That is, the Operation for Create might return before the Get operation shows the
	// completed state of the resource.
	time.Sleep(5 * time.Second)
	return resourceIAMWorkforcePoolOauthClientRead(d, meta)
}

func resourceIAMWorkforcePoolOauthClientDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for OauthClient: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{IAMWorkforcePoolBasePath}}projects/{{project}}/locations/{{location}}/oauthClients/{{oauth_client_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting OauthClient %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "OauthClient")
	}

	// This is useful if the resource in question doesn't have a perfectly consistent API
	// That is, the Operation for Create might return before the Get operation shows the
	// completed state of the resource.
	time.Sleep(5 * time.Second)

	log.Printf("[DEBUG] Finished deleting OauthClient %q: %#v", d.Id(), res)
	return nil
}

func resourceIAMWorkforcePoolOauthClientImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/oauthClients/(?P<oauth_client_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<oauth_client_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<oauth_client_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/oauthClients/{{oauth_client_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenIAMWorkforcePoolOauthClientAllowedScopes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientDisabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientClientId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientDisplayName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientAllowedGrantTypes(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientExpireTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientClientType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenIAMWorkforcePoolOauthClientAllowedRedirectUris(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandIAMWorkforcePoolOauthClientAllowedScopes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolOauthClientDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolOauthClientDisplayName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolOauthClientDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolOauthClientAllowedGrantTypes(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolOauthClientClientType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandIAMWorkforcePoolOauthClientAllowedRedirectUris(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func resourceIAMWorkforcePoolOauthClientDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	if v := res["state"]; v == "DELETED" {
		return nil, nil
	}

	return res, nil
}
