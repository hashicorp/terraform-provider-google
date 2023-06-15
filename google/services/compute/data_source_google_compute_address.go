// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	computeAddressIdTemplate = "projects/%s/regions/%s/addresses/%s"
	computeAddressLinkRegex  = regexp.MustCompile("projects/(.+)/regions/(.+)/addresses/(.+)$")
)

func DataSourceGoogleComputeAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeAddressRead,

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

			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleComputeAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	address, err := config.NewComputeClient(userAgent).Addresses.Get(project, region, name).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Address Not Found : %s", name))
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
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/addresses/%s", project, region, name))
	return nil
}

type ComputeAddressId struct {
	Project string
	Region  string
	Name    string
}

func (s ComputeAddressId) CanonicalId() string {
	return fmt.Sprintf(computeAddressIdTemplate, s.Project, s.Region, s.Name)
}

func ParseComputeAddressId(id string, config *transport_tpg.Config) (*ComputeAddressId, error) {
	var parts []string
	if computeAddressLinkRegex.MatchString(id) {
		parts = computeAddressLinkRegex.FindStringSubmatch(id)

		return &ComputeAddressId{
			Project: parts[1],
			Region:  parts[2],
			Name:    parts[3],
		}, nil
	} else {
		parts = strings.Split(id, "/")
	}

	if len(parts) == 3 {
		return &ComputeAddressId{
			Project: parts[0],
			Region:  parts[1],
			Name:    parts[2],
		}, nil
	} else if len(parts) == 2 {
		// Project is optional.
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{region}/{name}` id format.")
		}

		return &ComputeAddressId{
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

		return &ComputeAddressId{
			Project: config.Project,
			Region:  config.Region,
			Name:    parts[0],
		}, nil
	}

	return nil, fmt.Errorf("Invalid compute address id. Expecting resource link, `{project}/{region}/{name}`, `{region}/{name}` or `{name}` format.")
}
