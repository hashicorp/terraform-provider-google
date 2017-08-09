package google

import (
	"fmt"
	"log"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

var SubnetworkBaseApiVersion = v1
var SubnetworkVersionedFeatures = []Feature{
	{
		Version: v0beta,
		Item:    "secondary_ip_range",
	},
}

func resourceComputeSubnetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeSubnetworkCreate,
		Read:   resourceComputeSubnetworkRead,
		Update: resourceComputeSubnetworkUpdate,
		Delete: resourceComputeSubnetworkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeSubnetworkImportState,
		},

		Schema: map[string]*schema.Schema{
			"ip_cidr_range": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				DiffSuppressFunc: compareGlobalSelfLinkOrResourceName,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"gateway_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"private_ip_google_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"secondary_ip_range": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range_name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateGCPName,
						},
						"ip_cidr_range": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeSubnetworkCreate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, SubnetworkBaseApiVersion, SubnetworkVersionedFeatures)
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	network, err := getNetworkLink(d, config, "network")
	if err != nil {
		return err
	}

	// Build the subnetwork parameters
	subnetwork := &computeBeta.Subnetwork{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		IpCidrRange:           d.Get("ip_cidr_range").(string),
		PrivateIpGoogleAccess: d.Get("private_ip_google_access").(bool),
		SecondaryIpRanges:     expandSecondaryRanges(d.Get("secondary_ip_range").([]interface{})),
		Network:               network,
	}

	log.Printf("[DEBUG] Subnetwork insert request: %#v", subnetwork)

	var op interface{}
	switch computeApiVersion {
	case v1:
		subnetworkV1 := &compute.Subnetwork{}
		err := Convert(subnetwork, subnetworkV1)
		if err != nil {
			return err
		}

		op, err = config.clientCompute.Subnetworks.Insert(
			project, region, subnetworkV1).Do()
	case v0beta:
		op, err = config.clientComputeBeta.Subnetworks.Insert(
			project, region, subnetwork).Do()
	}

	if err != nil {
		return fmt.Errorf("Error creating subnetwork: %s", err)
	}

	// It probably maybe worked, so store the ID now. ID is a combination of region + subnetwork
	// name because subnetwork names are not unique in a project, per the Google docs:
	// "When creating a new subnetwork, its name has to be unique in that project for that region, even across networks.
	// The same name can appear twice in a project, as long as each one is in a different region."
	// https://cloud.google.com/compute/docs/subnetworks
	subnetwork.Region = region
	d.SetId(createBetaSubnetID(subnetwork))

	err = computeSharedOperationWait(config, op, project, "Creating Subnetwork")
	if err != nil {
		return err
	}

	return resourceComputeSubnetworkRead(d, meta)
}

func resourceComputeSubnetworkRead(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, SubnetworkBaseApiVersion, SubnetworkVersionedFeatures)
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

	subnetwork := &computeBeta.Subnetwork{}
	switch computeApiVersion {
	case v1:
		subnetworkV1, err := config.clientCompute.Subnetworks.Get(
			project, region, name).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Subnetwork %q", name))
		}

		err = Convert(subnetworkV1, subnetwork)
		if err != nil {
			return err
		}
	case v0beta:
		var err error
		subnetwork, err = config.clientComputeBeta.Subnetworks.Get(project, region, name).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Subnetwork %q", name))
		}
	}

	d.Set("name", subnetwork.Name)
	d.Set("ip_cidr_range", subnetwork.IpCidrRange)
	d.Set("network", subnetwork.Network)
	d.Set("description", subnetwork.Description)
	d.Set("private_ip_google_access", subnetwork.PrivateIpGoogleAccess)
	d.Set("gateway_address", subnetwork.GatewayAddress)
	d.Set("secondary_ip_range", flattenSecondaryRanges(subnetwork.SecondaryIpRanges))
	d.Set("self_link", ConvertSelfLinkToV1(subnetwork.SelfLink))

	return nil
}

func resourceComputeSubnetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersionUpdate(d, SubnetworkBaseApiVersion, SubnetworkVersionedFeatures, []Feature{})
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
		subnetworksSetPrivateIpGoogleAccessRequest := &computeBeta.SubnetworksSetPrivateIpGoogleAccessRequest{
			PrivateIpGoogleAccess: d.Get("private_ip_google_access").(bool),
		}

		log.Printf("[DEBUG] Updating Subnetwork PrivateIpGoogleAccess %q: %#v", d.Id(), subnetworksSetPrivateIpGoogleAccessRequest)

		var op interface{}
		switch computeApiVersion {
		case v1:
			subnetworksSetPrivateIpGoogleAccessRequestV1 := &compute.SubnetworksSetPrivateIpGoogleAccessRequest{}
			err := Convert(subnetworksSetPrivateIpGoogleAccessRequest, subnetworksSetPrivateIpGoogleAccessRequestV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.Subnetworks.SetPrivateIpGoogleAccess(
				project, region, d.Get("name").(string), subnetworksSetPrivateIpGoogleAccessRequestV1).Do()
		case v0beta:
			op, err = config.clientComputeBeta.Subnetworks.SetPrivateIpGoogleAccess(
				project, region, d.Get("name").(string), subnetworksSetPrivateIpGoogleAccessRequest).Do()

		}

		if err != nil {
			return fmt.Errorf("Error updating subnetwork PrivateIpGoogleAccess: %s", err)
		}

		err = computeSharedOperationWait(config, op, project, "Updating Subnetwork PrivateIpGoogleAccess")
		if err != nil {
			return err
		}

		d.SetPartial("private_ip_google_access")
	}

	d.Partial(false)

	return resourceComputeSubnetworkRead(d, meta)
}

func resourceComputeSubnetworkDelete(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, SubnetworkBaseApiVersion, SubnetworkVersionedFeatures)
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
	var op interface{}
	switch computeApiVersion {
	case v1:
		op, err = config.clientCompute.Subnetworks.Delete(
			project, region, d.Get("name").(string)).Do()
	case v0beta:
		op, err = config.clientComputeBeta.Subnetworks.Delete(
			project, region, d.Get("name").(string)).Do()
	}
	if err != nil {
		return fmt.Errorf("Error deleting subnetwork: %s", err)
	}

	err = computeSharedOperationWait(config, op, project, "Deleting Subnetwork")
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

func createBetaSubnetID(s *computeBeta.Subnetwork) string {
	return fmt.Sprintf("%s/%s", s.Region, s.Name)
}

func createSubnetID(s *compute.Subnetwork) string {
	return fmt.Sprintf("%s/%s", s.Region, s.Name)
}

func splitSubnetID(id string) (region string, name string) {
	parts := strings.Split(id, "/")
	region = parts[0]
	name = parts[1]
	return
}

func expandSecondaryRanges(configured []interface{}) []*computeBeta.SubnetworkSecondaryRange {
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

func flattenSecondaryRanges(secondaryRanges []*computeBeta.SubnetworkSecondaryRange) []map[string]interface{} {
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
