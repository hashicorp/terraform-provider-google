// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	"github.com/hashicorp/terraform-provider-google/google/services/storage"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var tfObjectAcl, errObjectAcl = ioutil.TempFile("", "tf-gce-test")

func testAclObjectName(t *testing.T) string {
	return fmt.Sprintf("%s-%d", "tf-test-acl-object", acctest.RandInt(t))
}

func TestAccStorageObjectAcl_basic(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_upgrade(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic2),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasic2(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasicDelete(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAclDelete(t, bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAclDelete(t, bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAclDelete(t, bucketName,
						objectName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_downgrade(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic2(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic3_owner),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasic3(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic3_reader),
				),
			},

			{
				Config: testGoogleStorageObjectsAclBasicDelete(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAclDelete(t, bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAclDelete(t, bucketName,
						objectName, roleEntityBasic2),
					testAccCheckGoogleStorageObjectAclDelete(t, bucketName,
						objectName, roleEntityBasic3_reader),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_predefined(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclPredefined(bucketName, objectName),
			},
		},
	})
}

func TestAccStorageObjectAcl_predefinedToExplicit(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclPredefined(bucketName, objectName),
			},
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic2),
				),
			},
		},
	})
}

func TestAccStorageObjectAcl_explicitToPredefined(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic1),
					testAccCheckGoogleStorageObjectAcl(t, bucketName,
						objectName, roleEntityBasic2),
				),
			},
			{
				Config: testGoogleStorageObjectsAclPredefined(bucketName, objectName),
			},
		},
	})
}

// Test that we allow the API to reorder our role entities without perma-diffing.
func TestAccStorageObjectAcl_unordered(t *testing.T) {
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccStorageObjectAclDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleStorageObjectAclUnordered(bucketName, objectName),
			},
		},
	})
}

// a round tripper that returns fake response for get object API removing `owner` attribute
// it only modifies the response once, since otherwise resource will fail to delete.
type testRoundTripper struct {
	http.RoundTripper
	bucketName, objectName string
	done                   bool
}

func (t *testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	response, err := t.RoundTripper.RoundTrip(r)
	if err != nil {
		return response, err
	}
	expectedPath := fmt.Sprintf("/storage/v1/b/%s/o/%s", t.bucketName, t.objectName)
	if t.done || r.URL.Path != expectedPath || r.Host != "storage.googleapis.com" {
		return response, err
	}
	t.done = true
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response, err
	}
	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(responseBytes, &responseMap)
	if err != nil {
		return response, err
	}
	delete(responseMap, "owner")
	responseBytes, err = json.Marshal(responseMap)
	if err != nil {
		return response, err
	}
	response.Body = io.NopCloser(bytes.NewBuffer(responseBytes))
	return response, nil
}

// Test that we don't fail if there's no owner for object
func TestAccStorageObjectAcl_noOwner(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	bucketName := acctest.TestBucketName(t)
	objectName := testAclObjectName(t)
	objectData := []byte("data data data")
	if err := ioutil.WriteFile(tfObjectAcl.Name(), objectData, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}

	// TODO (mbang) we can leave this one using the SDK provider as we need to overwrite the configure function,
	// which we can't do in the plugin-framework version of the provider. When this resource does get updated to
	// use plugin-framework, best I can guess we'll want to do something similar to NewFrameworkTestProvider where
	// we have a nested production version of the provider, we re-write configure to call the production version and
	// add the additional things inside there.
	provider := provider.Provider()
	oldConfigureFunc := provider.ConfigureContextFunc
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		c, diagnostics := oldConfigureFunc(ctx, d)
		config := c.(*transport_tpg.Config)
		roundTripper := &testRoundTripper{RoundTripper: config.Client.Transport, bucketName: bucketName, objectName: objectName}
		config.Client.Transport = roundTripper
		return c, diagnostics
	}
	providers := map[string]*schema.Provider{
		"google": provider,
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			if errObjectAcl != nil {
				panic(errObjectAcl)
			}
			acctest.AccTestPreCheck(t)
		},
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:             testGoogleStorageObjectsAclBasic1(bucketName, objectName),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckGoogleStorageObjectAcl(t *testing.T, bucket, object, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := storage.GetRoleEntityPair(roleEntityS)
		config := acctest.GoogleProviderConfig(t)

		res, err := config.NewStorageClient(config.UserAgent).ObjectAccessControls.Get(bucket,
			object, roleEntity.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error retrieving contents of acl for bucket %s: %s", bucket, err)
		}

		if res.Role != roleEntity.Role {
			return fmt.Errorf("Error, Role mismatch %s != %s", res.Role, roleEntity.Role)
		}

		return nil
	}
}

func testAccCheckGoogleStorageObjectAclDelete(t *testing.T, bucket, object, roleEntityS string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roleEntity, _ := storage.GetRoleEntityPair(roleEntityS)
		config := acctest.GoogleProviderConfig(t)

		_, err := config.NewStorageClient(config.UserAgent).ObjectAccessControls.Get(bucket,
			object, roleEntity.Entity).Do()

		if err != nil {
			return nil
		}

		return fmt.Errorf("Error, Entity still exists %s", roleEntity.Entity)
	}
}

func testAccStorageObjectAclDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_storage_bucket_acl" {
				continue
			}

			bucket := rs.Primary.Attributes["bucket"]
			object := rs.Primary.Attributes["object"]

			_, err := config.NewStorageClient(config.UserAgent).ObjectAccessControls.List(bucket, object).Do()

			if err == nil {
				return fmt.Errorf("Acl for bucket %s still exists", bucket)
			}
		}

		return nil
	}
}

func testGoogleStorageObjectsAclBasicDelete(bucketName string, objectName string) string {
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

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = []
}
`, bucketName, objectName, tfObjectAcl.Name())
}

func testGoogleStorageObjectsAclBasic1(bucketName string, objectName string) string {
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

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(),
		roleEntityBasic1, roleEntityBasic2)
}

func testGoogleStorageObjectsAclBasic2(bucketName string, objectName string) string {
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

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(),
		roleEntityBasic2, roleEntityBasic3_owner)
}

func testGoogleStorageObjectsAclBasic3(bucketName string, objectName string) string {
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

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(),
		roleEntityBasic2, roleEntityBasic3_reader)
}

func testGoogleStorageObjectsAclPredefined(bucketName string, objectName string) string {
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

resource "google_storage_object_acl" "acl" {
  object         = google_storage_bucket_object.object.name
  bucket         = google_storage_bucket.bucket.name
  predefined_acl = "projectPrivate"
}
`, bucketName, objectName, tfObjectAcl.Name())
}

func testGoogleStorageObjectAclUnordered(bucketName, objectName string) string {
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

resource "google_storage_object_acl" "acl" {
  object      = google_storage_bucket_object.object.name
  bucket      = google_storage_bucket.bucket.name
  role_entity = ["%s", "%s", "%s", "%s", "%s"]
}
`, bucketName, objectName, tfObjectAcl.Name(), roleEntityBasic1, roleEntityViewers, roleEntityOwners, roleEntityBasic2, roleEntityEditors)
}
