package google

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGoogleProjectIamRole_import(t *testing.T) {
	t.Parallel()

	roleId := "tfIamRole" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectIamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectIamRole_update(roleId),
			},
			{
				ResourceName:      "google_project_iam_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
