package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeFirewallPolicyRule_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", getTestOrgFromEnv(t)),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallPolicyRule_start(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeFirewallPolicyRule_update(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy", "target_resources"},
			},
			{
				Config: testAccComputeFirewallPolicyRule_removeConfigs(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy", "target_resources"},
			},
			{
				Config: testAccComputeFirewallPolicyRule_start(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
		},
	})
}

func testAccComputeFirewallPolicyRule_start(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
  account_id = "tf-test-sa-%{random_suffix}"
}

resource "google_service_account" "service_account2" {
  account_id = "tf-test-sa2-%{random_suffix}"
}

resource "google_compute_network" "network1" {
  name = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name = "tf-test-2-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.name
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_rule" "default" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 9000
  enable_logging = true
  action = "allow"
  direction = "EGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [80, 8080]
    }
    dest_ip_ranges = ["11.100.0.1/32"]
  }
}
`, context)
}

func testAccComputeFirewallPolicyRule_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
  account_id = "tf-test-sa-%{random_suffix}"
}

resource "google_service_account" "service_account2" {
  account_id = "tf-test-sa2-%{random_suffix}"
}

resource "google_compute_network" "network1" {
  name = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name = "tf-test-2-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.name
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_rule" "default" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 9000
  enable_logging = true
  action = "allow"
  direction = "EGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [8080]
    }
    layer4_configs {
      ip_protocol = "udp"
      ports = [22]
    }
    dest_ip_ranges = ["11.100.0.1/32", "10.0.0.0/24"]
  }
  target_resources = [google_compute_network.network1.self_link, google_compute_network.network2.self_link]
  target_service_accounts = [google_service_account.service_account.email]
}
`, context)
}

func testAccComputeFirewallPolicyRule_removeConfigs(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
  account_id = "tf-test-sa-%{random_suffix}"
}

resource "google_service_account" "service_account2" {
  account_id = "tf-test-sa2-%{random_suffix}"
}

resource "google_compute_network" "network1" {
  name = "tf-test-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network2" {
  name = "tf-test-2-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_rule" "default" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Test description"
  priority = 9000
  enable_logging = false
  action = "deny"
  direction = "INGRESS"
  disabled = true
  match {
    layer4_configs {
      ip_protocol = "udp"
      ports = [22]
    }
    src_ip_ranges = ["11.100.0.1/32", "10.0.0.0/24"]
  }
  target_resources = [google_compute_network.network1.self_link]
  target_service_accounts = [google_service_account.service_account.email, google_service_account.service_account2.email]
}
`, context)
}

func TestAccComputeFirewallPolicyRule_multipleRules(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"org_name":      fmt.Sprintf("organizations/%s", getTestOrgFromEnv(t)),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewallPolicyRule_multiple(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.rule1",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.rule2",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeFirewallPolicyRule_multipleAdd(context),
			},
			{
				ResourceName:      "google_compute_firewall_policy_rule.rule3",
				ImportState:       true,
				ImportStateVerify: true,
				// Referencing using ID causes import to fail
				ImportStateVerifyIgnore: []string{"firewall_policy"},
			},
			{
				Config: testAccComputeFirewallPolicyRule_multipleRemove(context),
			},
		},
	})
}

func testAccComputeFirewallPolicyRule_multiple(context map[string]interface{}) string {
	return Nprintf(`
resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.name
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_rule" "rule1" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 9000
  enable_logging = true
  action = "allow"
  direction = "EGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [80, 8080]
    }
    dest_ip_ranges = ["11.100.0.1/32"]
  }
}

resource "google_compute_firewall_policy_rule" "rule2" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 9001
  enable_logging = false
  action = "deny"
  direction = "INGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [80, 8080]
    }
    layer4_configs {
      ip_protocol = "all"
    }
    src_ip_ranges = ["11.100.0.1/32"]
  }
}
`, context)
}

func testAccComputeFirewallPolicyRule_multipleAdd(context map[string]interface{}) string {
	return Nprintf(`
resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.id
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Description Update"
}

resource "google_compute_firewall_policy_rule" "rule1" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 9000
  enable_logging = true
  action = "allow"
  direction = "EGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
    }
    dest_ip_ranges = ["11.100.0.1/32"]
  }
}

resource "google_compute_firewall_policy_rule" "rule2" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 9001
  enable_logging = false
  action = "deny"
  direction = "INGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [80, 8080]
    }
    layer4_configs {
      ip_protocol = "all"
    }
    src_ip_ranges = ["11.100.0.1/32"]
  }
}

resource "google_compute_firewall_policy_rule" "rule3" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 40
  enable_logging = true
  action = "allow"
  direction = "INGRESS"
  disabled = true
  match {
    layer4_configs {
      ip_protocol = "udp"
      ports = [8000]
    }
    src_ip_ranges = ["11.100.0.1/32", "10.0.0.0/24"]
  }
}
`, context)
}

func testAccComputeFirewallPolicyRule_multipleRemove(context map[string]interface{}) string {
	return Nprintf(`
resource "google_folder" "folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "%{org_name}"
}

resource "google_compute_firewall_policy" "default" {
  parent      = google_folder.folder.name
  short_name  = "tf-test-policy-%{random_suffix}"
  description = "Resource created for Terraform acceptance testing"
}

resource "google_compute_firewall_policy_rule" "rule1" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 9000
  enable_logging = true
  action = "allow"
  direction = "EGRESS"
  disabled = false
  match {
    layer4_configs {
      ip_protocol = "tcp"
      ports = [80, 8080]
    }
    dest_ip_ranges = ["11.100.0.1/32"]
  }
}

resource "google_compute_firewall_policy_rule" "rule3" {
  firewall_policy = google_compute_firewall_policy.default.id
  description = "Resource created for Terraform acceptance testing"
  priority = 40
  enable_logging = true
  action = "allow"
  direction = "INGRESS"
  disabled = true
  match {
    layer4_configs {
      ip_protocol = "udp"
      ports = [8000]
    }
    src_ip_ranges = ["11.100.0.1/32", "10.0.0.0/24"]
  }
}
`, context)
}
