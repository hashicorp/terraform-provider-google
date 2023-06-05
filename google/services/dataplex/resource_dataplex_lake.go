// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

package dataplex

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	dataplex "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataplex"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceDataplexLake() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataplexLakeCreate,
		Read:   resourceDataplexLakeRead,
		Update: resourceDataplexLakeUpdate,
		Delete: resourceDataplexLakeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDataplexLakeImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The name of the lake.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Description of the lake.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. User friendly display name.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. User-defined labels for the lake.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"metastore": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Settings to manage lake and Dataproc Metastore service instance association.",
				MaxItems:    1,
				Elem:        DataplexLakeMetastoreSchema(),
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"asset_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Aggregated status of the underlying assets of the lake.",
				Elem:        DataplexLakeAssetStatusSchema(),
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when the lake was created.",
			},

			"metastore_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Metastore status of the lake.",
				Elem:        DataplexLakeMetastoreStatusSchema(),
			},

			"service_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Service account associated with this lake. This service account must be authorized to access or operate on resources managed by the lake.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Current state of the lake. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. System generated globally unique ID for the lake. This ID will be different if the lake is deleted and re-created with the same name.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time when the lake was last updated.",
			},
		},
	}
}

func DataplexLakeMetastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. A relative reference to the Dataproc Metastore (https://cloud.google.com/dataproc-metastore/docs) service associated with the lake: `projects/{project_id}/locations/{location_id}/services/{service_id}`",
			},
		},
	}
}

func DataplexLakeAssetStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"active_assets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of active assets.",
			},

			"security_policy_applying_assets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of assets that are in process of updating the security policy on attached resources.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time of the status.",
			},
		},
	}
}

func DataplexLakeMetastoreStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI of the endpoint used to access the Metastore service.",
			},

			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the current status.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current state of association. Possible values: STATE_UNSPECIFIED, NONE, READY, UPDATING, ERROR",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time of the metastore status of the lake.",
			},
		},
	}
}

func resourceDataplexLakeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Lake{
		Location:    dcl.String(d.Get("location").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		DisplayName: dcl.String(d.Get("display_name").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Metastore:   expandDataplexLakeMetastore(d.Get("metastore")),
		Project:     dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyLake(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Lake: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Lake %q: %#v", d.Id(), res)

	return resourceDataplexLakeRead(d, meta)
}

func resourceDataplexLakeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Lake{
		Location:    dcl.String(d.Get("location").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		DisplayName: dcl.String(d.Get("display_name").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Metastore:   expandDataplexLakeMetastore(d.Get("metastore")),
		Project:     dcl.String(project),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetLake(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("DataplexLake %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("display_name", res.DisplayName); err != nil {
		return fmt.Errorf("error setting display_name in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("metastore", flattenDataplexLakeMetastore(res.Metastore)); err != nil {
		return fmt.Errorf("error setting metastore in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("asset_status", flattenDataplexLakeAssetStatus(res.AssetStatus)); err != nil {
		return fmt.Errorf("error setting asset_status in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("metastore_status", flattenDataplexLakeMetastoreStatus(res.MetastoreStatus)); err != nil {
		return fmt.Errorf("error setting metastore_status in state: %s", err)
	}
	if err = d.Set("service_account", res.ServiceAccount); err != nil {
		return fmt.Errorf("error setting service_account in state: %s", err)
	}
	if err = d.Set("state", res.State); err != nil {
		return fmt.Errorf("error setting state in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceDataplexLakeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Lake{
		Location:    dcl.String(d.Get("location").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		DisplayName: dcl.String(d.Get("display_name").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Metastore:   expandDataplexLakeMetastore(d.Get("metastore")),
		Project:     dcl.String(project),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyLake(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Lake: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Lake %q: %#v", d.Id(), res)

	return resourceDataplexLakeRead(d, meta)
}

func resourceDataplexLakeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &dataplex.Lake{
		Location:    dcl.String(d.Get("location").(string)),
		Name:        dcl.String(d.Get("name").(string)),
		Description: dcl.String(d.Get("description").(string)),
		DisplayName: dcl.String(d.Get("display_name").(string)),
		Labels:      tpgresource.CheckStringMap(d.Get("labels")),
		Metastore:   expandDataplexLakeMetastore(d.Get("metastore")),
		Project:     dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting Lake %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLDataplexClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteLake(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Lake: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Lake %q", d.Id())
	return nil
}

func resourceDataplexLakeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/lakes/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/lakes/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandDataplexLakeMetastore(o interface{}) *dataplex.LakeMetastore {
	if o == nil {
		return dataplex.EmptyLakeMetastore
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return dataplex.EmptyLakeMetastore
	}
	obj := objArr[0].(map[string]interface{})
	return &dataplex.LakeMetastore{
		Service: dcl.String(obj["service"].(string)),
	}
}

func flattenDataplexLakeMetastore(obj *dataplex.LakeMetastore) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"service": obj.Service,
	}

	return []interface{}{transformed}

}

func flattenDataplexLakeAssetStatus(obj *dataplex.LakeAssetStatus) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"active_assets":                   obj.ActiveAssets,
		"security_policy_applying_assets": obj.SecurityPolicyApplyingAssets,
		"update_time":                     obj.UpdateTime,
	}

	return []interface{}{transformed}

}

func flattenDataplexLakeMetastoreStatus(obj *dataplex.LakeMetastoreStatus) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"endpoint":    obj.Endpoint,
		"message":     obj.Message,
		"state":       obj.State,
		"update_time": obj.UpdateTime,
	}

	return []interface{}{transformed}

}
