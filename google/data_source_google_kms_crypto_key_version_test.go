package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleKmsCryptoKeyVersion_basic(t *testing.T) {
	asymSignKey := BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_SIGN")
	asymDecrKey := BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_DECRYPT")
	symKey := BootstrapKMSKey(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersion_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version.version", "version", "1"),
			},
			// Asymmetric keys should have a public key
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersion_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version.version", "public_key.#", "1"),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersion_basic(asymDecrKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version.version", "public_key.#", "1"),
			},
			// Symmetric key should have no public key
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersion_basic(symKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version.version", "public_key.#", "0"),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKeyVersion_basic(kmsKey string) string {
	return fmt.Sprintf(`
data "google_kms_crypto_key_version" "version" {
  crypto_key = "%s"
}
`, kmsKey)
}
