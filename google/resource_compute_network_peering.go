package google

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const peerNetworkLinkRegex = "projects/(" + ProjectRegex + ")/global/networks/((?:[a-z](?:[-a-z0-9]*[a-z0-9])?))$"

func resourceComputeNetworkPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkPeeringCreate,
		Read:   resourceComputeNetworkPeeringRead,
		Delete: resourceComputeNetworkPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeNetworkPeeringImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
				Description:  `Name of the peering.`,
			},

			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateRegexp(peerNetworkLinkRegex),
				DiffSuppressFunc: compareSelfLinkRelativePaths,
				Description:      `The primary network of the peering.`,
			},

			"peer_network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateRegexp(peerNetworkLinkRegex),
				DiffSuppressFunc: compareSelfLinkRelativePaths,
				Description:      `The peer network in the peering. The peer network may belong to a different project.`,
			},

			"export_custom_routes": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: `Whether to export the custom routes to the peer network. Defaults to false.`,
			},

			"import_custom_routes": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: `Whether to export the custom routes from the peer network. Defaults to false.`,
			},

			"export_subnet_routes_with_public_ip": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
			},

			"import_subnet_routes_with_public_ip": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `State for the peering, either ACTIVE or INACTIVE. The peering is ACTIVE when there's a matching configuration in the peer network.`,
			},

			"state_details": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Details about the current state of the peering.`,
			},

			"auto_create_routes": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "auto_create_routes has been removed because it's redundant and not user-configurable. It can safely be removed from your config",
				Computed: true,
			},
		},
	}
}

func resourceComputeNetworkPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkFieldValue, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}
	peerNetworkFieldValue, err := ParseNetworkFieldValue(d.Get("peer_network").(string), d, config)
	if err != nil {
		return err
	}

	request := &compute.NetworksAddPeeringRequest{}
	request.NetworkPeering = expandNetworkPeering(d)

	// Only one peering operation at a time can be performed for a given network.
	// Lock on both networks, sorted so we don't deadlock for A <--> B peering pairs.
	peeringLockNames := sortedNetworkPeeringMutexKeys(networkFieldValue, peerNetworkFieldValue)
	for _, kn := range peeringLockNames {
		mutexKV.Lock(kn)
		defer mutexKV.Unlock(kn)
	}

	addOp, err := config.clientCompute.Networks.AddPeering(networkFieldValue.Project, networkFieldValue.Name, request).Do()
	if err != nil {
		return fmt.Errorf("Error adding network peering: %s", err)
	}

	err = computeOperationWaitTime(config, addOp, networkFieldValue.Project, "Adding Network Peering", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", networkFieldValue.Name, d.Get("name").(string)))

	return resourceComputeNetworkPeeringRead(d, meta)
}

func resourceComputeNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	peeringName := d.Get("name").(string)
	networkFieldValue, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}

	network, err := config.clientCompute.Networks.Get(networkFieldValue.Project, networkFieldValue.Name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Network %q", networkFieldValue.Name))
	}

	peering := findPeeringFromNetwork(network, peeringName)
	if peering == nil {
		log.Printf("[WARN] Removing network peering %s from network %s because it's gone", peeringName, network.Name)
		d.SetId("")
		return nil
	}

	d.Set("peer_network", peering.Network)
	d.Set("name", peering.Name)
	d.Set("import_custom_routes", peering.ImportCustomRoutes)
	d.Set("export_custom_routes", peering.ExportCustomRoutes)
	d.Set("import_subnet_routes_with_public_ip", peering.ImportSubnetRoutesWithPublicIp)
	d.Set("export_subnet_routes_with_public_ip", peering.ExportSubnetRoutesWithPublicIp)
	d.Set("state", peering.State)
	d.Set("state_details", peering.StateDetails)

	return nil
}

func resourceComputeNetworkPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Remove the `network` to `peer_network` peering
	name := d.Get("name").(string)
	networkFieldValue, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}
	peerNetworkFieldValue, err := ParseNetworkFieldValue(d.Get("peer_network").(string), d, config)
	if err != nil {
		return err
	}

	request := &compute.NetworksRemovePeeringRequest{
		Name: name,
	}

	// Only one peering operation at a time can be performed for a given network.
	// Lock on both networks, sorted so we don't deadlock for A <--> B peering pairs.
	peeringLockNames := sortedNetworkPeeringMutexKeys(networkFieldValue, peerNetworkFieldValue)
	for _, kn := range peeringLockNames {
		mutexKV.Lock(kn)
		defer mutexKV.Unlock(kn)
	}

	removeOp, err := config.clientCompute.Networks.RemovePeering(networkFieldValue.Project, networkFieldValue.Name, request).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Peering `%s` already removed from network `%s`", name, networkFieldValue.Name)
		} else {
			return fmt.Errorf("Error removing peering `%s` from network `%s`: %s", name, networkFieldValue.Name, err)
		}
	} else {
		err = computeOperationWaitTime(config, removeOp, networkFieldValue.Project, "Removing Network Peering", d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return err
		}
	}

	return nil
}

func findPeeringFromNetwork(network *compute.Network, peeringName string) *compute.NetworkPeering {
	for _, p := range network.Peerings {
		if p.Name == peeringName {
			return p
		}
	}
	return nil
}
func expandNetworkPeering(d *schema.ResourceData) *compute.NetworkPeering {
	return &compute.NetworkPeering{
		ExchangeSubnetRoutes:           true,
		Name:                           d.Get("name").(string),
		Network:                        d.Get("peer_network").(string),
		ExportCustomRoutes:             d.Get("export_custom_routes").(bool),
		ImportCustomRoutes:             d.Get("import_custom_routes").(bool),
		ExportSubnetRoutesWithPublicIp: d.Get("export_subnet_routes_with_public_ip").(bool),
		ImportSubnetRoutesWithPublicIp: d.Get("import_subnet_routes_with_public_ip").(bool),
		ForceSendFields:                []string{"ExportSubnetRoutesWithPublicIp"},
	}
}

func sortedNetworkPeeringMutexKeys(networkName, peerNetworkName *GlobalFieldValue) []string {
	// Whether you delete the peering from network A to B or the one from B to A, they
	// cannot happen at the same time.
	networks := []string{
		fmt.Sprintf("%s/peerings", networkName.RelativeLink()),
		fmt.Sprintf("%s/peerings", peerNetworkName.RelativeLink()),
	}
	sort.Strings(networks)
	return networks
}

func resourceComputeNetworkPeeringImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	splits := strings.Split(d.Id(), "/")
	if len(splits) != 3 {
		return nil, fmt.Errorf("Error parsing network peering import format, expected: {project}/{network}/{name}")
	}

	// Build the template for the network self_link
	urlTemplate, err := replaceVars(d, config, "{{ComputeBasePath}}projects/%s/global/networks/%s")
	if err != nil {
		return nil, err
	}
	d.Set("network", ConvertSelfLinkToV1(fmt.Sprintf(urlTemplate, splits[0], splits[1])))
	d.Set("name", splits[2])

	// Replace import id for the resource id
	id := fmt.Sprintf("%s/%s", splits[1], splits[2])
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
