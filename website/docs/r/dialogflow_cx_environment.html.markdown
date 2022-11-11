---
subcategory: "Dialogflow CX"
page_title: "Google: google_dialogflow_cx_environment"
description: |-
  Represents an environment for an agent.
---

# google\_dialogflow\_cx\_environment

Represents an environment for an agent. You can create multiple versions of your agent and publish them to separate environments.
When you edit an agent, you are editing the draft agent. At any point, you can save the draft agent as an agent version, which is an immutable snapshot of your agent.
When you save the draft agent, it is published to the default environment. When you create agent versions, you can publish them to custom environments. You can create a variety of custom environments for testing, development, production, etc.


To get more information about Environment, see:

* [API documentation](https://cloud.google.com/dialogflow/cx/docs/reference/rest/v3/projects.locations.agents.environments)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/dialogflow/cx/docs)

<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_working_dir=dialogflowcx_environment_full&cloudshell_image=gcr.io%2Fgraphite-cloud-shell-images%2Fterraform%3Alatest&open_in_editor=main.tf&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Dialogflowcx Environment Full


```hcl
resource "google_dialogflow_cx_agent" "agent" {
  display_name = "dialogflowcx-agent"
  location = "global"
  default_language_code = "en"
  supported_language_codes = ["fr","de","es"]
  time_zone = "America/New_York"
  description = "Example description."
  avatar_uri = "https://cloud.google.com/_static/images/cloud/icons/favicons/onecloud/super_cloud.png"
  enable_stackdriver_logging = true
  enable_spell_correction    = true
	speech_to_text_settings {
		enable_speech_adaptation = true
	}
}

resource "google_dialogflow_cx_version" "version_1" {
  parent       = google_dialogflow_cx_agent.agent.start_flow
  display_name = "1.0.0"
  description  = "version 1.0.0"
}

resource "google_dialogflow_cx_environment" "development" {
  parent       = google_dialogflow_cx_agent.agent.id
  display_name = "Development"
  description  = "Development Environment"
  version_configs {
    version = google_dialogflow_cx_version.version_1.id
  }
}
```

## Argument Reference

The following arguments are supported:


* `display_name` -
  (Required)
  The human-readable name of the environment (unique in an agent). Limit of 64 characters.

* `version_configs` -
  (Required)
  A list of configurations for flow versions. You should include version configs for all flows that are reachable from [Start Flow][Agent.start_flow] in the agent. Otherwise, an error will be returned.
  Structure is [documented below](#nested_version_configs).


<a name="nested_version_configs"></a>The `version_configs` block supports:

* `version` -
  (Required)
  Format: projects/{{project}}/locations/{{location}}/agents/{{agent}}/flows/{{flow}}/versions/{{version}}.

- - -


* `description` -
  (Optional)
  The human-readable description of the environment. The maximum length is 500 characters. If exceeded, the request is rejected.

* `parent` -
  (Optional)
  The Agent to create an Environment for.
  Format: projects/<Project ID>/locations/<Location ID>/agents/<Agent ID>.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{parent}}/environments/{{name}}`

* `name` -
  The name of the environment.

* `update_time` -
  Update time of this environment. A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".


## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 40 minutes.
- `update` - Default is 40 minutes.
- `delete` - Default is 20 minutes.

## Import


Environment can be imported using any of these accepted formats:

```
$ terraform import google_dialogflow_cx_environment.default {{parent}}/environments/{{name}}
$ terraform import google_dialogflow_cx_environment.default {{parent}}/{{name}}
```
