package google

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	appengine "google.golang.org/api/appengine/v1"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

// resourceGoogleProject returns a *schema.Resource that allows a customer
// to declare a Google Cloud Project resource.
func resourceGoogleProject() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Create: resourceGoogleProjectCreate,
		Read:   resourceGoogleProjectRead,
		Update: resourceGoogleProjectUpdate,
		Delete: resourceGoogleProjectDelete,

		Importer: &schema.ResourceImporter{
			State: resourceProjectImportState,
		},
		MigrateState:  resourceGoogleProjectMigrateState,
		CustomizeDiff: resourceGoogleProjectCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateProjectID(),
			},
			"skip_delete": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_create_network": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateProjectName(),
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"folder_id": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: parseFolderId,
			},
			"policy_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use the 'google_project_iam_policy' resource to define policies for a Google Project",
			},
			"policy_etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Use the the 'google_project_iam_policy' resource to define policies for a Google Project",
			},
			"number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_account": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"app_engine": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     appEngineResource(),
				MaxItems: 1,
			},
		},
	}
}

func appEngineResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auth_domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"northamerica-northeast1",
					"us-central",
					"us-east1",
					"us-east4",
					"southamerica-east1",
					"europe-west",
					"europe-west2",
					"europe-west3",
					"asia-northeast1",
					"asia-south1",
					"australia-southeast1",
				}, false),
			},
			"serving_status": &schema.Schema{
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
			"feature_settings": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     appEngineFeatureSettingsResource(),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"url_dispatch_rule": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     appEngineURLDispatchRuleResource(),
			},
			"code_bucket": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_hostname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_bucket": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"gcr_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func appEngineURLDispatchRuleResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func appEngineFeatureSettingsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"split_health_checks": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceGoogleProjectCustomizeDiff(diff *schema.ResourceDiff, meta interface{}) error {
	if old, new := diff.GetChange("app_engine.#"); old != nil && new != nil && old.(int) > 0 && new.(int) < 1 {
		// if we're going from app_engine set to unset, we need to delete the project, app_engine has no delete
		return diff.ForceNew("app_engine")
	} else if old, _ := diff.GetChange("app_engine.0.location_id"); diff.HasChange("app_engine.0.location_id") && old != nil && old.(string) != "" {
		// if location_id was already set, and has a new value, that forces a new app
		// if location_id wasn't set, don't force a new value, as we're just enabling app engine
		return diff.ForceNew("app_engine.0.location_id")
	}
	return nil
}

func resourceGoogleProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var pid string
	var err error
	pid = d.Get("project_id").(string)

	log.Printf("[DEBUG]: Creating new project %q", pid)
	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
	}

	if err := getParentResourceId(d, project); err != nil {
		return err
	}

	if _, ok := d.GetOk("labels"); ok {
		project.Labels = expandLabels(d)
	}

	op, err := config.clientResourceManager.Projects.Create(project).Do()
	if err != nil {
		return fmt.Errorf("Error creating project %s (%s): %s.", project.ProjectId, project.Name, err)
	}

	d.SetId(pid)

	// Wait for the operation to complete
	waitErr := resourceManagerOperationWait(config.clientResourceManager, op, "project to create")
	if waitErr != nil {
		// The resource wasn't actually created
		d.SetId("")
		return waitErr
	}

	// Set the billing account
	if _, ok := d.GetOk("billing_account"); ok {
		err = updateProjectBillingAccount(d, config)
		if err != nil {
			return err
		}
	}

	// set up App Engine, too
	app, err := expandAppEngineApp(d)
	if err != nil {
		return err
	}
	if app != nil {
		log.Printf("[DEBUG] Enabling App Engine")
		// enable the app engine APIs so we can create stuff
		if err = enableService("appengine.googleapis.com", project.ProjectId, config); err != nil {
			return fmt.Errorf("Error enabling the App Engine Admin API required to configure App Engine applications: %s", err)
		}
		log.Printf("[DEBUG] Enabled App Engine")
		err = createAppEngineApp(config, pid, app)
		if err != nil {
			return err
		}
	}

	err = resourceGoogleProjectRead(d, meta)
	if err != nil {
		return err
	}

	// There's no such thing as "don't auto-create network", only "delete the network
	// post-creation" - but that's what it's called in the UI and let's not confuse
	// people if we don't have to.  The GCP Console is doing the same thing - creating
	// a network and deleting it in the background.
	if !d.Get("auto_create_network").(bool) {
		// The compute API has to be enabled before we can delete a network.
		if err = enableService("compute.googleapis.com", project.ProjectId, config); err != nil {
			return fmt.Errorf("Error enabling the Compute Engine API required to delete the default network: %s", err)
		}

		if err = forceDeleteComputeNetwork(project.ProjectId, "default", config); err != nil {
			return fmt.Errorf("Error deleting default network in project %s: %s", project.ProjectId, err)
		}
	}
	return nil
}

