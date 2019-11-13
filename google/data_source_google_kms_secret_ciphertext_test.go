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

	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)

	projectId := "terraform-" + acctest.RandString(10)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	plaintext := fmt.Sprintf("secret-%s", acctest.RandString(10))

	// The first test creates resources needed to encrypt plaintext and produce ciphertext
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
				Check: func(s *terraform.State) error {
					cryptoKeyId, err := getCryptoKeyId(s, "google_kms_crypto_key.crypto_key")

					if err != nil {
						return err
					}

					// The second test asserts that the data source created a ciphertext that can be decrypted to the correct plaintext
					resource.Test(t, resource.TestCase{
						PreCheck:  func() { testAccPreCheck(t) },
						Providers: testAccProviders,
						Steps: []resource.TestStep{
							{
								Config: testGoogleKmsSecretCiphertext_datasource(cryptoKeyId.terraformId(), plaintext),
								Check: func(s *terraform.State) error {
									plaintext, err := testAccDecryptSecretDataWithCryptoKey(s, cryptoKeyId, "data.google_kms_secret_ciphertext.acceptance")

									if err != nil {
										return err
									}

									return resource.TestCheckResourceAttr("data.google_kms_secret_ciphertext.acceptance", "plaintext", plaintext)(s)
								},
							},
						},
					})

					return nil
				},
			},
		},
	})
}

func getCryptoKeyId(s *terraform.State, cryptoKeyResourceName string) (*kmsCryptoKeyId, error) {
	config := testAccProvider.Meta().(*Config)
	rs, ok := s.RootModule().Resources[cryptoKeyResourceName]
	if !ok {
		return nil, fmt.Errorf("Resource not found: %s", cryptoKeyResourceName)
	}

	return parseKmsCryptoKeyId(rs.Primary.Attributes["id"], config)
}

func testAccDecryptSecretDataWithCryptoKey(s *terraform.State, cryptoKeyId *kmsCryptoKeyId, secretCiphertextResourceName string) (string, error) {
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

	decryptResponse, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Decrypt(cryptoKeyId.cryptoKeyId(), kmsDecryptRequest).Do()

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

func testGoogleKmsSecretCiphertext_datasource(cryptoKeyTerraformId, plaintext string) string {
	return fmt.Sprintf(`
data "google_kms_secret_ciphertext" "acceptance" {
  crypto_key = "%s"
  plaintext  = "%s"
}
`, cryptoKeyTerraformId, plaintext)
}
