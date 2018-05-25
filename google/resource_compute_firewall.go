package google

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
)

func resourceComputeFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeFirewallCreate,
		Read:   resourceComputeFirewallRead,
		Update: resourceComputeFirewallUpdate,
		Delete: resourceComputeFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		MigrateState:  resourceComputeFirewallMigrateState,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1000,
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"allow": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"deny"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},

						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Set: resourceComputeFirewallRuleHash,
			},

			"deny": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"allow"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							ForceNew: true,
						},
					},
				},
				Set: resourceComputeFirewallRuleHash,

				// Unlike allow, deny can't be updated upstream
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"direction": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"INGRESS", "EGRESS"}, false),
				ForceNew:     true,
			},

			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"source_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"destination_ranges": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"source_ranges", "source_tags"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ForceNew:      true,
			},

			"target_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"source_service_accounts": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"source_tags", "target_tags"},
			},

			"target_service_accounts": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"source_tags", "target_tags"},
			},
		},
	}
}

func resourceComputeFirewallRuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["protocol"].(string)))

	// We need to make sure to sort the strings below so that we always
	// generate the same hash code no matter what is in the set.
	if v, ok := m["ports"]; ok {
		s := convertStringArr(v.([]interface{}))
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	return hashcode.String(buf.String())
}

func resourceComputeFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	firewall, err := resourceFirewall(d, meta)
	if err != nil {
		return err
	}

	op, err := config.clientComputeBeta.Firewalls.Insert(project, firewall).Do()
	if err != nil {
		return fmt.Errorf("Error creating firewall: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(firewall.Name)

	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating Firewall")
	if err != nil {
		return err
	}

	return resourceComputeFirewallRead(d, meta)
}

func flattenFirewallAllowed(allowed []*computeBeta.FirewallAllowed) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(allowed))
	for _, allow := range allowed {
		allowMap := make(map[string]interface{})
		allowMap["protocol"] = allow.IPProtocol
		allowMap["ports"] = allow.Ports

		result = append(result, allowMap)
	}
	return result
}

func flattenFirewallDenied(denied []*computeBeta.FirewallDenied) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(denied))
	for _, deny := range denied {
		denyMap := make(map[string]interface{})
		denyMap["protocol"] = deny.IPProtocol
		denyMap["ports"] = deny.Ports

		result = append(result, denyMap)
	}
	return result
}

func resourceComputeFirewallRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	firewall, err := config.clientComputeBeta.Firewalls.Get(project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Firewall %q", d.Get("name").(string)))
	}

	d.Set("self_link", ConvertSelfLinkToV1(firewall.SelfLink))
	d.Set("name", firewall.Name)
	d.Set("network", ConvertSelfLinkToV1(firewall.Network))
	d.Set("direction", firewall.Direction)
	d.Set("description", firewall.Description)
	d.Set("project", project)
	d.Set("source_ranges", firewall.SourceRanges)
	d.Set("source_tags", firewall.SourceTags)
	d.Set("destination_ranges", firewall.DestinationRanges)
	d.Set("target_tags", firewall.TargetTags)
	d.Set("allow", flattenFirewallAllowed(firewall.Allowed))
	d.Set("deny", flattenFirewallDenied(firewall.Denied))
	d.Set("priority", int(firewall.Priority))
	d.Set("source_service_accounts", firewall.SourceServiceAccounts)
	d.Set("target_service_accounts", firewall.TargetServiceAccounts)
	d.Set("disabled", firewall.Disabled)
	return nil
}

func resourceComputeFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	firewall, err := resourceFirewall(d, meta)
	if err != nil {
		return err
	}

	op, err := config.clientComputeBeta.Firewalls.Patch(project, d.Id(), firewall).Do()
	if err != nil {
		return fmt.Errorf("Error updating firewall: %s", err)
	}

	err = computeSharedOperationWait(config.clientCompute, op, project, "Updating Firewall")
	if err != nil {
		return err
	}

	d.Partial(false)

	return resourceComputeFirewallRead(d, meta)
}

func resourceComputeFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the firewall
	op, err := config.clientComputeBeta.Firewalls.Delete(project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting firewall: %s", err)
	}

	err = computeSharedOperationWait(config.clientCompute, op, project, "Deleting Firewall")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceFirewall(d *schema.ResourceData, meta interface{}) (*computeBeta.Firewall, error) {
	config := meta.(*Config)

	network, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return nil, err
	}

	// Build up the list of allowed entries
	var allowed []*computeBeta.FirewallAllowed
	if v := d.Get("allow").(*schema.Set); v.Len() > 0 {
		allowed = make([]*computeBeta.FirewallAllowed, 0, v.Len())
		for _, v := range v.List() {
			m := v.(map[string]interface{})

			allowed = append(allowed, &computeBeta.FirewallAllowed{
				IPProtocol: m["protocol"].(string),
				Ports:      convertStringArr(m["ports"].([]interface{})),
			})
		}
	}

	// Build up the list of denied entries
	var denied []*computeBeta.FirewallDenied
	if v := d.Get("deny").(*schema.Set); v.Len() > 0 {
		denied = make([]*computeBeta.FirewallDenied, 0, v.Len())
		for _, v := range v.List() {
			m := v.(map[string]interface{})

			denied = append(denied, &computeBeta.FirewallDenied{
				IPProtocol: m["protocol"].(string),
				Ports:      convertStringArr(m["ports"].([]interface{})),
			})
		}
	}

	// Build up the list of sources
	var sourceRanges, sourceTags []string
	if v := d.Get("source_ranges").(*schema.Set); v.Len() > 0 {
		sourceRanges = make([]string, v.Len())
		for i, v := range v.List() {
			sourceRanges[i] = v.(string)
		}
	}
	if v := d.Get("source_tags").(*schema.Set); v.Len() > 0 {
		sourceTags = make([]string, v.Len())
		for i, v := range v.List() {
			sourceTags[i] = v.(string)
		}
	}

	// Build up the list of destinations
	var destinationRanges []string
	if v := d.Get("destination_ranges").(*schema.Set); v.Len() > 0 {
		destinationRanges = make([]string, v.Len())
		for i, v := range v.List() {
			destinationRanges[i] = v.(string)
		}
	}

	// Build up the list of targets
	var targetTags []string
	if v := d.Get("target_tags").(*schema.Set); v.Len() > 0 {
		targetTags = make([]string, v.Len())
		for i, v := range v.List() {
			targetTags[i] = v.(string)
		}
	}

	// Build the firewall parameter
	return &computeBeta.Firewall{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Direction:             d.Get("direction").(string),
		Network:               network.RelativeLink(),
		Allowed:               allowed,
		Denied:                denied,
		SourceRanges:          sourceRanges,
		SourceTags:            sourceTags,
		DestinationRanges:     destinationRanges,
		TargetTags:            targetTags,
		Priority:              int64(d.Get("priority").(int)),
		SourceServiceAccounts: convertStringSet(d.Get("source_service_accounts").(*schema.Set)),
		TargetServiceAccounts: convertStringSet(d.Get("target_service_accounts").(*schema.Set)),
		Disabled:              d.Get("disabled").(bool),
		ForceSendFields:       []string{"Disabled"},
	}, nil
}
