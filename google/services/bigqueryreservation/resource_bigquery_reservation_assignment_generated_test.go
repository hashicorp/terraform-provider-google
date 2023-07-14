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

package bigqueryreservation_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	bigqueryreservation "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/bigqueryreservation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccBigqueryReservationAssignment_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryReservationAssignmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryReservationAssignment_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_bigquery_reservation_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reservation"},
			},
		},
	})
}

func testAccBigqueryReservationAssignment_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "basic" {
  name  = "tf-test-my-reservation%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"
  slot_capacity = 0
  ignore_idle_slots = false
}

resource "google_bigquery_reservation_assignment" "primary" {
  assignee  = "projects/%{project_name}"
  job_type = "PIPELINE"
  reservation = google_bigquery_reservation.basic.id
}
`, context)
}

func testAccCheckBigqueryReservationAssignmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_bigquery_reservation_assignment" {
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

			obj := &bigqueryreservation.Assignment{
				Assignee:    dcl.String(rs.Primary.Attributes["assignee"]),
				JobType:     bigqueryreservation.AssignmentJobTypeEnumRef(rs.Primary.Attributes["job_type"]),
				Reservation: dcl.String(rs.Primary.Attributes["reservation"]),
				Location:    dcl.StringOrNil(rs.Primary.Attributes["location"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				Name:        dcl.StringOrNil(rs.Primary.Attributes["name"]),
				State:       bigqueryreservation.AssignmentStateEnumRef(rs.Primary.Attributes["state"]),
			}

			client := transport_tpg.NewDCLBigqueryReservationClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetAssignment(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_bigquery_reservation_assignment still exists %v", obj)
			}
		}
		return nil
	}
}
