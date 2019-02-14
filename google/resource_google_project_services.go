package google

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
)

func resourceGoogleProjectServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectServicesCreate,
		Read:   resourceGoogleProjectServicesRead,
		Update: resourceGoogleProjectServicesUpdate,
		Delete: resourceGoogleProjectServicesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

var ignoredProjectServices = []string{"dataproc-control.googleapis.com", "source.googleapis.com", "stackdriverprovisioning.googleapis.com"}

// These services can only be enabled as a side-effect of enabling other services,
// so don't bother storing them in the config or using them for diffing.
var ignoreProjectServices = golangSetFromStringSlice(ignoredProjectServices)

func resourceGoogleProjectServicesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	pid, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get services from config
	cfgServices := getConfigServices(d)

	// Get services from API
	apiServices, err := getApiServices(pid, config, ignoreProjectServices)
	if err != nil {
		return fmt.Errorf("Error creating services: %v", err)
	}

	// This call disables any APIs that aren't defined in cfgServices,
	// and enables all of those that are
	err = reconcileServices(cfgServices, apiServices, config, pid)
	if err != nil {
		return fmt.Errorf("Error creating services: %v", err)
	}

	d.SetId(pid)
	return resourceGoogleProjectServicesRead(d, meta)
}

func resourceGoogleProjectServicesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	services, err := getApiServices(d.Id(), config, ignoreProjectServices)
	if err != nil {
		return err
	}

	d.Set("project", d.Id())
	d.Set("services", services)
	return nil
}

func resourceGoogleProjectServicesUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Updating google_project_services")
	config := meta.(*Config)

	// Get services from config
	cfgServices := getConfigServices(d)

	// Get services from API
	apiServices, err := getApiServices(d.Id(), config, ignoreProjectServices)
	if err != nil {
		return fmt.Errorf("Error updating services: %v", err)
	}

	// This call disables any APIs that aren't defined in cfgServices,
	// and enables all of those that are
	err = reconcileServices(cfgServices, apiServices, config, d.Id())
	if err != nil {
		return fmt.Errorf("Error updating services: %v", err)
	}

	return resourceGoogleProjectServicesRead(d, meta)
}

func resourceGoogleProjectServicesDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Deleting google_project_services")

	if disable := d.Get("disable_on_destroy"); !(disable.(bool)) {
		log.Printf("Not disabling service '%s', because disable_on_destroy is false.", d.Id())
		d.SetId("")
		return nil
	}

	config := meta.(*Config)
	services := resourceServices(d)
	for _, s := range services {
		if err := disableService(s, d.Id(), config, true); err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}

// This function ensures that the services enabled for a project exactly match that
// in a config by disabling any services that are returned by the API but not present
// in the config
func reconcileServices(cfgServices, apiServices []string, config *Config, pid string) error {
	// Helper to convert slice to map
	m := func(vals []string) map[string]struct{} {
		sm := make(map[string]struct{})
		for _, s := range vals {
			sm[s] = struct{}{}
		}
		return sm
	}

	sort.Strings(cfgServices)
	cfgMap := m(cfgServices)
	log.Printf("[DEBUG]: Saw the following services in config: %v", cfgServices)
	apiMap := m(apiServices)
	log.Printf("[DEBUG]: Saw the following services enabled: %v", apiServices)

	for k := range apiMap {
		if _, ok := cfgMap[k]; !ok {
			log.Printf("[DEBUG]: Disabling %s as it's enabled upstream but not in config", k)
			err := disableService(k, pid, config, true)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[DEBUG]: Skipping %s as it's enabled in both config and upstream", k)
			delete(cfgMap, k)
		}
	}

	keys := make([]string, 0, len(cfgMap))
	for k := range cfgMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	log.Printf("[DEBUG]: Enabling the following services: %v", keys)
	err := enableServices(keys, pid, config)
	if err != nil {
		return err
	}
	return nil
}

// Retrieve services defined in a config
func getConfigServices(d *schema.ResourceData) (services []string) {
	if v, ok := d.GetOk("services"); ok {
		for _, svc := range v.(*schema.Set).List() {
			services = append(services, svc.(string))
		}
	}
	return
}

