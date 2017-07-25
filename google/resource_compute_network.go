package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"regexp"
)

const peerNetworkLinkRegex = "projects/(" + projectRegex + ")/global/networks/((?:[a-z](?:[-a-z0-9]*[a-z0-9])?))$"

func resourceComputeNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkCreate,
		Read:   resourceComputeNetworkRead,
		Update: resourceComputeNetworkUpdate,
		Delete: resourceComputeNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auto_create_subnetworks": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				/* Ideally this would default to true as per the API, but that would cause
				   existing Terraform configs which have not been updated to report this as
				   a change. Perhaps we can bump this for a minor release bump rather than
				   a point release.
				Default: false, */
				ConflictsWith: []string{"ipv4_range"},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"gateway_ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_range": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use google_compute_subnetwork resources instead.",
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"peering": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateGCPName,
						},
						"network": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validateRegexp(peerNetworkLinkRegex),
							DiffSuppressFunc: peerNetworkLinkDiffSuppress,
						},
						"auto_create_routes": &schema.Schema{
							Type:     schema.TypeBool,
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
				},
			},
		},
	}
}

func resourceComputeNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	//
	// Possible modes:
	// - 1 Legacy mode - Create a network in the legacy mode. ipv4_range is set. auto_create_subnetworks must not be
	//     set (enforced by ConflictsWith schema attribute)
	// - 2 Distributed Mode - Create a new generation network that supports subnetworks:
	//   - 2.a - Auto subnet mode - auto_create_subnetworks = true, Google will generate 1 subnetwork per region
	//   - 2.b - Custom subnet mode - auto_create_subnetworks = false & ipv4_range not set,
	//
	autoCreateSubnetworks := d.Get("auto_create_subnetworks").(bool)

	// Build the network parameter
	network := &compute.Network{
		Name: d.Get("name").(string),
		AutoCreateSubnetworks: autoCreateSubnetworks,
		Description:           d.Get("description").(string),
	}

	if v, ok := d.GetOk("ipv4_range"); ok {
		log.Printf("[DEBUG] Setting IPv4Range (%#v) for legacy network mode", v.(string))
		network.IPv4Range = v.(string)
	} else {
		// custom subnet mode, so make sure AutoCreateSubnetworks field is included in request otherwise
		// google will create a network in legacy mode.
		network.ForceSendFields = []string{"AutoCreateSubnetworks"}
	}

	log.Printf("[DEBUG] Network insert request: %#v", network)
	op, err := config.clientCompute.Networks.Insert(
		project, network).Do()
	if err != nil {
		return fmt.Errorf("Error creating network: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(network.Name)

	err = computeOperationWaitGlobal(config, op, project, "Creating Network")
	if err != nil {
		return err
	}

	// Add peering. Network peerings cannot be added using the Insert method.
	if d.Get("peering.#").(int) > 0 {
		err = addPeering(config, project, network.Name, convertSchemaArrayToMap(d.Get("peering").([]interface{})))
		if err != nil {
			return err
		}
	}

	return resourceComputeNetworkRead(d, meta)
}

func resourceComputeNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	network, err := config.clientCompute.Networks.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Network %q", d.Get("name").(string)))
	}

	d.Set("peering", flattenNetworkPeerings(network.Peerings))
	d.Set("gateway_ipv4", network.GatewayIPv4)
	d.Set("self_link", network.SelfLink)
	d.Set("ipv4_range", network.IPv4Range)
	d.Set("name", network.Name)
	d.Set("auto_create_subnetworks", network.AutoCreateSubnetworks)

	return nil
}

func resourceComputeNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	network := d.Get("name").(string)

	d.Partial(true)

	if d.HasChange("peering") {
		old, new := d.GetChange("peering")
		add, remove := calcAddRemoveNetworkPeerings(old, new)

		if len(remove) > 0 {
			for _, peering := range remove {
				request := &compute.NetworksRemovePeeringRequest{
					Name: peering["name"].(string),
				}

				addOp, err := config.clientCompute.Networks.RemovePeering(project, network, request).Do()
				if err != nil {
					if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
						log.Printf("[WARN] Peering `%s` already removed from Network `%s`", peering["name"], network)
					} else {
						return fmt.Errorf("Error removing peering `%s` from Network `%s`", peering["name"], network)
					}
				} else {
					err = computeOperationWaitGlobal(config, addOp, project, "Updating Network")
					if err != nil {
						return err
					}
				}
			}
		}

		err := addPeering(config, project, network, add)
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceComputeNetworkRead(d, meta)
}

func resourceComputeNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the network
	op, err := config.clientCompute.Networks.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting network: %s", err)
	}

	err = computeOperationWaitGlobal(config, op, project, "Deleting Network")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func flattenNetworkPeerings(peerings []*compute.NetworkPeering) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(peerings))
	for _, peering := range peerings {
		peeringMap := make(map[string]interface{})
		peeringMap["name"] = peering.Name
		peeringMap["network"] = peering.Network
		peeringMap["auto_create_routes"] = peering.AutoCreateRoutes
		peeringMap["state"] = peering.State
		peeringMap["state_details"] = peering.StateDetails

		result = append(result, peeringMap)
	}
	return result
}

func calcAddRemoveNetworkPeerings(old_ interface{}, new_ interface{}) ([]map[string]interface{}, []map[string]interface{}) {
	old := convertSchemaArrayToMap(old_.([]interface{}))
	new := convertSchemaArrayToMap(new_.([]interface{}))

	add := make([]map[string]interface{}, 0)
	remove := make([]map[string]interface{}, 0)

	for _, newPeering := range new {
		found := false
		for _, oldPeering := range old {
			if newPeering["name"] == oldPeering["name"] {
				found = true
				if newPeering["network"] != oldPeering["network"] || newPeering["auto_create_routes"] != oldPeering["auto_create_routes"] {
					// Update to the network peering, we must delete the old one and create the new one.
					remove = append(remove, oldPeering)
					add = append(add, newPeering)
				}
				break
			}
		}
		if !found {
			add = append(add, newPeering)
		}
	}

	for _, oldPeering := range old {
		found := false
		for _, newPeering := range new {
			if newPeering["name"] == oldPeering["name"] {
				found = true
				break
			}
		}
		if !found {
			remove = append(remove, oldPeering)
		}
	}

	return add, remove
}

func addPeering(config *Config, project, network string, add []map[string]interface{}) error {
	for _, peering := range add {
		request := &compute.NetworksAddPeeringRequest{
			Name:             peering["name"].(string),
			PeerNetwork:      peering["network"].(string),
			AutoCreateRoutes: peering["auto_create_routes"].(bool),
		}

		addOp, err := config.clientCompute.Networks.AddPeering(project, network, request).Do()
		if err != nil {
			return fmt.Errorf("Error adding peerings to network: %s", err)
		}

		err = computeOperationWaitGlobal(config, addOp, project, "Updating Network")
		if err != nil {
			return err
		}
	}

	return nil
}

func peerNetworkLinkDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	r := regexp.MustCompile(peerNetworkLinkRegex)

	m := r.FindStringSubmatch(old)
	if len(m) != 3 {
		return false
	}
	oldProject, oldPeeringNetworkName := m[1], m[2]

	m = r.FindStringSubmatch(new)
	if len(m) != 3 {
		return false
	}
	newProject, newPeeringNetworkName := m[1], m[2]

	if oldProject == newProject && oldPeeringNetworkName == newPeeringNetworkName {
		return true
	}
	return false
}
