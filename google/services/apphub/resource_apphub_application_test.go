// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apphub_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccApphubApplication_applicationUpdateFull(t *testing.T) {
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
				Config: testAccApphubApplication_apphubApplicationFullExample(context),
			},
			{
				ResourceName:            "google_apphub_application.example2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id"},
			},
			{
				Config: testAccApphubApplication_applicationUpdateDisplayName(context),
			},
			{
				ResourceName:            "google_apphub_application.example2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id"},
			},
			{
				Config: testAccApphubApplication_applicationUpdateEnvironment(context),
			},
			{
				ResourceName:            "google_apphub_application.example2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id"},
			},
			{
				Config: testAccApphubApplication_applicationUpdateCriticality(context),
			},
			{
				ResourceName:            "google_apphub_application.example2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id"},
			},
			{
				Config: testAccApphubApplication_applicationUpdateOwners(context),
			},
			{
				ResourceName:            "google_apphub_application.example2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "application_id"},
			},
		},
	})
}

func testAccApphubApplication_applicationUpdateDisplayName(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_apphub_application" "example2" {
  location = "us-east1"
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

func testAccApphubApplication_applicationUpdateEnvironment(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_apphub_application" "example2" {
  location = "us-east1"
  application_id = "tf-test-example-application%{random_suffix}"
  display_name = "Application Full New%{random_suffix}"
  scope {
    type = "REGIONAL"
  }
  attributes {
    environment {
      type = "TEST"
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

func testAccApphubApplication_applicationUpdateCriticality(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_apphub_application" "example2" {
  location = "us-east1"
  application_id = "tf-test-example-application%{random_suffix}"
  display_name = "Application Full New%{random_suffix}"
  scope {
    type = "REGIONAL"
  }
  attributes {
    environment {
      type = "TEST"
		}
		criticality {  
      type = "MEDIUM"
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

func testAccApphubApplication_applicationUpdateOwners(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_apphub_application" "example2" {
  location = "us-east1"
  application_id = "tf-test-example-application%{random_suffix}"
  display_name = "Application Full New%{random_suffix}"
  scope {
    type = "REGIONAL"
  }
  attributes {
    environment {
      type = "TEST"
		}
		criticality {  
      type = "MEDIUM"
		}
		business_owners {
		  display_name =  "Alice%{random_suffix}"
		  email        =  "alice@google.com%{random_suffix}"
		}
		developer_owners {
		  display_name =  "Bob%{random_suffix}"
		  email        =  "bob@google.com%{random_suffix}"
		}
		developer_owners {
			display_name =  "Derek%{random_suffix}"
			email        =  "derek@google.com%{random_suffix}"
		}
		operator_owners {
		  display_name =  "Charlie%{random_suffix}"
		  email        =  "charlie@google.com%{random_suffix}"
		}
  }
}
`, context)
}

func TestAccApphubApplication_invalidConfigFails(t *testing.T) {
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
				Config:      testAccApphubApplication_applicationInvalidConfig1(context),
				ExpectError: regexp.MustCompile("Error validating location global with REGIONAL scope type"),
			},
			{
				Config:      testAccApphubApplication_applicationInvalidConfig2(context),
				ExpectError: regexp.MustCompile("Error validating location us-east1 with GLOBAL scope type"),
			},
		},
	})
}

func testAccApphubApplication_applicationInvalidConfig1(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_apphub_application" "invalid_example" {
  location = "global"
  application_id = "tf-test-invalid-example-application%{random_suffix}"
  scope {
    type = "REGIONAL"
  }
}
`, context)
}

func testAccApphubApplication_applicationInvalidConfig2(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_apphub_application" "invalid_example" {
  location = "us-east1"
  application_id = "tf-test-invalid-example-application%{random_suffix}"
  scope {
    type = "GLOBAL"
  }
}
`, context)
}
