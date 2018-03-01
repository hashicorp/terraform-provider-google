package google

import (
	"fmt"
	"log"
	"regexp"

	"google.golang.org/api/dns/v1"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

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
			"managed_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateDNSRRsetName,
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
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

// This catches most invalid DNS names, though it might not catch some invlaid non-ASCII names
// it does not implement the TOASCII functionality defined in https://tools.ietf.org/html/rfc3454
// See:
// https://tools.ietf.org/html/rfc1034 DOMAIN NAMES - CONCEPTS AND FACILITIES
// https://tools.ietf.org/html/rfc3490 Internationalizing Domain Names in Applications (IDNA)
// https://tools.ietf.org/html/rfc3454 Preparation of Internationalized Strings ('stringprep')
func ValidateDNSRRsetName(i interface{}, k string) (s []string, errors []error) {
	name := i.(string)
	errors = append(errors, checkName(name)...)
	return
}

func checkName(name string) (errors []error) {

	if len(name) > 255 {
		errors = append(errors, fmt.Errorf("DNS name length %d greater than 254",
			len(name)))
	}

	// Compile regular expression for use in the loop

	// https://tools.ietf.org/html/rfc3490#section-3.1 Requirements
	// Whenever dots are used as label separators, the following
	// characters MUST be recognized as dots: U+002E (full stop),
	// U+3002 (ideographic full stop), U+FF0E (fullwidth full
	// stop), U+FF61 (halfwidth ideographic full stop).
	re, err := regexp.Compile("[\u002e\u3002\uFF0E\uFF61]$")
	if err != nil {
		errors = append(errors, fmt.Errorf("Internal error compiling regexp: %s",
			err))
		re = nil
	}
	if re != nil && !re.MatchString(name) {
		errors = append(errors,
			fmt.Errorf("DNS name must be fully qualified with a trailing dot"))
	}

	dotRegexp, err := regexp.Compile("[\u002e\u3002\uFF0E\uFF61]")
	if err != nil {
		errors = append(errors, fmt.Errorf("Internal error compiling regexp: %s",
			err))
		dotRegexp = nil
	}

	// Normally, if we were not accepting internationalized domain
	// names, we could just use isAscii here. But that will reject
	// legal domain names. Since ASCII and UTF-8 have the same
	// encoding in the range that ASCII covers, we still want to
	// flag characters that are not valid ASCII in the 0X00-0XFF
	// range. This regular expression excludes those characters
	verbotenRegexp, err := regexp.Compile("[\x00-\x2C\x2E-\x2F\x3A-\x40\x5B-\x60\x7B-\x7F]+")
	if err != nil {
		errors = append(errors, fmt.Errorf("Internal error compiling regexp: %s",
			err))
		verbotenRegexp = nil
	}

	allDigitRegexp, err := regexp.Compile("^[[:digit:]]+$")
	if err != nil {
		errors = append(errors, fmt.Errorf("Internal error compiling regexp: %s",
			err))
		allDigitRegexp = nil
	}

	spaceRegexp, err := regexp.Compile("[[:space:]]+")
	if err != nil {
		errors = append(errors, fmt.Errorf("Internal error compiling regexp: %s",
			err))
		spaceRegexp = nil
	}

	underscoreRegexp, err := regexp.Compile("_+")
	if err != nil {
		errors = append(errors, fmt.Errorf("Internal error compiling regexp: %s",
			err))
		underscoreRegexp = nil
	}

	// Split the DNS name into labels, and iterate over them with
	// precompile regular expressions.

	var labels []string
	if dotRegexp != nil {
		labels = dotRegexp.Split(name, -1)
		if len(labels) < 2 {
			errors = append(errors,
				fmt.Errorf("DNS name must not be a top level domain (%s).", name))
			labels = nil
		}
	} else {
		errors = append(errors, fmt.Errorf("Could not split DNS name."))
	}

	for index, label := range labels {
		// First, the simpler checks

		// Only the final label may be empty
		if index < len(labels)-1 && len(label) < 1 {
			errors = append(errors, fmt.Errorf("subdomains must not be empty"))
		}
		// errors = append(errors, fmt.Errorf("DEBUG: %d label is %s", index, label))

		s := []rune(label)
		if len(s) > 63 {
			errors = append(errors,
				fmt.Errorf("subdomains must not have more than 63 characters: %s - %d.",
					label, len(s)))
		}

		if len(s) > 0 {
			if string(s[0]) == "-" {
				errors = append(errors, fmt.Errorf("subdomains must not start with a hyphen (%s).",
					label))
			}
			if string(s[len(s)-1]) == "-" {
				errors = append(errors, fmt.Errorf("subdomains must not end with a hyphen. (%s)",
					label))
			}

			if len(s) > 4 && string(s[2:4]) == "--" {
				errors = append(errors, fmt.Errorf("subdomains must not mask encoded IDNA: %s.",
					label))
			}
			if underscoreRegexp != nil {
				rest := s[1:]
				if underscoreRegexp.MatchString(string(rest)) {
					errors = append(errors,
						fmt.Errorf("Only the leading character in a subdomain may be a '_' (%s).",
							label))
				}
			}
		}
		if spaceRegexp != nil && spaceRegexp.MatchString(label) {
			errors = append(errors,
				fmt.Errorf("Subdomain labels cannot contain whitespace: %s",
					label))
		}

		if allDigitRegexp != nil && allDigitRegexp.MatchString(label) {
			errors = append(errors,
				fmt.Errorf("Subdomain labels cannot all be digits: %s",
					label))
		}

		if verbotenRegexp != nil {
			badCharacters := verbotenRegexp.FindAllStringSubmatch(label, -1)
			if badCharacters != nil {
				errors = append(errors,
					fmt.Errorf("Illegal ASCII character(s) %q  in domain name.",
						badCharacters))
			}
		}

	}
	return
}

func resourceDnsRecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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
			&dns.ResourceRecordSet{
				Name:    name,
				Type:    rType,
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
	if rType == "NS" {
		log.Printf("[DEBUG] DNS record list request for %q", zone)
		res, err := config.clientDns.ResourceRecordSets.List(project, zone).Do()
		if err != nil {
			return fmt.Errorf("Error retrieving record sets for %q: %s", zone, err)
		}
		var deletions []*dns.ResourceRecordSet

		for _, record := range res.Rrsets {
			if record.Type != "NS" || record.Name != name {
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

	d.SetId(fmt.Sprintf("%s/%s/%s", zone, name, rType))

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

	d.Set("type", resp.Rrsets[0].Type)
	d.Set("ttl", resp.Rrsets[0].Ttl)
	d.Set("rrdatas", resp.Rrsets[0].Rrdatas)
	d.Set("project", project)

	return nil
}

func resourceDnsRecordSetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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
		mz, err := config.clientDns.ManagedZones.Get(project, zone).Do()
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

func resourceDnsRecordSetImportState(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid dns record specifier. Expecting {zone-name}/{record-name}/{record-type}. The record name must include a trailing '.' at the end.")
	}

	d.Set("managed_zone", parts[0])
	d.Set("name", parts[1])
	d.Set("type", parts[2])

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
