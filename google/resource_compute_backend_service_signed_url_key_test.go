package google

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeBackendServiceSignedUrlKey_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceSignedUrlKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendServiceSignedUrlKey_basic(context),
				Check:  testAccCheckComputeBackendServiceSignedUrlKeyCreated,
			},
		},
	})
}

func testAccComputeBackendServiceSignedUrlKey_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_backend_service_signed_url_key" "backend_key" {
  name            = "testkey-%{random_suffix}"
  key_value       = "iAmAFakeKeyRandomBytes=="
  backend_service = google_compute_backend_service.test_bs.name
}

resource "google_compute_backend_service" "test_bs" {
  name          = "testbs-%{random_suffix}"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "test-check-%{random_suffix}"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, context)
}

func testAccCheckComputeBackendServiceSignedUrlKeyDestroy(s *terraform.State) error {
	exists, err := checkComputeBackendServiceSignedUrlKeyExists(s)
	if err != nil && !isGoogleApiErrorWithCode(err, 404) {
		return err
	}
	if exists {
		return fmt.Errorf("ComputeBackendServiceSignedUrlKey still exists")
	}
	return nil
}

func testAccCheckComputeBackendServiceSignedUrlKeyCreated(s *terraform.State) error {
	exists, err := checkComputeBackendServiceSignedUrlKeyExists(s)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("expected ComputeBackendServiceSignedUrlKey to have been created")
	}
	return nil
}

func checkComputeBackendServiceSignedUrlKeyExists(s *terraform.State) (bool, error) {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_backend_service_signed_url_key" {
			continue
		}
		if strings.HasPrefix(name, "data.") {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		keyName := rs.Primary.Attributes["name"]

		url, err := replaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/backendServices/{{backend_service}}")
		if err != nil {
			return false, err
		}

		res, err := sendRequest(config, "GET", "", url, nil)
		if err != nil {
			return false, err
		}
		policyRaw, ok := res["cdnPolicy"]
		if !ok {
			return false, nil
		}

		policy := policyRaw.(map[string]interface{})
		keyNames, ok := policy["signedUrlKeyNames"]
		if !ok {
			return false, nil
		}

		// Because the sensitive key value is not returned, all we can do is verify a
		// key with this name exists and assume the key value hasn't been changed.
		for _, k := range keyNames.([]interface{}) {
			if k.(string) == keyName {
				// Just return empty map to indicate key was found
				return true, nil
			}
		}
	}

	return false, nil
}
