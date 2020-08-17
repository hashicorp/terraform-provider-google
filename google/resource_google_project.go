package google

import (
	"context"
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
	"google.golang.org/api/serviceusage/v1"
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
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		MigrateState: resourceGoogleProjectMigrateState,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateProjectID(),
				Description:  `The project ID. Changing this forces a new project to be created.`,
			},
			"skip_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `If true, the Terraform resource can be deleted without deleting the Project via the Google API.`,
			},
			"auto_create_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Create the 'default' network automatically.  Default true. If set to false, the default network will be deleted.  Note that, for quota purposes, you will still need to have 1 network slot available to create the project successfully, even if you set auto_create_network to false, since the network will exist momentarily.`,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateProjectName(),
				Description:  `The display name of the project.`,
			},
			"org_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The numeric ID of the organization this project belongs to. Changing this forces a new project to be created.  Only one of org_id or folder_id may be specified. If the org_id is specified then the project is created at the top level. Changing this forces the project to be migrated to the newly specified organization.`,
			},
			"folder_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				StateFunc:   parseFolderId,
				Description: `The numeric ID of the folder this project should be created under. Only one of org_id or folder_id may be specified. If the folder_id is specified, then the project is created under the specified folder. Changing this forces the project to be migrated to the newly specified folder.`,
			},
			"number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The numeric identifier of the project.`,
			},
			"billing_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The alphanumeric ID of the billing account this project belongs to. The user or service account performing this operation with Terraform must have Billing Account Administrator privileges (roles/billing.admin) in the organization. See Google Cloud Billing API Access Control for more details.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A set of key/value label pairs to assign to the project.`,
			},
		},
	}
}

func resourceGoogleProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceGoogleProjectCheckPreRequisites(config, d); err != nil {
		return fmt.Errorf("failed pre-requisites: %v", err)
	}

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

	d.SetId(fmt.Sprintf("projects/%s", pid))

	// Wait for the operation to complete
	opAsMap, err := ConvertToMap(op)
	if err != nil {
		return err
	}

	waitErr := resourceManagerOperationWaitTime(config, opAsMap, "creating folder", d.Timeout(schema.TimeoutCreate))
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

	// Sleep for 10s, letting the billing account settle before other resources
	// try to use this project.
	time.Sleep(10 * time.Second)

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
			return errwrap.Wrapf("Error enabling the Compute Engine API required to delete the default network: {{err}} ", err)
		}

		if err = forceDeleteComputeNetwork(d, config, project.ProjectId, "default"); err != nil {
			if isGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG] Default network not found for project %q, no need to delete it", project.ProjectId)
			} else {
				return errwrap.Wrapf(fmt.Sprintf("Error deleting default network in project %s: {{err}}", project.ProjectId), err)
			}
		}
	}
	return nil
}

func resourceGoogleProjectCheckPreRequisites(config *Config, d *schema.ResourceData) error {
	ib, ok := d.GetOk("billing_account")
	if !ok {
		return nil
	}
	ba := "billingAccounts/" + ib.(string)
	const perm = "billing.resourceAssociations.create"
	req := &cloudbilling.TestIamPermissionsRequest{
		Permissions: []string{perm},
	}
	resp, err := config.clientBilling.BillingAccounts.TestIamPermissions(ba, req).Do()
	if err != nil {
		return fmt.Errorf("failed to check permissions on billing account %q: %v", ba, err)
	}
	if !stringInSlice(resp.Permissions, perm) {
		return fmt.Errorf("missing permission on %q: %v", ba, perm)
	}
	return nil
}

func resourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]

	p, err := readGoogleProject(d, config)
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 403 && strings.Contains(gerr.Message, "caller does not have permission") {
			return fmt.Errorf("the user does not have permission to access Project %q or it may not exist", pid)
		}
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
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
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
		if p, err = updateProject(config, d, project_name, p); err != nil {
			return err
		}

		d.SetPartial("name")
	}

	// Project parent has changed
	if d.HasChange("org_id") || d.HasChange("folder_id") {
		if err := getParentResourceId(d, p); err != nil {
			return err
		}

		// Do update on project
		if p, err = updateProject(config, d, project_name, p); err != nil {
			return err
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
		if p, err = updateProject(config, d, project_name, p); err != nil {
			return err
		}
		d.SetPartial("labels")
	}

	d.Partial(false)
	return resourceGoogleProjectRead(d, meta)
}

func updateProject(config *Config, d *schema.ResourceData, projectName string, desiredProject *cloudresourcemanager.Project) (*cloudresourcemanager.Project, error) {
	var newProj *cloudresourcemanager.Project
	if err := retryTimeDuration(func() (updateErr error) {
		newProj, updateErr = config.clientResourceManager.Projects.Update(desiredProject.ProjectId, desiredProject).Do()
		return updateErr
	}, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return nil, fmt.Errorf("Error updating project %q: %s", projectName, err)
	}
	return newProj, nil
}

func resourceGoogleProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	// Only delete projects if skip_delete isn't set
	if !d.Get("skip_delete").(bool) {
		parts := strings.Split(d.Id(), "/")
		pid := parts[len(parts)-1]
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
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	// Prevent importing via project number, this will cause issues later
	matched, err := regexp.MatchString("^\\d+$", pid)
	if err != nil {
		return nil, fmt.Errorf("Error matching project %q: %s", pid, err)
	}

	if matched {
		return nil, fmt.Errorf("Error importing project %q, please use project_id", pid)
	}

	// Ensure the id format includes projects/
	d.SetId(fmt.Sprintf("projects/%s", pid))

	// Explicitly set to default as a workaround for `ImportStateVerify` tests, and so that users
	// don't see a diff immediately after import.
	d.Set("auto_create_network", true)
	return []*schema.ResourceData{d}, nil
}

