package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	compute "google.golang.org/api/compute/v1"
)

func TestAccComputeSslPolicy_basic(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslPolicyBasic(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.basic", &sslPolicy),
					// test attribute default values
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.basic", "profile", "COMPATIBLE"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.basic", "min_tls_version", "TLS_1_0"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSslPolicy_profile(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslPolicyProfile(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.profile", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.profile", "profile", "MODERN"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSslPolicy_update(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslUpdate1(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "profile", "MODERN"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "min_tls_version", "TLS_1_0"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccComputeSslUpdate2(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "profile", "RESTRICTED"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "min_tls_version", "TLS_1_2"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSslPolicy_tls_version(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslPolicyTlsVersion(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.tlsversion", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.tlsversion", "min_tls_version", "TLS_1_2"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.tlsversion",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSslPolicy_custom(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslPolicyCustom(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.custom", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.custom", "min_tls_version", "TLS_1_2"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.custom", "profile", "CUSTOM"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.custom",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSslPolicy_update_to_custom(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslUpdate1(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "profile", "MODERN"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "min_tls_version", "TLS_1_0"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccComputeSslUpdate3(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "profile", "CUSTOM"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "min_tls_version", "TLS_1_1"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeSslPolicy_update_from_custom(t *testing.T) {
	t.Parallel()

	var sslPolicy compute.SslPolicy
	sslPolicyName := fmt.Sprintf("test-ssl-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslUpdate3(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "profile", "CUSTOM"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "min_tls_version", "TLS_1_1"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccComputeSslUpdate1(sslPolicyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslPolicyExists(
						"google_compute_ssl_policy.update", &sslPolicy),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "profile", "MODERN"),
					resource.TestCheckResourceAttr(
						"google_compute_ssl_policy.update", "min_tls_version", "TLS_1_0"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_ssl_policy.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeSslPolicyExists(n string, sslPolicy *compute.SslPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		project, err := getTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		name := rs.Primary.Attributes["name"]

		found, err := config.clientCompute.SslPolicies.Get(
			project, name).Do()
		if err != nil {
			return fmt.Errorf("Error Reading SSL Policy %s: %s", name, err)
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("SSL Policy not found")
		}

		*sslPolicy = *found

		return nil
	}
}

func testAccComputeSslPolicyBasic(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "basic" {
	name = "%s"
	description = "Generated by TF provider acceptance test"
}
`, resourceName)
}

func testAccComputeSslPolicyProfile(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "profile" {
	name = "%s"
	description = "Generated by TF provider acceptance test"
	profile = "MODERN"
}
`, resourceName)
}

func testAccComputeSslUpdate1(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "update" {
	name = "%s"
	description = "Generated by TF provider acceptance test"
	min_tls_version = "TLS_1_0"
	profile = "MODERN"
}
`, resourceName)
}

func testAccComputeSslUpdate2(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "update" {
	name = "%s"
	description = "Generated by TF provider acceptance test"
	min_tls_version = "TLS_1_2"
	profile = "RESTRICTED"
}
`, resourceName)
}

func testAccComputeSslUpdate3(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "update" {
	name = "%s"
	description = "Generated by TF provider acceptance test"
	min_tls_version = "TLS_1_1"
	profile = "CUSTOM"
	custom_features = ["TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"]
}
`, resourceName)
}

func testAccComputeSslPolicyTlsVersion(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "tlsversion" {
	name = "%s"
	description = "Generated by TF provider acceptance test"
	min_tls_version = "TLS_1_2"
}
`, resourceName)
}

func testAccComputeSslPolicyCustom(resourceName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_policy" "custom" {
	name = "%s"
	description = "Generated by TF provider acceptance test"
	min_tls_version = "TLS_1_2"
	profile = "CUSTOM"
	custom_features = ["TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"]
}
`, resourceName)
}
