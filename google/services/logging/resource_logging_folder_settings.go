// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/logging/FolderSettings.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package logging

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceLoggingFolderSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingFolderSettingsCreate,
		Read:   resourceLoggingFolderSettingsRead,
		Update: resourceLoggingFolderSettingsUpdate,
		Delete: resourceLoggingFolderSettingsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceLoggingFolderSettingsImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"folder": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The folder for which to retrieve settings.`,
			},
			"disable_default_sink": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: `If set to true, the _Default sink in newly created projects and folders will created in a disabled state. This can be used to automatically disable log storage if there is already an aggregated sink configured in the hierarchy. The _Default sink can be re-enabled manually if needed.`,
			},
			"kms_key_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `The resource name for the configured Cloud KMS key.`,
			},
			"storage_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `The storage location that Cloud Logging will use to create new resources when a location is needed but not explicitly provided.`,
			},
			"kms_service_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service account that will be used by the Log Router to access your Cloud KMS key.`,
			},
			"logging_service_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service account for the given container. Sinks use this service account as their writerIdentity if no custom service account is provided.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name of the settings.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceLoggingFolderSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	kmsKeyNameProp, err := expandLoggingFolderSettingsKmsKeyName(d.Get("kms_key_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("kms_key_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(kmsKeyNameProp)) && (ok || !reflect.DeepEqual(v, kmsKeyNameProp)) {
		obj["kmsKeyName"] = kmsKeyNameProp
	}
	storageLocationProp, err := expandLoggingFolderSettingsStorageLocation(d.Get("storage_location"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("storage_location"); !tpgresource.IsEmptyValue(reflect.ValueOf(storageLocationProp)) && (ok || !reflect.DeepEqual(v, storageLocationProp)) {
		obj["storageLocation"] = storageLocationProp
	}
	disableDefaultSinkProp, err := expandLoggingFolderSettingsDisableDefaultSink(d.Get("disable_default_sink"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disable_default_sink"); !tpgresource.IsEmptyValue(reflect.ValueOf(disableDefaultSinkProp)) && (ok || !reflect.DeepEqual(v, disableDefaultSinkProp)) {
		obj["disableDefaultSink"] = disableDefaultSinkProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{LoggingBasePath}}folders/{{folder}}/settings?updateMask=disableDefaultSink,storageLocation,kmsKeyName")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new FolderSettings: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating FolderSettings: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "folders/{{folder}}/settings")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating FolderSettings %q: %#v", d.Id(), res)

	return resourceLoggingFolderSettingsRead(d, meta)
}

func resourceLoggingFolderSettingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{LoggingBasePath}}folders/{{folder}}/settings")
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("LoggingFolderSettings %q", d.Id()))
	}

	if err := d.Set("name", flattenLoggingFolderSettingsName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}
	if err := d.Set("kms_key_name", flattenLoggingFolderSettingsKmsKeyName(res["kmsKeyName"], d, config)); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}
	if err := d.Set("kms_service_account_id", flattenLoggingFolderSettingsKmsServiceAccountId(res["kmsServiceAccountId"], d, config)); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}
	if err := d.Set("storage_location", flattenLoggingFolderSettingsStorageLocation(res["storageLocation"], d, config)); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}
	if err := d.Set("disable_default_sink", flattenLoggingFolderSettingsDisableDefaultSink(res["disableDefaultSink"], d, config)); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}
	if err := d.Set("logging_service_account_id", flattenLoggingFolderSettingsLoggingServiceAccountId(res["loggingServiceAccountId"], d, config)); err != nil {
		return fmt.Errorf("Error reading FolderSettings: %s", err)
	}

	return nil
}

func resourceLoggingFolderSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	kmsKeyNameProp, err := expandLoggingFolderSettingsKmsKeyName(d.Get("kms_key_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("kms_key_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, kmsKeyNameProp)) {
		obj["kmsKeyName"] = kmsKeyNameProp
	}
	storageLocationProp, err := expandLoggingFolderSettingsStorageLocation(d.Get("storage_location"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("storage_location"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, storageLocationProp)) {
		obj["storageLocation"] = storageLocationProp
	}
	disableDefaultSinkProp, err := expandLoggingFolderSettingsDisableDefaultSink(d.Get("disable_default_sink"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disable_default_sink"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, disableDefaultSinkProp)) {
		obj["disableDefaultSink"] = disableDefaultSinkProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{LoggingBasePath}}folders/{{folder}}/settings?updateMask=disableDefaultSink,storageLocation,kmsKeyName")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating FolderSettings %q: %#v", d.Id(), obj)
	headers := make(http.Header)

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
		Headers:   headers,
	})

	if err != nil {
		return fmt.Errorf("Error updating FolderSettings %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating FolderSettings %q: %#v", d.Id(), res)
	}

	return resourceLoggingFolderSettingsRead(d, meta)
}

func resourceLoggingFolderSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] Logging FolderSettings resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceLoggingFolderSettingsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^folders/(?P<folder>[^/]+)/settings$",
		"^(?P<folder>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "folders/{{folder}}/settings")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenLoggingFolderSettingsName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingFolderSettingsKmsKeyName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingFolderSettingsKmsServiceAccountId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingFolderSettingsStorageLocation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingFolderSettingsDisableDefaultSink(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenLoggingFolderSettingsLoggingServiceAccountId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandLoggingFolderSettingsKmsKeyName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLoggingFolderSettingsStorageLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandLoggingFolderSettingsDisableDefaultSink(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
