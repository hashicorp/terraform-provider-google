// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package appengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	appengine "google.golang.org/api/appengine/v1"
)

func ResourceAppEngineApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppEngineApplicationCreate,
		Read:   resourceAppEngineApplicationRead,
		Update: resourceAppEngineApplicationUpdate,
		Delete: resourceAppEngineApplicationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			appEngineApplicationLocationIDCustomizeDiff,
		),

		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateProjectID(),
				Description:  `The project ID to create the application under.`,
			},
			"auth_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The domain to authenticate users with when using App Engine's User API.`,
			},
			"location_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The location to serve the app from.`,
			},
			"serving_status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"UNSPECIFIED",
					"SERVING",
					"USER_DISABLED",
					"SYSTEM_DISABLED",
				}, false),
				Computed:    true,
				Description: `The serving status of the app.`,
			},
			"database_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CLOUD_FIRESTORE",
					"CLOUD_DATASTORE_COMPATIBILITY",
					// NOTE: this is provided for compatibility with instances from
					// before CLOUD_DATASTORE_COMPATIBILITY - it cannot be set
					// for new instances.
					"CLOUD_DATASTORE",
				}, false),
				Computed: true,
			},
			"feature_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `A block of optional settings to configure specific App Engine features:`,
				Elem:        appEngineApplicationFeatureSettingsResource(),
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Unique name of the app.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Identifier of the app.`,
			},
			"url_dispatch_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `A list of dispatch rule blocks. Each block has a domain, path, and service field.`,
				Elem:        appEngineApplicationURLDispatchRuleResource(),
			},
			"code_bucket": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The GCS bucket code is being stored in for this app.`,
			},
			"default_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default hostname for this app.`,
			},
			"default_bucket": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The GCS bucket content is being stored in for this app.`,
			},
			"gcr_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The GCR domain used for storing managed Docker images for this app.`,
			},
			"iap": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Settings for enabling Cloud Identity Aware Proxy`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: `Adapted for use with the app`,
						},
						"oauth2_client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `OAuth2 client ID to use for the authentication flow.`,
						},
						"oauth2_client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `OAuth2 client secret to use for the authentication flow. The SHA-256 hash of the value is returned in the oauth2ClientSecretSha256 field.`,
						},
						"oauth2_client_secret_sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: `Hex-encoded SHA-256 hash of the client secret.`,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func appEngineApplicationURLDispatchRuleResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func appEngineApplicationFeatureSettingsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"split_health_checks": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func appEngineApplicationLocationIDCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	old, new := d.GetChange("location_id")
	if old != "" && old != new {
		return fmt.Errorf("Cannot change location_id once the resource is created.")
	}
	return nil
}

func resourceAppEngineApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	app, err := expandAppEngineApplication(d, project)
	if err != nil {
		return err
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "apps/{{project}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	log.Printf("[DEBUG] Creating App Engine App")
	op, err := config.NewAppEngineClient(userAgent).Apps.Create(app).Do()
	if err != nil {
		return fmt.Errorf("Error creating App Engine application: %s", err.Error())
	}

	d.SetId(project)

	// Wait for the operation to complete
	waitErr := AppEngineOperationWaitTime(config, op, project, "App Engine app to create", userAgent, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		d.SetId("")
		return waitErr
	}
	log.Printf("[DEBUG] Created App Engine App")

	return resourceAppEngineApplicationRead(d, meta)
}

func resourceAppEngineApplicationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	pid := d.Id()

	app, err := config.NewAppEngineClient(userAgent).Apps.Get(pid).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("App Engine Application %q", pid))
	}
	if err := d.Set("auth_domain", app.AuthDomain); err != nil {
		return fmt.Errorf("Error setting auth_domain: %s", err)
	}
	if err := d.Set("code_bucket", app.CodeBucket); err != nil {
		return fmt.Errorf("Error setting code_bucket: %s", err)
	}
	if err := d.Set("default_bucket", app.DefaultBucket); err != nil {
		return fmt.Errorf("Error setting default_bucket: %s", err)
	}
	if err := d.Set("default_hostname", app.DefaultHostname); err != nil {
		return fmt.Errorf("Error setting default_hostname: %s", err)
	}
	if err := d.Set("location_id", app.LocationId); err != nil {
		return fmt.Errorf("Error setting location_id: %s", err)
	}
	if err := d.Set("name", app.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("app_id", app.Id); err != nil {
		return fmt.Errorf("Error setting app_id: %s", err)
	}
	if err := d.Set("serving_status", app.ServingStatus); err != nil {
		return fmt.Errorf("Error setting serving_status: %s", err)
	}
	if err := d.Set("gcr_domain", app.GcrDomain); err != nil {
		return fmt.Errorf("Error setting gcr_domain: %s", err)
	}
	if err := d.Set("database_type", app.DatabaseType); err != nil {
		return fmt.Errorf("Error setting database_type: %s", err)
	}
	if err := d.Set("project", pid); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	dispatchRules, err := flattenAppEngineApplicationDispatchRules(app.DispatchRules)
	if err != nil {
		return err
	}
	err = d.Set("url_dispatch_rule", dispatchRules)
	if err != nil {
		return fmt.Errorf("Error setting dispatch rules in state. This is a bug, please report it at https://github.com/hashicorp/terraform-provider-google/issues. Error is:\n%s", err.Error())
	}
	featureSettings, err := flattenAppEngineApplicationFeatureSettings(app.FeatureSettings)
	if err != nil {
		return err
	}
	err = d.Set("feature_settings", featureSettings)
	if err != nil {
		return fmt.Errorf("Error setting feature settings in state. This is a bug, please report it at https://github.com/hashicorp/terraform-provider-google/issues. Error is:\n%s", err.Error())
	}
	iap, err := flattenAppEngineApplicationIap(d, app.Iap)
	if err != nil {
		return err
	}
	err = d.Set("iap", iap)
	if err != nil {
		return fmt.Errorf("Error setting iap in state. This is a bug, please report it at https://github.com/hashicorp/terraform-provider-google/issues. Error is:\n%s", err.Error())
	}
	return nil
}

func resourceAppEngineApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	pid := d.Id()
	app, err := expandAppEngineApplication(d, pid)
	if err != nil {
		return err
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "apps/{{project}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	log.Printf("[DEBUG] Updating App Engine App")
	op, err := config.NewAppEngineClient(userAgent).Apps.Patch(pid, app).UpdateMask("authDomain,databaseType,servingStatus,featureSettings.splitHealthChecks,iap").Do()
	if err != nil {
		return fmt.Errorf("Error updating App Engine application: %s", err.Error())
	}

	// Wait for the operation to complete
	waitErr := AppEngineOperationWaitTime(config, op, pid, "App Engine app to update", userAgent, d.Timeout(schema.TimeoutUpdate))
	if waitErr != nil {
		return waitErr
	}
	log.Printf("[DEBUG] Updated App Engine App")

	return resourceAppEngineApplicationRead(d, meta)
}

func resourceAppEngineApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[WARN] App Engine applications cannot be destroyed once created. The project must be deleted to delete the application.")
	return nil
}

func expandAppEngineApplication(d *schema.ResourceData, project string) (*appengine.Application, error) {
	result := &appengine.Application{
		AuthDomain:    d.Get("auth_domain").(string),
		LocationId:    d.Get("location_id").(string),
		Id:            project,
		GcrDomain:     d.Get("gcr_domain").(string),
		DatabaseType:  d.Get("database_type").(string),
		ServingStatus: d.Get("serving_status").(string),
	}
	featureSettings, err := expandAppEngineApplicationFeatureSettings(d)
	if err != nil {
		return nil, err
	}
	result.FeatureSettings = featureSettings
	iap, err := expandAppEngineApplicationIap(d)
	if err != nil {
		return nil, err
	}
	result.Iap = iap
	return result, nil
}

func expandAppEngineApplicationFeatureSettings(d *schema.ResourceData) (*appengine.FeatureSettings, error) {
	blocks := d.Get("feature_settings").([]interface{})
	if len(blocks) < 1 {
		return nil, nil
	}
	return &appengine.FeatureSettings{
		SplitHealthChecks: d.Get("feature_settings.0.split_health_checks").(bool),
		// force send SplitHealthChecks, so if it's set to false it still gets disabled
		ForceSendFields: []string{"SplitHealthChecks"},
	}, nil
}

func expandAppEngineApplicationIap(d *schema.ResourceData) (*appengine.IdentityAwareProxy, error) {
	blocks := d.Get("iap").([]interface{})
	if len(blocks) < 1 {
		return nil, nil
	}
	return &appengine.IdentityAwareProxy{
		Enabled:                  d.Get("iap.0.enabled").(bool),
		Oauth2ClientId:           d.Get("iap.0.oauth2_client_id").(string),
		Oauth2ClientSecret:       d.Get("iap.0.oauth2_client_secret").(string),
		Oauth2ClientSecretSha256: d.Get("iap.0.oauth2_client_secret_sha256").(string),
	}, nil
}

func flattenAppEngineApplicationFeatureSettings(settings *appengine.FeatureSettings) ([]map[string]interface{}, error) {
	if settings == nil {
		return []map[string]interface{}{}, nil
	}
	result := map[string]interface{}{
		"split_health_checks": settings.SplitHealthChecks,
	}
	return []map[string]interface{}{result}, nil
}

func flattenAppEngineApplicationIap(d *schema.ResourceData, iap *appengine.IdentityAwareProxy) ([]map[string]interface{}, error) {
	if iap == nil {
		return []map[string]interface{}{}, nil
	}
	result := map[string]interface{}{
		"enabled":                     iap.Enabled,
		"oauth2_client_id":            iap.Oauth2ClientId,
		"oauth2_client_secret":        d.Get("iap.0.oauth2_client_secret"),
		"oauth2_client_secret_sha256": iap.Oauth2ClientSecretSha256,
	}
	return []map[string]interface{}{result}, nil
}

func flattenAppEngineApplicationDispatchRules(rules []*appengine.UrlDispatchRule) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		results = append(results, map[string]interface{}{
			"domain":  rule.Domain,
			"path":    rule.Path,
			"service": rule.Service,
		})
	}
	return results, nil
}
