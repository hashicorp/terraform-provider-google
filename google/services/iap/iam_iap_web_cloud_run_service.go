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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/iap/CloudRunService.yaml
//     Template:      https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/iam_policy.go.tmpl
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------

package iap

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var IapWebCloudRunServiceIamSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"location": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"cloud_run_service_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
	},
}

type IapWebCloudRunServiceIamUpdater struct {
	project             string
	location            string
	cloudRunServiceName string
	d                   tpgresource.TerraformResourceData
	Config              *transport_tpg.Config
}

func IapWebCloudRunServiceIamUpdaterProducer(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	values := make(map[string]string)

	project, _ := tpgresource.GetProject(d, config)
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
	}
	values["project"] = project
	location, _ := tpgresource.GetLocation(d, config)
	if location != "" {
		if err := d.Set("location", location); err != nil {
			return nil, fmt.Errorf("Error setting location: %s", err)
		}
	}
	values["location"] = location
	if v, ok := d.GetOk("cloud_run_service_name"); ok {
		values["cloud_run_service_name"] = v.(string)
	}

	// We may have gotten either a long or short name, so attempt to parse long name if possible
	m, err := tpgresource.GetImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/iap_web/cloud_run-(?P<location>[^/]+)/services/(?P<cloud_run_service_name>[^/]+)", "(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<cloud_run_service_name>[^/]+)", "(?P<location>[^/]+)/(?P<cloud_run_service_name>[^/]+)", "(?P<cloud_run_service_name>[^/]+)"}, d, config, d.Get("cloud_run_service_name").(string))
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		values[k] = v
	}

	u := &IapWebCloudRunServiceIamUpdater{
		project:             values["project"],
		location:            values["location"],
		cloudRunServiceName: values["cloud_run_service_name"],
		d:                   d,
		Config:              config,
	}

	if err := d.Set("project", u.project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("location", u.location); err != nil {
		return nil, fmt.Errorf("Error setting location: %s", err)
	}
	if err := d.Set("cloud_run_service_name", u.GetResourceId()); err != nil {
		return nil, fmt.Errorf("Error setting cloud_run_service_name: %s", err)
	}

	return u, nil
}

func IapWebCloudRunServiceIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	values := make(map[string]string)

	project, _ := tpgresource.GetProject(d, config)
	if project != "" {
		values["project"] = project
	}

	location, _ := tpgresource.GetLocation(d, config)
	if location != "" {
		values["location"] = location
	}

	m, err := tpgresource.GetImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/iap_web/cloud_run-(?P<location>[^/]+)/services/(?P<cloud_run_service_name>[^/]+)", "(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<cloud_run_service_name>[^/]+)", "(?P<location>[^/]+)/(?P<cloud_run_service_name>[^/]+)", "(?P<cloud_run_service_name>[^/]+)"}, d, config, d.Id())
	if err != nil {
		return err
	}

	for k, v := range m {
		values[k] = v
	}

	u := &IapWebCloudRunServiceIamUpdater{
		project:             values["project"],
		location:            values["location"],
		cloudRunServiceName: values["cloud_run_service_name"],
		d:                   d,
		Config:              config,
	}
	if err := d.Set("cloud_run_service_name", u.GetResourceId()); err != nil {
		return fmt.Errorf("Error setting cloud_run_service_name: %s", err)
	}
	d.SetId(u.GetResourceId())
	return nil
}

func (u *IapWebCloudRunServiceIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	url, err := u.qualifyWebCloudRunServiceUrl("getIamPolicy")
	if err != nil {
		return nil, err
	}

	project, err := tpgresource.GetProject(u.d, u.Config)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	obj = map[string]interface{}{
		"options": map[string]interface{}{
			"requestedPolicyVersion": tpgiamresource.IamPolicyVersion,
		},
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	policy, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    u.Config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
	})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	out := &cloudresourcemanager.Policy{}
	err = tpgresource.Convert(policy, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a policy to a resource manager policy: {{err}}", err)
	}

	return out, nil
}

func (u *IapWebCloudRunServiceIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	json, err := tpgresource.ConvertToMap(policy)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	obj["policy"] = json

	url, err := u.qualifyWebCloudRunServiceUrl("setIamPolicy")
	if err != nil {
		return err
	}
	project, err := tpgresource.GetProject(u.d, u.Config)
	if err != nil {
		return err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    u.Config,
		Method:    "POST",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   u.d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *IapWebCloudRunServiceIamUpdater) qualifyWebCloudRunServiceUrl(methodIdentifier string) (string, error) {
	urlTemplate := fmt.Sprintf("{{IapBasePath}}%s:%s", fmt.Sprintf("projects/%s/iap_web/cloud_run-%s/services/%s", u.project, u.location, u.cloudRunServiceName), methodIdentifier)
	url, err := tpgresource.ReplaceVars(u.d, u.Config, urlTemplate)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (u *IapWebCloudRunServiceIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/iap_web/cloud_run-%s/services/%s", u.project, u.location, u.cloudRunServiceName)
}

func (u *IapWebCloudRunServiceIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-iap-webcloudrunservice-%s", u.GetResourceId())
}

func (u *IapWebCloudRunServiceIamUpdater) DescribeResource() string {
	return fmt.Sprintf("iap webcloudrunservice %q", u.GetResourceId())
}
