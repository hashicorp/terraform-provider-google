// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firebasedataconnect_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFirebaseDataConnectService_Update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirebaseDataConnectServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaseDataConnectService_update(context, "Original display name"),
			},
			{
				ResourceName:            "google_firebase_data_connect_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "service_id", "terraform_labels"},
			},
			{
				Config: testAccFirebaseDataConnectService_update(context, "Updated display name"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_firebase_data_connect_service.default", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_firebase_data_connect_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "service_id", "terraform_labels"},
			},
		},
	})
}

// TODO(b/394642094): Cover force deletion once it's supported
func testAccFirebaseDataConnectService_update(context map[string]interface{}, display_name string) string {
	context["display_name"] = display_name
	return acctest.Nprintf(`
# Enable Firebase Data Connect API
resource "google_project_service" "fdc" {
  project = "%{project_id}"
  service = "firebasedataconnect.googleapis.com"
  disable_on_destroy = false
}

# Create an FDC service
resource "google_firebase_data_connect_service" "default" {
  project = "%{project_id}"
  location = "us-central1"
  service_id = "tf-fdc-%{random_suffix}"
  display_name = "%{display_name}"

  depends_on = [google_project_service.fdc]
}
`, context)
}
