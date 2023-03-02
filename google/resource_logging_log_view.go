// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	logging "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/logging"
)

func ResourceLoggingLogView() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingLogViewCreate,
		Read:   resourceLoggingLogViewRead,
		Update: resourceLoggingLogViewUpdate,
		Delete: resourceLoggingLogViewDelete,

		Importer: &schema.ResourceImporter{
			State: resourceLoggingLogViewImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The bucket of the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource name of the view. For example: `projects/my-project/locations/global/buckets/my-bucket/views/my-view`",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Describes this view.",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter that restricts which log entries in a bucket are visible in this view. Filters are restricted to be a logical AND of ==/!= of any of the following: - originating project/folder/organization/billing account. - resource type - log id For example: SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\" AND LOG_ID(\"stdout\")",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of the resource. The supported locations are: global, us-central1, us-east1, us-west1, asia-east1, europe-west1.",
			},

			"parent": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The parent of the resource.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The creation timestamp of the view.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The last update timestamp of the view.",
			},
		},
	}
}

func resourceLoggingLogViewCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &logging.LogView{
		Bucket:      dcl.String(d.Get("bucket").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Filter:      dcl.String(d.Get("filter").(string)),
		Location:    dcl.StringOrNil(d.Get("location").(string)),
		Parent:      dcl.StringOrNil(d.Get("parent").(string)),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := CreateDirective
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLLoggingClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyLogView(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating LogView: %s", err)
	}

	log.Printf("[DEBUG] Finished creating LogView %q: %#v", d.Id(), res)

	return resourceLoggingLogViewRead(d, meta)
}

func resourceLoggingLogViewRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &logging.LogView{
		Bucket:      dcl.String(d.Get("bucket").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Filter:      dcl.String(d.Get("filter").(string)),
		Location:    dcl.StringOrNil(d.Get("location").(string)),
		Parent:      dcl.StringOrNil(d.Get("parent").(string)),
	}

	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLLoggingClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetLogView(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("LoggingLogView %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("bucket", res.Bucket); err != nil {
		return fmt.Errorf("error setting bucket in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("filter", res.Filter); err != nil {
		return fmt.Errorf("error setting filter in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("parent", res.Parent); err != nil {
		return fmt.Errorf("error setting parent in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceLoggingLogViewUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &logging.LogView{
		Bucket:      dcl.String(d.Get("bucket").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Filter:      dcl.String(d.Get("filter").(string)),
		Location:    dcl.StringOrNil(d.Get("location").(string)),
		Parent:      dcl.StringOrNil(d.Get("parent").(string)),
	}
	directive := UpdateDirective
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLLoggingClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyLogView(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating LogView: %s", err)
	}

	log.Printf("[DEBUG] Finished creating LogView %q: %#v", d.Id(), res)

	return resourceLoggingLogViewRead(d, meta)
}

func resourceLoggingLogViewDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := &logging.LogView{
		Bucket:      dcl.String(d.Get("bucket").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		Filter:      dcl.String(d.Get("filter").(string)),
		Location:    dcl.StringOrNil(d.Get("location").(string)),
		Parent:      dcl.StringOrNil(d.Get("parent").(string)),
	}

	log.Printf("[DEBUG] Deleting LogView %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLLoggingClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteLogView(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting LogView: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting LogView %q", d.Id())
	return nil
}

func resourceLoggingLogViewImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"(?P<parent>.+)/locations/(?P<location>.+)/buckets/(?P<bucket>.+)/views/(?P<name>.+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{parent}}/locations/{{location}}/buckets/{{bucket}}/views/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
