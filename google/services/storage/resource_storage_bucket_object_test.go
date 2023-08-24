// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"os"

	"google.golang.org/api/storage/v1"
)

const (
	objectName = "tf-gce-test"
	content    = "now this is content!"
)

func TestAccStorageObject_basic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte("data data data")
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectBasic(bucketName, testFile.Name()),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
			},
		},
	})
}

func TestAccStorageObject_recreate(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)

	writeFile := func(name string, data []byte) string {
		h := md5.New()
		if _, err := h.Write(data); err != nil {
			t.Errorf("error calculating md5: %v", err)
		}
		dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

		if err := ioutil.WriteFile(name, data, 0644); err != nil {
			t.Errorf("error writing file: %v", err)
		}
		return dataMd5
	}
	testFile := getNewTmpTestFile(t, "tf-test")
	dataMd5 := writeFile(testFile.Name(), []byte("data data data"))
	updatedName := testFile.Name() + ".update"
	updatedDataMd5 := writeFile(updatedName, []byte("datum"))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectBasic(bucketName, testFile.Name()),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
			},
			{
				PreConfig: func() {
					err := os.Rename(updatedName, testFile.Name())
					if err != nil {
						t.Errorf("Failed to rename %s to %s", updatedName, testFile.Name())
					}
				},
				Config: testGoogleStorageBucketsObjectBasic(bucketName, testFile.Name()),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, updatedDataMd5),
			},
		},
	})
}

func TestAccStorageObject_content(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectContent(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "content_type", "text/plain; charset=utf-8"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func TestAccStorageObject_folder(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	folderName := "tf-gce-folder-test/"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsFolder(bucketName, folderName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageFolder(t, bucketName, folderName),
				),
			},
		},
	})
}

func TestAccStorageObject_withContentCharacteristics(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	disposition, encoding, language, content_type := "inline", "compress", "en", "binary/octet-stream"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectOptionalContentFields(
					bucketName, disposition, encoding, language, content_type),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "content_disposition", disposition),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "content_encoding", encoding),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "content_language", language),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "content_type", content_type),
				),
			},
		},
	})
}

func TestAccStorageObject_dynamicContent(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectDynamicContent(acctest.TestBucketName(t)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "content_type", "text/plain; charset=utf-8"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func TestAccStorageObject_cacheControl(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	cacheControl := "private"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectCacheControl(bucketName, testFile.Name(), cacheControl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "cache_control", cacheControl),
				),
			},
		},
	})
}

func TestAccStorageObject_storageClass(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	storageClass := "MULTI_REGIONAL"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectStorageClass(bucketName, storageClass),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "storage_class", storageClass),
				),
			},
		},
	})
}

func TestAccStorageObject_metadata(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectMetadata(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "metadata.customKey", "custom_value"),
				),
			},
		},
	})
}

func TestAccStorageObjectKms(t *testing.T) {
	t.Parallel()

	kms := acctest.BootstrapKMSKeyInLocation(t, "us")
	bucketName := acctest.TestBucketName(t)
	data := []byte("data data data")
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectKms(bucketName, testFile.Name(), kms.CryptoKey.Name),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
			},
		},
	})
}

func TestAccStorageObject_customerEncryption(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	customerEncryptionKey := "qI6+xvCZE9jUm94nJWIulFc8rthN64ybkGCsLUY9Do4="
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectCustomerEncryption(bucketName, customerEncryptionKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectWithEncryption(t, bucketName, objectName, dataMd5, customerEncryptionKey),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "customer_encryption.0.encryption_key", customerEncryptionKey),
				),
			},
		},
	})
}

func TestAccStorageObject_holds(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	dataMd5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectHolds(bucketName, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "event_based_hold", "true"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "temporary_hold", "true"),
				),
			},
			{
				Config: testGoogleStorageBucketsObjectHolds(bucketName, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, dataMd5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "event_based_hold", "false"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "temporary_hold", "false"),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageObject(t *testing.T, bucket, object, md5 string) resource.TestCheckFunc {
	return testAccCheckGoogleStorageObjectWithEncryption(t, bucket, object, md5, "")
}

func testAccCheckGoogleStorageObjectWithEncryption(t *testing.T, bucket, object, md5 string, customerEncryptionKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		objectsService := storage.NewObjectsService(config.NewStorageClient(config.UserAgent))

		getCall := objectsService.Get(bucket, object)
		if customerEncryptionKey != "" {
			decodedKey, _ := base64.StdEncoding.DecodeString(customerEncryptionKey)
			keyHash := sha256.Sum256(decodedKey)
			headers := getCall.Header()
			headers.Set("x-goog-encryption-algorithm", "AES256")
			headers.Set("x-goog-encryption-key", customerEncryptionKey)
			headers.Set("x-goog-encryption-key-sha256", base64.StdEncoding.EncodeToString(keyHash[:]))
		}
		res, err := getCall.Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of object %s: %s", object, err)
		}

		if md5 != res.Md5Hash {
			return fmt.Errorf("Error contents of %s garbled, md5 hashes don't match (%s, %s)", object, md5, res.Md5Hash)
		}

		return nil
	}
}

