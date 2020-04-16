package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRegionInstanceGroupManager_basic(t *testing.T) {
	t.Parallel()

	template := fmt.Sprintf("igm-test-%s", randString(t, 10))
	target := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igm1 := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igm2 := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_basic(template, target, igm1, igm2),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-no-tp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_targetSizeZero(t *testing.T) {
	t.Parallel()

	templateName := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igmName := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_targetSizeZero(templateName, igmName),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_update(t *testing.T) {
	t.Parallel()

	template1 := fmt.Sprintf("igm-test-%s", randString(t, 10))
	target1 := fmt.Sprintf("igm-test-%s", randString(t, 10))
	target2 := fmt.Sprintf("igm-test-%s", randString(t, 10))
	template2 := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igm := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_update(template1, target1, igm),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRegionInstanceGroupManager_update2(template1, target1, target2, template2, igm),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_updateLifecycle(t *testing.T) {
	t.Parallel()

	tag1 := "tag1"
	tag2 := "tag2"
	igm := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_updateLifecycle(tag1, igm),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRegionInstanceGroupManager_updateLifecycle(tag2, igm),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_rollingUpdatePolicy(t *testing.T) {
	t.Parallel()

	igm := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_rollingUpdatePolicy(igm),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-rolling-update-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccRegionInstanceGroupManager_rollingUpdatePolicySetToDefault(igm),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccRegionInstanceGroupManager_rollingUpdatePolicy2(igm),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-rolling-update-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_separateRegions(t *testing.T) {
	t.Parallel()

	igm1 := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igm2 := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_separateRegions(igm1, igm2),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic-2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_versions(t *testing.T) {
	t.Parallel()

	primaryTemplate := fmt.Sprintf("igm-test-%s", randString(t, 10))
	canaryTemplate := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igm := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_versions(primaryTemplate, canaryTemplate, igm),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_autoHealingPolicies(t *testing.T) {
	t.Parallel()

	template := fmt.Sprintf("igm-test-%s", randString(t, 10))
	target := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igm := fmt.Sprintf("igm-test-%s", randString(t, 10))
	hck := fmt.Sprintf("igm-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_autoHealingPolicies(template, target, igm, hck),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRegionInstanceGroupManager_autoHealingPoliciesRemoved(template, target, igm, hck),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_distributionPolicy(t *testing.T) {
	t.Parallel()

	template := fmt.Sprintf("igm-test-%s", randString(t, 10))
	igm := fmt.Sprintf("igm-test-%s", randString(t, 10))
	zones := []string{"us-central1-a", "us-central1-b"}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRegionInstanceGroupManager_distributionPolicy(template, igm, zones),
			},
			{
				ResourceName:      "google_compute_region_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRegionInstanceGroupManagerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_region_instance_group_manager" {
				continue
			}
			_, err := config.clientCompute.RegionInstanceGroupManagers.Get(
				rs.Primary.Attributes["project"], rs.Primary.Attributes["region"], rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("RegionInstanceGroupManager still exists")
			}
		}

		return nil
	}
}

func testAccRegionInstanceGroupManager_basic(template, target, igm1, igm2 string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "igm-basic" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "primary"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  target_pools       = [google_compute_target_pool.igm-basic.self_link]
  base_instance_name = "igm-basic"
  region             = "us-central1"
  target_size        = 2
}

resource "google_compute_region_instance_group_manager" "igm-no-tp" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "primary"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  base_instance_name = "igm-no-tp"
  region             = "us-central1"
  target_size        = 2
}
`, template, target, igm1, igm2)
}

func testAccRegionInstanceGroupManager_targetSizeZero(template, igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "primary"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  base_instance_name = "igm-basic"
  region             = "us-central1"
}
`, template, igm)
}

func testAccRegionInstanceGroupManager_update(template, target, igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-update" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "igm-update" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "igm-update" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "primary"
    instance_template = google_compute_instance_template.igm-update.self_link
  }

  target_pools       = [google_compute_target_pool.igm-update.self_link]
  base_instance_name = "igm-update"
  region             = "us-central1"
  target_size        = 2
  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, template, target, igm)
}

// Change IGM's instance template and target size
func testAccRegionInstanceGroupManager_update2(template1, target1, target2, template2, igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-update" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "igm-update" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_target_pool" "igm-update2" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_instance_template" "igm-update2" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_region_instance_group_manager" "igm-update" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    instance_template = google_compute_instance_template.igm-update2.self_link
    name              = "primary"
  }

  target_pools = [
    google_compute_target_pool.igm-update.self_link,
    google_compute_target_pool.igm-update2.self_link,
  ]
  base_instance_name = "igm-update"
  region             = "us-central1"
  target_size        = 3
  named_port {
    name = "customhttp"
    port = 8080
  }
  named_port {
    name = "customhttps"
    port = 8443
  }
}
`, template1, target1, target2, template2, igm)
}

func testAccRegionInstanceGroupManager_updateLifecycle(tag, igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-update" {
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["%s"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_region_instance_group_manager" "igm-update" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    instance_template = google_compute_instance_template.igm-update.self_link
    name              = "primary"
  }

  base_instance_name = "igm-update"
  region             = "us-central1"
  target_size        = 2
  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, tag, igm)
}

func testAccRegionInstanceGroupManager_separateRegions(igm1, igm2 string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "primary"
  }

  base_instance_name = "igm-basic"
  region             = "us-central1"
  target_size        = 2
}

resource "google_compute_region_instance_group_manager" "igm-basic-2" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "primary"
  }

  base_instance_name = "igm-basic-2"
  region             = "us-west1"
  target_size        = 2
}
`, igm1, igm2)
}

