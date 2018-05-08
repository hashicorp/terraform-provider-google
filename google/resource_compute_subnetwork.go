package google

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform/helper/customdiff"
	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func resourceComputeSubnetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSubnetworkCreate,
		Read:   resourceComputeSubnetworkRead,
		Update: resourceComputeSubnetworkUpdate,
		Delete: resourceComputeSubnetworkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeSubnetworkImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"ip_cidr_range": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIpCidrRange,
				// ForceNew only if it shrinks the CIDR range, this is set in CustomizeDiff below.
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"private_ip_google_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"secondary_ip_range": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateGCPName,
						},
						"ip_cidr_range": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"enable_flow_logs": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("ip_cidr_range", isShrinkageIpCidr),
		),
	}
}

func resourceComputeSubnetworkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	network, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the subnetwork parameters
	subnetwork := &computeBeta.Subnetwork{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		IpCidrRange:           d.Get("ip_cidr_range").(string),
		PrivateIpGoogleAccess: d.Get("private_ip_google_access").(bool),
		SecondaryIpRanges:     expandSecondaryRangesV0Beta(d.Get("secondary_ip_range").([]interface{})),
		Network:               network.RelativeLink(),
		EnableFlowLogs:        d.Get("enable_flow_logs").(bool),
	}

	log.Printf("[DEBUG] Subnetwork insert request: %#v", subnetwork)

	op, err := config.clientComputeBeta.Subnetworks.Insert(project, region, subnetwork).Do()
	if err != nil {
		return fmt.Errorf("Error creating subnetwork: %s", err)
	}

	// It probably maybe worked, so store the ID now. ID is a combination of region + subnetwork
	// name because subnetwork names are not unique in a project, per the Google docs:
	// "When creating a new subnetwork, its name has to be unique in that project for that region, even across networks.
	// The same name can appear twice in a project, as long as each one is in a different region."
	// https://cloud.google.com/compute/docs/subnetworks
	subnetwork.Region = region
	d.SetId(createSubnetIDBeta(subnetwork))

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutCreate).Minutes()), "Creating Subnetwork")
	if err != nil {
		return err
	}

	return resourceComputeSubnetworkRead(d, meta)
}

func resourceComputeSubnetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	subnetwork, err := config.clientComputeBeta.Subnetworks.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Subnetwork %q", name))
	}

	d.Set("name", subnetwork.Name)
	d.Set("ip_cidr_range", subnetwork.IpCidrRange)
	d.Set("network", subnetwork.Network)
	d.Set("description", subnetwork.Description)
	d.Set("private_ip_google_access", subnetwork.PrivateIpGoogleAccess)
	d.Set("gateway_address", subnetwork.GatewayAddress)
	d.Set("secondary_ip_range", flattenSecondaryRangesV0Beta(subnetwork.SecondaryIpRanges))
	d.Set("project", project)
	d.Set("region", region)
	d.Set("enable_flow_logs", subnetwork.EnableFlowLogs)
	d.Set("self_link", ConvertSelfLinkToV1(subnetwork.SelfLink))
	d.Set("fingerprint", subnetwork.Fingerprint)

	return nil
}

func resourceComputeSubnetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("private_ip_google_access") {
		subnetworksSetPrivateIpGoogleAccessRequest := &compute.SubnetworksSetPrivateIpGoogleAccessRequest{
			PrivateIpGoogleAccess: d.Get("private_ip_google_access").(bool),
		}

		log.Printf("[DEBUG] Updating Subnetwork PrivateIpGoogleAccess %q: %#v", d.Id(), subnetworksSetPrivateIpGoogleAccessRequest)

		op, err := config.clientCompute.Subnetworks.SetPrivateIpGoogleAccess(
			project, region, d.Get("name").(string), subnetworksSetPrivateIpGoogleAccessRequest).Do()

		if err != nil {
			return fmt.Errorf("Error updating subnetwork PrivateIpGoogleAccess: %s", err)
		}

		err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "Updating Subnetwork PrivateIpGoogleAccess")
		if err != nil {
			return err
		}

		d.SetPartial("private_ip_google_access")
	}

	if d.HasChange("ip_cidr_range") {
		r := &compute.SubnetworksExpandIpCidrRangeRequest{
			IpCidrRange: d.Get("ip_cidr_range").(string),
		}

		op, err := config.clientCompute.Subnetworks.ExpandIpCidrRange(project, region, d.Get("name").(string), r).Do()

		if err != nil {
			return fmt.Errorf("Error expanding the ip cidr range: %s", err)
		}

		err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "Expanding Subnetwork IP CIDR range")
		if err != nil {
			return err
		}

		d.SetPartial("ip_cidr_range")
	}

	if d.HasChange("secondary_ip_range") || d.HasChange("enable_flow_logs") {
		v0BetaSubnetwork := &computeBeta.Subnetwork{
			Fingerprint: d.Get("fingerprint").(string),
		}
		if d.HasChange("secondary_ip_range") {
			v0BetaSubnetwork.SecondaryIpRanges = expandSecondaryRangesV0Beta(d.Get("secondary_ip_range").([]interface{}))
		}
		if d.HasChange("enable_flow_logs") {
			v0BetaSubnetwork.EnableFlowLogs = d.Get("enable_flow_logs").(bool)
			v0BetaSubnetwork.ForceSendFields = append(v0BetaSubnetwork.ForceSendFields, "EnableFlowLogs")
		}

		op, err := config.clientComputeBeta.Subnetworks.Patch(
			project, region, d.Get("name").(string), v0BetaSubnetwork).Do()
		if err != nil {
			return fmt.Errorf("Error updating subnetwork %q: %s", d.Id(), err)
		}

		err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutUpdate).Minutes()), "Updating Subnetwork")
		if err != nil {
			return err
		}

		d.SetPartial("secondary_ip_range")
		d.SetPartial("enable_flow_logs")
	}

	d.Partial(false)

	return resourceComputeSubnetworkRead(d, meta)
}

func resourceComputeSubnetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the subnetwork
	op, err := config.clientCompute.Subnetworks.Delete(
		project, region, d.Get("name").(string)).Do()
	if err != nil {
		return fmt.Errorf("Error deleting subnetwork: %s", err)
	}

	err = computeSharedOperationWaitTime(config.clientCompute, op, project, int(d.Timeout(schema.TimeoutDelete).Minutes()), "Deleting Subnetwork")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceComputeSubnetworkImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid compute subnetwork specifier. Expecting {region}/{name}")
	}

	region, name := parts[0], parts[1]
	d.Set("region", region)
	d.Set("name", name)

	d.SetId(createSubnetID(&compute.Subnetwork{
		Region: region,
		Name:   name,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitSubnetID(id string) (region string, name string) {
	parts := strings.Split(id, "/")
	region = parts[0]
	name = parts[1]
	return
}

func expandSecondaryRanges(configured []interface{}) []*compute.SubnetworkSecondaryRange {
	secondaryRanges := make([]*compute.SubnetworkSecondaryRange, 0, len(configured))
	for _, raw := range configured {
		data := raw.(map[string]interface{})
		secondaryRange := compute.SubnetworkSecondaryRange{
			RangeName:   data["range_name"].(string),
			IpCidrRange: data["ip_cidr_range"].(string),
		}

		secondaryRanges = append(secondaryRanges, &secondaryRange)
	}
	return secondaryRanges
}

func expandSecondaryRangesV0Beta(configured []interface{}) []*computeBeta.SubnetworkSecondaryRange {
	secondaryRanges := make([]*computeBeta.SubnetworkSecondaryRange, 0, len(configured))
	for _, raw := range configured {
		data := raw.(map[string]interface{})
		secondaryRange := computeBeta.SubnetworkSecondaryRange{
			RangeName:   data["range_name"].(string),
			IpCidrRange: data["ip_cidr_range"].(string),
		}

		secondaryRanges = append(secondaryRanges, &secondaryRange)
	}
	return secondaryRanges
}

func flattenSecondaryRangesV0Beta(secondaryRanges []*computeBeta.SubnetworkSecondaryRange) []map[string]interface{} {
	secondaryRangesSchema := make([]map[string]interface{}, 0, len(secondaryRanges))
	for _, secondaryRange := range secondaryRanges {
		data := map[string]interface{}{
			"range_name":    secondaryRange.RangeName,
			"ip_cidr_range": secondaryRange.IpCidrRange,
		}

		secondaryRangesSchema = append(secondaryRangesSchema, data)
	}
	return secondaryRangesSchema
}

// Whether the IP CIDR change shrinks the block.
func isShrinkageIpCidr(old, new, _ interface{}) bool {
	_, oldCidr, oldErr := net.ParseCIDR(old.(string))
	_, newCidr, newErr := net.ParseCIDR(new.(string))

	if oldErr != nil || newErr != nil {
		// This should never happen. The ValidateFunc on the field ensures it.
		return false
	}

	oldStart, oldEnd := cidr.AddressRange(oldCidr)

	if newCidr.Contains(oldStart) && newCidr.Contains(oldEnd) {
		// This is a CIDR range expansion, no need to ForceNew, we have an update method for it.
		return false
	}

	return true
}
