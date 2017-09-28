package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

// Unit tests

func TestComputeAddressIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedCanonicalId string
		Config              *Config
	}{
		"id is a full self link": {
			ImportId:            "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/addresses/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/regions/us-central1/addresses/test-address",
		},
		"id is a partial self link": {
			ImportId:            "projects/test-project/regions/us-central1/addresses/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/regions/us-central1/addresses/test-address",
		},
		"id is project/region/address": {
			ImportId:            "test-project/us-central1/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/regions/us-central1/addresses/test-address",
		},
		"id is region/address": {
			ImportId:            "us-central1/test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/default-project/regions/us-central1/addresses/test-address",
			Config:              &Config{Project: "default-project"},
		},
		"id is address": {
			ImportId:            "test-address",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/default-project/regions/us-east1/addresses/test-address",
			Config:              &Config{Project: "default-project", Region: "us-east1"},
		},
		"id has invalid format": {
			ImportId:      "i/n/v/a/l/i/d",
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		addressId, err := parseComputeAddressId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if addressId.canonicalId() != tc.ExpectedCanonicalId {
			t.Fatalf("bad: %s, expected canonical id to be `%s` but is `%s`", tn, tc.ExpectedCanonicalId, addressId.canonicalId())
		}
	}
}

// Acceptance tests

func TestAccComputeAddress_basic(t *testing.T) {
	t.Parallel()

	var addr compute.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeAddress_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeAddressExists(
						"google_compute_address.foobar", &addr),
				),
			},
		},
	})
}

func TestAccComputeAddress_internal(t *testing.T) {
	var addr computeBeta.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeAddress_internal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBetaAddressExists(
						"google_compute_address.foobar", &addr),
					resource.TestCheckResourceAttr(
						"google_compute_address.foobar", "address_type", "INTERNAL"),
				),
			},
		},
	})
}

func testAccCheckComputeAddressDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_address" {
			continue
		}

		addressId, err := parseComputeAddressId(rs.Primary.ID, nil)

		_, err = config.clientCompute.Addresses.Get(
			config.Project, addressId.Region, addressId.Name).Do()
		if err == nil {
			return fmt.Errorf("Address still exists")
		}
	}

	return nil
}

func testAccCheckComputeAddressExists(n string, addr *compute.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		addressId, err := parseComputeAddressId(rs.Primary.ID, nil)

		found, err := config.clientCompute.Addresses.Get(
			config.Project, addressId.Region, addressId.Name).Do()
		if err != nil {
			return err
		}

		if found.Name != addressId.Name {
			return fmt.Errorf("Addr not found")
		}

		*addr = *found

		return nil
	}
}

func testAccCheckComputeBetaAddressExists(n string, addr *computeBeta.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		addressId, err := parseComputeAddressId(rs.Primary.ID, nil)

		found, err := config.clientComputeBeta.Addresses.Get(
			config.Project, addressId.Region, addressId.Name).Do()
		if err != nil {
			return err
		}

		if found.Name != addressId.Name {
			return fmt.Errorf("Addr not found")
		}

		*addr = *found

		return nil
	}
}

var testAccComputeAddress_basic = fmt.Sprintf(`
resource "google_compute_address" "foobar" {
	name = "address-test-%s"
}`, acctest.RandString(10))

var testAccComputeAddress_internal = fmt.Sprintf(`
resource "google_compute_network" "default" {
  name = "network-test-%s"
}

resource "google_compute_subnetwork" "foo" {
  name          = "subnetwork-test-%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-east1"
  network       = "${google_compute_network.default.self_link}"
}

resource "google_compute_address" "foobar" {
  name         = "address-test-%s"
  subnetwork   = "${google_compute_subnetwork.foo.self_link}"
  address_type = "INTERNAL"
  region       = "us-east1"
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
