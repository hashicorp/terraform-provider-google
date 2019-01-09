package google

import (
	"fmt"
	"log"
	"os"
	"testing"

	"google.golang.org/api/cloudkms/v1"
)

var SharedKeyRing = "tftest-shared-keyring-1"
var SharedCyptoKey = "tftest-shared-key-1"

type bootstrappedKMS struct {
	*cloudkms.KeyRing
	*cloudkms.CryptoKey
}

/**
* BootstrapKMSkey will return a KMS key that can be used in tests that are
* testing KMS integration with other resources.
*
* This will either return an existing key or create one if it hasn't been created
* in the project yet. The motivation is because keyrings don't get deleted and we
* don't want a linear growth of disabled keyrings in a project. We also don't want
* to incur the overhead of creating a new project for each test that needs to use
* a KMS key.
**/
func BootstrapKMSKey(t *testing.T) bootstrappedKMS {
	if v := os.Getenv("TF_ACC"); v == "" {
		log.Println("Acceptance tests and bootstrapping skipped unless env 'TF_ACC' set")

		// If not running acceptance tests, return an empty object
		return bootstrappedKMS{
			&cloudkms.KeyRing{},
			&cloudkms.CryptoKey{},
		}
	}

	projectID := getTestProjectFromEnv()
	locationID := "global"
	keyRingParent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)
	keyRingName := fmt.Sprintf("%s/keyRings/%s", keyRingParent, SharedKeyRing)
	keyParent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", projectID, locationID, SharedKeyRing)
	keyName := fmt.Sprintf("%s/cryptoKeys/%s", keyParent, SharedCyptoKey)

	config := Config{
		Credentials: getTestCredsFromEnv(),
		Project:     getTestProjectFromEnv(),
		Region:      getTestRegionFromEnv(),
		Zone:        getTestZoneFromEnv(),
	}

	if err := config.loadAndValidate(); err != nil {
		t.Errorf("Unable to bootstrap KMS key: %s", err)
	}

	// Get or Create the hard coded shared keyring for testing
	kmsClient := config.clientKms
	keyRing, err := kmsClient.Projects.Locations.KeyRings.Get(keyRingName).Do()
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			keyRing, err = kmsClient.Projects.Locations.KeyRings.Create(keyRingParent, &cloudkms.KeyRing{}).
				KeyRingId(SharedKeyRing).Do()
			if err != nil {
				t.Errorf("Unable to bootstrap KMS key. Cannot create keyRing: %s", err)
			}
		} else {
			t.Errorf("Unable to bootstrap KMS key. Cannot retrieve keyRing: %s", err)
		}
	}

	if keyRing == nil {
		t.Fatalf("Unable to bootstrap KMS key. keyRing is nil!")
	}

	// Get or Create the hard coded, shared crypto key for testing
	cryptoKey, err := kmsClient.Projects.Locations.KeyRings.CryptoKeys.Get(keyName).Do()
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			newKey := cloudkms.CryptoKey{
				Purpose: "ENCRYPT_DECRYPT",
			}

			cryptoKey, err = kmsClient.Projects.Locations.KeyRings.CryptoKeys.Create(keyParent, &newKey).
				CryptoKeyId(SharedCyptoKey).Do()
			if err != nil {
				t.Errorf("Unable to bootstrap KMS key. Cannot create new CryptoKey: %s", err)
			}

		} else {
			t.Errorf("Unable to bootstrap KMS key. Cannot call CryptoKey service: %s", err)
		}
	}

	if cryptoKey == nil {
		t.Fatalf("Unable to bootstrap KMS key. CryptoKey is nil!")
	}

	return bootstrappedKMS{
		keyRing,
		cryptoKey,
	}
}
