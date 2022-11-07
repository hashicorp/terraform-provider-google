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
	clouddeploy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/clouddeploy"
)

func resourceClouddeployDeliveryPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceClouddeployDeliveryPipelineCreate,
		Read:   resourceClouddeployDeliveryPipelineRead,
		Update: resourceClouddeployDeliveryPipelineUpdate,
		Delete: resourceClouddeployDeliveryPipelineDelete,

		Importer: &schema.ResourceImporter{
			State: resourceClouddeployDeliveryPipelineImport,
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the `DeliveryPipeline`. Format is [a-z][a-z0-9\\-]{0,62}.",
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the `DeliveryPipeline`. Max length is 255 characters.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"serial_pipeline": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "SerialPipeline defines a sequential set of stages for a `DeliveryPipeline`.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineSchema(),
			},

			"suspended": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When suspended, no new releases or rollouts can be created, but in-progress ones will complete.",
			},

			"condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Information around the state of the Delivery Pipeline.",
				Elem:        ClouddeployDeliveryPipelineConditionSchema(),
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time at which the pipeline was created.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Unique identifier of the `DeliveryPipeline`.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Most recent time at which the pipeline was updated.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"stages": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Each stage specifies configuration for a `Target`. The ordering of this list defines the promotion flow.",
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"profiles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Skaffold profiles to use when rendering the manifest for this stage's `Target`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"target_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The target_id to which this stage points. This field refers exclusively to the last segment of a target name. For example, this field would just be `my-target` (rather than `projects/project/locations/location/targets/my-target`). The location of the `Target` is inferred to be the same as the location of the `DeliveryPipeline` that contains this `Stage`.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipeline_ready_condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details around the Pipeline's overall status.",
				Elem:        ClouddeployDeliveryPipelineConditionPipelineReadyConditionSchema(),
			},

			"targets_present_condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details around targets enumerated in the pipeline.",
				Elem:        ClouddeployDeliveryPipelineConditionTargetsPresentConditionSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineConditionPipelineReadyConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the Pipeline is in a valid state. Otherwise at least one condition in `PipelineCondition` is in an invalid state. Iterate over those conditions and see which condition(s) has status = false to find out what is wrong with the Pipeline.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last time the condition was updated.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineConditionTargetsPresentConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"missing_targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of Target names that are missing. For example, projects/{project_id}/locations/{location_name}/targets/{target_name}.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if there aren't any missing Targets.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last time the condition was updated.",
			},
		},
	}
}

func resourceClouddeployDeliveryPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Annotations:    checkStringMap(d.Get("annotations")),
		Description:    dcl.String(d.Get("description").(string)),
		Labels:         checkStringMap(d.Get("labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := CreateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyDeliveryPipeline(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating DeliveryPipeline: %s", err)
	}

	log.Printf("[DEBUG] Finished creating DeliveryPipeline %q: %#v", d.Id(), res)

	return resourceClouddeployDeliveryPipelineRead(d, meta)
}

func resourceClouddeployDeliveryPipelineRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Annotations:    checkStringMap(d.Get("annotations")),
		Description:    dcl.String(d.Get("description").(string)),
		Labels:         checkStringMap(d.Get("labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetDeliveryPipeline(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ClouddeployDeliveryPipeline %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("annotations", res.Annotations); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("serial_pipeline", flattenClouddeployDeliveryPipelineSerialPipeline(res.SerialPipeline)); err != nil {
		return fmt.Errorf("error setting serial_pipeline in state: %s", err)
	}
	if err = d.Set("suspended", res.Suspended); err != nil {
		return fmt.Errorf("error setting suspended in state: %s", err)
	}
	if err = d.Set("condition", flattenClouddeployDeliveryPipelineCondition(res.Condition)); err != nil {
		return fmt.Errorf("error setting condition in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceClouddeployDeliveryPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Annotations:    checkStringMap(d.Get("annotations")),
		Description:    dcl.String(d.Get("description").(string)),
		Labels:         checkStringMap(d.Get("labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}
	directive := UpdateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyDeliveryPipeline(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating DeliveryPipeline: %s", err)
	}

	log.Printf("[DEBUG] Finished creating DeliveryPipeline %q: %#v", d.Id(), res)

	return resourceClouddeployDeliveryPipelineRead(d, meta)
}

func resourceClouddeployDeliveryPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Annotations:    checkStringMap(d.Get("annotations")),
		Description:    dcl.String(d.Get("description").(string)),
		Labels:         checkStringMap(d.Get("labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}

	log.Printf("[DEBUG] Deleting DeliveryPipeline %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteDeliveryPipeline(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting DeliveryPipeline: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting DeliveryPipeline %q", d.Id())
	return nil
}

func resourceClouddeployDeliveryPipelineImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/deliveryPipelines/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandClouddeployDeliveryPipelineSerialPipeline(o interface{}) *clouddeploy.DeliveryPipelineSerialPipeline {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipeline
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipeline
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipeline{
		Stages: expandClouddeployDeliveryPipelineSerialPipelineStagesArray(obj["stages"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipeline(obj *clouddeploy.DeliveryPipelineSerialPipeline) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"stages": flattenClouddeployDeliveryPipelineSerialPipelineStagesArray(obj.Stages),
	}

	return []interface{}{transformed}

}
func expandClouddeployDeliveryPipelineSerialPipelineStagesArray(o interface{}) []clouddeploy.DeliveryPipelineSerialPipelineStages {
	if o == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStages, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStages, 0)
	}

	items := make([]clouddeploy.DeliveryPipelineSerialPipelineStages, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployDeliveryPipelineSerialPipelineStages(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployDeliveryPipelineSerialPipelineStages(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStages {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStages
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStages{
		Profiles: expandStringArray(obj["profiles"]),
		TargetId: dcl.String(obj["target_id"].(string)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesArray(objs []clouddeploy.DeliveryPipelineSerialPipelineStages) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployDeliveryPipelineSerialPipelineStages(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployDeliveryPipelineSerialPipelineStages(obj *clouddeploy.DeliveryPipelineSerialPipelineStages) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"profiles":  obj.Profiles,
		"target_id": obj.TargetId,
	}

	return transformed

}

func flattenClouddeployDeliveryPipelineCondition(obj *clouddeploy.DeliveryPipelineCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"pipeline_ready_condition":  flattenClouddeployDeliveryPipelineConditionPipelineReadyCondition(obj.PipelineReadyCondition),
		"targets_present_condition": flattenClouddeployDeliveryPipelineConditionTargetsPresentCondition(obj.TargetsPresentCondition),
	}

	return []interface{}{transformed}

}

func flattenClouddeployDeliveryPipelineConditionPipelineReadyCondition(obj *clouddeploy.DeliveryPipelineConditionPipelineReadyCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"status":      obj.Status,
		"update_time": obj.UpdateTime,
	}

	return []interface{}{transformed}

}

func flattenClouddeployDeliveryPipelineConditionTargetsPresentCondition(obj *clouddeploy.DeliveryPipelineConditionTargetsPresentCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"missing_targets": obj.MissingTargets,
		"status":          obj.Status,
		"update_time":     obj.UpdateTime,
	}

	return []interface{}{transformed}

}
