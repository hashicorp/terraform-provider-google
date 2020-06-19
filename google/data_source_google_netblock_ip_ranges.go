package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func dataSourceGoogleNetblockIpRanges() *schema.Resource {
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
		// https://cloud.google.com/compute/docs/faq#where_can_i_find_product_name_short_ip_ranges
		const CLOUD_NETBLOCK_DNS = "_cloud-netblocks.googleusercontent.com"
		CidrBlocks, err := getCidrBlocksFromDns(CLOUD_NETBLOCK_DNS)

		if err != nil {
			return err
		}
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
		d.Set("cidr_blocks_ipv6", CidrBlocks["cidr_blocks_ipv6"])
	case "google-netblocks":
		// https://cloud.google.com/vpc/docs/configure-private-google-access?hl=en#ip-addr-defaults
		const GOOGLE_NETBLOCK_URL = "http://www.gstatic.com/ipranges/goog.json"
		CidrBlocks, err := getCidrBlocksFromUrl(GOOGLE_NETBLOCK_URL)

		if err != nil {
			return err
		}
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
		d.Set("cidr_blocks_ipv6", CidrBlocks["cidr_blocks_ipv6"])
	// Static ranges
	case "restricted-googleapis":
		// https://cloud.google.com/vpc/docs/private-access-options#domain-vips
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "199.36.153.4/30")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
	case "private-googleapis":
		// https://cloud.google.com/vpc/docs/private-access-options#domain-vips
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "199.36.153.8/30")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
	case "dns-forwarders":
		// https://cloud.google.com/dns/zones/#creating-forwarding-zones
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.199.192.0/19")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
	case "iap-forwarders":
		// https://cloud.google.com/iap/docs/using-tcp-forwarding
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.235.240.0/20")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
	case "health-checkers":
		// https://cloud.google.com/load-balancing/docs/health-checks#fw-ruleh
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.191.0.0/16")
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "130.211.0.0/22")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
	case "legacy-health-checkers":
		// https://cloud.google.com/load-balancing/docs/health-check#fw-netlbs
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "35.191.0.0/16")
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "209.85.152.0/22")
		CidrBlocks["cidr_blocks_ipv4"] = append(CidrBlocks["cidr_blocks_ipv4"], "209.85.204.0/22")
		CidrBlocks["cidr_blocks"] = CidrBlocks["cidr_blocks_ipv4"]
		d.Set("cidr_blocks", CidrBlocks["cidr_blocks"])
		d.Set("cidr_blocks_ipv4", CidrBlocks["cidr_blocks_ipv4"])
	default:
		return fmt.Errorf("Unknown range_type: %s", rt)
	}

	d.SetId("netblock-ip-ranges-" + rt)

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

func getCidrBlocksFromDns(netblock string) (map[string][]string, error) {
	var dnsNetblockList []string
	cidrBlocks := make(map[string][]string)

	response, err := netblock_request(netblock)

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