func testAccCheckGoogleStorageFolder(t *testing.T, bucket, folderName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		objectsService := storage.NewObjectsService(config.NewStorageClient(config.UserAgent))

		getCall := objectsService.Get(bucket, folderName)
		res, err := getCall.Do()

		if err != nil {
			return fmt.Errorf("Error retrieving folder %s: %s", folderName, err)
		}

		if folderName != res.Name {
			return fmt.Errorf("Error folder name don't match (%s, %s)", folderName, res.Name)
		}

		return nil
	}
}

func testAccStorageObjectDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_bucket_object" {
				continue
			}

			bucket := rs.Primary.Attributes["bucket"]
			name := rs.Primary.Attributes["name"]

			objectsService := storage.NewObjectsService(config.NewStorageClient(config.UserAgent))

			getCall := objectsService.Get(bucket, name)
			_, err := getCall.Do()

			if err == nil {
				return fmt.Errorf("Object %s still exists", name)
			}
		}

		return nil
	}
}

func testGoogleStorageBucketsObjectContent(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "object" {
  name    = "%s"
  bucket  = google_storage_bucket.bucket.name
  content = "%s"
}
`, bucketName, objectName, content)
}

func testGoogleStorageBucketsFolder(bucketName, folderName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "object" {
  name    = "%s"
  bucket  = google_storage_bucket.bucket.name
  content = " "
}
`, bucketName, folderName)
}

func testGoogleStorageBucketsObjectDynamicContent(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name    = "%s"
  bucket  = google_storage_bucket.bucket.name
  content = google_storage_bucket.bucket.project
}
`, bucketName, objectName)
}

func testGoogleStorageBucketsObjectBasic(bucketName, sourceFilename string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}
`, bucketName, objectName, sourceFilename)
}

func testGoogleStorageBucketsObjectOptionalContentFields(
	bucketName, disposition, encoding, language, content_type string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name                = "%s"
  bucket              = google_storage_bucket.bucket.name
  content             = "%s"
  content_disposition = "%s"
  content_encoding    = "%s"
  content_language    = "%s"
  content_type        = "%s"
}
`, bucketName, objectName, content, disposition, encoding, language, content_type)
}

func testGoogleStorageBucketsObjectCacheControl(bucketName, sourceFilename, cacheControl string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name          = "%s"
  bucket        = google_storage_bucket.bucket.name
  source        = "%s"
  cache_control = "%s"
}
`, bucketName, objectName, sourceFilename, cacheControl)
}

func testGoogleStorageBucketsObjectStorageClass(bucketName string, storageClass string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name          = "%s"
  bucket        = google_storage_bucket.bucket.name
  content       = "%s"
  storage_class = "%s"
}
`, bucketName, objectName, content, storageClass)
}

func testGoogleStorageBucketsObjectMetadata(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name          = "%s"
  bucket        = google_storage_bucket.bucket.name
  content       = "%s"

  metadata = {
    "customKey" = "custom_value"
  }
}
`, bucketName, objectName, content)
}

func testGoogleStorageBucketsObjectCustomerEncryption(bucketName string, customerEncryptionKey string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "object" {
  name                = "%s"
  bucket              = google_storage_bucket.bucket.name
  content             = "%s"
  customer_encryption {
    encryption_key = "%s"
  }
}
`, bucketName, objectName, content, customerEncryptionKey)
}

func testGoogleStorageBucketsObjectHolds(bucketName string, eventBasedHold bool, temporaryHold bool) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name          = "%s"
  location      = "US"
  force_destroy = true
}

resource "google_storage_bucket_object" "object" {
  name             = "%s"
  bucket           = google_storage_bucket.bucket.name
  content          = "%s"
  event_based_hold = %t
  temporary_hold   = %t
}
`, bucketName, objectName, content, eventBasedHold, temporaryHold)
}

func testGoogleStorageBucketsObjectKms(bucketName, sourceFilename, kmsKey string) string {
	return fmt.Sprintf(`

resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

data "google_storage_project_service_account" "gcs" {
}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%s"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_storage_project_service_account.gcs.email_address}"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
  kms_key_name = google_kms_crypto_key_iam_member.crypto_key.crypto_key_id
}
`, bucketName, kmsKey, objectName, sourceFilename)
}

// Creates a new tmp test file. Fails the current test if we cannot create
// new tmp file in the filesystem.
func getNewTmpTestFile(t *testing.T, prefix string) *os.File {
	testFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		t.Fatalf("Cannot create temp file: %s", err)
	}
	return testFile
}
