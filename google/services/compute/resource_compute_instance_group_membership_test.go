// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeInstanceGroupMembership_instanceGroupMembershipBasic(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix": suffix,
		"zone":          envvar.GetTestZoneFromEnv(),
	}

	igId := fmt.Sprintf("projects/%s/zones/%s/instanceGroups/instance-group-%s",
		envvar.GetTestProjectFromEnv(), envvar.GetTestZoneFromEnv(), context["random_suffix"])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Create one membership resource
				Config: testAccComputeInstanceGroupMembership_instanceGroupMembershipBasic(context),
			},
			{
				ResourceName:      "google_compute_instance_group_membership.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Add two new members
				Config: testAccComputeInstanceGroupMembership_instanceGroupMembershipAdditional(context),
			},
			{
				ResourceName:            "google_compute_instance_group_membership.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "zone", "instance_group"},
			},
			{
				ResourceName:            "google_compute_instance_group_membership.add1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "zone", "instance_group"},
			},
			{
				ResourceName:            "google_compute_instance_group_membership.add2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "zone", "instance_group"},
			},
			{
				// Remove add1 and add2 membership resources
				Config: testAccComputeInstanceGroupMembership_instanceGroupMembershipBasic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceGroupMembershipDestroyed(
						t, igId,
						testAccComputeInstanceGroupMembershipGetInstanceName("add1-instance", suffix),
						testAccComputeInstanceGroupMembershipGetInstanceName("add2-instance", suffix),
					),
				),
			},
			{
				ResourceName:            "google_compute_instance_group_membership.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "zone", "instance_group"},
			},
			{
				// Delete all membership resources
				Config: testAccComputeInstanceGroupMembership_noInstanceGroupMembership(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceGroupMembershipDestroyed(
						t, igId,
						testAccComputeInstanceGroupMembershipGetInstanceName("default-instance", suffix)),
				),
			},
		},
	})
}

func testAccComputeInstanceGroupMembership_instanceGroupMembershipBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
    resource "google_compute_instance_group_membership" "default" {
      zone = "%{zone}"
      instance_group = google_compute_instance_group.default.name
      instance = google_compute_instance.default.self_link
    }
  `, context) + testAccComputeInstanceGroupMembership_noInstanceGroupMembership(context)
}

func testAccComputeInstanceGroupMembership_instanceGroupMembershipAdditional(context map[string]interface{}) string {
	return acctest.Nprintf(`
    resource "google_compute_instance_group_membership" "add1" {
      instance_group = google_compute_instance_group.default.name
      instance = google_compute_instance.add1.self_link
    }

    resource "google_compute_instance_group_membership" "add2" {
      instance_group = google_compute_instance_group.default.name
      instance = google_compute_instance.add2.self_link
    }
  `, context) + testAccComputeInstanceGroupMembership_instanceGroupMembershipBasic(context)
}

func testAccComputeInstanceGroupMembership_noInstanceGroupMembership(context map[string]interface{}) string {
	return acctest.Nprintf(`
    resource "google_compute_network" "default-network" {
      name = "default-%{random_suffix}"
    }

    resource "google_compute_instance" "default" {
      name         = "default-instance-%{random_suffix}"
      machine_type = "e2-medium"

      boot_disk {
        initialize_params {
          image = "debian-cloud/debian-11"
        }
      }

      network_interface {
        network = google_compute_network.default-network.name
      }
    }

    resource "google_compute_instance" "add1" {
      name         = "add1-instance-%{random_suffix}"
      machine_type = "e2-medium"

      boot_disk {
        initialize_params {
          image = "debian-cloud/debian-11"
        }
      }

      network_interface {
        network = google_compute_network.default-network.name
      }
    }

    resource "google_compute_instance" "add2" {
      name         = "add2-instance-%{random_suffix}"
      machine_type = "e2-medium"

      boot_disk {
        initialize_params {
          image = "debian-cloud/debian-11"
        }
      }

      network_interface {
        network = google_compute_network.default-network.name
      }
    }

    resource "google_compute_instance_group" "default" {
      name      = "instance-group-%{random_suffix}"
    }
  `, context)
}

func testAccCheckComputeInstanceGroupMembershipDestroyed(t *testing.T, instanceGroupId string, instances ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		foundInstances, err := testAccComputeInstanceGroupMembershipListMembership(t, instanceGroupId)
		if err != nil {
			return fmt.Errorf("unable to confirm instance group members with instances %+v was destroyed: %v", instances, err)
		}
		for _, p := range instances {
			if _, ok := foundInstances[p]; ok {
				return fmt.Errorf("instance group with instance %s still exists", p)
			}
		}
		return nil
	}
}

func testAccComputeInstanceGroupMembershipListMembership(t *testing.T, instanceGroupId string) (map[string]struct{}, error) {
	config := acctest.GoogleProviderConfig(t)

	url := fmt.Sprintf("https://www.googleapis.com/compute/v1/%s/listInstances", instanceGroupId)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		RawURL:    url,
		UserAgent: config.UserAgent,
	})

	if err != nil {
		return nil, err
	}

	v, ok := res["items"]
	if !ok || v == nil {
		return nil, nil
	}

	items := v.([]interface{})
	instances := make(map[string]struct{})
	for _, item := range items {
		instanceWithStatus := item.(map[string]interface{})
		v, ok := instanceWithStatus["instance"]
		if !ok || v == nil {
			continue
		}
		instance := v.(string)
		instances[fmt.Sprintf("%v", instance)] = struct{}{}
	}
	return instances, nil
}

func testAccComputeInstanceGroupMembershipGetInstanceName(instanceName string, suffix string) string {
	return fmt.Sprintf("projects/%s/zones/%s/instances/%s-%s",
		envvar.GetTestProjectFromEnv(),
		envvar.GetTestZoneFromEnv(),
		instanceName,
		suffix)
}
