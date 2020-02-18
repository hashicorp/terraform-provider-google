package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/dns/v1"
)

// DNSSEC Algorithm Numbers: https://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml
// The following are algorithms that are supported by Cloud DNS
var dnssecAlgoNums = map[string]int{
	"rsasha1":         5,
	"rsasha256":       8,
	"rsasha512":       10,
	"ecdsap256sha256": 13,
	"ecdsap384sha384": 14,
}

// DS RR Digest Types: https://www.iana.org/assignments/ds-rr-types/ds-rr-types.xhtml
// The following are digests that are supported by Cloud DNS
var dnssecDigestType = map[string]int{
	"sha1":   1,
	"sha256": 2,
	"sha384": 4,
}

func dataSourceDNSKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSKeysRead,

		Schema: map[string]*schema.Schema{
			"managed_zone": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_signing_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     kskResource(),
			},
			"zone_signing_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dnsKeyResource(),
			},
		},
	}
}

func dnsKeyResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"digests": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"digest": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"key_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key_tag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func kskResource() *schema.Resource {
	resource := dnsKeyResource()

	resource.Schema["ds_record"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return resource
}

func generateDSRecord(signingKey *dns.DnsKey) (string, error) {
	algoNum, found := dnssecAlgoNums[signingKey.Algorithm]
	if !found {
		return "", fmt.Errorf("DNSSEC Algorithm number for %s not found", signingKey.Algorithm)
	}

	digestType, found := dnssecDigestType[signingKey.Digests[0].Type]
	if !found {
		return "", fmt.Errorf("DNSSEC Digest type for %s not found", signingKey.Digests[0].Type)
	}

	return fmt.Sprintf("%d %d %d %s",
		signingKey.KeyTag,
		algoNum,
		digestType,
		signingKey.Digests[0].Digest), nil
}

func flattenSigningKeys(signingKeys []*dns.DnsKey, keyType string) []map[string]interface{} {
	var keys []map[string]interface{}

	for _, signingKey := range signingKeys {
		if signingKey != nil && signingKey.Type == keyType {
			data := map[string]interface{}{
				"algorithm":     signingKey.Algorithm,
				"creation_time": signingKey.CreationTime,
				"description":   signingKey.Description,
				"digests":       flattenDigests(signingKey.Digests),
				"id":            signingKey.Id,
				"is_active":     signingKey.IsActive,
				"key_length":    signingKey.KeyLength,
				"key_tag":       signingKey.KeyTag,
				"public_key":    signingKey.PublicKey,
			}

			if signingKey.Type == "keySigning" && len(signingKey.Digests) > 0 {
				dsRecord, err := generateDSRecord(signingKey)
				if err == nil {
					data["ds_record"] = dsRecord
				}
			}

			keys = append(keys, data)
		}
	}

	return keys
}

func flattenDigests(dnsKeyDigests []*dns.DnsKeyDigest) []map[string]interface{} {
	var digests []map[string]interface{}

	for _, dnsKeyDigest := range dnsKeyDigests {
		if dnsKeyDigest != nil {
			data := map[string]interface{}{
				"digest": dnsKeyDigest.Digest,
				"type":   dnsKeyDigest.Type,
			}

			digests = append(digests, data)
		}
	}

	return digests
}

func dataSourceDNSKeysRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	fv, err := parseProjectFieldValue("managedZones", d.Get("managed_zone").(string), "project", d, config, false)
	if err != nil {
		return err
	}
	project := fv.Project
	managedZone := fv.Name

	d.Set("project", project)
	d.SetId(fmt.Sprintf("projects/%s/managedZones/%s", project, managedZone))

	log.Printf("[DEBUG] Fetching DNS keys from managed zone %s", managedZone)

	response, err := config.clientDns.DnsKeys.List(project, managedZone).Do()
	if err != nil && !isGoogleApiErrorWithCode(err, 404) {
		return fmt.Errorf("error retrieving DNS keys: %s", err)
	} else if isGoogleApiErrorWithCode(err, 404) {
		return nil
	}

	log.Printf("[DEBUG] Fetched DNS keys from managed zone %s", managedZone)

	d.Set("key_signing_keys", flattenSigningKeys(response.DnsKeys, "keySigning"))
	d.Set("zone_signing_keys", flattenSigningKeys(response.DnsKeys, "zoneSigning"))

	return nil
}
