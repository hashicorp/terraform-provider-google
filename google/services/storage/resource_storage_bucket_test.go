// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

func TestAccStorageBucket_basic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
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
				ResourceName:            "google_storage_bucket.bucket",
				ImportStateId:           fmt.Sprintf("%s/%s", envvar.GetTestProjectFromEnv(), bucketName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccStorageBucket_basicWithAutoclass(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basicWithAutoclass(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
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
			// Autoclass is ForceNew, so this destroys & recreates, but does test the explicitly disabled config
			{
				Config: testAccStorageBucket_basicWithAutoclass(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
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
		},
	})
}

func TestAccStorageBucket_requesterPays(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-requester-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_requesterPays(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "requester_pays", "true"),
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

func TestAccStorageBucket_lowercaseLocation(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lowercaseLocation(bucketName),
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

func TestAccStorageBucket_dualLocation(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_dualLocation(bucketName),
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

func TestAccStorageBucket_customAttributes(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
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
	// multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRulesMultiple(bucketName),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_lifecycleRulesMultiple_update(bucketName),
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

func TestAccStorageBucket_lifecycleRuleStateLive(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRule_withStateLive(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(true), &bucket),
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

func TestAccStorageBucket_lifecycleRuleStateArchived(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRule_emptyArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(nil, &bucket),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(false), &bucket),
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

func TestAccStorageBucket_lifecycleRuleStateAny(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lifecycleRule_withStateArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(false), &bucket),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateLive(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(true), &bucket),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateAny(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(nil, &bucket),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_lifecycleRule_withStateArchived(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketLifecycleConditionState(googleapi.Bool(false), &bucket),
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

func TestAccStorageBucket_storageClass(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_storageClass(bucketName, "MULTI_REGIONAL", "US"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_storageClass(bucketName, "NEARLINE", "US"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &updated),
					// storage_class-only change should not recreate
					testAccCheckStorageBucketWasUpdated(&updated, &bucket),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_storageClass(bucketName, "REGIONAL", "US-CENTRAL1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &updated),
					// Location change causes recreate
					testAccCheckStorageBucketWasRecreated(&updated, &bucket),
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

func TestAccStorageBucket_update_requesterPays(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-requester-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_requesterPays(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_requesterPays(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &updated),
					testAccCheckStorageBucketWasUpdated(&updated, &bucket),
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

func TestAccStorageBucket_update(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var recreated storage.Bucket
	var updated storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
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
						t, "google_storage_bucket.bucket", bucketName, &recreated),
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
						t, "google_storage_bucket.bucket", bucketName, &updated),
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
						t, "google_storage_bucket.bucket", bucketName, &updated),
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
				Config: testAccStorageBucket_customAttributes_withLifecycle1Update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &updated),
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
						t, "google_storage_bucket.bucket", bucketName, &updated),
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
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				Config: testAccStorageBucket_customAttributes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(t, bucketName),
				),
			},
			{
				Config: testAccStorageBucket_customAttributes(fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketMissing(t, bucketName),
				),
			},
		},
	})
}

func TestAccStorageBucket_forceDestroyWithVersioning(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
				),
			},
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(t, bucketName),
				),
			},
			{
				Config: testAccStorageBucket_forceDestroyWithVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(t, bucketName),
				),
			},
		},
	})
}

func TestAccStorageBucket_forceDestroyObjectDeleteError(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_forceDestroyWithRetentionPolicy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketPutItem(t, bucketName),
				),
			},
			{
				Config:      testAccStorageBucket_forceDestroyWithRetentionPolicy(bucketName),
				Destroy:     true,
				ExpectError: regexp.MustCompile("could not delete non-empty bucket due to error when deleting contents"),
			},
			{
				Config: testAccStorageBucket_forceDestroy(bucketName),
			},
		},
	})
}

