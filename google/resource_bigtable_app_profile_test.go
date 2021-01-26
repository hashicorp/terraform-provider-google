package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigtableAppProfile_update(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableAppProfileDestroyProducer(t),
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
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableAppProfileDestroyProducer(t),
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

func testAccBigtableAppProfile_update1(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 3
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
`, instanceName, instanceName, instanceName)
}

func testAccBigtableAppProfile_update2(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 3
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
`, instanceName, instanceName, instanceName)
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
