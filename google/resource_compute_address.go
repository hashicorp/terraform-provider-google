package google

import (
	"fmt"
	"log"

	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

var (
	computeAddressIdTemplate = "projects/%s/regions/%s/addresses/%s"
	computeAddressLinkRegex  = regexp.MustCompile("projects/(.+)/regions/(.+)/addresses/(.+)$")
	AddressBaseApiVersion    = v1
	AddressVersionedFeatures = []Feature{
		{Version: v0beta, Item: "address_type"},
		{Version: v0beta, Item: "subnetwork"},
	}
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
				ValidateFunc: validation.StringInSlice(
					[]string{"INTERNAL", "EXTERNAL"}, false),
			},

			"subnetwork": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: linkDiffSuppress,
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
	computeApiVersion := getComputeApiVersion(d, AddressBaseApiVersion, AddressVersionedFeatures)
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
	addr := &computeBeta.Address{
		Name:        d.Get("name").(string),
		AddressType: d.Get("address_type").(string),
		Subnetwork:  d.Get("subnetwork").(string),
	}

	var op interface{}
	switch computeApiVersion {
	case v1:
		v1Address := &compute.Address{}
		err = Convert(addr, v1Address)
		if err != nil {
			return err
		}
		op, err = config.clientCompute.Addresses.Insert(
			project, region, v1Address).Do()
		if err != nil {
			return fmt.Errorf("Error creating address: %s", err)
		}
	case v0beta:
		v0BetaAddress := &computeBeta.Address{}
		err = Convert(addr, v0BetaAddress)
		if err != nil {
			return err
		}
		op, err = config.clientComputeBeta.Addresses.Insert(
			project, region, v0BetaAddress).Do()
		if err != nil {
			return fmt.Errorf("Error creating address: %s", err)
		}
	}

	// It probably maybe worked, so store the ID now
	d.SetId(computeAddressId{
		Project: project,
		Region:  region,
		Name:    addr.Name,
	}.canonicalId())

	err = computeSharedOperationWait(config, op, project, "Creating Address")
	if err != nil {
		return err
	}

	return resourceComputeAddressRead(d, meta)
}

func resourceComputeAddressRead(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, AddressBaseApiVersion, AddressVersionedFeatures)
	config := meta.(*Config)

	addressId, err := parseComputeAddressId(d.Id(), config)
	if err != nil {
		return err
	}

	addr := &computeBeta.Address{}
	switch computeApiVersion {
	case v1:
		v1Address, err := config.clientCompute.Addresses.Get(
			addressId.Project, addressId.Region, addressId.Name).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Address %q", d.Get("name").(string)))
		}

		err = Convert(v1Address, addr)
		if err != nil {
			return err
		}
	case v0beta:
		v0BetaAddr, err := config.clientComputeBeta.Addresses.Get(
			addressId.Project, addressId.Region, addressId.Name).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Address %q", d.Get("name").(string)))
		}

		err = Convert(v0BetaAddr, addr)
		if err != nil {
			return err
		}
	}

	d.Set("address", addr.Address)
	d.Set("self_link", addr.SelfLink)
	d.Set("name", addr.Name)
	d.Set("region", GetResourceNameFromSelfLink(addr.Region))

	return nil
}

func resourceComputeAddressDelete(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, AddressBaseApiVersion, AddressVersionedFeatures)
	config := meta.(*Config)

	addressId, err := parseComputeAddressId(d.Id(), config)
	if err != nil {
		return err
	}

	var op interface{}
	switch computeApiVersion {
	case v1:
		// Delete the address
		log.Printf("[DEBUG] address delete request")
		op, err = config.clientCompute.Addresses.Delete(
			addressId.Project, addressId.Region, addressId.Name).Do()
		if err != nil {
			return fmt.Errorf("Error deleting address: %s", err)
		}
	case v0beta:
		op, err = config.clientComputeBeta.Addresses.Delete(
			addressId.Project, addressId.Region, addressId.Name).Do()
		if err != nil {
			return fmt.Errorf("Error deleting address: %s", err)
		}
	}

	err = computeSharedOperationWait(config, op, addressId.Project, "Deleting Address")
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

type computeAddressId struct {
	Project string
	Region  string
	Name    string
}

func (s computeAddressId) canonicalId() string {
	return fmt.Sprintf(computeAddressIdTemplate, s.Project, s.Region, s.Name)
}

func parseComputeAddressId(id string, config *Config) (*computeAddressId, error) {
	var parts []string
	if computeAddressLinkRegex.MatchString(id) {
		parts = computeAddressLinkRegex.FindStringSubmatch(id)

		return &computeAddressId{
			Project: parts[1],
			Region:  parts[2],
			Name:    parts[3],
		}, nil
	} else {
		parts = strings.Split(id, "/")
	}

	if len(parts) == 3 {
		return &computeAddressId{
			Project: parts[0],
			Region:  parts[1],
			Name:    parts[2],
		}, nil
	} else if len(parts) == 2 {
		// Project is optional.
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{region}/{name}` id format.")
		}

		return &computeAddressId{
			Project: config.Project,
			Region:  parts[0],
			Name:    parts[1],
		}, nil
	} else if len(parts) == 1 {
		// Project and region is optional
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{name}` id format.")
		}
		if config.Region == "" {
			return nil, fmt.Errorf("The default region for the provider must be set when using the `{name}` id format.")
		}

		return &computeAddressId{
			Project: config.Project,
			Region:  config.Region,
			Name:    parts[0],
		}, nil
	}

	return nil, fmt.Errorf("Invalid compute address id. Expecting resource link, `{project}/{region}/{name}`, `{region}/{name}` or `{name}` format.")
}
