package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"google.golang.org/api/iam/v1"
)

func TestAccCloudIdentityGroupMembership_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_domain":    getTestOrgDomainFromEnv(t),
		"cust_id":       getTestCustIdFromEnv(t),
		"identity_user": getTestIdentityUserFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
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
	return Nprintf(`
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
	return Nprintf(`
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

func TestAccCloudIdentityGroupMembership_import(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_domain":    getTestOrgDomainFromEnv(t),
		"cust_id":       getTestCustIdFromEnv(t),
		"identity_user": getTestIdentityUserFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
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
	return Nprintf(`
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

func TestAccCloudIdentityGroupMembership_membershipDoesNotExist(t *testing.T) {
	// Skip VCR because the service account needs to be created/deleted out of
	// band, and so those calls aren't recorded
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"org_domain":    getTestOrgDomainFromEnv(t),
		"cust_id":       getTestCustIdFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	saId := "tf-test-sa-" + randString(t, 10)
	project := getTestProjectFromEnv()
	config := BootstrapConfig(t)

	r := &iam.CreateServiceAccountRequest{
		AccountId:      saId,
		ServiceAccount: &iam.ServiceAccount{},
	}

	sa, err := config.NewIamClient(config.userAgent).Projects.ServiceAccounts.Create("projects/"+project, r).Do()
	if err != nil {
		t.Fatalf("Error creating service account: %s", err)
	}

	context["member_id"] = sa.Email

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIdentityGroupMembershipDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIdentityGroupMembership_dne(context),
			},
			{
				PreConfig: func() {
					config := googleProviderConfig(t)

					_, err := config.NewIamClient(config.userAgent).Projects.ServiceAccounts.Delete(sa.Name).Do()
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
	return Nprintf(`
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
