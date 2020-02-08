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

func TestAccKmsSecret_basic(t *testing.T) {
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
					ciphertext, cryptoKeyId, err := testAccEncryptSecretDataWithCryptoKey(s, "google_kms_crypto_key.crypto_key", plaintext)

					if err != nil {
						return err
					}

					// The second test asserts that the data source has the correct plaintext, given the created ciphertext
					resource.Test(t, resource.TestCase{
						PreCheck:  func() { testAccPreCheck(t) },
						Providers: testAccProviders,
						Steps: []resource.TestStep{
							{
								Config: testGoogleKmsSecret_datasource(cryptoKeyId.terraformId(), ciphertext),
								Check:  resource.TestCheckResourceAttr("data.google_kms_secret.acceptance", "plaintext", plaintext),
							},
						},
					})

					return nil
				},
			},
		},
	})
}

func testAccEncryptSecretDataWithCryptoKey(s *terraform.State, cryptoKeyResourceName, plaintext string) (string, *kmsCryptoKeyId, error) {
	config := testAccProvider.Meta().(*Config)

	rs, ok := s.RootModule().Resources[cryptoKeyResourceName]
	if !ok {
		return "", nil, fmt.Errorf("Resource not found: %s", cryptoKeyResourceName)
	}

	cryptoKeyId, err := parseKmsCryptoKeyId(rs.Primary.Attributes["id"], config)

	if err != nil {
		return "", nil, err
	}

	kmsEncryptRequest := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	encryptResponse, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Encrypt(cryptoKeyId.cryptoKeyId(), kmsEncryptRequest).Do()

	if err != nil {
		return "", nil, fmt.Errorf("Error encrypting plaintext: %s", err)
	}

	log.Printf("[INFO] Successfully encrypted plaintext and got ciphertext: %s", encryptResponse.Ciphertext)

	return encryptResponse.Ciphertext, cryptoKeyId, nil
}

func testGoogleKmsSecret_datasource(cryptoKeyTerraformId, ciphertext string) string {
	return fmt.Sprintf(`
data "google_kms_secret" "acceptance" {
  crypto_key = "%s"
  ciphertext = "%s"
}
`, cryptoKeyTerraformId, ciphertext)
}
