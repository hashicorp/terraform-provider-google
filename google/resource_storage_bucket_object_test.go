package google

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

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

	bucketName := testBucketName(t)
	data := []byte("data data data")
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	data_md5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectBasic(bucketName, testFile.Name()),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, data_md5),
			},
		},
	})
}

func TestAccStorageObject_recreate(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)

	writeFile := func(name string, data []byte) string {
		h := md5.New()
		if _, err := h.Write(data); err != nil {
			t.Errorf("error calculating md5: %v", err)
		}
		data_md5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

		if err := ioutil.WriteFile(name, data, 0644); err != nil {
			t.Errorf("error writing file: %v", err)
		}
		return data_md5
	}
	testFile := getNewTmpTestFile(t, "tf-test")
	data_md5 := writeFile(testFile.Name(), []byte("data data data"))
	updatedName := testFile.Name() + ".update"
	updated_data_md5 := writeFile(updatedName, []byte("datum"))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectBasic(bucketName, testFile.Name()),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, data_md5),
			},
			{
				PreConfig: func() {
					err := os.Rename(updatedName, testFile.Name())
					if err != nil {
						t.Errorf("Failed to rename %s to %s", updatedName, testFile.Name())
					}
				},
				Config: testGoogleStorageBucketsObjectBasic(bucketName, testFile.Name()),
				Check:  testAccCheckGoogleStorageObject(t, bucketName, objectName, updated_data_md5),
			},
		},
	})
}

func TestAccStorageObject_content(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	data_md5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectContent(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, data_md5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "content_type", "text/plain; charset=utf-8"),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "storage_class", "STANDARD"),
				),
			},
		},
	})
}

func TestAccStorageObject_withContentCharacteristics(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	data_md5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	disposition, encoding, language, content_type := "inline", "compress", "en", "binary/octet-stream"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObject_optionalContentFields(
					bucketName, disposition, encoding, language, content_type),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, data_md5),
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

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObjectDynamicContent(testBucketName(t)),
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

	bucketName := testBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	data_md5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	cacheControl := "private"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObject_cacheControl(bucketName, testFile.Name(), cacheControl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, data_md5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "cache_control", cacheControl),
				),
			},
		},
	})
}

func TestAccStorageObject_storageClass(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	data_md5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	storageClass := "MULTI_REGIONAL"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObject_storageClass(bucketName, storageClass),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, data_md5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "storage_class", storageClass),
				),
			},
		},
	})
}

func TestAccStorageObject_metadata(t *testing.T) {
	t.Parallel()

	bucketName := testBucketName(t)
	data := []byte(content)
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		t.Errorf("error calculating md5: %v", err)
	}
	data_md5 := base64.StdEncoding.EncodeToString(h.Sum(nil))
	testFile := getNewTmpTestFile(t, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccStorageObjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageBucketsObject_metadata(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObject(t, bucketName, objectName, data_md5),
					resource.TestCheckResourceAttr(
						"google_storage_bucket_object.object", "metadata.customKey", "custom_value"),
				),
			},
		},
	})
}

func testAccCheckGoogleStorageObject(t *testing.T, bucket, object, md5 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		objectsService := storage.NewObjectsService(config.clientStorage)

		getCall := objectsService.Get(bucket, object)
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

func testAccStorageObjectDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_bucket_object" {
				continue
			}

			bucket := rs.Primary.Attributes["bucket"]
			name := rs.Primary.Attributes["name"]

			objectsService := storage.NewObjectsService(config.clientStorage)

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
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name    = "%s"
  bucket  = google_storage_bucket.bucket.name
  content = "%s"
}
`, bucketName, objectName, content)
}

func testGoogleStorageBucketsObjectDynamicContent(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
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
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name   = "%s"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}
`, bucketName, objectName, sourceFilename)
}

func testGoogleStorageBucketsObject_optionalContentFields(
	bucketName, disposition, encoding, language, content_type string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
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

func testGoogleStorageBucketsObject_cacheControl(bucketName, sourceFilename, cacheControl string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name          = "%s"
  bucket        = google_storage_bucket.bucket.name
  source        = "%s"
  cache_control = "%s"
}
`, bucketName, objectName, sourceFilename, cacheControl)
}

func testGoogleStorageBucketsObject_storageClass(bucketName string, storageClass string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "object" {
  name          = "%s"
  bucket        = google_storage_bucket.bucket.name
  content       = "%s"
  storage_class = "%s"
}
`, bucketName, objectName, content, storageClass)
}

func testGoogleStorageBucketsObject_metadata(bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
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

// Creates a new tmp test file. Fails the current test if we cannot create
// new tmp file in the filesystem.
func getNewTmpTestFile(t *testing.T, prefix string) *os.File {
	testFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		t.Fatalf("Cannot create temp file: %s", err)
	}
	return testFile
}