func createAppEngineApp(config *Config, pid string, app *appengine.Application) error {
	app.Id = pid
	log.Printf("[DEBUG] Creating App Engine App")
	op, err := config.clientAppEngine.Apps.Create(app).Do()
	if err != nil {
		return fmt.Errorf("Error creating App Engine application: %s", err.Error())
	}

	// Wait for the operation to complete
	waitErr := appEngineOperationWait(config.clientAppEngine, op, pid, "App Engine app to create")
	if waitErr != nil {
		return waitErr
	}
	log.Printf("[DEBUG] Created App Engine App")
	return nil
}

func resourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid := d.Id()

	// Read the project
	p, err := config.clientResourceManager.Projects.Get(pid).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project %q", pid))
	}

	// If the project has been deleted from outside Terraform, remove it from state file.
	if p.LifecycleState != "ACTIVE" {
		log.Printf("[WARN] Removing project '%s' because its state is '%s' (requires 'ACTIVE').", pid, p.LifecycleState)
		d.SetId("")
		return nil
	}

	d.Set("project_id", pid)
	d.Set("number", strconv.FormatInt(int64(p.ProjectNumber), 10))
	d.Set("name", p.Name)
	d.Set("labels", p.Labels)

	if p.Parent != nil {
		switch p.Parent.Type {
		case "organization":
			d.Set("org_id", p.Parent.Id)
			d.Set("folder_id", "")
		case "folder":
			d.Set("folder_id", p.Parent.Id)
			d.Set("org_id", "")
		}
	}

	// Read the billing account
	ba, err := config.clientBilling.Projects.GetBillingInfo(prefixedProject(pid)).Do()
	if err != nil && !isApiNotEnabledError(err) {
		return fmt.Errorf("Error reading billing account for project %q: %v", prefixedProject(pid), err)
	} else if isApiNotEnabledError(err) {
		log.Printf("[WARN] Billing info API not enabled, please enable it to read billing info about project %q: %s", pid, err.Error())
	} else if ba.BillingAccountName != "" {
		// BillingAccountName is contains the resource name of the billing account
		// associated with the project, if any. For example,
		// `billingAccounts/012345-567890-ABCDEF`. We care about the ID and not
		// the `billingAccounts/` prefix, so we need to remove that. If the
		// prefix ever changes, we'll validate to make sure it's something we
		// recognize.
		_ba := strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/")
		if ba.BillingAccountName == _ba {
			return fmt.Errorf("Error parsing billing account for project %q. Expected value to begin with 'billingAccounts/' but got %s", prefixedProject(pid), ba.BillingAccountName)
		}
		d.Set("billing_account", _ba)
	}

	// read the App Engine app, if one exists
	// we don't have the config available for import, so we can't rely on
	// that to read it. And honestly, we want to know if an App exists that
	// shouldn't. So this tries to read it, sets it to empty if none exists,
	// or sets it in state if one does exist.
	app, err := config.clientAppEngine.Apps.Get(pid).Do()
	if err != nil && !isGoogleApiErrorWithCode(err, 404) && !isApiNotEnabledError(err) {
		return fmt.Errorf("Error retrieving App Engine application %q: %s", pid, err.Error())
	} else if isGoogleApiErrorWithCode(err, 404) {
		d.Set("app_engine", []map[string]interface{}{})
	} else if isApiNotEnabledError(err) {
		log.Printf("[WARN] App Engine Admin API not enabled, please enable it to read App Engine info about project %q: %s", pid, err.Error())
	} else {
		appBlocks, err := flattenAppEngineApp(app)
		if err != nil {
			return fmt.Errorf("Error serializing App Engine app: %s", err.Error())
		}
		err = d.Set("app_engine", appBlocks)
		if err != nil {
			return fmt.Errorf("Error setting App Engine application in state. This is a bug, please report it at https://github.com/terraform-providers/terraform-provider-google/issues. Error is:\n%s", err.Error())
		}
	}
	return nil
}

func prefixedProject(pid string) string {
	return "projects/" + pid
}

func getParentResourceId(d *schema.ResourceData, p *cloudresourcemanager.Project) error {
	orgId := d.Get("org_id").(string)
	folderId := d.Get("folder_id").(string)

	if orgId != "" && folderId != "" {
		return fmt.Errorf("'org_id' and 'folder_id' cannot be both set.")
	}

	if orgId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   orgId,
			Type: "organization",
		}
	}

	if folderId != "" {
		p.Parent = &cloudresourcemanager.ResourceId{
			Id:   parseFolderId(folderId),
			Type: "folder",
		}
	}

	return nil
}

