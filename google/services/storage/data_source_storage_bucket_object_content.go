// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/storage/v1"
)

func DataSourceGoogleStorageBucketObjectContent() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceStorageBucketObject().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "bucket")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// The field must be optional for backward compatibility.
	dsSchema["content"].Optional = true
	dsSchema["content_base64"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Base64 encoded version of the object content. Use this when dealing with binary data.",
		Computed:    true,
		Optional:    false,
		Required:    false,
	}

	dsSchema["content_hexsha512"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Hex encoded SHA512 checksum of object content.",
		Computed:    true,
		Optional:    false,
		Required:    false,
	}

	dsSchema["content_base64sha512"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Base64 encoded SHA512 checksum of object content.",
		Computed:    true,
		Optional:    false,
		Required:    false,
	}

	return &schema.Resource{
		Read:   dataSourceGoogleStorageBucketObjectContentRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleStorageBucketObjectContentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	objectsService := storage.NewObjectsService(config.NewStorageClient(userAgent))
	getCall := objectsService.Get(bucket, name)

	res, err := getCall.Download()
	if err != nil {
		return fmt.Errorf("Error downloading storage bucket object: %s", err)
	}

	defer res.Body.Close()
	var objectBytes []byte

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading all  from res.Body: %s", err)
		}
		objectBytes = bodyBytes
	}

	if err := d.Set("content", string(objectBytes)); err != nil {
		return fmt.Errorf("Error setting content: %s", err)
	}

	if err := d.Set("content_base64", base64.StdEncoding.EncodeToString(objectBytes)); err != nil {
		return fmt.Errorf("Error setting content_base64: %s", err)
	}

	sha512Sum := sha512.Sum512(objectBytes)
	if err := d.Set("content_hexsha512", hex.EncodeToString(sha512Sum[:])); err != nil {
		return fmt.Errorf("Error setting content_hexsha512: %s", err)
	}

	if err := d.Set("content_base64sha512", base64.StdEncoding.EncodeToString(sha512Sum[:])); err != nil {
		return fmt.Errorf("Error setting content_base64sha512: %s", err)
	}

	d.SetId(bucket + "-" + name)
	return nil
}
