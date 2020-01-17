package google

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

func testBucketName() string {
	return fmt.Sprintf("%s-%d", "tf-test-bucket", acctest.RandInt())
}

func TestAccStorageBucket_basic(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "false"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportStateId:     fmt.Sprintf("%s/%s", getTestProjectFromEnv(), bucketName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_requesterPays(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-requester-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_requesterPays(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "requester_pays", "true"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_lowercaseLocation(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lowercaseLocation(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_customAttributes(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccStorageBucket_lifecycleRulesMultiple(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRulesMultiple(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_lifecycleRuleStateLive(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())
	hashK := resourceGCSBucketLifecycleRuleConditionHash(map[string]interface{}{
		"age":                10,
		"with_state":         "LIVE",
		"num_newer_versions": 0,
		"created_before":     "",
	})
	attrPrefix := fmt.Sprintf("lifecycle_rule.0.condition.%d.", hashK)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRule_withStateLive(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(true), &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", attrPrefix+"with_state", "LIVE"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_lifecycleRuleStateArchived(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())
	hashK := resourceGCSBucketLifecycleRuleConditionHash(map[string]interface{}{
		"age":                10,
		"with_state":         "ARCHIVED",
		"num_newer_versions": 0,
		"created_before":     "",
	})
	attrPrefix := fmt.Sprintf("lifecycle_rule.0.condition.%d.", hashK)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRule_emptyArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(nil, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(false), &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", attrPrefix+"with_state", "ARCHIVED"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_lifecycleRuleStateAny(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())

	hashKLive := resourceGCSBucketLifecycleRuleConditionHash(map[string]interface{}{
		"age":                10,
		"with_state":         "LIVE",
		"num_newer_versions": 0,
		"created_before":     "",
	})
	hashKArchived := resourceGCSBucketLifecycleRuleConditionHash(map[string]interface{}{
		"age":                10,
		"with_state":         "ARCHIVED",
		"num_newer_versions": 0,
		"created_before":     "",
	})
	hashKAny := resourceGCSBucketLifecycleRuleConditionHash(map[string]interface{}{
		"age":                10,
		"with_state":         "ANY",
		"num_newer_versions": 0,
		"created_before":     "",
	})

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRule_withStateArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(false), &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", fmt.Sprintf("lifecycle_rule.0.condition.%d.with_state", hashKArchived), "ARCHIVED"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateLive(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(true), &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", fmt.Sprintf("lifecycle_rule.0.condition.%d.with_state", hashKLive), "LIVE"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateAny(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(nil, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", fmt.Sprintf("lifecycle_rule.0.condition.%d.with_state", hashKAny), "ANY"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(false), &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", fmt.Sprintf("lifecycle_rule.0.condition.%d.with_state", hashKArchived), "ARCHIVED"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_storageClass(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_storageClass(bucketName, "MULTI_REGIONAL", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_storageClass(bucketName, "NEARLINE", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					// storage_class-only change should not recreate
					testAccCheckStorageBucketWasUpdated(&updated, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_storageClass(bucketName, "REGIONAL", "US-CENTRAL1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					// Location change causes recreate
					testAccCheckStorageBucketWasRecreated(&updated, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_update_requesterPays(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-requester-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_requesterPays(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_requesterPays(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &bucket),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_update(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var recreated storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "false"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &recreated),
					testAccCheckStorageBucketWasRecreated(&recreated, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes_withLifecycle1(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &recreated),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes_withLifecycle2(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &recreated),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &recreated),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "force_destroy", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccStorageBucket_forceDestroy(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(bucketName),
				),
			},
			{
				Config: testAccStorageBucket_customAttributes(acctest.RandomWithPrefix("tf-test-acl-bucket")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketMissing(bucketName),
				),
			},
		},
	})
}

func TestAccStorageBucket_forceDestroyWithVersioning(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(bucketName),
				),
			},
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(bucketName),
				),
			},
		},
	})
}

func TestAccStorageBucket_versioning(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_versioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.0.enabled", "true"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_logging(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_logging(bucketName, "log-bucket"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_bucket", "log-bucket"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_object_prefix", bucketName),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_loggingWithPrefix(bucketName, "another-log-bucket", "object-prefix"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_bucket", "another-log-bucket"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.0.log_object_prefix", "object-prefix"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.#", "0"),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_cors(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsCors(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_defaultEventBasedHold(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_defaultEventBasedHold(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_encryption(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"organization":    getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(10),
		"random_int":      acctest.RandInt(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_encryption(context),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_bucketPolicyOnly(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_bucketPolicyOnly(bucketName, true),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_bucketPolicyOnly(bucketName, false),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_labels(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			// Going from two labels
			{
				Config: testAccStorageBucket_updateLabels(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Down to only one label (test single label deletion)
			{
				Config: testAccStorageBucket_labels(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// And make sure deleting all labels work
			{
				Config: testAccStorageBucket_basic(bucketName),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_retentionPolicy(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_retentionPolicy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketRetentionPolicy(bucketName),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_website(t *testing.T) {
	t.Parallel()

	bucketSuffix := acctest.RandomWithPrefix("tf-website-test")
	errRe := regexp.MustCompile("one of `((website.0.main_page_suffix,website.0.not_found_page)|(website.0.not_found_page,website.0.main_page_suffix))` must be specified")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccStorageBucket_websiteNoAttributes(bucketSuffix),
				ExpectError: errRe,
			},
			{
				Config: testAccStorageBucket_websiteOneAttribute(bucketSuffix),
			},
			{
				ResourceName:      "google_storage_bucket.website",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_website(bucketSuffix),
			},
			{
				ResourceName:      "google_storage_bucket.website",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_retentionPolicyLocked(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var newBucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lockedRetentionPolicy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketRetentionPolicy(bucketName),
				),
			},
			{
				ResourceName:      "google_storage_bucket.bucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_retentionPolicy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						"google_storage_bucket.bucket", bucketName, &newBucket),
					testAccCheckStorageBucketWasRecreated(&newBucket, &bucket),
				),
			},
		},
	})
}

func testAccCheckStorageBucketExists(n string, bucketName string, bucket *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project_ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientStorage.Buckets.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Bucket not found")
		}

		if found.Name != bucketName {
			return fmt.Errorf("expected name %s, got %s", bucketName, found.Name)
		}

		*bucket = *found
		return nil
	}
}

func testAccCheckStorageBucketWasUpdated(newBucket *storage.Bucket, b *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if newBucket.TimeCreated != b.TimeCreated {
			return fmt.Errorf("expected storage bucket to have been updated (had same creation time), instead was recreated - old creation time %s, new creation time %s", newBucket.TimeCreated, b.TimeCreated)
		}
		return nil
	}
}

func testAccCheckStorageBucketWasRecreated(newBucket *storage.Bucket, b *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if newBucket.TimeCreated == b.TimeCreated {
			return fmt.Errorf("expected storage bucket to have been recreated, instead had same creation time (%s)", b.TimeCreated)
		}
		return nil
	}
}

func testAccCheckStorageBucketPutItem(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		data := bytes.NewBufferString("test")
		dataReader := bytes.NewReader(data.Bytes())
		object := &storage.Object{Name: "bucketDestroyTestFile"}

		// This needs to use Media(io.Reader) call, otherwise it does not go to /upload API and fails
		if res, err := config.clientStorage.Objects.Insert(bucketName, object).Media(dataReader).Do(); err == nil {
			log.Printf("[INFO] Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			return fmt.Errorf("Objects.Insert failed: %v", err)
		}

		return nil
	}
}

func testAccCheckStorageBucketRetentionPolicy(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		data := bytes.NewBufferString("test")
		dataReader := bytes.NewReader(data.Bytes())
		object := &storage.Object{Name: "bucketDestroyTestFile"}

		// This needs to use Media(io.Reader) call, otherwise it does not go to /upload API and fails
		if res, err := config.clientStorage.Objects.Insert(bucketName, object).Media(dataReader).Do(); err == nil {
			log.Printf("[INFO] Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			return fmt.Errorf("Objects.Insert failed: %v", err)
		}

		// Test deleting immediately, this should fail because of the 10 second retention
		if err := config.clientStorage.Objects.Delete(bucketName, objectName).Do(); err == nil {
			return fmt.Errorf("Objects.Delete succeeded: %v", object.Name)
		}

		// Wait 10 seconds and delete again
		time.Sleep(10000 * time.Millisecond)

		if err := config.clientStorage.Objects.Delete(bucketName, object.Name).Do(); err == nil {
			log.Printf("[INFO] Deleted object %v at location %v\n\n", object.Name, object.SelfLink)
		} else {
			return fmt.Errorf("Objects.Delete failed: %v", err)
		}

		return nil
	}
}

func testAccCheckStorageBucketMissing(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		_, err := config.clientStorage.Buckets.Get(bucketName).Do()
		if err == nil {
			return fmt.Errorf("Found %s", bucketName)
		}

		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return nil
		}

		return err
	}
}

func testAccCheckStorageBucketLifecycleConditionState(expected *bool, b *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual := b.Lifecycle.Rule[0].Condition.IsLive
		if expected == nil && b.Lifecycle.Rule[0].Condition.IsLive == nil {
			return nil
		}
		if expected == nil {
			return fmt.Errorf("expected condition isLive to be unset, instead got %t", *actual)
		}
		if actual == nil {
			return fmt.Errorf("expected condition isLive to be %t, instead got nil (unset)", *expected)
		}
		if *expected != *actual {
			return fmt.Errorf("expected condition isLive to be %t, instead got %t", *expected, *actual)
		}
		return nil
	}
}

func testAccStorageBucketDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_bucket" {
			continue
		}

		_, err := config.clientStorage.Buckets.Get(rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Bucket still exists")
		}
	}

	return nil
}

func testAccStorageBucket_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}
`, bucketName)
}

func testAccStorageBucket_requesterPays(bucketName string, pays bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name           = "%s"
  requester_pays = %t
}
`, bucketName, pays)
}

func testAccStorageBucket_lowercaseLocation(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "eu"
}
`, bucketName)
}

func testAccStorageBucket_customAttributes(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "EU"
  force_destroy = "true"
}
`, bucketName)
}

func testAccStorageBucket_customAttributes_withLifecycle1(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "EU"
  force_destroy = "true"
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age = 10
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_customAttributes_withLifecycle2(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "EU"
  force_destroy = "true"
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
    condition {
      age = 2
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age                = 10
      num_newer_versions = 2
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_storageClass(bucketName, storageClass, location string) string {
	var locationBlock string
	if location != "" {
		locationBlock = fmt.Sprintf(`
	location = "%s"`, location)
	}
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  storage_class = "%s"%s
}
`, bucketName, storageClass, locationBlock)
}

func testGoogleStorageBucketsCors(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  cors {
    origin          = ["abc", "def"]
    method          = ["a1a"]
    response_header = ["123", "456", "789"]
    max_age_seconds = 10
  }

  cors {
    origin          = ["ghi", "jkl"]
    method          = ["z9z"]
    response_header = ["000"]
    max_age_seconds = 5
  }
}
`, bucketName)
}

func testAccStorageBucket_defaultEventBasedHold(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  default_event_based_hold = true
}
`, bucketName)
}

func testAccStorageBucket_forceDestroyWithVersioning(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  force_destroy = "true"
  versioning {
    enabled = "true"
  }
}
`, bucketName)
}

func testAccStorageBucket_versioning(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  versioning {
    enabled = "true"
  }
}
`, bucketName)
}

func testAccStorageBucket_logging(bucketName string, logBucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  logging {
    log_bucket = "%s"
  }
}
`, bucketName, logBucketName)
}

