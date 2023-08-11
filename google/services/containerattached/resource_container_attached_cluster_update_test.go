// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package containerattached_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccContainerAttachedCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContainerAttachedClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerAttachedCluster_containerAttachedCluster_full(context),
			},
			{
				ResourceName:            "google_container_attached_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAttachedCluster_containerAttachedCluster_update(context),
			},
			{
				ResourceName:            "google_container_attached_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContainerAttachedCluster_containerAttachedCluster_destroy(context),
			},
			{
				ResourceName:            "google_container_attached_cluster.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccContainerAttachedCluster_containerAttachedCluster_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_attached_versions" "versions" {
	location       = "us-west1"
	project        = data.google_project.project.project_id
}

resource "google_container_attached_cluster" "primary" {
  name     = "update%{random_suffix}"
  project = data.google_project.project.project_id
  location = "us-west1"
  description = "Test cluster"
  distribution = "aks"
  annotations = {
    label-one = "value-one"
  }
  authorization {
    admin_users = [ "user1@example.com", "user2@example.com"]
  }
  oidc_config {
      issuer_url = "https://oidc.issuer.url"
      jwks = base64encode("{\"keys\":[{\"use\":\"sig\",\"kty\":\"RSA\",\"kid\":\"testid\",\"alg\":\"RS256\",\"n\":\"somedata\",\"e\":\"AQAB\"}]}")
  }
  platform_version = data.google_container_attached_versions.versions.valid_versions[0]
  fleet {
      project = "projects/${data.google_project.project.number}"
  }
  logging_config {
    component_config {
      enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    }
  }
  monitoring_config {
    managed_prometheus_config {
      enabled = true
    }
  }
}
`, context)
}

func testAccContainerAttachedCluster_containerAttachedCluster_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_attached_versions" "versions" {
	location       = "us-west1"
	project        = data.google_project.project.project_id
}

resource "google_container_attached_cluster" "primary" {
  name     = "update%{random_suffix}"
  project = data.google_project.project.project_id
  location = "us-west1"
  description = "Test cluster updated"
  distribution = "aks"
  annotations = {
    label-one = "value-one"
  label-two = "value-two"
  }
  authorization {
    admin_users = [ "user2@example.com", "user3@example.com"]
  }
  oidc_config {
      issuer_url = "https://oidc.issuer.url"
      jwks = base64encode("{\"keys\":[{\"use\":\"sig\",\"kty\":\"RSA\",\"kid\":\"testid\",\"alg\":\"RS256\",\"n\":\"somedata\",\"e\":\"AQAB\"}]}")
  }
  platform_version = data.google_container_attached_versions.versions.valid_versions[0]
  fleet {
    project = "projects/${data.google_project.project.number}"
  }
  monitoring_config {
    managed_prometheus_config {}
  }
  lifecycle {
    prevent_destroy = true
  }
}
`, context)
}

// Duplicate of testAccContainerAttachedCluster_containerAttachedCluster_update without lifecycle.prevent_destroy set
// so the test can clean up the resource after the update.
func testAccContainerAttachedCluster_containerAttachedCluster_destroy(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

data "google_container_attached_versions" "versions" {
	location       = "us-west1"
	project        = data.google_project.project.project_id
}

resource "google_container_attached_cluster" "primary" {
  name     = "update%{random_suffix}"
  project = data.google_project.project.project_id
  location = "us-west1"
  description = "Test cluster updated"
  distribution = "aks"
  annotations = {
    label-one = "value-one"
  label-two = "value-two"
  }
  authorization {
    admin_users = [ "user2@example.com", "user3@example.com"]
  }
  oidc_config {
      issuer_url = "https://oidc.issuer.url"
      jwks = base64encode("{\"keys\":[{\"use\":\"sig\",\"kty\":\"RSA\",\"kid\":\"testid\",\"alg\":\"RS256\",\"n\":\"somedata\",\"e\":\"AQAB\"}]}")
  }
  platform_version = data.google_container_attached_versions.versions.valid_versions[0]
  fleet {
    project = "projects/${data.google_project.project.number}"
  }
  monitoring_config {
    managed_prometheus_config {}
  }
}
`, context)
}
