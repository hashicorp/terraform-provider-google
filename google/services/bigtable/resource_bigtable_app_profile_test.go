// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigtableAppProfile_update(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableAppProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAppProfile_update1(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
			{
				Config: testAccBigtableAppProfile_update2(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
		},
	})
}

func TestAccBigtableAppProfile_ignoreWarnings(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableAppProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAppProfile_warningsProduced(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.gae-profile1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
		},
	})
}

func TestAccBigtableAppProfile_multiClusterIds(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableAppProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAppProfile_updateMC1(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
			{
				Config: testAccBigtableAppProfile_updateMC2(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
		},
	})
}

func TestAccBigtableAppProfile_updateSingleToMulti(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableAppProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAppProfile_update1(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
			{
				Config: testAccBigtableAppProfile_updateMC2(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
			{
				Config: testAccBigtableAppProfile_update1(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
		},
	})
}

func testAccBigtableAppProfile_update1(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s2"
    zone         = "us-central1-a"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s3"
    zone         = "us-central1-c"
    num_nodes    = 1
    storage_type = "HDD"
  }

  deletion_protection = false
}

resource "google_bigtable_app_profile" "ap" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "test"

  single_cluster_routing {
    cluster_id                 = %q
    allow_transactional_writes = true
  }

  ignore_warnings               = true
}
`, instanceName, instanceName, instanceName, instanceName, instanceName)
}

func testAccBigtableAppProfile_update2(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s2"
    zone         = "us-central1-a"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s3"
    zone         = "us-central1-c"
    num_nodes    = 1
    storage_type = "HDD"
  }

  deletion_protection = false
}

resource "google_bigtable_app_profile" "ap" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "test"
  description    = "add a description"

  single_cluster_routing {
    cluster_id                 = %q
    allow_transactional_writes = false
  }

  ignore_warnings               = true
}
`, instanceName, instanceName, instanceName, instanceName, instanceName)
}

func testAccBigtableAppProfile_warningsProduced(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  instance_type = "PRODUCTION"
  cluster {
    cluster_id   = "%s1"
    zone         = "us-central1-b"
    num_nodes    = 3
  }

  cluster {
    cluster_id   = "%s2"
    zone         = "us-west1-a"
    num_nodes    = 3
  }

  cluster {
    cluster_id   = "%s3"
    zone         = "us-west1-b"
    num_nodes    = 3
  }

  deletion_protection = false
}

resource "google_bigtable_app_profile" "gae-profile1" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "bigtableinstance-sample1"

  single_cluster_routing {
    cluster_id                 = "%s1"
    allow_transactional_writes = true
  }

  ignore_warnings               = true
}

resource "google_bigtable_app_profile" "gae-profile2" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "bigtableinstance-sample2"

  single_cluster_routing {
    cluster_id                 = "%s2"
    allow_transactional_writes = true
  }

  ignore_warnings               = true
}

resource "google_bigtable_app_profile" "gae-profile3" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "bigtableinstance-sample3"

  single_cluster_routing {
    cluster_id                 = "%s3"
    allow_transactional_writes = true
  }

  ignore_warnings               = true
}
`, instanceName, instanceName, instanceName, instanceName, instanceName, instanceName, instanceName)
}

func testAccBigtableAppProfile_updateMC1(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s2"
    zone         = "us-central1-a"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s3"
    zone         = "us-central1-c"
    num_nodes    = 1
    storage_type = "HDD"
  }

  deletion_protection = false
}

resource "google_bigtable_app_profile" "ap" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "test"

  multi_cluster_routing_use_any     = true
  multi_cluster_routing_cluster_ids = ["%s", "%s2", "%s3"]

  ignore_warnings               = true
}
`, instanceName, instanceName, instanceName, instanceName, instanceName, instanceName, instanceName)
}

func testAccBigtableAppProfile_updateMC2(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s2"
    zone         = "us-central1-a"
    num_nodes    = 1
    storage_type = "HDD"
  }

  cluster {
    cluster_id   = "%s3"
    zone         = "us-central1-c"
    num_nodes    = 1
    storage_type = "HDD"
  }

  deletion_protection = false
}

resource "google_bigtable_app_profile" "ap" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "test"
  description    = "add a description"

  multi_cluster_routing_use_any     = true
  multi_cluster_routing_cluster_ids = ["%s", "%s2"]

  ignore_warnings               = true
}
`, instanceName, instanceName, instanceName, instanceName, instanceName, instanceName)
}
