package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleCloudBuildTrigger_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCloudBuildTrigger_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_cloudbuild_trigger.foo", "google_cloudbuild_trigger.test-trigger"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCloudBuildTrigger_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloudbuild_trigger" "test-trigger" {
	location = "us-central1"
	name        = "manual-build%{random_suffix}"
	trigger_template {
		branch_name = "main"
		repo_name   = "my-repo"
	}
	
	substitutions = {
		_FOO = "bar"
		_BAZ = "qux"
	}
	
	filename = "cloudbuild.yaml"
}

data "google_cloudbuild_trigger" "foo" {
	location = google_cloudbuild_trigger.test-trigger.location
	trigger_id = google_cloudbuild_trigger.test-trigger.trigger_id
}`, context)

}
