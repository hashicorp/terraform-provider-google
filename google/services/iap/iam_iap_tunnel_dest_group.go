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
//     Configuration: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/iap/TunnelDestGroup.yaml
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

var IapTunnelDestGroupIamSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"region": {
		Type:     schema.TypeString,
		Computed: true,
		Optional: true,
		ForceNew: true,
	},
	"dest_group": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
	},
}

type IapTunnelDestGroupIamUpdater struct {
	project   string
	region    string
	destGroup string
	d         tpgresource.TerraformResourceData
	Config    *transport_tpg.Config
}

func IapTunnelDestGroupIamUpdaterProducer(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	values := make(map[string]string)

	project, _ := tpgresource.GetProject(d, config)
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
	}
	values["project"] = project
	region, _ := tpgresource.GetRegion(d, config)
	if region != "" {
		if err := d.Set("region", region); err != nil {
			return nil, fmt.Errorf("Error setting region: %s", err)
		}
	}
	values["region"] = region
	if v, ok := d.GetOk("dest_group"); ok {
		values["dest_group"] = v.(string)
	}

	// We may have gotten either a long or short name, so attempt to parse long name if possible
	m, err := tpgresource.GetImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/iap_tunnel/locations/(?P<region>[^/]+)/destGroups/(?P<dest_group>[^/]+)", "(?P<project>[^/]+)/iap_tunnel/locations/(?P<region>[^/]+)/destGroups/(?P<dest_group>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<dest_group>[^/]+)", "(?P<region>[^/]+)/(?P<dest_group>[^/]+)", "(?P<dest_group>[^/]+)"}, d, config, d.Get("dest_group").(string))
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		values[k] = v
	}

	u := &IapTunnelDestGroupIamUpdater{
		project:   values["project"],
		region:    values["region"],
		destGroup: values["dest_group"],
		d:         d,
		Config:    config,
	}

	if err := d.Set("project", u.project); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("region", u.region); err != nil {
		return nil, fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("dest_group", u.GetResourceId()); err != nil {
		return nil, fmt.Errorf("Error setting dest_group: %s", err)
	}

	return u, nil
}

func IapTunnelDestGroupIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	values := make(map[string]string)

	project, _ := tpgresource.GetProject(d, config)
	if project != "" {
		values["project"] = project
	}

	region, _ := tpgresource.GetRegion(d, config)
	if region != "" {
		values["region"] = region
	}

	m, err := tpgresource.GetImportIdQualifiers([]string{"projects/(?P<project>[^/]+)/iap_tunnel/locations/(?P<region>[^/]+)/destGroups/(?P<dest_group>[^/]+)", "(?P<project>[^/]+)/iap_tunnel/locations/(?P<region>[^/]+)/destGroups/(?P<dest_group>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<dest_group>[^/]+)", "(?P<region>[^/]+)/(?P<dest_group>[^/]+)", "(?P<dest_group>[^/]+)"}, d, config, d.Id())
	if err != nil {
		return err
	}

	for k, v := range m {
		values[k] = v
	}

	u := &IapTunnelDestGroupIamUpdater{
		project:   values["project"],
		region:    values["region"],
		destGroup: values["dest_group"],
		d:         d,
		Config:    config,
	}
	if err := d.Set("dest_group", u.GetResourceId()); err != nil {
		return fmt.Errorf("Error setting dest_group: %s", err)
	}
	d.SetId(u.GetResourceId())
	return nil
}

func (u *IapTunnelDestGroupIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	url, err := u.qualifyTunnelDestGroupUrl("getIamPolicy")
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

func (u *IapTunnelDestGroupIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	json, err := tpgresource.ConvertToMap(policy)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	obj["policy"] = json

	url, err := u.qualifyTunnelDestGroupUrl("setIamPolicy")
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

func (u *IapTunnelDestGroupIamUpdater) qualifyTunnelDestGroupUrl(methodIdentifier string) (string, error) {
	urlTemplate := fmt.Sprintf("{{IapBasePath}}%s:%s", fmt.Sprintf("projects/%s/iap_tunnel/locations/%s/destGroups/%s", u.project, u.region, u.destGroup), methodIdentifier)
	url, err := tpgresource.ReplaceVars(u.d, u.Config, urlTemplate)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (u *IapTunnelDestGroupIamUpdater) GetResourceId() string {
	return fmt.Sprintf("projects/%s/iap_tunnel/locations/%s/destGroups/%s", u.project, u.region, u.destGroup)
}

func (u *IapTunnelDestGroupIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-iap-tunneldestgroup-%s", u.GetResourceId())
}

func (u *IapTunnelDestGroupIamUpdater) DescribeResource() string {
	return fmt.Sprintf("iap tunneldestgroup %q", u.GetResourceId())
}
