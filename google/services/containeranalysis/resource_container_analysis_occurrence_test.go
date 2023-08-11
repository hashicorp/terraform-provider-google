// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package containeranalysis_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"crypto/sha512"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"google.golang.org/api/cloudkms/v1"
)

const testAttestationOccurrenceImageUrl = "gcr.io/cloud-marketplace/google/ubuntu1804"
const testAttestationOccurrenceImageDigest = "sha256:3593cd4ac7d782d460dc86ba9870a3beaf81c8f5cdbcc8880bf9a5ef6af10c5a"
const testAttestationOccurrencePayloadTemplate = "test-fixtures/generated_payload.json.tmpl"

var testAttestationOccurrenceFullImagePath = fmt.Sprintf("%s@%s", testAttestationOccurrenceImageUrl, testAttestationOccurrenceImageDigest)

func getTestOccurrenceAttestationPayload(t *testing.T) string {
	payloadTmpl, err := ioutil.ReadFile(testAttestationOccurrencePayloadTemplate)
	if err != nil {
		t.Fatal(err.Error())
	}
	return fmt.Sprintf(string(payloadTmpl),
		testAttestationOccurrenceImageUrl,
		testAttestationOccurrenceImageDigest)
}

func getSignedTestOccurrenceAttestationPayload(
	t *testing.T, config *transport_tpg.Config,
	signingKey acctest.BootstrappedKMS, rawPayload string) string {
	pbytes := []byte(rawPayload)
	ssum := sha512.Sum512(pbytes)
	hashed := base64.StdEncoding.EncodeToString(ssum[:])
	signed, err := config.NewKmsClient(config.UserAgent).Projects.Locations.KeyRings.CryptoKeys.
		CryptoKeyVersions.AsymmetricSign(
		fmt.Sprintf("%s/cryptoKeyVersions/1", signingKey.CryptoKey.Name),
		&cloudkms.AsymmetricSignRequest{
			Digest: &cloudkms.Digest{
				Sha512: hashed,
			},
		}).Do()
	if err != nil {
		t.Fatalf("Unable to sign attestation payload with KMS key: %s", err)
	}

	return signed.Signature
}

