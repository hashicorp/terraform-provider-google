// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeGlobalAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeGlobalAddressRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"address_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_tier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"prefix_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"purpose": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnetwork": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"users": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleComputeGlobalAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	address, err := config.NewComputeClient(userAgent).GlobalAddresses.Get(project, name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Global Address Not Found : %s", name))
	}

	if err := d.Set("address", address.Address); err != nil {
		return fmt.Errorf("Error setting address: %s", err)
	}
	if err := d.Set("address_type", address.AddressType); err != nil {
		return fmt.Errorf("Error setting address_type: %s", err)
	}
	if err := d.Set("network", address.Network); err != nil {
		return fmt.Errorf("Error setting network: %s", err)
	}
	if err := d.Set("network_tier", address.NetworkTier); err != nil {
		return fmt.Errorf("Error setting network_tier: %s", err)
	}
	if err := d.Set("prefix_length", address.PrefixLength); err != nil {
		return fmt.Errorf("Error setting prefix_length: %s", err)
	}
	if err := d.Set("purpose", address.Purpose); err != nil {
		return fmt.Errorf("Error setting purpose: %s", err)
	}
	if err := d.Set("subnetwork", address.Subnetwork); err != nil {
		return fmt.Errorf("Error setting subnetwork: %s", err)
	}
	if err := d.Set("status", address.Status); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}
	if err := d.Set("self_link", address.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/global/addresses/%s", project, name))
	return nil
}
