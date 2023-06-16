// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicenetworking

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/servicenetworking/v1"
)

func ResourceGoogleServiceNetworkingPeeredDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleServiceNetworkingPeeredDNSDomainCreate,
		Read:   resourceGoogleServiceNetworkingPeeredDNSDomainRead,
		Delete: resourceGoogleServiceNetworkingPeeredDNSDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceGoogleServiceNetworkingPeeredDNSDomainImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the project that the service account will be created in. Defaults to the provider project configuration.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the peered DNS domain.",
			},
			"dns_suffix": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The DNS domain name suffix of the peered DNS domain.",
			},
			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "servicenetworking.googleapis.com",
				Description: "The name of the service to create a peered DNS domain for, e.g. servicenetworking.googleapis.com",
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network in the consumer project to peer with.",
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceGoogleServiceNetworkingPeeredDNSDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 9 {
		return nil, fmt.Errorf("Invalid google_project_service_peered_dns_domain id format for import, expecting `services/{service}/projects/{project}/global/networks/{network}/peeredDnsDomains/{name}`, found %s", d.Id())
	}
	if err := d.Set("service", parts[1]); err != nil {
		return nil, fmt.Errorf("Error setting service: %s", err)
	}
	if err := d.Set("project", parts[3]); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("network", parts[6]); err != nil {
		return nil, fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("name", parts[8]); err != nil {
		return nil, fmt.Errorf("Error setting name: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceGoogleServiceNetworkingPeeredDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	projectNumber, err := getProjectNumber(d, config, project, userAgent)
	if err != nil {
		return err
	}

	service := d.Get("service").(string)
	network := d.Get("network").(string)
	parent := fmt.Sprintf("services/%s/projects/%s/global/networks/%s", service, projectNumber, network)

	name := d.Get("name").(string)
	dnsSuffix := d.Get("dns_suffix").(string)
	r := &servicenetworking.PeeredDnsDomain{
		DnsSuffix: dnsSuffix,
		Name:      name,
	}

	apiService := config.NewServiceNetworkingClient(userAgent)
	peeredDnsDomainsService := servicenetworking.NewServicesProjectsGlobalNetworksPeeredDnsDomainsService(apiService)
	createCall := peeredDnsDomainsService.Create(parent, r)
	if config.UserProjectOverride {
		createCall.Header().Add("X-Goog-User-Project", project)
	}
	op, err := createCall.Do()
	if err != nil {
		return err
	}

	if err := ServiceNetworkingOperationWaitTime(config, op, "Create Service Networking Peered DNS Domain", userAgent, project, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	if err := d.Set("parent", parent); err != nil {
		return fmt.Errorf("Error setting parent: %s", err)
	}
	id := fmt.Sprintf("%s/peeredDnsDomains/%s", parent, name)
	d.SetId(id)
	return resourceGoogleServiceNetworkingPeeredDNSDomainRead(d, meta)
}

func resourceGoogleServiceNetworkingPeeredDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	projectNumber, err := getProjectNumber(d, config, project, userAgent)
	if err != nil {
		return err
	}

	service := d.Get("service").(string)
	network := d.Get("network").(string)
	parent := fmt.Sprintf("services/%s/projects/%s/global/networks/%s", service, projectNumber, network)

	apiService := config.NewServiceNetworkingClient(userAgent)
	peeredDnsDomainsService := servicenetworking.NewServicesProjectsGlobalNetworksPeeredDnsDomainsService(apiService)
	readCall := peeredDnsDomainsService.List(parent)
	if config.UserProjectOverride {
		readCall.Header().Add("X-Goog-User-Project", project)
	}
	response, err := readCall.Do()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	id := fmt.Sprintf("%s/peeredDnsDomains/%s", parent, name)
	d.SetId(id)

	var peeredDnsDomain *servicenetworking.PeeredDnsDomain
	for _, c := range response.PeeredDnsDomains {
		if c.Name == name {
			peeredDnsDomain = c
			break
		}
	}

	if peeredDnsDomain == nil {
		d.SetId("")
		log.Printf("[WARNING] Failed to find Service Peered DNS Domain, service: %s, project: %s, network: %s, name: %s", service, project, network, name)
		return nil
	}

	if err := d.Set("network", network); err != nil {
		return fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("name", peeredDnsDomain.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("dns_suffix", peeredDnsDomain.DnsSuffix); err != nil {
		return fmt.Errorf("Error setting peering: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("parent", parent); err != nil {
		return fmt.Errorf("Error setting parent: %s", err)
	}

	return nil
}

func resourceGoogleServiceNetworkingPeeredDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	apiService := config.NewServiceNetworkingClient(userAgent)
	peeredDnsDomainsService := servicenetworking.NewServicesProjectsGlobalNetworksPeeredDnsDomainsService(apiService)

	if err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			_, delErr := peeredDnsDomainsService.Delete(d.Id()).Do()
			return delErr
		},
		Timeout: d.Timeout(schema.TimeoutDelete),
	}); err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Peered DNS domain %s", name))
	}

	d.SetId("")
	return nil
}

// NOTE(deviavir): An out of band aspect of this API is that it uses a unique formatting of network
// different from the standard self_link URI. It requires a call to the resource manager to get the project
// number for the current project.
func getProjectNumber(d *schema.ResourceData, config *transport_tpg.Config, project, userAgent string) (string, error) {
	log.Printf("[DEBUG] Retrieving project number by doing a GET with the project id, as required by service networking")
	// err == nil indicates that the billing_project value was found
	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	getProjectCall := config.NewResourceManagerClient(userAgent).Projects.Get(project)
	if config.UserProjectOverride {
		getProjectCall.Header().Add("X-Goog-User-Project", billingProject)
	}
	projectCall, err := getProjectCall.Do()
	if err != nil {
		// note: returning a wrapped error is part of this method's contract!
		// https://blog.golang.org/go1.13-errors
		return "", fmt.Errorf("Failed to retrieve project, project: %s, err: %w", project, err)
	}

	return strconv.FormatInt(projectCall.ProjectNumber, 10), nil
}
