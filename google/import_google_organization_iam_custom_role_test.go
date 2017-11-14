package google

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGoogleOrganizationIamCustomRole_import(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_ORG")
	roleId := "tfIamRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationIamCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganizationIamCustomRole_update(org, roleId),
			},
			{
				ResourceName:      "google_organization_iam_custom_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
