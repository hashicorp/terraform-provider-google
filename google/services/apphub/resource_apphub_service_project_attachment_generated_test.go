// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package apphub_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccApphubServiceProjectAttachment_serviceProjectAttachmentBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"host_project":  envvar.GetTestProjectFromEnv(),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckApphubServiceProjectAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApphubServiceProjectAttachment_serviceProjectAttachmentBasicExample(context),
			},
			{
				ResourceName:            "google_apphub_service_project_attachment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
		},
	})
}

func testAccApphubServiceProjectAttachment_serviceProjectAttachmentBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apphub_service_project_attachment" "example" {
  service_project_attachment_id = google_project.service_project.project_id
  depends_on = [time_sleep.wait_120s]
}

resource "google_project" "service_project" {
  project_id ="tf-test-project-1%{random_suffix}"
  name = "Service Project"
  org_id = "%{org_id}"
}

resource "time_sleep" "wait_120s" {
  depends_on = [google_project.service_project]

  create_duration = "120s"
}
`, context)
}

func TestAccApphubServiceProjectAttachment_serviceProjectAttachmentFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"host_project":  envvar.GetTestProjectFromEnv(),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckApphubServiceProjectAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApphubServiceProjectAttachment_serviceProjectAttachmentFullExample(context),
			},
			{
				ResourceName:            "google_apphub_service_project_attachment.example2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_project_attachment_id"},
			},
		},
	})
}

func testAccApphubServiceProjectAttachment_serviceProjectAttachmentFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apphub_service_project_attachment" "example2" {
  service_project_attachment_id = google_project.service_project_full.project_id
  service_project = google_project.service_project_full.project_id
  depends_on = [time_sleep.wait_120s]
}

resource "google_project" "service_project_full" {
  project_id ="tf-test-project-1%{random_suffix}"
  name = "Service Project Full"
  org_id = "%{org_id}"
}

resource "time_sleep" "wait_120s" {
  depends_on = [google_project.service_project_full]

  create_duration = "120s"
}
`, context)
}

func testAccCheckApphubServiceProjectAttachmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apphub_service_project_attachment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ApphubBasePath}}projects/{{project}}/locations/global/serviceProjectAttachments/{{service_project_attachment_id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ApphubServiceProjectAttachment still exists at %s", url)
			}
		}

		return nil
	}
}
