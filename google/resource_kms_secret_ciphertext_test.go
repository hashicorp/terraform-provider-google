package google

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/cloudkms/v1"
)

func TestAccKmsSecretCiphertext_basic(t *testing.T) {
	t.Parallel()

	kms := BootstrapKMSKey(t)

	plaintext := fmt.Sprintf("secret-%s", acctest.RandString(10))
	aad := "plainaad"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsSecretCiphertext(kms.CryptoKey.Name, plaintext),
				Check: func(s *terraform.State) error {
					plaintext, err := testAccDecryptSecretDataWithCryptoKey(s, kms.CryptoKey.Name, "google_kms_secret_ciphertext.acceptance", "")

					if err != nil {
						return err
					}

					return resource.TestCheckResourceAttr("google_kms_secret_ciphertext.acceptance", "plaintext", plaintext)(s)
				},
			},
			// With AAD
			{
				Config: testGoogleKmsSecretCiphertext_withAAD(kms.CryptoKey.Name, plaintext, aad),
				Check: func(s *terraform.State) error {
					plaintext, err := testAccDecryptSecretDataWithCryptoKey(s, kms.CryptoKey.Name, "google_kms_secret_ciphertext.acceptance", aad)

					if err != nil {
						return err
					}

					return resource.TestCheckResourceAttr("google_kms_secret_ciphertext.acceptance", "plaintext", plaintext)(s)
				},
			},
		},
	})
}

func testAccDecryptSecretDataWithCryptoKey(s *terraform.State, cryptoKeyId string, secretCiphertextResourceName, aad string) (string, error) {
	config := testAccProvider.Meta().(*Config)
	rs, ok := s.RootModule().Resources[secretCiphertextResourceName]
	if !ok {
		return "", fmt.Errorf("Resource not found: %s", secretCiphertextResourceName)
	}
	ciphertext, ok := rs.Primary.Attributes["ciphertext"]
	if !ok {
		return "", fmt.Errorf("Attribute 'ciphertext' not found in resource '%s'", secretCiphertextResourceName)
	}

	kmsDecryptRequest := &cloudkms.DecryptRequest{
		Ciphertext: ciphertext,
	}

	if aad != "" {
		kmsDecryptRequest.AdditionalAuthenticatedData = base64.StdEncoding.EncodeToString([]byte(aad))
	}

	decryptResponse, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Decrypt(cryptoKeyId, kmsDecryptRequest).Do()

	if err != nil {
		return "", fmt.Errorf("Error decrypting ciphertext: %s", err)
	}

	plaintextBytes, err := base64.StdEncoding.DecodeString(decryptResponse.Plaintext)

	if err != nil {
		return "", err
	}

	plaintext := string(plaintextBytes)
	log.Printf("[INFO] Successfully decrypted ciphertext and got plaintext: %s", plaintext)

	return plaintext, nil
}

func testGoogleKmsSecretCiphertext(cryptoKeyTerraformId, plaintext string) string {
	return fmt.Sprintf(`
resource "google_kms_secret_ciphertext" "acceptance" {
  crypto_key = "%s"
  plaintext  = "%s"
}
`, cryptoKeyTerraformId, plaintext)
}

func testGoogleKmsSecretCiphertext_withAAD(cryptoKeyTerraformId, plaintext, aad string) string {
	return fmt.Sprintf(`
resource "google_kms_secret_ciphertext" "acceptance" {
  crypto_key                    = "%s"
  plaintext                     = "%s"
  additional_authenticated_data = "%s"
}
`, cryptoKeyTerraformId, plaintext, aad)
}
