// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apphub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestDataSourceApphubApplication_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApphubApplicationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceApphubApplication_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_apphub_application.example_data", "google_apphub_application.example"),
				),
			},
		},
	})
}

func testDataSourceApphubApplication_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_apphub_application" "example_data" {
	project = google_apphub_application.example.project
	application_id = google_apphub_application.example.application_id
	location = google_apphub_application.example.location
}

resource "google_apphub_application" "example" {
  location = "us-central1"
  application_id = "tf-test-example-application%{random_suffix}"
  display_name = "Application Full New%{random_suffix}"
  scope {
    type = "REGIONAL"
  }
  attributes {
    environment {
      type = "STAGING"
	  }
    criticality {  
      type = "MISSION_CRITICAL"
    }
    business_owners {
      display_name =  "Alice%{random_suffix}"
      email        =  "alice@google.com%{random_suffix}"
    }
    developer_owners {
      display_name =  "Bob%{random_suffix}"
      email        =  "bob@google.com%{random_suffix}"
    }
    operator_owners {
      display_name =  "Charlie%{random_suffix}"
      email        =  "charlie@google.com%{random_suffix}"
    }
  }
}
`, context)
}
