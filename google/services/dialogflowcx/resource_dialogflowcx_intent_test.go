// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDialogflowCXIntent_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowCXIntent_basic(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_intent.my_intent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccDialogflowCXIntent_full(context),
			},
			{
				ResourceName:            "google_dialogflow_cx_intent.my_intent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccDialogflowCXIntent_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dialogflow_cx_agent" "agent_intent" {
		display_name = "tf-test-%{random_suffix}"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["fr","de","es"]
		time_zone = "America/New_York"
		description = "Description 1."
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
	}
    
	resource "google_dialogflow_cx_intent" "my_intent" {
        parent       = google_dialogflow_cx_agent.agent_intent.id
        display_name = "Example"
        priority     = 1
        description  = "Intent example"
        training_phrases {
            parts {
                text = "training"
            }

            parts {
                text = "phrase"
            }

            parts {
                text = "example"
            }

			repeat_count = 1
        }

        parameters {
            id          = "param1"
            entity_type = "projects/-/locations/-/agents/-/entityTypes/sys.date"
        }

        labels  = {
            label1 = "value1",
            label2 = "value2"
        } 
    } 
    `, context)
}

func testAccDialogflowCXIntent_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_dialogflow_cx_agent" "agent_intent" {
		display_name = "tf-test-%{random_suffix}update"
		location = "global"
		default_language_code = "en"
		supported_language_codes = ["no"]
		time_zone = "Europe/London"
		description = "Description 2!"
		avatar_uri = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo-2.png"
		enable_stackdriver_logging = true
        enable_spell_correction    = true
		speech_to_text_settings {
			enable_speech_adaptation = true
		}
	}
    
	resource "google_dialogflow_cx_intent" "my_intent" {
        parent       = google_dialogflow_cx_agent.agent_intent.id
        display_name = "Example"
        priority     = 1
        description  = "Intent example"
        training_phrases {
            parts {
                text = "training"
            }

            parts {
                text = "phrase"
            }

            parts {
                text = "example"
            }

			repeat_count = 1
        }

        parameters {
            id          = "param1"
            entity_type = "projects/-/locations/-/agents/-/entityTypes/sys.date"
        }

        labels  = {
            label1 = "value1",
            label2 = "value2"
        } 
    } 
	  `, context)
}

func TestAccDialogflowCXIntent_defaultIntents(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Note: this isn't actually a "create" test; it creates resources in the TF state, but is actually importing the default objects GCP has created, then updating them.
				Config: testAccDialogflowCXIntent_defaultIntents_create(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dialogflow_cx_intent.default_negative_intent", "name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("google_dialogflow_cx_intent.default_welcome_intent", "name", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				ResourceName:            "google_dialogflow_cx_intent.default_negative_intent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				ResourceName:            "google_dialogflow_cx_intent.default_welcome_intent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				// This is testing updating the default objects without having to create them in the TF state first.
				Config: testAccDialogflowCXIntent_defaultIntents_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dialogflow_cx_intent.default_negative_intent", "name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("google_dialogflow_cx_intent.default_welcome_intent", "name", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				ResourceName:            "google_dialogflow_cx_intent.default_negative_intent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				ResourceName:            "google_dialogflow_cx_intent.default_welcome_intent",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccDialogflowCXIntent_defaultIntents_create(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
  display_name          = "tf-test-dialogflowcx-agent%{random_suffix}"
  location              = "global"
  default_language_code = "en"
  time_zone             = "America/New_York"
}

resource "google_dialogflow_cx_intent" "default_negative_intent" {
  parent                     = google_dialogflow_cx_agent.agent.id
  is_default_negative_intent = true
  display_name               = "Default Negative Intent"
  priority                   = 1
  is_fallback                = true
  training_phrases {
     parts {
         text = "Never match this phrase"
     }
     repeat_count = 1
  }
}

resource "google_dialogflow_cx_intent" "default_welcome_intent" {
  parent                    = google_dialogflow_cx_agent.agent.id
  is_default_welcome_intent = true
  display_name              = "Default Welcome Intent"
  priority                  = 1
  training_phrases {
     parts {
         text = "Hello"
     }
     repeat_count = 1
  }
}
`, context)
}

func testAccDialogflowCXIntent_defaultIntents_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
  display_name          = "tf-test-dialogflowcx-agent%{random_suffix}"
  location              = "global"
  default_language_code = "en"
  time_zone             = "America/New_York"
}

resource "google_dialogflow_cx_intent" "default_negative_intent" {
  parent                     = google_dialogflow_cx_agent.agent.id
  is_default_negative_intent = true
  display_name               = "Default Negative Intent"
  priority                   = 1
  is_fallback                = true
  training_phrases {
     parts {
         text = "An updated phrase to never match."
     }
     repeat_count = 2
  }
}

resource "google_dialogflow_cx_intent" "default_welcome_intent" {
  parent                    = google_dialogflow_cx_agent.agent.id
  is_default_welcome_intent = true
  display_name              = "Default Welcome Intent"
  priority                  = 1
  training_phrases {
     parts {
         text = "An updated hello."
     }
     repeat_count = 2
  }
}
`, context)
}
