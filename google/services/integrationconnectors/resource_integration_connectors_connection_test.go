// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package integrationconnectors_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since connection resources can't be created in parallel, we need to run the tests one by one.
func TestAccIntegrationConnectorsConnection(t *testing.T) {
	t.Parallel()
	testCases := map[string]func(t *testing.T){
		"basic":    testAccIntegrationConnectorsConnection_integrationConnectorsConnectionBasicResource,
		"advanced": testAccIntegrationConnectorsConnection_integrationConnectorsConnectionAdvancedResource,
		"sa":       testAccIntegrationConnectorsConnection_integrationConnectorsConnectionSaResource,
		"oauth":    testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthResource,
		"ssh":      testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthSshResource,
		"cc":       testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthCcResource,
		"jwt":      testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthJwtResource,
		"update":   testAccIntegrationConnectorsConnection_updateResource,
	}

	for name, tc := range testCases {
		// shadows the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionBasicResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_integrationConnectorsConnectionBasic(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.pubsubconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_integration_connectors_connection" "pubsubconnection" {
  name     = "tf-test-test-pubsub%{random_suffix}"
  location = "us-central1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/gcp/connectors/pubsub/versions/1"
  description = "tf created description"
  config_variable {
      key = "project_id"
      string_value = "connectors-example"
  }
  config_variable {
      key = "topic_id"
      string_value = "test"
  }
}
`, context)
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionAdvancedResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_integrationConnectorsConnectionAdvanced(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.zendeskconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionAdvanced(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id     = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "zendeskconnection" {
  name     = "tf-test-test-zendesk%{random_suffix}"
  description = "tf updated description"
  location = "us-central1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/zendesk/connectors/zendesk/versions/1"
  config_variable {
      key = "proxy_enabled"
      boolean_value = false
  }
  config_variable {
    key = "sample_integer_value"
    integer_value = 1
  }

  config_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
  }

  config_variable {
    key = "sample_secret_value"
    secret_value {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
  }

  suspended = false
  auth_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    auth_type = "USER_PASSWORD"
    auth_key = "sampleAuthKey"
    user_password {
      username = "user@xyz.com"
      password {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
  }

  destination_config {
    key = "url"
    destination {
        host = "https://test.zendesk.com"
        port = 80
    }
  }
  lock_config {
    locked = false
    reason = "Its not locked"
  }
  log_config {
    enabled = true
  }
  node_config {
    min_node_count = 2
    max_node_count = 50
  }
  labels = {
    foo = "bar"
  }
  ssl_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    client_cert_type = "PEM"
    client_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key_pass {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    private_server_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    server_cert_type = "PEM"
    trust_model      = "PRIVATE"
    type             = "TLS"
    use_ssl          = true
  }

  eventing_enablement_type = "EVENTING_AND_CONNECTION"
  eventing_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    registration_destination_config {
      key = "registration_destination_config"
      destination {
          host = "https://test.zendesk.com"
          port = 80
        }
    }
    auth_config {
      auth_type = "USER_PASSWORD"
      auth_key = "sampleAuthKey"
      user_password {
        username = "user@xyz.com"
        password {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_string"
        string_value = "sampleString"
      }
      additional_variable {
        key = "sample_boolean"
        boolean_value = false
      }
      additional_variable {
        key = "sample_integer"
        integer_value = 1
      }
      additional_variable {
        key = "sample_secret_value"
        secret_value {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "GOOGLE_MANAGED"
          kms_key_name = "sampleKMSKkey"
        }
      }
    }
    enrichment_enabled = true
  }
}
`, context)
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionSaResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_integrationConnectorsConnectionSa(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.zendeskconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionSa(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id     = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "zendeskconnection" {
  name     = "tf-test-test-zendesk%{random_suffix}"
  description = "tf updated description"
  location = "us-central1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/zendesk/connectors/zendesk/versions/1"
  config_variable {
    key = "proxy_enabled"
    boolean_value = false
  }
  config_variable {
    key = "sample_integer_value"
    integer_value = 1
  }

  config_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
  }

  config_variable {
    key = "sample_secret_value"
    secret_value {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
  }

  suspended = false
  auth_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    auth_type = "USER_PASSWORD"
    auth_key = "sampleAuthKey"
    user_password {
      username = "user@xyz.com"
      password {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
  }

  destination_config {
    key = "url"
    destination {
      service_attachment = "projects/connectors-example/regions/us-central1/serviceAttachments/test"
    }
  }
  lock_config {
    locked = false
    reason = "Its not locked"
  }
  log_config {
    enabled = true
  }
  node_config {
    min_node_count = 2
    max_node_count = 50
  }
  labels = {
    foo = "bar"
  }
  ssl_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    client_cert_type = "PEM"
    client_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    client_private_key_pass {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    private_server_certificate {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
    }
    server_cert_type = "PEM"
    trust_model      = "PRIVATE"
    type             = "TLS"
    use_ssl          = true
  }

  eventing_enablement_type = "EVENTING_AND_CONNECTION"
  eventing_config {
    additional_variable {
      key = "sample_string"
      string_value = "sampleString"
    }
    additional_variable {
      key = "sample_boolean"
      boolean_value = false
    }
    additional_variable {
      key = "sample_integer"
      integer_value = 1
    }
    additional_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }
    additional_variable {
      key = "sample_encryption_key_value"
      encryption_key_value {
        type = "GOOGLE_MANAGED"
        kms_key_name = "sampleKMSKkey"
      }
    }
    registration_destination_config {
      key = "registration_destination_config"
      destination {
          service_attachment = "projects/connectors-example/regions/us-central1/serviceAttachments/test"
        }
    }
    auth_config {
      auth_type = "USER_PASSWORD"
      auth_key = "sampleAuthKey"
      user_password {
        username = "user@xyz.com"
        password {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_string"
        string_value = "sampleString"
      }
      additional_variable {
        key = "sample_boolean"
        boolean_value = false
      }
      additional_variable {
        key = "sample_integer"
        integer_value = 1
      }
      additional_variable {
        key = "sample_secret_value"
        secret_value {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "GOOGLE_MANAGED"
          kms_key_name = "sampleKMSKkey"
        }
      }
    }
    enrichment_enabled = true
  }
}
`, context)
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauth(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.boxconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauth(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id     = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "boxconnection" {
  name     = "tf-test-test-box%{random_suffix}"
  location = "us-central1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/box/connectors/box/versions/1"
  description = "tf created description"
  config_variable {
      key = "impersonate_user_mode"
      string_value = "User"
  }
  config_variable {
      key = "proxy_enabled"
      boolean_value = false
  }
  auth_config{
    auth_type = "OAUTH2_AUTH_CODE_FLOW"
    oauth2_auth_code_flow {
        auth_uri  = "sampleauthuri"
        client_id = "sampleclientid"
        client_secret {
            secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
        enable_pkce   = true
        scopes = [
            "sample_scope_1",
            "sample_scope_2"
        ]
    }
  }

}
`, context)
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthSshResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthSsh(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.boxconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthSsh(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id     = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "boxconnection" {
  name     = "tf-test-test-box%{random_suffix}"
  location = "us-central1"
  connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/box/connectors/box/versions/1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  description = "tf created description"
  config_variable {
      key = "impersonate_user_mode"
      string_value = "User"
  }
  config_variable {
      key = "proxy_enabled"
      boolean_value = false
  }
  auth_config{
    auth_type = "SSH_PUBLIC_KEY"
    ssh_public_key {
      cert_type = "PEMKEY_BLOB"
      ssh_client_cert {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      ssh_client_cert_pass {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      username = "abc"
    }
  }

}
`, context)
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthCcResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthCc(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.boxconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthCc(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id     = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "boxconnection" {
  name     = "tf-test-test-box%{random_suffix}"
  location = "us-central1"
  connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/box/connectors/box/versions/1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  description = "tf created description"
  config_variable {
      key = "impersonate_user_mode"
      string_value = "User"
  }
  config_variable {
      key = "proxy_enabled"
      boolean_value = false
  }
  auth_config {
   auth_type = "OAUTH2_CLIENT_CREDENTIALS"
   oauth2_client_credentials {
     client_id = "testclientid"
     client_secret {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
     }
   }
   additional_variable {
     key = "oauth_jwt_cert"
     secret_value {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
     }
   }
   additional_variable {
     key = "oauth_jwt_cert_password"
     secret_value {
      secret_version = google_secret_manager_secret_version.secret-version-basic.name
     }
   }
   additional_variable {
     key = "oauth_jwt_subject_type"
     string_value = "sample"
   }
   additional_variable {
     key = "oauth_jwt_subject"
     string_value = "sample"
   }
   additional_variable {
     key = "oauth_jwt_public_key_id"
     string_value = "sample"
   }
   additional_variable {
     key = "auth_scheme"
     string_value = "sample"
   }
   additional_variable {
     key = "initiate_oauth"
     string_value = "sample"
   }
   additional_variable {
     key = "oauth_jwt_cert_type"
     string_value = "PEMKEY_BLOB"
   }
 }
}
`, context)
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthJwtResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthJwt(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.boxconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_integrationConnectorsConnectionOauthJwt(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_secret_manager_secret" "secret-basic" {
  secret_id     = "tf-test-test-secret%{random_suffix}"
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data = "dummypassword"
}

resource "google_secret_manager_secret_iam_member" "secret_iam" {
  secret_id  = google_secret_manager_secret.secret-basic.id
  role       = "roles/secretmanager.admin"
  member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_secret_manager_secret_version.secret-version-basic]
}

resource "google_integration_connectors_connection" "boxconnection" {
  name     = "tf-test-test-box%{random_suffix}"
  location = "us-central1"
  connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/box/connectors/box/versions/1"
  service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
  description = "tf created description"
  config_variable {
      key = "impersonate_user_mode"
      string_value = "User"
  }
  config_variable {
      key = "proxy_enabled"
      boolean_value = false
  }
  auth_config {
    auth_type = "OAUTH2_JWT_BEARER"
    oauth2_jwt_bearer {
      client_key {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      jwt_claims {
        issuer = "test"
        subject = "johndoe@example.org"
        audience  = "test"
      }
    }
  }
}
`, context)
}

func testAccIntegrationConnectorsConnection_updateResource(t *testing.T) {

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsConnection_full(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.zendeskconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
			{
				Config: testAccIntegrationConnectorsConnection_update(context),
			},
			{
				ResourceName:            "google_integration_connectors_connection.zendeskconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels", "status.0.description"},
			},
		},
	})
}

func testAccIntegrationConnectorsConnection_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
  data "google_project" "test_project" {
  }

  resource "google_secret_manager_secret" "secret-basic" {
    secret_id = "tf-test-test-secret%{random_suffix}"
    replication {
      user_managed {
        replicas {
          location = "us-central1"
        }
      }
    }
  }


  resource "google_secret_manager_secret_version" "secret-version-basic" {
    secret = google_secret_manager_secret.secret-basic.id
    secret_data = "dummypassword"
  }

  resource "google_secret_manager_secret_iam_member" "secret_iam" {
    secret_id  = google_secret_manager_secret.secret-basic.id
    role       = "roles/secretmanager.admin"
    member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
    depends_on = [google_secret_manager_secret_version.secret-version-basic]
  }

  resource "google_integration_connectors_connection" "zendeskconnection" {
    name     = "tf-test-test-zendesk%{random_suffix}"
    description = "tf description"
    location = "us-central1"
    service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
    connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/zendesk/connectors/zendesk/versions/1"
    config_variable {
        key = "proxy_enabled"
        boolean_value = false
    }
    config_variable {
      key = "sample_integer_value"
      integer_value = 1
    }

    config_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "GOOGLE_MANAGED"
          kms_key_name = "sampleKMSKkey"
        }
    }

    config_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }

    suspended = false
    auth_config {
      auth_type = "USER_PASSWORD"
      auth_key = "sampleAuthKey"
      user_password {
        username = "user@xyz.com"
        password {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
    }

    destination_config {
      key = "url"
      destination {
          host = "https://test.zendesk.com"
      }
    }
    lock_config {
      locked = false
      reason = "Its for sure not locked"
    }
    log_config {
      enabled = true
    }
    node_config {
      min_node_count = 2
      max_node_count = 50
    }
    labels = {
      foo = "bar"
    }
    ssl_config {
      additional_variable {
        key = "sample_string"
        string_value = "sampleString"
      }
      additional_variable {
        key = "sample_boolean"
        boolean_value = false
      }
      additional_variable {
        key = "sample_integer"
        integer_value = 1
      }
      additional_variable {
        key = "sample_secret_value"
        secret_value {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "GOOGLE_MANAGED"
          kms_key_name = "sampleKMSKkey"
        }
      }
      client_cert_type = "PEM"
      client_certificate {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      client_private_key {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      client_private_key_pass {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      private_server_certificate {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      server_cert_type = "PEM"
      trust_model      = "PRIVATE"
      type             = "TLS"
      use_ssl          = true
    }

    eventing_enablement_type = "EVENTING_AND_CONNECTION"
    eventing_config {
      additional_variable {
        key = "sample_string"
        string_value = "sampleString"
      }
      additional_variable {
        key = "sample_boolean"
        boolean_value = false
      }
      additional_variable {
        key = "sample_integer"
        integer_value = 1
      }
      additional_variable {
        key = "sample_secret_value"
        secret_value {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "GOOGLE_MANAGED"
          kms_key_name = "sampleKMSKkey"
        }
      }
      registration_destination_config {
        key = "registration_destination_config"
        destination {
          host = "https://test.zendesk.com"
        }
      }
      auth_config {
        auth_type = "USER_PASSWORD"
        auth_key = "sampleAuthKey"
        user_password {
          username = "user@xyz.com"
          password {
            secret_version = google_secret_manager_secret_version.secret-version-basic.name
          }
        }
        additional_variable {
          key = "sample_string"
          string_value = "sampleString"
        }
        additional_variable {
          key = "sample_boolean"
          boolean_value = false
        }
        additional_variable {
          key = "sample_integer"
          integer_value = 1
        }
        additional_variable {
          key = "sample_secret_value"
          secret_value {
            secret_version = google_secret_manager_secret_version.secret-version-basic.name
          }
        }
        additional_variable {
          key = "sample_encryption_key_value"
          encryption_key_value {
            type = "GOOGLE_MANAGED"
            kms_key_name = "sampleKMSKkey"
          }
        }
      }
      enrichment_enabled = true
    }
  }
`, context)
}

func testAccIntegrationConnectorsConnection_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
  data "google_project" "test_project" {
  }

  resource "google_secret_manager_secret" "secret-basic" {
    secret_id = "tf-test-test-secret%{random_suffix}"
    replication {
      user_managed {
        replicas {
          location = "us-central1"
        }
      }
    }
  }


  resource "google_secret_manager_secret_version" "secret-version-basic" {
    secret = google_secret_manager_secret.secret-basic.id
    secret_data = "dummypassword"
  }

  resource "google_secret_manager_secret_iam_member" "secret_iam" {
    secret_id  = google_secret_manager_secret.secret-basic.id
    role       = "roles/secretmanager.admin"
    member     = "serviceAccount:${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
    depends_on = [google_secret_manager_secret_version.secret-version-basic]
  }

  resource "google_integration_connectors_connection" "zendeskconnection" {
    name     = "tf-test-test-zendesk%{random_suffix}"
    description = "tf updated description"
    location = "us-central1"
    service_account = "${data.google_project.test_project.number}-compute@developer.gserviceaccount.com"
    connector_version = "projects/${data.google_project.test_project.project_id}/locations/global/providers/zendesk/connectors/zendesk/versions/1"
    config_variable {
        key = "proxy_enabled"
        boolean_value = true
    }
    config_variable {
      key = "sample_integer_value"
      integer_value = 2
    }

    config_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "CUSTOMER_MANAGED"
          kms_key_name = "sampleKMSKkey1"
        }
    }

    config_variable {
      key = "sample_secret_value"
      secret_value {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
    }

    suspended = false
    auth_config {
      auth_type = "USER_PASSWORD"
      auth_key = "sampleNewAuthKey"
      user_password {
        username = "user1@xyz.com"
        password {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
    }

    destination_config {
      key = "url"
      destination {
          host = "https://test1.zendesk.com"
      }
    }
    lock_config {
      locked = false
      reason = "Its for sure not locked"
    }
    log_config {
      enabled = true
    }
    node_config {
      min_node_count = 3
      max_node_count = 49
    }
    labels = {
      bar = "foo"
    }
    ssl_config {
      additional_variable {
        key = "sample_string"
        string_value = "sampleString1"
      }
      additional_variable {
        key = "sample_boolean"
        boolean_value = true
      }
      additional_variable {
        key = "sample_integer"
        integer_value = 2
      }
      additional_variable {
        key = "sample_secret_value"
        secret_value {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "CUSTOMER_MANAGED"
          kms_key_name = "sampleNewKMSKkey"
        }
      }
      client_cert_type = "PEM"
      client_certificate {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      client_private_key {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      client_private_key_pass {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      private_server_certificate {
        secret_version = google_secret_manager_secret_version.secret-version-basic.name
      }
      server_cert_type = "PEM"
      trust_model      = "INSECURE"
      type             = "MTLS"
      use_ssl          = false
    }

    eventing_enablement_type = "EVENTING_AND_CONNECTION"
    eventing_config {
      additional_variable {
        key = "sample_string"
        string_value = "sampleString1"
      }
      additional_variable {
        key = "sample_boolean"
        boolean_value = true
      }
      additional_variable {
        key = "sample_integer"
        integer_value = 2
      }
      additional_variable {
        key = "sample_secret_value"
        secret_value {
          secret_version = google_secret_manager_secret_version.secret-version-basic.name
        }
      }
      additional_variable {
        key = "sample_encryption_key_value"
        encryption_key_value {
          type = "CUSTOMER_MANAGED"
          kms_key_name = "sampleNewKMSKkey"
        }
      }
      registration_destination_config {
        key = "registration_destination_config"
        destination {
          host = "https://test1.zendesk.com"
        }
      }
      auth_config {
        auth_type = "USER_PASSWORD"
        auth_key = "sampleAuthKey1"
        user_password {
          username = "user1@xyz.com"
          password {
            secret_version = google_secret_manager_secret_version.secret-version-basic.name
          }
        }
        additional_variable {
          key = "sample_string"
          string_value = "sampleString1"
        }
        additional_variable {
          key = "sample_boolean"
          boolean_value = true
        }
        additional_variable {
          key = "sample_integer"
          integer_value = 2
        }
        additional_variable {
          key = "sample_secret_value"
          secret_value {
            secret_version = google_secret_manager_secret_version.secret-version-basic.name
          }
        }
        additional_variable {
          key = "sample_encryption_key_value"
          encryption_key_value {
            type = "CUSTOMER_MANAGED"
            kms_key_name = "sampleNewKMSKkey"
          }
        }
      }
      enrichment_enabled = false
    }
  }
`, context)
}

func testAccCheckIntegrationConnectorsConnectionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_integration_connectors_connection" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{IntegrationConnectorsBasePath}}projects/{{project}}/locations/{{location}}/connections/{{name}}")
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
				return fmt.Errorf("IntegrationConnectorsConnection still exists at %s", url)
			}
		}

		return nil
	}
}
