// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVertexAIFeaturestoreEntitytypeFeature_vertexAiFeaturestoreEntitytypeFeatureExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIFeaturestoreEntitytypeFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIFeaturestoreEntitytypeFeature_vertexAiFeaturestoreEntitytypeFeatureExample(context),
			},
			{
				ResourceName:            "google_vertex_ai_featurestore_entitytype_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "entitytype"},
			},
		},
	})
}

func testAccVertexAIFeaturestoreEntitytypeFeature_vertexAiFeaturestoreEntitytypeFeatureExample(context map[string]interface{}) string {
	return tpgresource.Nprintf(`
resource "google_vertex_ai_featurestore" "featurestore" {
  name     = "terraform%{random_suffix}"
  labels = {
    foo = "bar"
  }
  region   = "us-central1"
  online_serving_config {
    fixed_node_count = 2
  }
}

resource "google_vertex_ai_featurestore_entitytype" "entity" {
  name     = "terraform%{random_suffix}"
  labels = {
    foo = "bar"
  }
  featurestore = google_vertex_ai_featurestore.featurestore.id
}

resource "google_vertex_ai_featurestore_entitytype_feature" "feature" {
  name     = "terraform%{random_suffix}"
  labels = {
    foo = "bar"
  }
  entitytype = google_vertex_ai_featurestore_entitytype.entity.id

  value_type = "INT64_ARRAY"
}
`, context)
}

func testAccCheckVertexAIFeaturestoreEntitytypeFeatureDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_featurestore_entitytype_feature" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VertexAIBasePath}}{{entitytype}}/features/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("VertexAIFeaturestoreEntitytypeFeature still exists at %s", url)
			}
		}

		return nil
	}
}
