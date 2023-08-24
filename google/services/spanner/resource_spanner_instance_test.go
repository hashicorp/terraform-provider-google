// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Acceptance Tests

func TestAccSpannerInstance_basic(t *testing.T) {
	t.Parallel()

	idName := fmt.Sprintf("spanner-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basic(idName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_noNodeCountSpecified(t *testing.T) {
	t.Parallel()

	idName := fmt.Sprintf("spanner-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSpannerInstance_noNodeCountSpecified(idName),
				ExpectError: regexp.MustCompile("one of `num_nodes,processing_units` must be specified"),
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutogenName(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	displayName := fmt.Sprintf("spanner-test-%s-dname", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutogenName(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "name"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_update(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	dName1 := fmt.Sprintf("spanner-dname1-%s", acctest.RandString(t, 10))
	dName2 := fmt.Sprintf("spanner-dname2-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_update(dName1, 1, false),
			},
			{
				ResourceName:      "google_spanner_instance.updater",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstance_update(dName2, 2, true),
			},
			{
				ResourceName:      "google_spanner_instance.updater",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_virtualUpdate(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	dName := fmt.Sprintf("spanner-dname1-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_virtualUpdate(dName, "true"),
			},
			{
				ResourceName: "google_spanner_instance.basic",
				ImportState:  true,
			},
			{
				Config: testAccSpannerInstance_virtualUpdate(dName, "false"),
			},
			{
				ResourceName: "google_spanner_instance.basic",
				ImportState:  true,
			},
		},
	})
}

func testAccSpannerInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-dname"
  num_nodes    = 1
}
`, name, name)
}

func testAccSpannerInstance_noNodeCountSpecified(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-dname"
}
`, name, name)
}

func testAccSpannerInstance_basicWithAutogenName(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = 1
}
`, name)
}

func testAccSpannerInstance_update(name string, nodes int, addLabel bool) string {
	extraLabel := ""
	if addLabel {
		extraLabel = "\"key2\" = \"value2\""
	}
	return fmt.Sprintf(`
resource "google_spanner_instance" "updater" {
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = %d

  labels = {
    "key1" = "value1"
    %s
  }
}
`, name, nodes, extraLabel)
}

func testAccSpannerInstance_virtualUpdate(name, virtual string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  processing_units = 100
  force_destroy    = "%s"
}
`, name, name, virtual)
}
