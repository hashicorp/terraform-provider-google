package google

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("ComputeInstanceGroupManager", &resource.Sweeper{
		Name: "ComputeInstanceGroupManager",
		F:    testSweepComputeInstanceGroupManager,
	})
}

// At the time of writing, the CI only passes us-central1 as the region.
// Since we can read all instances across zones, we don't really use this param.
func testSweepComputeInstanceGroupManager(region string) error {
	resourceName := "ComputeInstanceGroupManager"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	found, err := config.clientCompute.InstanceGroupManagers.AggregatedList(config.Project).Do()
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request: %s", err)
		return nil
	}

	// Keep count of items that aren't sweepable for logging.
	nonPrefixCount := 0
	for zone, itemList := range found.Items {
		for _, igm := range itemList.InstanceGroupManagers {
			if !isSweepableTestResource(igm.Name) {
				nonPrefixCount++
				continue
			}

			// Don't wait on operations as we may have a lot to delete
			_, err := config.clientCompute.InstanceGroupManagers.Delete(config.Project, GetResourceNameFromSelfLink(zone), igm.Name).Do()
			if err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error deleting %s resource %s : %s", resourceName, igm.Name, err)
			} else {
				log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, igm.Name)
			}
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items were non-sweepable and skipped.", nonPrefixCount)
	}

	return nil
}

