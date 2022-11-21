package google

import (
	"fmt"
	"log"

	"strings"

	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/dns/v1"
)

func rrdatasDnsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "rrdatas.#" && (new == "0" || new == "") && old != new {
		return false
	}

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRecordNameTrailingDot,
				Description:  `The DNS name this record set will apply to.`,
			},

			"rrdatas": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: rrdatasDnsDiffSuppress,
				Description:      `The string data for the records in this record set whose meaning depends on the DNS type. For TXT record, if the string data contains spaces, add surrounding \" if you don't want your string to get split on spaces. To specify a single record value longer than 255 characters such as a TXT record for DKIM, add \"\" inside the Terraform configuration string (e.g. "first255characters\"\"morecharacters").`,
				ExactlyOneOf:     []string{"rrdatas", "routing_policy"},
			},

			"routing_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The configuration for steering traffic based on query. You can specify either Weighted Round Robin(WRR) type or Geolocation(GEO) type.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"wrr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The configuration for Weighted Round Robin based routing policy.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"weight": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: `The ratio of traffic routed to the target.`,
									},
									"rrdatas": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"health_checked_targets": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The list of targets to be health checked. Note that if DNSSEC is enabled for this zone, only one of `rrdatas` or `health_checked_targets` can be set.",
										MaxItems:    1,
										Elem:        healthCheckedTargetSchema,
									},
								},
							},
							ExactlyOneOf:  []string{"routing_policy.0.wrr", "routing_policy.0.geo", "routing_policy.0.primary_backup"},
							ConflictsWith: []string{"routing_policy.0.enable_geo_fencing"},
						},
						"geo": {
							Type:         schema.TypeList,
							Optional:     true,
							Description:  `The configuration for Geo location based routing policy.`,
							Elem:         geoPolicySchema,
							ExactlyOneOf: []string{"routing_policy.0.wrr", "routing_policy.0.geo", "routing_policy.0.primary_backup"},
						},
						"enable_geo_fencing": {
							Type:          schema.TypeBool,
							Optional:      true,
							Description:   "Specifies whether to enable fencing for geo queries.",
							ConflictsWith: []string{"routing_policy.0.wrr", "routing_policy.0.primary_backup"},
						},
						"primary_backup": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The configuration for a primary-backup policy with global to regional failover. Queries are responded to with the global primary targets, but if none of the primary targets are healthy, then we fallback to a regional failover policy.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"primary": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The list of global primary targets to be health checked.",
										MaxItems:    1,
										Elem:        healthCheckedTargetSchema,
									},
									"backup_geo": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The backup geo targets, which provide a regional failover policy for the otherwise global primary targets.",
										Elem:        geoPolicySchema,
									},
									"enable_geo_fencing_for_backups": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies whether to enable fencing for backup geo queries.",
									},
									"trickle_ratio": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Specifies the percentage of traffic to send to the backup targets even when the primary targets are healthy.",
									},
								},
							},
							ExactlyOneOf:  []string{"routing_policy.0.wrr", "routing_policy.0.geo", "routing_policy.0.primary_backup"},
							ConflictsWith: []string{"routing_policy.0.enable_geo_fencing"},
						},
					},
				},
				ExactlyOneOf: []string{"rrdatas", "routing_policy"},
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

var geoPolicySchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"location": {
			Type:        schema.TypeString,
			Required:    true,
			Description: `The location name defined in Google Cloud.`,
		},
		"rrdatas": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"health_checked_targets": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "For A and AAAA types only. The list of targets to be health checked. These can be specified along with `rrdatas` within this item.",
			MaxItems:    1,
			Elem:        healthCheckedTargetSchema,
		},
	},
}

