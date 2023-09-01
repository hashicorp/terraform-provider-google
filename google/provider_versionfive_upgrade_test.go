// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestProvider_versionfive_upgrade(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	name1 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	name2 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	name3 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	name4 := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.80.0",
						Source:            "hashicorp/google-beta",
					},
				},
				Config: testProvider_versionfive_upgrades(project, org, billingId, name1, name2, name3, name4),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testProvider_versionfive_upgrades(project, org, billingId, name1, name2, name3, name4),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_data_fusion_instance.unset",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_data_fusion_instance.set",
				ImportState:              true,
				ImportStateVerify:        true,
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				ResourceName:             "google_data_fusion_instance.reference",
				ImportState:              true,
				ImportStateVerify:        true,
			},
		},
	})
}

func TestProvider_versionfive_ignorereads_upgrade(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	var itName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerName = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	diskName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	policyName := fmt.Sprintf("tf-test-policy-%s", acctest.RandString(t, 10))

	endpointContext := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	var itNameRegion = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var tpNameRegion = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var igmNameRegion = fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	var autoscalerNameRegion = fmt.Sprintf("tf-test-region-autoscaler-%s", acctest.RandString(t, 10))

	policyContext := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	attachmentContext := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: resource.ComposeTestCheckFunc(testAccCheckComputeResourcePolicyDestroyProducer(t),
			testAccCheckComputeRegionAutoscalerDestroyProducer(t),
			testAccCheckComputeNetworkEndpointGroupDestroyProducer(t),
			testAccCheckComputeAutoscalerDestroyProducer(t),
		),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"google": {
						VersionConstraint: "4.80.0",
						Source:            "hashicorp/google-beta",
					},
				},
				Config: testProvider_versionfive_upgrades_ignorereads(itName, tpName, igmName, autoscalerName, diskName, policyName,
					itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion, endpointContext,
					policyContext, attachmentContext),
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config: testProvider_versionfive_upgrades_ignorereads(itName, tpName, igmName, autoscalerName, diskName, policyName,
					itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion, endpointContext,
					policyContext, attachmentContext),
				PlanOnly: true,
			},
		},
	})
}

func testProvider_versionfive_upgrades(project, org, billing, name1, name2, name3, name4 string) string {
	return fmt.Sprintf(`
resource "google_project" "host" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "dfapi" {
  project = google_project.host.project_id
  service = "datafusion.googleapis.com"

  disable_dependent_services = false
}

resource "google_data_fusion_instance" "unset" {
  name   = "%s"
  type   = "BASIC"
  options = {
  	prober_test_run = "true"
  }
}

resource "google_data_fusion_instance" "set" {
  name   = "%s"
  region = "us-west1"
  type   = "BASIC"
  options = {
  	prober_test_run = "true"
  }
}

resource "google_data_fusion_instance" "reference" {
  project = google_project.host.project_id
  name   = "%s"
  type   = "DEVELOPER"
  options = {
  	prober_test_run = "true"
  }
  zone   = "us-west1-a"
  depends_on = [
    google_project_service.dfapi
  ]
}

resource "google_redis_instance" "overridewithnonstandardlogic" {
  name           = "%s"
  memory_size_gb = 1
  location_id    = "us-south1-a"
}


`, project, project, org, billing, name1, name2, name3, name4)
}

func testProvider_versionfive_upgrades_ignorereads(itName, tpName, igmName, autoscalerName, diskName, policyName, itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion string, endpointContext, policyContext, attachmentContext map[string]interface{}) string {
	return testAccComputeAutoscaler_basic(itName, tpName, igmName, autoscalerName) +
		testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName) +
		testAccComputeNetworkEndpointGroup_networkEndpointGroup(endpointContext) +
		testAccComputeRegionAutoscaler_basic(itNameRegion, tpNameRegion, igmNameRegion, autoscalerNameRegion) +
		testAccComputeResourcePolicy_resourcePolicyBasicExample(policyContext) +
		testAccComputeServiceAttachment_serviceAttachmentBasicExample(attachmentContext)
}

