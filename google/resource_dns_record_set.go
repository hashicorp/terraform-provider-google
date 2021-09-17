package google

import (
	"fmt"
	"log"

	"strings"

	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/dns/v1"
)

func rrdatasDnsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange("rrdatas")
	if o == nil || n == nil {
		return false
	}

	oList := convertStringArr(o.([]interface{}))
	nList := convertStringArr(n.([]interface{}))

	parseFunc := func(record string) string {
		switch d.Get("type") {
		case "AAAA":
			// parse ipv6 to a key from one list
			return net.ParseIP(record).String()
		case "MX", "DS":
			return strings.ToLower(record)
		case "TXT":
			return strings.ToLower(strings.Trim(record, `"`))
		default:
			return record
		}
	}
	return rrdatasListDiffSuppress(oList, nList, parseFunc, d)
}

// suppress on a list when 1) its items have dups that need to be ignored
// and 2) string comparison on the items may need a special parse function
// example of usage can be found ../../../third_party/terraform/tests/resource_dns_record_set_test.go.erb
func rrdatasListDiffSuppress(oldList, newList []string, fun func(x string) string, _ *schema.ResourceData) bool {
	// compare two lists of unordered records
	diff := make(map[string]bool, len(oldList))
	for _, oldRecord := range oldList {
		// set all new IPs to true
		diff[fun(oldRecord)] = true
	}
	for _, newRecord := range newList {
		// set matched IPs to false otherwise can't suppress
		if diff[fun(newRecord)] {
			diff[fun(newRecord)] = false
		} else {
			return false
		}
	}
	// can't suppress if unmatched records are found
	for _, element := range diff {
		if element {
			return false
		}
	}
	return true
}

func resourceDnsRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsRecordSetCreate,
		Read:   resourceDnsRecordSetRead,
		Delete: resourceDnsRecordSetDelete,
		Update: resourceDnsRecordSetUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceDnsRecordSetImportState,
		},

		Schema: map[string]*schema.Schema{
			"managed_zone": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name of the zone in which this record set will reside.`,
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The DNS name this record set will apply to.`,
			},

			"rrdatas": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: rrdatasDnsDiffSuppress,
				Description:      `The string data for the records in this record set whose meaning depends on the DNS type. For TXT record, if the string data contains spaces, add surrounding \" if you don't want your string to get split on spaces. To specify a single record value longer than 255 characters such as a TXT record for DKIM, add \"\" inside the Terraform configuration string (e.g. "first255characters\"\"morecharacters").`,
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The time-to-live of this record set (seconds).`,
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The DNS record set type.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceDnsRecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	zone := d.Get("managed_zone").(string)
	rType := d.Get("type").(string)

	// Build the change
	chg := &dns.Change{
		Additions: []*dns.ResourceRecordSet{
			{
				Name:    name,
				Type:    rType,
				Ttl:     int64(d.Get("ttl").(int)),
				Rrdatas: rrdata(d),
			},
		},
	}

	// The terraform provider is authoritative, so what we do here is check if
	// any records that we are trying to create already exist and make sure we
	// delete them, before adding in the changes requested.  Normally this would
	// result in an AlreadyExistsError.
	log.Printf("[DEBUG] DNS record list request for %q", zone)
	res, err := config.NewDnsClient(userAgent).ResourceRecordSets.List(project, zone).Do()
	if err != nil {
		return fmt.Errorf("Error retrieving record sets for %q: %s", zone, err)
	}
	var deletions []*dns.ResourceRecordSet

	for _, record := range res.Rrsets {
		if record.Type != rType || record.Name != name {
			continue
		}
		deletions = append(deletions, record)
	}
	if len(deletions) > 0 {
		chg.Deletions = deletions
	}

	log.Printf("[DEBUG] DNS Record create request: %#v", chg)
	chg, err = config.NewDnsClient(userAgent).Changes.Create(project, zone, chg).Do()
	if err != nil {
		return fmt.Errorf("Error creating DNS RecordSet: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/managedZones/%s/rrsets/%s/%s", project, zone, name, rType))

	w := &DnsChangeWaiter{
		Service:     config.NewDnsClient(userAgent),
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("managed_zone").(string)

	// name and type are effectively the 'key'
	name := d.Get("name").(string)
	dnsType := d.Get("type").(string)

	var resp *dns.ResourceRecordSetsListResponse
	err = retry(func() error {
		var reqErr error
		resp, reqErr = config.NewDnsClient(userAgent).ResourceRecordSets.List(
			project, zone).Name(name).Type(dnsType).Do()
		return reqErr
	})
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

	if err := d.Set("type", resp.Rrsets[0].Type); err != nil {
		return fmt.Errorf("Error setting type: %s", err)
	}
	if err := d.Set("ttl", resp.Rrsets[0].Ttl); err != nil {
		return fmt.Errorf("Error setting ttl: %s", err)
	}
	if err := d.Set("rrdatas", resp.Rrsets[0].Rrdatas); err != nil {
		return fmt.Errorf("Error setting rrdatas: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

func resourceDnsRecordSetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("managed_zone").(string)

	// NS records must always have a value, so we short-circuit delete
	// this allows terraform delete to work, but may have unexpected
	// side-effects when deleting just that record set.
	// Unfortunately, you can set NS records on subdomains, and those
	// CAN and MUST be deleted, so we need to retrieve the managed zone,
	// check if what we're looking at is a subdomain, and only not delete
	// if it's not actually a subdomain
	if d.Get("type").(string) == "NS" {
		mz, err := config.NewDnsClient(userAgent).ManagedZones.Get(project, zone).Do()
		if err != nil {
			return fmt.Errorf("Error retrieving managed zone %q from %q: %s", zone, project, err)
		}
		domain := mz.DnsName

		if domain == d.Get("name").(string) {
			log.Println("[DEBUG] NS records can't be deleted due to API restrictions, so they're being left in place. See https://www.terraform.io/docs/providers/google/r/dns_record_set.html for more information.")
			return nil
		}
	}

	// Build the change
	chg := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			{
				Name:    d.Get("name").(string),
				Type:    d.Get("type").(string),
				Ttl:     int64(d.Get("ttl").(int)),
				Rrdatas: rrdata(d),
			},
		},
	}

	log.Printf("[DEBUG] DNS Record delete request: %#v", chg)
	chg, err = config.NewDnsClient(userAgent).Changes.Create(project, zone, chg).Do()
	if err != nil {
		return handleNotFoundError(err, d, "google_dns_record_set")
	}

	w := &DnsChangeWaiter{
		Service:     config.NewDnsClient(userAgent),
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

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
			{
				Name:    recordName,
				Type:    oldType.(string),
				Ttl:     int64(oldTtl.(int)),
				Rrdatas: make([]string, oldCount),
			},
		},
		Additions: []*dns.ResourceRecordSet{
			{
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
	chg, err = config.NewDnsClient(userAgent).Changes.Create(project, zone, chg).Do()
	if err != nil {
		return fmt.Errorf("Error changing DNS RecordSet: %s", err)
	}

	w := &DnsChangeWaiter{
		Service:     config.NewDnsClient(userAgent),
		Change:      chg,
		Project:     project,
		ManagedZone: zone,
	}
	if _, err = w.Conf().WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Google DNS change: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/managedZones/%s/rrsets/%s/%s", project, zone, recordName, newType))

	return resourceDnsRecordSetRead(d, meta)
}

func resourceDnsRecordSetImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/managedZones/(?P<managed_zone>[^/]+)/rrsets/(?P<name>[^/]+)/(?P<type>[^/]+)",
		"(?P<project>[^/]+)/(?P<managed_zone>[^/]+)/(?P<name>[^/]+)/(?P<type>[^/]+)",
		"(?P<managed_zone>[^/]+)/(?P<name>[^/]+)/(?P<type>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/managedZones/{{managed_zone}}/rrsets/{{name}}/{{type}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
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
