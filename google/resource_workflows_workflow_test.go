package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWorkflowsWorkflow_Update(t *testing.T) {
	// Custom test written to test diffs
	t.Parallel()

	workflowName := fmt.Sprintf("tf-test-acc-workflow-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWorkflowsWorkflowDestroyProducer(t),
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
			Meta: &Config{},
		},
		"short name stays": {
			Attributes: map[string]interface{}{
				"name": "my-workflow",
			},
			Expected: map[string]string{
				"name": "my-workflow",
			},
			Meta: &Config{},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := resourceWorkflowsWorkflowUpgradeV0(context.Background(), tc.Attributes, tc.Meta)

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
