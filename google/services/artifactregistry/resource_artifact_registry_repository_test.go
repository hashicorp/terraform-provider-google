// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package artifactregistry_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccArtifactRegistryRepository_update(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccArtifactRegistryRepository_update(repositoryID),
			},
			{
				ResourceName:            "google_artifact_registry_repository.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccArtifactRegistryRepository_update2(repositoryID),
			},
			{
				ResourceName:            "google_artifact_registry_repository.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccArtifactRegistryRepository_createMvnSnapshot(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccArtifactRegistryRepository_createMvnWithVersionPolicy(repositoryID, "SNAPSHOT"),
			},
			{
				ResourceName:      "google_artifact_registry_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccArtifactRegistryRepository_createMvnRelease(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccArtifactRegistryRepository_createMvnWithVersionPolicy(repositoryID, "RELEASE"),
			},
			{
				ResourceName:      "google_artifact_registry_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccArtifactRegistryRepository_updateEmptyMvn(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccArtifactRegistryRepository_createMvnEmpty1(repositoryID),
			},
			{
				ResourceName:      "google_artifact_registry_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccArtifactRegistryRepository_createMvnEmpty2(repositoryID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config:             testAccArtifactRegistryRepository_createMvnEmpty3(repositoryID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccArtifactRegistryRepository_kfp(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccArtifactRegistryRepository_kfp(repositoryID),
			},
			{
				ResourceName:      "google_artifact_registry_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccArtifactRegistryRepository_update(repositoryID string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  repository_id = "%s"
  location = "us-central1"
  description = "pre-update"
  format = "DOCKER"

  labels = {
    my_key    = "my_val"
    other_key = "other_val"
  }
}
`, repositoryID)
}

func testAccArtifactRegistryRepository_update2(repositoryID string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  repository_id = "%s"
  location = "us-central1"
  description = "post-update"
  format = "DOCKER"

  labels = {
    my_key    = "my_val"
    other_key = "new_val"
  }
}
`, repositoryID)
}

func testAccArtifactRegistryRepository_createMvnWithVersionPolicy(repositoryID string, versionPolicy string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  repository_id = "%s"
  location = "us-central1"
  description = "post-update"
  format = "MAVEN"
  maven_config {
    version_policy = "%s"
  }
}
`, repositoryID, versionPolicy)
}

func testAccArtifactRegistryRepository_createMvnEmpty1(repositoryID string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  repository_id = "%s"
  location = "us-central1"
  description = "maven repository with empty maven_config"
  format = "MAVEN"
}
`, repositoryID)
}

func testAccArtifactRegistryRepository_createMvnEmpty2(repositoryID string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  repository_id = "%s"
  location = "us-central1"
  description = "maven repository with empty maven_config"
  format = "MAVEN"
  maven_config { }
}
`, repositoryID)
}

func testAccArtifactRegistryRepository_createMvnEmpty3(repositoryID string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  repository_id = "%s"
  location = "us-central1"
  description = "maven repository with empty maven_config"
  format = "MAVEN"
  maven_config {
    allow_snapshot_overwrites = false
  }
}
`, repositoryID)
}

func testAccArtifactRegistryRepository_kfp(repositoryID string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  repository_id = "%s"
  location = "us-central1"
  description = "my-kfp-repository"
  format = "KFP"
}
`, repositoryID)
}