func parseFolderId(v interface{}) string {
	folderId := v.(string)
	if strings.HasPrefix(folderId, "folders/") {
		return folderId[8:]
	}
	return folderId
}

func resourceGoogleProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid := d.Id()
	project_name := d.Get("name").(string)

	// Read the project
	// we need the project even though refresh has already been called
	// because the API doesn't support patch, so we need the actual object
	p, err := config.clientResourceManager.Projects.Get(pid).Do()
	if err != nil {
		if v, ok := err.(*googleapi.Error); ok && v.Code == http.StatusNotFound {
			return fmt.Errorf("Project %q does not exist.", pid)
		}
		return fmt.Errorf("Error checking project %q: %s", pid, err)
	}

	d.Partial(true)

	// Project display name has changed
	if ok := d.HasChange("name"); ok {
		p.Name = project_name
		// Do update on project
		p, err = config.clientResourceManager.Projects.Update(p.ProjectId, p).Do()
		if err != nil {
			return fmt.Errorf("Error updating project %q: %s", project_name, err)
		}
		d.SetPartial("name")
	}

	// Project parent has changed
	if d.HasChange("org_id") || d.HasChange("folder_id") {
		if err := getParentResourceId(d, p); err != nil {
			return err
		}

		// Do update on project
		p, err = config.clientResourceManager.Projects.Update(p.ProjectId, p).Do()
		if err != nil {
			return fmt.Errorf("Error updating project %q: %s", project_name, err)
		}
		d.SetPartial("org_id")
		d.SetPartial("folder_id")
	}

	// Billing account has changed
	if ok := d.HasChange("billing_account"); ok {
		err = updateProjectBillingAccount(d, config)
		if err != nil {
			return err
		}
	}

	// Project Labels have changed
	if ok := d.HasChange("labels"); ok {
		p.Labels = expandLabels(d)

		// Do Update on project
		p, err = config.clientResourceManager.Projects.Update(p.ProjectId, p).Do()
		if err != nil {
			return fmt.Errorf("Error updating project %q: %s", project_name, err)
		}
		d.SetPartial("labels")
	}

	// App Engine App has changed
	if ok := d.HasChange("app_engine"); ok {
		app, err := expandAppEngineApp(d)
		if err != nil {
			return err
		}
		// ignore if app is now not set; that should force new resource using customizediff
		if app != nil {
			if old, new := d.GetChange("app_engine.#"); (old == nil || old.(int) < 1) && new != nil && new.(int) > 0 {
				err = createAppEngineApp(config, pid, app)
				if err != nil {
					return err
				}
			} else {
				log.Printf("[DEBUG] Updating App Engine App")
				op, err := config.clientAppEngine.Apps.Patch(pid, app).UpdateMask("authDomain,servingStatus,featureSettings.splitHealthChecks").Do()
				if err != nil {
					return fmt.Errorf("Error creating App Engine application: %s", err.Error())
				}

				// Wait for the operation to complete
				waitErr := appEngineOperationWait(config.clientAppEngine, op, pid, "App Engine app to update")
				if waitErr != nil {
					return waitErr
				}
				log.Printf("[DEBUG] Updated App Engine App")
			}
			d.SetPartial("app_engine")
		}
	}

	d.Partial(false)

	return resourceGoogleProjectRead(d, meta)
}

func resourceGoogleProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	// Only delete projects if skip_delete isn't set
	if !d.Get("skip_delete").(bool) {
		pid := d.Id()
		_, err := config.clientResourceManager.Projects.Delete(pid).Do()
		if err != nil {
			return fmt.Errorf("Error deleting project %q: %s", pid, err)
		}
	}
	d.SetId("")
	return nil
}

func resourceProjectImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Explicitly set to default as a workaround for `ImportStateVerify` tests, and so that users
	// don't see a diff immediately after import.
	d.Set("auto_create_network", true)
	return []*schema.ResourceData{d}, nil
}

