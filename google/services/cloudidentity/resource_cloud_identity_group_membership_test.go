// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/iam/v1"
)

func testAccCloudIdentityGroupMembership_updateTest(t *testing.T) {
	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"identity_user": envvar.GetTestIdentityUserFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembership_update1(context),
			},
			{
				ResourceName:      "google_cloud_identity_group_membership.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIdentityGroupMembership_update2(context),
			},
			{
				ResourceName:      "google_cloud_identity_group_membership.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIdentityGroupMembership_update1(context),
			},
			{
				ResourceName:      "google_cloud_identity_group_membership.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudIdentityGroupMembership_update1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group_membership" "basic" {
  group    = google_cloud_identity_group.group.id

  preferred_member_key {
    id = "%{identity_user}@%{org_domain}"
  }

  roles {
    name = "MEMBER"
  }

}
`, context)
}

func testAccCloudIdentityGroupMembership_update2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group_membership" "basic" {
  group    = google_cloud_identity_group.group.id

  preferred_member_key {
    id = "%{identity_user}@%{org_domain}"
  }

  roles {
    name = "MEMBER"
  }

  roles {
    name = "MANAGER"
  }
}
`, context)
}

func testAccCloudIdentityGroupMembership_importTest(t *testing.T) {
	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"identity_user": envvar.GetTestIdentityUserFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembership_import(context),
			},
			{
				ResourceName:      "google_cloud_identity_group_membership.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudIdentityGroupMembership_import(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group_membership" "basic" {
  group    = google_cloud_identity_group.group.id

  preferred_member_key {
    id = "%{identity_user}@%{org_domain}"
  }

  roles {
    name = "MEMBER"
  }

  roles {
    name = "MANAGER"
  }
}
`, context)
}

func testAccCloudIdentityGroupMembership_membershipDoesNotExistTest(t *testing.T) {
	// Skip VCR because the service account needs to be created/deleted out of
	// band, and so those calls aren't recorded
	acctest.SkipIfVcr(t)

	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	saId := "tf-test-sa-" + acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()
	config := acctest.BootstrapConfig(t)

	r := &iam.CreateServiceAccountRequest{
		AccountId:      saId,
		ServiceAccount: &iam.ServiceAccount{},
	}

	sa, err := config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.Create("projects/"+project, r).Do()
	if err != nil {
		t.Fatalf("Error creating service account: %s", err)
	}

	context["member_id"] = sa.Email

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembership_dne(context),
			},
			{
				PreConfig: func() {
					config := acctest.GoogleProviderConfig(t)

					_, err := config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.Delete(sa.Name).Do()
					if err != nil {
						t.Errorf("cannot delete service account %s: %v", sa.Name, err)
						return
					}
				},
				Config:             testAccCloudIdentityGroupMembership_dne(context),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCloudIdentityGroupMembership_dne(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group-%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group-%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group_membership" "basic" {
  group = google_cloud_identity_group.group.id

  preferred_member_key {
    id = "%{member_id}"
  }

  roles {
    name = "MEMBER"
  }
}
`, context)
}

func testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipExampleTest(t *testing.T) {
	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipExample(context),
			},
			{
				ResourceName:            "google_cloud_identity_group_membership.cloud_identity_group_membership_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group"},
			},
		},
	})
}

func testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
  	id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group" "child-group" {
  display_name = "tf-test-my-identity-group%{random_suffix}-child"

  parent = "customers/%{cust_id}"

  group_key {
  	id = "tf-test-my-identity-group%{random_suffix}-child@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group_membership" "cloud_identity_group_membership_basic" {
  group    = google_cloud_identity_group.group.id

  preferred_member_key {
    id = google_cloud_identity_group.child-group.group_key[0].id
  }

  roles {
  	name = "MEMBER"
  }
}
`, context)
}

func testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipUserExampleTest(t *testing.T) {
	context := map[string]interface{}{
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"identity_user": envvar.GetTestIdentityUserFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipUserExample(context),
			},
			{
				ResourceName:            "google_cloud_identity_group_membership.cloud_identity_group_membership_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group"},
			},
		},
	})
}

func testAccCloudIdentityGroupMembership_cloudIdentityGroupMembershipUserExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group_membership" "cloud_identity_group_membership_basic" {
  group    = google_cloud_identity_group.group.id

  preferred_member_key {
    id = "%{identity_user}@%{org_domain}"
  }

  roles {
    name = "MEMBER"
  }

  roles {
    name = "MANAGER"
  }
}
`, context)
}

func testAccCheckCloudIdentityGroupMembershipDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_cloud_identity_group_membership" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{CloudIdentityBasePath}}{{name}}")
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
				return fmt.Errorf("CloudIdentityGroupMembership still exists at %s", url)
			}
		}

		return nil
	}
}