func TestAccInstanceGroupManager_basic(t *testing.T) {
	t.Parallel()

	template := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	target := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	igm1 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	igm2 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_basic(template, target, igm1, igm2),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-no-tp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_targetSizeZero(t *testing.T) {
	t.Parallel()

	templateName := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	igmName := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_targetSizeZero(templateName, igmName),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_update(t *testing.T) {
	t.Parallel()

	template1 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	target1 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	target2 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	template2 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	igm := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_update(template1, target1, igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceGroupManager_update2(template1, target1, target2, template2, igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceGroupManager_update3(template1, target1, target2, template2, igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_updateLifecycle(t *testing.T) {
	// Randomness in instance template
	skipIfVcr(t)
	t.Parallel()

	tag1 := "tag1"
	tag2 := "tag2"
	igm := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_updateLifecycle(tag1, igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceGroupManager_updateLifecycle(tag2, igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_updatePolicy(t *testing.T) {
	// Randomness in instance template
	skipIfVcr(t)
	t.Parallel()

	igm := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_rollingUpdatePolicy(igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-rolling-update-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceGroupManager_rollingUpdatePolicy2(igm),
			},

			{
				ResourceName:      "google_compute_instance_group_manager.igm-rolling-update-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceGroupManager_rollingUpdatePolicy3(igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-rolling-update-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceGroupManager_rollingUpdatePolicy4(igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-rolling-update-policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_separateRegions(t *testing.T) {
	// Randomness in instance template
	skipIfVcr(t)
	t.Parallel()

	igm1 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	igm2 := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_separateRegions(igm1, igm2),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-basic-2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_versions(t *testing.T) {
	t.Parallel()

	primaryTemplate := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	canaryTemplate := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	igm := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_versions(primaryTemplate, canaryTemplate, igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_autoHealingPolicies(t *testing.T) {
	t.Parallel()

	template := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	target := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	igm := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))
	hck := fmt.Sprintf("tf-test-igm-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_autoHealingPolicies(template, target, igm, hck),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceGroupManager_autoHealingPoliciesRemoved(template, target, igm, hck),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckInstanceGroupManagerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_instance_group_manager" {
				continue
			}
			_, err := config.clientCompute.InstanceGroupManagers.Get(
				config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("InstanceGroupManager still exists")
			}
		}

		return nil
	}
}

func testAccInstanceGroupManager_basic(template, target, igm1, igm2 string) string {
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

resource "google_compute_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  target_pools       = [google_compute_target_pool.igm-basic.self_link]
  base_instance_name = "igm-basic"
  zone               = "us-central1-c"
  target_size        = 2
}

resource "google_compute_instance_group_manager" "igm-no-tp" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  base_instance_name = "igm-no-tp"
  zone               = "us-central1-c"
  target_size        = 2
}
`, template, target, igm1, igm2)
}

func testAccInstanceGroupManager_targetSizeZero(template, igm string) string {
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

resource "google_compute_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  base_instance_name = "igm-basic"
  zone               = "us-central1-c"
}
`, template, igm)
}

func testAccInstanceGroupManager_update(template, target, igm string) string {
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

resource "google_compute_instance_group_manager" "igm-update" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-update.self_link
  }

  target_pools       = [google_compute_target_pool.igm-update.self_link]
  base_instance_name = "igm-update"
  zone               = "us-central1-c"
  target_size        = 2
  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, template, target, igm)
}

// Change IGM's instance template and target size
func testAccInstanceGroupManager_update2(template1, target1, target2, template2, igm string) string {
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

resource "google_compute_instance_group_manager" "igm-update" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-update2.self_link
  }

  target_pools = [
    google_compute_target_pool.igm-update.self_link,
    google_compute_target_pool.igm-update2.self_link,
  ]
  base_instance_name = "igm-update"
  zone               = "us-central1-c"
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

// Remove target pools
func testAccInstanceGroupManager_update3(template1, target1, target2, template2, igm string) string {
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

resource "google_compute_instance_group_manager" "igm-update" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-update2.self_link
  }

  base_instance_name = "igm-update"
  zone               = "us-central1-c"
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

func testAccInstanceGroupManager_updateLifecycle(tag, igm string) string {
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

resource "google_compute_instance_group_manager" "igm-update" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-update.self_link
  }

  base_instance_name = "igm-update"
  zone               = "us-central1-c"
  target_size        = 2
  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, tag, igm)
}

func testAccInstanceGroupManager_rollingUpdatePolicy(igm string) string {
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

resource "google_compute_instance_group_manager" "igm-rolling-update-policy" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-rolling-update-policy.self_link
  }
  base_instance_name = "igm-rolling-update-policy"
  zone               = "us-central1-c"
  target_size        = 3
  update_policy {
    type                    = "PROACTIVE"
    minimal_action          = "REPLACE"
    max_surge_percent       = 50
    max_unavailable_percent = 50
    min_ready_sec           = 20
  }
  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, igm)
}

func testAccInstanceGroupManager_rollingUpdatePolicy2(igm string) string {
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

resource "google_compute_instance_group_manager" "igm-rolling-update-policy" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    name              = "prod2"
    instance_template = google_compute_instance_template.igm-rolling-update-policy.self_link
  }
  base_instance_name = "igm-rolling-update-policy"
  zone               = "us-central1-c"
  target_size        = 3
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

func testAccInstanceGroupManager_rollingUpdatePolicy3(igm string) string {
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

resource "google_compute_instance_group_manager" "igm-rolling-update-policy" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    name              = "prod2"
    instance_template = google_compute_instance_template.igm-rolling-update-policy.self_link
  }
  base_instance_name = "igm-rolling-update-policy"
  zone               = "us-central1-c"
  target_size        = 3
  update_policy {
    type                  = "PROACTIVE"
    minimal_action        = "REPLACE"
    max_surge_fixed       = 0
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

func testAccInstanceGroupManager_rollingUpdatePolicy4(igm string) string {
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

resource "google_compute_instance_group_manager" "igm-rolling-update-policy" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    name              = "prod2"
    instance_template = google_compute_instance_template.igm-rolling-update-policy.self_link
  }
  base_instance_name = "igm-rolling-update-policy"
  zone               = "us-central1-c"
  target_size        = 3
  update_policy {
    type                  = "PROACTIVE"
    minimal_action        = "REPLACE"
    max_surge_fixed       = 2
    max_unavailable_fixed = 0
    min_ready_sec         = 20
  }
  named_port {
    name = "customhttp"
    port = 8080
  }
}
`, igm)
}

func testAccInstanceGroupManager_separateRegions(igm1, igm2 string) string {
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

resource "google_compute_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "prod"
  }

  base_instance_name = "igm-basic"
  zone               = "us-central1-c"
  target_size        = 2
}

resource "google_compute_instance_group_manager" "igm-basic-2" {
  description = "Terraform test instance group manager"
  name        = "%s"

  version {
    name              = "prod"
    instance_template = google_compute_instance_template.igm-basic.self_link
  }

  base_instance_name = "igm-basic-2"
  zone               = "us-west1-b"
  target_size        = 2
}
`, igm1, igm2)
}

func testAccInstanceGroupManager_autoHealingPolicies(template, target, igm, hck string) string {
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

resource "google_compute_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "prod"
  }
  target_pools       = [google_compute_target_pool.igm-basic.self_link]
  base_instance_name = "igm-basic"
  zone               = "us-central1-c"
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

func testAccInstanceGroupManager_autoHealingPoliciesRemoved(template, target, igm, hck string) string {
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

resource "google_compute_instance_group_manager" "igm-basic" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.igm-basic.self_link
    name              = "prod"
  }
  target_pools       = [google_compute_target_pool.igm-basic.self_link]
  base_instance_name = "igm-basic"
  zone               = "us-central1-c"
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

func testAccInstanceGroupManager_versions(primaryTemplate string, canaryTemplate string, igm string) string {
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

resource "google_compute_instance_group_manager" "igm-basic" {
  description        = "Terraform test instance group manager"
  name               = "%s"
  base_instance_name = "igm-basic"
  zone               = "us-central1-c"
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
