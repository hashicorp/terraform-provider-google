package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccStorageHmacKey_update(t *testing.T) {
	t.Parallel()

	saName := fmt.Sprintf("%v%v", "service-account", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageHmacKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleStorageHmacKeyBasic(saName, "ACTIVE"),
			},
			{
				ResourceName:            "google_storage_hmac_key.key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
			{
				Config: testAccGoogleStorageHmacKeyBasic(saName, "INACTIVE"),
			},
			{
				ResourceName:            "google_storage_hmac_key.key",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
		},
	})
}

func testAccGoogleStorageHmacKeyBasic(saName, state string) string {
	return fmt.Sprintf(`
resource "google_service_account" "service_account" {
  account_id = "%s"
}

resource "google_storage_hmac_key" "key" {
	service_account_email = google_service_account.service_account.email
	state = "%s"
}
`, saName, state)
}
