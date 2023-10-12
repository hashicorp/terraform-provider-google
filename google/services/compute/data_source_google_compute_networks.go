// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeNetworksRead,

		Schema: map[string]*schema.Schema{

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleComputeNetworksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	networkList, err := config.NewComputeClient(userAgent).Networks.List(project).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Network Not Found : %s", project))
	}

	var networks = make([]string, len(networkList.Items))

	for i := 0; i < len(networkList.Items); i++ {
		networks[i] = networkList.Items[i].Name
	}

	if err := d.Set("networks", networks); err != nil {
		return fmt.Errorf("Error setting the network names: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting the network names: %s", err)
	}

	if err := d.Set("self_link", networkList.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/global/networks", project))
	return nil
}
