// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudkms/v1"
)

func TestAccKmsSecretCiphertext_basic(t *testing.T) {
	t.Parallel()

	kms := acctest.BootstrapKMSKey(t)

	plaintext := fmt.Sprintf("secret-%s", acctest.RandString(t, 10))
	aad := "plainaad"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsSecretCiphertext(kms.CryptoKey.Name, plaintext),
				Check: func(s *terraform.State) error {
					plaintext, err := testAccDecryptSecretDataWithCryptoKey(t, s, kms.CryptoKey.Name, "google_kms_secret_ciphertext.acceptance", "")

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
					plaintext, err := testAccDecryptSecretDataWithCryptoKey(t, s, kms.CryptoKey.Name, "google_kms_secret_ciphertext.acceptance", aad)

					if err != nil {
						return err
					}

					return resource.TestCheckResourceAttr("google_kms_secret_ciphertext.acceptance", "plaintext", plaintext)(s)
				},
			},
		},
	})
}

func testAccDecryptSecretDataWithCryptoKey(t *testing.T, s *terraform.State, cryptoKeyId string, secretCiphertextResourceName, aad string) (string, error) {
	config := acctest.GoogleProviderConfig(t)
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

	decryptResponse, err := config.NewKmsClient(config.UserAgent).Projects.Locations.KeyRings.CryptoKeys.Decrypt(cryptoKeyId, kmsDecryptRequest).Do()

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
