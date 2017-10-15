package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"google.golang.org/api/compute/v1"
)

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
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the address parameter
	addr := &compute.Address{
		Name:      d.Get("name").(string),
		IpVersion: d.Get("ip_version").(string),
	}

	op, err := config.clientCompute.GlobalAddresses.Insert(project, addr).Do()
	if err != nil {
		return fmt.Errorf("Error creating address: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(addr.Name)

	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating Global Address")
	if err != nil {
		return err
	}

	return resourceComputeGlobalAddressRead(d, meta)
}

func resourceComputeGlobalAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	addr, err := config.clientCompute.GlobalAddresses.Get(project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Global Address %q", d.Get("name").(string)))
	}

	d.Set("name", addr.Name)
	d.Set("ip_version", addr.IpVersion)
	d.Set("address", addr.Address)
	d.Set("self_link", ConvertSelfLinkToV1(addr.SelfLink))

	return nil
}

func resourceComputeGlobalAddressDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the address
	log.Printf("[DEBUG] address delete request")
	op, err := config.clientCompute.GlobalAddresses.Delete(project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting address: %s", err)
	}

	err = computeSharedOperationWait(config.clientCompute, op, project, "Deleting Global Address")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
