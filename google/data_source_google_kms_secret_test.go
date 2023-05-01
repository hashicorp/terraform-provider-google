package google

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"google.golang.org/api/cloudkms/v1"
)

func TestAccKmsSecret_basic(t *testing.T) {
	// Nested tests confuse VCR
	acctest.SkipIfVcr(t)
	t.Parallel()

	projectOrg := acctest.GetTestOrgFromEnv(t)
	projectBillingAccount := acctest.GetTestBillingAccountFromEnv(t)

	projectId := "tf-test-" + RandString(t, 10)
	keyRingName := fmt.Sprintf("tf-test-%s", RandString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", RandString(t, 10))

	plaintext := fmt.Sprintf("secret-%s", RandString(t, 10))
	aad := "plainaad"

	// The first test creates resources needed to encrypt plaintext and produce ciphertext
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
				Check: func(s *terraform.State) error {
					ciphertext, cryptoKeyId, err := testAccEncryptSecretDataWithCryptoKey(t, s, "google_kms_crypto_key.crypto_key", plaintext, "")

					if err != nil {
						return err
					}

					// The second test asserts that the data source has the correct plaintext, given the created ciphertext
					VcrTest(t, resource.TestCase{
						PreCheck:                 func() { acctest.AccTestPreCheck(t) },
						ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
						Steps: []resource.TestStep{
							{
								Config: testGoogleKmsSecret_datasource(cryptoKeyId.TerraformId(), ciphertext),
								Check:  resource.TestCheckResourceAttr("data.google_kms_secret.acceptance", "plaintext", plaintext),
							},
						},
					})

					return nil
				},
			},
			// With AAD
			{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
				Check: func(s *terraform.State) error {
					ciphertext, cryptoKeyId, err := testAccEncryptSecretDataWithCryptoKey(t, s, "google_kms_crypto_key.crypto_key", plaintext, aad)

					if err != nil {
						return err
					}

					// The second test asserts that the data source has the correct plaintext, given the created ciphertext
					VcrTest(t, resource.TestCase{
						PreCheck:                 func() { acctest.AccTestPreCheck(t) },
						ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
						Steps: []resource.TestStep{
							{
								Config: testGoogleKmsSecret_aadDatasource(cryptoKeyId.TerraformId(), ciphertext, base64.StdEncoding.EncodeToString([]byte(aad))),
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

func testAccEncryptSecretDataWithCryptoKey(t *testing.T, s *terraform.State, cryptoKeyResourceName, plaintext, aad string) (string, *KmsCryptoKeyId, error) {
	config := GoogleProviderConfig(t)

	rs, ok := s.RootModule().Resources[cryptoKeyResourceName]
	if !ok {
		return "", nil, fmt.Errorf("Resource not found: %s", cryptoKeyResourceName)
	}

	cryptoKeyId, err := ParseKmsCryptoKeyId(rs.Primary.Attributes["id"], config)

	if err != nil {
		return "", nil, err
	}

	kmsEncryptRequest := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	if aad != "" {
		kmsEncryptRequest.AdditionalAuthenticatedData = base64.StdEncoding.EncodeToString([]byte(aad))
	}

	encryptResponse, err := config.NewKmsClient(config.UserAgent).Projects.Locations.KeyRings.CryptoKeys.Encrypt(cryptoKeyId.CryptoKeyId(), kmsEncryptRequest).Do()

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

func testGoogleKmsSecret_aadDatasource(cryptoKeyTerraformId, ciphertext, aad string) string {
	return fmt.Sprintf(`
data "google_kms_secret" "acceptance" {
  crypto_key                    = "%s"
  ciphertext                    = "%s"
  additional_authenticated_data = "%s"
}
`, cryptoKeyTerraformId, ciphertext, aad)
}
