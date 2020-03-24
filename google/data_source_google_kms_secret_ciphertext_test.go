package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataKmsSecretCiphertext_basic(t *testing.T) {
	t.Parallel()

	kms := BootstrapKMSKey(t)

	plaintext := fmt.Sprintf("secret-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsSecretCiphertext_datasource(kms.CryptoKey.Name, plaintext),
				Check: func(s *terraform.State) error {
					plaintext, err := testAccDecryptSecretDataWithCryptoKey(s, kms.CryptoKey.Name, "data.google_kms_secret_ciphertext.acceptance", "")

					if err != nil {
						return err
					}

					return resource.TestCheckResourceAttr("data.google_kms_secret_ciphertext.acceptance", "plaintext", plaintext)(s)
				},
			},
		},
	})
}

func testGoogleKmsSecretCiphertext_datasource(cryptoKeyTerraformId, plaintext string) string {
	return fmt.Sprintf(`
data "google_kms_secret_ciphertext" "acceptance" {
  crypto_key = "%s"
  plaintext  = "%s"
}
`, cryptoKeyTerraformId, plaintext)
}
