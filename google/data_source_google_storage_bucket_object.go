package google

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	// see https://github.com/terraform-providers/terraform-provider-google/issues/3176
	if strings.Contains(name, "/") {
		name = url.QueryEscape(name)
	}
	// Using REST apis because the storage go client doesn't support folders
	url := fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o/%s", bucket, name)

	res, err := sendRequest(config, "GET", "", url, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving storage bucket object: %s", err)
	}

	d.Set("cache_control", res["cacheControl"])
	d.Set("content_disposition", res["contentDisposition"])
	d.Set("content_encoding", res["contentEncoding"])
	d.Set("content_language", res["contentLanguage"])
	d.Set("content_type", res["contentType"])
	d.Set("crc32c", res["crc32c"])
	d.Set("self_link", res["selfLink"])
	d.Set("storage_class", res["storageClass"])
	d.Set("md5hash", res["md5Hash"])

	d.SetId(bucket + "-" + name)

	return nil
}
