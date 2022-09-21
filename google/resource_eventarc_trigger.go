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
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
)

func resourceEventarcTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventarcTriggerCreate,
		Read:   resourceEventarcTriggerRead,
		Update: resourceEventarcTriggerUpdate,
		Delete: resourceEventarcTriggerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceEventarcTriggerImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"destination": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Destination specifies where the events should be sent to.",
				MaxItems:    1,
				Elem:        EventarcTriggerDestinationSchema(),
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"matching_criteria": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Required. null The list of filters that applies to event attributes. Only events that match all the provided filters will be sent to the destination.",
				Elem:        EventarcTriggerMatchingCriteriaSchema(),
				Set:         schema.HashResource(EventarcTriggerMatchingCriteriaSchema()),
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The resource name of the trigger. Must be unique within the location on the project.",
			},

			"channel": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Optional. The name of the channel associated with the trigger in `projects/{project}/locations/{location}/channels/{channel}` format. You must provide a channel to receive events from Eventarc SaaS partners.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. User labels attached to the triggers that can be used to group resources.",
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

			"service_account": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Optional. The IAM service account email associated with the trigger. The service account represents the identity of the trigger. The principal who calls this API must have `iam.serviceAccounts.actAs` permission in the service account. See https://cloud.google.com/iam/docs/understanding-service-accounts#sa_common for more information. For Cloud Run destinations, this service account is used to generate identity tokens when invoking the service. See https://cloud.google.com/run/docs/triggering/pubsub-push#create-service-account for information on how to invoke authenticated Cloud Run services. In order to create Audit Log triggers, the service account should also have `roles/eventarc.eventReceiver` IAM role.",
			},

			"transport": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional. In order to deliver messages, Eventarc may use other GCP products as transport intermediary. This field contains a reference to that transport intermediary. This information can be used for debugging purposes.",
				MaxItems:    1,
				Elem:        EventarcTriggerTransportSchema(),
			},

			"conditions": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Output only. The reason(s) why a trigger is in FAILED state.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The creation time.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. This checksum is computed by the server based on the value of other fields, and may be sent only on create requests to ensure the client has an up-to-date value before proceeding.",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Server assigned unique identifier for the trigger. The value is a UUID4 string and guaranteed to remain unchanged until the resource is deleted.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The last-modified time.",
			},
		},
	}
}

func EventarcTriggerDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_function": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "[WARNING] Configuring a Cloud Function in Trigger is not supported as of today. The Cloud Function resource name. Format: projects/{project}/locations/{location}/functions/{function}",
			},

			"cloud_run_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cloud Run fully-managed service that receives the events. The service should be running in the same project of the trigger.",
				MaxItems:    1,
				Elem:        EventarcTriggerDestinationCloudRunServiceSchema(),
			},

			"gke": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A GKE service capable of receiving events. The service should be running in the same project as the trigger.",
				MaxItems:    1,
				Elem:        EventarcTriggerDestinationGkeSchema(),
			},

			"workflow": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The resource name of the Workflow whose Executions are triggered by the events. The Workflow resource should be deployed in the same project as the trigger. Format: `projects/{project}/locations/{location}/workflows/{workflow}`",
			},
		},
	}
}

func EventarcTriggerDestinationCloudRunServiceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Required. The name of the Cloud Run service being addressed. See https://cloud.google.com/run/docs/reference/rest/v1/namespaces.services. Only services located in the same project of the trigger object can be addressed.",
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The relative path on the Cloud Run service the events should be sent to. The value must conform to the definition of URI path segment (section 3.3 of RFC2396). Examples: \"/route\", \"route\", \"route/subroute\".",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Required. The region the Cloud Run service is deployed in.",
			},
		},
	}
}

func EventarcTriggerDestinationGkeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Required. The name of the cluster the GKE service is running in. The cluster must be running in the same project as the trigger being created.",
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The name of the Google Compute Engine in which the cluster resides, which can either be compute zone (for example, us-central1-a) for the zonal clusters or region (for example, us-central1) for regional clusters.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The namespace the GKE service is running in.",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the GKE service.",
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The relative path on the GKE service the events should be sent to. The value must conform to the definition of a URI path segment (section 3.3 of RFC2396). Examples: \"/route\", \"route\", \"route/subroute\".",
			},
		},
	}
}

func EventarcTriggerMatchingCriteriaSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The name of a CloudEvents attribute. Currently, only a subset of attributes are supported for filtering. All triggers MUST provide a filter for the 'type' attribute.",
			},

			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The value for the attribute. See https://cloud.google.com/eventarc/docs/creating-triggers#trigger-gcloud for available values.",
			},

			"operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The operator used for matching the events with the value of the filter. If not specified, only events that have an exact key-value pair specified in the filter are matched. The only allowed value is `match-path-pattern`.",
			},
		},
	}
}

func EventarcTriggerTransportSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pubsub": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The Pub/Sub topic and subscription used by Eventarc as delivery intermediary.",
				MaxItems:    1,
				Elem:        EventarcTriggerTransportPubsubSchema(),
			},
		},
	}
}

func EventarcTriggerTransportPubsubSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "Optional. The name of the Pub/Sub topic created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/topics/{TOPIC_NAME}. You may set an existing topic for triggers of the type google.cloud.pubsub.topic.v1.messagePublished` only. The topic you provide here will not be deleted by Eventarc at trigger deletion.",
			},

			"subscription": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The name of the Pub/Sub subscription created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/subscriptions/{SUBSCRIPTION_NAME}`.",
			},
		},
	}
}

func resourceEventarcTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &eventarc.Trigger{
		Destination:      expandEventarcTriggerDestination(d.Get("destination")),
		Location:         dcl.String(d.Get("location").(string)),
		MatchingCriteria: expandEventarcTriggerMatchingCriteriaArray(d.Get("matching_criteria")),
		Name:             dcl.String(d.Get("name").(string)),
		Channel:          dcl.String(d.Get("channel").(string)),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		ServiceAccount:   dcl.String(d.Get("service_account").(string)),
		Transport:        expandEventarcTriggerTransport(d.Get("transport")),
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
	client := NewDCLEventarcClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyTrigger(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Trigger: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Trigger %q: %#v", d.Id(), res)

	return resourceEventarcTriggerRead(d, meta)
}

func resourceEventarcTriggerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &eventarc.Trigger{
		Destination:      expandEventarcTriggerDestination(d.Get("destination")),
		Location:         dcl.String(d.Get("location").(string)),
		MatchingCriteria: expandEventarcTriggerMatchingCriteriaArray(d.Get("matching_criteria")),
		Name:             dcl.String(d.Get("name").(string)),
		Channel:          dcl.String(d.Get("channel").(string)),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		ServiceAccount:   dcl.String(d.Get("service_account").(string)),
		Transport:        expandEventarcTriggerTransport(d.Get("transport")),
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
	client := NewDCLEventarcClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetTrigger(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("EventarcTrigger %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("destination", flattenEventarcTriggerDestination(res.Destination)); err != nil {
		return fmt.Errorf("error setting destination in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("matching_criteria", flattenEventarcTriggerMatchingCriteriaArray(res.MatchingCriteria)); err != nil {
		return fmt.Errorf("error setting matching_criteria in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("channel", res.Channel); err != nil {
		return fmt.Errorf("error setting channel in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("service_account", res.ServiceAccount); err != nil {
		return fmt.Errorf("error setting service_account in state: %s", err)
	}
	if err = d.Set("transport", flattenEventarcTriggerTransport(res.Transport)); err != nil {
		return fmt.Errorf("error setting transport in state: %s", err)
	}
	if err = d.Set("conditions", res.Conditions); err != nil {
		return fmt.Errorf("error setting conditions in state: %s", err)
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
func resourceEventarcTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &eventarc.Trigger{
		Destination:      expandEventarcTriggerDestination(d.Get("destination")),
		Location:         dcl.String(d.Get("location").(string)),
		MatchingCriteria: expandEventarcTriggerMatchingCriteriaArray(d.Get("matching_criteria")),
		Name:             dcl.String(d.Get("name").(string)),
		Channel:          dcl.String(d.Get("channel").(string)),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		ServiceAccount:   dcl.String(d.Get("service_account").(string)),
		Transport:        expandEventarcTriggerTransport(d.Get("transport")),
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
	client := NewDCLEventarcClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyTrigger(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Trigger: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Trigger %q: %#v", d.Id(), res)

	return resourceEventarcTriggerRead(d, meta)
}

func resourceEventarcTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &eventarc.Trigger{
		Destination:      expandEventarcTriggerDestination(d.Get("destination")),
		Location:         dcl.String(d.Get("location").(string)),
		MatchingCriteria: expandEventarcTriggerMatchingCriteriaArray(d.Get("matching_criteria")),
		Name:             dcl.String(d.Get("name").(string)),
		Channel:          dcl.String(d.Get("channel").(string)),
		Labels:           checkStringMap(d.Get("labels")),
		Project:          dcl.String(project),
		ServiceAccount:   dcl.String(d.Get("service_account").(string)),
		Transport:        expandEventarcTriggerTransport(d.Get("transport")),
	}

	log.Printf("[DEBUG] Deleting Trigger %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLEventarcClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteTrigger(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Trigger: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Trigger %q", d.Id())
	return nil
}

func resourceEventarcTriggerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/triggers/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/triggers/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandEventarcTriggerDestination(o interface{}) *eventarc.TriggerDestination {
	if o == nil {
		return eventarc.EmptyTriggerDestination
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return eventarc.EmptyTriggerDestination
	}
	obj := objArr[0].(map[string]interface{})
	return &eventarc.TriggerDestination{
		CloudFunction:   dcl.String(obj["cloud_function"].(string)),
		CloudRunService: expandEventarcTriggerDestinationCloudRunService(obj["cloud_run_service"]),
		Gke:             expandEventarcTriggerDestinationGke(obj["gke"]),
		Workflow:        dcl.String(obj["workflow"].(string)),
	}
}

func flattenEventarcTriggerDestination(obj *eventarc.TriggerDestination) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cloud_function":    obj.CloudFunction,
		"cloud_run_service": flattenEventarcTriggerDestinationCloudRunService(obj.CloudRunService),
		"gke":               flattenEventarcTriggerDestinationGke(obj.Gke),
		"workflow":          obj.Workflow,
	}

	return []interface{}{transformed}

}

func expandEventarcTriggerDestinationCloudRunService(o interface{}) *eventarc.TriggerDestinationCloudRunService {
	if o == nil {
		return eventarc.EmptyTriggerDestinationCloudRunService
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return eventarc.EmptyTriggerDestinationCloudRunService
	}
	obj := objArr[0].(map[string]interface{})
	return &eventarc.TriggerDestinationCloudRunService{
		Service: dcl.String(obj["service"].(string)),
		Path:    dcl.String(obj["path"].(string)),
		Region:  dcl.StringOrNil(obj["region"].(string)),
	}
}

func flattenEventarcTriggerDestinationCloudRunService(obj *eventarc.TriggerDestinationCloudRunService) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"service": obj.Service,
		"path":    obj.Path,
		"region":  obj.Region,
	}

	return []interface{}{transformed}

}

func expandEventarcTriggerDestinationGke(o interface{}) *eventarc.TriggerDestinationGke {
	if o == nil {
		return eventarc.EmptyTriggerDestinationGke
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return eventarc.EmptyTriggerDestinationGke
	}
	obj := objArr[0].(map[string]interface{})
	return &eventarc.TriggerDestinationGke{
		Cluster:   dcl.String(obj["cluster"].(string)),
		Location:  dcl.String(obj["location"].(string)),
		Namespace: dcl.String(obj["namespace"].(string)),
		Service:   dcl.String(obj["service"].(string)),
		Path:      dcl.String(obj["path"].(string)),
	}
}

func flattenEventarcTriggerDestinationGke(obj *eventarc.TriggerDestinationGke) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster":   obj.Cluster,
		"location":  obj.Location,
		"namespace": obj.Namespace,
		"service":   obj.Service,
		"path":      obj.Path,
	}

	return []interface{}{transformed}

}
func expandEventarcTriggerMatchingCriteriaArray(o interface{}) []eventarc.TriggerMatchingCriteria {
	if o == nil {
		return make([]eventarc.TriggerMatchingCriteria, 0)
	}

	o = o.(*schema.Set).List()

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]eventarc.TriggerMatchingCriteria, 0)
	}

	items := make([]eventarc.TriggerMatchingCriteria, 0, len(objs))
	for _, item := range objs {
		i := expandEventarcTriggerMatchingCriteria(item)
		items = append(items, *i)
	}

	return items
}

func expandEventarcTriggerMatchingCriteria(o interface{}) *eventarc.TriggerMatchingCriteria {
	if o == nil {
		return eventarc.EmptyTriggerMatchingCriteria
	}

	obj := o.(map[string]interface{})
	return &eventarc.TriggerMatchingCriteria{
		Attribute: dcl.String(obj["attribute"].(string)),
		Value:     dcl.String(obj["value"].(string)),
		Operator:  dcl.String(obj["operator"].(string)),
	}
}

func flattenEventarcTriggerMatchingCriteriaArray(objs []eventarc.TriggerMatchingCriteria) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenEventarcTriggerMatchingCriteria(&item)
		items = append(items, i)
	}

	return items
}

func flattenEventarcTriggerMatchingCriteria(obj *eventarc.TriggerMatchingCriteria) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"attribute": obj.Attribute,
		"value":     obj.Value,
		"operator":  obj.Operator,
	}

	return transformed

}

func expandEventarcTriggerTransport(o interface{}) *eventarc.TriggerTransport {
	if o == nil {
		return nil
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return nil
	}
	obj := objArr[0].(map[string]interface{})
	return &eventarc.TriggerTransport{
		Pubsub: expandEventarcTriggerTransportPubsub(obj["pubsub"]),
	}
}

func flattenEventarcTriggerTransport(obj *eventarc.TriggerTransport) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"pubsub": flattenEventarcTriggerTransportPubsub(obj.Pubsub),
	}

	return []interface{}{transformed}

}

func expandEventarcTriggerTransportPubsub(o interface{}) *eventarc.TriggerTransportPubsub {
	if o == nil {
		return eventarc.EmptyTriggerTransportPubsub
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return eventarc.EmptyTriggerTransportPubsub
	}
	obj := objArr[0].(map[string]interface{})
	return &eventarc.TriggerTransportPubsub{
		Topic: dcl.String(obj["topic"].(string)),
	}
}

func flattenEventarcTriggerTransportPubsub(obj *eventarc.TriggerTransportPubsub) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"topic":        obj.Topic,
		"subscription": obj.Subscription,
	}

	return []interface{}{transformed}

}