// Delete a compute network along with the firewall rules inside it.
func forceDeleteComputeNetwork(projectId, networkName string, config *Config) error {
	networkLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", projectId, networkName)

	token := ""
	for paginate := true; paginate; {
		filter := fmt.Sprintf("network eq %s", networkLink)
		resp, err := config.clientCompute.Firewalls.List(projectId).Filter(filter).Do()
		if err != nil {
			return fmt.Errorf("Error listing firewall rules in proj: %s", err)
		}

		log.Printf("[DEBUG] Found %d firewall rules in %q network", len(resp.Items), networkName)

		for _, firewall := range resp.Items {
			op, err := config.clientCompute.Firewalls.Delete(projectId, firewall.Name).Do()
			if err != nil {
				return fmt.Errorf("Error deleting firewall: %s", err)
			}
			err = computeSharedOperationWait(config.clientCompute, op, projectId, "Deleting Firewall")
			if err != nil {
				return err
			}
		}

		token = resp.NextPageToken
		paginate = token != ""
	}

	return deleteComputeNetwork(projectId, networkName, config)
}

func updateProjectBillingAccount(d *schema.ResourceData, config *Config) error {
	pid := d.Id()
	name := d.Get("billing_account").(string)
	ba := cloudbilling.ProjectBillingInfo{}
	// If we're unlinking an existing billing account, an empty request does that, not an empty-string billing account.
	if name != "" {
		ba.BillingAccountName = "billingAccounts/" + name
	}
	_, err := config.clientBilling.Projects.UpdateBillingInfo(prefixedProject(pid), &ba).Do()
	if err != nil {
		d.Set("billing_account", "")
		if _err, ok := err.(*googleapi.Error); ok {
			return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), _err)
		}
		return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), err)
	}
	for retries := 0; retries < 3; retries++ {
		err = resourceGoogleProjectRead(d, config)
		if err != nil {
			return err
		}
		if d.Get("billing_account").(string) == name {
			break
		}
		time.Sleep(3)
	}
	if d.Get("billing_account").(string) != name {
		return fmt.Errorf("Timed out waiting for billing account to return correct value.  Waiting for %s, got %s.",
			d.Get("billding_account").(string), name)
	}
	return nil
}

func expandAppEngineApp(d *schema.ResourceData) (*appengine.Application, error) {
	blocks := d.Get("app_engine").([]interface{})
	if len(blocks) < 1 {
		return nil, nil
	}
	if len(blocks) > 1 {
		return nil, fmt.Errorf("only one app_engine block may be defined per project")
	}
	result := &appengine.Application{
		AuthDomain:    d.Get("app_engine.0.auth_domain").(string),
		LocationId:    d.Get("app_engine.0.location_id").(string),
		Id:            d.Get("project_id").(string),
		GcrDomain:     d.Get("app_engine.0.gcr_domain").(string),
		ServingStatus: d.Get("app_engine.0.serving_status").(string),
	}
	featureSettings, err := expandAppEngineFeatureSettings(d, "app_engine.0.")
	if err != nil {
		return nil, err
	}
	result.FeatureSettings = featureSettings
	return result, nil
}

func flattenAppEngineApp(app *appengine.Application) ([]map[string]interface{}, error) {
	result := map[string]interface{}{
		"auth_domain":      app.AuthDomain,
		"code_bucket":      app.CodeBucket,
		"default_bucket":   app.DefaultBucket,
		"default_hostname": app.DefaultHostname,
		"location_id":      app.LocationId,
		"name":             app.Name,
		"serving_status":   app.ServingStatus,
	}
	dispatchRules, err := flattenAppEngineDispatchRules(app.DispatchRules)
	if err != nil {
		return nil, err
	}
	result["url_dispatch_rule"] = dispatchRules
	featureSettings, err := flattenAppEngineFeatureSettings(app.FeatureSettings)
	if err != nil {
		return nil, err
	}
	result["feature_settings"] = featureSettings
	return []map[string]interface{}{result}, nil
}

func expandAppEngineFeatureSettings(d *schema.ResourceData, prefix string) (*appengine.FeatureSettings, error) {
	blocks := d.Get(prefix + "feature_settings").([]interface{})
	if len(blocks) < 1 {
		return nil, nil
	}
	if len(blocks) > 1 {
		return nil, fmt.Errorf("only one feature_settings block may be defined per app")
	}
	return &appengine.FeatureSettings{
		SplitHealthChecks: d.Get(prefix + "feature_settings.0.split_health_checks").(bool),
		// force send SplitHealthChecks, so if it's set to false it still gets disabled
		ForceSendFields: []string{"SplitHealthChecks"},
	}, nil
}

func flattenAppEngineFeatureSettings(settings *appengine.FeatureSettings) ([]map[string]interface{}, error) {
	if settings == nil {
		return []map[string]interface{}{}, nil
	}
	result := map[string]interface{}{
		"split_health_checks": settings.SplitHealthChecks,
	}
	return []map[string]interface{}{result}, nil
}

func flattenAppEngineDispatchRules(rules []*appengine.UrlDispatchRule) ([]map[string]interface{}, error) {
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
