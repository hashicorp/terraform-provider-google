// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package filestore_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/filestore"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				ImportStateVerifyIgnore: []string{"zone", "location"},
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
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
			{
				Config: testAccFilestoreInstance_reservedIpRange_update2(name),
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
    reserved_ip_range = "172.19.30.0/29"
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
    reserved_ip_range = "172.19.31.0/29"
  }
}
`, name)
}
