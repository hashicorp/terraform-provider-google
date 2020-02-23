package google

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleKmsKeyRing_basic(t *testing.T) {
	kms := BootstrapKMSKey(t)

	keyParts := strings.Split(kms.KeyRing.Name, "/")
	keyRingId := keyParts[len(keyParts)-1]

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyRing_basic(keyRingId),
				Check:  resource.TestMatchResourceAttr("data.google_kms_key_ring.kms_key_ring", "self_link", regexp.MustCompile(kms.KeyRing.Name)),
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