func TestAccStorageBucket_versioning(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_versioning(bucketName, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_versioning_empty(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.0.enabled", "true"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_versioning(bucketName, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.0.enabled", "false"),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_versioning_empty(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.#", "1"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "versioning.0.enabled", "false"),
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

func TestAccStorageBucket_logging(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
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
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
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
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket.bucket", "logging.#", "0"),
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

func TestAccStorageBucket_cors(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsCors(bucketName),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_basic(bucketName),
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

func TestAccStorageBucket_defaultEventBasedHold(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_defaultEventBasedHold(bucketName),
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

func TestAccStorageBucket_encryption(t *testing.T) {
	// when rotation is set, next rotation time is set using time.Now
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"organization":    envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
		"random_int":      acctest.RandInt(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_encryption(context),
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

func TestAccStorageBucket_publicAccessPrevention(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_publicAccessPrevention(bucketName, "enforced"),
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

func TestAccStorageBucket_uniformBucketAccessOnly(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_uniformBucketAccessOnly(bucketName, true),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_uniformBucketAccessOnly(bucketName, false),
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

func TestAccStorageBucket_labels(t *testing.T) {
	t.Parallel()

	bucketName := fmt.Sprintf("tf-test-acl-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			// Going from two labels
			{
				Config: testAccStorageBucket_updateLabels(bucketName),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			// Down to only one label (test single label deletion)
			{
				Config: testAccStorageBucket_labels(bucketName),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			// And make sure deleting all labels work
			{
				Config: testAccStorageBucket_basic(bucketName),
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

func TestAccStorageBucket_retentionPolicy(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_retentionPolicy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketRetentionPolicy(t, bucketName),
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

func TestAccStorageBucket_website(t *testing.T) {
	t.Parallel()

	bucketSuffix := fmt.Sprintf("tf-website-test-%d", acctest.RandInt(t))
	errRe := regexp.MustCompile("one of\n`website.0.main_page_suffix,website.0.not_found_page` must be specified")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccStorageBucket_websiteNoAttributes(bucketSuffix),
				ExpectError: errRe,
			},
			{
				Config: testAccStorageBucket_websiteOneAttribute(bucketSuffix),
			},
			{
				ResourceName:            "google_storage_bucket.website",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_websiteOneAttributeUpdate(bucketSuffix),
			},
			{
				ResourceName:            "google_storage_bucket.website",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_website(bucketSuffix),
			},
			{
				ResourceName:            "google_storage_bucket.website",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_websiteRemoved(bucketSuffix),
			},
			{
				ResourceName:            "google_storage_bucket.website",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccStorageBucket_retentionPolicyLocked(t *testing.T) {
	t.Parallel()

	var bucket storage.Bucket
	var newBucket storage.Bucket
	bucketName := fmt.Sprintf("tf-test-acc-bucket-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageBucketDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_lockedRetentionPolicy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &bucket),
					testAccCheckStorageBucketRetentionPolicy(t, bucketName),
				),
			},
			{
				ResourceName:            "google_storage_bucket.bucket",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccStorageBucket_retentionPolicy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageBucketExists(
						t, "google_storage_bucket.bucket", bucketName, &newBucket),
					testAccCheckStorageBucketWasRecreated(&newBucket, &bucket),
				),
			},
		},
	})
}

func testAccCheckStorageBucketExists(t *testing.T, n string, bucketName string, bucket *storage.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project_ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		found, err := config.NewStorageClient(config.UserAgent).Buckets.Get(rs.Primary.ID).Do()
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

func testAccCheckStorageBucketPutItem(t *testing.T, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		data := bytes.NewBufferString("test")
		dataReader := bytes.NewReader(data.Bytes())
		object := &storage.Object{Name: "bucketDestroyTestFile"}

		// This needs to use Media(io.Reader) call, otherwise it does not go to /upload API and fails
		if res, err := config.NewStorageClient(config.UserAgent).Objects.Insert(bucketName, object).Media(dataReader).Do(); err == nil {
			log.Printf("[INFO] Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			return fmt.Errorf("Objects.Insert failed: %v", err)
		}

		return nil
	}
}

func testAccCheckStorageBucketRetentionPolicy(t *testing.T, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		data := bytes.NewBufferString("test")
		dataReader := bytes.NewReader(data.Bytes())
		object := &storage.Object{Name: "bucketDestroyTestFile"}

		// This needs to use Media(io.Reader) call, otherwise it does not go to /upload API and fails
		if res, err := config.NewStorageClient(config.UserAgent).Objects.Insert(bucketName, object).Media(dataReader).Do(); err == nil {
			log.Printf("[INFO] Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			return fmt.Errorf("Objects.Insert failed: %v", err)
		}

		// Test deleting immediately, this should fail because of the 10 second retention
		if err := config.NewStorageClient(config.UserAgent).Objects.Delete(bucketName, objectName).Do(); err == nil {
			return fmt.Errorf("Objects.Delete succeeded: %v", object.Name)
		}

		// Wait 10 seconds and delete again
		time.Sleep(10000 * time.Millisecond)

		if err := config.NewStorageClient(config.UserAgent).Objects.Delete(bucketName, object.Name).Do(); err == nil {
			log.Printf("[INFO] Deleted object %v at location %v\n\n", object.Name, object.SelfLink)
		} else {
			return fmt.Errorf("Objects.Delete failed: %v", err)
		}

		return nil
	}
}

func testAccCheckStorageBucketMissing(t *testing.T, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		_, err := config.NewStorageClient(config.UserAgent).Buckets.Get(bucketName).Do()
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

func testAccStorageBucketDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_bucket" {
				continue
			}

			_, err := config.NewStorageClient(config.UserAgent).Buckets.Get(rs.Primary.ID).Do()
			if err == nil {
				return fmt.Errorf("Bucket still exists")
			}
		}

		return nil
	}
}

func testAccStorageBucket_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}
`, bucketName)
}

func testAccStorageBucket_basicWithAutoclass(bucketName string, autoclass bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
  autoclass  {
    enabled  = %t
  }
}
`, bucketName, autoclass)
}

func testAccStorageBucket_requesterPays(bucketName string, pays bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name           = "%s"
  location       = "US"
  requester_pays = %t
  force_destroy  = true
}
`, bucketName, pays)
}

func testAccStorageBucket_lowercaseLocation(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "eu"
  force_destroy = true
}
`, bucketName)
}

func testAccStorageBucket_dualLocation(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "ASIA"
  force_destroy = true
  custom_placement_config {
    data_locations = ["ASIA-EAST1", "ASIA-SOUTHEAST1"]
  }
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

func testAccStorageBucket_customAttributes_withLifecycle1Update(bucketName string) string {
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
      age = 0
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
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  storage_class = "%s"
  location      = "%s"
  force_destroy = true
}
`, bucketName, storageClass, location)
}

func testGoogleStorageBucketsCors(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
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
  name                     = "%s"
  location                 = "US"
  default_event_based_hold = true
  force_destroy            = true
}
`, bucketName)
}

func testAccStorageBucket_forceDestroyWithVersioning(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = "true"
  versioning {
    enabled = "true"
  }
}
`, bucketName)
}

func testAccStorageBucket_versioning(bucketName, enabled string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  versioning {
    enabled = "%s"
  }
}
`, bucketName, enabled)
}

func testAccStorageBucket_versioning_empty(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}
`, bucketName)
}

func testAccStorageBucket_logging(bucketName string, logBucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
  logging {
    log_bucket = "%s"
  }
}
`, bucketName, logBucketName)
}

