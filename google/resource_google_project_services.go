package google

import (
	"context"
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
	"log"
	"time"
)

const maxServiceUsageBatchSize = 20

var ignoredProjectServices = []string{"dataproc-control.googleapis.com", "source.googleapis.com", "stackdriverprovisioning.googleapis.com"}

// These services can only be enabled as a side-effect of enabling other services,
// so don't bother storing them in the config or using them for diffing.
var ignoredProjectServicesSet = golangSetFromStringSlice(ignoredProjectServices)

func resourceGoogleProjectServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectServicesCreateUpdate,
		Read:   resourceGoogleProjectServicesRead,
		Update: resourceGoogleProjectServicesCreateUpdate,
		Delete: resourceGoogleProjectServicesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"services": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: StringNotInSlice(ignoredProjectServices, false),
				},
			},
			"disable_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGoogleProjectServicesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get services from config
	services, err := expandServiceUsageProjectServicesServices(d.Get("services"), d, config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]: Enabling Project Services for %s: %+v", d.Id(), services)
	if err := setServiceUsageProjectEnabledServices(services, project, d, config); err != nil {
		return fmt.Errorf("Error authoritatively enabling Project %s Services: %v", project, err)
	}
	log.Printf("[DEBUG]: Finished enabling Project Services for %s: %+v", d.Id(), services)

	d.SetId(project)
	return resourceGoogleProjectServicesRead(d, meta)
}

func resourceGoogleProjectServicesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	enabledSet, err := listCurrentlyEnabledServices(d.Id(), config, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return err
	}
	services := stringSliceFromGolangSet(enabledSet)

	d.Set("project", d.Id())
	d.Set("services", flattenServiceUsageProjectServicesServices(services, d))
	return nil
}

func resourceGoogleProjectServicesDelete(d *schema.ResourceData, meta interface{}) error {
	if disable := d.Get("disable_on_destroy"); !(disable.(bool)) {
		log.Printf("[WARN] Project Services disable_on_destroy set to false, skip disabling services for %s.", d.Id())
		d.SetId("")
		return nil
	}

	config := meta.(*Config)

	// Get services from config
	services, err := expandServiceUsageProjectServicesServices(d.Get("services"), d, config)
	if err != nil {
		return err
	}
	project := d.Id()

	log.Printf("[DEBUG]: Disabling Project Services %s: %+v", project, services)
	for _, s := range services {
		if err := disableServiceUsageProjectService(s, project, d, config, true); err != nil {
			return fmt.Errorf("Unable to destroy google_project_services for %s: %s", d.Id(), err)
		}
	}
	log.Printf("[DEBUG] Finished disabling Project Services %s: %+v", project, services)

	d.SetId("")
	return nil
}

// *Authoritatively* sets enabled services.
func setServiceUsageProjectEnabledServices(services []string, project string, d *schema.ResourceData, config *Config) error {
	currentlyEnabled, err := listCurrentlyEnabledServices(project, config, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return err
	}

	toEnable := make([]string, 0, len(services))
	for _, srv := range services {
		// We don't have to enable a service if it's already enabled.
		if _, ok := currentlyEnabled[srv]; !ok {
			toEnable = append(toEnable, srv)
		}
	}

	if err := BatchRequestEnableServices(toEnable, project, d, config); err != nil {
		return fmt.Errorf("unable to enable Project Services %s (%+v): %s", project, services, err)
	}

	srvSet := golangSetFromStringSlice(services)
	for srv := range currentlyEnabled {
		// Disable any services that are currently enabled for project but are not
		// in our list of acceptable services.
		if _, ok := srvSet[srv]; !ok {
			log.Printf("[DEBUG] Disabling project %s service %s", project, srv)
			if err := disableServiceUsageProjectService(srv, project, d, config, true); err != nil {
				return fmt.Errorf("unable to disable unwanted Project Service %s %s): %s", project, srv, err)
			}
		}
	}
	return nil
}

// Disables a project service.
func disableServiceUsageProjectService(service, project string, d *schema.ResourceData, config *Config, disableDependentServices bool) error {
	err := retryTimeDuration(func() error {
		name := fmt.Sprintf("projects/%s/services/%s", project, service)
		sop, err := config.clientServiceUsage.Services.Disable(name, &serviceusage.DisableServiceRequest{
			DisableDependentServices: disableDependentServices,
		}).Do()
		if err != nil {
			return err
		}
		// Wait for the operation to complete
		waitErr := serviceUsageOperationWait(config, sop, "api to disable")
		if waitErr != nil {
			return waitErr
		}
		return nil
	}, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error disabling service %q for project %q: %v", service, project, err)
	}
	return nil
}

// Retrieve a project's services from the API
func listCurrentlyEnabledServices(project string, config *Config, timeout time.Duration) (map[string]struct{}, error) {
	// Verify project for services still exists
	p, err := config.clientResourceManager.Projects.Get(project).Do()
	if err != nil {
		return nil, err
	}
	if p.LifecycleState == "DELETE_REQUESTED" {
		// Construct a 404 error for handleNotFoundError
		return nil, &googleapi.Error{
			Code:    404,
			Message: "Project deletion was requested",
		}
	}

	log.Printf("[DEBUG] Listing enabled services for project %s", project)
	apiServices := make(map[string]struct{})
	err = retryTimeDuration(func() error {
		ctx := context.Background()
		return config.clientServiceUsage.Services.
			List(fmt.Sprintf("projects/%s", project)).
			Fields("services/name,nextPageToken").
			Filter("state:ENABLED").
			Pages(ctx, func(r *serviceusage.ListServicesResponse) error {
				for _, v := range r.Services {
					// services are returned as "projects/PROJECT/services/NAME"
					name := GetResourceNameFromSelfLink(v.Name)
					if _, ok := ignoredProjectServicesSet[name]; !ok {
						apiServices[name] = struct{}{}
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
	waitErr := serviceUsageOperationWait(config, op, fmt.Sprintf("Enable Project %q Services: %+v", project, services))
	if waitErr != nil {
		return waitErr
	}
	return nil
}

func handleServiceUsageRetryableError(err error) error {
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

func flattenServiceUsageProjectServicesServices(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	if strV, ok := v.([]string); ok {
		v = convertStringArrToInterface(strV)
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func expandServiceUsageProjectServicesServices(v interface{}, d TerraformResourceData, config *Config) ([]string, error) {
	if v == nil {
		return nil, nil
	}
	return convertStringArr(v.(*schema.Set).List()), nil
}
