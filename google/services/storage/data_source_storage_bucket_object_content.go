// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
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
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "content")

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
	var bodyString string

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading all  from res.Body: %s", err)
		}
		bodyString = string(bodyBytes)
	}

	if err := d.Set("content", bodyString); err != nil {
		return fmt.Errorf("Error setting content: %s", err)
	}

	d.SetId(bucket + "-" + name)
	return nil
}
