package google

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateGCPName,
			},
			"network": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateRegexp(peerNetworkLinkRegex),
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},
			"peer_network": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateRegexp(peerNetworkLinkRegex),
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},
			"auto_create_routes": &schema.Schema{
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_details": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeNetworkPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	networkLink := d.Get("network").(string)
	networkName := getNameFromNetworkLink(networkLink)

	request := &compute.NetworksAddPeeringRequest{
		Name:             d.Get("name").(string),
		PeerNetwork:      d.Get("peer_network").(string),
		AutoCreateRoutes: d.Get("auto_create_routes").(bool),
	}

	addOp, err := config.clientCompute.Networks.AddPeering(project, networkName, request).Do()
	if err != nil {
		return fmt.Errorf("Error adding network peering: %s", err)
	}

	err = computeOperationWait(config, addOp, project, "Adding Network Peering")
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", networkName, d.Get("name").(string)))

	return resourceComputeNetworkPeeringRead(d, meta)
}

func resourceComputeNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	peeringName := d.Get("name").(string)
	networkLink := d.Get("network").(string)
	networkName := getNameFromNetworkLink(networkLink)

	network, err := config.clientCompute.Networks.Get(project, networkName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Network %q", networkName))
	}

	peering := findPeeringFromNetwork(network, peeringName)
	if peering == nil {
		log.Printf("[WARN] Removing network peering %s from network %s because it's gone", peeringName, networkName)
		d.SetId("")
		return nil
	}

	d.Set("peer_network", peering.Network)
	d.Set("auto_create_routes", peering.AutoCreateRoutes)
	d.Set("state", peering.State)
	d.Set("state_details", peering.StateDetails)

	return nil
}

func resourceComputeNetworkPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Remove the `network` to `peer_network` peering
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	networkLink := d.Get("network").(string)
	peerNetworkLink := d.Get("peer_network").(string)
	networkName := getNameFromNetworkLink(networkLink)
	peerNetworkName := getNameFromNetworkLink(peerNetworkLink)

	request := &compute.NetworksRemovePeeringRequest{
		Name: name,
	}

	// Only one delete peering operation at a time can be performed inside any peered VPCs.
	peeringLockName := getNetworkPeeringLockName(networkName, peerNetworkName)
	mutexKV.Lock(peeringLockName)
	defer mutexKV.Unlock(peeringLockName)

	removeOp, err := config.clientCompute.Networks.RemovePeering(project, networkName, request).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Peering `%s` already removed from network `%s`", name, networkName)
		} else {
			return fmt.Errorf("Error removing peering `%s` from network `%s`: %s", name, networkName, err)
		}
	} else {
		err = computeOperationWait(config, removeOp, project, "Removing Network Peering")
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

func getNameFromNetworkLink(network string) string {
	r := regexp.MustCompile(peerNetworkLinkRegex)

	m := r.FindStringSubmatch(network)
	return m[2]
}

func getNetworkPeeringLockName(networkName, peerNetworkName string) string {
	// Whether you delete the peering from network A to B or the one from B to A, they
	// cannot happen at the same time.
	networks := []string{networkName, peerNetworkName}
	sort.Strings(networks)

	return fmt.Sprintf("network_peering/%s/%s", networks[0], networks[1])
}
