package google

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"

	"sort"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRegionInstanceGroupManager_basic(t *testing.T) {
	var manager compute.InstanceGroupManager

	template := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	target := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_basic(template, target, igm1, igm2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-basic", &manager),
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-no-tp", &manager),
				),
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_targetSizeZero(t *testing.T) {
	var manager compute.InstanceGroupManager

	templateName := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igmName := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_targetSizeZero(templateName, igmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-basic", &manager),
				),
			},
		},
	})

	if manager.TargetSize != 0 {
		t.Errorf("Expected target_size to be 0, got %d", manager.TargetSize)
	}
}

func TestAccRegionInstanceGroupManager_update(t *testing.T) {
	var manager compute.InstanceGroupManager

	template1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	target1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	target2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	template2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_update(template1, target1, igm),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-update", &manager),
					testAccCheckRegionInstanceGroupManagerNamedPorts(
						"google_compute_region_instance_group_manager.igm-update",
						map[string]int64{"customhttp": 8080},
						&manager),
				),
			},
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_update2(template1, target1, target2, template2, igm),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-update", &manager),
					testAccCheckRegionInstanceGroupManagerUpdated(
						"google_compute_region_instance_group_manager.igm-update", 3,
						[]string{target1, target2}, template2),
					testAccCheckRegionInstanceGroupManagerNamedPorts(
						"google_compute_region_instance_group_manager.igm-update",
						map[string]int64{"customhttp": 8080, "customhttps": 8443},
						&manager),
				),
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_updateLifecycle(t *testing.T) {
	var manager compute.InstanceGroupManager

	tag1 := "tag1"
	tag2 := "tag2"
	igm := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_updateLifecycle(tag1, igm),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-update", &manager),
				),
			},
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_updateLifecycle(tag2, igm),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-update", &manager),
					testAccCheckRegionInstanceGroupManagerTemplateTags(
						"google_compute_region_instance_group_manager.igm-update", []string{tag2}),
				),
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_separateRegions(t *testing.T) {
	var manager compute.InstanceGroupManager

	igm1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_separateRegions(igm1, igm2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-basic", &manager),
					testAccCheckRegionInstanceGroupManagerExists(
						"google_compute_region_instance_group_manager.igm-basic-2", &manager),
				),
			},
		},
	})
}

func TestAccRegionInstanceGroupManager_autoHealingPolicies(t *testing.T) {
	var manager computeBeta.InstanceGroupManager

	template := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	target := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	hck := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRegionInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRegionInstanceGroupManager_autoHealingPolicies(template, target, igm, hck),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegionInstanceGroupManagerBetaExists(
						"google_compute_region_instance_group_manager.igm-basic", &manager),
					testAccCheckRegionInstanceGroupManagerAutoHealingPolicies("google_compute_region_instance_group_manager.igm-basic", hck, 10),
				),
			},
		},
	})
}

func testAccCheckRegionInstanceGroupManagerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_region_instance_group_manager" {
			continue
		}
		_, err := config.clientCompute.RegionInstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("RegionInstanceGroupManager still exists")
		}
	}

	return nil
}

func testAccCheckRegionInstanceGroupManagerExists(n string, manager *compute.InstanceGroupManager) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.RegionInstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("RegionInstanceGroupManager not found")
		}

		*manager = *found

		return nil
	}
}

func testAccCheckRegionInstanceGroupManagerBetaExists(n string, manager *computeBeta.InstanceGroupManager) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientComputeBeta.RegionInstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("RegionInstanceGroupManager not found")
		}

		*manager = *found

		return nil
	}
}

func testAccCheckRegionInstanceGroupManagerUpdated(n string, size int64, targetPools []string, template string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		manager, err := config.clientCompute.RegionInstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		// Cannot check the target pool as the instance creation is asynchronous.  However, can
		// check the target_size.
		if manager.TargetSize != size {
			return fmt.Errorf("instance count incorrect")
		}

		tpNames := make([]string, 0, len(manager.TargetPools))
		for _, targetPool := range manager.TargetPools {
			targetPoolParts := strings.Split(targetPool, "/")
			tpNames = append(tpNames, targetPoolParts[len(targetPoolParts)-1])
		}

		sort.Strings(tpNames)
		sort.Strings(targetPools)
		if !reflect.DeepEqual(tpNames, targetPools) {
			return fmt.Errorf("target pools incorrect. Expected %s, got %s", targetPools, tpNames)
		}

		// check that the instance template updated
		instanceTemplate, err := config.clientCompute.InstanceTemplates.Get(
			config.Project, template).Do()
		if err != nil {
			return fmt.Errorf("Error reading instance template: %s", err)
		}

		if instanceTemplate.Name != template {
			return fmt.Errorf("instance template not updated")
		}

		return nil
	}
}