// Delete a compute network along with the firewall rules inside it.
func forceDeleteComputeNetwork(d *schema.ResourceData, config *Config, projectId, networkName string) error {
	// Read the network from the API so we can get the correct self link format. We can't construct it from the
	// base path because it might not line up exactly (compute.googleapis.com vs www.googleapis.com)
	net, err := config.clientCompute.Networks.Get(projectId, networkName).Do()
	if err != nil {
		return err
	}

	token := ""
	for paginate := true; paginate; {
		filter := fmt.Sprintf("network eq %s", net.SelfLink)
		resp, err := config.clientCompute.Firewalls.List(projectId).Filter(filter).Do()
		if err != nil {
			return errwrap.Wrapf("Error listing firewall rules in proj: {{err}}", err)
		}

		log.Printf("[DEBUG] Found %d firewall rules in %q network", len(resp.Items), networkName)

		for _, firewall := range resp.Items {
			op, err := config.clientCompute.Firewalls.Delete(projectId, firewall.Name).Do()
			if err != nil {
				return errwrap.Wrapf("Error deleting firewall: {{err}}", err)
			}
			err = computeOperationWaitTime(config, op, projectId, "Deleting Firewall", d.Timeout(schema.TimeoutCreate))
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
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	name := d.Get("billing_account").(string)
	ba := &cloudbilling.ProjectBillingInfo{}
	// If we're unlinking an existing billing account, an empty request does that, not an empty-string billing account.
	if name != "" {
		ba.BillingAccountName = "billingAccounts/" + name
	}
	updateBillingInfoFunc := func() error {
		_, err := config.clientBilling.Projects.UpdateBillingInfo(prefixedProject(pid), ba).Do()
		return err
	}
	err := retryTimeDuration(updateBillingInfoFunc, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		d.Set("billing_account", "")
		if _err, ok := err.(*googleapi.Error); ok {
			return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), _err)
		}
		return fmt.Errorf("Error setting billing account %q for project %q: %v", name, prefixedProject(pid), err)
	}
	for retries := 0; retries < 3; retries++ {
		var ba *cloudbilling.ProjectBillingInfo
		err = retryTimeDuration(func() (reqErr error) {
			ba, reqErr = config.clientBilling.Projects.GetBillingInfo(prefixedProject(pid)).Do()
			return reqErr
		}, d.Timeout(schema.TimeoutRead))
		if err != nil {
			return fmt.Errorf("Error getting billing info for project %q: %v", prefixedProject(pid), err)
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
		return errwrap.Wrapf("Error deleting network: {{err}}", err)
	}

	err = computeOperationWaitTime(config, op, project, "Deleting Network", 10*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func readGoogleProject(d *schema.ResourceData, config *Config) (*cloudresourcemanager.Project, error) {
	var p *cloudresourcemanager.Project
	// Read the project
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	err := retryTimeDuration(func() (reqErr error) {
		p, reqErr = config.clientResourceManager.Projects.Get(pid).Do()
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

func doEnableServicesRequest(services []string, project string, config *Config, timeout time.Duration) error {
	var op *serviceusage.Operation

	err := retryTimeDuration(func() error {
		var rerr error
		if len(services) == 1 {
			// BatchEnable returns an error for a single item, so just enable
			// using service endpoint.
			name := fmt.Sprintf("projects/%s/services/%s", project, services[0])
			req := &serviceusage.EnableServiceRequest{}
			op, rerr = config.clientServiceUsage.Services.Enable(name, req).Do()
		} else {
			// Batch enable for multiple services.
			name := fmt.Sprintf("projects/%s", project)
			req := &serviceusage.BatchEnableServicesRequest{ServiceIds: services}
			op, rerr = config.clientServiceUsage.Services.BatchEnable(name, req).Do()
		}
		return handleServiceUsageRetryableError(rerr)
	}, timeout)
	if err != nil {
		return errwrap.Wrapf("failed to send enable services request: {{err}}", err)
	}
	// Poll for the API to return
	waitErr := serviceUsageOperationWait(config, op, project, fmt.Sprintf("Enable Project %q Services: %+v", project, services), timeout)
	if waitErr != nil {
		return waitErr
	}
	return nil
}

// Retrieve a project's services from the API
// if a service has been renamed, this function will list both the old and new
// forms of the service. LIST responses are expected to return only the old or
// new form, but we'll always return both.
func listCurrentlyEnabledServices(project string, config *Config, timeout time.Duration) (map[string]struct{}, error) {
	log.Printf("[DEBUG] Listing enabled services for project %s", project)
	apiServices := make(map[string]struct{})
	err := retryTimeDuration(func() error {
		ctx := context.Background()
		return config.clientServiceUsage.Services.
			List(fmt.Sprintf("projects/%s", project)).
			Fields("services/name,nextPageToken").
			Filter("state:ENABLED").
			Pages(ctx, func(r *serviceusage.ListServicesResponse) error {
				for _, v := range r.Services {
					// services are returned as "projects/{{project}}/services/{{name}}"
					name := GetResourceNameFromSelfLink(v.Name)

					// if name not in ignoredProjectServicesSet
					if _, ok := ignoredProjectServicesSet[name]; !ok {
						apiServices[name] = struct{}{}

						// if a service has been renamed, set both. We'll deal
						// with setting the right values later.
						if v, ok := renamedServicesByOldAndNewServiceNames[name]; ok {
							log.Printf("[DEBUG] Adding service alias for %s to enabled services: %s", name, v)
							apiServices[v] = struct{}{}
						}
					}
				}
				return nil
			})
	}, timeout)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Failed to list enabled services for project %s: {{err}}", project), err)
	}
	return apiServices, nil
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
