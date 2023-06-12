// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeAddresses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGoogleComputeAddressesRead,

		Schema: map[string]*schema.Schema{
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the IP address.`,
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address.`,
						},
						"address_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address type.`,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"filter": {
				Type: schema.TypeString,
				Description: `Filter sets the optional parameter "filter": A filter expression that
filters resources listed in the response. The expression must specify
the field name, an operator, and the value that you want to use for
filtering. The value must be a string, a number, or a boolean. The
operator must be either "=", "!=", ">", "<", "<=", ">=" or ":". For
example, if you are filtering Compute Engine instances, you can
exclude instances named "example-instance" by specifying "name !=
example-instance". The ":" operator can be used with string fields to
match substrings. For non-string fields it is equivalent to the "="
operator. The ":*" comparison can be used to test whether a key has
been defined. For example, to find all objects with "owner" label
use: """ labels.owner:* """ You can also filter nested fields. For
example, you could specify "scheduling.automaticRestart = false" to
include instances only if they are not scheduled for automatic
restarts. You can use filtering on nested fields to filter based on
resource labels. To filter on multiple expressions, provide each
separate expression within parentheses. For example: """
(scheduling.automaticRestart = true) (cpuPlatform = "Intel Skylake")
""" By default, each expression is an "AND" expression. However, you
can include "AND" and "OR" expressions explicitly. For example: """
(cpuPlatform = "Intel Skylake") OR (cpuPlatform = "Intel Broadwell")
AND (scheduling.automaticRestart = true) """`,
				Optional: true,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Region that should be considered to search addresses. All regions are considered if missing.`,
			},

			"project": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `The google project in which addresses are listed. Defaults to provider's configuration if missing.`,
			},
		},
	}
}

func dataSourceGoogleComputeAddressesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return diag.FromErr(err)
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	allAddresses := make([]map[string]interface{}, 0)

	client := config.NewComputeClient(userAgent).Addresses
	if region, has_region := d.GetOk("region"); has_region {
		request := client.List(project, region.(string))
		if filter, has_filter := d.GetOk("filter"); has_filter {
			request = request.Filter(filter.(string))
		}
		err = request.Pages(context, func(addresses *compute.AddressList) error {
			for _, address := range addresses.Items {
				allAddresses = append(allAddresses, generateTfAddress(address))
			}
			return nil
		})
	} else {
		request := client.AggregatedList(project)
		if filter, has_filter := d.GetOk("filter"); has_filter {
			request = request.Filter(filter.(string))
		}
		err = request.Pages(context, func(addresses *compute.AddressAggregatedList) error {
			for _, items := range addresses.Items {
				for _, address := range items.Addresses {
					allAddresses = append(allAddresses, generateTfAddress(address))
				}
			}
			return nil
		})
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("addresses", allAddresses); err != nil {
		return diag.FromErr(fmt.Errorf("error setting addresses: %s", err))
	}

	if err := d.Set("project", project); err != nil {
		return diag.FromErr(fmt.Errorf("error setting project: %s", err))
	}
	d.SetId(computeId(project, d))
	return nil
}

func generateTfAddress(address *compute.Address) map[string]interface{} {
	return map[string]interface{}{
		"name":         address.Name,
		"address":      address.Address,
		"address_type": address.AddressType,
		"description":  address.Description,
		"region":       regionFromUrl(address.Region),
		"status":       address.Status,
		"self_link":    address.SelfLink,
	}
}

func computeId(project string, d *schema.ResourceData) string {
	region := "ALL"
	filter := "ALL"
	if p_region, has_region := d.GetOk("region"); has_region {
		region = p_region.(string)
	}
	if p_filter, has_filter := d.GetOk("filter"); has_filter {
		filter = p_filter.(string)
	}
	return fmt.Sprintf("%s-%s-%s", project, region, filter)
}

func regionFromUrl(url string) string {
	parts := strings.Split(url, "/")
	if count := len(parts); count > 0 {
		return parts[count-1]
	}
	return ""
}
