// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingBucketConfigFolder_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"folder_name":   "tf-test-" + acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigFolder_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_folder_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
			{
				Config: testAccLoggingBucketConfigFolder_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_folder_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
		},
	})
}

func TestAccLoggingBucketConfigProject_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_name":  "tf-test-" + acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_basic(context, 40),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccLoggingBucketConfigProject_analyticsEnabled(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_name":  "tf-test-" + acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_analyticsEnabled(context, true),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_analyticsEnabled(context, false),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccLoggingBucketConfigProject_locked(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_locked(context, false),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.variable_locked",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_locked(context, true),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.variable_locked",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccLoggingBucketConfigProject_cmekSettings(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	bucketId := fmt.Sprintf("tf-test-bucket-%s", acctest.RandString(t, 10))
	keyRingName := fmt.Sprintf("tf-test-key-ring-%s", acctest.RandString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-crypto-key-%s", acctest.RandString(t, 10))
	cryptoKeyNameUpdate := fmt.Sprintf("tf-test-crypto-key-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_cmekSettings(context, bucketId, keyRingName, cryptoKeyName, cryptoKeyNameUpdate),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_cmekSettingsUpdate(context, bucketId, keyRingName, cryptoKeyName, cryptoKeyNameUpdate),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccLoggingBucketConfigBillingAccount_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(t, 10),
		"billing_account_name": "billingAccounts/" + envvar.GetTestMasterBillingAccountFromEnv(t),
		"org_id":               envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigBillingAccount_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_billing_account_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account"},
			},
			{
				Config: testAccLoggingBucketConfigBillingAccount_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_billing_account_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account"},
			},
		},
	})
}

func TestAccLoggingBucketConfigOrganization_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigOrganization_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_organization_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
			{
				Config: testAccLoggingBucketConfigOrganization_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_organization_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
		},
	})
}

func testAccLoggingBucketConfigFolder_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_folder" "default" {
	display_name = "%{folder_name}"
	parent       = "organizations/%{org_id}"
}

resource "google_logging_folder_bucket_config" "basic" {
	folder    = google_folder.default.name
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingBucketConfigProject_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_project" "default" {
	project_id = "%{project_name}"
	name       = "%{project_name}"
	org_id     = "%{org_id}"
}

resource "google_logging_project_bucket_config" "basic" {
	project    = google_project.default.name
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingBucketConfigProject_analyticsEnabled(context map[string]interface{}, analytics bool) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_project" "default" {
	project_id = "%{project_name}"
	name       = "%{project_name}"
	org_id     = "%{org_id}"
}

resource "google_logging_project_bucket_config" "basic" {
	project    = google_project.default.name
	location  = "global"
	enable_analytics = %t
	bucket_id = "_Default"
}
`, context), analytics)
}

func testAccLoggingBucketConfigProject_locked(context map[string]interface{}, locked bool) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_project" "default" {
	project_id = "%{project_name}"
	name       = "%{project_name}"
	org_id     = "%{org_id}"
	billing_account = "%{billing_account}"
}

resource "google_logging_project_bucket_config" "fixed_locked" {
	project    = google_project.default.name
	location  = "global"
	locked = true
	bucket_id = "fixed-locked"
}

resource "google_logging_project_bucket_config" "variable_locked" {
	project    = google_project.default.name
	location  = "global"
	description = "lock status is %v" # test simultaneous update
	locked = %t
	bucket_id = "variable-locked"
}
`, context), locked, locked)
}

func testAccLoggingBucketConfigProject_preCmekSettings(context map[string]interface{}, keyRingName, cryptoKeyName, cryptoKeyNameUpdate string) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_project" "default" {
	project_id      = "%{project_name}"
	name            = "%{project_name}"
	org_id          = "%{org_id}"
	billing_account = "%{billing_account}"
}

resource "google_project_service" "logging_service" {
	project = google_project.default.project_id
	service = "logging.googleapis.com"
}

data "google_logging_project_cmek_settings" "cmek_settings" {
	project = google_project_service.logging_service.project
}