var healthCheckedTargetSchema *schema.Resource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"internal_load_balancers": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The list of internal load balancers to health check.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"load_balancer_type": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  `The type of load balancer. This value is case-sensitive. Possible values: ["regionalL4ilb"]`,
						ValidateFunc: validation.StringInSlice([]string{"regionalL4ilb"}, false),
					},
					"ip_address": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The frontend IP address of the load balancer.",
					},
					"port": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The configured port of the load balancer.",
					},
					"ip_protocol": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  `The configured IP protocol of the load balancer. This value is case-sensitive. Possible values: ["tcp", "udp"]`,
						ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
					},
					"network_url": {
						Type:             schema.TypeString,
						Required:         true,
						DiffSuppressFunc: compareSelfLinkOrResourceName,
						Description:      "The fully qualified url of the network in which the load balancer belongs. This should be formatted like `https://www.googleapis.com/compute/v1/projects/{project}/global/networks/{network}`.",
					},
					"project": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The ID of the project in which the load balancer belongs.",
					},
					"region": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The region of the load balancer. Only needed for regional load balancers.",
					},
				},
			},
		},
	},
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
	rset := &dns.ResourceRecordSet{
		Name: name,
		Type: rType,
		Ttl:  int64(d.Get("ttl").(int)),
	}
	if rrdatas := expandDnsRecordSetRrdata(d.Get("rrdatas").([]interface{})); len(rrdatas) > 0 {
		rset.Rrdatas = rrdatas
	}

	rp, err := expandDnsRecordSetRoutingPolicy(d.Get("routing_policy").([]interface{}), d, config)
	if err != nil {
		return err
	}
	if rp != nil {
		rset.RoutingPolicy = rp
	}
	chg := &dns.Change{
		Additions: []*dns.ResourceRecordSet{rset},
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
	rrset := resp.Rrsets[0]
	if err := d.Set("type", rrset.Type); err != nil {
		return fmt.Errorf("Error setting type: %s", err)
	}
	if err := d.Set("ttl", rrset.Ttl); err != nil {
		return fmt.Errorf("Error setting ttl: %s", err)
	}
	if len(rrset.Rrdatas) > 0 {
		if err := d.Set("rrdatas", rrset.Rrdatas); err != nil {
			return fmt.Errorf("Error setting rrdatas: %s", err)
		}
	}
	if rrset.RoutingPolicy != nil {
		if err := d.Set("routing_policy", flattenDnsRecordSetRoutingPolicy(rrset.RoutingPolicy)); err != nil {
			return fmt.Errorf("Error setting routing_policy: %s", err)
		}
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
			log.Println("[DEBUG] NS records can't be deleted due to API restrictions, so they're being left in place. See https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dns_record_set for more information.")
			return nil
		}
	}

	routingPolicy, err := expandDnsRecordSetRoutingPolicy(d.Get("routing_policy").([]interface{}), d, config)
	if err != nil {
		return err
	}

	// Build the change
	chg := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			{
				Name:          d.Get("name").(string),
				Type:          d.Get("type").(string),
				Ttl:           int64(d.Get("ttl").(int)),
				Rrdatas:       expandDnsRecordSetRrdata(d.Get("rrdatas").([]interface{})),
				RoutingPolicy: routingPolicy,
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

	oldRoutingPolicyRaw, _ := d.GetChange("routing_policy")
	oldRoutingPolicyList := oldRoutingPolicyRaw.([]interface{})

	oldRoutingPolicy, err := expandDnsRecordSetRoutingPolicy(oldRoutingPolicyList, d, config)
	if err != nil {
		return err
	}

	newRoutingPolicy, err := expandDnsRecordSetRoutingPolicy(d.Get("routing_policy").([]interface{}), d, config)
	if err != nil {
		return err
	}

	chg := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			{
				Name:          recordName,
				Type:          oldType.(string),
				Ttl:           int64(oldTtl.(int)),
				Rrdatas:       make([]string, oldCount),
				RoutingPolicy: oldRoutingPolicy,
			},
		},
		Additions: []*dns.ResourceRecordSet{
			{
				Name:          recordName,
				Type:          newType.(string),
				Ttl:           int64(newTtl.(int)),
				Rrdatas:       expandDnsRecordSetRrdata(d.Get("rrdatas").([]interface{})),
				RoutingPolicy: newRoutingPolicy,
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

func expandDnsRecordSetRrdata(configured []interface{}) []string {
	return convertStringArr(configured)
}

func expandDnsRecordSetRoutingPolicy(configured []interface{}, d TerraformResourceData, config *Config) (*dns.RRSetRoutingPolicy, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	data := configured[0].(map[string]interface{})
	wrrRawItems, _ := data["wrr"].([]interface{})
	geoRawItems, _ := data["geo"].([]interface{})
	rawPrimaryBackup, _ := data["primary_backup"].([]interface{})

	if len(wrrRawItems) > 0 {
		wrrItems, err := expandDnsRecordSetRoutingPolicyWrrItems(wrrRawItems, d, config)
		if err != nil {
			return nil, err
		}
		return &dns.RRSetRoutingPolicy{
			Wrr: &dns.RRSetRoutingPolicyWrrPolicy{
				Items: wrrItems,
			},
		}, nil
	}

	if len(geoRawItems) > 0 {
		geoItems, err := expandDnsRecordSetRoutingPolicyGeoItems(geoRawItems, d, config)
		if err != nil {
			return nil, err
		}
		return &dns.RRSetRoutingPolicy{
			Geo: &dns.RRSetRoutingPolicyGeoPolicy{
				Items:         geoItems,
				EnableFencing: data["enable_geo_fencing"].(bool),
			},
		}, nil
	}

	if len(rawPrimaryBackup) > 0 {
		primaryBackup, err := expandDnsRecordSetRoutingPolicyPrimaryBackup(rawPrimaryBackup, d, config)
		if err != nil {
			return nil, err
		}
		return &dns.RRSetRoutingPolicy{
			PrimaryBackup: primaryBackup,
		}, nil
	}

	return nil, nil // unreachable here if ps is valid data
}

func expandDnsRecordSetRoutingPolicyWrrItems(configured []interface{}, d TerraformResourceData, config *Config) ([]*dns.RRSetRoutingPolicyWrrPolicyWrrPolicyItem, error) {
	items := make([]*dns.RRSetRoutingPolicyWrrPolicyWrrPolicyItem, 0, len(configured))
	for _, raw := range configured {
		item, err := expandDnsRecordSetRoutingPolicyWrrItem(raw, d, config)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func expandDnsRecordSetRoutingPolicyWrrItem(configured interface{}, d TerraformResourceData, config *Config) (*dns.RRSetRoutingPolicyWrrPolicyWrrPolicyItem, error) {
	data := configured.(map[string]interface{})
	healthCheckedTargets, err := expandDnsRecordSetHealthCheckedTargets(data["health_checked_targets"].([]interface{}), d, config)
	if err != nil {
		return nil, err
	}
	return &dns.RRSetRoutingPolicyWrrPolicyWrrPolicyItem{
		Rrdatas:              convertStringArr(data["rrdatas"].([]interface{})),
		Weight:               data["weight"].(float64),
		HealthCheckedTargets: healthCheckedTargets,
	}, nil
}

func expandDnsRecordSetRoutingPolicyGeoItems(configured []interface{}, d TerraformResourceData, config *Config) ([]*dns.RRSetRoutingPolicyGeoPolicyGeoPolicyItem, error) {
	items := make([]*dns.RRSetRoutingPolicyGeoPolicyGeoPolicyItem, 0, len(configured))
	for _, raw := range configured {
		item, err := expandDnsRecordSetRoutingPolicyGeoItem(raw, d, config)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func expandDnsRecordSetRoutingPolicyGeoItem(configured interface{}, d TerraformResourceData, config *Config) (*dns.RRSetRoutingPolicyGeoPolicyGeoPolicyItem, error) {
	data := configured.(map[string]interface{})
	healthCheckedTargets, err := expandDnsRecordSetHealthCheckedTargets(data["health_checked_targets"].([]interface{}), d, config)
	if err != nil {
		return nil, err
	}
	return &dns.RRSetRoutingPolicyGeoPolicyGeoPolicyItem{
		Rrdatas:              convertStringArr(data["rrdatas"].([]interface{})),
		Location:             data["location"].(string),
		HealthCheckedTargets: healthCheckedTargets,
	}, nil
}

func expandDnsRecordSetHealthCheckedTargets(configured []interface{}, d TerraformResourceData, config *Config) (*dns.RRSetRoutingPolicyHealthCheckTargets, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	data := configured[0].(map[string]interface{})
	internalLoadBalancers, err := expandDnsRecordSetHealthCheckedTargetsInternalLoadBalancers(data["internal_load_balancers"].([]interface{}), d, config)
	if err != nil {
		return nil, err
	}
	return &dns.RRSetRoutingPolicyHealthCheckTargets{
		InternalLoadBalancers: internalLoadBalancers,
	}, nil
}

func expandDnsRecordSetHealthCheckedTargetsInternalLoadBalancers(configured []interface{}, d TerraformResourceData, config *Config) ([]*dns.RRSetRoutingPolicyLoadBalancerTarget, error) {
	ilbs := make([]*dns.RRSetRoutingPolicyLoadBalancerTarget, 0, len(configured))
	for _, raw := range configured {
		ilb, err := expandDnsRecordSetHealthCheckedTargetsInternalLoadBalancer(raw, d, config)
		if err != nil {
			return nil, err
		}
		ilbs = append(ilbs, ilb)
	}
	return ilbs, nil
}

func expandDnsRecordSetHealthCheckedTargetsInternalLoadBalancer(configured interface{}, d TerraformResourceData, config *Config) (*dns.RRSetRoutingPolicyLoadBalancerTarget, error) {
	data := configured.(map[string]interface{})
	networkUrl, err := expandDnsRecordSetHealthCheckedTargetsInternalLoadBalancerNetworkUrl(data["network_url"], d, config)
	if err != nil {
		return nil, err
	}
	return &dns.RRSetRoutingPolicyLoadBalancerTarget{
		LoadBalancerType: data["load_balancer_type"].(string),
		IpAddress:        data["ip_address"].(string),
		Port:             data["port"].(string),
		IpProtocol:       data["ip_protocol"].(string),
		NetworkUrl:       networkUrl.(string),
		Project:          data["project"].(string),
		Region:           data["region"].(string),
	}, nil
}

func expandDnsRecordSetHealthCheckedTargetsInternalLoadBalancerNetworkUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	if v == nil || v.(string) == "" {
		return "", nil
	} else if strings.HasPrefix(v.(string), "https://") {
		return v, nil
	}
	url, err := replaceVars(d, config, "{{ComputeBasePath}}"+v.(string))
	if err != nil {
		return "", err
	}
	return ConvertSelfLinkToV1(url), nil
}

func expandDnsRecordSetRoutingPolicyPrimaryBackup(configured []interface{}, d TerraformResourceData, config *Config) (*dns.RRSetRoutingPolicyPrimaryBackupPolicy, error) {
	if len(configured) == 0 || configured[0] == nil {
		return nil, nil
	}

	data := configured[0].(map[string]interface{})

	geoRawItems, _ := data["backup_geo"].([]interface{})

	primaryTargets, err := expandDnsRecordSetHealthCheckedTargets(data["primary"].([]interface{}), d, config)
	if err != nil {
		return nil, err
	}

	items, err := expandDnsRecordSetRoutingPolicyGeoItems(geoRawItems, d, config)
	if err != nil {
		return nil, err
	}

	return &dns.RRSetRoutingPolicyPrimaryBackupPolicy{
		PrimaryTargets: primaryTargets,
		TrickleTraffic: data["trickle_ratio"].(float64),
		BackupGeoTargets: &dns.RRSetRoutingPolicyGeoPolicy{
			Items:         items,
			EnableFencing: data["enable_geo_fencing_for_backups"].(bool),
		},
	}, nil
}

func flattenDnsRecordSetRoutingPolicy(policy *dns.RRSetRoutingPolicy) []interface{} {
	if policy == nil {
		return []interface{}{}
	}
	ps := make([]interface{}, 0, 1)
	p := make(map[string]interface{})
	if policy.Wrr != nil {
		p["wrr"] = flattenDnsRecordSetRoutingPolicyWRR(policy.Wrr)
	}
	if policy.Geo != nil {
		p["geo"] = flattenDnsRecordSetRoutingPolicyGEO(policy.Geo)
		p["enable_geo_fencing"] = policy.Geo.EnableFencing
	}
	if policy.PrimaryBackup != nil {
		p["primary_backup"] = flattenDnsRecordSetRoutingPolicyPrimaryBackup(policy.PrimaryBackup)
	}
	return append(ps, p)
}

func flattenDnsRecordSetRoutingPolicyWRR(wrr *dns.RRSetRoutingPolicyWrrPolicy) []interface{} {
	ris := make([]interface{}, 0, len(wrr.Items))
	for _, item := range wrr.Items {
		ri := make(map[string]interface{})
		ri["weight"] = item.Weight
		ri["rrdatas"] = item.Rrdatas
		ri["health_checked_targets"] = flattenDnsRecordSetHealthCheckedTargets(item.HealthCheckedTargets)
		ris = append(ris, ri)
	}
	return ris
}

func flattenDnsRecordSetRoutingPolicyGEO(geo *dns.RRSetRoutingPolicyGeoPolicy) []interface{} {
	ris := make([]interface{}, 0, len(geo.Items))
	for _, item := range geo.Items {
		ri := make(map[string]interface{})
		ri["location"] = item.Location
		ri["rrdatas"] = item.Rrdatas
		ri["health_checked_targets"] = flattenDnsRecordSetHealthCheckedTargets(item.HealthCheckedTargets)
		ris = append(ris, ri)
	}
	return ris
}

func flattenDnsRecordSetHealthCheckedTargets(targets *dns.RRSetRoutingPolicyHealthCheckTargets) []map[string]interface{} {
	if targets == nil {
		return nil
	}

	data := map[string]interface{}{
		"internal_load_balancers": flattenDnsRecordSetInternalLoadBalancers(targets.InternalLoadBalancers),
	}

	return []map[string]interface{}{data}
}

func flattenDnsRecordSetInternalLoadBalancers(ilbs []*dns.RRSetRoutingPolicyLoadBalancerTarget) []map[string]interface{} {
	ilbsSchema := make([]map[string]interface{}, 0, len(ilbs))
	for _, ilb := range ilbs {
		data := map[string]interface{}{
			"load_balancer_type": ilb.LoadBalancerType,
			"ip_address":         ilb.IpAddress,
			"port":               ilb.Port,
			"ip_protocol":        ilb.IpProtocol,
			"network_url":        ilb.NetworkUrl,
			"project":            ilb.Project,
			"region":             ilb.Region,
		}
		ilbsSchema = append(ilbsSchema, data)
	}
	return ilbsSchema
}

func flattenDnsRecordSetRoutingPolicyPrimaryBackup(primaryBackup *dns.RRSetRoutingPolicyPrimaryBackupPolicy) []map[string]interface{} {
	if primaryBackup == nil {
		return nil
	}

	data := map[string]interface{}{
		"primary":                        flattenDnsRecordSetHealthCheckedTargets(primaryBackup.PrimaryTargets),
		"trickle_ratio":                  primaryBackup.TrickleTraffic,
		"backup_geo":                     flattenDnsRecordSetRoutingPolicyGEO(primaryBackup.BackupGeoTargets),
		"enable_geo_fencing_for_backups": primaryBackup.BackupGeoTargets.EnableFencing,
	}

	return []map[string]interface{}{data}
}

func validateRecordNameTrailingDot(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	len_value := len(value)
	if len_value == 0 {
		errors = append(errors, fmt.Errorf("the empty string is not a valid name field value"))
		return nil, errors
	}
	last1 := value[len_value-1:]
	if last1 != "." {
		errors = append(errors, fmt.Errorf("%q (%q) doesn't end with %q, name field must end with trailing dot, for example test.example.com. (note the trailing dot)", k, value, "."))
		return nil, errors
	}
	return nil, nil
}
