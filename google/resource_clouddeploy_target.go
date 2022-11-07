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

func resourceClouddeployTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceClouddeployTargetCreate,
		Read:   resourceClouddeployTargetRead,
		Update: resourceClouddeployTargetUpdate,
		Delete: resourceClouddeployTargetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceClouddeployTargetImport,
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
				Description: "Name of the `Target`. Format is [a-z][a-z0-9\\-]{0,62}.",
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"anthos_cluster": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Information specifying an Anthos Cluster.",
				MaxItems:      1,
				Elem:          ClouddeployTargetAnthosClusterSchema(),
				ConflictsWith: []string{"gke"},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Description of the `Target`. Max length is 255 characters.",
			},

			"execution_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Configurations for all execution that relates to this `Target`. Each `ExecutionEnvironmentUsage` value may only be used in a single configuration; using the same value multiple times is an error. When one or more configurations are specified, they must include the `RENDER` and `DEPLOY` `ExecutionEnvironmentUsage` values. When no configurations are specified, execution will use the default specified in `DefaultPool`.",
				Elem:        ClouddeployTargetExecutionConfigsSchema(),
			},

			"gke": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Information specifying a GKE Cluster.",
				MaxItems:      1,
				Elem:          ClouddeployTargetGkeSchema(),
				ConflictsWith: []string{"anthos_cluster"},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.",
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

			"require_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Whether or not the `Target` requires approval.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time at which the `Target` was created.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional. This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Resource id of the `Target`.",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Unique identifier of the `Target`.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Most recent time at which the `Target` was updated.",
			},
		},
	}
}

func ClouddeployTargetAnthosClusterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"membership": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Membership of the GKE Hub-registered cluster to which to apply the Skaffold configuration. Format is `projects/{project}/locations/{location}/memberships/{membership_name}`.",
			},
		},
	}
}

func ClouddeployTargetExecutionConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"usages": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Usages when this configuration should be applied.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"artifact_storage": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Cloud Storage location in which to store execution outputs. This can either be a bucket (\"gs://my-bucket\") or a path within a bucket (\"gs://my-bucket/my-dir\"). If unspecified, a default bucket located in the same region will be used.",
			},

			"execution_timeout": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Execution timeout for a Cloud Build Execution. This must be between 10m and 24h in seconds format. If unspecified, a default timeout of 1h is used.",
			},

			"service_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Google service account to use for execution. If unspecified, the project execution service account (-compute@developer.gserviceaccount.com) is used.",
			},

			"worker_pool": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Optional. The resource name of the `WorkerPool`, with the format `projects/{project}/locations/{location}/workerPools/{worker_pool}`. If this optional field is unspecified, the default Cloud Build pool will be used.",
			},
		},
	}
}

func ClouddeployTargetGkeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Information specifying a GKE Cluster. Format is `projects/{project_id}/locations/{location_id}/clusters/{cluster_id}.",
			},

			"internal_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. If true, `cluster` is accessed using the private IP address of the control plane endpoint. Otherwise, the default IP address of the control plane endpoint is used. The default IP address is the private IP address for clusters with private control-plane endpoints and the public IP address otherwise. Only specify this option when `cluster` is a [private GKE cluster](https://cloud.google.com/kubernetes-engine/docs/concepts/private-cluster-concept).",
			},
		},
	}
}

func resourceClouddeployTargetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		Annotations:      checkStringMap(d.Get("annotations")),
		AnthosCluster:    expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		Description:      dcl.String(d.Get("description").(string)),
		ExecutionConfigs: expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:              expandClouddeployTargetGke(d.Get("gke")),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		RequireApproval:  dcl.Bool(d.Get("require_approval").(bool)),
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
	res, err := client.ApplyTarget(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Target: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Target %q: %#v", d.Id(), res)

	return resourceClouddeployTargetRead(d, meta)
}

func resourceClouddeployTargetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		Annotations:      checkStringMap(d.Get("annotations")),
		AnthosCluster:    expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		Description:      dcl.String(d.Get("description").(string)),
		ExecutionConfigs: expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:              expandClouddeployTargetGke(d.Get("gke")),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		RequireApproval:  dcl.Bool(d.Get("require_approval").(bool)),
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
	res, err := client.GetTarget(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ClouddeployTarget %q", d.Id())
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
	if err = d.Set("anthos_cluster", flattenClouddeployTargetAnthosCluster(res.AnthosCluster)); err != nil {
		return fmt.Errorf("error setting anthos_cluster in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("execution_configs", flattenClouddeployTargetExecutionConfigsArray(res.ExecutionConfigs)); err != nil {
		return fmt.Errorf("error setting execution_configs in state: %s", err)
	}
	if err = d.Set("gke", flattenClouddeployTargetGke(res.Gke)); err != nil {
		return fmt.Errorf("error setting gke in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("require_approval", res.RequireApproval); err != nil {
		return fmt.Errorf("error setting require_approval in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("target_id", res.TargetId); err != nil {
		return fmt.Errorf("error setting target_id in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceClouddeployTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		Annotations:      checkStringMap(d.Get("annotations")),
		AnthosCluster:    expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		Description:      dcl.String(d.Get("description").(string)),
		ExecutionConfigs: expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:              expandClouddeployTargetGke(d.Get("gke")),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		RequireApproval:  dcl.Bool(d.Get("require_approval").(bool)),
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
	res, err := client.ApplyTarget(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Target: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Target %q: %#v", d.Id(), res)

	return resourceClouddeployTargetRead(d, meta)
}

func resourceClouddeployTargetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		Annotations:      checkStringMap(d.Get("annotations")),
		AnthosCluster:    expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		Description:      dcl.String(d.Get("description").(string)),
		ExecutionConfigs: expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:              expandClouddeployTargetGke(d.Get("gke")),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		RequireApproval:  dcl.Bool(d.Get("require_approval").(bool)),
	}

	log.Printf("[DEBUG] Deleting Target %q", d.Id())
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
	if err := client.DeleteTarget(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Target: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Target %q", d.Id())
	return nil
}

func resourceClouddeployTargetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/targets/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/targets/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandClouddeployTargetAnthosCluster(o interface{}) *clouddeploy.TargetAnthosCluster {
	if o == nil {
		return clouddeploy.EmptyTargetAnthosCluster
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyTargetAnthosCluster
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.TargetAnthosCluster{
		Membership: dcl.String(obj["membership"].(string)),
	}
}

func flattenClouddeployTargetAnthosCluster(obj *clouddeploy.TargetAnthosCluster) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"membership": obj.Membership,
	}

	return []interface{}{transformed}

}
func expandClouddeployTargetExecutionConfigsArray(o interface{}) []clouddeploy.TargetExecutionConfigs {
	if o == nil {
		return nil
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return nil
	}

	items := make([]clouddeploy.TargetExecutionConfigs, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployTargetExecutionConfigs(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployTargetExecutionConfigs(o interface{}) *clouddeploy.TargetExecutionConfigs {
	if o == nil {
		return nil
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.TargetExecutionConfigs{
		Usages:           expandClouddeployTargetExecutionConfigsUsagesArray(obj["usages"]),
		ArtifactStorage:  dcl.StringOrNil(obj["artifact_storage"].(string)),
		ExecutionTimeout: dcl.StringOrNil(obj["execution_timeout"].(string)),
		ServiceAccount:   dcl.StringOrNil(obj["service_account"].(string)),
		WorkerPool:       dcl.String(obj["worker_pool"].(string)),
	}
}

func flattenClouddeployTargetExecutionConfigsArray(objs []clouddeploy.TargetExecutionConfigs) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployTargetExecutionConfigs(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployTargetExecutionConfigs(obj *clouddeploy.TargetExecutionConfigs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"usages":            flattenClouddeployTargetExecutionConfigsUsagesArray(obj.Usages),
		"artifact_storage":  obj.ArtifactStorage,
		"execution_timeout": obj.ExecutionTimeout,
		"service_account":   obj.ServiceAccount,
		"worker_pool":       obj.WorkerPool,
	}

	return transformed

}

func expandClouddeployTargetGke(o interface{}) *clouddeploy.TargetGke {
	if o == nil {
		return clouddeploy.EmptyTargetGke
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyTargetGke
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.TargetGke{
		Cluster:    dcl.String(obj["cluster"].(string)),
		InternalIP: dcl.Bool(obj["internal_ip"].(bool)),
	}
}

func flattenClouddeployTargetGke(obj *clouddeploy.TargetGke) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster":     obj.Cluster,
		"internal_ip": obj.InternalIP,
	}

	return []interface{}{transformed}

}
func flattenClouddeployTargetExecutionConfigsUsagesArray(obj []clouddeploy.TargetExecutionConfigsUsagesEnum) interface{} {
	if obj == nil {
		return nil
	}
	items := []string{}
	for _, item := range obj {
		items = append(items, string(item))
	}
	return items
}
func expandClouddeployTargetExecutionConfigsUsagesArray(o interface{}) []clouddeploy.TargetExecutionConfigsUsagesEnum {
	objs := o.([]interface{})
	items := make([]clouddeploy.TargetExecutionConfigsUsagesEnum, 0, len(objs))
	for _, item := range objs {
		i := clouddeploy.TargetExecutionConfigsUsagesEnumRef(item.(string))
		items = append(items, *i)
	}
	return items
}
