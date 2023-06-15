// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	tpgserviceusage "github.com/hashicorp/terraform-provider-google/google/services/serviceusage"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
)

type ServicesCall interface {
	Header() http.Header
	Do(opts ...googleapi.CallOption) (*serviceusage.Operation, error)
}

// ResourceGoogleProject returns a *schema.Resource that allows a customer
// to declare a Google Cloud Project resource.
func ResourceGoogleProject() *schema.Resource {
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
				ValidateFunc: verify.ValidateProjectID(),
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
				ValidateFunc: verify.ValidateProjectName(),
				Description:  `The display name of the project.`,
			},
			"org_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"folder_id"},
				Description:   `The numeric ID of the organization this project belongs to. Changing this forces a new project to be created.  Only one of org_id or folder_id may be specified. If the org_id is specified then the project is created at the top level. Changing this forces the project to be migrated to the newly specified organization.`,
			},
			"folder_id": {
				Type:          schema.TypeString,
				Optional:      true,
				StateFunc:     ParseFolderId,
				ConflictsWith: []string{"org_id"},
				Description:   `The numeric ID of the folder this project should be created under. Only one of org_id or folder_id may be specified. If the folder_id is specified, then the project is created under the specified folder. Changing this forces the project to be migrated to the newly specified folder.`,
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
		UseJSONNumber: true,
	}
}

func resourceGoogleProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	if err = resourceGoogleProjectCheckPreRequisites(config, d, userAgent); err != nil {
		return fmt.Errorf("failed pre-requisites: %v", err)
	}

	var pid string
	pid = d.Get("project_id").(string)

	log.Printf("[DEBUG]: Creating new project %q", pid)
	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      d.Get("name").(string),
	}

	if err = getParentResourceId(d, project); err != nil {
		return err
	}

	if _, ok := d.GetOk("labels"); ok {
		project.Labels = tpgresource.ExpandLabels(d)
	}

	var op *cloudresourcemanager.Operation
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			op, reqErr = config.NewResourceManagerClient(userAgent).Projects.Create(project).Do()
			return reqErr
		},
		Timeout: d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("error creating project %s (%s): %s. "+
			"If you received a 403 error, make sure you have the"+
			" `roles/resourcemanager.projectCreator` permission",
			project.ProjectId, project.Name, err)
	}

	d.SetId(fmt.Sprintf("projects/%s", pid))

	// Wait for the operation to complete
	opAsMap, err := tpgresource.ConvertToMap(op)
	if err != nil {
		return err
	}

	waitErr := ResourceManagerOperationWaitTime(config, opAsMap, "creating folder", userAgent, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		// The resource wasn't actually created
		d.SetId("")
		return waitErr
	}

	// Set the billing account
	if _, ok := d.GetOk("billing_account"); ok {
		err = updateProjectBillingAccount(d, config, userAgent)
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

		billingProject := project.ProjectId
		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}

		if err = EnableServiceUsageProjectServices([]string{"compute.googleapis.com"}, project.ProjectId, billingProject, userAgent, config, d.Timeout(schema.TimeoutCreate)); err != nil {
			return errwrap.Wrapf("Error enabling the Compute Engine API required to delete the default network: {{err}} ", err)
		}

		if err = forceDeleteComputeNetwork(d, config, project.ProjectId, "default"); err != nil {
			if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
				log.Printf("[DEBUG] Default network not found for project %q, no need to delete it", project.ProjectId)
			} else {
				return errwrap.Wrapf(fmt.Sprintf("Error deleting default network in project %s: {{err}}", project.ProjectId), err)
			}
		}
	}
	return nil
}

