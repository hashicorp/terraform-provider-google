// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDatastoreIndex_datastoreIndexExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckDatastoreIndexDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastoreIndex_datastoreIndexExample(context),
			},
			{
				ResourceName:      "google_datastore_index.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDatastoreIndex_datastoreIndexExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_datastore_index" "default" {
  kind = "foo"
  properties {
    name = "tf_test_property_a%{random_suffix}"
    direction = "ASCENDING"
  }
  properties {
    name = "tf_test_property_b%{random_suffix}"
    direction = "ASCENDING"
  }
}
`, context)
}

func testAccCheckDatastoreIndexDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_datastore_index" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{DatastoreBasePath}}projects/{{project}}/indexes/{{index_id}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, nil, datastoreIndex409Contention)
			if err == nil {
				return fmt.Errorf("DatastoreIndex still exists at %s", url)
			}
		}

		return nil
	}
}
