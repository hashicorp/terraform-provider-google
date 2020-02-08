package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleKmsCryptoKeyVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsCryptoKeyVersionRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pem": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleKmsCryptoKeyVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting attributes for CryptoKeyVersion: %#v", url)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Get("crypto_key").(string), config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", cryptoKeyId.KeyRingId.Project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("KmsCryptoKeyVersion %q", d.Id()))
	}

	if err := d.Set("version", flattenKmsCryptoKeyVersionVersion(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading CryptoKeyVersion: %s", err)
	}
	if err := d.Set("state", flattenKmsCryptoKeyVersionState(res["state"], d)); err != nil {
		return fmt.Errorf("Error reading CryptoKeyVersion: %s", err)
	}
	if err := d.Set("protection_level", flattenKmsCryptoKeyVersionProtectionLevel(res["protectionLevel"], d)); err != nil {
		return fmt.Errorf("Error reading CryptoKeyVersion: %s", err)
	}
	if err := d.Set("algorithm", flattenKmsCryptoKeyVersionAlgorithm(res["algorithm"], d)); err != nil {
		return fmt.Errorf("Error reading CryptoKeyVersion: %s", err)
	}

	url, err = replaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting purpose of CryptoKey: %#v", url)
	res, err = sendRequest(config, "GET", cryptoKeyId.KeyRingId.Project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("KmsCryptoKey %q", d.Id()))
	}

	if res["purpose"] == "ASYMMETRIC_SIGN" || res["purpose"] == "ASYMMETRIC_DECRYPT" {
		url, err = replaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}/publicKey")
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Getting public key of CryptoKeyVersion: %#v", url)
		res, _ = sendRequest(config, "GET", cryptoKeyId.KeyRingId.Project, url, nil)

		if err := d.Set("public_key", flattenKmsCryptoKeyVersionPublicKey(res, d)); err != nil {
			return fmt.Errorf("Error reading CryptoKeyVersion public key: %s", err)
		}
	}
	d.SetId(fmt.Sprintf("//cloudkms.googleapis.com/%s/cryptoKeyVersions/%d", d.Get("crypto_key"), d.Get("version")))

	return nil
}

func flattenKmsCryptoKeyVersionVersion(v interface{}, d *schema.ResourceData) interface{} {
	parts := strings.Split(v.(string), "/")
	version := parts[len(parts)-1]
	// Handles the string fixed64 format
	if intVal, err := strconv.ParseInt(version, 10, 64); err == nil {
		return intVal
	} // let terraform core handle it if we can't convert the string to an int.
	return v
}

func flattenKmsCryptoKeyVersionState(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionProtectionLevel(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionAlgorithm(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionPublicKey(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["pem"] =
		flattenKmsCryptoKeyVersionPublicKeyPem(original["pem"], d)
	transformed["algorithm"] =
		flattenKmsCryptoKeyVersionPublicKeyAlgorithm(original["algorithm"], d)
	return []interface{}{transformed}
}
func flattenKmsCryptoKeyVersionPublicKeyPem(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionPublicKeyAlgorithm(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
