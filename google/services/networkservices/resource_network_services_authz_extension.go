// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
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

package networkservices

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
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func ResourceNetworkServicesAuthzExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkServicesAuthzExtensionCreate,
		Read:   resourceNetworkServicesAuthzExtensionRead,
		Update: resourceNetworkServicesAuthzExtensionUpdate,
		Delete: resourceNetworkServicesAuthzExtensionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceNetworkServicesAuthzExtensionImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.SetLabelsDiff,
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"authority": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The :authority header in the gRPC request sent from Envoy to the extension service.`,
			},
			"load_balancing_scheme": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidateEnum([]string{"INTERNAL_MANAGED", "EXTERNAL_MANAGED"}),
				Description: `All backend services and forwarding rules referenced by this extension must share the same load balancing scheme.
For more information, refer to [Backend services overview](https://cloud.google.com/load-balancing/docs/backend-service). Possible values: ["INTERNAL_MANAGED", "EXTERNAL_MANAGED"]`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The location of the resource.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Identifier. Name of the AuthzExtension resource.`,
			},
			"service": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.ProjectNumberDiffSuppress,
				Description: `The reference to the service that runs the extension.
To configure a callout extension, service must be a fully-qualified reference to a [backend service](https://cloud.google.com/compute/docs/reference/rest/v1/backendServices) in the format:
https://www.googleapis.com/compute/v1/projects/{project}/regions/{region}/backendServices/{backendService} or https://www.googleapis.com/compute/v1/projects/{project}/global/backendServices/{backendService}.`,
			},
			"timeout": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.DurationDiffSuppress,
				Description:      `Specifies the timeout for each individual message on the stream. The timeout must be between 10-10000 milliseconds.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A human-readable description of the resource.`,
			},
			"fail_open": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				Description: `Determines how the proxy behaves if the call to the extension fails or times out.
When set to TRUE, request or response processing continues without error. Any subsequent extensions in the extension chain are also executed. When set to FALSE or the default setting of FALSE is used, one of the following happens:
* If response headers have not been delivered to the downstream client, a generic 500 error is returned to the client. The error response can be tailored by configuring a custom error response in the load balancer.
* If response headers have been delivered, then the HTTP stream to the downstream client is reset.`,
			},
			"forward_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `List of the HTTP headers to forward to the extension (from the client). If omitted, all headers are sent. Each element is a string indicating the header name.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Set of labels associated with the AuthzExtension resource.


