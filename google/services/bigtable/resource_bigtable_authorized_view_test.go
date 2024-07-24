// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBigtableAuthorizedView_basic(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	authorizedViewName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	familyName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigtableAuthorizedViewInvalidDeletionProtection(instanceName, tableName, authorizedViewName, familyName),
				ExpectError: regexp.MustCompile(".*expected deletion_protection to be one of.*"),
			},
			{
				Config:      testAccBigtableAuthorizedViewInvalidSubsetView(instanceName, tableName, authorizedViewName, familyName),
				ExpectError: regexp.MustCompile(".*subset_view must be specified for authorized view.*"),
			},
			{
				Config:      testAccBigtableAuthorizedViewInvalidEncoding(instanceName, tableName, authorizedViewName, familyName),
				ExpectError: regexp.MustCompile(".*illegal base64 data.*"),
			},
			{
				Config: testAccBigtableAuthorizedViewBasic(instanceName, tableName, authorizedViewName, familyName),
			},
			{
				ResourceName:      "google_bigtable_authorized_view.authorized_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableAuthorizedViewWithRowPrefixesOnly(instanceName, tableName, authorizedViewName, familyName),
			},
			{
				ResourceName:      "google_bigtable_authorized_view.authorized_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableAuthorizedView_update(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	authorizedViewName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	familyName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAuthorizedViewWithQualifiersOnly(instanceName, tableName, authorizedViewName, familyName),
			},
			{
				ResourceName:      "google_bigtable_authorized_view.authorized_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableAuthorizedViewWithFamilySubsetsOnly(instanceName, tableName, authorizedViewName, familyName),
			},
			{
				ResourceName:            "google_bigtable_authorized_view.authorized_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"subset_view.0.family_subsets"}, // The order of the two family subsets is indeterministic.
			},
		},
	})
}

func TestAccBigtableAuthorizedView_destroy(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	authorizedViewName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	familyName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAuthorizedViewWithQualifierPrefixesOnly(instanceName, tableName, authorizedViewName, familyName),
			},
			{
				ResourceName:      "google_bigtable_authorized_view.authorized_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableAuthorizedViewDestroy(instanceName, tableName, familyName),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigtableAuthorizedViewDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_authorized_view" {
				continue
			}

			config := acctest.GoogleProviderConfig(t)
			c, err := config.BigTableClientFactory(config.UserAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
			if err != nil {
				// The instance is already gone
				return nil
			}

			_, err = c.AuthorizedViewInfo(ctx, rs.Primary.Attributes["table_name"], rs.Primary.Attributes["name"])
			if err == nil {
				return fmt.Errorf("AuthorizedView still present. Found %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["table_name"])
			}

			c.Close()
		}

		return nil
	}
}

func testAccBigtableAuthorizedViewInvalidDeletionProtection(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}
resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "random"
  subset_view {}
}
`, instanceName, instanceName, tableName, familyName, authorizedViewName)
}

func testAccBigtableAuthorizedViewInvalidSubsetView(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}
resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "UNPROTECTED"
}
`, instanceName, instanceName, tableName, familyName, authorizedViewName)
}

func testAccBigtableAuthorizedViewInvalidEncoding(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}
resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "UNPROTECTED"
  subset_view {
	  row_prefixes = ["#"]
  }
}
`, instanceName, instanceName, tableName, familyName, authorizedViewName)
}

func testAccBigtableAuthorizedViewBasic(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}

resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "UNPROTECTED"
  subset_view {}
}
`, instanceName, instanceName, tableName, familyName, authorizedViewName)
}

func testAccBigtableAuthorizedViewWithRowPrefixesOnly(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}

resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "UNPROTECTED"

  subset_view {
    row_prefixes = [base64encode("row1#"), base64encode("row2#")]
  }
}
`, instanceName, instanceName, tableName, familyName, authorizedViewName)
}

func testAccBigtableAuthorizedViewWithFamilySubsetsOnly(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
  column_family {
	  family = "%s-second"
  }
}

resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "UNPROTECTED"

  subset_view {
    family_subsets {
      family_name = "%s"
      qualifiers = [base64encode("qualifier"), base64encode("qualifier-second")]
    }
	family_subsets {
	  family_name = "%s-second"
	  qualifier_prefixes = [""]
	}
  }
}
`, instanceName, instanceName, tableName, familyName, familyName, authorizedViewName, familyName, familyName)
}

func testAccBigtableAuthorizedViewWithQualifiersOnly(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}

resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "UNPROTECTED"

  subset_view {
    family_subsets {
      family_name = "%s"
      qualifiers = [base64encode("qualifier")]
    }
  }
}
`, instanceName, instanceName, tableName, familyName, authorizedViewName, familyName)
}

func testAccBigtableAuthorizedViewWithQualifierPrefixesOnly(instanceName, tableName, authorizedViewName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}

resource "google_bigtable_authorized_view" "authorized_view" {
  name         = "%s"
  instance_name = google_bigtable_instance.instance.id
  table_name = google_bigtable_table.table.name
  deletion_protection = "UNPROTECTED"

  subset_view {
	  family_subsets {
	    family_name = "%s"
	    qualifier_prefixes = [""]
	  }
  }
}
`, instanceName, instanceName, tableName, familyName, authorizedViewName, familyName)
}

func testAccBigtableAuthorizedViewDestroy(instanceName, tableName, familyName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
	num_nodes  = 1
  }
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  column_family {
	  family = "%s"
  }
}
`, instanceName, instanceName, tableName, familyName)
}
