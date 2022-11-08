// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVertexAIFeaturestoreEntitytype_vertexAiFeaturestoreEntitytypeExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
		"kms_key_name":    BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
		"random_suffix":   randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVertexAIFeaturestoreEntitytypeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIFeaturestoreEntitytype_vertexAiFeaturestoreEntitytypeExample(context),
			},
			{
				ResourceName:            "google_vertex_ai_featurestore_entitytype.entity",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "featurestore"},
			},
		},
	})
}

func testAccVertexAIFeaturestoreEntitytype_vertexAiFeaturestoreEntitytypeExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_vertex_ai_featurestore" "featurestore" {
  name     = "terraform%{random_suffix}"
  labels = {
    foo = "bar"
  }
  region   = "us-central1"
  online_serving_config {
    fixed_node_count = 2
  }
  encryption_spec {
    kms_key_name = "%{kms_key_name}"
  }
}

resource "google_vertex_ai_featurestore_entitytype" "entity" {
  name     = "terraform%{random_suffix}"
  labels = {
    foo = "bar"
  }
  featurestore = google_vertex_ai_featurestore.featurestore.id
  monitoring_config {
    snapshot_analysis {
      disabled = false
      monitoring_interval_days = 1
      staleness_days = 21
    }
    numerical_threshold_config {
      value = 0.8
    }
    categorical_threshold_config {
      value = 10.0
    }
    import_features_analysis {
      state = "ENABLED"
      anomaly_detection_baseline = "PREVIOUS_IMPORT_FEATURES_STATS"
    }
  }
}
`, context)
}

func testAccCheckVertexAIFeaturestoreEntitytypeDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_featurestore_entitytype" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{VertexAIBasePath}}{{featurestore}}/entityTypes/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("VertexAIFeaturestoreEntitytype still exists at %s", url)
			}
		}

		return nil
	}
}
