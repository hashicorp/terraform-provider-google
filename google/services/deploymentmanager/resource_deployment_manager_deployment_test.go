// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package deploymentmanager_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDeploymentManagerDeployment_basicFile(t *testing.T) {
	t.Parallel()

	randSuffix := acctest.RandString(t, 10)
	deploymentId := "tf-dm-" + randSuffix
	accountId := "tf-dm-account-" + randSuffix
	yamlPath := createYamlConfigFileForTest(t, "test-fixtures/service_account.yml.tmpl", map[string]interface{}{
		"account_id": accountId,
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckDeploymentManagerDeploymentDestroyProducer(t),
			testDeploymentManagerDeploymentVerifyServiceAccountMissing(t, accountId)),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentManagerDeployment_basicFile(deploymentId, yamlPath),
			},
			{
				ResourceName:            "google_deployment_manager_deployment.deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target", "create_policy", "delete_policy", "preview"},
			},
		},
	})
}

func TestAccDeploymentManagerDeployment_deleteInvalidOnCreate(t *testing.T) {
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	deploymentName := "tf-dm-" + randStr
	accountId := "tf-dm-" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDeploymentManagerDestroyInvalidDeployment(t, deploymentName),
		Steps: []resource.TestStep{
			{
				Config:      testAccDeploymentManagerDeployment_invalidCreatePolicy(deploymentName, accountId),
				ExpectError: regexp.MustCompile("BAD REQUEST"),
			},
		},
	})
}

func TestAccDeploymentManagerDeployment_createDeletePolicy(t *testing.T) {
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	deploymentName := "tf-dm-" + randStr
	accountId := "tf-dm-" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDeploymentManagerDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentManagerDeployment_createDeletePolicy(deploymentName, accountId),
			},
			{
				ResourceName:            "google_deployment_manager_deployment.deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target", "create_policy", "delete_policy", "preview"},
			},
		},
	})
}

func TestAccDeploymentManagerDeployment_imports(t *testing.T) {
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	deploymentName := "tf-dm-" + randStr
	accountId := "tf-dm-" + randStr
	importFilepath := createYamlConfigFileForTest(t, "test-fixtures/service_account.yml.tmpl", map[string]interface{}{
		"account_id": "{{ env['name'] }}",
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckDeploymentManagerDeploymentDestroyProducer(t),
			testDeploymentManagerDeploymentVerifyServiceAccountMissing(t, accountId)),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentManagerDeployment_imports(deploymentName, accountId, importFilepath),
				Check:  testDeploymentManagerDeploymentVerifyServiceAccountExists(t, accountId),
			},
			{
				ResourceName:            "google_deployment_manager_deployment.deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target", "create_policy", "delete_policy", "preview"},
			},
		},
	})
}

func TestAccDeploymentManagerDeployment_update(t *testing.T) {
	t.Parallel()

	randStr := acctest.RandString(t, 10)
	deploymentName := "tf-dm-" + randStr
	accountId := "tf-dm-first" + randStr
	accountId2 := "tf-dm-second" + randStr

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckDeploymentManagerDeploymentDestroyProducer(t),
			testDeploymentManagerDeploymentVerifyServiceAccountMissing(t, accountId)),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentManagerDeployment_preview(deploymentName, accountId),
				Check:  testDeploymentManagerDeploymentVerifyServiceAccountMissing(t, accountId),
			},
			{
				ResourceName:            "google_deployment_manager_deployment.deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target", "create_policy", "delete_policy", "preview"},
			},
			{
				Config: testAccDeploymentManagerDeployment_previewUpdated(deploymentName, accountId2),
				Check:  testDeploymentManagerDeploymentVerifyServiceAccountMissing(t, accountId2),
			},
			{
				ResourceName:            "google_deployment_manager_deployment.deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target", "create_policy", "delete_policy", "preview"},
			},
			{
				// Turn preview to false
				Config: testAccDeploymentManagerDeployment_deployed(deploymentName, accountId),
				Check:  testDeploymentManagerDeploymentVerifyServiceAccountExists(t, accountId),
			},
			{
				ResourceName:            "google_deployment_manager_deployment.deployment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target", "create_policy", "delete_policy", "preview"},
			},
		},
	})
}

func testAccDeploymentManagerDeployment_basicFile(deploymentName, yamlPath string) string {
	return fmt.Sprintf(`
resource "google_deployment_manager_deployment" "deployment" {
  name = "%s"

  target {
    config {
      content = file("%s")
    }
  }

  labels {
    key = "foo"
    value = "bar"
  }
}
`, deploymentName, yamlPath)
}

func testAccDeploymentManagerDeployment_invalidCreatePolicy(deployment, accountId string) string {
	// The service account doesn't exist, so create policy acquire fails
	return fmt.Sprintf(`
resource "google_deployment_manager_deployment" "deployment" {
  name = "%s"
  create_policy = "ACQUIRE"

  target {
    config {
      content = <<EOF
resources:
- name: %s
  type: iam.v1.serviceAccount
  properties:
    accountId: %s
    displayName: Test service account created by a DM Deployment, created in Terraform
EOF
    }
  }
}
`, deployment, accountId, accountId)
}

