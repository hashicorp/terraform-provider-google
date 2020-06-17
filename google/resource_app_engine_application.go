package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	appengine "google.golang.org/api/appengine/v1"
)

func resourceAppEngineApplication() *schema.Resource {
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
				ValidateFunc: validateProjectID(),
			},
			"auth_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Computed: true,
			},
			"database_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CLOUD_FIRESTORE",
					"CLOUD_DATASTORE_COMPATIBILITY",
				}, false),
				Computed: true,
			},
			"feature_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     appEngineApplicationFeatureSettingsResource(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url_dispatch_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     appEngineApplicationURLDispatchRuleResource(),
			},
			"code_bucket": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_bucket": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gcr_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iap": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Settings for enabling Cloud Identity Aware Proxy`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"oauth2_client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"oauth2_client_secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"oauth2_client_secret_sha256": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
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

func appEngineApplicationLocationIDCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	old, new := d.GetChange("location_id")
	if old != "" && old != new {
		return fmt.Errorf("Cannot change location_id once the resource is created.")
	}
	return nil
}

func resourceAppEngineApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	app, err := expandAppEngineApplication(d, project)
	if err != nil {
		return err
	}

	lockName, err := replaceVars(d, config, "apps/{{project}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	log.Printf("[DEBUG] Creating App Engine App")
	op, err := config.clientAppEngine.Apps.Create(app).Do()
	if err != nil {
		return fmt.Errorf("Error creating App Engine application: %s", err.Error())
	}

	d.SetId(project)

	// Wait for the operation to complete
	waitErr := appEngineOperationWaitTime(config, op, project, "App Engine app to create", d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		d.SetId("")
		return waitErr
	}
	log.Printf("[DEBUG] Created App Engine App")

	return resourceAppEngineApplicationRead(d, meta)
}

func resourceAppEngineApplicationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid := d.Id()

	app, err := config.clientAppEngine.Apps.Get(pid).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("App Engine Application %q", pid))
	}
	d.Set("auth_domain", app.AuthDomain)
	d.Set("code_bucket", app.CodeBucket)
	d.Set("default_bucket", app.DefaultBucket)
	d.Set("default_hostname", app.DefaultHostname)
	d.Set("location_id", app.LocationId)
	d.Set("name", app.Name)
	d.Set("app_id", app.Id)
	d.Set("serving_status", app.ServingStatus)
	d.Set("gcr_domain", app.GcrDomain)
	d.Set("database_type", app.DatabaseType)
	d.Set("project", pid)
	dispatchRules, err := flattenAppEngineApplicationDispatchRules(app.DispatchRules)
	if err != nil {
		return err
	}
	err = d.Set("url_dispatch_rule", dispatchRules)
	if err != nil {
		return fmt.Errorf("Error setting dispatch rules in state. This is a bug, please report it at https://github.com/terraform-providers/terraform-provider-google/issues. Error is:\n%s", err.Error())
	}
	featureSettings, err := flattenAppEngineApplicationFeatureSettings(app.FeatureSettings)
	if err != nil {
		return err
	}
	err = d.Set("feature_settings", featureSettings)
	if err != nil {
		return fmt.Errorf("Error setting feature settings in state. This is a bug, please report it at https://github.com/terraform-providers/terraform-provider-google/issues. Error is:\n%s", err.Error())
	}
	iap, err := flattenAppEngineApplicationIap(d, app.Iap)
	if err != nil {
		return err
	}
	err = d.Set("iap", iap)
	if err != nil {
		return fmt.Errorf("Error setting iap in state. This is a bug, please report it at https://github.com/terraform-providers/terraform-provider-google/issues. Error is:\n%s", err.Error())
	}
	return nil
}

func resourceAppEngineApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid := d.Id()
	app, err := expandAppEngineApplication(d, pid)
	if err != nil {
		return err
	}

	lockName, err := replaceVars(d, config, "apps/{{project}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	log.Printf("[DEBUG] Updating App Engine App")
	op, err := config.clientAppEngine.Apps.Patch(pid, app).UpdateMask("authDomain,databaseType,servingStatus,featureSettings.splitHealthChecks,iap").Do()
	if err != nil {
		return fmt.Errorf("Error updating App Engine application: %s", err.Error())
	}

	// Wait for the operation to complete
	waitErr := appEngineOperationWaitTime(config, op, pid, "App Engine app to update", d.Timeout(schema.TimeoutUpdate))
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
