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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/firebaseapphosting/DefaultDomain.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package firebaseapphosting

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceFirebaseAppHostingDefaultDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirebaseAppHostingDefaultDomainCreate,
		Read:   resourceFirebaseAppHostingDefaultDomainRead,
		Update: resourceFirebaseAppHostingDefaultDomainUpdate,
		Delete: resourceFirebaseAppHostingDefaultDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceFirebaseAppHostingDefaultDomainImport,
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
			"backend": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the Backend that this Domain is associated with`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Id of the domain. For default domain, it should be {{backend}}--{{project_id}}.{{location}}.hosted.app`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The location of the Backend that this Domain is associated with`,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: `Whether the domain is disabled. Defaults to false.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time at which the domain was created.`,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Server-computed checksum based on other values; may be sent
on update or delete to ensure operation is done on expected resource.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Identifier. The resource name of the domain, e.g.
'projects/{project}/locations/{locationId}/backends/{backendId}/domains/{domainId}'`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `System-assigned, unique identifier.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time at which the domain was last updated.`,
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

func resourceFirebaseAppHostingDefaultDomainCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	disabledProp, err := expandFirebaseAppHostingDefaultDomainDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(disabledProp)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{FirebaseAppHostingBasePath}}projects/{{project}}/locations/{{location}}/backends/{{backend}}/domains/{{domain_id}}?update_mask=disabled")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new DefaultDomain: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DefaultDomain: %s", err)
	}
	billingProject = project

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
		return fmt.Errorf("Error creating DefaultDomain: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/backends/{{backend}}/domains/{{domain_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = FirebaseAppHostingOperationWaitTime(
		config, res, project, "Creating DefaultDomain", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create DefaultDomain: %s", err)
	}

	log.Printf("[DEBUG] Finished creating DefaultDomain %q: %#v", d.Id(), res)

	return resourceFirebaseAppHostingDefaultDomainRead(d, meta)
}

func resourceFirebaseAppHostingDefaultDomainRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{FirebaseAppHostingBasePath}}projects/{{project}}/locations/{{location}}/backends/{{backend}}/domains/{{domain_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DefaultDomain: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("FirebaseAppHostingDefaultDomain %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading DefaultDomain: %s", err)
	}

	if err := d.Set("disabled", flattenFirebaseAppHostingDefaultDomainDisabled(res["disabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultDomain: %s", err)
	}
	if err := d.Set("name", flattenFirebaseAppHostingDefaultDomainName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultDomain: %s", err)
	}
	if err := d.Set("uid", flattenFirebaseAppHostingDefaultDomainUid(res["uid"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultDomain: %s", err)
	}
	if err := d.Set("etag", flattenFirebaseAppHostingDefaultDomainEtag(res["etag"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultDomain: %s", err)
	}
	if err := d.Set("update_time", flattenFirebaseAppHostingDefaultDomainUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultDomain: %s", err)
	}
	if err := d.Set("create_time", flattenFirebaseAppHostingDefaultDomainCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultDomain: %s", err)
	}

	return nil
}

func resourceFirebaseAppHostingDefaultDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for DefaultDomain: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	disabledProp, err := expandFirebaseAppHostingDefaultDomainDisabled(d.Get("disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disabled"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, disabledProp)) {
		obj["disabled"] = disabledProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{FirebaseAppHostingBasePath}}projects/{{project}}/locations/{{location}}/backends/{{backend}}/domains/{{domain_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating DefaultDomain %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("disabled") {
		updateMask = append(updateMask, "disabled")
	}
	// updateMask is a URL parameter but not present in the schema, so ReplaceVars
	// won't set it
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// if updateMask is empty we are not updating anything so skip the post
	if len(updateMask) > 0 {
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
			return fmt.Errorf("Error updating DefaultDomain %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating DefaultDomain %q: %#v", d.Id(), res)
		}

		err = FirebaseAppHostingOperationWaitTime(
			config, res, project, "Updating DefaultDomain", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceFirebaseAppHostingDefaultDomainRead(d, meta)
}

func resourceFirebaseAppHostingDefaultDomainDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] FirebaseAppHosting DefaultDomain resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceFirebaseAppHostingDefaultDomainImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/backends/(?P<backend>[^/]+)/domains/(?P<domain_id>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<backend>[^/]+)/(?P<domain_id>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<backend>[^/]+)/(?P<domain_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/backends/{{backend}}/domains/{{domain_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenFirebaseAppHostingDefaultDomainDisabled(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirebaseAppHostingDefaultDomainName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirebaseAppHostingDefaultDomainUid(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirebaseAppHostingDefaultDomainEtag(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirebaseAppHostingDefaultDomainUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenFirebaseAppHostingDefaultDomainCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandFirebaseAppHostingDefaultDomainDisabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
