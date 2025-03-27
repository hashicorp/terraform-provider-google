// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanagerregional

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceParameterManagerRegionalRegionalParameterVersionRender() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceParameterManagerRegionalRegionalParameterVersionRenderRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": {
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

func dataSourceParameterManagerRegionalRegionalParameterVersionRenderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	// Check if the parameter is provided as a resource reference or a parameter id.
	parameterRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/parameters/(.+)$")
	dParameter, ok := d.Get("parameter").(string)
	if !ok {
		return fmt.Errorf("wrong type for parameter field (%T), expected string", d.Get("parameter"))
	}

	parts := parameterRegex.FindStringSubmatch(dParameter)
	var project string

	// if reference of the regional parameter is provided in the parameter field
	if len(parts) == 4 {
		// Stores value of project to set in state
		project = parts[1]
		if dProject, ok := d.Get("project").(string); !ok {
			return fmt.Errorf("wrong type for project (%T), expected string", d.Get("project"))
		} else if dProject != "" && dProject != project {
			return fmt.Errorf("project field value (%s) does not match project of regional parameter (%s).", dProject, project)
		}

		if dLocation, ok := d.Get("location").(string); !ok {
			return fmt.Errorf("wrong type for location (%T), expected string", d.Get("location"))
		} else if dLocation != "" && dLocation != parts[2] {
			return fmt.Errorf("location field value (%s) does not match location of regional parameter (%s).", dLocation, parts[2])
		}

		if err := d.Set("location", parts[2]); err != nil {
			return fmt.Errorf("error setting location: %s", err)
		}
		if err := d.Set("parameter", parts[3]); err != nil {
			return fmt.Errorf("error setting parameter: %s", err)
		}
	} else { // if regional parameter name is provided in the parameter field
		// Stores value of project to set in state
		project, err = tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("error fetching project for regional parameter: %s", err)
		}
		if dLocation, ok := d.Get("location").(string); ok && dLocation == "" {
			return fmt.Errorf("location must be set when providing only regional parameter name")
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	dParameterVersionId, ok := d.Get("parameter_version_id").(string)
	if !ok {
		return fmt.Errorf("wrong type for parameter version id field (%T), expected string", d.Get("parameter_version_id"))
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerRegionalBasePath}}projects/{{project}}/locations/{{location}}/parameters/{{parameter}}/versions/{{parameter_version_id}}:render")
	if err != nil {
		return err
	}

	headers := make(http.Header)
	regionalParameterVersion, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("error retrieving available parameter manager regional parameter version: %s", err.Error())
	}

	// If the response contains the disabled value, return an error stating that the regional parameter version is currently disabled
	isDisabled, ok := regionalParameterVersion["disabled"]
	if ok && isDisabled.(bool) {
		return fmt.Errorf("regional parameter version %s is in DISABLED state.", dParameterVersionId)
	}

	nameValue, ok := regionalParameterVersion["parameterVersion"]
	if !ok {
		return fmt.Errorf("read response didn't contain critical fields. Read may not have succeeded.")
	}

	if err := d.Set("name", nameValue.(string)); err != nil {
		return fmt.Errorf("error reading regionalParameterVersion: %s", err)
	}

	if err := d.Set("disabled", false); err != nil {
		return fmt.Errorf("error setting disabled: %s", err)
	}

	data := regionalParameterVersion["payload"].(map[string]interface{})
	parameterData, err := base64.StdEncoding.DecodeString(data["data"].(string))
	if err != nil {
		return fmt.Errorf("error decoding parameter manager regional parameter version data: %s", err.Error())
	}
	if err := d.Set("parameter_data", string(parameterData)); err != nil {
		return fmt.Errorf("error setting parameter_data: %s", err)
	}

	renderedParameterData, err := base64.StdEncoding.DecodeString(regionalParameterVersion["renderedPayload"].(string))
	if err != nil {
		return fmt.Errorf("error decoding parameter manager regional parameter version rendered payload data: %s", err.Error())
	}
	if err := d.Set("rendered_parameter_data", string(renderedParameterData)); err != nil {
		return fmt.Errorf("error setting rendered_parameter_data: %s", err)
	}
	d.SetId(nameValue.(string))
	return nil
}
