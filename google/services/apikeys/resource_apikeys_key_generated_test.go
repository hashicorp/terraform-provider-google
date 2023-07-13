// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package apikeys_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	apikeys "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/apikeys"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccApikeysKey_AndroidKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApikeysKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApikeysKey_AndroidKey(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApikeysKey_AndroidKeyUpdate0(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccApikeysKey_BasicKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApikeysKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApikeysKey_BasicKey(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApikeysKey_BasicKeyUpdate0(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccApikeysKey_IosKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApikeysKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApikeysKey_IosKey(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApikeysKey_IosKeyUpdate0(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccApikeysKey_MinimalKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApikeysKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApikeysKey_MinimalKey(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccApikeysKey_ServerKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  envvar.GetTestBillingAccountFromEnv(t),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckApikeysKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApikeysKey_ServerKey(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApikeysKey_ServerKeyUpdate0(context),
			},
			{
				ResourceName:      "google_apikeys_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApikeysKey_AndroidKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    android_key_restrictions {
      allowed_applications {
        package_name     = "com.example.app123"
        sha1_fingerprint = "1699466a142d4682a5f91b50fdf400f2358e2b0b"
      }
    }

    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_AndroidKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    android_key_restrictions {
      allowed_applications {
        package_name     = "com.example.app124"
        sha1_fingerprint = "1cf89aa28625da86a7e5a7550cf7fd33d611f6fd"
      }
    }

    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_BasicKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    browser_key_restrictions {
      allowed_referrers = [".*"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_BasicKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key-update"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "maps.googleapis.com"
      methods = ["POST*"]
    }

    browser_key_restrictions {
      allowed_referrers = [".*com"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_IosKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    ios_key_restrictions {
      allowed_bundle_ids = ["com.google.app.macos"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_IosKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    ios_key_restrictions {
      allowed_bundle_ids = ["com.google.alex.ios"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_MinimalKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_ServerKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    server_key_restrictions {
      allowed_ips = ["127.0.0.1"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccApikeysKey_ServerKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_apikeys_key" "primary" {
  name         = "tf-test-key%{random_suffix}"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    server_key_restrictions {
      allowed_ips = ["127.0.0.2", "192.168.1.1"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "tf-test-app%{random_suffix}"
  name       = "tf-test-app%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccCheckApikeysKeyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_apikeys_key" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &apikeys.Key{
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				DisplayName: dcl.String(rs.Primary.Attributes["display_name"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				KeyString:   dcl.StringOrNil(rs.Primary.Attributes["key_string"]),
				Uid:         dcl.StringOrNil(rs.Primary.Attributes["uid"]),
			}

			client := transport_tpg.NewDCLApikeysClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetKey(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_apikeys_key still exists %v", obj)
			}
		}
		return nil
	}
}
