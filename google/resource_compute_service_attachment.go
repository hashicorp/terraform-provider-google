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
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
)

func resourceComputeServiceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeServiceAttachmentCreate,
		Read:   resourceComputeServiceAttachmentRead,
		Update: resourceComputeServiceAttachmentUpdate,
		Delete: resourceComputeServiceAttachmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeServiceAttachmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"connection_preference": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The connection preference of service attachment. The value can be set to `ACCEPT_AUTOMATIC`. An `ACCEPT_AUTOMATIC` service attachment is one that always accepts the connection from consumer forwarding rules. Possible values: CONNECTION_PREFERENCE_UNSPECIFIED, ACCEPT_AUTOMATIC, ACCEPT_MANUAL",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the name must be 1-63 characters long and match the regular expression `)?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
			},

			"nat_subnets": {
				Type:             schema.TypeList,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceNameList,
				Description:      "An array of URLs where each entry is the URL of a subnet provided by the service producer to use for NAT in this service attachment.",
				Elem:             &schema.Schema{Type: schema.TypeString},
			},

			"target_service": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The URL of a service serving the endpoint identified by this service attachment.",
			},

			"consumer_accept_lists": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Projects that are allowed to connect to this service attachment.",
				Elem:        ComputeServiceAttachmentConsumerAcceptListsSchema(),
			},

			"consumer_reject_lists": {
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceNameList,
				Description:      "Projects that are not allowed to connect to this service attachment. The project can be specified using its id or number.",
				Elem:             &schema.Schema{Type: schema.TypeString},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description of this resource. Provide this property when you create the resource.",
			},

			"enable_proxy_protocol": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "If true, enable the proxy protocol which is for supplying client TCP/IP address data in TCP connections that traverse proxies on their way to destination servers.",
			},

			"project": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The project for the resource",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"connected_endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An array of connections for all the consumers connected to this service attachment.",
				Elem:        ComputeServiceAttachmentConnectedEndpointsSchema(),
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fingerprint of this resource. A hash of the contents stored in this object. This field is used in optimistic locking. This field will be ignored when inserting a `ServiceAttachment`. An up-to-date fingerprint must be provided in order to patch/update the ServiceAttachment; otherwise, the request will fail with error `412 conditionNotMet`. To see the latest fingerprint, make a `get()` request to retrieve the ServiceAttachment.",
			},

			"psc_service_attachment_id": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 128-bit global unique ID of the PSC service attachment.",
				Elem:        ComputeServiceAttachmentPscServiceAttachmentIdSchema(),
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-defined URL for the resource.",
			},

			"service_attachment_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the resource type. The server generates this identifier.",
			},
		},
	}
}

func ComputeServiceAttachmentConsumerAcceptListsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id_or_num": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The project id or number for the project to set the limit for.",
			},

			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The value of the limit to set.",
			},
		},
	}
}

func ComputeServiceAttachmentConnectedEndpointsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The url of a connected endpoint.",
			},

			"psc_connection_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The PSC connection id of the connected endpoint.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of a connected endpoint to this service attachment. Possible values: PENDING, RUNNING, DONE",
			},
		},
	}
}

func ComputeServiceAttachmentPscServiceAttachmentIdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"high": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "",
			},

			"low": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func resourceComputeServiceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ServiceAttachment{
		ConnectionPreference: compute.ServiceAttachmentConnectionPreferenceEnumRef(d.Get("connection_preference").(string)),
		Name:                 dcl.String(d.Get("name").(string)),
		NatSubnets:           expandStringArray(d.Get("nat_subnets")),
		TargetService:        dcl.String(d.Get("target_service").(string)),
		ConsumerAcceptLists:  expandComputeServiceAttachmentConsumerAcceptListsArray(d.Get("consumer_accept_lists")),
		ConsumerRejectLists:  expandStringArray(d.Get("consumer_reject_lists")),
		Description:          dcl.String(d.Get("description").(string)),
		EnableProxyProtocol:  dcl.Bool(d.Get("enable_proxy_protocol").(bool)),
		Project:              dcl.String(project),
		Location:             dcl.String(region),
	}

	id, err := replaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/serviceAttachments/{{name}}")
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	createDirective := CreateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	res, err := client.ApplyServiceAttachment(context.Background(), obj, createDirective...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating ServiceAttachment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating ServiceAttachment %q: %#v", d.Id(), res)

	return resourceComputeServiceAttachmentRead(d, meta)
}

func resourceComputeServiceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ServiceAttachment{
		ConnectionPreference: compute.ServiceAttachmentConnectionPreferenceEnumRef(d.Get("connection_preference").(string)),
		Name:                 dcl.String(d.Get("name").(string)),
		NatSubnets:           expandStringArray(d.Get("nat_subnets")),
		TargetService:        dcl.String(d.Get("target_service").(string)),
		ConsumerAcceptLists:  expandComputeServiceAttachmentConsumerAcceptListsArray(d.Get("consumer_accept_lists")),
		ConsumerRejectLists:  expandStringArray(d.Get("consumer_reject_lists")),
		Description:          dcl.String(d.Get("description").(string)),
		EnableProxyProtocol:  dcl.Bool(d.Get("enable_proxy_protocol").(bool)),
		Project:              dcl.String(project),
		Location:             dcl.String(region),
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
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	res, err := client.GetServiceAttachment(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ComputeServiceAttachment %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("connection_preference", res.ConnectionPreference); err != nil {
		return fmt.Errorf("error setting connection_preference in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("nat_subnets", res.NatSubnets); err != nil {
		return fmt.Errorf("error setting nat_subnets in state: %s", err)
	}
	if err = d.Set("target_service", res.TargetService); err != nil {
		return fmt.Errorf("error setting target_service in state: %s", err)
	}
	if err = d.Set("consumer_accept_lists", flattenComputeServiceAttachmentConsumerAcceptListsArray(res.ConsumerAcceptLists)); err != nil {
		return fmt.Errorf("error setting consumer_accept_lists in state: %s", err)
	}
	if err = d.Set("consumer_reject_lists", res.ConsumerRejectLists); err != nil {
		return fmt.Errorf("error setting consumer_reject_lists in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("enable_proxy_protocol", res.EnableProxyProtocol); err != nil {
		return fmt.Errorf("error setting enable_proxy_protocol in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("region", res.Location); err != nil {
		return fmt.Errorf("error setting region in state: %s", err)
	}
	if err = d.Set("connected_endpoints", flattenComputeServiceAttachmentConnectedEndpointsArray(res.ConnectedEndpoints)); err != nil {
		return fmt.Errorf("error setting connected_endpoints in state: %s", err)
	}
	if err = d.Set("fingerprint", res.Fingerprint); err != nil {
		return fmt.Errorf("error setting fingerprint in state: %s", err)
	}
	if err = d.Set("psc_service_attachment_id", flattenComputeServiceAttachmentPscServiceAttachmentId(res.PscServiceAttachmentId)); err != nil {
		return fmt.Errorf("error setting psc_service_attachment_id in state: %s", err)
	}
	if err = d.Set("self_link", res.SelfLink); err != nil {
		return fmt.Errorf("error setting self_link in state: %s", err)
	}
	if err = d.Set("service_attachment_id", res.Id); err != nil {
		return fmt.Errorf("error setting service_attachment_id in state: %s", err)
	}

	return nil
}
func resourceComputeServiceAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ServiceAttachment{
		ConnectionPreference: compute.ServiceAttachmentConnectionPreferenceEnumRef(d.Get("connection_preference").(string)),
		Name:                 dcl.String(d.Get("name").(string)),
		NatSubnets:           expandStringArray(d.Get("nat_subnets")),
		TargetService:        dcl.String(d.Get("target_service").(string)),
		ConsumerAcceptLists:  expandComputeServiceAttachmentConsumerAcceptListsArray(d.Get("consumer_accept_lists")),
		ConsumerRejectLists:  expandStringArray(d.Get("consumer_reject_lists")),
		Description:          dcl.String(d.Get("description").(string)),
		EnableProxyProtocol:  dcl.Bool(d.Get("enable_proxy_protocol").(bool)),
		Project:              dcl.String(project),
		Location:             dcl.String(region),
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
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	res, err := client.ApplyServiceAttachment(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating ServiceAttachment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating ServiceAttachment %q: %#v", d.Id(), res)

	return resourceComputeServiceAttachmentRead(d, meta)
}

func resourceComputeServiceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	obj := &compute.ServiceAttachment{
		ConnectionPreference: compute.ServiceAttachmentConnectionPreferenceEnumRef(d.Get("connection_preference").(string)),
		Name:                 dcl.String(d.Get("name").(string)),
		NatSubnets:           expandStringArray(d.Get("nat_subnets")),
		TargetService:        dcl.String(d.Get("target_service").(string)),
		ConsumerAcceptLists:  expandComputeServiceAttachmentConsumerAcceptListsArray(d.Get("consumer_accept_lists")),
		ConsumerRejectLists:  expandStringArray(d.Get("consumer_reject_lists")),
		Description:          dcl.String(d.Get("description").(string)),
		EnableProxyProtocol:  dcl.Bool(d.Get("enable_proxy_protocol").(bool)),
		Project:              dcl.String(project),
		Location:             dcl.String(region),
	}

	log.Printf("[DEBUG] Deleting ServiceAttachment %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLComputeClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if err := client.DeleteServiceAttachment(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting ServiceAttachment: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting ServiceAttachment %q", d.Id())
	return nil
}

func resourceComputeServiceAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/serviceAttachments/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/regions/{{region}}/serviceAttachments/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandComputeServiceAttachmentConsumerAcceptListsArray(o interface{}) []compute.ServiceAttachmentConsumerAcceptLists {
	if o == nil {
		return make([]compute.ServiceAttachmentConsumerAcceptLists, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]compute.ServiceAttachmentConsumerAcceptLists, 0)
	}

	items := make([]compute.ServiceAttachmentConsumerAcceptLists, 0, len(objs))
	for _, item := range objs {
		i := expandComputeServiceAttachmentConsumerAcceptLists(item)
		items = append(items, *i)
	}

	return items
}

func expandComputeServiceAttachmentConsumerAcceptLists(o interface{}) *compute.ServiceAttachmentConsumerAcceptLists {
	if o == nil {
		return compute.EmptyServiceAttachmentConsumerAcceptLists
	}

	obj := o.(map[string]interface{})
	return &compute.ServiceAttachmentConsumerAcceptLists{
		ProjectIdOrNum:  dcl.String(obj["project_id_or_num"].(string)),
		ConnectionLimit: dcl.Int64(int64(obj["connection_limit"].(int))),
	}
}

func flattenComputeServiceAttachmentConsumerAcceptListsArray(objs []compute.ServiceAttachmentConsumerAcceptLists) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeServiceAttachmentConsumerAcceptLists(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeServiceAttachmentConsumerAcceptLists(obj *compute.ServiceAttachmentConsumerAcceptLists) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"project_id_or_num": obj.ProjectIdOrNum,
		"connection_limit":  obj.ConnectionLimit,
	}

	return transformed

}

func flattenComputeServiceAttachmentConnectedEndpointsArray(objs []compute.ServiceAttachmentConnectedEndpoints) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenComputeServiceAttachmentConnectedEndpoints(&item)
		items = append(items, i)
	}

	return items
}

func flattenComputeServiceAttachmentConnectedEndpoints(obj *compute.ServiceAttachmentConnectedEndpoints) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"endpoint":          obj.Endpoint,
		"psc_connection_id": obj.PscConnectionId,
		"status":            obj.Status,
	}

	return transformed

}

func flattenComputeServiceAttachmentPscServiceAttachmentId(obj *compute.ServiceAttachmentPscServiceAttachmentId) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"high": obj.High,
		"low":  obj.Low,
	}

	return []interface{}{transformed}

}
