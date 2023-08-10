// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleSQLCaCerts_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("data-ssl-ca-cert-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleSQLCaCertsConfig(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleSQLCaCertsCheck("data.google_sql_ca_certs.ca_certs", "google_sql_database_instance.foo"),
					testAccDataSourceGoogleSQLCaCertsCheck("data.google_sql_ca_certs.ca_certs_self_link", "google_sql_database_instance.foo"),
					resource.TestCheckResourceAttr("data.google_sql_ca_certs.ca_certs", "certs.#", "1"),
					resource.TestCheckResourceAttr("data.google_sql_ca_certs.ca_certs_self_link", "certs.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleSQLCaCertsCheck(datasourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", datasourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		datasourceAttributes := ds.Primary.Attributes
		resourceAttributes := rs.Primary.Attributes

		instanceToDatasourceAttrsMapping := map[string]string{
			"server_ca_cert.0.cert":             "certs.0.cert",
			"server_ca_cert.0.common_name":      "certs.0.common_name",
			"server_ca_cert.0.create_time":      "certs.0.create_time",
			"server_ca_cert.0.expiration_time":  "certs.0.expiration_time",
			"server_ca_cert.0.sha1_fingerprint": "certs.0.sha1_fingerprint",
		}

		for resourceAttr, datasourceAttr := range instanceToDatasourceAttrsMapping {
			if resourceAttributes[resourceAttr] != datasourceAttributes[datasourceAttr] {
				return fmt.Errorf(
					"%s is %s; want %s",
					datasourceAttr,
					datasourceAttributes[datasourceAttr],
					resourceAttributes[resourceAttr],
				)
			}
		}

		return nil
	}
}

func testAccDataSourceGoogleSQLCaCertsConfig(instanceName string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "foo" {
  name             = "%s"
  region           = "us-central1"
  database_version = "MYSQL_5_7"
  settings {
    tier                   = "db-f1-micro"
  }

  deletion_protection = false
}

data "google_sql_ca_certs" "ca_certs" {
  instance = google_sql_database_instance.foo.name
}

data "google_sql_ca_certs" "ca_certs_self_link" {
  instance = google_sql_database_instance.foo.self_link
}
`, instanceName)
}
