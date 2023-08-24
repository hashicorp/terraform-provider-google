// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package workflows_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/workflows"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccWorkflowsWorkflow_Update(t *testing.T) {
	// Custom test written to test diffs
	t.Parallel()

	workflowName := fmt.Sprintf("tf-test-acc-workflow-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkflowsWorkflowDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowsWorkflow_Update(workflowName),
			},
			{
				Config: testAccWorkflowsWorkflow_Updated(workflowName),
			},
		},
	})
}

func testAccWorkflowsWorkflow_Update(name string) string {
	return fmt.Sprintf(`
resource "google_workflows_workflow" "example" {
  name          = "%s"
  region        = "us-central1"
  description   = "Magic"
  source_contents = <<-EOF
  # This is a sample workflow, feel free to replace it with your source code
  #
  # This workflow does the following:
  # - reads current time and date information from an external API and stores
  #   the response in CurrentDateTime variable
  # - retrieves a list of Wikipedia articles related to the day of the week
  #   from CurrentDateTime
  # - returns the list of articles as an output of the workflow
  # FYI, In terraform you need to escape the $$ or it will cause errors.

  - getCurrentTime:
      call: http.get
      args:
          url: https://us-central1-workflowsample.cloudfunctions.net/datetime
      result: CurrentDateTime
  - readWikipedia:
      call: http.get
      args:
          url: https://en.wikipedia.org/w/api.php
          query:
              action: opensearch
              search: $${CurrentDateTime.body.dayOfTheWeek}
      result: WikiResult
  - returnOutput:
      return: $${WikiResult.body[1]}
EOF
}
`, name)
}

func testAccWorkflowsWorkflow_Updated(name string) string {
	return fmt.Sprintf(`
resource "google_workflows_workflow" "example" {
  name          = "%s"
  region        = "us-central1"
  description   = "Magic"
  source_contents = <<-EOF
  # This is a sample workflow, feel free to replace it with your source code
  #
  # This workflow does the following:
  # - reads current time and date information from an external API and stores
  #   the response in CurrentDateTime variable
  # - retrieves a list of Wikipedia articles related to the day of the week
  #   from CurrentDateTime
  # - returns the list of articles as an output of the workflow
  # FYI, In terraform you need to escape the $$ or it will cause errors.

  - getCurrentTime:
      call: http.get
      args:
          url: https://us-central1-workflowsample.cloudfunctions.net/datetime
      result: CurrentDateTime
  - readWikipedia:
      call: http.get
      args:
          url: https:/fi.wikipedia.org/w/api.php
          query:
              action: opensearch
              search: $${CurrentDateTime.body.dayOfTheWeek}
      result: WikiResult
  - returnOutput:
      return: $${WikiResult.body[1]}
EOF
}
`, name)
}

func TestWorkflowsWorkflowStateUpgradeV0(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Attributes map[string]interface{}
		Expected   map[string]string
		Meta       interface{}
	}{
		"shorten long name": {
			Attributes: map[string]interface{}{
				"name": "projects/my-project/locations/us-central1/workflows/my-workflow",
			},
			Expected: map[string]string{
				"name": "my-workflow",
			},
			Meta: &transport_tpg.Config{},
		},
		"short name stays": {
			Attributes: map[string]interface{}{
				"name": "my-workflow",
			},
			Expected: map[string]string{
				"name": "my-workflow",
			},
			Meta: &transport_tpg.Config{},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := workflows.ResourceWorkflowsWorkflowUpgradeV0(context.Background(), tc.Attributes, tc.Meta)

			if err != nil {
				t.Error(err)
			}

			for _, expectedName := range tc.Expected {
				if actual["name"] != expectedName {
					t.Errorf("expected: name -> %#v\n got: name -> %#v\n in: %#v",
						expectedName, actual["name"], actual)
				}
			}
		})
	}
}

func TestAccWorkflowsWorkflow_CMEK(t *testing.T) {
	// Custom test written to test diffs
	t.Parallel()

	workflowName := fmt.Sprintf("tf-test-acc-workflow-%d", acctest.RandInt(t))
	kms := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	if acctest.BootstrapPSARole(t, "service-", "gcp-sa-workflows", "roles/cloudkms.cryptoKeyEncrypterDecrypter") {
		t.Fatal("Stopping the test because a role was added to the policy.")
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkflowsWorkflowDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowsWorkflow_CMEK(workflowName, kms.CryptoKey.Name),
			},
		},
	})
}

func testAccWorkflowsWorkflow_CMEK(workflowName, kmsKeyName string) string {
	return fmt.Sprintf(`
resource "google_workflows_workflow" "example" {
  name          = "%s"
  region        = "us-central1"
  description   = "Magic"
  crypto_key_name = "%s"
  source_contents = <<-EOF
  # This is a sample workflow, feel free to replace it with your source code
  #
  # This workflow does the following:
  # - reads current time and date information from an external API and stores
  #   the response in CurrentDateTime variable
  # - retrieves a list of Wikipedia articles related to the day of the week
  #   from CurrentDateTime
  # - returns the list of articles as an output of the workflow
  # FYI, In terraform you need to escape the $$ or it will cause errors.

  - getCurrentTime:
      call: http.get
      args:
          url: https://us-central1-workflowsample.cloudfunctions.net/datetime
      result: CurrentDateTime
  - readWikipedia:
      call: http.get
      args:
          url: https:/fi.wikipedia.org/w/api.php
          query:
              action: opensearch
              search: $${CurrentDateTime.body.dayOfTheWeek}
      result: WikiResult
  - returnOutput:
      return: $${WikiResult.body[1]}
EOF
}
`, workflowName, kmsKeyName)
}
