// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccKmsKeyRingImportJob_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsKeyRingImportJob_basic(context),
			},
			{
				ResourceName:            "google_kms_key_ring_import_job.import-job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key_ring", "import_job_id", "stateß"},
			},
		},
	})
}

func testGoogleKmsKeyRingImportJob_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_key_ring" "keyring" {
  name     = "tf-test-import-job-%{random_suffix}"
  location = "global"
}

resource "google_kms_key_ring_import_job" "import-job" {
  key_ring = google_kms_key_ring.keyring.id
  import_job_id = "my-import-job"

  import_method = "RSA_OAEP_3072_SHA1_AES_256"
  protection_level = "SOFTWARE"
}
`, context)
}