func resourceGoogleProjectCheckPreRequisites(config *transport_tpg.Config, d *schema.ResourceData, userAgent string) error {
	ib, ok := d.GetOk("billing_account")
	if !ok {
		return nil
	}
	ba := "billingAccounts/" + ib.(string)
	const perm = "billing.resourceAssociations.create"
	req := &cloudbilling.TestIamPermissionsRequest{
		Permissions: []string{perm},
	}
	resp, err := config.NewBillingClient(userAgent).BillingAccounts.TestIamPermissions(ba, req).Do()
	if err != nil {
		return fmt.Errorf("failed to check permissions on billing account %q: %v", ba, err)
	}
	if !tpgresource.StringInSlice(resp.Permissions, perm) {
		return fmt.Errorf("missing permission on %q: %v", ba, perm)
	}
	if !d.Get("auto_create_network").(bool) {
		call := config.NewServiceUsageClient(userAgent).Services.Get("projects/00000000000/services/serviceusage.googleapis.com")
		if config.UserProjectOverride {
			if billingProject, err := tpgresource.GetBillingProject(d, config); err == nil {
				call.Header().Add("X-Goog-User-Project", billingProject)
			}
		}
		_, err := call.Do()
		switch {
		// We are querying a dummy project since the call is already coming from the quota project.
		// If the API is enabled we get a not found message or accessNotConfigured if API is not enabled.
		case err.Error() == "googleapi: Error 403: Project '00000000000' not found or permission denied., forbidden":
			return nil
		case strings.Contains(err.Error(), "accessNotConfigured"):
			return fmt.Errorf("API serviceusage not enabled.\nFound error: %v", err)
		}
	}
	return nil
}

func resourceGoogleProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]

	p, err := readGoogleProject(d, config, userAgent)
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 403 && strings.Contains(gerr.Message, "caller does not have permission") {
			return fmt.Errorf("the user does not have permission to access Project %q or it may not exist", pid)
		}
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project %q", pid))
	}

	// If the project has been deleted from outside Terraform, remove it from state file.
	if p.LifecycleState != "ACTIVE" {
		log.Printf("[WARN] Removing project '%s' because its state is '%s' (requires 'ACTIVE').", pid, p.LifecycleState)
		d.SetId("")
		return nil
	}

	if err := d.Set("project_id", pid); err != nil {
		return fmt.Errorf("Error setting project_id: %s", err)
	}
	if err := d.Set("number", strconv.FormatInt(p.ProjectNumber, 10)); err != nil {
		return fmt.Errorf("Error setting number: %s", err)
	}
	if err := d.Set("name", p.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("labels", p.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}

	if p.Parent != nil {
		switch p.Parent.Type {
		case "organization":
			if err := d.Set("org_id", p.Parent.Id); err != nil {
				return fmt.Errorf("Error setting org_id: %s", err)
			}
			if err := d.Set("folder_id", ""); err != nil {
				return fmt.Errorf("Error setting folder_id: %s", err)
			}
		case "folder":
			if err := d.Set("folder_id", p.Parent.Id); err != nil {
				return fmt.Errorf("Error setting folder_id: %s", err)
			}
			if err := d.Set("org_id", ""); err != nil {
				return fmt.Errorf("Error setting org_id: %s", err)
			}
		}
	}

	var ba *cloudbilling.ProjectBillingInfo
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			ba, reqErr = config.NewBillingClient(userAgent).Projects.GetBillingInfo(PrefixedProject(pid)).Do()
			return reqErr
		},
		Timeout: d.Timeout(schema.TimeoutRead),
	})
	// Read the billing account
	if err != nil && !transport_tpg.IsApiNotEnabledError(err) {
		return fmt.Errorf("Error reading billing account for project %q: %v", PrefixedProject(pid), err)
	} else if transport_tpg.IsApiNotEnabledError(err) {
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
			return fmt.Errorf("Error parsing billing account for project %q. Expected value to begin with 'billingAccounts/' but got %s", PrefixedProject(pid), ba.BillingAccountName)
		}
		if err := d.Set("billing_account", _ba); err != nil {
			return fmt.Errorf("Error setting billing_account: %s", err)
		}
	}

	return nil
}

func PrefixedProject(pid string) string {
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
			Id:   ParseFolderId(folderId),
			Type: "folder",
		}
	}

	return nil
}

