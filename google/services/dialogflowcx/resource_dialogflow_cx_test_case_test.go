// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dialogflowcx_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDialogflowCXTestCase_update(t *testing.T) {
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
				Config: testAccDialogflowCXTestCase_full(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_test_case.basic_test_case",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDialogflowCXTestCase_update(context),
			},
			{
				ResourceName:      "google_dialogflow_cx_test_case.basic_test_case",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDialogflowCXTestCase_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
  display_name               = "tf-test-dialogflowcx-agent%{random_suffix}"
  location                   = "global"
  default_language_code      = "en"
  supported_language_codes   = ["fr", "de", "es"]
  time_zone                  = "America/New_York"
  description                = "Example description."
  avatar_uri                 = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
  enable_stackdriver_logging = true
  enable_spell_correction    = true
  speech_to_text_settings {
    enable_speech_adaptation = true
  }
}

resource "google_dialogflow_cx_page" "page" {
  parent       = google_dialogflow_cx_agent.agent.start_flow
  display_name = "MyPage"

  transition_routes {
    intent = google_dialogflow_cx_intent.intent.id
    trigger_fulfillment {
      messages {
        text {
          text = ["Training phrase response"]
        }
      }
    }
  }

  event_handlers {
    event = "some-event"
    trigger_fulfillment {
      messages {
        text {
          text = ["Handling some event"]
        }
      }
    }
  }
}

resource "google_dialogflow_cx_intent" "intent" {
  parent       = google_dialogflow_cx_agent.agent.id
  display_name = "MyIntent"
  priority     = 1
  training_phrases {
    parts {
      text = "training phrase"
    }
    repeat_count = 1
  }
}

resource "google_dialogflow_cx_test_case" "basic_test_case" {
  parent       = google_dialogflow_cx_agent.agent.id
  display_name = "MyTestCase"

  test_config {
    tracking_parameters = []
    flow                = google_dialogflow_cx_agent.agent.start_flow
  }

  test_case_conversation_turns {
    user_input {
      input {
        language_code = "en"
        text {
          text = "some phrase"
        }
      }
    }
    virtual_agent_output {
      text_responses {
        text = ["Some response"]
      }
    }
  }
}`, context)
}

func testAccDialogflowCXTestCase_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_cx_agent" "agent" {
  display_name               = "tf-test-dialogflowcx-agent%{random_suffix}"
  location                   = "global"
  default_language_code      = "en"
  supported_language_codes   = ["fr", "de", "es"]
  time_zone                  = "America/New_York"
  description                = "Example description."
  avatar_uri                 = "https://storage.cloud.google.com/dialogflow-test-host-image/cloud-logo.png"
  enable_stackdriver_logging = true
  enable_spell_correction    = true
  speech_to_text_settings {
    enable_speech_adaptation = true
  }
}

resource "google_dialogflow_cx_page" "page" {
  parent       = google_dialogflow_cx_agent.agent.start_flow
  display_name = "MyPage"

  transition_routes {
    intent = google_dialogflow_cx_intent.intent.id
    trigger_fulfillment {
      messages {
        text {
          text = ["Training phrase response"]
        }
      }
    }
  }

  event_handlers {
    event = "some-event"
    trigger_fulfillment {
      messages {
        text {
          text = ["Handling some event"]
        }
      }
    }
  }
}

resource "google_dialogflow_cx_intent" "intent" {
  parent       = google_dialogflow_cx_agent.agent.id
  display_name = "MyIntent"
  priority     = 1
  training_phrases {
    parts {
      text = "training phrase"
    }
    repeat_count = 1
  }
}

resource "google_dialogflow_cx_test_case" "basic_test_case" {
  parent       = google_dialogflow_cx_agent.agent.id
  display_name = "MyTestCase"
  tags         = ["#tag1"]
  notes        = "demonstrates a simple training phrase response"

  test_config {
    tracking_parameters = ["some_param"]
    page                = google_dialogflow_cx_page.page.id
  }

  test_case_conversation_turns {
    user_input {
      input {
        language_code = "en"
        text {
          text = "training phrase"
        }
      }
      injected_parameters       = jsonencode({ some_param = "1" })
      is_webhook_enabled        = true
      enable_sentiment_analysis = true
    }
    virtual_agent_output {
      session_parameters = jsonencode({ some_param = "1" })
      triggered_intent {
        name = google_dialogflow_cx_intent.intent.id
      }
      current_page {
        name = google_dialogflow_cx_page.page.id
      }
      text_responses {
        text = ["Training phrase response"]
      }
    }
  }

  test_case_conversation_turns {
    user_input {
      input {
        event {
          event = "some-event"
        }
      }
    }
    virtual_agent_output {
      current_page {
        name = google_dialogflow_cx_page.page.id
      }
      text_responses {
        text = ["Handling some event"]
      }
    }
  }

  test_case_conversation_turns {
    user_input {
      input {
        dtmf {
          digits       = "12"
          finish_digit = "3"
        }
      }
    }
    virtual_agent_output {
      text_responses {
        text = ["I didn't get that. Can you say it again?"]
      }
    }
  }
}`, context)
}