func testAccCheckRegionInstanceGroupManagerNamedPorts(n string, np map[string]int64, instanceGroupManager *compute.InstanceGroupManager) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		manager, err := config.clientCompute.RegionInstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		var found bool
		for _, namedPort := range manager.NamedPorts {
			found = false
			for name, port := range np {
				if namedPort.Name == name && namedPort.Port == port {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("named port incorrect")
			}
		}

		return nil
	}
}

func testAccCheckRegionInstanceGroupManagerAutoHealingPolicies(n, hck string, initialDelaySec int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		manager, err := config.clientComputeBeta.RegionInstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if len(manager.AutoHealingPolicies) != 1 {
			return fmt.Errorf("Expected # of auto healing policies to be 1, got %d", len(manager.AutoHealingPolicies))
		}
		autoHealingPolicy := manager.AutoHealingPolicies[0]

		if !strings.Contains(autoHealingPolicy.HealthCheck, hck) {
			return fmt.Errorf("Expected string \"%s\" to appear in \"%s\"", hck, autoHealingPolicy.HealthCheck)
		}

		if autoHealingPolicy.InitialDelaySec != initialDelaySec {
			return fmt.Errorf("Expected auto healing policy inital delay to be %d, got %d", initialDelaySec, autoHealingPolicy.InitialDelaySec)
		}
		return nil
	}
}

func testAccCheckRegionInstanceGroupManagerTemplateTags(n string, tags []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		manager, err := config.clientCompute.RegionInstanceGroupManagers.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		// check that the instance template updated
		instanceTemplate, err := config.clientCompute.InstanceTemplates.Get(
			config.Project, resourceSplitter(manager.InstanceTemplate)).Do()
		if err != nil {
			return fmt.Errorf("Error reading instance template: %s", err)
		}

		if !reflect.DeepEqual(instanceTemplate.Properties.Tags.Items, tags) {
			return fmt.Errorf("instance template not updated")
		}

		return nil
	}
}

func testAccRegionInstanceGroupManager_basic(template, target, igm1, igm2 string) string {
	return fmt.Sprintf(`
	resource "google_compute_instance_template" "igm-basic" {
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "debian-cloud/debian-8-jessie-v20160803"
			auto_delete = true
			boot = true
		}

		network_interface {
			network = "default"
		}

		metadata {
			foo = "bar"
		}

		service_account {
			scopes = ["userinfo-email", "compute-ro", "storage-ro"]
		}
	}

	resource "google_compute_target_pool" "igm-basic" {
		description = "Resource created for Terraform acceptance testing"
		name = "%s"
		session_affinity = "CLIENT_IP_PROTO"
	}

	resource "google_compute_region_instance_group_manager" "igm-basic" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		target_pools = ["${google_compute_target_pool.igm-basic.self_link}"]
		base_instance_name = "igm-basic"
		region = "us-central1"
		target_size = 2
	}

	resource "google_compute_region_instance_group_manager" "igm-no-tp" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-no-tp"
		region = "us-central1"
		target_size = 2
	}
	`, template, target, igm1, igm2)
}

func testAccRegionInstanceGroupManager_targetSizeZero(template, igm string) string {
	return fmt.Sprintf(`
	resource "google_compute_instance_template" "igm-basic" {
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "debian-cloud/debian-8-jessie-v20160803"
			auto_delete = true
			boot = true
		}

		network_interface {
			network = "default"
		}

		metadata {
			foo = "bar"
		}

		service_account {
			scopes = ["userinfo-email", "compute-ro", "storage-ro"]
		}
	}

	resource "google_compute_region_instance_group_manager" "igm-basic" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-basic"
		region = "us-central1"
	}
	`, template, igm)
}

func testAccRegionInstanceGroupManager_update(template, target, igm string) string {
	return fmt.Sprintf(`
	resource "google_compute_instance_template" "igm-update" {
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "debian-cloud/debian-8-jessie-v20160803"
			auto_delete = true
			boot = true
		}

		network_interface {
			network = "default"
		}

		metadata {
			foo = "bar"
		}

		service_account {
			scopes = ["userinfo-email", "compute-ro", "storage-ro"]
		}
	}

	resource "google_compute_target_pool" "igm-update" {
		description = "Resource created for Terraform acceptance testing"
		name = "%s"
		session_affinity = "CLIENT_IP_PROTO"
	}

	resource "google_compute_region_instance_group_manager" "igm-update" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-update.self_link}"
		target_pools = ["${google_compute_target_pool.igm-update.self_link}"]
		base_instance_name = "igm-update"
		region = "us-central1"
		target_size = 2
		named_port {
			name = "customhttp"
			port = 8080
		}
	}`, template, target, igm)
}

