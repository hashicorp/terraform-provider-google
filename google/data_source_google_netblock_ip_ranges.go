package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"net/http"
	"strings"
)

func dataSourceGoogleNetblockIpRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleNetblockIpRangesRead,

		Schema: map[string]*schema.Schema{
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
	d.SetId("netblock-ip-ranges")

	// https://cloud.google.com/compute/docs/faq#where_can_i_find_product_name_short_ip_ranges
	CidrBlocks, err := getCidrBlocks()

	if err != nil {
		return err
	}

	d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
	d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
	d.Set("cidr_blocks_ipv6", CidrBlocks["cidr_blocks_ipv6"])

	return nil
}

func netblock_request(name string) (string, error) {
	response, err := http.Get(fmt.Sprintf("https://dns.google.com/resolve?name=%s&type=TXT", name))

	if err != nil {
		return "", fmt.Errorf("Error from _cloud-netblocks: %s", err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("Error to retrieve the domains list: %s", err)
	}

	return string(body), nil
}

func getCidrBlocks() (map[string][]string, error) {
	const INITIAL_NETBLOCK_DNS = "_cloud-netblocks.googleusercontent.com"
	var dnsNetblockList []string
	cidrBlocks := make(map[string][]string)

	response, err := netblock_request(INITIAL_NETBLOCK_DNS)

	if err != nil {
		return nil, err
	}

	splitedResponse := strings.Split(response, " ")

	for _, sp := range splitedResponse {
		if strings.HasPrefix(sp, "include:") {
			dnsNetblock := strings.Replace(sp, "include:", "", 1)
			dnsNetblockList = append(dnsNetblockList, dnsNetblock)
		}
	}

	for len(dnsNetblockList) > 0 {

		dnsNetblock := dnsNetblockList[0]

		dnsNetblockList[0] = ""
		dnsNetblockList = dnsNetblockList[1:]

		response, err = netblock_request(dnsNetblock)

		if err != nil {
			return nil, err
		}

		splitedResponse = strings.Split(response, " ")

		for _, sp := range splitedResponse {
			if strings.HasPrefix(sp, "ip4") {
				cdrBlock := strings.Replace(sp, "ip4:", "", 1)
				cidrBlocks["cidr_blocks_ipv4"] = append(cidrBlocks["cidr_blocks_ipv4"], cdrBlock)
				cidrBlocks["cidr_blocks"] = append(cidrBlocks["cidr_blocks"], cdrBlock)

			} else if strings.HasPrefix(sp, "ip6") {
				cdrBlock := strings.Replace(sp, "ip6:", "", 1)
				cidrBlocks["cidr_blocks_ipv6"] = append(cidrBlocks["cidr_blocks_ipv6"], cdrBlock)
				cidrBlocks["cidr_blocks"] = append(cidrBlocks["cidr_blocks"], cdrBlock)

			} else if strings.HasPrefix(sp, "include:") {
				cidr_block := strings.Replace(sp, "include:", "", 1)
				dnsNetblockList = append(dnsNetblockList, cidr_block)
			}
		}
	}

	return cidrBlocks, nil
}
