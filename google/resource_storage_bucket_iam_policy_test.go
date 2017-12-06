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

func testBucketIamName() string {
	return acctest.RandomWithPrefix("tf-test-iam-bucket")
}

func TestAccGoogleStorageBucketIAMPolicy_basic(t *testing.T) {
	bucketName := testBucketIamName()
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
		CheckDestroy: testAccGoogleStorageBucketIAMPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageBucketsIAMPolicyBasic(bucketName, policy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketIAMPolicy(bucketName, policy),
				),
			},
		},
	})
}

func TestAccGoogleStorageBucketIAMPolicy_upgrade(t *testing.T) {
	bucketName := testBucketIamName()
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
		CheckDestroy: testAccGoogleStorageBucketIAMPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageBucketsIAMPolicyBasic(bucketName, policy1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketIAMPolicy(bucketName, policy1),
				),
			},

			resource.TestStep{
				Config: testGoogleStorageBucketsIAMPolicyBasic(bucketName, policy2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageBucketIAMPolicy(bucketName, policy2),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageBucketIAMPolicy(bucket string, policy *storagev1.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[bucket]
		if !ok {
			return fmt.Errorf("Not found: %s", bucket)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		p, err := config.clientStorage.Buckets.GetIamPolicy(rs.Primary.ID).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of IAMPolicy for bucket %s: %s", bucket, err)
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

func testAccGoogleStorageBucketIAMPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_bucket_iam_policy" {
			continue
		}

		bucket := rs.Primary.Attributes["bucket"]

		policy, err := config.clientStorage.Buckets.GetIamPolicy(bucket).Do()

		if err != nil && len(policy.Bindings) > 0 {
			return fmt.Errorf("Bucket %s policy hasn't been deleted", bucket)
		}
	}

	return nil
}

func testGoogleStorageBucketsIAMPolicyBasic(bucketName string, policy *storagev1.Policy) string {
	var bindingBuffer bytes.Buffer

	// Generate binding for google_iam_policy based on policy.Bindings
	// Example:
	// data "google_iam_policy" "admin" {
	// 	binding {
	//  role = "roles/editor"
	//    members = [
	//      "user:jane@example.com",
	//    ]
	//  }
	//}

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
resource "google_storage_bucket" "bucket" {
  display_name = "%s"
}
data "google_iam_policy" "test" {
  %s
}
resource "google_storage_bucket_iam_policy" "test" {
  bucket = "${google_storage_bucket.bucket.name}"
  policy_data = "${data.google_iam_policy.test.policy_data}"
}
`, bucketName, bindingBuffer.String())

}
