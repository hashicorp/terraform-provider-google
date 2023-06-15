// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV3 "google.golang.org/api/cloudresourcemanager/v3"
)

var IamFolderSchema = map[string]*schema.Schema{
	"folder": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type FolderIamUpdater struct {
	folderId string
	d        tpgresource.TerraformResourceData
	Config   *transport_tpg.Config
}

func NewFolderIamUpdater(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	return &FolderIamUpdater{
		folderId: CanonicalFolderId(d.Get("folder").(string)),
		d:        d,
		Config:   config,
	}, nil
}

func FolderIdParseFunc(d *schema.ResourceData, _ *transport_tpg.Config) error {
	if !strings.HasPrefix(d.Id(), "folders/") {
		d.SetId(fmt.Sprintf("folders/%s", d.Id()))
	}
	if err := d.Set("folder", d.Id()); err != nil {
		return fmt.Errorf("Error setting folder: %s", err)
	}
	return nil
}

func (u *FolderIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	return GetFolderIamPolicyByFolderName(u.folderId, userAgent, u.Config)
}

func (u *FolderIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	v2Policy, err := v1PolicyToV2(policy)
	if err != nil {
		return err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = u.Config.NewResourceManagerV3Client(userAgent).Folders.SetIamPolicy(u.folderId, &resourceManagerV3.SetIamPolicyRequest{
		Policy:     v2Policy,
		UpdateMask: "bindings,etag,auditConfigs",
	}).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *FolderIamUpdater) GetResourceId() string {
	return u.folderId
}

func (u *FolderIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-folder-%s", u.folderId)
}

func (u *FolderIamUpdater) DescribeResource() string {
	return fmt.Sprintf("folder %q", u.folderId)
}

func CanonicalFolderId(folder string) string {
	if strings.HasPrefix(folder, "folders/") {
		return folder
	}

	return "folders/" + folder
}

// v1 and v2 policy are identical
func v1PolicyToV2(in *cloudresourcemanager.Policy) (*resourceManagerV3.Policy, error) {
	out := &resourceManagerV3.Policy{}
	err := tpgresource.Convert(in, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a v2 policy: {{err}}", err)
	}
	return out, nil
}

func v2PolicyToV1(in *resourceManagerV3.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := tpgresource.Convert(in, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v2 policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}

// Retrieve the existing IAM Policy for a folder
func GetFolderIamPolicyByFolderName(folderName, userAgent string, config *transport_tpg.Config) (*cloudresourcemanager.Policy, error) {
	p, err := config.NewResourceManagerV3Client(userAgent).Folders.GetIamPolicy(folderName,
		&resourceManagerV3.GetIamPolicyRequest{
			Options: &resourceManagerV3.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for folder %q: {{err}}", folderName), err)
	}

	v1Policy, err := v2PolicyToV1(p)
	if err != nil {
		return nil, err
	}

	return v1Policy, nil
}
