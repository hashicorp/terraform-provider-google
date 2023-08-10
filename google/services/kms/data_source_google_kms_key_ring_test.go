// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsKeyRing_basic(t *testing.T) {
	kms := acctest.BootstrapKMSKey(t)

	keyParts := strings.Split(kms.KeyRing.Name, "/")
	keyRingId := keyParts[len(keyParts)-1]

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyRing_basic(keyRingId),
				Check:  resource.TestMatchResourceAttr("data.google_kms_key_ring.kms_key_ring", "id", regexp.MustCompile(kms.KeyRing.Name)),
			},
		},
	})
}

func testAccDataSourceGoogleKmsKeyRing_basic(keyRingName string) string {
	return fmt.Sprintf(`
data "google_kms_key_ring" "kms_key_ring" {
  name     = "%s"
  location = "global"
}
`, keyRingName)
}