func ParseFolderId(v interface{}) string {
	folderId := v.(string)
	if strings.HasPrefix(folderId, "folders/") {
		return folderId[8:]
	}
	return folderId
}

func resourceGoogleProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	project_name := d.Get("name").(string)

	// Read the project
	// we need the project even though refresh has already been called
	// because the API doesn't support patch, so we need the actual object
	p, err := readGoogleProject(d, config, userAgent)
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return fmt.Errorf("Project %q does not exist.", pid)
		}
		return fmt.Errorf("Error checking project %q: %s", pid, err)
	}

	d.Partial(true)

	// Project display name has changed
	if ok := d.HasChange("name"); ok {
		p.Name = project_name
		// Do update on project
		if p, err = updateProject(config, d, project_name, userAgent, p); err != nil {
			return err
		}
	}

	// Project parent has changed
	if d.HasChange("org_id") || d.HasChange("folder_id") {
		if err := getParentResourceId(d, p); err != nil {
			return err
		}

		// Do update on project
		if p, err = updateProject(config, d, project_name, userAgent, p); err != nil {
			return err
		}
	}

	// Billing account has changed
	if ok := d.HasChange("billing_account"); ok {
		err = updateProjectBillingAccount(d, config, userAgent)
		if err != nil {
			return err
		}
	}

	// Project Labels have changed
	if ok := d.HasChange("labels"); ok {
		p.Labels = tpgresource.ExpandLabels(d)

		// Do Update on project
		if p, err = updateProject(config, d, project_name, userAgent, p); err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceGoogleProjectRead(d, meta)
}

func updateProject(config *transport_tpg.Config, d *schema.ResourceData, projectName, userAgent string, desiredProject *cloudresourcemanager.Project) (*cloudresourcemanager.Project, error) {
	var newProj *cloudresourcemanager.Project
	if err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (updateErr error) {
			newProj, updateErr = config.NewResourceManagerClient(userAgent).Projects.Update(desiredProject.ProjectId, desiredProject).Do()
			return updateErr
		},
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}); err != nil {
		return nil, fmt.Errorf("Error updating project %q: %s", projectName, err)
	}
	return newProj, nil
}

func resourceGoogleProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	// Only delete projects if skip_delete isn't set
	if !d.Get("skip_delete").(bool) {
		parts := strings.Split(d.Id(), "/")
		pid := parts[len(parts)-1]
		if err := transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() error {
				_, delErr := config.NewResourceManagerClient(userAgent).Projects.Delete(pid).Do()
				return delErr
			},
			Timeout: d.Timeout(schema.TimeoutDelete),
		}); err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project %s", pid))
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
	if err := d.Set("auto_create_network", true); err != nil {
		return nil, fmt.Errorf("Error setting auto_create_network: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

// Delete a compute network along with the firewall rules inside it.
func forceDeleteComputeNetwork(d *schema.ResourceData, config *transport_tpg.Config, projectId, networkName string) error {
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	// Read the network from the API so we can get the correct self link format. We can't construct it from the
	// base path because it might not line up exactly (compute.googleapis.com vs www.googleapis.com)
	net, err := config.NewComputeClient(userAgent).Networks.Get(projectId, networkName).Do()
	if err != nil {
		return err
	}

	token := ""
	for paginate := true; paginate; {
		filter := fmt.Sprintf("network eq %s", net.SelfLink)
		resp, err := config.NewComputeClient(userAgent).Firewalls.List(projectId).Filter(filter).Do()
		if err != nil {
			return errwrap.Wrapf("Error listing firewall rules in proj: {{err}}", err)
		}

		log.Printf("[DEBUG] Found %d firewall rules in %q network", len(resp.Items), networkName)

		for _, firewall := range resp.Items {
			op, err := config.NewComputeClient(userAgent).Firewalls.Delete(projectId, firewall.Name).Do()
			if err != nil {
				return errwrap.Wrapf("Error deleting firewall: {{err}}", err)
			}
			err = tpgcompute.ComputeOperationWaitTime(config, op, projectId, "Deleting Firewall", userAgent, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return err
			}
		}

		token = resp.NextPageToken
		paginate = token != ""
	}

	return deleteComputeNetwork(projectId, networkName, userAgent, config)
}

func updateProjectBillingAccount(d *schema.ResourceData, config *transport_tpg.Config, userAgent string) error {
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	name := d.Get("billing_account").(string)
	ba := &cloudbilling.ProjectBillingInfo{}
	// If we're unlinking an existing billing account, an empty request does that, not an empty-string billing account.
	if name != "" {
		ba.BillingAccountName = "billingAccounts/" + name
	}
	updateBillingInfoFunc := func() error {
		_, err := config.NewBillingClient(userAgent).Projects.UpdateBillingInfo(PrefixedProject(pid), ba).Do()
		return err
	}
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: updateBillingInfoFunc,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
	})
	if err != nil {
		if err := d.Set("billing_account", ""); err != nil {
			return fmt.Errorf("Error setting billing_account: %s", err)
		}
		if _err, ok := err.(*googleapi.Error); ok {
			return fmt.Errorf("Error setting billing account %q for project %q: %v", name, PrefixedProject(pid), _err)
		}
		return fmt.Errorf("Error setting billing account %q for project %q: %v", name, PrefixedProject(pid), err)
	}
	for retries := 0; retries < 3; retries++ {
		var ba *cloudbilling.ProjectBillingInfo
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (reqErr error) {
				ba, reqErr = config.NewBillingClient(userAgent).Projects.GetBillingInfo(PrefixedProject(pid)).Do()
				return reqErr
			},
			Timeout: d.Timeout(schema.TimeoutRead),
		})
		if err != nil {
			return fmt.Errorf("Error getting billing info for project %q: %v", PrefixedProject(pid), err)
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

func deleteComputeNetwork(project, network, userAgent string, config *transport_tpg.Config) error {
	op, err := config.NewComputeClient(userAgent).Networks.Delete(
		project, network).Do()
	if err != nil {
		return errwrap.Wrapf("Error deleting network: {{err}}", err)
	}

	err = tpgcompute.ComputeOperationWaitTime(config, op, project, "Deleting Network", userAgent, 10*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func readGoogleProject(d *schema.ResourceData, config *transport_tpg.Config, userAgent string) (*cloudresourcemanager.Project, error) {
	var p *cloudresourcemanager.Project
	// Read the project
	parts := strings.Split(d.Id(), "/")
	pid := parts[len(parts)-1]
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			p, reqErr = config.NewResourceManagerClient(userAgent).Projects.Get(pid).Do()
			return reqErr
		},
		Timeout: d.Timeout(schema.TimeoutRead),
	})
	return p, err
}

// Enables services. WARNING: Use BatchRequestEnableServices for better batching if possible.
func EnableServiceUsageProjectServices(services []string, project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) error {
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

		if err := doEnableServicesRequest(nextBatch, project, billingProject, userAgent, config, timeout); err != nil {
			return err
		}
		log.Printf("[DEBUG] Finished enabling next batch of %d project services: %+v", len(nextBatch), nextBatch)
	}

	log.Printf("[DEBUG] Verifying that all services are enabled")
	return waitForServiceUsageEnabledServices(services, project, billingProject, userAgent, config, timeout)
}