// Change IGM's instance template and target size
func testAccRegionInstanceGroupManager_update2(template1, target1, target2, template2, igm string) string {
	return fmt.Sprintf(`
	resource "google_compute_instance_template" "igm-update" {
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "debian-cloud/debian-8-jessie-v20160803"
			auto_delete = true
			boot = true
		}

		network_interface {
			network = "default"
		}

		metadata {
			foo = "bar"
		}

		service_account {
			scopes = ["userinfo-email", "compute-ro", "storage-ro"]
		}
	}

	resource "google_compute_target_pool" "igm-update" {
		description = "Resource created for Terraform acceptance testing"
		name = "%s"
		session_affinity = "CLIENT_IP_PROTO"
	}

	resource "google_compute_target_pool" "igm-update2" {
		description = "Resource created for Terraform acceptance testing"
		name = "%s"
		session_affinity = "CLIENT_IP_PROTO"
	}

	resource "google_compute_instance_template" "igm-update2" {
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "debian-cloud/debian-8-jessie-v20160803"
			auto_delete = true
			boot = true
		}

		network_interface {
			network = "default"
		}

		metadata {
			foo = "bar"
		}

		service_account {
			scopes = ["userinfo-email", "compute-ro", "storage-ro"]
		}
	}

	resource "google_compute_region_instance_group_manager" "igm-update" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-update2.self_link}"
		target_pools = [
			"${google_compute_target_pool.igm-update.self_link}",
			"${google_compute_target_pool.igm-update2.self_link}",
		]
		base_instance_name = "igm-update"
		region = "us-central1"
		target_size = 3
		named_port {
			name = "customhttp"
			port = 8080
		}
		named_port {
			name = "customhttps"
			port = 8443
		}
	}`, template1, target1, target2, template2, igm)
}

func testAccRegionInstanceGroupManager_updateLifecycle(tag, igm string) string {
	return fmt.Sprintf(`
	resource "google_compute_instance_template" "igm-update" {
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["%s"]

		disk {
			source_image = "debian-cloud/debian-8-jessie-v20160803"
			auto_delete = true
			boot = true
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
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-update.self_link}"
		base_instance_name = "igm-update"
		region = "us-central1"
		target_size = 2
		named_port {
			name = "customhttp"
			port = 8080
		}
	}`, tag, igm)
}

func testAccRegionInstanceGroupManager_separateRegions(igm1, igm2 string) string {
	return fmt.Sprintf(`
	resource "google_compute_instance_template" "igm-basic" {
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "debian-cloud/debian-8-jessie-v20160803"
			auto_delete = true
			boot = true
		}

		network_interface {
			network = "default"
		}

		metadata {
			foo = "bar"
		}

		service_account {
			scopes = ["userinfo-email", "compute-ro", "storage-ro"]
		}
	}

	resource "google_compute_region_instance_group_manager" "igm-basic" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-basic"
		region = "us-central1"
		target_size = 2
	}

	resource "google_compute_region_instance_group_manager" "igm-basic-2" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-basic-2"
		region = "us-west1"
		target_size = 2
	}
	`, igm1, igm2)
}

func testAccRegionInstanceGroupManager_autoHealingPolicies(template, target, igm, hck string) string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "igm-basic" {
	name = "%s"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]
	disk {
		source_image = "debian-cloud/debian-8-jessie-v20160803"
		auto_delete = true
		boot = true
	}
	network_interface {
		network = "default"
	}
	metadata {
		foo = "bar"
	}
	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}

resource "google_compute_target_pool" "igm-basic" {
	description = "Resource created for Terraform acceptance testing"
	name = "%s"
	session_affinity = "CLIENT_IP_PROTO"
}

resource "google_compute_region_instance_group_manager" "igm-basic" {
	description = "Terraform test instance group manager"
	name = "%s"
	instance_template = "${google_compute_instance_template.igm-basic.self_link}"
	target_pools = ["${google_compute_target_pool.igm-basic.self_link}"]
	base_instance_name = "igm-basic"
	region = "us-central1"
	target_size = 2
	auto_healing_policies {
		health_check = "${google_compute_http_health_check.zero.self_link}"
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