// NOTE: This is not recommended for use as actual Terraform config.
// This is just meant to test non-default createPolicy/deletePolicy parameters, but
// users shouldn't be managing resources in both Terraform and DM.
func testAccDeploymentManagerDeployment_createDeletePolicy(deployment, accountId string) string {
	return fmt.Sprintf(`
resource google_service_account "deployment_account" {
  account_id = "%s"
  display_name = "test account for Terraform DeploymentManager deployment"
}

resource "google_deployment_manager_deployment" "deployment" {
  name = "%s"

  // Deployment Manager will not create or delete resources
  create_policy = "ACQUIRE"
  delete_policy = "ABANDON"

  target {
    config {
      content = <<EOF
resources:
- name: "${google_service_account.deployment_account.account_id}"
  type: iam.v1.serviceAccount
  properties:
    accountId: "${google_service_account.deployment_account.account_id}"
    displayName: "${google_service_account.deployment_account.display_name}"
EOF
    }
  }
}
`, deployment, accountId)
}

func testAccDeploymentManagerDeployment_imports(deployment, accountId, importYamlPath string) string {
	return fmt.Sprintf(`
resource "google_deployment_manager_deployment" "deployment" {
  name = "%s"
  target {
    config {
      content = <<EOF
imports:
- path: service_account.jinja

resources:
- name: %s
  type: service_account.jinja
EOF
    }

    imports {
      name = "service_account.jinja"
      content = file("%s")
    }
  }
}
`, deployment, accountId, importYamlPath)
}

func testAccDeploymentManagerDeployment_preview(deployment, accountId string) string {
	return fmt.Sprintf(`
resource "google_deployment_manager_deployment" "deployment" {
  name = "%s"
  preview = true
  target {
    config {
      content = <<EOF
resources:
- name: %s
  type: iam.v1.serviceAccount
  properties:
    accountId: %s
    displayName: Test service account created by a DM Deployment, created in Terraform
EOF
    }
  }

  labels {
    key = "foo"
    value = "one"
  }
}
`, deployment, accountId, accountId)
}

func testAccDeploymentManagerDeployment_previewUpdated(deployment, accountId string) string {
	return fmt.Sprintf(`
resource "google_deployment_manager_deployment" "deployment" {
  name = "%s"
  preview = true
  target {
    config {
      content = <<EOF
resources:
- name: %s
  type: iam.v1.serviceAccount
  properties:
    accountId: %s
    displayName: Test service account created by a Terraform DeploymentManager Deployment
EOF
    }
  }

  labels {
    key = "foo"
    value = "one"
  }

  labels {
    key = "bar"
    value = "two"
  }
}
`, deployment, accountId, accountId)
}

func testAccDeploymentManagerDeployment_deployed(deployment, accountId string) string {
	return fmt.Sprintf(`
resource "google_deployment_manager_deployment" "deployment" {
  name = "%s"

  target {
    config {
      content = <<EOF
resources:
- name: %s
  type: iam.v1.serviceAccount
  properties:
    accountId: %s
    displayName: Test service account created by a DM Deployment, created in Terraform
EOF
    }
  }

  labels {
    key = "foo"
    value = "one"
  }
}
`, deployment, accountId, accountId)
}

func testDeploymentManagerDeploymentVerifyServiceAccountMissing(t *testing.T, accountId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		exists, err := testCheckDeploymentServiceAccountExists(accountId, config)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("service account %s found, should not exist", accountId)
		}
		return nil
	}
}

func testDeploymentManagerDeploymentVerifyServiceAccountExists(t *testing.T, accountId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		exists, err := testCheckDeploymentServiceAccountExists(accountId, config)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("service account %s not found", accountId)
		}
		return nil
	}
}

func testCheckDeploymentServiceAccountExists(accountId string, config *transport_tpg.Config) (exists bool, err error) {
	_, err = config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.Get(
		fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", envvar.GetTestProjectFromEnv(), accountId, envvar.GetTestProjectFromEnv())).Do()
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return false, nil
		}
		return false, fmt.Errorf("unexpected error while trying to confirm deployment service account %q exists: %v", accountId, err)
	}
	return true, nil
}

func testAccCheckDeploymentManagerDestroyInvalidDeployment(t *testing.T, deploymentName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type == "google_deployment_manager_deployment" {
				return fmt.Errorf("unexpected invalid deployment %q was saved in state", name)
			}
		}

		config := acctest.GoogleProviderConfig(t)
		url := fmt.Sprintf("%sprojects/%s/global/deployments/%s", config.DeploymentManagerBasePath, envvar.GetTestProjectFromEnv(), deploymentName)
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return fmt.Errorf("Unexpected error while trying to confirm DeploymentManagerDeployment deleted: %v", err)
		}
		if err == nil {
			return fmt.Errorf("DeploymentManagerDeployment still exists at %s", url)
		}
		return nil
	}
}

func testAccCheckDeploymentManagerDeploymentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_deployment_manager_deployment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DeploymentManagerBasePath}}projects/{{project}}/global/deployments/{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("DeploymentManagerDeployment still exists at %s", url)
			}
		}

		return nil
	}
}

func createYamlConfigFileForTest(t *testing.T, sourcePath string, context map[string]interface{}) string {
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	buf.WriteString(acctest.Nprintf(string(source), context))
	// Create temp file to write zip to
	tmpfile, err := ioutil.TempFile("", "*.yml")
	if err != nil {
		t.Fatal(err.Error())
	}
	if _, err := tmpfile.Write(buf.Bytes()); err != nil {
		t.Fatal(err.Error())
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err.Error())
	}
	return tmpfile.Name()
}
