// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sourcerepo_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSourceRepoRepository_basic(t *testing.T) {
	t.Parallel()

	repositoryName := fmt.Sprintf("source-repo-repository-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSourceRepoRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceRepoRepository_basic(repositoryName),
			},
			{
				ResourceName:      "google_sourcerepo_repository.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSourceRepoRepository_update(t *testing.T) {
	t.Parallel()

	repositoryName := fmt.Sprintf("source-repo-repository-test-%s", acctest.RandString(t, 10))
	accountId := fmt.Sprintf("account-id-%s", acctest.RandString(t, 10))
	topicName := fmt.Sprintf("topic-name-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSourceRepoRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceRepoRepository_basic(repositoryName),
			},
			{
				ResourceName:      "google_sourcerepo_repository.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSourceRepoRepository_extended(accountId, topicName, repositoryName),
			},
			{
				ResourceName:      "google_sourcerepo_repository.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSourceRepoRepository_basic(repositoryName string) string {
	return fmt.Sprintf(`
resource "google_sourcerepo_repository" "acceptance" {
  name = "%s"
}
`, repositoryName)
}

func testAccSourceRepoRepository_extended(accountId string, topicName string, repositoryName string) string {
	return fmt.Sprintf(`
	resource "google_service_account" "test-account" {
		account_id   = "%s"
		display_name = "Test Service Account"
	  }
	  
	  resource "google_pubsub_topic" "topic" {
		name     = "%s"
	  }
	  
	  resource "google_sourcerepo_repository" "acceptance" {
		name = "%s"
		pubsub_configs {
			topic = google_pubsub_topic.topic.id
			message_format = "JSON"
			service_account_email = google_service_account.test-account.email
		}
	  }
`, accountId, topicName, repositoryName)
}

// Test setting create_ignore_already_exists on an existing resource
func TestAccSourceRepoRepository_existingResourceCreateIgnoreAlreadyExists(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	repositoryName := fmt.Sprintf("source-repo-repository-test-%s", acctest.RandString(t, 10))
	id := fmt.Sprintf("projects/%s/repos/%s", project, repositoryName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSourceRepoRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			// The first step creates a new resource with create_ignore_already_exists=false
			{
				Config: testAccSourceRepoRepositoryCreateIgnoreAlreadyExists(repositoryName, false),
				Check:  resource.TestCheckResourceAttr("google_sourcerepo_repository.acceptance", "id", id),
			},
			{
				ResourceName:            "google_sourcerepo_repository.acceptance",
				ImportStateId:           id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_ignore_already_exists"}, // Import leaves this field out when false
			},
			// The second step updates the resource to have create_ignore_already_exists=true
			{
				Config: testAccSourceRepoRepositoryCreateIgnoreAlreadyExists(repositoryName, true),
				Check:  resource.TestCheckResourceAttr("google_sourcerepo_repository.acceptance", "id", id),
			},
		},
	})
}

// Test the option to ignore ALREADY_EXISTS error from creating a Source Repository.
func TestAccSourceRepoRepository_createIgnoreAlreadyExists(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	repositoryName := fmt.Sprintf("source-repo-repository-test-%s", acctest.RandString(t, 10))
	id := fmt.Sprintf("projects/%s/repos/%s", project, repositoryName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSourceRepoRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			// The first step creates a basic Source Repository
			{
				Config: testAccSourceRepoRepository_basic(repositoryName),
				Check:  resource.TestCheckResourceAttr("google_sourcerepo_repository.acceptance", "id", id),
			},
			{
				ResourceName:      "google_sourcerepo_repository.acceptance",
				ImportStateId:     id,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step creates a new resource that duplicates with the existing Source Repository.
			{
				Config: testAccSourceRepoRepositoryDuplicateIgnoreAlreadyExists(repositoryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_sourcerepo_repository.acceptance", "id", id),
					resource.TestCheckResourceAttr("google_sourcerepo_repository.duplicate", "id", id),
				),
			},
		},
	})
}

func testAccSourceRepoRepositoryCreateIgnoreAlreadyExists(repositoryName string, ignore_already_exists bool) string {
	return fmt.Sprintf(`
resource "google_sourcerepo_repository" "acceptance" {
  name = "%s"
  create_ignore_already_exists = %t
}
`, repositoryName, ignore_already_exists)
}

func testAccSourceRepoRepositoryDuplicateIgnoreAlreadyExists(repositoryName string) string {
	return fmt.Sprintf(`
resource "google_sourcerepo_repository" "acceptance" {
  name = "%s"
}

resource "google_sourcerepo_repository" "duplicate" {
  name = "%s"
  create_ignore_already_exists = true
}
`, repositoryName, repositoryName)
}
