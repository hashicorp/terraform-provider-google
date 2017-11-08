package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/dns/v1"
)

func resourceDnsRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsRecordSetCreate,
		Read:   resourceDnsRecordSetRead,
		Delete: resourceDnsRecordSetDelete,
		Update: resourceDnsRecordSetUpdate,

		Schema: map[string]*schema.Schema{
			"managed_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"rrdatas": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDnsRecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("managed_zone").(string)

	// Build the change
	chg := &dns.Change{
		Additions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Name:    d.Get("name").(string),
				Type:    d.Get("type").(string),
				Ttl:     int64(d.Get("ttl").(int)),
				Rrdatas: rrdata(d),
			},
		},
	}

	// we need to replace NS record sets in the same call. That means
	// we need to list all the current NS record sets attached to the
	// zone and add them to the change as deletions. We can't just add
	// new NS record sets, or we'll get an error about the NS record set
	// already existing; see terraform-providers/terraform-provider-google#95.
	// We also can't just remove the NS recordsets on creation, as at
	// least one is required. So the solution is to "update in place" by
	// putting the addition and the removal in the same API call.
	if d.Get("type").(string) == "NS" {
		log.Printf("[DEBUG] DNS record list request for %q", zone)
		res, err := config.clientDns.ResourceRecordSets.List(project, zone).Do()
		if err != nil {
			return fmt.Errorf("Error retrieving record sets for %q: %s", zone, err)
		}
		var deletions []*dns.ResourceRecordSet

		for _, record := range res.Rrsets {
			if record.Type != "NS" {
				continue
			}
			deletions = append(deletions, record)
		}
		if len(deletions) > 0 {
			chg.Deletions = deletions
		}
	}

	log.Printf("[DEBUG] DNS Record create request: %#v", chg)
	chg, err = config.clientDns.Changes.Create(project, zone, chg).Do()
	if err != nil {
		return fmt.Errorf("Error creating DNS RecordSet: %s", err)
	}

	d.SetId(chg.Id)

	w := &DnsChangeWaiter{
		Service:     config.clientDns,
		Change:      chg,
		Project:     project,
		ManagedZone: zone,
	}
	_, err = w.Conf().WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Google DNS change: %s", err)
	}

	return resourceDnsRecordSetRead(d, meta)
}

func resourceDnsRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("managed_zone").(string)

	// name and type are effectively the 'key'
	name := d.Get("name").(string)
	dnsType := d.Get("type").(string)

	resp, err := config.clientDns.ResourceRecordSets.List(
		project, zone).Name(name).Type(dnsType).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("DNS Record Set %q", d.Get("name").(string)))
	}
	if len(resp.Rrsets) == 0 {
		// The resource doesn't exist anymore
		d.SetId("")
		return nil
	}

	if len(resp.Rrsets) > 1 {
		return fmt.Errorf("Only expected 1 record set, got %d", len(resp.Rrsets))
	}

	d.Set("ttl", resp.Rrsets[0].Ttl)
	d.Set("rrdatas", resp.Rrsets[0].Rrdatas)

	return nil
}

func resourceDnsRecordSetDelete(d *schema.ResourceData, meta interface{}) error {

	// NS records must always have a value, so we short-circuit delete
	// this allows terraform delete to work, but may have unexpected
	// side-effects when deleting just that record set.
	if d.Get("type").(string) == "NS" {
		log.Println("[DEBUG] NS records can't be deleted due to API restrictions, so they're being left in place. See https://www.terraform.io/docs/providers/google/r/dns_record_set.html for more information.")
		return nil
	}
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("managed_zone").(string)

	// Build the change
	chg := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Name:    d.Get("name").(string),
				Type:    d.Get("type").(string),
				Ttl:     int64(d.Get("ttl").(int)),
				Rrdatas: rrdata(d),
			},
		},
	}

	log.Printf("[DEBUG] DNS Record delete request: %#v", chg)
	chg, err = config.clientDns.Changes.Create(project, zone, chg).Do()
	if err != nil {
		return fmt.Errorf("Error deleting DNS RecordSet: %s", err)
	}

	w := &DnsChangeWaiter{
		Service:     config.clientDns,
		Change:      chg,
		Project:     project,
		ManagedZone: zone,
	}
	_, err = w.Conf().WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Google DNS change: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDnsRecordSetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("managed_zone").(string)
	recordName := d.Get("name").(string)

	oldTtl, newTtl := d.GetChange("ttl")
	oldType, newType := d.GetChange("type")

	oldCountRaw, _ := d.GetChange("rrdatas.#")
	oldCount := oldCountRaw.(int)

	chg := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Name:    recordName,
				Type:    oldType.(string),
				Ttl:     int64(oldTtl.(int)),
				Rrdatas: make([]string, oldCount),
			},
		},
		Additions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Name:    recordName,
				Type:    newType.(string),
				Ttl:     int64(newTtl.(int)),
				Rrdatas: rrdata(d),
			},
		},
	}

	for i := 0; i < oldCount; i++ {
		rrKey := fmt.Sprintf("rrdatas.%d", i)
		oldRR, _ := d.GetChange(rrKey)
		chg.Deletions[0].Rrdatas[i] = oldRR.(string)
	}
	log.Printf("[DEBUG] DNS Record change request: %#v old: %#v new: %#v", chg, chg.Deletions[0], chg.Additions[0])
	chg, err = config.clientDns.Changes.Create(project, zone, chg).Do()
	if err != nil {
		return fmt.Errorf("Error changing DNS RecordSet: %s", err)
	}

	w := &DnsChangeWaiter{
		Service:     config.clientDns,
		Change:      chg,
		Project:     project,
		ManagedZone: zone,
	}
	if _, err = w.Conf().WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Google DNS change: %s", err)
	}

	return resourceDnsRecordSetRead(d, meta)
}

func rrdata(
	d *schema.ResourceData,
) []string {
	rrdatasCount := d.Get("rrdatas.#").(int)
	data := make([]string, rrdatasCount)
	for i := 0; i < rrdatasCount; i++ {
		data[i] = d.Get(fmt.Sprintf("rrdatas.%d", i)).(string)
	}
	return data
}
