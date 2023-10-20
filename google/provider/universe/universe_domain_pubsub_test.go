// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package universe_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccUniverseDomainPubSub(t *testing.T) {
	// Skip this test in all env since this can only run in specific test project.
	t.Skip()

	universeDomain := envvar.GetTestUniverseDomainFromEnv(t)
	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(t, 10))
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUniverseDomain_basic_pubsub(universeDomain, topic, subscription),
			},
		},
	})
}

func testAccUniverseDomain_basic_pubsub(universeDomain, topic, subscription string) string {
	return fmt.Sprintf(`
provider "google" {
  universe_domain = "%s"
}
	  
resource "google_pubsub_topic" "foo" {
  name = "%s"
}
  
resource "google_pubsub_subscription" "foo" {
  name  = "%s"
  topic = google_pubsub_topic.foo.id
  
  message_retention_duration = "1200s"
  retain_acked_messages      = true
  ack_deadline_seconds       = 20
  expiration_policy {
	ttl = ""
  }
  enable_message_ordering    = false
}
`, universeDomain, topic, subscription)
}

func testAccCheckPubsubSubscriptionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_pubsub_subscription" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{PubsubBasePath}}projects/{{project}}/subscriptions/{{name}}")
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
				return fmt.Errorf("PubsubSubscription still exists at %s", url)
			}
		}

		return nil
	}
}
