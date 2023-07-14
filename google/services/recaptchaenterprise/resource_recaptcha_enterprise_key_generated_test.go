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

package recaptchaenterprise_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	recaptchaenterprise "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/recaptchaenterprise"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccRecaptchaEnterpriseKey_AndroidKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRecaptchaEnterpriseKey_AndroidKey(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRecaptchaEnterpriseKey_AndroidKeyUpdate0(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccRecaptchaEnterpriseKey_IosKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRecaptchaEnterpriseKey_IosKey(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRecaptchaEnterpriseKey_IosKeyUpdate0(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccRecaptchaEnterpriseKey_MinimalKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRecaptchaEnterpriseKey_MinimalKey(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccRecaptchaEnterpriseKey_WebKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRecaptchaEnterpriseKey_WebKey(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRecaptchaEnterpriseKey_WebKeyUpdate0(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccRecaptchaEnterpriseKey_WebScoreKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRecaptchaEnterpriseKey_WebScoreKey(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRecaptchaEnterpriseKey_WebScoreKeyUpdate0(context),
			},
			{
				ResourceName:      "google_recaptcha_enterprise_key.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRecaptchaEnterpriseKey_AndroidKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-one"

  android_settings {
    allow_all_package_names = true
    allowed_package_names   = []
  }

  labels = {
    label-one = "value-one"
  }

  project = "%{project_name}"

  testing_options {
    testing_score = 0.8
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_AndroidKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-two"

  android_settings {
    allow_all_package_names = false
    allowed_package_names   = ["com.android.application"]
  }

  labels = {
    label-two = "value-two"
  }

  project = "%{project_name}"

  testing_options {
    testing_score = 0.8
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_IosKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-one"

  ios_settings {
    allow_all_bundle_ids = true
    allowed_bundle_ids   = []
  }

  labels = {
    label-one = "value-one"
  }

  project = "%{project_name}"

  testing_options {
    testing_score = 1
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_IosKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-two"

  ios_settings {
    allow_all_bundle_ids = false
    allowed_bundle_ids   = ["com.companyname.appname"]
  }

  labels = {
    label-two = "value-two"
  }

  project = "%{project_name}"

  testing_options {
    testing_score = 1
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_MinimalKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-one"
  labels       = {}
  project      = "%{project_name}"

  web_settings {
    integration_type  = "SCORE"
    allow_all_domains = true
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_WebKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-one"

  labels = {
    label-one = "value-one"
  }

  project = "%{project_name}"

  testing_options {
    testing_challenge = "NOCAPTCHA"
    testing_score     = 0.5
  }

  web_settings {
    integration_type              = "CHECKBOX"
    allow_all_domains             = true
    allowed_domains               = []
    challenge_security_preference = "USABILITY"
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_WebKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-two"

  labels = {
    label-two = "value-two"
  }

  project = "%{project_name}"

  testing_options {
    testing_challenge = "NOCAPTCHA"
    testing_score     = 0.5
  }

  web_settings {
    integration_type              = "CHECKBOX"
    allow_all_domains             = false
    allowed_domains               = ["subdomain.example.com"]
    challenge_security_preference = "SECURITY"
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_WebScoreKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-one"

  labels = {
    label-one = "value-one"
  }

  project = "%{project_name}"

  testing_options {
    testing_score = 0.5
  }

  web_settings {
    integration_type  = "SCORE"
    allow_all_domains = true
    allow_amp_traffic = false
    allowed_domains   = []
  }
}


`, context)
}

func testAccRecaptchaEnterpriseKey_WebScoreKeyUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name-two"

  labels = {
    label-two = "value-two"
  }

  project = "%{project_name}"

  testing_options {
    testing_score = 0.5
  }

  web_settings {
    integration_type  = "SCORE"
    allow_all_domains = false
    allow_amp_traffic = true
    allowed_domains   = ["subdomain.example.com"]
  }
}


`, context)
}

func testAccCheckRecaptchaEnterpriseKeyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_recaptcha_enterprise_key" {
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

			obj := &recaptchaenterprise.Key{
				DisplayName: dcl.String(rs.Primary.Attributes["display_name"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Name:        dcl.StringOrNil(rs.Primary.Attributes["name"]),
			}

			client := transport_tpg.NewDCLRecaptchaEnterpriseClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetKey(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_recaptcha_enterprise_key still exists %v", obj)
			}
		}
		return nil
	}
}
