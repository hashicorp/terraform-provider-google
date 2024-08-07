// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsKeyRings_basic(t *testing.T) {
	kms := acctest.BootstrapKMSKey(t)
	idPath := strings.Split(kms.KeyRing.Name, "/")
	location := idPath[3]
	randomString := acctest.RandString(t, 10)
	filterNameFindSharedKeys := "name:tftest-shared-"
	filterNameFindsNoKeys := fmt.Sprintf("name:%s", randomString)

	keyRingsID := fmt.Sprintf("projects/%s/locations/%s/keyRings", idPath[1], location)
	findSharedKeysId := fmt.Sprintf("%s/filter=%s", keyRingsID, filterNameFindSharedKeys)
	findsNoKeysId := fmt.Sprintf("%s/filter=%s", keyRingsID, filterNameFindsNoKeys)

	context := map[string]interface{}{
		"filter":   "", // Can be overridden using 2nd argument to config funcs
		"location": location,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyRings_basic(context, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_kms_key_rings.all_key_rings", "id", keyRingsID),
					resource.TestMatchResourceAttr("data.google_kms_key_rings.all_key_rings", "key_rings.#", regexp.MustCompile("[1-9]+[0-9]*")),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsKeyRings_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindSharedKeys)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve the bootstrapped KMS key rings used by the test
					resource.TestCheckResourceAttr("data.google_kms_key_rings.all_key_rings", "id", findSharedKeysId),
					resource.TestMatchResourceAttr("data.google_kms_key_rings.all_key_rings", "key_rings.#", regexp.MustCompile("[1-9]+[0-9]*")),
				),
			},
			{
				Config: testAccDataSourceGoogleKmsKeyRings_basic(context, fmt.Sprintf("filter = \"%s\"", filterNameFindsNoKeys)),
				Check: resource.ComposeTestCheckFunc(
					// This filter should retrieve no keys
					resource.TestCheckResourceAttr("data.google_kms_key_rings.all_key_rings", "id", findsNoKeysId),
					resource.TestCheckResourceAttr("data.google_kms_key_rings.all_key_rings", "key_rings.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleKmsKeyRings_basic(context map[string]interface{}, filter string) string {
	context["filter"] = filter

	return acctest.Nprintf(`
data "google_kms_key_rings" "all_key_rings" {
  location = "%{location}"
  %{filter}
}
`, context)
}
