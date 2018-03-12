package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceComputeBackendService_basic(t *testing.T) {
	t.Parallel()

	serviceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	checkName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceComputeBackendService_basic(serviceName, checkName),
				Check:  testAccDataSourceComputeBackendServiceCheck("data.google_compute_backend_service.baz", "google_compute_backend_service.foobar"),
			},
		},
	})
}

func testAccDataSourceComputeBackendServiceCheck(dsName, rsName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return fmt.Errorf("can't find resource called %s in state", rsName)
		}

		ds, ok := s.RootModule().Resources[dsName]
		if !ok {
			return fmt.Errorf("can't find data source called %s in state", dsName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		attrsToTest := []string{
			"id",
			"name",
			"description",
			"self_link",
			"fingerprint",
			"port_name",
			"protocol",
		}

		for _, attrToTest := range attrsToTest {
			if dsAttr[attrToTest] != rsAttr[attrToTest] {
				return fmt.Errorf("%s is %s; want %s", attrToTest, dsAttr[attrToTest], rsAttr[attrToTest])
			}
		}

		return nil
	}
}

func testAccDataSourceComputeBackendService_basic(serviceName, checkName string) string {
	return fmt.Sprintf(`
resource "google_compute_backend_service" "foobar" {
  name          = "%s"
  description   = "foobar backend service"
  health_checks = ["${google_compute_http_health_check.zero.self_link}"]
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

data "google_compute_backend_service" "baz" {
  name = "${google_compute_backend_service.foobar.name}"
}
`, serviceName, checkName)
}
