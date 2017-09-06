package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
	"strings"
)

func resourceComputeAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeAddressCreate,
		Read:   resourceComputeAddressRead,
		Delete: resourceComputeAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  resourceComputeAddressMigrateState,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address": &schema.Schema{
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
				Computed: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeAddressCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the address parameter
	addr := &compute.Address{Name: d.Get("name").(string)}
	op, err := config.clientCompute.Addresses.Insert(
		project, region, addr).Do()
	if err != nil {
		return fmt.Errorf("Error creating address: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(computeAddressId{
		Project: project,
		Region:  region,
		Name:    addr.Name,
	}.terraformId())

	err = computeOperationWait(config, op, project, "Creating Address")
	if err != nil {
		return err
	}

	return resourceComputeAddressRead(d, meta)
}

func resourceComputeAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	addressId, err := parseComputeAddressId(d.Id())
	if err != nil {
		return err
	}

	addr, err := config.clientCompute.Addresses.Get(
		addressId.Project, addressId.Region, addressId.Name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Address %q", d.Get("name").(string)))
	}

	d.Set("address", addr.Address)
	d.Set("self_link", addr.SelfLink)
	d.Set("name", addr.Name)
	d.Set("region", addr.Region)

	return nil
}

func resourceComputeAddressDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	addressId, err := parseComputeAddressId(d.Id())
	if err != nil {
		return err
	}

	// Delete the address
	log.Printf("[DEBUG] address delete request")
	op, err := config.clientCompute.Addresses.Delete(
		addressId.Project, addressId.Region, addressId.Name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting address: %s", err)
	}

	err = computeOperationWait(config, op, addressId.Project, "Deleting Address")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

type computeAddressId struct {
	Project string
	Region  string
	Name    string
}

func (s computeAddressId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Region, s.Name)
}

func parseComputeAddressId(id string) (*computeAddressId, error) {
	parts := strings.Split(id, "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid compute address id. Expecting {project}/{region}/{name} format.")
	}

	return &computeAddressId{
		Project: parts[0],
		Region:  parts[1],
		Name:    parts[2],
	}, nil
}
