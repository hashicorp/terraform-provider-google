// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oslogin_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccOSLoginSSHPublicKey_osLoginSshKeyExpiry(t *testing.T) {
	// Uses time provider
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckOSLoginSSHPublicKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOSLoginSSHPublicKey_osLoginSshKeyExpiry(context),
			},
			{
				ResourceName:            "google_os_login_ssh_public_key.cache",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user", "project"},
			},
		},
	})
}

func testAccOSLoginSSHPublicKey_osLoginSshKeyExpiry(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}
resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "oslogin" {
  project = google_project.project.project_id
  service = "oslogin.googleapis.com"
  disable_dependent_services = true
}

data "google_client_openid_userinfo" "me" {
}

resource "time_offset" "expiry" {
  offset_hours = 1
}

resource "google_os_login_ssh_public_key" "cache" {
  project = google_project.project.project_id
  user    =  data.google_client_openid_userinfo.me.email
  key     = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPM4pxpbPpjuBocS6qlW0BHRYgH5xmv/yVrANZR9lc1N"
  expiration_time_usec = time_offset.expiry.unix * 1000000
  depends_on = [
	google_project_service.compute,
	google_project_service.oslogin,
  ]
}
`, context)
}