func testAccStorageBucket_loggingWithPrefix(bucketName string, logBucketName string, prefix string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
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
  name          = "%s"
  location      = "US"
  force_destroy = true
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
      matches_storage_class = []
      age = 10
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      custom_time_before = "2019-01-01"
    }
  }
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
    condition {
      created_before = "2019-01-01"
      days_since_custom_time = 3
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
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "ARCHIVE"
    }
    condition {
      with_state = "ARCHIVED"
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      matches_prefix = ["test"]
      age            = 2
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      matches_suffix = ["test"]
      age            = 2
    }
  }
  lifecycle_rule {
    action {
      type = "AbortIncompleteMultipartUpload"
    }
    condition {
      age = 1
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_lifecycleRulesMultiple_update(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
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
      matches_storage_class = []
      age = 10
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      custom_time_before = "2019-01-01"
    }
  }
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
    condition {
      created_before = "2019-01-01"
      days_since_custom_time = 5
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
  lifecycle_rule {
    action {
      type          = "SetStorageClass"
      storage_class = "ARCHIVE"
    }
    condition {
      with_state = "ARCHIVED"
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      matches_prefix = ["test"]
      age            = 2
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      matches_suffix = ["test"]
      age            = 2
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_lifecycleRule_emptyArchived(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
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
  name          = "%s"
  location      = "US"
  force_destroy = true
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
  name          = "%s"
  location      = "US"
  force_destroy = true
  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age        = 10
      with_state = "LIVE"
	  days_since_noncurrent_time = 5
    }
  }
  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age        = 2
	  noncurrent_time_before = "2019-01-01"
    }
  }
}
`, bucketName)
}

func testAccStorageBucket_lifecycleRule_withStateAny(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
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
  name          = "%s"
  location      = "US"
  force_destroy = true
  labels = {
    my-label = "my-label-value"
  }
}
`, bucketName)
}

