// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeHealthCheck_tcp_update(t *testing.T) {
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_tcp(hckName),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeHealthCheck_tcp_update(hckName),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_ssl_port_spec(t *testing.T) {
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_ssl_fixed_port(hckName),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_http_port_spec(t *testing.T) {
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeHealthCheck_http_port_spec(hckName),
				ExpectError: regexp.MustCompile("Error in http_health_check: Must specify port_name when using USE_NAMED_PORT as port_specification."),
			},
			{
				Config: testAccComputeHealthCheck_http_named_port(hckName),
			},
		},
	})
}

func TestAccComputeHealthCheck_https_serving_port(t *testing.T) {
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_https_serving_port(hckName),
			},
			{
				ResourceName:      "google_compute_health_check.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeHealthCheck_typeTransition(t *testing.T) {
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeHealthCheck_https(hckName),
			},
			{
				Config: testAccComputeHealthCheck_http(hckName),
			},
			{
				Config: testAccComputeHealthCheck_ssl(hckName),
			},
			{
				Config: testAccComputeHealthCheck_tcp(hckName),
			},
			{
				Config: testAccComputeHealthCheck_http2(hckName),
			},
			{
				Config: testAccComputeHealthCheck_https(hckName),
			},
		},
	})
}

func TestAccComputeHealthCheck_tcpAndSsl_shouldFail(t *testing.T) {
	// No HTTP interactions, is a unit test
	acctest.SkipIfVcr(t)
	t.Parallel()

	hckName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeHealthCheckDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeHealthCheck_tcpAndSsl_shouldFail(hckName),
				ExpectError: regexp.MustCompile("only one of\n`grpc_health_check,http2_health_check,http_health_check,https_health_check,ssl_health_check,tcp_health_check`\ncan be specified, but `ssl_health_check,tcp_health_check` were specified"),
			},
		},
	})
}

func testAccComputeHealthCheck_tcp(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  tcp_health_check {
    port = 443
  }
}
`, hckName)
}

func testAccComputeHealthCheck_tcp_update(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  healthy_threshold   = 10
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 10
  tcp_health_check {
    port = "8080"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_ssl(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  ssl_health_check {
    port = "443"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_ssl_fixed_port(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  ssl_health_check {
    port               = "443"
    port_specification = "USE_FIXED_PORT"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_http(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  http_health_check {
    port = "80"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_http_port_spec(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  http_health_check {
    port_specification = "USE_NAMED_PORT"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_http_named_port(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  http_health_check {
    port_name          = "http"
    port_specification = "USE_NAMED_PORT"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_https(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  https_health_check {
    port = "443"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_https_serving_port(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  https_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_http2(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3
  http2_health_check {
    port = "443"
  }
}
`, hckName)
}

func testAccComputeHealthCheck_tcpAndSsl_shouldFail(hckName string) string {
	return fmt.Sprintf(`
resource "google_compute_health_check" "foobar" {
  check_interval_sec  = 3
  description         = "Resource created for Terraform acceptance testing"
  healthy_threshold   = 3
  name                = "health-test-%s"
  timeout_sec         = 2
  unhealthy_threshold = 3

  tcp_health_check {
    port = 443
  }
  ssl_health_check {
    port = 443
  }
}
`, hckName)
}
