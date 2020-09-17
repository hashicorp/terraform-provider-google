package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGameServicesGameServerDeploymentRollout_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGameServicesGameServerDeploymentRolloutDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGameServicesGameServerDeploymentRollout_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_game_services_game_server_deployment_rollout.qa", "google_game_services_game_server_deployment_rollout.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGameServicesGameServerDeploymentRollout_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_game_services_game_server_deployment" "default" {
  deployment_id  = "tf-test-deployment-%{random_suffix}"
  description = "a deployment description"
}

resource "google_game_services_game_server_config" "default" {
  config_id     = "tf-test-config-%{random_suffix}"
  deployment_id = google_game_services_game_server_deployment.default.deployment_id
  description   = "a config description"

  fleet_configs {
    name       = "some-non-guid"
    fleet_spec = jsonencode({ "replicas" : 1, "scheduling" : "Packed", "template" : { "metadata" : { "name" : "tf-test-game-server-template" }, "spec" : { "template" : { "spec" : { "containers" : [{ "name" : "simple-udp-server", "image" : "gcr.io/agones-images/udp-server:0.14" }] } } } } })

    // Alternate usage:
    // fleet_spec = file(fleet_configs.json)
  }
}

resource "google_game_services_game_server_deployment_rollout" "foo" {
  deployment_id              = google_game_services_game_server_deployment.default.deployment_id
  default_game_server_config = google_game_services_game_server_config.default.name
}

data "google_game_services_game_server_deployment_rollout" "qa" {
    deployment_id = google_game_services_game_server_deployment_rollout.foo.deployment_id
}
`, context)
}
