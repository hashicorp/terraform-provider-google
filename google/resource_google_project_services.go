package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1beta1"
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
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"services": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"disable_on_destroy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

// These services can only be enabled as a side-effect of enabling other services,
// so don't bother storing them in the config or using them for diffing.
var ignoreProjectServices = map[string]struct{}{
	"containeranalysis.googleapis.com": struct{}{},
	"dataproc-control.googleapis.com":  struct{}{},
	"source.googleapis.com":            struct{}{},
}

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
		disableService(s, d.Id(), config)
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

	cfgMap := m(cfgServices)
	apiMap := m(apiServices)

	for k, _ := range apiMap {
		if _, ok := cfgMap[k]; !ok {
			// The service in the API is not in the config; disable it.
			err := disableService(k, pid, config)
			if err != nil {
				return err
			}
		} else {
			// The service exists in the config and the API, so we don't need
			// to re-enable it
			delete(cfgMap, k)
		}
	}

	keys := make([]string, 0, len(cfgMap))
	for k, _ := range cfgMap {
		keys = append(keys, k)
	}
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
	apiServices := make([]string, 0)
	// Get services from the API
	token := ""
	for paginate := true; paginate; {
		svcResp, err := config.clientServiceUsage.Services.List("projects/" + pid).PageToken(token).Filter("state:ENABLED").Do()
		if err != nil {
			return apiServices, err
		}
		for _, v := range svcResp.Services {
			// names are returned as projects/{project-number}/services/{service-name}
			nameParts := strings.Split(v.Name, "/")
			name := nameParts[len(nameParts)-1]
			if _, ok := ignore[name]; !ok {
				apiServices = append(apiServices, name)
			}
		}
		token = svcResp.NextPageToken
		paginate = token != ""
	}
	return apiServices, nil
}

func enableService(s, pid string, config *Config) error {
	return enableServices([]string{s}, pid, config)
}

func enableServices(s []string, pid string, config *Config) error {
	err := retryTime(func() error {
		var sop *serviceusage.Operation
		var err error
		if len(s) > 1 {
			req := &serviceusage.BatchEnableServicesRequest{ServiceIds: s}
			sop, err = config.clientServiceUsage.Services.BatchEnable("projects/"+pid, req).Do()
		} else if len(s) == 1 {
			name := fmt.Sprintf("projects/%s/services/%s", pid, s[0])
			sop, err = config.clientServiceUsage.Services.Enable(name, &serviceusage.EnableServiceRequest{}).Do()
		} else {
			// No services to enable
			return nil
		}
		if err != nil {
			return err
		}
		_, waitErr := serviceUsageOperationWait(config, sop, "api to enable")
		if waitErr != nil {
			return waitErr
		}
		services, err := getApiServices(pid, config, map[string]struct{}{})
		if err != nil {
			return err
		}
		var missing []string
		for _, toEnable := range s {
			var found bool
			for _, service := range services {
				if service == toEnable {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, toEnable)
			}
		}
		if len(missing) > 0 {
			// spoof a googleapi Error so retryTime will try again
			return &googleapi.Error{
				Code:    503, // haha, get it, service unavailable
				Message: fmt.Sprintf("The services %s are still being enabled for project %q. This isn't a real API error, this is just eventual consistency.", strings.Join(missing, ", "), pid),
			}
		}
		return nil
	}, 10)
	if err != nil {
		return fmt.Errorf("Error enabling service %q for project %q: %v", s, pid, err)
	}
	return nil
}

func disableService(s, pid string, config *Config) error {
	err := retryTime(func() error {
		name := fmt.Sprintf("projects/%s/services/%s", pid, s)
		sop, err := config.clientServiceUsage.Services.Disable(name, &serviceusage.DisableServiceRequest{}).Do()
		if err != nil {
			return err
		}
		// Wait for the operation to complete
		_, waitErr := serviceUsageOperationWait(config, sop, "api to disable")
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
