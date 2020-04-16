package google

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleKmsCryptoKey_basic(t *testing.T) {
	kms := BootstrapKMSKey(t)

	// Name in the KMS client is in the format projects/<project>/locations/<location>/keyRings/<keyRingName>/cryptoKeys/<keyId>
	keyParts := strings.Split(kms.CryptoKey.Name, "/")
	cryptoKeyId := keyParts[len(keyParts)-1]

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKey_basic(kms.KeyRing.Name, cryptoKeyId),
				Check:  resource.TestMatchResourceAttr("data.google_kms_crypto_key.kms_crypto_key", "self_link", regexp.MustCompile(kms.CryptoKey.Name)),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKey_basic(keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
data "google_kms_crypto_key" "kms_crypto_key" {
  key_ring = "%s"
  name     = "%s"
}
`, keyRingName, cryptoKeyName)
}