func TestAccContainerAnalysisOccurrence_basic(t *testing.T) {
	t.Parallel()
	randSuffix := acctest.RandString(t, 10)

	config := acctest.BootstrapConfig(t)
	if config == nil {
		return
	}

	signKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_SIGN")
	payload := getTestOccurrenceAttestationPayload(t)
	signed := getSignedTestOccurrenceAttestationPayload(t, config, signKey, payload)
	params := map[string]interface{}{
		"random_suffix": randSuffix,
		"image_url":     testAttestationOccurrenceFullImagePath,
		"key_ring":      tpgresource.GetResourceNameFromSelfLink(signKey.KeyRing.Name),
		"crypto_key":    tpgresource.GetResourceNameFromSelfLink(signKey.CryptoKey.Name),
		"payload":       base64.StdEncoding.EncodeToString([]byte(payload)),
		"signature":     base64.StdEncoding.EncodeToString([]byte(signed)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAnalysisNoteDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAnalysisOccurence_basic(params),
			},
			{
				ResourceName:      "google_container_analysis_occurrence.occurrence",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContainerAnalysisOccurrence_multipleSignatures(t *testing.T) {
	t.Parallel()
	randSuffix := acctest.RandString(t, 10)

	config := acctest.BootstrapConfig(t)
	if config == nil {
		return
	}

	payload := getTestOccurrenceAttestationPayload(t)
	key1 := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ASYMMETRIC_SIGN", "global", "tf-bootstrap-binauthz-key1")
	signature1 := getSignedTestOccurrenceAttestationPayload(t, config, key1, payload)

	key2 := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ASYMMETRIC_SIGN", "global", "tf-bootstrap-binauthz-key2")
	signature2 := getSignedTestOccurrenceAttestationPayload(t, config, key2, payload)

	paramsMultipleSignatures := map[string]interface{}{
		"random_suffix": randSuffix,
		"image_url":     testAttestationOccurrenceFullImagePath,
		"key_ring":      tpgresource.GetResourceNameFromSelfLink(key1.KeyRing.Name),
		"payload":       base64.StdEncoding.EncodeToString([]byte(payload)),
		"key1":          tpgresource.GetResourceNameFromSelfLink(key1.CryptoKey.Name),
		"signature1":    base64.StdEncoding.EncodeToString([]byte(signature1)),
		"key2":          tpgresource.GetResourceNameFromSelfLink(key2.CryptoKey.Name),
		"signature2":    base64.StdEncoding.EncodeToString([]byte(signature2)),
	}
	paramsSingle := map[string]interface{}{
		"random_suffix": randSuffix,
		"image_url":     testAttestationOccurrenceFullImagePath,
		"key_ring":      tpgresource.GetResourceNameFromSelfLink(key1.KeyRing.Name),
		"crypto_key":    tpgresource.GetResourceNameFromSelfLink(key1.CryptoKey.Name),
		"payload":       base64.StdEncoding.EncodeToString([]byte(payload)),
		"signature":     base64.StdEncoding.EncodeToString([]byte(signature1)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAnalysisNoteDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAnalysisOccurence_multipleSignatures(paramsMultipleSignatures),
			},
			{
				ResourceName:      "google_container_analysis_occurrence.occurrence",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContainerAnalysisOccurence_basic(paramsSingle),
			},
			{
				ResourceName:      "google_container_analysis_occurrence.occurrence",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccContainerAnalysisOccurence_basic(params map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_binary_authorization_attestor" "attestor" {
  name = "test-attestor%{random_suffix}"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
    public_keys {
      id = data.google_kms_crypto_key_version.version.id
      pkix_public_key {
        public_key_pem      = data.google_kms_crypto_key_version.version.public_key[0].pem
        signature_algorithm = data.google_kms_crypto_key_version.version.public_key[0].algorithm
      }
    }
  }
}

resource "google_container_analysis_note" "note" {
  name = "test-attestor-note%{random_suffix}"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}

data "google_kms_key_ring" "keyring" {
  name = "%{key_ring}"
  location = "global"
}

data "google_kms_crypto_key" "crypto-key" {
  name     = "%{crypto_key}"
  key_ring = data.google_kms_key_ring.keyring.id
}

data "google_kms_crypto_key_version" "version" {
  crypto_key = data.google_kms_crypto_key.crypto-key.id
}

resource "google_container_analysis_occurrence" "occurrence" {
  resource_uri = "%{image_url}"
  note_name = google_container_analysis_note.note.id

  attestation {
    serialized_payload = "%{payload}"
    signatures {
      public_key_id = data.google_kms_crypto_key_version.version.id
      signature = "%{signature}"
    }
  }
}
`, params)
}

func testAccContainerAnalysisOccurence_multipleSignatures(params map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_binary_authorization_attestor" "attestor" {
  name = "test-attestor%{random_suffix}"
  attestation_authority_note {
    note_reference = google_container_analysis_note.note.name
    public_keys {
      id = data.google_kms_crypto_key_version.version-key1.id
      pkix_public_key {
        public_key_pem      = data.google_kms_crypto_key_version.version-key1.public_key[0].pem
        signature_algorithm = data.google_kms_crypto_key_version.version-key1.public_key[0].algorithm
      }
    }

		public_keys {
      id = data.google_kms_crypto_key_version.version-key2.id
      pkix_public_key {
        public_key_pem      = data.google_kms_crypto_key_version.version-key2.public_key[0].pem
        signature_algorithm = data.google_kms_crypto_key_version.version-key2.public_key[0].algorithm
      }
    }
  }
}

resource "google_container_analysis_note" "note" {
  name = "test-attestor-note%{random_suffix}"
  attestation_authority {
    hint {
      human_readable_name = "Attestor Note"
    }
  }
}

data "google_kms_key_ring" "keyring" {
  name = "%{key_ring}"
  location = "global"
}

data "google_kms_crypto_key" "crypto-key1" {
  name     = "%{key1}"
  key_ring = data.google_kms_key_ring.keyring.id
}

data "google_kms_crypto_key" "crypto-key2" {
  name     = "%{key2}"
  key_ring = data.google_kms_key_ring.keyring.id
}

data "google_kms_crypto_key_version" "version-key1" {
  crypto_key = data.google_kms_crypto_key.crypto-key1.id
}

data "google_kms_crypto_key_version" "version-key2" {
  crypto_key = data.google_kms_crypto_key.crypto-key2.id
}

resource "google_container_analysis_occurrence" "occurrence" {
  resource_uri = "%{image_url}"
  note_name = google_container_analysis_note.note.id

  attestation {
    serialized_payload = "%{payload}"
    signatures {
      public_key_id = data.google_kms_crypto_key_version.version-key1.id
      signature = "%{signature1}"
    }

		signatures {
      public_key_id = data.google_kms_crypto_key_version.version-key2.id
      signature = "%{signature2}"
    }
  }
}
`, params)
}
