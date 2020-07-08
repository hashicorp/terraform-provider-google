---
subcategory: "Game Servers"
layout: "google"
page_title: "Google: google_game_services_game_server_deployment_rollout"
sidebar_current: "docs-google-datasource-game-services-game-server-deployment-rollout"
description: |-
  Get the rollout state.
---

# google\_game\_services\_game\_server\_deployment\_rollout

Use this data source to get the rollout state. 

https://cloud.google.com/game-servers/docs/reference/rest/v1beta/GameServerDeploymentRollout

## Example Usage 


```hcl
data "google_game_services_game_server_deployment_rollout" "qa" {
    provider = google-beta
    deployment_id = "tf-test-deployment-s8sn12jt2c"
}
```

## Argument Reference

The following arguments are supported:


* `deployment_id` - (Required)
  The deployment to get the rollout state from. Only 1 rollout must be associated with each deployment.


## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `default_game_server_config` -
  This field points to the game server config that is
  applied by default to all realms and clusters. For example,
  `projects/my-project/locations/global/gameServerDeployments/my-game/configs/my-config`.


* `game_server_config_overrides` -
  The game_server_config_overrides contains the per game server config
  overrides. The overrides are processed in the order they are listed. As
  soon as a match is found for a cluster, the rest of the list is not
  processed.  Structure is documented below.

* `project` - The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


The `game_server_config_overrides` block contains:

* `realms_selector` -
  Selection by realms.  Structure is documented below.

* `config_version` -
  Version of the configuration.

The `realms_selector` block contains:

* `realms` -
  List of realms to match against.

* `id` - an identifier for the resource with format `projects/{{project}}/locations/global/gameServerDeployments/{{deployment_id}}/rollout`

* `name` -
  The resource id of the game server deployment
  eg: `projects/my-project/locations/global/gameServerDeployments/my-deployment/rollout`.
