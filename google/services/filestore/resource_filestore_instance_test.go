// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package filestore_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/filestore"
)

func testResourceFilestoreInstanceStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"zone": "us-central1-a",
	}
}

func testResourceFilestoreInstanceStateDataV1() map[string]interface{} {
	v0 := testResourceFilestoreInstanceStateDataV0()
	return map[string]interface{}{
		"location": v0["zone"],
		"zone":     v0["zone"],
	}
}

func TestFilestoreInstanceStateUpgradeV0(t *testing.T) {
	expected := testResourceFilestoreInstanceStateDataV1()
	// linter complains about nil context even in a test setting
	actual, err := filestore.ResourceFilestoreInstanceUpgradeV0(context.Background(), testResourceFilestoreInstanceStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

func TestAccFilestoreInstance_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_update(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccFilestoreInstance_update2(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
		},
	})
}

func testAccFilestoreInstance_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "tf-instance-%s"
  zone        = "us-central1-b"
  tier        = "BASIC_HDD"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }

  labels = {
    baz = "qux"
  }
}
`, name)
}

func testAccFilestoreInstance_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "tf-instance-%s"
  zone        = "us-central1-b"
  tier        = "BASIC_HDD"
  description = "A modified instance created during testing."

  file_shares {
    capacity_gb = 1536
    name        = "share"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}
`, name)
}

func TestAccFilestoreInstance_reservedIpRange_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_reservedIpRange_update(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location", "networks.0.reserved_ip_range"},
			},
			{
				Config: testAccFilestoreInstance_reservedIpRange_update2(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location", "networks.0.reserved_ip_range"},
			},
		},
	})
}

func testAccFilestoreInstance_reservedIpRange_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.31.0/29"
  }
}
`, name)
}

func testAccFilestoreInstance_reservedIpRange_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.31.8/29"
  }
}
`, name)
}

func TestAccFilestoreInstance_deletionProtection_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	location := "us-central1-a"
	tier := "ZONAL"

	deletionProtection := true

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_deletionProtection_create(name, location, tier),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccFilestoreInstance_deletionProtection_update(name, location, tier, deletionProtection),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccFilestoreInstance_deletionProtection_update(name, location, tier, !deletionProtection),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

func testAccFilestoreInstance_deletionProtection_create(name, location, tier string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "%s"
  zone        = "%s"
  tier        = "%s"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
		modes   = ["MODE_IPV4"]
  }
}
`, name, location, tier)
}

func testAccFilestoreInstance_deletionProtection_update(name, location, tier string, deletionProtection bool) string {
	deletionProtectionReason := ""
	if deletionProtection {
		deletionProtectionReason = "A reason for deletion protection"
	}

	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "%s"
  zone        = "%s"
  tier        = "%s"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
		modes   = ["MODE_IPV4"]
  }

	deletion_protection_enabled = %t
	deletion_protection_reason = "%s"
}
`, name, location, tier, deletionProtection, deletionProtectionReason)
}

func TestAccFilestoreInstance_performanceConfig(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	location := "us-central1"
	tier := "REGIONAL"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_fixedIopsPerformanceConfig(name, location, tier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_filestore_instance.instance", "performance_config.0.fixed_iops.0.max_iops", "17000"),
				),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccFilestoreInstance_iopsPerTbPerformanceConfig(name, location, tier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_filestore_instance.instance", "performance_config.0.iops_per_tb.0.max_iops_per_tb", "17000"),
				),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccFilestoreInstance_defaultConfig(name, location, tier),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

func testAccFilestoreInstance_fixedIopsPerformanceConfig(name, location, tier string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "%s"
  location    = "%s"
  tier        = "%s"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
		modes   = ["MODE_IPV4"]
  }

  performance_config {
    fixed_iops {
      max_iops = 17000
    }
  }
}
`, name, location, tier)
}

func testAccFilestoreInstance_iopsPerTbPerformanceConfig(name, location, tier string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "%s"
  zone        = "%s"
  tier        = "%s"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
		modes   = ["MODE_IPV4"]
  }

  performance_config {
    iops_per_tb {
      max_iops_per_tb = 17000
    }
  }
}
`, name, location, tier)
}

func testAccFilestoreInstance_defaultConfig(name, location, tier string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "%s"
  zone        = "%s"
  tier        = "%s"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share"
  }

  networks {
    network = "default"
		modes   = ["MODE_IPV4"]
  }
}
`, name, location, tier)
}

func TestAccFilestoreInstance_tags(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	org := envvar.GetTestOrgFromEnv(t)
	tagKey := acctest.BootstrapSharedTestTagKey(t, "filestore-instances-tagkey")
	tagValue := acctest.BootstrapSharedTestTagValue(t, "filestore-instances-tagvalue", tagKey)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFileInstanceTags(name, map[string]string{org + "/" + tagKey: tagValue}),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location", "networks.0.reserved_ip_range", "tags"},
			},
		},
	})
}

func testAccFileInstanceTags(name string, tags map[string]string) string {
	r := fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-test-instance-%s"
  zone = "us-central1-b"
  tier = "BASIC_HDD"
  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }
  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.31.8/29"
  }
tags = {`, name)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}
