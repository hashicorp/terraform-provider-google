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

func TestAccDataSourceGoogleKmsCryptoKeyVersions_basic(t *testing.T) {
	asymSignKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_SIGN")

	id := asymSignKey.CryptoKey.Name + "/cryptoKeyVersions"

	randomString := acctest.RandString(t, 10)
	filterNameFindSharedCryptoKeyVersions := "name:tftest-shared-"
	filterNameFindsNoCryptoKeyVersions := fmt.Sprintf("name:%s", randomString)
	filterNameFindEnabledCryptoKeyVersions := "state:enabled"
	filterNameFindDisabledCryptoKeyVersions := "state:disabled"

	findSharedCryptoKeyVersionsId := fmt.Sprintf("%s/filter=%s", id, filterNameFindSharedCryptoKeyVersions)
	findsNoCryptoKeyVersionsId := fmt.Sprintf("%s/filter=%s", id, filterNameFindsNoCryptoKeyVersions)
	findsEnabledCryptoKeyVersionsId := fmt.Sprintf("%s/filter=%s", id, filterNameFindEnabledCryptoKeyVersions)
	findsDisabledCryptoKeyVersionsId := fmt.Sprintf("%s/filter=%s", id, filterNameFindDisabledCryptoKeyVersions)

	context := map[string]interface{}{
		"crypto_key": asymSignKey.CryptoKey.Name,
		"filter":     "", // Can be overridden using 2nd argument to config funcs
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersions_basic(context, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "id", id),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "crypto_key", asymSignKey.CryptoKey.Name),
					resource.TestMatchResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "versions.#", regexp.MustCompile("[1-9]+[0-9]*")),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersions_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindSharedCryptoKeyVersions)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve cryptoKeyVersions in the bootstrapped KMS crypto key used by the test
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "id", findSharedCryptoKeyVersionsId),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "crypto_key", asymSignKey.CryptoKey.Name),
					resource.TestMatchResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "versions.#", regexp.MustCompile("[1-9]+[0-9]*")),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersions_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindsNoCryptoKeyVersions)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve no cryptoKeyVersions
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "id", findsNoCryptoKeyVersionsId),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "crypto_key", asymSignKey.CryptoKey.Name),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "versions.#", "0"),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersions_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindEnabledCryptoKeyVersions)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve versions that are enabled
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "id", findsEnabledCryptoKeyVersionsId),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "crypto_key", asymSignKey.CryptoKey.Name),
					resource.TestMatchResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "versions.#", regexp.MustCompile("[1-9]+[0-9]*")),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersions_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindDisabledCryptoKeyVersions)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve versions that are disabled
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "id", findsDisabledCryptoKeyVersionsId),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "crypto_key", asymSignKey.CryptoKey.Name),
					resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.all_versions_in_key", "versions.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKeyVersions_basic(context map[string]interface{}, filter string) string {
	context["filter"] = filter

	return acctest.Nprintf(`
data "google_kms_crypto_key_versions" "all_versions_in_key" {
	crypto_key = "%{crypto_key}"
	%{filter}
}
`, context)
}