// need to make copies of all the respective resource functions within here as *_test.go files can not be imported
// checkdestroys
func testAccCheckComputeResourcePolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_resource_policy" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/resourcePolicies/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeResourcePolicy still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccCheckComputeRegionAutoscalerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_region_autoscaler" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/autoscalers/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeRegionAutoscaler still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccCheckComputeNetworkEndpointGroupDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_network_endpoint_group" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/networkEndpointGroups/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeNetworkEndpointGroup still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccCheckComputeAutoscalerDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_autoscaler" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/autoscalers/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeAutoscaler still exists at %s", url)
			}
		}

		return nil
	}
}

// tests
func testAccComputeAutoscaler_scaffolding(itName, tpName, igmName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image1" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar1" {
  name           = "%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image1.self_link
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

resource "google_compute_target_pool" "foobar1" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
  region					 = "us-west1"
}

resource "google_compute_instance_group_manager" "foobar1" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.foobar1.self_link
    name              = "primary"
  }
  target_pools       = [google_compute_target_pool.foobar1.self_link]
  base_instance_name = "foobar1"
  zone               = "us-west1-a"
}
`, itName, tpName, igmName)

}

func testAccComputeAutoscaler_basic(itName, tpName, igmName, autoscalerName string) string {
	return testAccComputeAutoscaler_scaffolding(itName, tpName, igmName) + fmt.Sprintf(`
resource "google_compute_autoscaler" "foobar1" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  zone        = "us-west1-a"
  target      = google_compute_instance_group_manager.foobar1.self_link
  autoscaling_policy {
    max_replicas    = 5
    min_replicas    = 1
    cooldown_period = 60
    cpu_utilization {
      target = 0.5
    }
  }
}
`, autoscalerName)
}

func testAccComputeDiskResourcePolicyAttachment_basic(diskName, policyName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image2" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar2" {
  name  = "%s"
  image = data.google_compute_image.my_image2.self_link
  size  = 1000
  type  = "pd-extreme"
  zone  = "us-west1-c"
  labels = {
    my-label = "my-label-value"
  }
  provisioned_iops = 90000
}

resource "google_compute_resource_policy" "foobar2" {
  name = "%s"
  region = "us-west1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time = "04:00"
      }
    }
  }
}

resource "google_compute_disk_resource_policy_attachment" "foobar2" {
  name = google_compute_resource_policy.foobar2.name
  disk = google_compute_disk.foobar2.name
  zone = "us-west1-c"
}
`, diskName, policyName)
}

func testAccComputeNetworkEndpointGroup_networkEndpointGroup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network_endpoint_group" "neg3" {
  name         = "tf-test-my-lb-neg%{random_suffix}"
  network      = google_compute_network.default3.id
  default_port = "90"
  zone         = "us-west1-a"
}

resource "google_compute_network" "default3" {
  name                    = "tf-test-neg-network%{random_suffix}"
  auto_create_subnetworks = true
}
`, context)
}

func testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName string) string {
	return fmt.Sprintf(`
