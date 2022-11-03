package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStorageTransferAgentPool_agentPoolUpdate(t *testing.T) {
	t.Parallel()

	agentPoolName := fmt.Sprintf("tf-test-agent-pool-%s", randString(t, 10))
	displayName := fmt.Sprintf("tf-test-display-name-%s", randString(t, 10))
	displayNameUpdate := fmt.Sprintf("tf-test-display-name-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStorageTransferAgentPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageTransferAgentPool_agentPoolBasic(getTestProjectFromEnv(), agentPoolName, displayName),
			},
			{
				ResourceName:      "google_storage_transfer_agent_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferAgentPool_agentPoolBasic(getTestProjectFromEnv(), agentPoolName, displayNameUpdate),
			},
			{
				ResourceName:      "google_storage_transfer_agent_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferAgentPool_updateLimitMbps(getTestProjectFromEnv(), agentPoolName, displayNameUpdate),
			},
			{
				ResourceName:      "google_storage_transfer_agent_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferAgentPool_omitDisplayName(getTestProjectFromEnv(), agentPoolName),
			},
			{
				ResourceName:      "google_storage_transfer_agent_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageTransferAgentPool_omitBandwidthLimit(getTestProjectFromEnv(), agentPoolName, displayNameUpdate),
			},
			{
				ResourceName:      "google_storage_transfer_agent_pool.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStorageTransferAgentPool_agentPoolBasic(project, agentPoolName, displayName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_project_iam_member" "agent_pool" {
  project = "%s"
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  display_name = "%s"
  bandwidth_limit {
    limit_mbps = "120"
  }

  depends_on = [google_project_iam_member.agent_pool]
}
`, project, project, agentPoolName, displayName)
}

func testAccStorageTransferAgentPool_updateLimitMbps(project, agentPoolName, displayName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_project_iam_member" "agent_pool" {
  project = "%s"
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  display_name = "%s"
  bandwidth_limit {
    limit_mbps = "150"
  }

  depends_on = [google_project_iam_member.agent_pool]
}
`, project, project, agentPoolName, displayName)
}

func testAccStorageTransferAgentPool_omitDisplayName(project string, agentPoolName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_project_iam_member" "agent_pool" {
  project = "%s"
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  bandwidth_limit {
    limit_mbps = "120"
  }

  depends_on = [google_project_iam_member.agent_pool]
}
`, project, project, agentPoolName)
}

func testAccStorageTransferAgentPool_omitBandwidthLimit(project string, agentPoolName string, displayName string) string {
	return fmt.Sprintf(`
data "google_storage_transfer_project_service_account" "default" {
  project = "%s"
}

resource "google_project_iam_member" "agent_pool" {
  project = "%s"
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${data.google_storage_transfer_project_service_account.default.email}"
}

resource "google_storage_transfer_agent_pool" "foo" {
  name         = "%s"
  display_name = "%s"

  depends_on = [google_project_iam_member.agent_pool]
}
`, project, project, agentPoolName, displayName)
}

func testAccCheckStorageTransferAgentPoolDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_transfer_agent_pool" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{StorageTransferBasePath}}projects/{{project}}/agentPools/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("StorageTransferAgentPool still exists at %s", url)
			}
		}

		return nil
	}
}