resource "google_kms_key_ring" "keyring" {
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "key1" {
	name            = "%s"
	key_ring        = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key" "key2" {
	name            = "%s"
	key_ring        = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_binding1" {
	crypto_key_id = google_kms_crypto_key.key1.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	
	members = [
		"serviceAccount:${data.google_logging_project_cmek_settings.cmek_settings.service_account_id}",
	]
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_binding2" {
	crypto_key_id = google_kms_crypto_key.key2.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	
	members = [
		"serviceAccount:${data.google_logging_project_cmek_settings.cmek_settings.service_account_id}",
	]
}
`, context), keyRingName, cryptoKeyName, cryptoKeyNameUpdate)
}

func testAccLoggingBucketConfigProject_cmekSettings(context map[string]interface{}, bucketId, keyRingName, cryptoKeyName, cryptoKeyNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "google_logging_project_bucket_config" "basic" {
	project        = google_project.default.name
	location       = "us-central1"
	retention_days = 30
	description    = "retention test 30 days"
	bucket_id      = "%s"

	cmek_settings {
		kms_key_name = google_kms_crypto_key.key1.id
	}

	depends_on   = [google_kms_crypto_key_iam_binding.crypto_key_binding1]
}
`, testAccLoggingBucketConfigProject_preCmekSettings(context, keyRingName, cryptoKeyName, cryptoKeyNameUpdate), bucketId)
}

func testAccLoggingBucketConfigProject_cmekSettingsUpdate(context map[string]interface{}, bucketId, keyRingName, cryptoKeyName, cryptoKeyNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "google_logging_project_bucket_config" "basic" {
	project        = google_project.default.name
	location       = "us-central1"
	retention_days = 30
	description    = "retention test 30 days"
	bucket_id      = "%s"

	cmek_settings {
		kms_key_name = google_kms_crypto_key.key2.id
	}

	depends_on   = [google_kms_crypto_key_iam_binding.crypto_key_binding2]
}
`, testAccLoggingBucketConfigProject_preCmekSettings(context, keyRingName, cryptoKeyName, cryptoKeyNameUpdate), bucketId)
}

func TestAccLoggingBucketConfig_CreateBuckets_withCustomId(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(t, 10),
		"billing_account_name": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":               envvar.GetTestOrgFromEnv(t),
		"project_name":         "tf-test-" + acctest.RandString(t, 10),
		"bucket_id":            "tf-test-bucket-" + acctest.RandString(t, 10),
	}

	configList := getLoggingBucketConfigs(context)

	for res, config := range configList {
		acctest.VcrTest(t, resource.TestCase{
			PreCheck:                 func() { acctest.AccTestPreCheck(t) },
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			Steps: []resource.TestStep{
				{
					Config: config,
				},
				{
					ResourceName:            fmt.Sprintf("google_logging_%s_bucket_config.basic", res),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{res},
				},
			},
		})
	}
}

func testAccLoggingBucketConfigBillingAccount_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`

data "google_billing_account" "default" {
	billing_account = "%{billing_account_name}"
}

resource "google_logging_billing_account_bucket_config" "basic" {
	billing_account    = data.google_billing_account.default.billing_account
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingBucketConfigOrganization_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`
data "google_organization" "default" {
	organization = "%{org_id}"
}

resource "google_logging_organization_bucket_config" "basic" {
	organization    = data.google_organization.default.organization
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func getLoggingBucketConfigs(context map[string]interface{}) map[string]string {
	return map[string]string{
		"project": acctest.Nprintf(`resource "google_project" "default" {
				project_id = "%{project_name}"
				name       = "%{project_name}"
				org_id     = "%{org_id}"
				billing_account = "%{billing_account_name}"
			}
			
			resource "google_logging_project_bucket_config" "basic" {
				project    = google_project.default.name
				location  = "global"
				retention_days = 10
				description = "retention test 10 days"
				bucket_id = "%{bucket_id}"
			}`, context),
	}

}
