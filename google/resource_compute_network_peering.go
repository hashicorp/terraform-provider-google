package google

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const peerNetworkLinkRegex = "projects/(" + ProjectRegex + ")/global/networks/((?:[a-z](?:[-a-z0-9]*[a-z0-9])?))$"

func resourceComputeNetworkPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkPeeringCreate,
		Read:   resourceComputeNetworkPeeringRead,
		Delete: resourceComputeNetworkPeeringDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateRegexp(peerNetworkLinkRegex),
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},
			"peer_network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateRegexp(peerNetworkLinkRegex),
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},
			// The API only accepts true as a value for exchange_subnet_routes or auto_create_routes (of which only one can be set in a valid request).
			// Also, you can't set auto_create_routes if you use the networkPeering object. auto_create_routes is also removed
			"auto_create_routes": {
				Type:     schema.TypeBool,
				Optional: true,
				Removed:  "auto_create_routes has been removed because it's redundant and not user-configurable. It can safely be removed from your config",
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_details": {
				Type:     schema.TypeString,
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

	request := &computeBeta.NetworksAddPeeringRequest{}
	request.NetworkPeering = expandNetworkPeering(d)

	addOp, err := config.clientComputeBeta.Networks.AddPeering(networkFieldValue.Project, networkFieldValue.Name, request).Do()
	if err != nil {
		return fmt.Errorf("Error adding network peering: %s", err)
	}

	err = computeSharedOperationWait(config.clientCompute, addOp, networkFieldValue.Project, "Adding Network Peering")
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

	network, err := config.clientComputeBeta.Networks.Get(networkFieldValue.Project, networkFieldValue.Name).Do()
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

	// Only one delete peering operation at a time can be performed inside any peered VPCs.
	peeringLockName := getNetworkPeeringLockName(networkFieldValue.Name, peerNetworkFieldValue.Name)
	mutexKV.Lock(peeringLockName)
	defer mutexKV.Unlock(peeringLockName)

	removeOp, err := config.clientCompute.Networks.RemovePeering(networkFieldValue.Project, networkFieldValue.Name, request).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Peering `%s` already removed from network `%s`", name, networkFieldValue.Name)
		} else {
			return fmt.Errorf("Error removing peering `%s` from network `%s`: %s", name, networkFieldValue.Name, err)
		}
	} else {
		err = computeSharedOperationWait(config.clientCompute, removeOp, networkFieldValue.Project, "Removing Network Peering")
		if err != nil {
			return err
		}
	}

	return nil
}

func findPeeringFromNetwork(network *computeBeta.Network, peeringName string) *computeBeta.NetworkPeering {
	for _, p := range network.Peerings {
		if p.Name == peeringName {
			return p
		}
	}
	return nil
}
func expandNetworkPeering(d *schema.ResourceData) *computeBeta.NetworkPeering {
	return &computeBeta.NetworkPeering{
		// auto_create_routes was replaced by exchange_subnet_routes in the network peering object
		ExchangeSubnetRoutes: true,
		Name:                 d.Get("name").(string),
		Network:              d.Get("peer_network").(string),
	}
}

func getNetworkPeeringLockName(networkName, peerNetworkName string) string {
	// Whether you delete the peering from network A to B or the one from B to A, they
	// cannot happen at the same time.
	networks := []string{networkName, peerNetworkName}
	sort.Strings(networks)

	return fmt.Sprintf("network_peering/%s/%s", networks[0], networks[1])
}
