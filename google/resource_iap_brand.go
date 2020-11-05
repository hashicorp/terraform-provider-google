// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
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

package google

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIapBrand() *schema.Resource {
	return &schema.Resource{
		Create: resourceIapBrandCreate,
		Read:   resourceIapBrandRead,
		Delete: resourceIapBrandDelete,

		Importer: &schema.ResourceImporter{
			State: resourceIapBrandImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"application_title": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Application name displayed on OAuth consent screen.`,
			},
			"support_email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Support email displayed on the OAuth consent screen. Can be either a
user or group email. When a user email is specified, the caller must
be the user with the associated email address. When a group email is
specified, the caller can be either a user or a service account which
is an owner of the specified group in Cloud Identity.`,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Output only. Identifier of the brand, in the format
'projects/{project_number}/brands/{brand_id}'. NOTE: The brand
identification corresponds to the project number as only one
brand per project can be created.`,
			},
			"org_internal_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the brand is only intended for usage inside the GSuite organization only.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIapBrandCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	supportEmailProp, err := expandIapBrandSupportEmail(d.Get("support_email"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("support_email"); !isEmptyValue(reflect.ValueOf(supportEmailProp)) && (ok || !reflect.DeepEqual(v, supportEmailProp)) {
		obj["supportEmail"] = supportEmailProp
	}
	applicationTitleProp, err := expandIapBrandApplicationTitle(d.Get("application_title"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("application_title"); !isEmptyValue(reflect.ValueOf(applicationTitleProp)) && (ok || !reflect.DeepEqual(v, applicationTitleProp)) {
		obj["applicationTitle"] = applicationTitleProp
	}

	url, err := replaceVars(d, config, "{{IapBasePath}}projects/{{project}}/brands")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Brand: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Brand: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Brand: %s", err)
	}
	if err := d.Set("name", flattenIapBrandName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = PollingWaitTime(resourceIapBrandPollRead(d, meta), PollCheckForExistence, "Creating Brand", d.Timeout(schema.TimeoutCreate), 5)
	if err != nil {
		return fmt.Errorf("Error waiting to create Brand: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Brand %q: %#v", d.Id(), res)

	// `name` is autogenerated from the api so needs to be set post-create
	name, ok := res["name"]
	if !ok {
		respBody, ok := res["response"]
		if !ok {
			return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
		}

		name, ok = respBody.(map[string]interface{})["name"]
		if !ok {
			return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
		}
	}
	if err := d.Set("name", name.(string)); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	d.SetId(name.(string))

	return resourceIapBrandRead(d, meta)
}

func resourceIapBrandPollRead(d *schema.ResourceData, meta interface{}) PollReadFunc {
	return func() (map[string]interface{}, error) {
		config := meta.(*Config)

		url, err := replaceVars(d, config, "{{IapBasePath}}{{name}}")
		if err != nil {
			return nil, err
		}

		billingProject := ""

		project, err := getProject(d, config)
		if err != nil {
			return nil, fmt.Errorf("Error fetching project for Brand: %s", err)
		}
		billingProject = project

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return nil, err
		}

		res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
		if err != nil {
			return res, err
		}
		return res, nil
	}
}

func resourceIapBrandRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{IapBasePath}}{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Brand: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("IapBrand %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Brand: %s", err)
	}

	if err := d.Set("support_email", flattenIapBrandSupportEmail(res["supportEmail"], d, config)); err != nil {
		return fmt.Errorf("Error reading Brand: %s", err)
	}
	if err := d.Set("application_title", flattenIapBrandApplicationTitle(res["applicationTitle"], d, config)); err != nil {
		return fmt.Errorf("Error reading Brand: %s", err)
	}
	if err := d.Set("org_internal_only", flattenIapBrandOrgInternalOnly(res["orgInternalOnly"], d, config)); err != nil {
		return fmt.Errorf("Error reading Brand: %s", err)
	}
	if err := d.Set("name", flattenIapBrandName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Brand: %s", err)
	}

	return nil
}

func resourceIapBrandDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] Iap Brand resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceIapBrandImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	// current import_formats can't import fields with forward slashes in their value
	if err := parseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
		return nil, err
	}

	nameParts := strings.Split(d.Get("name").(string), "/")
	if len(nameParts) != 4 {
		return nil, fmt.Errorf(
			"Saw %s when the name is expected to have shape %s",
			d.Get("name"),
			"projects/{{project}}/brands/{{name}}",
		)
	}

	if err := d.Set("project", nameParts[1]); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func flattenIapBrandSupportEmail(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapBrandApplicationTitle(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapBrandOrgInternalOnly(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapBrandName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandIapBrandSupportEmail(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapBrandApplicationTitle(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
