package google

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	storagev1 "google.golang.org/api/storage/v1"
)

func testObjectIamName() string {
	return fmt.Sprintf("%s-%d", "tf-test-iam-Object", acctest.RandInt())
}

func TestAccGoogleStorageObjectIAMPolicy_basic(t *testing.T) {
	bucketName := testBucketIamName()
	objectName := testObjectIamName()
	skipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")

	policy := &storagev1.Policy{
		Bindings: []*storagev1.PolicyBindings{
			{
				Role: "roles/viewer",
				Members: []string{
					"user:admin@hashicorptest.com",
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleStorageObjectIAMPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageObjectsIAMPolicyBasic(bucketName, objectName, policy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectIAMPolicy(bucketName, objectName, policy),
				),
			},
		},
	})
}

func TestAccGoogleStorageObjectIAMPolicy_upgrade(t *testing.T) {
	bucketName := testBucketIamName()
	objectName := testObjectIamName()
	skipIfEnvNotSet(t, "GOOGLE_PROJECT_NUMBER")

	policy1 := &storagev1.Policy{
		Bindings: []*storagev1.PolicyBindings{
			{
				Role: "roles/viewer",
				Members: []string{
					"user:admin@hashicorptest.com",
				},
			},
		},
	}
	policy2 := &storagev1.Policy{
		Bindings: []*storagev1.PolicyBindings{
			{
				Role: "roles/editor",
				Members: []string{
					"user:admin@hashicorptest.com",
				},
			},
			{
				Role: "roles/viewer",
				Members: []string{
					"user:admin@hashicorptest.com",
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleStorageObjectIAMPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageObjectsIAMPolicyBasic(bucketName, objectName, policy1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectIAMPolicy(bucketName, objectName, policy1),
				),
			},

			resource.TestStep{
				Config: testGoogleStorageObjectsIAMPolicyBasic(bucketName, objectName, policy2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectIAMPolicy(bucketName, objectName, policy2),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageObjectIAMPolicy(bucket string, object string, policy *storagev1.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[object]
		if !ok {
			return fmt.Errorf("Not found: %s", object)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		p, err := config.clientStorage.Objects.GetIamPolicy(bucket, rs.Primary.ID).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of IAMPolicy for Object %s: %s", object, err)
		}

		if !reflect.DeepEqual(p.Bindings, policy.Bindings) {
			return fmt.Errorf("Incorrect iam policy bindings. Expected '%s', got '%s'", policy.Bindings, p.Bindings)
		}

		if _, ok = rs.Primary.Attributes["etag"]; !ok {
			return fmt.Errorf("Etag should be set.")
		}

		if rs.Primary.Attributes["etag"] != p.Etag {
			return fmt.Errorf("Incorrect etag value. Expected '%s', got '%s'", p.Etag, rs.Primary.Attributes["etag"])
		}

		return nil
	}
}

func testAccGoogleStorageObjectIAMPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_Object_iam_policy" {
			continue
		}

		object := rs.Primary.Attributes["object"]
		bucket := rs.Primary.Attributes["bucket"]

		policy, err := config.clientStorage.Objects.GetIamPolicy(bucket, object).Do()

		if err != nil && len(policy.Bindings) > 0 {
			return fmt.Errorf("Object %s policy hasn't been deleted", object)
		}
	}

	return nil
}

func testGoogleStorageObjectsIAMPolicyBasic(bucketName string, objectName string, policy *storagev1.Policy) string {
	var bindingBuffer bytes.Buffer

	for _, binding := range policy.Bindings {
		bindingBuffer.WriteString("binding {\n")
		bindingBuffer.WriteString(fmt.Sprintf("role = \"%s\"\n", binding.Role))
		bindingBuffer.WriteString(fmt.Sprintf("members = [\n"))
		for _, member := range binding.Members {
			bindingBuffer.WriteString(fmt.Sprintf("\"%s\",\n", member))
		}
		bindingBuffer.WriteString("]}\n")
	}
	return fmt.Sprintf(`
resource "google_storage_object" "object" {
  display_name = "%s"
	bucket = "google_storage_bucket.bucket.name"
}
resource "google_storage_bucket" "bucket" {
  display_name = "%s"
}
data "google_iam_policy" "test" {
  %s
}
resource "google_storage_Object_iam_policy" "test" {
  Object = "${google_storage_Object.Object.name}"
  policy_data = "${data.google_iam_policy.test.policy_data}"
}
`, objectName, bucketName, bindingBuffer.String())

}
