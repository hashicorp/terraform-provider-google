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

package firebaserules_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	firebaserules "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/firebaserules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccFirebaserulesRelease_FirestoreReleaseHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFirebaserulesReleaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaserulesRelease_FirestoreReleaseHandWritten(context),
			},
			{
				ResourceName:      "google_firebaserules_release.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirebaserulesRelease_FirestoreReleaseHandWrittenUpdate0(context),
			},
			{
				ResourceName:      "google_firebaserules_release.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirebaserulesRelease_FirestoreReleaseHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firebaserules_release" "primary" {
  name         = "cloud.firestore"
  ruleset_name = "projects/%{project_name}/rulesets/${google_firebaserules_ruleset.firestore.name}"
  project      = "%{project_name}"

  lifecycle {
    replace_triggered_by = [
      google_firebaserules_ruleset.firestore
    ]
  }
}

resource "google_firebaserules_ruleset" "firestore" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name    = "firestore.rules"
    }
  }

  project = "%{project_name}"
}

`, context)
}

func testAccFirebaserulesRelease_FirestoreReleaseHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_firebaserules_release" "primary" {
  name         = "cloud.firestore"
  ruleset_name = "projects/%{project_name}/rulesets/${google_firebaserules_ruleset.firestore.name}"
  project      = "%{project_name}"

  lifecycle {
    replace_triggered_by = [
      google_firebaserules_ruleset.firestore
    ]
  }
}

resource "google_firebaserules_ruleset" "firestore" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if request.auth != null; } } }"
      name    = "firestore.rules"
    }
  }

  project = "%{project_name}"
}

`, context)
}

func testAccCheckFirebaserulesReleaseDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_firebaserules_release" {
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

			obj := &firebaserules.Release{
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				RulesetName: dcl.String(rs.Primary.Attributes["ruleset_name"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Disabled:    dcl.Bool(rs.Primary.Attributes["disabled"] == "true"),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLFirebaserulesClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetRelease(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_firebaserules_release still exists %v", obj)
			}
		}
		return nil
	}
}
