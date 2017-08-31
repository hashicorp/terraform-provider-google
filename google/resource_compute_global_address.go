package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

var GlobalAddressBaseApiVersion = v1
var GlobalAddressVersionedFeatures = []Feature{Feature{Version: v0beta, Item: "ip_version"}}

func resourceComputeGlobalAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeGlobalAddressCreate,
		Read:   resourceComputeGlobalAddressRead,
		Delete: resourceComputeGlobalAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6"}, false),
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeGlobalAddressCreate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, GlobalAddressBaseApiVersion, GlobalAddressVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the address parameter
	addr := &computeBeta.Address{
		Name:      d.Get("name").(string),
		IpVersion: d.Get("ip_version").(string),
	}

	var op interface{}
	switch computeApiVersion {
	case v1:
		v1Addr := &compute.Address{}
		err = Convert(addr, v1Addr)
		if err != nil {
			return err
		}

		op, err = config.clientCompute.GlobalAddresses.Insert(project, v1Addr).Do()
		if err != nil {
			return fmt.Errorf("Error creating address: %s", err)
		}
	case v0beta:
		v0BetaAddr := &computeBeta.Address{}
		err = Convert(addr, v0BetaAddr)
		if err != nil {
			return err
		}

		op, err = config.clientComputeBeta.GlobalAddresses.Insert(project, v0BetaAddr).Do()
		if err != nil {
			return fmt.Errorf("Error creating address: %s", err)
		}
	}

	// It probably maybe worked, so store the ID now
	d.SetId(addr.Name)

	err = computeSharedOperationWait(config, op, project, "Creating Global Address")
	if err != nil {
		return err
	}

	return resourceComputeGlobalAddressRead(d, meta)
}

func resourceComputeGlobalAddressRead(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, GlobalAddressBaseApiVersion, GlobalAddressVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	addr := &computeBeta.Address{}
	switch computeApiVersion {
	case v1:
		v1Addr, err := config.clientCompute.GlobalAddresses.Get(project, d.Id()).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Global Address %q", d.Get("name").(string)))
		}

		err = Convert(v1Addr, addr)
		if err != nil {
			return err
		}
	case v0beta:
		v0BetaAddr, err := config.clientComputeBeta.GlobalAddresses.Get(project, d.Id()).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Global Address %q", d.Get("name").(string)))
		}

		err = Convert(v0BetaAddr, addr)
		if err != nil {
			return err
		}
	}

	d.Set("name", addr.Name)
	d.Set("ip_version", addr.IpVersion)
	d.Set("address", addr.Address)
	d.Set("self_link", ConvertSelfLinkToV1(addr.SelfLink))

	return nil
}

func resourceComputeGlobalAddressDelete(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, GlobalAddressBaseApiVersion, GlobalAddressVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the address
	log.Printf("[DEBUG] address delete request")
	var op interface{}
	switch computeApiVersion {
	case v1:
		op, err = config.clientCompute.GlobalAddresses.Delete(project, d.Id()).Do()
		if err != nil {
			return fmt.Errorf("Error deleting address: %s", err)
		}
	case v0beta:
		op, err = config.clientComputeBeta.GlobalAddresses.Delete(project, d.Id()).Do()
		if err != nil {
			return fmt.Errorf("Error deleting address: %s", err)
		}
	}

	err = computeSharedOperationWait(config, op, project, "Deleting Global Address")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