**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field 'effective_labels' for all of the labels present on the resource.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `The metadata provided here is included as part of the metadata_context (of type google.protobuf.Struct) in the ProcessingRequest message sent to the extension server. The metadata is available under the namespace com.google.authz_extension.<resourceName>. The following variables are supported in the metadata Struct:

{forwarding_rule_id} - substituted with the forwarding rule's fully qualified resource name.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"wire_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidateEnum([]string{"WIRE_FORMAT_UNSPECIFIED", "EXT_PROC_GRPC", ""}),
				Description:  `The format of communication supported by the callout extension. Default value: "EXT_PROC_GRPC" Possible values: ["WIRE_FORMAT_UNSPECIFIED", "EXT_PROC_GRPC"]`,
				Default:      "EXT_PROC_GRPC",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The timestamp when the resource was created.`,
			},
			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"terraform_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Description: `The combination of labels configured directly on the resource
 and default labels configured on the provider.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The timestamp when the resource was updated.`,
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

func resourceNetworkServicesAuthzExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkServicesAuthzExtensionDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	loadBalancingSchemeProp, err := expandNetworkServicesAuthzExtensionLoadBalancingScheme(d.Get("load_balancing_scheme"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("load_balancing_scheme"); !tpgresource.IsEmptyValue(reflect.ValueOf(loadBalancingSchemeProp)) && (ok || !reflect.DeepEqual(v, loadBalancingSchemeProp)) {
		obj["loadBalancingScheme"] = loadBalancingSchemeProp
	}
	authorityProp, err := expandNetworkServicesAuthzExtensionAuthority(d.Get("authority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("authority"); !tpgresource.IsEmptyValue(reflect.ValueOf(authorityProp)) && (ok || !reflect.DeepEqual(v, authorityProp)) {
		obj["authority"] = authorityProp
	}
	serviceProp, err := expandNetworkServicesAuthzExtensionService(d.Get("service"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("service"); !tpgresource.IsEmptyValue(reflect.ValueOf(serviceProp)) && (ok || !reflect.DeepEqual(v, serviceProp)) {
		obj["service"] = serviceProp
	}
	timeoutProp, err := expandNetworkServicesAuthzExtensionTimeout(d.Get("timeout"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("timeout"); !tpgresource.IsEmptyValue(reflect.ValueOf(timeoutProp)) && (ok || !reflect.DeepEqual(v, timeoutProp)) {
		obj["timeout"] = timeoutProp
	}
	failOpenProp, err := expandNetworkServicesAuthzExtensionFailOpen(d.Get("fail_open"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("fail_open"); ok || !reflect.DeepEqual(v, failOpenProp) {
		obj["failOpen"] = failOpenProp
	}
	metadataProp, err := expandNetworkServicesAuthzExtensionMetadata(d.Get("metadata"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("metadata"); !tpgresource.IsEmptyValue(reflect.ValueOf(metadataProp)) && (ok || !reflect.DeepEqual(v, metadataProp)) {
		obj["metadata"] = metadataProp
	}
	forwardHeadersProp, err := expandNetworkServicesAuthzExtensionForwardHeaders(d.Get("forward_headers"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("forward_headers"); !tpgresource.IsEmptyValue(reflect.ValueOf(forwardHeadersProp)) && (ok || !reflect.DeepEqual(v, forwardHeadersProp)) {
		obj["forwardHeaders"] = forwardHeadersProp
	}
	wireFormatProp, err := expandNetworkServicesAuthzExtensionWireFormat(d.Get("wire_format"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("wire_format"); !tpgresource.IsEmptyValue(reflect.ValueOf(wireFormatProp)) && (ok || !reflect.DeepEqual(v, wireFormatProp)) {
		obj["wireFormat"] = wireFormatProp
	}
	labelsProp, err := expandNetworkServicesAuthzExtensionEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	nameProp, err := expandNetworkServicesAuthzExtensionName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/authzExtensions?authzExtensionId={{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new AuthzExtension: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AuthzExtension: %s", err)
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
		return fmt.Errorf("Error creating AuthzExtension: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/authzExtensions/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = NetworkServicesOperationWaitTime(
		config, res, project, "Creating AuthzExtension", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create AuthzExtension: %s", err)
	}

	log.Printf("[DEBUG] Finished creating AuthzExtension %q: %#v", d.Id(), res)

	return resourceNetworkServicesAuthzExtensionRead(d, meta)
}

func resourceNetworkServicesAuthzExtensionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/authzExtensions/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AuthzExtension: %s", err)
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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("NetworkServicesAuthzExtension %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}

	if err := d.Set("create_time", flattenNetworkServicesAuthzExtensionCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("update_time", flattenNetworkServicesAuthzExtensionUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("description", flattenNetworkServicesAuthzExtensionDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("labels", flattenNetworkServicesAuthzExtensionLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("load_balancing_scheme", flattenNetworkServicesAuthzExtensionLoadBalancingScheme(res["loadBalancingScheme"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("authority", flattenNetworkServicesAuthzExtensionAuthority(res["authority"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("service", flattenNetworkServicesAuthzExtensionService(res["service"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("timeout", flattenNetworkServicesAuthzExtensionTimeout(res["timeout"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("fail_open", flattenNetworkServicesAuthzExtensionFailOpen(res["failOpen"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("metadata", flattenNetworkServicesAuthzExtensionMetadata(res["metadata"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("forward_headers", flattenNetworkServicesAuthzExtensionForwardHeaders(res["forwardHeaders"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("wire_format", flattenNetworkServicesAuthzExtensionWireFormat(res["wireFormat"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("terraform_labels", flattenNetworkServicesAuthzExtensionTerraformLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("effective_labels", flattenNetworkServicesAuthzExtensionEffectiveLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}
	if err := d.Set("name", flattenNetworkServicesAuthzExtensionName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading AuthzExtension: %s", err)
	}

	return nil
}

func resourceNetworkServicesAuthzExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AuthzExtension: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	descriptionProp, err := expandNetworkServicesAuthzExtensionDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	loadBalancingSchemeProp, err := expandNetworkServicesAuthzExtensionLoadBalancingScheme(d.Get("load_balancing_scheme"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("load_balancing_scheme"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, loadBalancingSchemeProp)) {
		obj["loadBalancingScheme"] = loadBalancingSchemeProp
	}
	authorityProp, err := expandNetworkServicesAuthzExtensionAuthority(d.Get("authority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("authority"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, authorityProp)) {
		obj["authority"] = authorityProp
	}
	serviceProp, err := expandNetworkServicesAuthzExtensionService(d.Get("service"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("service"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, serviceProp)) {
		obj["service"] = serviceProp
	}
	timeoutProp, err := expandNetworkServicesAuthzExtensionTimeout(d.Get("timeout"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("timeout"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, timeoutProp)) {
		obj["timeout"] = timeoutProp
	}
	failOpenProp, err := expandNetworkServicesAuthzExtensionFailOpen(d.Get("fail_open"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("fail_open"); ok || !reflect.DeepEqual(v, failOpenProp) {
		obj["failOpen"] = failOpenProp
	}
	metadataProp, err := expandNetworkServicesAuthzExtensionMetadata(d.Get("metadata"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("metadata"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, metadataProp)) {
		obj["metadata"] = metadataProp
	}
	forwardHeadersProp, err := expandNetworkServicesAuthzExtensionForwardHeaders(d.Get("forward_headers"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("forward_headers"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, forwardHeadersProp)) {
		obj["forwardHeaders"] = forwardHeadersProp
	}
	wireFormatProp, err := expandNetworkServicesAuthzExtensionWireFormat(d.Get("wire_format"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("wire_format"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, wireFormatProp)) {
		obj["wireFormat"] = wireFormatProp
	}
	labelsProp, err := expandNetworkServicesAuthzExtensionEffectiveLabels(d.Get("effective_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("effective_labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	nameProp, err := expandNetworkServicesAuthzExtensionName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/authzExtensions/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating AuthzExtension %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}

	if d.HasChange("load_balancing_scheme") {
		updateMask = append(updateMask, "loadBalancingScheme")
	}

	if d.HasChange("authority") {
		updateMask = append(updateMask, "authority")
	}

	if d.HasChange("service") {
		updateMask = append(updateMask, "service")
	}

	if d.HasChange("timeout") {
		updateMask = append(updateMask, "timeout")
	}

	if d.HasChange("fail_open") {
		updateMask = append(updateMask, "failOpen")
	}

	if d.HasChange("metadata") {
		updateMask = append(updateMask, "metadata")
	}

	if d.HasChange("forward_headers") {
		updateMask = append(updateMask, "forwardHeaders")
	}

	if d.HasChange("wire_format") {
		updateMask = append(updateMask, "wireFormat")
	}

	if d.HasChange("effective_labels") {
		updateMask = append(updateMask, "labels")
	}

	if d.HasChange("name") {
		updateMask = append(updateMask, "name")
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
			return fmt.Errorf("Error updating AuthzExtension %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating AuthzExtension %q: %#v", d.Id(), res)
		}

		err = NetworkServicesOperationWaitTime(
			config, res, project, "Updating AuthzExtension", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceNetworkServicesAuthzExtensionRead(d, meta)
}

func resourceNetworkServicesAuthzExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for AuthzExtension: %s", err)
	}
	billingProject = project

	url, err := tpgresource.ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/authzExtensions/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting AuthzExtension %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "AuthzExtension")
	}

	err = NetworkServicesOperationWaitTime(
		config, res, project, "Deleting AuthzExtension", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting AuthzExtension %q: %#v", d.Id(), res)
	return nil
}

func resourceNetworkServicesAuthzExtensionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/authzExtensions/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<location>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/authzExtensions/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNetworkServicesAuthzExtensionCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenNetworkServicesAuthzExtensionLoadBalancingScheme(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionAuthority(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionService(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenNetworkServicesAuthzExtensionTimeout(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionFailOpen(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionMetadata(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionForwardHeaders(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionWireFormat(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil || tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
		return "EXT_PROC_GRPC"
	}

	return v
}

func flattenNetworkServicesAuthzExtensionTerraformLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}

	transformed := make(map[string]interface{})
	if l, ok := d.GetOkExists("terraform_labels"); ok {
		for k := range l.(map[string]interface{}) {
			transformed[k] = v.(map[string]interface{})[k]
		}
	}

	return transformed
}

func flattenNetworkServicesAuthzExtensionEffectiveLabels(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenNetworkServicesAuthzExtensionName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}

func expandNetworkServicesAuthzExtensionDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionLoadBalancingScheme(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionAuthority(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionTimeout(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionFailOpen(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionMetadata(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandNetworkServicesAuthzExtensionForwardHeaders(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionWireFormat(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandNetworkServicesAuthzExtensionEffectiveLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandNetworkServicesAuthzExtensionName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return fmt.Sprintf("projects/%s/locations/%s/authzExtensions/%s", d.Get("project"), d.Get("location"), v), nil
}
