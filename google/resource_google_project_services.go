package google

import (
	"context"
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
	"log"
	"strings"
	"time"
)

const maxServiceUsageBatchSize = 20

func resourceGoogleProjectServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectServicesCreateUpdate,
		Read:   resourceGoogleProjectServicesRead,
		Update: resourceGoogleProjectServicesCreateUpdate,
		Delete: resourceGoogleProjectServicesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "google_project_services is deprecated - many users reported " +
			"issues with dependent services that were not resolvable.  Please use google_project_service or the " +
			"https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/project_services" +
			" module.  It's recommended that you use a provider version of 2.13.0 or higher when you migrate so that requests are" +
			" batched to the API, reducing the request rate. This resource will be removed in version 3.0.0.",

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

	// use old services to set the correct renamed service names in state
	s, _ := expandServiceUsageProjectServicesServices(d.Get("services"), d, config)
	log.Printf("[DEBUG] Saw services in state on Read: %s ", s)
	sset := golangSetFromStringSlice(s)
	for ov, nv := range renamedServices {
		_, ook := sset[ov]
		_, nok := sset[nv]

		// preserve the values set in prior state if they're identical. If none
		// were set, we  delete the new value if it exists. By doing that that
		// we only store the old value if the service is enabled, and no value
		// if it isn't.
		if ook && nok {
			continue
		} else if ook {
			delete(enabledSet, nv)
		} else if nok {
			delete(enabledSet, ov)
		} else {
			delete(enabledSet, nv)
		}
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

	toEnable := map[string]struct{}{}
	for _, srv := range services {
		// We don't have to enable a service if it's already enabled.
		if _, ok := currentlyEnabled[srv]; !ok {
			toEnable[srv] = struct{}{}
		}
	}

	if len(toEnable) > 0 {
		log.Printf("[DEBUG] Enabling services: %s", toEnable)
		if err := BatchRequestEnableServices(toEnable, project, d, config); err != nil {
			return fmt.Errorf("unable to enable Project Services %s (%+v): %s", project, services, err)
		}
	} else {
		log.Printf("[DEBUG] No services to enable.")
	}

	srvSet := golangSetFromStringSlice(services)

	srvSetWithRenames := map[string]struct{}{}

	// we'll always list both names for renamed services, so allow both forms if
	// we see both.
	for k := range srvSet {
		srvSetWithRenames[k] = struct{}{}
		if v, ok := renamedServicesByOldAndNewServiceNames[k]; ok {
			srvSetWithRenames[v] = struct{}{}
		}
	}

	for srv := range currentlyEnabled {
		// Disable any services that are currently enabled for project but are not
		// in our list of acceptable services.
		if _, ok := srvSetWithRenames[srv]; !ok {
			// skip deleting services by their new names and prefer the old name.
			if _, ok := renamedServicesByNewServiceNames[srv]; ok {
				continue
			}

			log.Printf("[DEBUG] Disabling project %s service %s", project, srv)
			err := disableServiceUsageProjectService(srv, project, d, config, true)
			if err != nil {
				log.Printf("[DEBUG] Saw error %s deleting service %s", err, srv)

				// if we got the right error and the service is renamed, delete by the new name
				if n, ok := renamedServices[srv]; ok && strings.Contains(err.Error(), "not found or permission denied.") {
					log.Printf("[DEBUG] Failed to delete service %s, it doesn't exist. Trying %s", srv, n)
					err = disableServiceUsageProjectService(n, project, d, config, true)
					if err == nil {
						return nil
					}
				}

				return fmt.Errorf("unable to disable unwanted Project Service %s %s): %s", project, srv, err)
			}
		}
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

// Retrieve a project's services from the API
// if a service has been renamed, this function will list both the old and new
// forms of the service. LIST responses are expected to return only the old or
// new form, but we'll always return both.
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