func testAccStorageBucket_uniformBucketAccessOnly(bucketName string, enabled bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                        = "%s"
  location                    = "US"
  uniform_bucket_level_access = %t
  force_destroy               = true
}
`, bucketName, enabled)
}

func testAccStorageBucket_publicAccessPrevention(bucketName string, prevention string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name                      = "%s"
  location                  = "US"
  public_access_prevention  = "%s"
  force_destroy             = true
}
`, bucketName, prevention)
}

func testAccStorageBucket_encryption(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "acceptance" {
  name            = "tf-test-%{random_suffix}"
  project_id      = "tf-test-%{random_suffix}"
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

data "google_storage_project_service_account" "gcs_account" {
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"
}

resource "google_storage_bucket" "bucket" {
  name          = "tf-test-crypto-bucket-%{random_int}"
  location      = "US"
  force_destroy = true
  encryption {
    default_kms_key_name = google_kms_crypto_key.crypto_key.id
  }

  depends_on = [google_kms_crypto_key_iam_member.iam]
}
`, context)
}

func testAccStorageBucket_updateLabels(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
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
  storage_class = "STANDARD"
  force_destroy = true

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
  name          = "%s"
  location      = "US"
  force_destroy = true

  retention_policy {
    retention_period = 10
  }
}
`, bucketName)
}

func testAccStorageBucket_lockedRetentionPolicy(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true

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
  storage_class = "STANDARD"
  force_destroy = true

  website {
  }
}
`, bucketName)
}

func testAccStorageBucket_websiteRemoved(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "website" {
  name          = "%s.gcp.tfacc.hashicorptest.com"
  location      = "US"
  storage_class = "STANDARD"
  force_destroy = true
}
`, bucketName)
}

func testAccStorageBucket_websiteOneAttribute(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "website" {
  name          = "%s.gcp.tfacc.hashicorptest.com"
  location      = "US"
  storage_class = "STANDARD"
  force_destroy = true

  website {
    main_page_suffix = "index.html"
  }
}
`, bucketName)
}

func testAccStorageBucket_websiteOneAttributeUpdate(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "website" {
  name          = "%s.gcp.tfacc.hashicorptest.com"
  location      = "US"
  storage_class = "STANDARD"
  force_destroy = true

  website {
    main_page_suffix = "default.html"
  }
}
`, bucketName)
}

func testAccStorageBucket_forceDestroy(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}
`, bucketName)
}

func testAccStorageBucket_forceDestroyWithRetentionPolicy(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true

  retention_policy {
    retention_period = 3600
  }
}
`, bucketName)
}
