package google

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleStorageBucketObject() *schema.Resource {

	dsSchema := datasourceSchemaFromResourceSchema(resourceStorageBucketObject().Schema)

	addOptionalFieldsToSchema(dsSchema, "bucket")
	addOptionalFieldsToSchema(dsSchema, "name")

	return &schema.Resource{
		Read:   dataSourceGoogleStorageBucketObjectRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleStorageBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	// URL encode folder names, but to ensure backward compatibility don't url encode
	// them if they were already encoded manually in config.
	// see https://github.com/hashicorp/terraform-provider-google/issues/3176
	if strings.Contains(name, "/") {
		name = url.QueryEscape(name)
	}
	// Using REST apis because the storage go client doesn't support folders
	url := fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o/%s", bucket, name)

	res, err := sendRequest(config, "GET", "", url, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving storage bucket object: %s", err)
	}

	if err := d.Set("cache_control", res["cacheControl"]); err != nil {
		return fmt.Errorf("Error reading cache_control: %s", err)
	}
	if err := d.Set("content_disposition", res["contentDisposition"]); err != nil {
		return fmt.Errorf("Error reading content_disposition: %s", err)
	}
	if err := d.Set("content_encoding", res["contentEncoding"]); err != nil {
		return fmt.Errorf("Error reading content_encoding: %s", err)
	}
	if err := d.Set("content_language", res["contentLanguage"]); err != nil {
		return fmt.Errorf("Error reading content_language: %s", err)
	}
	if err := d.Set("content_type", res["contentType"]); err != nil {
		return fmt.Errorf("Error reading content_type: %s", err)
	}
	if err := d.Set("crc32c", res["crc32c"]); err != nil {
		return fmt.Errorf("Error reading crc32c: %s", err)
	}
	if err := d.Set("self_link", res["selfLink"]); err != nil {
		return fmt.Errorf("Error reading self_link: %s", err)
	}
	if err := d.Set("storage_class", res["storageClass"]); err != nil {
		return fmt.Errorf("Error reading storage_class: %s", err)
	}
	if err := d.Set("md5hash", res["md5Hash"]); err != nil {
		return fmt.Errorf("Error reading md5hash: %s", err)
	}
	if err := d.Set("media_link", res["mediaLink"]); err != nil {
		return fmt.Errorf("Error reading media_link: %s", err)
	}
	if err := d.Set("metadata", res["metadata"]); err != nil {
		return fmt.Errorf("Error reading metadata: %s", err)
	}

	d.SetId(bucket + "-" + name)

	return nil
}
