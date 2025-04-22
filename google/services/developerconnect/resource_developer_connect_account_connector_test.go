// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package developerconnect_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDeveloperConnectAccountConnector_developerConnectAccountConnectorGithubUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectAccountConnector_Github(context),
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
			{
				Config: testAccDeveloperConnectAccountConnector_GithubUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_developer_connect_account_connector.my-account-connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
		},
	})
}

func testAccDeveloperConnectAccountConnector_Github(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"

  provider_oauth_config {
    system_provider_id = "GITHUB"
    scopes = ["repo"]
  }
}
`, context)
}

func testAccDeveloperConnectAccountConnector_GithubUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"
  annotations = {
    "foo": "bar"
  }
  labels = {
    "bar": "foo"
  }

  provider_oauth_config {
    system_provider_id = "GITHUB"
    scopes = ["repo", "public_repo"]
  }
}
`, context)
}

func TestAccDeveloperConnectAccountConnector_developerConnectAccountConnectorGitlabUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectAccountConnector_Gitlab(context),
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
			{
				Config: testAccDeveloperConnectAccountConnector_GitlabUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_developer_connect_account_connector.my-account-connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
		},
	})
}

func testAccDeveloperConnectAccountConnector_Gitlab(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"

  provider_oauth_config {
    system_provider_id = "GITLAB"
    scopes = ["api"]
  }
}
`, context)
}

func testAccDeveloperConnectAccountConnector_GitlabUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"
  
  annotations = {
    "foo": "bar"
  }
  
  labels = {
    "bar": "foo"
  }
  
  provider_oauth_config {
    system_provider_id = "GITLAB"
    scopes = ["api", "read_api"]
  }
}
`, context)
}

func TestAccDeveloperConnectAccountConnector_developerConnectAccountConnectorGoogleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectAccountConnector_Google(context),
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
			{
				Config: testAccDeveloperConnectAccountConnector_GoogleUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_developer_connect_account_connector.my-account-connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
		},
	})
}

func testAccDeveloperConnectAccountConnector_Google(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"

  provider_oauth_config {
    system_provider_id = "GOOGLE"
    scopes = ["https://www.googleapis.com/auth/drive.readonly"]
  }
}
`, context)
}

func testAccDeveloperConnectAccountConnector_GoogleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"
  
  annotations = {
    "foo": "bar"
  }
  
  labels = {
    "bar": "foo"
  }

  provider_oauth_config {
    system_provider_id = "GOOGLE"
    scopes = ["https://www.googleapis.com/auth/drive.readonly", "https://www.googleapis.com/auth/documents.readonly"]
  }
}
`, context)
}

func TestAccDeveloperConnectAccountConnector_developerConnectAccountConnectorSentryUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectAccountConnector_Sentry(context),
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
			{
				Config: testAccDeveloperConnectAccountConnector_SentryUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_developer_connect_account_connector.my-account-connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
		},
	})
}

func testAccDeveloperConnectAccountConnector_Sentry(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"

  provider_oauth_config {
    system_provider_id = "SENTRY"
    scopes = ["org:read"]
  }
}
`, context)
}

func testAccDeveloperConnectAccountConnector_SentryUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"
  
  annotations = {
    "foo": "bar"
  }
  
  labels = {
    "bar": "foo"
  }

  provider_oauth_config {
    system_provider_id = "SENTRY"
    scopes = ["org:read", "org:write"]
  }
}
`, context)
}

func TestAccDeveloperConnectAccountConnector_developerConnectAccountConnectorRovoUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectAccountConnector_Rovo(context),
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
			{
				Config: testAccDeveloperConnectAccountConnector_RovoUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_developer_connect_account_connector.my-account-connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
		},
	})
}

func testAccDeveloperConnectAccountConnector_Rovo(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"

  provider_oauth_config {
    system_provider_id = "ROVO"
    scopes = ["rovo"]
  }
}
`, context)
}

func testAccDeveloperConnectAccountConnector_RovoUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"
  
  annotations = {
     "foo": "bar"
  }
  
  labels = {
    "bar": "foo"
  }

  provider_oauth_config {
    system_provider_id = "ROVO"
    scopes = ["rovo"]
  }
}
`, context)
}

func TestAccDeveloperConnectAccountConnector_developerConnectAccountConnectorNewRelicUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectAccountConnector_NewRelic(context),
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
			{
				Config: testAccDeveloperConnectAccountConnector_NewRelicUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_developer_connect_account_connector.my-account-connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
		},
	})
}

func testAccDeveloperConnectAccountConnector_NewRelic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"

  provider_oauth_config {
    system_provider_id = "NEW_RELIC"
    scopes = []
  }
}
`, context)
}

func testAccDeveloperConnectAccountConnector_NewRelicUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"
  
  annotations = {
     "foo": "bar"
  }
  
  labels = {
    "bar": "foo"
  }

  provider_oauth_config {
    system_provider_id = "NEW_RELIC"
    scopes = []
  }
}
`, context)
}

func TestAccDeveloperConnectAccountConnector_developerConnectAccountConnectorDatastaxUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeveloperConnectAccountConnector_Datastax(context),
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
			{
				Config: testAccDeveloperConnectAccountConnector_DatastaxUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_developer_connect_account_connector.my-account-connector", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_developer_connect_account_connector.my-account-connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_connector_id", "annotations", "labels"},
			},
		},
	})
}

func testAccDeveloperConnectAccountConnector_Datastax(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"

  provider_oauth_config {
    system_provider_id = "DATASTAX"
    scopes = []
  }
}
`, context)
}

func testAccDeveloperConnectAccountConnector_DatastaxUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_developer_connect_account_connector" "my-account-connector" {
  location = "us-central1"
  account_connector_id = "tf-test-ac%{random_suffix}"
  
  annotations = {
     "foo": "bar"
  }
  
  labels = {
    "bar": "foo"
  }

  provider_oauth_config {
    system_provider_id = "DATASTAX"
    scopes = []
  }
}
`, context)
}
