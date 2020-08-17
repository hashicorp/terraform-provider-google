package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSqlClientCert_mysql(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlClientCertDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlClientCert_mysql(instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlClientCertExists(t, "google_sql_ssl_cert.cert1"),
					testAccCheckGoogleSqlClientCertExists(t, "google_sql_ssl_cert.cert2"),
				),
			},
		},
	})
}

func TestAccSqlClientCert_postgres(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlClientCertDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlClientCert_postgres(instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlClientCertExists(t, "google_sql_ssl_cert.cert"),
				),
			},
		},
	})
}

func testAccCheckGoogleSqlClientCertExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		instance := rs.Primary.Attributes["instance"]
		fingerprint := rs.Primary.Attributes["sha1_fingerprint"]
		sslClientCert, err := config.clientSqlAdmin.SslCerts.Get(config.Project, instance, fingerprint).Do()

		if err != nil {
			return err
		}

		if sslClientCert.Instance == instance && sslClientCert.Sha1Fingerprint == fingerprint {
			return nil
		}

		return fmt.Errorf("Not found: %s: %s", n, err)
	}
}

func testAccSqlClientCertDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			config := googleProviderConfig(t)
			if rs.Type != "google_sql_ssl_cert" {
				continue
			}

			fingerprint := rs.Primary.Attributes["sha1_fingerprint"]
			instance := rs.Primary.Attributes["instance"]
			sslCert, _ := config.clientSqlAdmin.SslCerts.Get(config.Project, instance, fingerprint).Do()

			commonName := rs.Primary.Attributes["common_name"]
			if sslCert != nil {
				return fmt.Errorf("Client cert %q still exists, should have been destroyed", commonName)
			}

			return nil
		}

		return nil
	}
}

func testGoogleSqlClientCert_mysql(instance string) string {
	return fmt.Sprintf(`
	resource "google_sql_database_instance" "instance" {
		name = "%s"
		region = "us-central1"
		settings {
			tier = "db-f1-micro"
		}
	}

	resource "google_sql_ssl_cert" "cert1" {
		common_name = "cert1"
		instance = "${google_sql_database_instance.instance.name}"
	}

	resource "google_sql_ssl_cert" "cert2" {
		common_name = "cert2"
		instance = "${google_sql_database_instance.instance.name}"
	}
	`, instance)
}

func testGoogleSqlClientCert_postgres(instance string) string {
	return fmt.Sprintf(`
	resource "google_sql_database_instance" "instance" {
		name = "%s"
		region = "us-central1"
		database_version = "POSTGRES_9_6"

		settings {
			tier = "db-f1-micro"
		}

		timeouts {
			create = "20m"
			delete = "20m"
		}
	}

	resource "google_sql_ssl_cert" "cert" {
		common_name = "cert"
		instance = "${google_sql_database_instance.instance.name}"
	}
	`, instance)
}
