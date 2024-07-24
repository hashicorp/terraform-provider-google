// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccPubsubSchema_update(t *testing.T) {
	t.Parallel()

	schema := fmt.Sprintf("tf-test-schema-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSchema_basic(schema),
			},
			{
				ResourceName:      "google_pubsub_schema.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubSchema_updated(schema),
			},
			{
				ResourceName:      "google_pubsub_schema.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPubsubSchema_basic(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\n}"
	}
`, schema)
}

func testAccPubsubSchema_updated(schema string) string {
	return fmt.Sprintf(`
	resource "google_pubsub_schema" "foo" {
		name = "%s"
		type = "PROTOCOL_BUFFER"
		definition = "syntax = \"proto3\";\nmessage Results {\nstring message_request = 1;\nstring message_response = 2;\nstring timestamp_request = 3;\n}"
	}
`, schema)
}
