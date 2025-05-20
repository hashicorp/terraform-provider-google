// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(t *testing.T) {
	asymSignKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_SIGN")

	id := asymSignKey.CryptoKey.Name + "/latestCryptoKeyVersion"

	randomString := acctest.RandString(t, 10)
	filterNameFindsNoLatestCryptoKeyVersion := fmt.Sprintf("name:%s", randomString)
	filterNameFindEnabledLatestCryptoKeyVersion := "state:enabled"

	findsNoLatestCryptoKeyVersionId := fmt.Sprintf("%s/filter=%s", id, filterNameFindsNoLatestCryptoKeyVersion)
	findsEnabledLatestCryptoKeyVersionId := fmt.Sprintf("%s/filter=%s", id, filterNameFindEnabledLatestCryptoKeyVersion)

	context := map[string]interface{}{
		"crypto_key": asymSignKey.CryptoKey.Name,
		"filter":     "", // Can be overridden using 2nd argument to config funcs
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(context, ""),
				// Test will attempt to get the latest version from the list of cryptoKeyVersions, if the latest is not enabled it will return an error
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "id", id),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "crypto_key", asymSignKey.CryptoKey.Name),
					resource.TestMatchResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "version", regexp.MustCompile("[1-9]+[0-9]*")),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindEnabledLatestCryptoKeyVersion)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve the latest ENABLED cryptoKeyVersion in the bootstrapped KMS crypto key used by the test
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "id", findsEnabledLatestCryptoKeyVersionId),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "crypto_key", asymSignKey.CryptoKey.Name),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "state", "ENABLED"),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindsNoLatestCryptoKeyVersion)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve no latest version
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "id", findsNoLatestCryptoKeyVersionId),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.latest_version", "crypto_key", asymSignKey.CryptoKey.Name),
				),
				ExpectError: regexp.MustCompile("Error: No CryptoVersions found in crypto key"),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(context map[string]interface{}, filter string) string {
	context["filter"] = filter

	return acctest.Nprintf(`
data "google_kms_crypto_key_latest_version" "latest_version" {
	crypto_key = "%{crypto_key}"
	%{filter}
}
`, context)
}
