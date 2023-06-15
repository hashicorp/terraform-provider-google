// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleKmsCryptoKeyVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsCryptoKeyVersionRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting attributes for CryptoKeyVersion: %#v", url)

	cryptoKeyId, err := ParseKmsCryptoKeyId(d.Get("crypto_key").(string), config)
	if err != nil {
		return err
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   cryptoKeyId.KeyRingId.Project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("KmsCryptoKeyVersion %q", d.Id()))
	}

	if err := d.Set("version", flattenKmsCryptoKeyVersionVersion(res["name"], d)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("name", flattenKmsCryptoKeyVersionName(res["name"], d)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("state", flattenKmsCryptoKeyVersionState(res["state"], d)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("protection_level", flattenKmsCryptoKeyVersionProtectionLevel(res["protectionLevel"], d)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("algorithm", flattenKmsCryptoKeyVersionAlgorithm(res["algorithm"], d)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}

	url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting purpose of CryptoKey: %#v", url)
	res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   cryptoKeyId.KeyRingId.Project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("KmsCryptoKey %q", d.Id()))
	}

	if res["purpose"] == "ASYMMETRIC_SIGN" || res["purpose"] == "ASYMMETRIC_DECRYPT" {
		url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}/publicKey")
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Getting public key of CryptoKeyVersion: %#v", url)

		res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              cryptoKeyId.KeyRingId.Project,
			RawURL:               url,
			UserAgent:            userAgent,
			Timeout:              d.Timeout(schema.TimeoutRead),
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsCryptoKeyVersionsPendingGeneration},
		})

		if err != nil {
			log.Printf("Error generating public key: %s", err)
			return err
		}

		if err := d.Set("public_key", flattenKmsCryptoKeyVersionPublicKey(res, d)); err != nil {
			return fmt.Errorf("Error setting CryptoKeyVersion public key: %s", err)
		}
	}
	d.SetId(fmt.Sprintf("//cloudkms.googleapis.com/v1/%s/cryptoKeyVersions/%d", d.Get("crypto_key"), d.Get("version")))

	return nil
}

func flattenKmsCryptoKeyVersionVersion(v interface{}, d *schema.ResourceData) interface{} {
	parts := strings.Split(v.(string), "/")
	version := parts[len(parts)-1]
	// Handles the string fixed64 format
	if intVal, err := tpgresource.StringToFixed64(version); err == nil {
		return intVal
	} // let terraform core handle it if we can't convert the string to an int.
	return v
}

func flattenKmsCryptoKeyVersionName(v interface{}, d *schema.ResourceData) interface{} {
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
