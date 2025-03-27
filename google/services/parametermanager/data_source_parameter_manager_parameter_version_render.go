// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanager

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceParameterManagerParameterVersionRender() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceParameterManagerParameterVersionRenderRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"parameter": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
			},
			"parameter_version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameter_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"rendered_parameter_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceParameterManagerParameterVersionRenderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	// Check if the parameter is provided as a resource reference or a parameter id.
	parameterRegex := regexp.MustCompile("projects/(.+)/locations/global/parameters/(.+)$")
	dParameter, ok := d.Get("parameter").(string)
	if !ok {
		return fmt.Errorf("wrong type for parameter field (%T), expected string", d.Get("parameter"))
	}

	parts := parameterRegex.FindStringSubmatch(dParameter)
	var project string

	// if reference of the parameter is provided in the parameter field
	if len(parts) == 3 {
		// Stores value of project to set in state
		project = parts[1]
		if dProject, ok := d.Get("project").(string); !ok {
			return fmt.Errorf("wrong type for project (%T), expected string", d.Get("project"))
		} else if dProject != "" && dProject != project {
			return fmt.Errorf("project field value (%s) does not match project of parameter (%s).", dProject, project)
		}
		if err := d.Set("parameter", parts[2]); err != nil {
			return fmt.Errorf("error setting parameter: %s", err)
		}
	} else { // if parameter name is provided in the parameter field
		// Stores value of project to set in state
		project, err = tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("error fetching project for parameter: %s", err)
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	dParameterVersionId, ok := d.Get("parameter_version_id").(string)
	if !ok {
		return fmt.Errorf("wrong type for parameter version id field (%T), expected string", d.Get("parameter_version_id"))
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerBasePath}}projects/{{project}}/locations/global/parameters/{{parameter}}/versions/{{parameter_version_id}}:render")
	if err != nil {
		return err
	}

	headers := make(http.Header)
	parameterVersion, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("error retrieving available parameter manager parameter version: %s", err.Error())
	}

	// If the response contains the disabled value, return an error stating that the parameter version is currently disabled
	isDisabled, ok := parameterVersion["disabled"]
	if ok && isDisabled.(bool) {
		return fmt.Errorf("parameter version %s is in DISABLED state.", dParameterVersionId)
	}

	nameValue, ok := parameterVersion["parameterVersion"]
	if !ok {
		return fmt.Errorf("read response didn't contain critical fields. Read may not have succeeded.")
	}

	if err := d.Set("name", nameValue.(string)); err != nil {
		return fmt.Errorf("error reading parameterVersion: %s", err)
	}

	if err := d.Set("disabled", false); err != nil {
		return fmt.Errorf("error setting disabled: %s", err)
	}

	data := parameterVersion["payload"].(map[string]interface{})
	parameterData, err := base64.StdEncoding.DecodeString(data["data"].(string))
	if err != nil {
		return fmt.Errorf("error decoding parameter manager parameter version data: %s", err.Error())
	}
	if err := d.Set("parameter_data", string(parameterData)); err != nil {
		return fmt.Errorf("error setting parameter_data: %s", err)
	}

	renderedParameterData, err := base64.StdEncoding.DecodeString(parameterVersion["renderedPayload"].(string))
	if err != nil {
		return fmt.Errorf("error decoding parameter manager parameter version rendered payload data: %s", err.Error())
	}
	if err := d.Set("rendered_parameter_data", string(renderedParameterData)); err != nil {
		return fmt.Errorf("error setting rendered_parameter_data: %s", err)
	}
	d.SetId(nameValue.(string))
	return nil
}
