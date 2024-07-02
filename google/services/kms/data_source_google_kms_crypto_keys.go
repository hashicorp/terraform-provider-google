// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleKmsCryptoKeys() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceKMSCryptoKey().Schema)
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "key_ring")

	// We need to explicitly add the id field to the schema used for individual keys
	// Currently the id field in the google_kms_crypto_key resource is implied/added by the SDK
	dsSchema["id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Read: dataSourceGoogleKmsCryptoKeysRead,
		Schema: map[string]*schema.Schema{
			"key_ring": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The key ring that the keys belongs to. Format: 'projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}'.`,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `
					The filter argument is used to add a filter query parameter that limits which keys are retrieved by the data source: ?filter={{filter}}.
					Example values:
					
					* "name:my-key-" will retrieve keys that contain "my-key-" anywhere in their name. Note: names take the form projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{cryptoKey}}.
					* "name=projects/my-project/locations/global/keyRings/my-key-ring/cryptoKeys/my-key-1" will only retrieve a key with that exact name.
					
					[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)
				`,
			},
			"keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all the retrieved keys from the provided key ring",
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
		},
	}
}

func dataSourceGoogleKmsCryptoKeysRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	keyRingId, err := parseKmsKeyRingId(d.Get("key_ring").(string), config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s/cryptoKeys", keyRingId.KeyRingId())
	if filter, ok := d.GetOk("filter"); ok {
		id += "/filter=" + filter.(string)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Searching for keys in key ring %s", keyRingId.KeyRingId())
	keys, err := dataSourceKMSCryptoKeysList(d, meta, keyRingId.KeyRingId())
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		log.Printf("[DEBUG] Found %d keys in key ring %s", len(keys), keyRingId.KeyRingId())
		value, err := flattenKMSKeysList(d, config, keys, keyRingId.KeyRingId())
		if err != nil {
			return fmt.Errorf("error flattening keys list: %s", err)
		}
		if err := d.Set("keys", value); err != nil {
			return fmt.Errorf("error setting keys: %s", err)
		}
	} else {
		log.Printf("[DEBUG] Found 0 keys in key ring %s", keyRingId.KeyRingId())
	}

	return nil
}

// dataSourceKMSCryptoKeysList calls the list endpoint for Crypto Key resources and collects all keys in a slice.
// This function handles pagination by collecting the resources returned by multiple calls to the list endpoint.
// This function also handles server-side filtering by setting the filter query parameter on each API call.
func dataSourceKMSCryptoKeysList(d *schema.ResourceData, meta interface{}, keyRingId string) ([]interface{}, error) {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{key_ring}}/cryptoKeys")
	if err != nil {
		return nil, err
	}

	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// Always include the filter param, and optionally include the pageToken parameter for subsequent requests
	var params = make(map[string]string, 0)
	if filter, ok := d.GetOk("filter"); ok {
		log.Printf("[DEBUG] Search for keys in key ring %s is using filter ?filter=%s", keyRingId, filter.(string))
		params["filter"] = filter.(string)
	}

	cryptoKeys := make([]interface{}, 0)
	for {
		// Depending on previous iterations, params might contain a pageToken param
		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return nil, err
		}

		headers := make(http.Header)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Headers:   headers,
			// ErrorRetryPredicates used to allow retrying if rate limits are hit when requesting multiple pages in a row
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
		})
		if err != nil {
			return nil, transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("KMSCryptoKeys %q", d.Id()))
		}

		if res == nil {
			// Decoding the object has resulted in it being gone. It may be marked deleted
			log.Printf("[DEBUG] Removing KMSCryptoKey because it no longer exists.")
			d.SetId("")
			return nil, nil
		}

		// Store info from this page
		if v, ok := res["cryptoKeys"].([]interface{}); ok {
			cryptoKeys = append(cryptoKeys, v...)
		}

		// Handle pagination for next loop, or break loop
		v, ok := res["nextPageToken"]
		if ok {
			params["pageToken"] = v.(string)
		}
		if !ok {
			break
		}
	}
	return cryptoKeys, nil
}

// flattenKMSKeysList flattens a list of crypto keys from a given crypto key ring
func flattenKMSKeysList(d *schema.ResourceData, config *transport_tpg.Config, keysList []interface{}, keyRingId string) ([]interface{}, error) {
	var keys []interface{}
	for _, k := range keysList {
		key := k.(map[string]interface{})
		parsedId, err := ParseKmsCryptoKeyId(key["name"].(string), config)
		if err != nil {
			return nil, err
		}

		data := map[string]interface{}{}
		// The google_kms_crypto_key resource and dataset set
		// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}})
		// and set name is set as just {{name}}.
		data["id"] = key["name"]
		data["name"] = parsedId.Name
		data["key_ring"] = keyRingId

		data["labels"] = flattenKMSCryptoKeyLabels(key["labels"], d, config)
		data["primary"] = flattenKMSCryptoKeyPrimary(key["primary"], d, config)
		data["purpose"] = flattenKMSCryptoKeyPurpose(key["purpose"], d, config)
		data["rotation_period"] = flattenKMSCryptoKeyRotationPeriod(key["rotationPeriod"], d, config)
		data["version_template"] = flattenKMSCryptoKeyVersionTemplate(key["versionTemplate"], d, config)
		data["destroy_scheduled_duration"] = flattenKMSCryptoKeyDestroyScheduledDuration(key["destroyScheduledDuration"], d, config)
		data["import_only"] = flattenKMSCryptoKeyImportOnly(key["importOnly"], d, config)
		data["crypto_key_backend"] = flattenKMSCryptoKeyCryptoKeyBackend(key["cryptoKeyBackend"], d, config)
		keys = append(keys, data)
	}

	return keys, nil
}