func doEnableServicesRequest(services []string, project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) error {
	var op *serviceusage.Operation

	// errors can come up at multiple points, so there are a few levels of
	// retrying here.
	// logicalErr / waitErr: overall error on the logical operation (enabling services)
	// but possibly also errors when retrieving the LRO (these are rare)
	// err / reqErr: precondition errors when sending the request received instead of an LRO
	logicalErr := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			err := transport_tpg.Retry(transport_tpg.RetryOptions{
				RetryFunc: func() error {
					var reqErr error
					var call ServicesCall
					if len(services) == 1 {
						// BatchEnable returns an error for a single item, so enable with single endpoint
						name := fmt.Sprintf("projects/%s/services/%s", project, services[0])
						req := &serviceusage.EnableServiceRequest{}
						call = config.NewServiceUsageClient(userAgent).Services.Enable(name, req)
					} else {
						// Batch enable for multiple services.
						name := fmt.Sprintf("projects/%s", project)
						req := &serviceusage.BatchEnableServicesRequest{ServiceIds: services}
						call = config.NewServiceUsageClient(userAgent).Services.BatchEnable(name, req)
					}

					if config.UserProjectOverride && billingProject != "" {
						call.Header().Add("X-Goog-User-Project", billingProject)
					}

					op, reqErr = call.Do()
					return handleServiceUsageRetryablePreconditionError(reqErr)
				},
				Timeout:              timeout,
				ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.ServiceUsageServiceBeingActivated},
			})
			if err != nil {
				return errwrap.Wrapf("failed on request preconditions: {{err}}", err)
			}

			waitErr := tpgserviceusage.ServiceUsageOperationWait(config, op, billingProject, fmt.Sprintf("Enable Project %q Services: %+v", project, services), userAgent, timeout)
			if waitErr != nil {
				return waitErr
			}

			return nil
		},
		Timeout:              timeout,
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.ServiceUsageInternalError160009},
	})

	if logicalErr != nil {
		return errwrap.Wrapf("failed to enable services: {{err}}", logicalErr)
	}

	return nil
}

// Handle errors that are retryable at call time for serviceusage
// Specifically, errors in https://cloud.google.com/service-usage/docs/reference/rest/v1/services/batchEnable#response-body
// Errors in operations are handled separately.
// NOTE(rileykarson): This should probably be turned into a retry predicate
func handleServiceUsageRetryablePreconditionError(err error) error {
	if err == nil {
		return nil
	}
	if gerr, ok := err.(*googleapi.Error); ok {
		if (gerr.Code == 400 || gerr.Code == 412) && gerr.Message == "Precondition check failed." {
			return &googleapi.Error{
				Code:    503,
				Message: "api returned \"precondition failed\" while enabling service",
			}
		}
	}
	return err
}

// Retrieve a project's services from the API
// if a service has been renamed, this function will list both the old and new
// forms of the service. LIST responses are expected to return only the old or
// new form, but we'll always return both.
func ListCurrentlyEnabledServices(project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) (map[string]struct{}, error) {
	log.Printf("[DEBUG] Listing enabled services for project %s", project)
	apiServices := make(map[string]struct{})
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			ctx := context.Background()
			call := config.NewServiceUsageClient(userAgent).Services.List(fmt.Sprintf("projects/%s", project))
			if config.UserProjectOverride && billingProject != "" {
				call.Header().Add("X-Goog-User-Project", billingProject)
			}
			return call.Fields("services/name,nextPageToken").Filter("state:ENABLED").
				Pages(ctx, func(r *serviceusage.ListServicesResponse) error {
					for _, v := range r.Services {
						// services are returned as "projects/{{project}}/services/{{name}}"
						name := tpgresource.GetResourceNameFromSelfLink(v.Name)

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
		},
		Timeout: timeout,
	})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Failed to list enabled services for project %s: {{err}}", project), err)
	}
	return apiServices, nil
}

// waitForServiceUsageEnabledServices doesn't resend enable requests - it just
// waits for service enablement status to propagate. Essentially, it waits until
// all services show up as enabled when listing services on the project.
func waitForServiceUsageEnabledServices(services []string, project, billingProject, userAgent string, config *transport_tpg.Config, timeout time.Duration) error {
	missing := make([]string, 0, len(services))
	delay := time.Duration(0)
	interval := time.Second
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			// Get the list of services that are enabled on the project
			enabledServices, err := ListCurrentlyEnabledServices(project, billingProject, userAgent, config, timeout)
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
		},
		Timeout: timeout,
	})
	if err != nil {
		return errwrap.Wrap(err, fmt.Errorf("failed to enable some service(s) %q for project %s", missing, project))
	}
	return nil
}
