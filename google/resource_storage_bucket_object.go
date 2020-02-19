package google

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"crypto/md5"
	"encoding/base64"
	"io/ioutil"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

func resourceStorageBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageBucketObjectCreate,
		Read:   resourceStorageBucketObjectRead,
		Delete: resourceStorageBucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cache_control": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"content_language": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source"},
				Sensitive:     true,
			},

			"crc32c": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"md5hash": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"content"},
			},

			// Detect changes to local file or changes made outside of Terraform to the file stored on the server.
			"detect_md5hash": {
				Type: schema.TypeString,
				// This field is not Computed because it needs to trigger a diff.
				Optional: true,
				ForceNew: true,
				// Makes the diff message nicer:
				// detect_md5hash:       "1XcnP/iFw/hNrbhXi7QTmQ==" => "different hash" (forces new resource)
				// Instead of the more confusing:
				// detect_md5hash:       "1XcnP/iFw/hNrbhXi7QTmQ==" => "" (forces new resource)
				Default: "different hash",
				// 1. Compute the md5 hash of the local file
				// 2. Compare the computed md5 hash with the hash stored in Cloud Storage
				// 3. Don't suppress the diff iff they don't match
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					localMd5Hash := ""
					if source, ok := d.GetOkExists("source"); ok {
						localMd5Hash = getFileMd5Hash(source.(string))
					}

					if content, ok := d.GetOkExists("content"); ok {
						localMd5Hash = getContentMd5Hash([]byte(content.(string)))
					}

					// If `source` or `content` is dynamically set, both field will be empty.
					// We should not suppress the diff to avoid the following error:
					// 'Mismatch reason: extra attributes: detect_md5hash'
					if localMd5Hash == "" {
						return false
					}

					// `old` is the md5 hash we retrieved from the server in the ReadFunc
					if old != localMd5Hash {
						return false
					}

					return true
				},
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// https://github.com/hashicorp/terraform/issues/19052
			"output_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func objectGetId(object *storage.Object) string {
	return object.Bucket + "-" + object.Name
}

func resourceStorageBucketObjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	var media io.Reader

	if v, ok := d.GetOk("source"); ok {
		var err error
		media, err = os.Open(v.(string))
		if err != nil {
			return err
		}
	} else if v, ok := d.GetOk("content"); ok {
		media = bytes.NewReader([]byte(v.(string)))
	} else {
		return fmt.Errorf("Error, either \"content\" or \"source\" must be specified")
	}

	objectsService := storage.NewObjectsService(config.clientStorage)
	object := &storage.Object{Bucket: bucket}

	if v, ok := d.GetOk("cache_control"); ok {
		object.CacheControl = v.(string)
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		object.ContentDisposition = v.(string)
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		object.ContentEncoding = v.(string)
	}

	if v, ok := d.GetOk("content_language"); ok {
		object.ContentLanguage = v.(string)
	}

	if v, ok := d.GetOk("content_type"); ok {
		object.ContentType = v.(string)
	}

	if v, ok := d.GetOk("storage_class"); ok {
		object.StorageClass = v.(string)
	}

	insertCall := objectsService.Insert(bucket, object)
	insertCall.Name(name)
	insertCall.Media(media)

	_, err := insertCall.Do()

	if err != nil {
		return fmt.Errorf("Error uploading object %s: %s", name, err)
	}

	return resourceStorageBucketObjectRead(d, meta)
}

func resourceStorageBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	objectsService := storage.NewObjectsService(config.clientStorage)
	getCall := objectsService.Get(bucket, name)

	res, err := getCall.Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Storage Bucket Object %q", d.Get("name").(string)))
	}

	d.Set("md5hash", res.Md5Hash)
	d.Set("detect_md5hash", res.Md5Hash)
	d.Set("crc32c", res.Crc32c)
	d.Set("cache_control", res.CacheControl)
	d.Set("content_disposition", res.ContentDisposition)
	d.Set("content_encoding", res.ContentEncoding)
	d.Set("content_language", res.ContentLanguage)
	d.Set("content_type", res.ContentType)
	d.Set("storage_class", res.StorageClass)
	d.Set("self_link", res.SelfLink)
	d.Set("output_name", res.Name)

	d.SetId(objectGetId(res))

	return nil
}

func resourceStorageBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	objectsService := storage.NewObjectsService(config.clientStorage)

	DeleteCall := objectsService.Delete(bucket, name)
	err := DeleteCall.Do()

	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing Bucket Object %q because it's gone", name)
			// The resource doesn't exist anymore
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error deleting contents of object %s: %s", name, err)
	}

	return nil
}

func getFileMd5Hash(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("[WARN] Failed to read source file %q. Cannot compute md5 hash for it.", filename)
		return ""
	}

	return getContentMd5Hash(data)
}

func getContentMd5Hash(content []byte) string {
	h := md5.New()
	if _, err := h.Write(content); err != nil {
		log.Printf("[WARN] Failed to compute md5 hash for content: %v", err)
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
