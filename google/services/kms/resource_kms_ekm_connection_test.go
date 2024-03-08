// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccKMSEkmConnection_kmsEkmConnectionBasicExample_full(context),
			},
			{
				ResourceName:            "google_kms_ekm_connection.example-ekmconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
			{
				Config: testAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(context),
			},
			{
				ResourceName:            "google_kms_ekm_connection.example-ekmconnection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
		},
	})
}

func testAccKMSEkmConnection_kmsEkmConnectionBasicExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_secret_manager_secret_version" "raw_der" {
  secret = "playground-cert"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "hostname" {
  secret = "external-uri"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "servicedirectoryservice" {
  secret = "external-servicedirectoryservice"
  project = "315636579862"
}
data "google_project" "vpc-project" {
  project_id = "cloud-ekm-refekm-playground"
}
data "google_project" "project" {
}
resource "google_project_iam_member" "add_sdviewer" {
  project = data.google_project.vpc-project.number
  role    = "roles/servicedirectory.viewer"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-ekms.iam.gserviceaccount.com"
}
resource "google_project_iam_member" "add_pscAuthorizedService" {
  project = data.google_project.vpc-project.number
  role    = "roles/servicedirectory.pscAuthorizedService"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-ekms.iam.gserviceaccount.com"
}
resource "google_kms_ekm_connection" "example-ekmconnection" {
  name            	= "tf_test_ekmconnection_example%{random_suffix}"
  location		= "us-central1"
  key_management_mode 	= "MANUAL"
  service_resolvers  	{
      service_directory_service  = data.google_secret_manager_secret_version.servicedirectoryservice.secret_data
      hostname 			 = data.google_secret_manager_secret_version.hostname.secret_data
      server_certificates        {
      		raw_der	= data.google_secret_manager_secret_version.raw_der.secret_data
      	}
    }
  depends_on = [
    	google_project_iam_member.add_pscAuthorizedService,
   	google_project_iam_member.add_sdviewer
   ]
}
`, context)
}

func testAccKMSEkmConnection_kmsEkmConnectionBasicExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "vpc-project" {
  project_id = "cloud-ekm-refekm-playground"
}
data "google_project" "project" {
}
data "google_secret_manager_secret_version" "raw_der" {
  secret = "playground-cert"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "hostname" {
  secret = "external-uri"
  project = "315636579862"
}
data "google_secret_manager_secret_version" "servicedirectoryservice" {
  secret = "external-servicedirectoryservice"
  project = "315636579862"
}
resource "google_project_iam_member" "add_sdviewer_updateekmconnection" {
  project = data.google_project.vpc-project.number
  role    = "roles/servicedirectory.viewer"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-ekms.iam.gserviceaccount.com"
}
resource "google_project_iam_member" "add_pscAuthorizedService_updateekmconnection" {
  project = data.google_project.vpc-project.number
  role    = "roles/servicedirectory.pscAuthorizedService"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-ekms.iam.gserviceaccount.com"
}
resource "google_kms_ekm_connection" "example-ekmconnection" {
  name            	= "tf_test_ekmconnection_example%{random_suffix}"
  location     		= "us-central1"
  key_management_mode 	= "CLOUD_KMS"
  crypto_space_path	= "v0/longlived/crypto-space-placeholder"
  service_resolvers  	{
      service_directory_service  = data.google_secret_manager_secret_version.servicedirectoryservice.secret_data
      hostname 			 = data.google_secret_manager_secret_version.hostname.secret_data
      server_certificates        {
      		raw_der	= data.google_secret_manager_secret_version.raw_der.secret_data
      	}
  }
  depends_on = [
    	google_project_iam_member.add_pscAuthorizedService_updateekmconnection,
   	google_project_iam_member.add_sdviewer_updateekmconnection
   ]
}
`, context)
}
