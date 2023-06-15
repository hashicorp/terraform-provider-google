// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tpgserviceusage "github.com/hashicorp/terraform-provider-google/google/services/serviceusage"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/serviceusage/v1"
)

// These services can only be enabled as a side-effect of enabling other services,
// so don't bother storing them in the config or using them for diffing.
var ignoredProjectServices = []string{"dataproc-control.googleapis.com", "source.googleapis.com", "stackdriverprovisioning.googleapis.com"}
var ignoredProjectServicesSet = tpgresource.GolangSetFromStringSlice(ignoredProjectServices)

// Services that can't be user-specified but are otherwise valid. Renamed
// services should be added to this set during major releases.
var bannedProjectServices = []string{"bigquery-json.googleapis.com"}

// Service Renames
// we expect when a service is renamed:
// - both service names will continue to be able to be set
// - setting one will effectively enable the other as a dependent
// - GET will return whichever service name is requested
// - LIST responses will not contain the old service name
// renames may be reverted, though, so we should canonicalise both ways until
// the old service is fully removed from the provider
//
// We handle service renames in the provider by pretending that we've read both
// the old and new service names from the API if we see either, and only setting
// the one(s) that existed in prior state in config (if any). If neither exists,
// we'll set the old service name in state.
// Additionally, in case of service rename rollbacks or unexpected early
// removals of services, if we fail to create or delete a service that's been
// renamed we'll retry using an alternate name.
// We try creation by the user-specified value followed by the other value.
// We try deletion by the old value followed by the new value.

// map from old -> new names of services that have been renamed
// these should be removed during major provider versions. comment here with
// "DEPRECATED FOR {{version}} next to entries slated for removal in {{version}}
// upon removal, we should disallow the old name from being used even if it's
// not gone from the underlying API yet
var RenamedServices = map[string]string{}

// RenamedServices in reverse (new -> old)
var renamedServicesByNewServiceNames = tpgresource.ReverseStringMap(RenamedServices)

// RenamedServices expressed as both old -> new and new -> old
var renamedServicesByOldAndNewServiceNames = tpgresource.MergeStringMaps(RenamedServices, renamedServicesByNewServiceNames)

const maxServiceUsageBatchSize = 20

func validateProjectServiceService(val interface{}, key string) (warns []string, errs []error) {
	bannedServicesFunc := verify.StringNotInSlice(append(ignoredProjectServices, bannedProjectServices...), false)
	warns, errs = bannedServicesFunc(val, key)
	if len(errs) > 0 {
		return
	}

	// StringNotInSlice already validates that this is a string
	v, _ := val.(string)
	if !strings.Contains(v, ".") {
		errs = append(errs, fmt.Errorf("expected %s to be a domain like serviceusage.googleapis.com", v))
	}
	return
}

func ResourceGoogleProjectService() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectServiceCreate,
		Read:   resourceGoogleProjectServiceRead,
		Delete: resourceGoogleProjectServiceDelete,
		Update: resourceGoogleProjectServiceUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceGoogleProjectServiceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateProjectServiceService,
			},
			"project": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
			},

			"disable_dependent_services": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"disable_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceGoogleProjectServiceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid google_project_service id format for import, expecting `{project}/{service}`, found %s", d.Id())
	}
	if err := d.Set("project", parts[0]); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("service", parts[1]); err != nil {
		return nil, fmt.Errorf("Error setting service: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceGoogleProjectServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	project = tpgresource.GetResourceNameFromSelfLink(project)

	srv := d.Get("service").(string)
	id := project + "/" + srv

	// Check if the service has already been enabled
	servicesRaw, err := BatchRequestReadServices(project, d, config)
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}
	servicesList := servicesRaw.(map[string]struct{})
	if _, ok := servicesList[srv]; ok {
		log.Printf("[DEBUG] service %s was already found to be enabled in project %s", srv, project)
		d.SetId(id)
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("service", srv); err != nil {
			return fmt.Errorf("Error setting service: %s", err)
		}
		return nil
	}

	err = BatchRequestEnableService(srv, project, d, config)
	if err != nil {
		return err
	}
	d.SetId(id)
	return resourceGoogleProjectServiceRead(d, meta)
}

func resourceGoogleProjectServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	project = tpgresource.GetResourceNameFromSelfLink(project)

	// Verify project for services still exists
	projectGetCall := config.NewResourceManagerClient(userAgent).Projects.Get(project)
	if config.UserProjectOverride {
		billingProject := project

		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}
		projectGetCall.Header().Add("X-Goog-User-Project", billingProject)
	}
	p, err := projectGetCall.Do()

	if err == nil && p.LifecycleState == "DELETE_REQUESTED" {
		// Construct a 404 error for transport_tpg.HandleNotFoundError
		err = &googleapi.Error{
			Code:    404,
			Message: "Project deletion was requested",
		}
	}
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}

	servicesRaw, err := BatchRequestReadServices(project, d, config)
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}
	servicesList := servicesRaw.(map[string]struct{})

	srv := d.Get("service").(string)
	if _, ok := servicesList[srv]; ok {
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("service", srv); err != nil {
			return fmt.Errorf("Error setting service: %s", err)
		}
		return nil
	}

	log.Printf("[DEBUG] service %s not in enabled services for project %s, removing from state", srv, project)
	d.SetId("")
	return nil
}

func resourceGoogleProjectServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	if disable := d.Get("disable_on_destroy"); !(disable.(bool)) {
		log.Printf("[WARN] Project service %q disable_on_destroy is false, skip disabling service", d.Id())
		d.SetId("")
		return nil
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	project = tpgresource.GetResourceNameFromSelfLink(project)

	service := d.Get("service").(string)
	disableDependencies := d.Get("disable_dependent_services").(bool)
	if err = disableServiceUsageProjectService(service, project, d, config, disableDependencies); err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}

	d.SetId("")
	return nil
}

func resourceGoogleProjectServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	// This update method is no-op because the only updatable fields
	// are state/config-only, i.e. they aren't sent in requests to the API.
	return nil
}

// Disables a project service.
func disableServiceUsageProjectService(service, project string, d *schema.ResourceData, config *transport_tpg.Config, disableDependentServices bool) error {
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			billingProject := project
			userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
			if err != nil {
				return err
			}
			name := fmt.Sprintf("projects/%s/services/%s", project, service)
			servicesDisableCall := config.NewServiceUsageClient(userAgent).Services.Disable(name, &serviceusage.DisableServiceRequest{
				DisableDependentServices: disableDependentServices,
			})
			if config.UserProjectOverride {
				// err == nil indicates that the billing_project value was found
				if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
					billingProject = bp
				}
				servicesDisableCall.Header().Add("X-Goog-User-Project", billingProject)
			}
			sop, err := servicesDisableCall.Do()
			if err != nil {
				return err
			}
			// Wait for the operation to complete
			waitErr := tpgserviceusage.ServiceUsageOperationWait(config, sop, billingProject, "api to disable", userAgent, d.Timeout(schema.TimeoutDelete))
			if waitErr != nil {
				return waitErr
			}
			return nil
		},
		Timeout:              d.Timeout(schema.TimeoutDelete),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.ServiceUsageServiceBeingActivated},
	})
	if err != nil {
		return fmt.Errorf("Error disabling service %q for project %q: %v", service, project, err)
	}
	return nil
}
