package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKmsKeyRingImportJob_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	return Nprintf(`
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