func testAccRegionInstanceGroupManager_autoHealingPolicies(template, target, igm, hck string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]
  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "igm-basic" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "primary"
  }
  target_pools       = [google_compute_target_pool.igm-basic.self_link]
  base_instance_name = "igm-basic"
  region             = "us-central1"
  target_size        = 2
  auto_healing_policies {
    health_check      = google_compute_http_health_check.zero.self_link
    initial_delay_sec = "10"
  }
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, template, target, igm, hck)
}

func testAccRegionInstanceGroupManager_autoHealingPoliciesRemoved(template, target, igm, hck string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]
  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_target_pool" "igm-basic" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "primary"
  }
  target_pools       = [google_compute_target_pool.igm-basic.self_link]
  base_instance_name = "igm-basic"
  region             = "us-central1"
  target_size        = 2
}

resource "google_compute_http_health_check" "zero" {
  name               = "%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}
`, template, target, igm, hck)
}

func testAccRegionInstanceGroupManager_versions(primaryTemplate string, canaryTemplate string, igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-primary" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]
  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_instance_template" "igm-canary" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]
  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
  description        = "Terraform test region instance group manager"
  name               = "%s"
  base_instance_name = "igm-basic"
  region             = "us-central1"
  target_size        = 2

  version {
    name              = "primary"
    instance_template = google_compute_instance_template.igm-primary.self_link
  }

  version {
    name              = "canary"
    instance_template = google_compute_instance_template.igm-canary.self_link
    target_size {
      fixed = 1
    }
  }
}
`, primaryTemplate, canaryTemplate, igm)
}

func testAccRegionInstanceGroupManager_distributionPolicy(template, igm string, zones []string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-basic" {
  name           = "%s"
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["foo", "bar"]
  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }
  network_interface {
    network = "default"
  }
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "primary"
  }

  base_instance_name        = "igm-basic"
  region                    = "us-central1"
  target_size               = 2
  distribution_policy_zones = ["%s"]
}
`, template, igm, strings.Join(zones, "\",\""))
}

func testAccRegionInstanceGroupManager_rollingUpdatePolicy(igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-rolling-update-policy" {
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["terraform-testing"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_region_instance_group_manager" "igm-rolling-update-policy" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.igm-rolling-update-policy.self_link
    name              = "primary"
  }
  base_instance_name        = "igm-rolling-update-policy"
  region                    = "us-central1"
  target_size               = 4
  distribution_policy_zones = ["us-central1-a", "us-central1-f"]

  update_policy {
    type                  = "PROACTIVE"
    minimal_action        = "REPLACE"
    max_surge_fixed       = 2
    max_unavailable_fixed = 2
    min_ready_sec         = 20
  }

  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, igm)
}

func testAccRegionInstanceGroupManager_rollingUpdatePolicySetToDefault(igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-rolling-update-policy" {
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["terraform-testing"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_region_instance_group_manager" "igm-rolling-update-policy" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.igm-rolling-update-policy.self_link
    name              = "primary"
  }
  base_instance_name        = "igm-rolling-update-policy"
  region                    = "us-central1"
  target_size               = 4
  distribution_policy_zones = ["us-central1-a", "us-central1-f"]

  update_policy {
    type                         = "PROACTIVE"
    instance_redistribution_type = "PROACTIVE"
    minimal_action               = "REPLACE"
    max_surge_fixed              = 2
    max_unavailable_fixed        = 2
    min_ready_sec                = 20
  }

  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, igm)
}

func testAccRegionInstanceGroupManager_rollingUpdatePolicy2(igm string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "igm-rolling-update-policy" {
  machine_type   = "n1-standard-1"
  can_ip_forward = false
  tags           = ["terraform-testing"]

  disk {
    source_image = data.google_compute_image.my_image.self_link
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_region_instance_group_manager" "igm-rolling-update-policy" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    name              = "primary"
    instance_template = google_compute_instance_template.igm-rolling-update-policy.self_link
  }
  base_instance_name        = "igm-rolling-update-policy"
  region                    = "us-central1"
  distribution_policy_zones = ["us-central1-a", "us-central1-f"]
  target_size               = 3
  update_policy {
    type                         = "PROACTIVE"
    instance_redistribution_type = "NONE"
    minimal_action               = "REPLACE"
    max_surge_fixed              = 2
    max_unavailable_fixed        = 0
    min_ready_sec                = 10
  }
  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, igm)
}
