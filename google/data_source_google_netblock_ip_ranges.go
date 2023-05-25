// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type googRanges struct {
	SyncToken    string     `json:"syncToken"`
	CreationTime string     `json:"creationTime"`
	Prefixes     []prefixes `json:"prefixes"`
}

type prefixes struct {
	Ipv4Prefix string `json:"ipv4Prefix"`
	Ipv6Prefix string `json:"ipv6Prefix"`
}

func DataSourceGoogleNetblockIpRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleNetblockIpRangesRead,

		Schema: map[string]*schema.Schema{
			"range_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "cloud-netblocks",
			},
			"cidr_blocks": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cidr_blocks_ipv4": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cidr_blocks_ipv6": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleNetblockIpRangesRead(d *schema.ResourceData, meta interface{}) error {

	rt := d.Get("range_type").(string)
	CidrBlocks := make(map[string][]string)

	switch rt {
	// Dynamic ranges
	case "cloud-netblocks":
		// https://cloud.google.com/compute/docs/faq#find_ip_range
		const CLOUD_NETBLOCK_URL = "https://www.gstatic.com/ipranges/cloud.json"
		CidrBlocks, err := getCidrBlocksFromUrl(CLOUD_NETBLOCK_URL)

		if err != nil {
			return err
		}
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv6", CidrBlocks["cidr_blocks_ipv6"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv6: %s", err)
		}
	case "google-netblocks":
		// https://cloud.google.com/vpc/docs/configure-private-google-access?hl=en#ip-addr-defaults
		const GOOGLE_NETBLOCK_URL = "https://www.gstatic.com/ipranges/goog.json"
		CidrBlocks, err := getCidrBlocksFromUrl(GOOGLE_NETBLOCK_URL)

		if err != nil {
			return err
		}
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv6", CidrBlocks["cidr_blocks_ipv6"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv6: %s", err)
		}
	// Static ranges
	case "restricted-googleapis":
		// https://cloud.google.com/vpc/docs/private-access-options#domain-vips
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "199.36.153.4/30")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
	case "private-googleapis":
		// https://cloud.google.com/vpc/docs/private-access-options#domain-vips
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "199.36.153.8/30")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
	case "dns-forwarders":
		// https://cloud.google.com/dns/zones/#creating-forwarding-zones
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.199.192.0/19")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
	case "iap-forwarders":
		// https://cloud.google.com/iap/docs/using-tcp-forwarding
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.235.240.0/20")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
	case "health-checkers":
		// https://cloud.google.com/load-balancing/docs/health-checks#fw-ruleh
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.191.0.0/16")
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "130.211.0.0/22")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
	case "legacy-health-checkers":
		// https://cloud.google.com/load-balancing/docs/health-check#fw-netlbs
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.191.0.0/16")
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "209.85.152.0/22")
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "209.85.204.0/22")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		if err := d.Set("cidr_blocks", CidrBlocks["cidr_blocks"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks: %s", err)
		}
		if err := d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"]); err != nil {
			return fmt.Errorf("Error setting cidr_blocks_ipv4: %s", err)
		}
	default:
		return fmt.Errorf("Unknown range_type: %s", rt)
	}

	d.SetId("netblock-ip-ranges-" + rt)

	return nil
}

func getCidrBlocksFromUrl(url string) (map[string][]string, error) {
	cidrBlocks := make(map[string][]string)

	response, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Error: %s", err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("Error to retrieve the CIDR list: %s", err)
	}

	ranges := googRanges{}
	jsonErr := json.Unmarshal(body, &ranges)
	if jsonErr != nil {
		return nil, fmt.Errorf("Error reading JSON list: %s", jsonErr)
	}

	for _, element := range ranges.Prefixes {

		if len(element.Ipv4Prefix) > 0 {
			cidrBlocks["cidr_blocks_ipv4"] = append(cidrBlocks["cidr_blocks_ipv4"], element.Ipv4Prefix)
			cidrBlocks["cidr_blocks"] = append(cidrBlocks["cidr_blocks"], element.Ipv4Prefix)
		} else if len(element.Ipv6Prefix) > 0 {
			cidrBlocks["cidr_blocks_ipv6"] = append(cidrBlocks["cidr_blocks_ipv6"], element.Ipv6Prefix)
			cidrBlocks["cidr_blocks"] = append(cidrBlocks["cidr_blocks"], element.Ipv6Prefix)
		}

	}

	return cidrBlocks, nil
}
