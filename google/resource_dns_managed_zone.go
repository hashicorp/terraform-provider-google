package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/dns/v1"
)

func resourceDnsManagedZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsManagedZoneCreate,
		Read:   resourceDnsManagedZoneRead,
		Update: resourceDnsManagedZoneUpdate,
		Delete: resourceDnsManagedZoneDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDnsManagedZoneImport,
		},
		Schema: map[string]*schema.Schema{
			"dns_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},

			"name_servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Google Cloud DNS ManagedZone resources do not have a SelfLink attribute.

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDnsManagedZoneCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the parameter
	zone := &dns.ManagedZone{
		Name:        d.Get("name").(string),
		DnsName:     d.Get("dns_name").(string),
		Description: d.Get("description").(string),
	}

	if _, ok := d.GetOk("labels"); ok {
		zone.Labels = expandLabels(d)
	}

	log.Printf("[DEBUG] DNS ManagedZone create request: %#v", zone)

	zone, err = config.clientDns.ManagedZones.Create(project, zone).Do()
	if err != nil {
		return fmt.Errorf("Error creating DNS ManagedZone: %s", err)
	}

	d.SetId(zone.Name)

	return resourceDnsManagedZoneRead(d, meta)
}

func resourceDnsManagedZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := config.clientDns.ManagedZones.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("DNS Managed Zone %q", d.Get("name").(string)))
	}

	d.Set("name_servers", zone.NameServers)
	d.Set("name", zone.Name)
	d.Set("dns_name", zone.DnsName)
	d.Set("description", zone.Description)
	d.Set("project", project)
	d.Set("labels", zone.Labels)

	return nil
}

func resourceDnsManagedZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := &dns.ManagedZone{
		Name:        d.Get("name").(string),
		DnsName:     d.Get("dns_name").(string),
		Description: d.Get("description").(string),
	}

	if _, ok := d.GetOk("labels"); ok {
		zone.Labels = expandLabels(d)
	}

	op, err := config.clientDns.ManagedZones.Patch(project, d.Id(), zone).Do()
	if err != nil {
		return err
	}

	err = dnsOperationWait(config.clientDns, op, project, "Updating DNS Managed Zone")
	if err != nil {
		return err
	}

	return resourceDnsManagedZoneRead(d, meta)
}

func resourceDnsManagedZoneDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	err = config.clientDns.ManagedZones.Delete(project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting DNS ManagedZone: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDnsManagedZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	parseImportId([]string{"projects/(?P<project>[^/]+)/managedZones/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/managedZones/(?P<name>[^/]+)",
		"(?P<name>[^/]+)"}, d, config)

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
