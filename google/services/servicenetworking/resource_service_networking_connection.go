// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package servicenetworking

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/servicenetworking/v1"
)

func ResourceServiceNetworkingConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceNetworkingConnectionCreate,
		Read:   resourceServiceNetworkingConnectionRead,
		Update: resourceServiceNetworkingConnectionUpdate,
		Delete: resourceServiceNetworkingConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServiceNetworkingConnectionImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `Name of VPC network connected with service producers using VPC peering.`,
			},
			// NOTE(craigatgoogle): This field is weird, it's required to make the Insert/List calls as a parameter
			// named "parent", however it's also defined in the response as an output field called "peering", which
			// uses "-" as a delimiter instead of ".". To alleviate user confusion I've opted to model the gcloud
			// CLI's approach, calling the field "service" and accepting the same format as the CLI with the "."
			// delimiter.
			// See: https://cloud.google.com/vpc/docs/configure-private-services-access#creating-connection
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Provider peering service that is managing peering connectivity for a service provider organization. For Google services that support this functionality it is 'servicenetworking.googleapis.com'.`,
			},
			"reserved_peering_ranges": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Named IP address range(s) of PEERING type reserved for this service provider. Note that invoking this method with a different range when connection is already established will not reallocate already provisioned service producer subnetworks.`,
			},
			"peering": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceServiceNetworkingConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	network := d.Get("network").(string)
	serviceNetworkingNetworkName, err := RetrieveServiceNetworkingNetworkName(d, config, network, userAgent)
	if err != nil {
		return errwrap.Wrapf("Failed to find Service Networking Connection, err: {{err}}", err)
	}

	connection := &servicenetworking.Connection{
		Network:               serviceNetworkingNetworkName,
		ReservedPeeringRanges: tpgresource.ConvertStringArr(d.Get("reserved_peering_ranges").([]interface{})),
	}

	networkFieldValue, err := tpgresource.ParseNetworkFieldValue(network, d, config)
	if err != nil {
		return errwrap.Wrapf("Failed to retrieve network field value, err: {{err}}", err)
	}
	project := networkFieldValue.Project

	parentService := formatParentService(d.Get("service").(string))
	// We use Patch instead of Create, because we're getting
	//  "Error waiting for Create Service Networking Connection:
	//   Error code 9, message: Cannot modify allocated ranges in
	//   CreateConnection. Please use UpdateConnection."
	// if we're creating peerings to more than one VPC (like two
	// CloudSQL instances within one project, peered with two
	// clusters.)
	//
	// This is a workaround for:
	// https://issuetracker.google.com/issues/131908322
	//
	// The API docs don't specify that you can do connections/-,
	// but that's what gcloud does, and it's easier than grabbing
	// the connection name.

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		project = bp
	}

	createCall := config.NewServiceNetworkingClient(userAgent).Services.Connections.Patch(parentService+"/connections/-", connection).UpdateMask("reservedPeeringRanges").Force(true)
	if config.UserProjectOverride {
		createCall.Header().Add("X-Goog-User-Project", project)
	}
	op, err := createCall.Do()
	if err != nil {
		return err
	}

	if err := ServiceNetworkingOperationWaitTime(config, op, "Create Service Networking Connection", userAgent, project, d.Timeout(schema.TimeoutCreate)); err != nil {
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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	connectionId, err := parseConnectionId(d.Id())
	if err != nil {
		return errwrap.Wrapf("Unable to parse Service Networking Connection id, err: {{err}}", err)
	}

	serviceNetworkingNetworkName, err := RetrieveServiceNetworkingNetworkName(d, config, connectionId.Network, userAgent)
	if err != nil {
		return errwrap.Wrapf("Failed to find Service Networking Connection, err: {{err}}", err)
	}

	network := d.Get("network").(string)
	networkFieldValue, err := tpgresource.ParseNetworkFieldValue(network, d, config)
	if err != nil {
		return errwrap.Wrapf("Failed to retrieve network field value, err: {{err}}", err)
	}
	project := networkFieldValue.Project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		project = bp
	}

	parentService := formatParentService(connectionId.Service)
	readCall := config.NewServiceNetworkingClient(userAgent).Services.Connections.List(parentService).Network(serviceNetworkingNetworkName)
	if config.UserProjectOverride {
		readCall.Header().Add("X-Goog-User-Project", project)
	}
	response, err := readCall.Do()
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
		d.SetId("")
		log.Printf("[WARNING] Failed to find Service Networking Connection, network: %s service: %s", connectionId.Network, connectionId.Service)
		return nil
	}

	if err := d.Set("network", connectionId.Network); err != nil {
		return fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("service", connectionId.Service); err != nil {
		return fmt.Errorf("Error setting service: %s", err)
	}
	if err := d.Set("peering", connection.Peering); err != nil {
		return fmt.Errorf("Error setting peering: %s", err)
	}
	if err := d.Set("reserved_peering_ranges", connection.ReservedPeeringRanges); err != nil {
		return fmt.Errorf("Error setting reserved_peering_ranges: %s", err)
	}
	return nil
}

func resourceServiceNetworkingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	connectionId, err := parseConnectionId(d.Id())
	if err != nil {
		return errwrap.Wrapf("Unable to parse Service Networking Connection id, err: {{err}}", err)
	}

	parentService := formatParentService(connectionId.Service)

	if d.HasChange("reserved_peering_ranges") {
		network := d.Get("network").(string)
		serviceNetworkingNetworkName, err := RetrieveServiceNetworkingNetworkName(d, config, network, userAgent)
		if err != nil {
			return errwrap.Wrapf("Failed to find Service Networking Connection, err: {{err}}", err)
		}

		connection := &servicenetworking.Connection{
			Network:               serviceNetworkingNetworkName,
			ReservedPeeringRanges: tpgresource.ConvertStringArr(d.Get("reserved_peering_ranges").([]interface{})),
		}

		networkFieldValue, err := tpgresource.ParseNetworkFieldValue(network, d, config)
		if err != nil {
			return errwrap.Wrapf("Failed to retrieve network field value, err: {{err}}", err)
		}
		project := networkFieldValue.Project

		// The API docs don't specify that you can do connections/-, but that's what gcloud does,
		// and it's easier than grabbing the connection name.

		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			project = bp
		}

		patchCall := config.NewServiceNetworkingClient(userAgent).Services.Connections.Patch(parentService+"/connections/-", connection).UpdateMask("reservedPeeringRanges").Force(true)
		if config.UserProjectOverride {
			patchCall.Header().Add("X-Goog-User-Project", project)
		}
		op, err := patchCall.Do()
		if err != nil {
			return err
		}
		if err := ServiceNetworkingOperationWaitTime(config, op, "Update Service Networking Connection", userAgent, project, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	return resourceServiceNetworkingConnectionRead(d, meta)
}

func resourceServiceNetworkingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	network := d.Get("network").(string)
	serviceNetworkingNetworkName, err := RetrieveServiceNetworkingNetworkName(d, config, network, userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	peering := d.Get("peering").(string)
	obj["name"] = peering
	url := fmt.Sprintf("%s%s/removePeering", config.ComputeBasePath, serviceNetworkingNetworkName)

	networkFieldValue, err := tpgresource.ParseNetworkFieldValue(network, d, config)
	if err != nil {
		return errwrap.Wrapf("Failed to retrieve network field value, err: {{err}}", err)
	}

	project := networkFieldValue.Project
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ServiceNetworkingConnection %q", d.Id()))
	}

	op := &compute.Operation{}
	err = tpgresource.Convert(res, op)
	if err != nil {
		return err
	}

	err = tpgcompute.ComputeOperationWaitTime(
		config, op, project, "Updating Network", userAgent, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] Service network connection removed.")

	return nil
}

func resourceServiceNetworkingConnectionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	connectionId, err := parseConnectionId(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("network", connectionId.Network); err != nil {
		return nil, fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("service", connectionId.Service); err != nil {
		return nil, fmt.Errorf("Error setting service: %s", err)
	}
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
		return nil, errwrap.Wrapf("Failed to parse service networking connection id, invalid network, err: {{err}}", err)
	} else if len(network) == 0 {
		return nil, fmt.Errorf("Failed to parse service networking connection id, empty network")
	}

	service, err := url.QueryUnescape(res[1])
	if err != nil {
		return nil, errwrap.Wrapf("Failed to parse service networking connection id, invalid service, err: {{err}}", err)
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
func RetrieveServiceNetworkingNetworkName(d *schema.ResourceData, config *transport_tpg.Config, network, userAgent string) (string, error) {
	networkFieldValue, err := tpgresource.ParseNetworkFieldValue(network, d, config)
	if err != nil {
		return "", errwrap.Wrapf("Failed to retrieve network field value, err: {{err}}", err)
	}

	pid := networkFieldValue.Project
	if pid == "" {
		return "", fmt.Errorf("Could not determine project")
	}
	log.Printf("[DEBUG] Retrieving project number by doing a GET with the project id, as required by service networking")
	// err == nil indicates that the billing_project value was found
	billingProject := pid
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	getProjectCall := config.NewResourceManagerClient(userAgent).Projects.Get(pid)
	if config.UserProjectOverride {
		getProjectCall.Header().Add("X-Goog-User-Project", billingProject)
	}
	project, err := getProjectCall.Do()
	if err != nil {
		// note: returning a wrapped error is part of this method's contract!
		// https://blog.golang.org/go1.13-errors
		return "", fmt.Errorf("Failed to retrieve project, pid: %s, err: %w", pid, err)
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
