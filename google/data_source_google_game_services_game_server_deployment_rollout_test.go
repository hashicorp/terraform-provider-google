package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGameServicesGameServerDeploymentRollout_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGameServicesGameServerDeploymentRolloutDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGameServicesGameServerDeploymentRollout_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_game_services_game_server_deployment_rollout.qa", "google_game_services_game_server_deployment_rollout.foo"),
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
    fleet_spec = jsonencode({ "replicas" : 1, "scheduling" : "Packed", "template" : { "metadata" : { "name" : "tf-test-game-server-template" }, "spec" : { "ports": [{"name": "default", "portPolicy": "Dynamic", "containerPort": 7654, "protocol": "UDP"}], "template" : { "spec" : { "containers" : [{ "name" : "simple-udp-server", "image" : "gcr.io/agones-images/udp-server:0.14" }] } } } } })

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
