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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/securityposture/PostureDeployment.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package securityposture

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceSecurityposturePostureDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityposturePostureDeploymentCreate,
		Read:   resourceSecurityposturePostureDeploymentRead,
		Update: resourceSecurityposturePostureDeploymentUpdate,
		Delete: resourceSecurityposturePostureDeploymentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSecurityposturePostureDeploymentImport,
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
				Description: `The location of the resource, eg. global'.`,
			},
			"parent": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The parent of the resource, an organization. Format should be 'organizations/{organization_id}'.`,
			},
			"posture_deployment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `ID of the posture deployment.`,
			},
			"posture_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Relative name of the posture which needs to be deployed. It should be in the format:
  organizations/{organization_id}/locations/{location}/postures/{posture_id}`,
			},
			"posture_revision_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Revision_id the posture which needs to be deployed.`,
			},
			"target_resource": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The resource on which the posture should be deployed. This can be in one of the following formats:
projects/{project_number},
folders/{folder_number},
organizations/{organization_id}`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Description of the posture deployment.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time the posture deployment was created in UTC.`,
			},
			"desired_posture_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `This is an output only optional field which will be filled in case when
PostureDeployment state is UPDATE_FAILED or CREATE_FAILED or DELETE_FAILED.
It denotes the desired posture to be deployed.`,
			},
			"desired_posture_revision_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `This is an output only optional field which will be filled in case when
