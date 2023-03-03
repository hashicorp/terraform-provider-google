package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccArtifactRegistryRepository_update(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccArtifactRegistryRepository_update(repositoryID),
			},
			{
				ResourceName:      "google_artifact_registry_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccArtifactRegistryRepository_update2(repositoryID),
			},
			{
				ResourceName:      "google_artifact_registry_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccArtifactRegistryRepository_createMvnSnapshot(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
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

	repositoryID := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
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

func TestAccArtifactRegistryRepository_kfp(t *testing.T) {
	t.Parallel()

	repositoryID := fmt.Sprintf("tf-test-%d", RandInt(t))

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckArtifactRegistryRepositoryDestroyProducer(t),
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