// Retrieve a project's services from the API
func getApiServices(pid string, config *Config, ignore map[string]struct{}) ([]string, error) {
	if ignore == nil {
		ignore = make(map[string]struct{})
	}

	var apiServices []string

	if err := retryTime(func() error {
		// Reset the list of apiServices in case of a retry. A partial page failure
		// could result in duplicate services.
		apiServices = make([]string, 0, 10)

		ctx := context.Background()
		return config.clientServiceUsage.Services.
			List("projects/"+pid).
			Fields("services/name,nextPageToken").
			Filter("state:ENABLED").
			Pages(ctx, func(r *serviceusage.ListServicesResponse) error {
				for _, v := range r.Services {
					// services are returned as "projects/PROJECT/services/NAME"
					parts := strings.Split(v.Name, "/")
					if len(parts) > 0 {
						name := parts[len(parts)-1]
						if _, ok := ignore[name]; !ok {
							apiServices = append(apiServices, name)
						}
					}
				}

				return nil
			})
	}, 10); err != nil {
		return nil, errwrap.Wrapf("failed to list services: {{err}}", err)
	}

	return apiServices, nil
}

func enableService(s, pid string, config *Config) error {
	return enableServices([]string{s}, pid, config)
}

func enableServices(s []string, pid string, config *Config) error {
	// It's not permitted to enable more than 20 services in one API call (even
	// for batch).
	//
	// https://godoc.org/google.golang.org/api/serviceusage/v1#BatchEnableServicesRequest
	batchSize := 20

	for i := 0; i < len(s); i += batchSize {
		j := i + batchSize
		if j > len(s) {
			j = len(s)
		}

		services := s[i:j]

		if err := retryTime(func() error {
			var sop *serviceusage.Operation
			var err error

			if len(services) < 1 {
				// No more services to enable
				return nil
			} else if len(services) == 1 {
				// Use the singular enable - can't use batch for a single item
				name := fmt.Sprintf("projects/%s/services/%s", pid, services[0])
				req := &serviceusage.EnableServiceRequest{}
				sop, err = config.clientServiceUsage.Services.Enable(name, req).Do()
			} else {
				// Batch enable 2+ services
				name := fmt.Sprintf("projects/%s", pid)
				req := &serviceusage.BatchEnableServicesRequest{ServiceIds: services}
				sop, err = config.clientServiceUsage.Services.BatchEnable(name, req).Do()
			}
			if err != nil {
				// Check for a "precondition failed" error. The API seems to randomly
				// (although more than 50%) return this error when enabling certain
				// APIs. It's transient, so we catch it and re-raise it as an error that
				// is retryable instead.
				if gerr, ok := err.(*googleapi.Error); ok {
					if (gerr.Code == 400 || gerr.Code == 412) && gerr.Message == "Precondition check failed." {
						return &googleapi.Error{
							Code:    503,
							Message: "api returned \"precondition failed\" while enabling service",
						}
					}
				}
				return errwrap.Wrapf("failed to issue request: {{err}}", err)
			}

			// Poll for the API to return
			activity := fmt.Sprintf("apis %q to be enabled for %s", services, pid)
			waitErr := serviceUsageOperationWait(config, sop, activity)
			if waitErr != nil {
				return waitErr
			}

			// Accumulate the list of services that are enabled on the project
			enabledServices, err := getApiServices(pid, config, nil)
			if err != nil {
				return err
			}

			// Diff the list of requested services to enable against the list of
			// services on the project.
			missing := diffStringSlice(services, enabledServices)

			// If there are any missing, force a retry
			if len(missing) > 0 {
				// Spoof a googleapi Error so retryTime will try again
				return &googleapi.Error{
					Code:    503,
					Message: fmt.Sprintf("The service(s) %q are still being enabled for project %s. This isn't a real API error, this is just eventual consistency.", missing, pid),
				}
			}

			return nil
		}, 10); err != nil {
			return errwrap.Wrap(err, fmt.Errorf("failed to enable service(s) %q for project %s", services, pid))
		}
	}

	return nil
}

func diffStringSlice(wanted, actual []string) []string {
	var missing []string

	for _, want := range wanted {
		found := false

		for _, act := range actual {
			if want == act {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, want)
		}
	}

	return missing
}

func disableService(s, pid string, config *Config, disableDependentServices bool) error {
	err := retryTime(func() error {
		name := fmt.Sprintf("projects/%s/services/%s", pid, s)
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
	}, 10)
	if err != nil {
		return fmt.Errorf("Error disabling service %q for project %q: %v", s, pid, err)
	}
	return nil
}

func resourceServices(d *schema.ResourceData) []string {
	// Calculate the tags
	var services []string
	if s := d.Get("services"); s != nil {
		ss := s.(*schema.Set)
		services = make([]string, ss.Len())
		for i, v := range ss.List() {
			services[i] = v.(string)
		}
	}
	return services
}
