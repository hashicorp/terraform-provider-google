package google

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"

	"google.golang.org/api/compute/v1"
)

const peerNetworkLinkRegex = "projects/(" + ProjectRegex + ")/global/networks/((?:[a-z](?:[-a-z0-9]*[a-z0-9])?))$"

func ResourceComputeNetworkPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkPeeringCreate,
		Read:   resourceComputeNetworkPeeringRead,
		Update: resourceComputeNetworkPeeringUpdate,
		Delete: resourceComputeNetworkPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeNetworkPeeringImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCEName,
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
				Optional:    true,
				Default:     false,
				Description: `Whether to export the custom routes to the peer network. Defaults to false.`,
			},

			"import_custom_routes": {
				Type:        schema.TypeBool,
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
		},
		UseJSONNumber: true,
	}
}

func resourceComputeNetworkPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

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

	addOp, err := config.NewComputeClient(userAgent).Networks.AddPeering(networkFieldValue.Project, networkFieldValue.Name, request).Do()
	if err != nil {
		return fmt.Errorf("Error adding network peering: %s", err)
	}

	err = ComputeOperationWaitTime(config, addOp, networkFieldValue.Project, "Adding Network Peering", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", networkFieldValue.Name, d.Get("name").(string)))

	return resourceComputeNetworkPeeringRead(d, meta)
}

func resourceComputeNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	peeringName := d.Get("name").(string)
	networkFieldValue, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}

	network, err := config.NewComputeClient(userAgent).Networks.Get(networkFieldValue.Project, networkFieldValue.Name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Network %q", networkFieldValue.Name))
	}

	peering := findPeeringFromNetwork(network, peeringName)
	if peering == nil {
		log.Printf("[WARN] Removing network peering %s from network %s because it's gone", peeringName, network.Name)
		d.SetId("")
		return nil
	}

	if err := d.Set("peer_network", peering.Network); err != nil {
		return fmt.Errorf("Error setting peer_network: %s", err)
	}
	if err := d.Set("name", peering.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("import_custom_routes", peering.ImportCustomRoutes); err != nil {
		return fmt.Errorf("Error setting import_custom_routes: %s", err)
	}
	if err := d.Set("export_custom_routes", peering.ExportCustomRoutes); err != nil {
		return fmt.Errorf("Error setting export_custom_routes: %s", err)
	}
	if err := d.Set("import_subnet_routes_with_public_ip", peering.ImportSubnetRoutesWithPublicIp); err != nil {
		return fmt.Errorf("Error setting import_subnet_routes_with_public_ip: %s", err)
	}
	if err := d.Set("export_subnet_routes_with_public_ip", peering.ExportSubnetRoutesWithPublicIp); err != nil {
		return fmt.Errorf("Error setting export_subnet_routes_with_public_ip: %s", err)
	}
	if err := d.Set("state", peering.State); err != nil {
		return fmt.Errorf("Error setting state: %s", err)
	}
	if err := d.Set("state_details", peering.StateDetails); err != nil {
		return fmt.Errorf("Error setting state_details: %s", err)
	}

	return nil
}

func resourceComputeNetworkPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	networkFieldValue, err := ParseNetworkFieldValue(d.Get("network").(string), d, config)
	if err != nil {
		return err
	}
	peerNetworkFieldValue, err := ParseNetworkFieldValue(d.Get("peer_network").(string), d, config)
	if err != nil {
		return err
	}

	request := &compute.NetworksUpdatePeeringRequest{}
	request.NetworkPeering = expandNetworkPeering(d)

	// Only one peering operation at a time can be performed for a given network.
	// Lock on both networks, sorted so we don't deadlock for A <--> B peering pairs.
	peeringLockNames := sortedNetworkPeeringMutexKeys(networkFieldValue, peerNetworkFieldValue)
	for _, kn := range peeringLockNames {
		mutexKV.Lock(kn)
		defer mutexKV.Unlock(kn)
	}

	updateOp, err := config.NewComputeClient(userAgent).Networks.UpdatePeering(networkFieldValue.Project, networkFieldValue.Name, request).Do()
	if err != nil {
		return fmt.Errorf("Error updating network peering: %s", err)
	}

	err = ComputeOperationWaitTime(config, updateOp, networkFieldValue.Project, "Updating Network Peering", userAgent, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return resourceComputeNetworkPeeringRead(d, meta)
}

func resourceComputeNetworkPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

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

	removeOp, err := config.NewComputeClient(userAgent).Networks.RemovePeering(networkFieldValue.Project, networkFieldValue.Name, request).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Peering `%s` already removed from network `%s`", name, networkFieldValue.Name)
		} else {
			return fmt.Errorf("Error removing peering `%s` from network `%s`: %s", name, networkFieldValue.Name, err)
		}
	} else {
		err = ComputeOperationWaitTime(config, removeOp, networkFieldValue.Project, "Removing Network Peering", userAgent, d.Timeout(schema.TimeoutDelete))
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
		ForceSendFields:                []string{"ExportSubnetRoutesWithPublicIp", "ImportCustomRoutes", "ExportCustomRoutes"},
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
	project := splits[0]
	network := splits[1]
	name := splits[2]

	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	// Since the format of the network URL in the peering might be different depending on the ComputeBasePath,
	// just read the network self link from the API.
	net, err := config.NewComputeClient(userAgent).Networks.Get(project, network).Do()
	if err != nil {
		return nil, handleNotFoundError(err, d, fmt.Sprintf("Network %q", splits[1]))
	}

	if err := d.Set("network", ConvertSelfLinkToV1(net.SelfLink)); err != nil {
		return nil, fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("name", name); err != nil {
		return nil, fmt.Errorf("Error setting name: %s", err)
	}

	// Replace import id for the resource id
	id := fmt.Sprintf("%s/%s", network, name)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
