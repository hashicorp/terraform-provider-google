package google

import (
	"fmt"

	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/compute/v1"
)

const (
	addressTypeExternal = "EXTERNAL"
	addressTypeInternal = "INTERNAL"
)

var (
	computeAddressIdTemplate = "projects/%s/regions/%s/addresses/%s"
	computeAddressLinkRegex  = regexp.MustCompile("projects/(.+)/regions/(.+)/addresses/(.+)$")
)

func resourceComputeAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeAddressCreate,
		Read:   resourceComputeAddressRead,
		Delete: resourceComputeAddressDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeAddressImportState,
		},

		SchemaVersion: 1,
		MigrateState:  resourceComputeAddressMigrateState,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  addressTypeExternal,
				ValidateFunc: validation.StringInSlice(
					[]string{addressTypeInternal, addressTypeExternal}, false),
			},

			"subnetwork": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: linkDiffSuppress,
			},

			// address will be computed unless it is specified explicitly.
			// address may only be specified for the INTERNAL address_type.
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	address := &compute.Address{
		Name:        d.Get("name").(string),
		AddressType: d.Get("address_type").(string),
		Subnetwork:  d.Get("subnetwork").(string),
		Address:     d.Get("address").(string),
	}

	op, err := config.clientCompute.Addresses.Insert(project, region, address).Do()
	if err != nil {
		return fmt.Errorf("Error creating address: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(computeAddressId{
		Project: project,
		Region:  region,
		Name:    address.Name,
	}.canonicalId())

	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating Address")
	if err != nil {
		return err
	}

	return resourceComputeAddressRead(d, meta)
}

func resourceComputeAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	addressId, err := parseComputeAddressId(d.Id(), config)
	if err != nil {
		return err
	}

	addr, err := config.clientCompute.Addresses.Get(
		addressId.Project, addressId.Region, addressId.Name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Address %q", d.Get("name").(string)))
	}

	d.Set("address_type", addr.AddressType)
	// The API returns an empty AddressType for EXTERNAL address.
	if addr.AddressType == "" {
		d.Set("address_type", addressTypeExternal)
	}
	d.Set("subnetwork", addr.Subnetwork)
	d.Set("address", addr.Address)
	d.Set("self_link", addr.SelfLink)
	d.Set("name", addr.Name)
	d.Set("project", addressId.Project)
	d.Set("region", GetResourceNameFromSelfLink(addr.Region))

	return nil
}

func resourceComputeAddressDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	addressId, err := parseComputeAddressId(d.Id(), config)
	if err != nil {
		return err
	}

	// Delete the address
	op, err := config.clientCompute.Addresses.Delete(
		addressId.Project, addressId.Region, addressId.Name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting address: %s", err)
	}

	err = computeSharedOperationWait(config.clientCompute, op, addressId.Project, "Deleting Address")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceComputeAddressImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	addressId, err := parseComputeAddressId(d.Id(), config)
	if err != nil {
		return nil, err
	}

	d.SetId(addressId.canonicalId())

	return []*schema.ResourceData{d}, nil
}
