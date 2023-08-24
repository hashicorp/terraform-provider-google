// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeBackendServiceSignedUrlKey_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeBackendServiceSignedUrlKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeBackendServiceSignedUrlKey_basic(context),
				Check:  testAccCheckComputeBackendServiceSignedUrlKeyCreatedProducer(t),
			},
		},
	})
}

func testAccComputeBackendServiceSignedUrlKey_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
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

func testAccCheckComputeBackendServiceSignedUrlKeyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		exists, err := checkComputeBackendServiceSignedUrlKeyExists(t, s)
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return err
		}
		if exists {
			return fmt.Errorf("ComputeBackendServiceSignedUrlKey still exists")
		}
		return nil
	}
}

func testAccCheckComputeBackendServiceSignedUrlKeyCreatedProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		exists, err := checkComputeBackendServiceSignedUrlKeyExists(t, s)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("expected ComputeBackendServiceSignedUrlKey to have been created")
		}
		return nil
	}
}

func checkComputeBackendServiceSignedUrlKeyExists(t *testing.T, s *terraform.State) (bool, error) {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_backend_service_signed_url_key" {
			continue
		}
		if strings.HasPrefix(name, "data.") {
			continue
		}

		config := acctest.GoogleProviderConfig(t)
		keyName := rs.Primary.Attributes["name"]

		url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/backendServices/{{backend_service}}")
		if err != nil {
			return false, err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
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
