package google

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/servicenetworking/v1"
)

func resourceServiceNetworkingConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceNetworkingConnectionCreate,
		Read:   resourceServiceNetworkingConnectionRead,
		Update: resourceServiceNetworkingConnectionUpdate,
		Delete: resourceServiceNetworkingConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServiceNetworkingConnectionImportState,
		},

		Schema: map[string]*schema.Schema{
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			// NOTE(craigatgoogle): This field is weird, it's required to make the Insert/List calls as a parameter
			// named "parent", however it's also defined in the response as an output field called "peering", which
			// uses "-" as a delimeter instead of ".". To alleviate user confusion I've opted to model the gcloud
			// CLI's approach, calling the field "service" and accepting the same format as the CLI with the "."
			// delimiter.
			// See: https://cloud.google.com/vpc/docs/configure-private-services-access#creating-connection
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reserved_peering_ranges": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceServiceNetworkingConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	network := d.Get("network").(string)
	serviceNetworkingNetworkName, err := retrieveServiceNetworkingNetworkName(d, config, network)
	if err != nil {
		return fmt.Errorf("Failed to find Service Networking Connection, err: %s", err)
	}

	connection := &servicenetworking.Connection{
		Network:               serviceNetworkingNetworkName,
		ReservedPeeringRanges: convertStringArr(d.Get("reserved_peering_ranges").([]interface{})),
	}

	parentService := formatParentService(d.Get("service").(string))
	op, err := config.clientServiceNetworking.Services.Connections.Create(parentService, connection).Do()
	if err != nil {
		return err
	}

	if err := serviceNetworkingOperationWait(config, op, "Create Service Networking Connection"); err != nil {
		return err
	}

	connectionId := &connectionId{
		Network: network,
		Service: d.Get("service").(string),
	}

	d.SetId(connectionId.Id())
	return resourceServiceNetworkingConnectionRead(d, meta)
}

func resourceServiceNetworkingConnectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	connectionId, err := parseConnectionId(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to find Service Networking Connection, err: %s", err)
	}

	serviceNetworkingNetworkName, err := retrieveServiceNetworkingNetworkName(d, config, connectionId.Network)
	if err != nil {
		return fmt.Errorf("Failed to find Service Networking Connection, err: %s", err)
	}

	parentService := formatParentService(connectionId.Service)
	listCall := config.clientServiceNetworking.Services.Connections.List(parentService)
	listCall.Network(serviceNetworkingNetworkName)
	response, err := listCall.Do()
	if err != nil {
		return err
	}

	var connection *servicenetworking.Connection
	for _, c := range response.Connections {
		if c.Network == serviceNetworkingNetworkName {
			connection = c
			break
		}
	}

	if connection == nil {
		return fmt.Errorf("Failed to find Service Networking Connection, network: %s service: %s", connectionId.Network, connectionId.Service)
	}

	d.Set("network", connectionId.Network)
	d.Set("service", connectionId.Service)
	d.Set("reserved_peering_ranges", connection.ReservedPeeringRanges)
	return nil
}

func resourceServiceNetworkingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	connectionId, err := parseConnectionId(d.Id())
	if err != nil {
		return fmt.Errorf("Failed to find Service Networking Connection, err: %s", err)
	}

	parentService := formatParentService(connectionId.Service)

	if d.HasChange("reserved_peering_ranges") {
		network := d.Get("network").(string)
		serviceNetworkingNetworkName, err := retrieveServiceNetworkingNetworkName(d, config, network)
		if err != nil {
			return fmt.Errorf("Failed to find Service Networking Connection, err: %s", err)
		}

		connection := &servicenetworking.Connection{
			Network:               serviceNetworkingNetworkName,
			ReservedPeeringRanges: convertStringArr(d.Get("reserved_peering_ranges").([]interface{})),
		}

		// The API docs don't specify that you can do connections/-, but that's what gcloud does,
		// and it's easier than grabbing the connection name.
		op, err := config.clientServiceNetworking.Services.Connections.Patch(parentService+"/connections/-", connection).UpdateMask("reservedPeeringRanges").Force(true).Do()
		if err != nil {
			return err
		}
		if err := serviceNetworkingOperationWait(config, op, "Update Service Networking Connection"); err != nil {
			return err
		}
	}
	return resourceServiceNetworkingConnectionRead(d, meta)
}

// NOTE(craigatgoogle): This resource doesn't have a defined Delete method, however an un-documented
// behavior is for the Connection to be deleted when its associated network is deleted. This is
// helpeful for acctest cleanup.
func resourceServiceNetworkingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	connectionId, err := parseConnectionId(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[WARNING] Service Networking Connection resources cannot be deleted from GCP. This Connection (network: %s, service: %s) will be removed from Terraform state, but will still be present on the server.", connectionId.Network, connectionId.Service)

	d.SetId("")

	return nil
}

func resourceServiceNetworkingConnectionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	connectionId, err := parseConnectionId(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("network", connectionId.Network)
	d.Set("service", connectionId.Service)
	return []*schema.ResourceData{d}, nil
}

// NOTE(craigatgoogle): The Connection resource in this API doesn't have an Id field, so inorder
// to support the Read method, we create an Id using the tuple(Network, Service).
type connectionId struct {
	Network string
	Service string
}

func (id *connectionId) Id() string {
	return fmt.Sprintf("%s:%s", url.QueryEscape(id.Network), url.QueryEscape(id.Service))
}

func parseConnectionId(id string) (*connectionId, error) {
	res := strings.Split(id, ":")

	if len(res) != 2 {
		return nil, fmt.Errorf("Failed to parse service networking connection id, value: %s", id)
	}

	network, err := url.QueryUnescape(res[0])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse service networking connection id, invalid network, err: %s", err)
	} else if len(network) == 0 {
		return nil, fmt.Errorf("Failed to parse service networking connection id, empty network")
	}

	service, err := url.QueryUnescape(res[1])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse service networking connection id, invalid service, err: %s", err)
	} else if len(service) == 0 {
		return nil, fmt.Errorf("Failed to parse service networking connection id, empty service")
	}

	return &connectionId{
		Network: network,
		Service: service,
	}, nil
}

// NOTE(craigatgoogle): An out of band aspect of this API is that it uses a unique formatting of network
// different from the standard self_link URI. It requires a call to the resource manager to get the project
// number for the current project.
func retrieveServiceNetworkingNetworkName(d *schema.ResourceData, config *Config, network string) (string, error) {
	networkFieldValue, err := ParseNetworkFieldValue(network, d, config)
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve network field value, err: %s", err)
	}

	pid := networkFieldValue.Project
	if pid == "" {
		return "", fmt.Errorf("Could not determine project")
	}

	project, err := config.clientResourceManager.Projects.Get(pid).Do()
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve project, pid: %s, err: %s", pid, err)
	}

	networkName := networkFieldValue.Name
	if networkName == "" {
		return "", fmt.Errorf("Failed to parse network")
	}

	// return the network name formatting unique to this API
	return fmt.Sprintf("projects/%v/global/networks/%v", project.ProjectNumber, networkName), nil

}

const parentServicePattern = "^services/.+$"

// NOTE(craigatgoogle): An out of band aspect of this API is that it requires the service name to be
// formatted as "services/<serviceName>"
func formatParentService(service string) string {
	r := regexp.MustCompile(parentServicePattern)
	if !r.MatchString(service) {
		return fmt.Sprintf("services/%s", service)
	} else {
		return service
	}
}
