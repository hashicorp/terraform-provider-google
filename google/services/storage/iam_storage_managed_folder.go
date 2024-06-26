// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

var StorageManagedFolderIamSchema = map[string]*schema.Schema{
	"bucket": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"managed_folder": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
		ValidateFunc:     verify.ValidateRegexp(`/$`),
	},
}

type StorageManagedFolderIamUpdater struct {
	bucket        string
	managedFolder string
	d             tpgresource.TerraformResourceData
	Config        *transport_tpg.Config
}

func StorageManagedFolderIamUpdaterProducer(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (tpgiamresource.ResourceIamUpdater, error) {
	values := make(map[string]string)

	if v, ok := d.GetOk("bucket"); ok {
		values["bucket"] = v.(string)
	}

	if v, ok := d.GetOk("managed_folder"); ok {
		values["managed_folder"] = v.(string)
	}

	u := &StorageManagedFolderIamUpdater{
		bucket:        values["bucket"],
		managedFolder: values["managed_folder"],
		d:             d,
		Config:        config,
	}

	if err := d.Set("bucket", u.bucket); err != nil {
		return nil, fmt.Errorf("Error setting bucket: %s", err)
	}
	if err := d.Set("managed_folder", u.managedFolder); err != nil {
		return nil, fmt.Errorf("Error setting managed_folder: %s", err)
	}

	return u, nil
}

func StorageManagedFolderIdParseFunc(d *schema.ResourceData, config *transport_tpg.Config) error {
	values := make(map[string]string)

	m, err := tpgresource.GetImportIdQualifiers([]string{"(?P<bucket>[^/]+)/managedFolders/(?P<managed_folder>.+)", "(?P<bucket>[^/]+)/(?P<managed_folder>.+)"}, d, config, d.Id())
	if err != nil {
		return err
	}

	for k, v := range m {
		values[k] = v
	}

	u := &StorageManagedFolderIamUpdater{
		bucket:        values["bucket"],
		managedFolder: values["managed_folder"],
		d:             d,
		Config:        config,
	}
	if err := d.Set("bucket", u.bucket); err != nil {
		return fmt.Errorf("Error setting bucket: %s", err)
	}
	if err := d.Set("managed_folder", u.managedFolder); err != nil {
		return fmt.Errorf("Error setting managed_folder: %s", err)
	}
	d.SetId(u.GetResourceId())
	return nil
}

func (u *StorageManagedFolderIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	url, err := u.qualifyManagedFolderUrl("iam")
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"optionsRequestedPolicyVersion": fmt.Sprintf("%d", tpgiamresource.IamPolicyVersion)})
	if err != nil {
		return nil, err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return nil, err
	}

	policy, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    u.Config,
		Method:    "GET",
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

func (u *StorageManagedFolderIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	json, err := tpgresource.ConvertToMap(policy)
	if err != nil {
		return err
	}

	obj := json

	url, err := u.qualifyManagedFolderUrl("iam")
	if err != nil {
		return err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(u.d, u.Config.UserAgent)
	if err != nil {
		return err
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    u.Config,
		Method:    "PUT",
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

func (u *StorageManagedFolderIamUpdater) qualifyManagedFolderUrl(methodIdentifier string) (string, error) {
	urlTemplate := fmt.Sprintf("{{StorageBasePath}}b/%s/managedFolders/%s/%s", u.bucket, url.PathEscape(u.managedFolder), methodIdentifier)
	url, err := tpgresource.ReplaceVars(u.d, u.Config, urlTemplate)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (u *StorageManagedFolderIamUpdater) GetResourceId() string {
	return fmt.Sprintf("b/%s/managedFolders/%s", u.bucket, u.managedFolder)
}

func (u *StorageManagedFolderIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-storage-managedfolder-%s", u.GetResourceId())
}

func (u *StorageManagedFolderIamUpdater) DescribeResource() string {
	return fmt.Sprintf("storage managedfolder %q", u.GetResourceId())
}
