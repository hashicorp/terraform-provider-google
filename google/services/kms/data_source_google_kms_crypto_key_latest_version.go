// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleKmsLatestCryptoKeyVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsLatestCryptoKeyVersionRead,
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
				Computed: true,
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
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `
					The filter argument is used to add a filter query parameter that limits which type of cryptoKeyVersion is retrieved as the latest by the data source: ?filter={{filter}}. When no value is provided there is no filtering.

					Example filter values if filtering on state.

					* "state:ENABLED" will retrieve the latest cryptoKeyVersion that has the state "ENABLED".

					[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)
				`,
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

func dataSourceGoogleKmsLatestCryptoKeyVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	cryptoKeyId, err := ParseKmsCryptoKeyId(d.Get("crypto_key").(string), config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s/latestCryptoKeyVersion", cryptoKeyId.CryptoKeyId())
	if filter, ok := d.GetOk("filter"); ok {
		id += "/filter=" + filter.(string)
	}
	d.SetId(id)

	versions, err := dataSourceKMSCryptoKeyVersionsList(d, meta, cryptoKeyId.CryptoKeyId(), userAgent)
	if err != nil {
		return err
	}

	// grab latest version
	lv := len(versions) - 1
	if lv < 0 {
		return fmt.Errorf("No CryptoVersions found in crypto key %s", cryptoKeyId.CryptoKeyId())
	}

	latestVersion := versions[lv].(map[string]interface{})

	// The google_kms_crypto_key resource and dataset set
	// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}})
	// and set name is set as just {{name}}.

	if err := d.Set("name", flattenKmsCryptoKeyVersionName(latestVersion["name"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}
	if err := d.Set("version", flattenKmsCryptoKeyVersionVersion(latestVersion["name"], d)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("state", flattenKmsCryptoKeyVersionState(latestVersion["state"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}
	if err := d.Set("protection_level", flattenKmsCryptoKeyVersionProtectionLevel(latestVersion["protectionLevel"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}
	if err := d.Set("algorithm", flattenKmsCryptoKeyVersionAlgorithm(latestVersion["algorithm"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting attributes for CryptoKeyVersion: %#v", url)

	url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting purpose of CryptoKey: %#v", url)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   cryptoKeyId.KeyRingId.Project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("KmsCryptoKey %q", d.Id()), url)
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

	return nil
}