PostureDeployment state is UPDATE_FAILED or CREATE_FAILED or DELETE_FAILED.
It denotes the desired posture revision_id to be deployed.`,
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `For Resource freshness validation (https://google.aip.dev/154)`,
			},
			"failure_message": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `This is a output only optional field which will be filled in case where
PostureDeployment enters a failure state like UPDATE_FAILED or
CREATE_FAILED or DELETE_FAILED. It will have the failure message for posture deployment's
CREATE/UPDATE/DELETE methods.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Name of the posture deployment instance.`,
			},
			"reconciling": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `If set, there are currently changes in flight to the posture deployment.`,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `State of the posture deployment. A posture deployment can be in the following terminal states:
ACTIVE, CREATE_FAILED, UPDATE_FAILED, DELETE_FAILED.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time the posture deployment was updated in UTC.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSecurityposturePostureDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	targetResourceProp, err := expandSecurityposturePostureDeploymentTargetResource(d.Get("target_resource"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("target_resource"); !tpgresource.IsEmptyValue(reflect.ValueOf(targetResourceProp)) && (ok || !reflect.DeepEqual(v, targetResourceProp)) {
		obj["targetResource"] = targetResourceProp
	}
	postureIdProp, err := expandSecurityposturePostureDeploymentPostureId(d.Get("posture_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("posture_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(postureIdProp)) && (ok || !reflect.DeepEqual(v, postureIdProp)) {
		obj["postureId"] = postureIdProp
	}
	postureRevisionIdProp, err := expandSecurityposturePostureDeploymentPostureRevisionId(d.Get("posture_revision_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("posture_revision_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(postureRevisionIdProp)) && (ok || !reflect.DeepEqual(v, postureRevisionIdProp)) {
		obj["postureRevisionId"] = postureRevisionIdProp
	}
	descriptionProp, err := expandSecurityposturePostureDeploymentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecuritypostureBasePath}}{{parent}}/locations/{{location}}/postureDeployments?postureDeploymentId={{posture_deployment_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new PostureDeployment: %#v", obj)
	billingProject := ""

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
		return fmt.Errorf("Error creating PostureDeployment: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/locations/{{location}}/postureDeployments/{{posture_deployment_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = SecuritypostureOperationWaitTime(
		config, res, "Creating PostureDeployment", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create PostureDeployment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating PostureDeployment %q: %#v", d.Id(), res)

	return resourceSecurityposturePostureDeploymentRead(d, meta)
}

func resourceSecurityposturePostureDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecuritypostureBasePath}}{{parent}}/locations/{{location}}/postureDeployments/{{posture_deployment_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

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
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SecurityposturePostureDeployment %q", d.Id()))
	}

	if err := d.Set("name", flattenSecurityposturePostureDeploymentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("target_resource", flattenSecurityposturePostureDeploymentTargetResource(res["targetResource"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("state", flattenSecurityposturePostureDeploymentState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("posture_id", flattenSecurityposturePostureDeploymentPostureId(res["postureId"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("posture_revision_id", flattenSecurityposturePostureDeploymentPostureRevisionId(res["postureRevisionId"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("create_time", flattenSecurityposturePostureDeploymentCreateTime(res["createTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("update_time", flattenSecurityposturePostureDeploymentUpdateTime(res["updateTime"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("description", flattenSecurityposturePostureDeploymentDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("etag", flattenSecurityposturePostureDeploymentEtag(res["etag"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("reconciling", flattenSecurityposturePostureDeploymentReconciling(res["reconciling"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("desired_posture_id", flattenSecurityposturePostureDeploymentDesiredPostureId(res["desiredPostureId"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("desired_posture_revision_id", flattenSecurityposturePostureDeploymentDesiredPostureRevisionId(res["desiredPostureRevisionId"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}
	if err := d.Set("failure_message", flattenSecurityposturePostureDeploymentFailureMessage(res["failureMessage"], d, config)); err != nil {
		return fmt.Errorf("Error reading PostureDeployment: %s", err)
	}

	return nil
}

func resourceSecurityposturePostureDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	obj := make(map[string]interface{})
	postureIdProp, err := expandSecurityposturePostureDeploymentPostureId(d.Get("posture_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("posture_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, postureIdProp)) {
		obj["postureId"] = postureIdProp
	}
	postureRevisionIdProp, err := expandSecurityposturePostureDeploymentPostureRevisionId(d.Get("posture_revision_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("posture_revision_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, postureRevisionIdProp)) {
		obj["postureRevisionId"] = postureRevisionIdProp
	}
	descriptionProp, err := expandSecurityposturePostureDeploymentDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SecuritypostureBasePath}}{{parent}}/locations/{{location}}/postureDeployments/{{posture_deployment_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating PostureDeployment %q: %#v", d.Id(), obj)
	headers := make(http.Header)
	updateMask := []string{}

	if d.HasChange("posture_id") {
		updateMask = append(updateMask, "postureId")
	}

	if d.HasChange("posture_revision_id") {
		updateMask = append(updateMask, "postureRevisionId")
	}

	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
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
			return fmt.Errorf("Error updating PostureDeployment %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating PostureDeployment %q: %#v", d.Id(), res)
		}

		err = SecuritypostureOperationWaitTime(
			config, res, "Updating PostureDeployment", userAgent,
			d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return err
		}
	}

	return resourceSecurityposturePostureDeploymentRead(d, meta)
}

func resourceSecurityposturePostureDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := tpgresource.ReplaceVars(d, config, "{{SecuritypostureBasePath}}{{parent}}/locations/{{location}}/postureDeployments/{{posture_deployment_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)

	log.Printf("[DEBUG] Deleting PostureDeployment %q", d.Id())
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
		return transport_tpg.HandleNotFoundError(err, d, "PostureDeployment")
	}

	err = SecuritypostureOperationWaitTime(
		config, res, "Deleting PostureDeployment", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting PostureDeployment %q: %#v", d.Id(), res)
	return nil
}

func resourceSecurityposturePostureDeploymentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^(?P<parent>.+)/locations/(?P<location>[^/]+)/postureDeployments/(?P<posture_deployment_id>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/locations/{{location}}/postureDeployments/{{posture_deployment_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenSecurityposturePostureDeploymentName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentTargetResource(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentState(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentPostureId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentPostureRevisionId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentCreateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentUpdateTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentEtag(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentReconciling(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentDesiredPostureId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentDesiredPostureRevisionId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSecurityposturePostureDeploymentFailureMessage(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandSecurityposturePostureDeploymentTargetResource(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityposturePostureDeploymentPostureId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityposturePostureDeploymentPostureRevisionId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSecurityposturePostureDeploymentDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
