// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeLbIpRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeLbIpRangesRead,

		Schema: map[string]*schema.Schema{
			"network": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"http_ssl_tcp_internal": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeLbIpRangesRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("compute-lb-ip-ranges")

	// https://cloud.google.com/compute/docs/load-balancing/health-checks#health_check_source_ips_and_firewall_rules

	networkIpRanges := []string{
		"209.85.152.0/22",
		"209.85.204.0/22",
		"35.191.0.0/16",
	}
	if err := d.Set("network", networkIpRanges); err != nil {
		return fmt.Errorf("Error setting network: %s", err)
	}

	httpSslTcpInternalRanges := []string{
		"130.211.0.0/22",
		"35.191.0.0/16",
	}
	if err := d.Set("http_ssl_tcp_internal", httpSslTcpInternalRanges); err != nil {
		return fmt.Errorf("Error setting http_ssl_tcp_internal: %s", err)
	}

	return nil
}