func testAccStorageBucket_loggingWithPrefix(bucketName string, logBucketName string, prefix string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  logging {
    log_bucket        = "%s"
    log_object_prefix = "%s"
  }
}
`, bucketName, logBucketName, prefix)
}

func testAccStorageBucket_lifecycleRulesMultiple(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
    condition {
      matches_storage_class = ["COLDLINE"]
      age                   = 2
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age = 10
    }
  }
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
    condition {
      created_before = "2019-01-01"
    }
  }
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
    condition {
      num_newer_versions = 10
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_lifecycleRule_emptyArchived(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age = 10
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_lifecycleRule_withStateArchived(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age        = 10
      with_state = "ARCHIVED"
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_lifecycleRule_withStateLive(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age        = 10
      with_state = "LIVE"
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_lifecycleRule_withStateAny(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age        = 10
      with_state = "ANY"
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_labels(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  labels = {
    my-label = "my-label-value"
  }
}
`, bucketName)
}

func testAccStorageBucket_bucketPolicyOnly(bucketName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name               = "%s"
  bucket_policy_only = %t
}
`, bucketName, enabled)
}

func testAccStorageBucket_encryption(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "acceptance" {
  name            = "terraform-%{random_suffix}"
  project_id      = "terraform-%{random_suffix}"
  org_id          = "%{organization}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_key_ring" "key_ring" {
  name     = "tf-test-%{random_suffix}"
  project  = google_project_service.acceptance.project
  location = "us"
}

resource "google_kms_crypto_key" "crypto_key" {
  name            = "tf-test-%{random_suffix}"
  key_ring        = google_kms_key_ring.key_ring.id
  rotation_period = "1000000s"
}

resource "google_storage_bucket" "bucket" {
  name = "tf-test-crypto-bucket-%{random_int}"
  encryption {
    default_kms_key_name = google_kms_crypto_key.crypto_key.self_link
  }
}
`, context)
}

func testAccStorageBucket_updateLabels(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
  labels = {
    my-label    = "my-updated-label-value"
    a-new-label = "a-new-label-value"
  }
}
`, bucketName)
}

func testAccStorageBucket_website(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "website" {
  name          = "%s.gcp.tfacc.hashicorptest.com"
  location      = "US"
  storage_class = "MULTI_REGIONAL"

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}
`, bucketName)
}

func testAccStorageBucket_retentionPolicy(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"

  retention_policy {
    retention_period = 10
  }
}
`, bucketName)
}

func testAccStorageBucket_lockedRetentionPolicy(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"

  retention_policy {
    is_locked        = true
    retention_period = 10
  }
}
`, bucketName)
}

func testAccStorageBucket_websiteNoAttributes(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "website" {
  name          = "%s.gcp.tfacc.hashicorptest.com"
  location      = "US"
  storage_class = "MULTI_REGIONAL"

  website {
  }
}
`, bucketName)
}

func testAccStorageBucket_websiteOneAttribute(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "website" {
  name          = "%s.gcp.tfacc.hashicorptest.com"
  location      = "US"
  storage_class = "MULTI_REGIONAL"

  website {
    main_page_suffix = "index.html"
  }
}
`, bucketName)
}
