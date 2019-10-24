package google

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Read:   schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		MigrateState: resourceGoogleProjectMigrateState,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateProjectID(),
			},
			"skip_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_create_network": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateProjectName(),
			},
			"org_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"folder_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: parseFolderId,
			},
			"policy_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Use the 'google_project_iam_policy' resource to define policies for a Google Project",
			},
			"policy_etag": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Use the the 'google_project_iam_policy' resource to define policies for a Google Project",
			},
			"number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"app_engine": {
				Type:     schema.TypeList,
				Elem:     appEngineResource(),
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
		},
	}
}

func appEngineResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auth_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"location_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"serving_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"feature_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
				Elem:     appEngineFeatureSettingsResource(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"url_dispatch_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
				Elem:     appEngineURLDispatchRuleResource(),
			},
			"code_bucket": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"default_hostname": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"default_bucket": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"gcr_domain": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
		},
	}
}

func appEngineURLDispatchRuleResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
		},
	}
}

func appEngineFeatureSettingsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"split_health_checks": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "This field has been removed. Use the google_app_engine_application resource instead.",
			},
		},
	}
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

	var op *cloudresourcemanager.Operation
	err = retryTimeDuration(func() (reqErr error) {
		op, reqErr = config.clientResourceManager.Projects.Create(project).Do()
		return reqErr
	}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating project %s (%s): %s. "+
			"If you received a 403 error, make sure you have the"+
			" `roles/resourcemanager.projectCreator` permission",
			project.ProjectId, project.Name, err)
	}

	d.SetId(pid)

	// Wait for the operation to complete
	opAsMap, err := ConvertToMap(op)
	if err != nil {
		return err
	}

	waitErr := resourceManagerOperationWaitTime(config, opAsMap, "creating folder", int(d.Timeout(schema.TimeoutCreate).Minutes()))
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
		if err = enableServiceUsageProjectServices([]string{"compute.googleapis.com"}, project.ProjectId, config, d.Timeout(schema.TimeoutCreate)); err != nil {
			return fmt.Errorf("Error enabling the Compute Engine API required to delete the default network: %s", err)
		}

		if err = forceDeleteComputeNetwork(d, config, project.ProjectId, "default"); err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG] Default network not found for project %q, no need to delete it", project.ProjectId)
			} else {
				return fmt.Errorf("Error deleting default network in project %s: %s", project.ProjectId, err)
			}
		}
	}
	return nil
}

func resourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	pid := d.Id()

	p, err := readGoogleProject(d, config)
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
	d.Set("number", strconv.FormatInt(p.ProjectNumber, 10))
	d.Set("name", p.Name)
	d.Set("labels", p.Labels)

	// We get app_engine.#: "" => "<computed>" without this set
	// Remove when app_engine field is removed from schema completely
	d.Set("app_engine", nil)

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

	var ba *cloudbilling.ProjectBillingInfo
	err = retryTimeDuration(func() (reqErr error) {
		ba, reqErr = config.clientBilling.Projects.GetBillingInfo(prefixedProject(pid)).Do()
		return reqErr
	}, d.Timeout(schema.TimeoutRead))
	// Read the billing account
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
	p, err := readGoogleProject(d, config)
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			return fmt.Errorf("Project %q does not exist.", pid)
		}
		return fmt.Errorf("Error checking project %q: %s", pid, err)
	}

	d.Partial(true)

	// Project display name has changed
	if ok := d.HasChange("name"); ok {
		p.Name = project_name
		// Do update on project
		if err = retryTimeDuration(func() (updateErr error) {
			p, updateErr = config.clientResourceManager.Projects.Update(p.ProjectId, p).Do()
			return updateErr
		}, d.Timeout(schema.TimeoutUpdate)); err != nil {
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
		if err = retryTimeDuration(func() (updateErr error) {
			p, updateErr = config.clientResourceManager.Projects.Update(p.ProjectId, p).Do()
			return updateErr
		}, d.Timeout(schema.TimeoutUpdate)); err != nil {
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
		if err = retryTimeDuration(func() (updateErr error) {
			p, updateErr = config.clientResourceManager.Projects.Update(p.ProjectId, p).Do()
			return updateErr
		}, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("Error updating project %q: %s", project_name, err)
		}
		d.SetPartial("labels")
	}

	d.Partial(false)
	return resourceGoogleProjectRead(d, meta)
}

func resourceGoogleProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	// Only delete projects if skip_delete isn't set
	if !d.Get("skip_delete").(bool) {
		pid := d.Id()
		if err := retryTimeDuration(func() error {
			_, delErr := config.clientResourceManager.Projects.Delete(pid).Do()
			return delErr
		}, d.Timeout(schema.TimeoutDelete)); err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Project %s", pid))
		}
	}
	d.SetId("")
	return nil
}

func resourceProjectImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	pid := d.Id()
	// Prevent importing via project number, this will cause issues later
	matched, err := regexp.MatchString("^\\d+$", pid)
	if err != nil {
		return nil, fmt.Errorf("Error matching project %q: %s", pid, err)
	}

	if matched {
		return nil, fmt.Errorf("Error importing project %q, please use project_id", pid)
	}

	// Explicitly set to default as a workaround for `ImportStateVerify` tests, and so that users
	// don't see a diff immediately after import.
	d.Set("auto_create_network", true)
	return []*schema.ResourceData{d}, nil
}

// Delete a compute network along with the firewall rules inside it.
func forceDeleteComputeNetwork(d *schema.ResourceData, config *Config, projectId, networkName string) error {
	networkLink, err := replaceVars(d, config, fmt.Sprintf("{{ComputeBasePath}}projects/%s/global/networks/%s", projectId, networkName))
	if err != nil {
		return err
	}

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
	ba := &cloudbilling.ProjectBillingInfo{}
	// If we're unlinking an existing billing account, an empty request does that, not an empty-string billing account.
	if name != "" {
		ba.BillingAccountName = "billingAccounts/" + name
	}
	_, err := config.clientBilling.Projects.UpdateBillingInfo(prefixedProject(pid), ba).Do()
	if err != nil {
		d.Set("billing_account", "")
		if _err, ok := err.(*googleapi.Error); ok {
			return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), _err)
		}
		return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), err)
	}
	for retries := 0; retries < 3; retries++ {
		ba, err = config.clientBilling.Projects.GetBillingInfo(prefixedProject(pid)).Do()
		if err != nil {
			return err
		}
		baName := strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/")
		if baName == name {
			return nil
		}
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("Timed out waiting for billing account to return correct value.  Waiting for %s, got %s.",
		name, strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/"))
}

func deleteComputeNetwork(project, network string, config *Config) error {
	op, err := config.clientCompute.Networks.Delete(
		project, network).Do()
	if err != nil {
		return fmt.Errorf("Error deleting network: %s", err)
	}

	err = computeOperationWaitTime(config.clientCompute, op, project, "Deleting Network", 10)
	if err != nil {
		return err
	}
	return nil
}

func readGoogleProject(d *schema.ResourceData, config *Config) (*cloudresourcemanager.Project, error) {
	var p *cloudresourcemanager.Project
	// Read the project
	err := retryTimeDuration(func() (reqErr error) {
		p, reqErr = config.clientResourceManager.Projects.Get(d.Id()).Do()
		return reqErr
	}, d.Timeout(schema.TimeoutRead))
	return p, err
}

// Enables services. WARNING: Use BatchRequestEnableServices for better batching if possible.
func enableServiceUsageProjectServices(services []string, project string, config *Config, timeout time.Duration) error {
	// ServiceUsage does not allow more than 20 services to be enabled per
	// batchEnable API call. See
	// https://cloud.google.com/service-usage/docs/reference/rest/v1/services/batchEnable
	for i := 0; i < len(services); i += maxServiceUsageBatchSize {
		j := i + maxServiceUsageBatchSize
		if j > len(services) {
			j = len(services)
		}
		nextBatch := services[i:j]
		if len(nextBatch) == 0 {
			// All batches finished, return.
			return nil
		}

		if err := doEnableServicesRequest(nextBatch, project, config, timeout); err != nil {
			return err
		}
		log.Printf("[DEBUG] Finished enabling next batch of %d project services: %+v", len(nextBatch), nextBatch)
	}

	log.Printf("[DEBUG] Verifying that all services are enabled")
	return waitForServiceUsageEnabledServices(services, project, config, timeout)
}

// waitForServiceUsageEnabledServices doesn't resend enable requests - it just
// waits for service enablement status to propagate. Essentially, it waits until
// all services show up as enabled when listing services on the project.
func waitForServiceUsageEnabledServices(services []string, project string, config *Config, timeout time.Duration) error {
	missing := make([]string, 0, len(services))
	delay := time.Duration(0)
	interval := time.Second
	err := retryTimeDuration(func() error {
		// Get the list of services that are enabled on the project
		enabledServices, err := listCurrentlyEnabledServices(project, config, timeout)
		if err != nil {
			return err
		}

		missing := make([]string, 0, len(services))
		for _, s := range services {
			if _, ok := enabledServices[s]; !ok {
				missing = append(missing, s)
			}
		}
		if len(missing) > 0 {
			log.Printf("[DEBUG] waiting %v before reading project %s services...", delay, project)
			time.Sleep(delay)
			delay += interval
			interval += delay

			// Spoof a googleapi Error so retryTime will try again
			return &googleapi.Error{
				Code:    503,
				Message: fmt.Sprintf("The service(s) %q are still being enabled for project %s. This isn't a real API error, this is just eventual consistency.", missing, project),
			}
		}
		return nil
	}, timeout)
	if err != nil {
		return errwrap.Wrap(err, fmt.Errorf("failed to enable some service(s) %q for project %s", missing, project))
	}
	return nil
}
