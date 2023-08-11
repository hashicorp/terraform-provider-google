// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package corebilling_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBillingProjectInfo_update(t *testing.T) {
	t.Parallel()

	projectId := "tf-test-" + acctest.RandString(t, 10)
	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCoreBillingProjectInfoDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBillingProjectInfo_basic(projectId, orgId, billingAccount),
			},
			{
				ResourceName:      "google_billing_project_info.info",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingProjectInfo_basic(projectId, orgId, ""),
			},
			{
				ResourceName:      "google_billing_project_info.info",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBillingProjectInfo_basic(projectId, orgId, billingAccount),
			},
			{
				ResourceName:      "google_billing_project_info.info",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBillingProjectInfo_basic(projectId, orgId, billingAccountId string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id = "%s"
  name       = "%[1]s"
  org_id     = "%s"
  lifecycle {
    ignore_changes = [billing_account]
  }
}

resource "google_billing_project_info" "info" {
  project         = google_project.project.project_id
  billing_account = "%s"
}
`, projectId, orgId, billingAccountId)
}
