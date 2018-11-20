package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccInstanceGroupManager_basic(t *testing.T) {
	t.Parallel()

	template := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	target := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
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

	templateName := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igmName := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
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

	template1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	target1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	target2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	template2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
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
		},
	})
}

func TestAccInstanceGroupManager_updateLifecycle(t *testing.T) {
	t.Parallel()

	tag1 := "tag1"
	tag2 := "tag2"
	igm := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
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

func TestAccInstanceGroupManager_updateStrategy(t *testing.T) {
	t.Parallel()

	igm := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroupManager_updateStrategy(igm),
			},
			{
				ResourceName:      "google_compute_instance_group_manager.igm-update-strategy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccInstanceGroupManager_separateRegions(t *testing.T) {
	t.Parallel()

	igm1 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))
	igm2 := fmt.Sprintf("igm-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceGroupManagerDestroy,
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

func testAccCheckInstanceGroupManagerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_instance_group_manager" {
			continue
		}
		id, err := parseInstanceGroupManagerId(rs.Primary.ID)
		if err != nil {
			return err
		}
		if id.Project == "" {
			id.Project = config.Project
		}
		if id.Zone == "" {
			id.Zone = rs.Primary.Attributes["zone"]
		}
		_, err = config.clientCompute.InstanceGroupManagers.Get(
			id.Project, id.Zone, id.Name).Do()
		if err == nil {
			return fmt.Errorf("InstanceGroupManager still exists")
		}
	}

	return nil
}

func testAccInstanceGroupManager_basic(template, target, igm1, igm2 string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-9"
		project = "debian-cloud"
	}

	resource "google_compute_instance_template" "igm-basic" {
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "${data.google_compute_image.my_image.self_link}"
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

	resource "google_compute_instance_group_manager" "igm-basic" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		target_pools = ["${google_compute_target_pool.igm-basic.self_link}"]
		base_instance_name = "igm-basic"
		zone = "us-central1-c"
		target_size = 2
	}

	resource "google_compute_instance_group_manager" "igm-no-tp" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-no-tp"
		zone = "us-central1-c"
		target_size = 2
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
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "${data.google_compute_image.my_image.self_link}"
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

	resource "google_compute_instance_group_manager" "igm-basic" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-basic"
		zone = "us-central1-c"
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
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "${data.google_compute_image.my_image.self_link}"
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

	resource "google_compute_instance_group_manager" "igm-update" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-update.self_link}"
		target_pools = ["${google_compute_target_pool.igm-update.self_link}"]
		base_instance_name = "igm-update"
		zone = "us-central1-c"
		target_size = 2
		named_port {
			name = "customhttp"
			port = 8080
		}
	}`, template, target, igm)
}

// Change IGM's instance template and target size
func testAccInstanceGroupManager_update2(template1, target1, target2, template2, igm string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-9"
		project = "debian-cloud"
	}

	resource "google_compute_instance_template" "igm-update" {
		name = "%s"
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "${data.google_compute_image.my_image.self_link}"
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
			source_image = "${data.google_compute_image.my_image.self_link}"
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

	resource "google_compute_instance_group_manager" "igm-update" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-update2.self_link}"
		target_pools = [
			"${google_compute_target_pool.igm-update.self_link}",
			"${google_compute_target_pool.igm-update2.self_link}",
		]
		base_instance_name = "igm-update"
		zone = "us-central1-c"
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

func testAccInstanceGroupManager_updateLifecycle(tag, igm string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-9"
		project = "debian-cloud"
	}

	resource "google_compute_instance_template" "igm-update" {
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["%s"]

		disk {
			source_image = "${data.google_compute_image.my_image.self_link}"
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

	resource "google_compute_instance_group_manager" "igm-update" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-update.self_link}"
		base_instance_name = "igm-update"
		zone = "us-central1-c"
		target_size = 2
		named_port {
			name = "customhttp"
			port = 8080
		}
	}`, tag, igm)
}

func testAccInstanceGroupManager_updateStrategy(igm string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-9"
		project = "debian-cloud"
	}

	resource "google_compute_instance_template" "igm-update-strategy" {
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["terraform-testing"]

		disk {
			source_image = "${data.google_compute_image.my_image.self_link}"
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

	resource "google_compute_instance_group_manager" "igm-update-strategy" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-update-strategy.self_link}"
		base_instance_name = "igm-update-strategy"
		zone = "us-central1-c"
		target_size = 2
		update_strategy = "REPLACE"
		named_port {
			name = "customhttp"
			port = 8080
		}
	}`, igm)
}

func testAccInstanceGroupManager_separateRegions(igm1, igm2 string) string {
	return fmt.Sprintf(`
	data "google_compute_image" "my_image" {
		family  = "debian-9"
		project = "debian-cloud"
	}

	resource "google_compute_instance_template" "igm-basic" {
		machine_type = "n1-standard-1"
		can_ip_forward = false
		tags = ["foo", "bar"]

		disk {
			source_image = "${data.google_compute_image.my_image.self_link}"
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

	resource "google_compute_instance_group_manager" "igm-basic" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-basic"
		zone = "us-central1-c"
		target_size = 2
	}

	resource "google_compute_instance_group_manager" "igm-basic-2" {
		description = "Terraform test instance group manager"
		name = "%s"
		instance_template = "${google_compute_instance_template.igm-basic.self_link}"
		base_instance_name = "igm-basic-2"
		zone = "us-west1-b"
		target_size = 2
	}
	`, igm1, igm2)
}