data "google_compute_image" "my_image4" {
  family  = "debian-11"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "foobar4" {
  name           = "%s"
  machine_type   = "e2-medium"
  can_ip_forward = false
  tags           = ["foo", "bar"]

  disk {
    source_image = data.google_compute_image.my_image4.self_link
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

resource "google_compute_target_pool" "foobar4" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  session_affinity = "CLIENT_IP_PROTO"
  region					 = "us-west1"
}

resource "google_compute_region_instance_group_manager" "foobar4" {
  description = "Terraform test instance group manager"
  name        = "%s"
  version {
    instance_template = google_compute_instance_template.foobar4.self_link
    name              = "primary"
  }
  target_pools       = [google_compute_target_pool.foobar4.self_link]
  base_instance_name = "tf-test-foobar4"
  region             = "us-west1"
}

`, itName, tpName, igmName)
}

func testAccComputeRegionAutoscaler_basic(itName, tpName, igmName, autoscalerName string) string {
	return testAccComputeRegionAutoscaler_scaffolding(itName, tpName, igmName) + fmt.Sprintf(`
resource "google_compute_region_autoscaler" "foobar4" {
  description = "Resource created for Terraform acceptance testing"
  name        = "%s"
  region      = "us-west1"
  target      = google_compute_region_instance_group_manager.foobar4.self_link
  autoscaling_policy {
    max_replicas    = 5
    min_replicas    = 0
    cooldown_period = 60
    cpu_utilization {
      target = 0.5
    }
  }
}
`, autoscalerName)
}

func testAccComputeResourcePolicy_resourcePolicyBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_resource_policy" "foo5" {
  name   = "tf-test-gce-policy%{random_suffix}"
  region = "us-west1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time    = "04:00"
      }
    }
  }
}
`, context)
}

func testAccComputeServiceAttachment_serviceAttachmentBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_service_attachment" "psc_ilb_service_attachment6" {
  name        = "tf-test-my-psc-ilb%{random_suffix}"
  region      = "us-west2"
  description = "A service attachment configured with Terraform"

  enable_proxy_protocol    = true
  connection_preference    = "ACCEPT_AUTOMATIC"
  nat_subnets              = [google_compute_subnetwork.psc_ilb_nat6.id]
  target_service           = google_compute_forwarding_rule.psc_ilb_target_service6.id
  reconcile_connections    = true
}

resource "google_compute_address" "psc_ilb_consumer_address6" {
  name   = "tf-test-psc-ilb-consumer-address%{random_suffix}"
  region = "us-west2"

  subnetwork   = "default"
  address_type = "INTERNAL"
}

resource "google_compute_forwarding_rule" "psc_ilb_consumer6" {
  name   = "tf-test-psc-ilb-consumer-forwarding-rule%{random_suffix}"
  region = "us-west2"

  target                = google_compute_service_attachment.psc_ilb_service_attachment6.id
  load_balancing_scheme = "" # need to override EXTERNAL default when target is a service attachment
  network               = "default"
  ip_address            = google_compute_address.psc_ilb_consumer_address6.id
}

resource "google_compute_forwarding_rule" "psc_ilb_target_service6" {
  name   = "tf-test-producer-forwarding-rule%{random_suffix}"
  region = "us-west2"

  load_balancing_scheme = "INTERNAL"
  backend_service       = google_compute_region_backend_service.producer_service_backend6.id
  all_ports             = true
  network               = google_compute_network.psc_ilb_network6.name
  subnetwork            = google_compute_subnetwork.psc_ilb_producer_subnetwork6.name
}

resource "google_compute_region_backend_service" "producer_service_backend6" {
  name   = "tf-test-producer-service%{random_suffix}"
  region = "us-west2"

  health_checks = [google_compute_health_check.producer_service_health_check6.id]
}

resource "google_compute_health_check" "producer_service_health_check6" {
  name = "tf-test-producer-service-health-check%{random_suffix}"

  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "80"
  }
}

resource "google_compute_network" "psc_ilb_network6" {
  name = "tf-test-psc-ilb-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "psc_ilb_producer_subnetwork6" {
  name   = "tf-test-psc-ilb-producer-subnetwork%{random_suffix}"
  region = "us-west2"

  network       = google_compute_network.psc_ilb_network6.id
  ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_subnetwork" "psc_ilb_nat6" {
  name   = "tf-test-psc-ilb-nat%{random_suffix}"
  region = "us-west2"

  network       = google_compute_network.psc_ilb_network6.id
  purpose       =  "PRIVATE_SERVICE_CONNECT"
  ip_cidr_range = "10.1.0.0/16"
}
`, context)
}
