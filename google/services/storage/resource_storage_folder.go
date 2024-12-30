// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package storage

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceStorageFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageFolderCreate,
		Read:   resourceStorageFolderRead,
		Update: resourceStorageFolderUpdate,
		Delete: resourceStorageFolderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceStorageFolderImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the bucket that contains the folder.`,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidateRegexp(`/$`),
				Description: `The name of the folder expressed as a path. Must include
trailing '/'. For example, 'example_dir/example_dir2/', 'example@#/', 'a-b/d-f/'.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The timestamp at which this folder was created.`,
			},
			"metageneration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The metadata generation of the folder.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The timestamp at which this folder was most recently updated.`,
			},
			"force_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `If set to true, items within folder if any will be force destroyed.`,
				Default:     false,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceStorageFolderCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	nameProp, err := expandStorageFolderName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/folders")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Folder: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating Folder: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{bucket}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Folder %q: %#v", d.Id(), res)

	return resourceStorageFolderRead(d, meta)
}

func resourceStorageFolderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/folders/{{%name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("StorageFolder %q", d.Id()))
	}

	// Explicitly set virtual fields to default values if unset
	if _, ok := d.GetOkExists("force_destroy"); !ok {
		if err := d.Set("force_destroy", false); err != nil {
			return fmt.Errorf("Error setting force_destroy: %s", err)
		}
	}

	if err := d.Set("create_time", flattenStorageFolderCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Folder: %s", err)
	}
	if err := d.Set("update_time", flattenStorageFolderUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Folder: %s", err)
	}
	if err := d.Set("metageneration", flattenStorageFolderMetageneration(res["metageneration"], d, config)); err != nil {
		return fmt.Errorf("Error reading Folder: %s", err)
	}
	if err := d.Set("name", flattenStorageFolderName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Folder: %s", err)
	}
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Folder: %s", err)
	}

	return nil
}

func resourceStorageFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	_ = config
	// we can only get here if force_destroy was updated
	if d.Get("force_destroy") != nil {
		if err := d.Set("force_destroy", d.Get("force_destroy")); err != nil {
			return fmt.Errorf("Error updating force_destroy: %s", err)
		}
	}

	// all other fields are immutable, don't do anything else
	return nil
}

func resourceStorageFolderDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	var listError, deleteObjectError error
	for deleteObjectError == nil {
		res, err := config.NewStorageClient(userAgent).Objects.List(bucket).Prefix(name).Do()
		if err != nil {
			log.Printf("Error listing contents of folder %s: %v", bucket, err)
			listError = err
			break
		}

		if len(res.Items) == 0 {
			break // 0 items, folder empty
		}

		if !d.Get("force_destroy").(bool) {
			deleteErr := fmt.Errorf("Error trying to delete folder %s containing objects without force_destroy set to true", bucket)
			log.Printf("Error! %s : %s\n\n", bucket, deleteErr)
			return deleteErr
		}
		// GCS requires that a folder be empty (have no objects or object
		// versions) before it can be deleted.
		log.Printf("[DEBUG] GCS Folder attempting to forceDestroy\n\n")

		// Create a workerpool for parallel deletion of resources. In the
		// future, it would be great to expose Terraform's global parallelism
		// flag here, but that's currently reserved for core use. Testing
		// shows that NumCPUs-1 is the most performant on average networks.
		//
		// The challenge with making this user-configurable is that the
		// configuration would reside in the Terraform configuration file,
		// decreasing its portability. Ideally we'd want this to connect to
		// Terraform's top-level -parallelism flag, but that's not plumbed nor
		// is it scheduled to be plumbed to individual providers.
		wp := workerpool.New(runtime.NumCPU() - 1)

		for _, object := range res.Items {
			log.Printf("[DEBUG] Found %s", object.Name)
			object := object

			wp.Submit(func() {
				log.Printf("[TRACE] Attempting to delete %s", object.Name)
				if err := config.NewStorageClient(userAgent).Objects.Delete(bucket, object.Name).Generation(object.Generation).Do(); err != nil {
					deleteObjectError = err
					log.Printf("[ERR] Failed to delete storage object %s: %s", object.Name, err)
				} else {
					log.Printf("[TRACE] Successfully deleted %s", object.Name)
				}
			})
		}

		// Wait for everything to finish.
		wp.StopWait()
	}

	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		err := config.NewStorageClient(userAgent).Folders.Delete(bucket, name).Do()
		if err == nil {
			return nil
		}
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 429 {
			return retry.RetryableError(gerr)
		}
		return retry.NonRetryableError(err)
	})

	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 && strings.Contains(gerr.Message, "not empty") && listError != nil {
		return fmt.Errorf("could not delete non-empty folder due to error when listing contents: %v", listError)
	}
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 && strings.Contains(gerr.Message, "not empty") && deleteObjectError != nil {
		return fmt.Errorf("could not delete non-empty folder due to error when deleting contents: %v", deleteObjectError)
	}
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 && strings.Contains(gerr.Message, "not empty") && !d.Get("force_destroy").(bool) {
		return fmt.Errorf("Sub folders or items may exist within folder, use force_destroy to true to delete all subfolders: %v", err)
	}

	if err == nil {
		log.Printf("[DEBUG] Deleted empty folder %v\n\n", name)
		return nil
	} else {
		log.Printf("[ERROR] Error deleting folder %v, %v\n\n", name, err)
	}

	// attempts to delete any sub folders within the folder
	foldersList, err := config.NewStorageClient(userAgent).Folders.List(bucket).Prefix(name).Do()
	if err != nil {
		return err
	}
	if d.Get("force_destroy").(bool) {
		log.Printf("[DEBUG] folder names to delete: %#v", name)
		items := foldersList.Items
		for i := len(items) - 1; i >= 0; i-- {
			err = transport_tpg.Retry(transport_tpg.RetryOptions{
				RetryFunc: func() error {
					err = config.NewStorageClient(userAgent).Folders.Delete(bucket, items[i].Name).Do()
					return err
				},
				Timeout:              d.Timeout(schema.TimeoutDelete),
				ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
			})
			if err != nil {
				return err
			}
		}

		log.Printf("[DEBUG] Finished deleting Folder %q: %#v", d.Id(), name)
	} else {
		deleteErr := fmt.Errorf("Sub folders exist within folder, use force_destroy to true to delete all subfolders")
		log.Printf("Error! %s : %s\n\n", name, deleteErr)
		return deleteErr
	}
	return nil
}

func resourceStorageFolderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^(?P<bucket>[^/]+)/folders/(?P<name>.+)$",
		"^(?P<bucket>[^/]+)/(?P<name>.+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{bucket}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	if err := d.Set("force_destroy", false); err != nil {
		return nil, fmt.Errorf("Error setting force_destroy: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenStorageFolderCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageFolderUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageFolderMetageneration(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenStorageFolderName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandStorageFolderName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}