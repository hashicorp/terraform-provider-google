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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/firestore/Document.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package firestore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceFirestoreDocument() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirestoreDocumentCreate,
		Read:   resourceFirestoreDocumentRead,
		Update: resourceFirestoreDocumentUpdate,
		Delete: resourceFirestoreDocumentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceFirestoreDocumentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"collection": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The collection ID, relative to database. For example: chatrooms or chatrooms/my-document/private-messages.`,
			},
			"document_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The client-assigned document ID to use for this document during creation.`,
			},
			"fields": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				StateFunc:    func(v interface{}) string { s, _ := structure.NormalizeJsonString(v); return s },
				Description:  `The document's [fields](https://cloud.google.com/firestore/docs/reference/rest/v1/projects.databases.documents) formated as a json string.`,
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The Firestore database id. Defaults to '"(default)"'.`,
				Default:     "(default)",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Creation timestamp in RFC3339 format.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `A server defined name for this document. Format:
'projects/{{project_id}}/databases/{{database_id}}/documents/{{path}}/{{document_id}}'`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `A relative path to the collection this document exists within`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Last update timestamp in RFC3339 format.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceFirestoreDocumentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	fieldsProp, err := expandFirestoreDocumentFields(d.Get("fields"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("fields"); !tpgresource.IsEmptyValue(reflect.ValueOf(fieldsProp)) && (ok || !reflect.DeepEqual(v, fieldsProp)) {
		obj["fields"] = fieldsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{FirestoreBasePath}}projects/{{project}}/databases/{{database}}/documents/{{collection}}?documentId={{document_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Document: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Document: %s", err)
	}
	billingProject = project

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
		return fmt.Errorf("Error creating Document: %s", err)
	}
	// Set computed resource properties from create API response so that they're available on the subsequent Read
	// call.
	res, err = resourceFirestoreDocumentDecoder(d, meta, res)
	if err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	if res == nil {
		return fmt.Errorf("decoding response, could not find object")
	}
	if err := d.Set("name", flattenFirestoreDocumentName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Document %q: %#v", d.Id(), res)

	return resourceFirestoreDocumentRead(d, meta)
}

func resourceFirestoreDocumentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{FirestoreBasePath}}{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Document: %s", err)
	}
	billingProject = project

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("FirestoreDocument %q", d.Id()))
	}

	res, err = resourceFirestoreDocumentDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing FirestoreDocument because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Document: %s", err)
	}

	if err := d.Set("name", flattenFirestoreDocumentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Document: %s", err)
	}
	if err := d.Set("path", flattenFirestoreDocumentPath(res["path"], d, config)); err != nil {
		return fmt.Errorf("Error reading Document: %s", err)
	}
	if err := d.Set("fields", flattenFirestoreDocumentFields(res["fields"], d, config)); err != nil {
		return fmt.Errorf("Error reading Document: %s", err)
	}
	if err := d.Set("create_time", flattenFirestoreDocumentCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Document: %s", err)
	}
	if err := d.Set("update_time", flattenFirestoreDocumentUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading Document: %s", err)
	}

	return nil
}

func resourceFirestoreDocumentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Document: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	fieldsProp, err := expandFirestoreDocumentFields(d.Get("fields"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("fields"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, fieldsProp)) {
		obj["fields"] = fieldsProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{FirestoreBasePath}}{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Document %q: %#v", d.Id(), obj)
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
		return fmt.Errorf("Error updating Document %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Document %q: %#v", d.Id(), res)
	}

	return resourceFirestoreDocumentRead(d, meta)
}

func resourceFirestoreDocumentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Document: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{FirestoreBasePath}}{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting Document %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Document")
	}

	log.Printf("[DEBUG] Finished deleting Document %q: %#v", d.Id(), res)
	return nil
}

func resourceFirestoreDocumentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	config := meta.(*transport_tpg.Config)

	// current import_formats can't import fields with forward slashes in their value
	if err := tpgresource.ParseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
		return nil, err
	}

	re := regexp.MustCompile("^projects/([^/]+)/databases/([^/]+)/documents/(.+)/([^/]+)$")
	match := re.FindStringSubmatch(d.Get("name").(string))
	if len(match) > 0 {
		if err := d.Set("project", match[1]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("database", match[2]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("collection", match[3]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("document_id", match[4]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
	} else {
		return nil, fmt.Errorf("import did not match the regex ^projects/([^/]+)/databases/([^/]+)/documents/(.+)/([^/]+)$")
	}

	return []*schema.ResourceData{d}, nil
}

func flattenFirestoreDocumentName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirestoreDocumentPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirestoreDocumentFields(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		// TODO: return error once https://github.com/GoogleCloudPlatform/magic-modules/issues/3257 is fixed.
		log.Printf("[ERROR] failed to marshal schema to JSON: %v", err)
	}
	return string(b)
}

func flattenFirestoreDocumentCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirestoreDocumentUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandFirestoreDocumentFields(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	b := []byte(v.(string))
	if len(b) == 0 {
		return nil, nil
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func resourceFirestoreDocumentDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	// We use this decoder to add the path field
	if name, ok := res["name"]; ok {
		re := regexp.MustCompile("^projects/[^/]+/databases/[^/]+/documents/(.+)$")
		match := re.FindStringSubmatch(name.(string))
		if len(match) > 0 {
			res["path"] = match[1]
		}
	}
	return res, nil
}
