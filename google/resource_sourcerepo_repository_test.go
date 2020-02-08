package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSourceRepoRepository_basic(t *testing.T) {
	t.Parallel()

	repositoryName := fmt.Sprintf("source-repo-repository-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceRepoRepositoryDestroy,
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

	repositoryName := fmt.Sprintf("source-repo-repository-test-%s", acctest.RandString(10))
	accountId := fmt.Sprintf("account-id-%s", acctest.RandString(10))
	topicName := fmt.Sprintf("topic-name-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceRepoRepositoryDestroy,
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
