package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
	"google.golang.org/api/googleapi"
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
	Config   *Config
}

func NewFolderIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	return &FolderIamUpdater{
		folderId: canonicalFolderId(d.Get("folder").(string)),
		Config:   config,
	}, nil
}

func FolderIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("folder", d.Id())
	return nil
}

func (u *FolderIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	return getFolderIamPolicyByFolderName(u.folderId, u.Config)
}

func (u *FolderIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	v2BetaPolicy, err := v1PolicyToV2Beta(policy)
	if err != nil {
		return err
	}

	_, err = u.Config.clientResourceManagerV2Beta1.Folders.SetIamPolicy(u.folderId, &resourceManagerV2Beta1.SetIamPolicyRequest{
		Policy: v2BetaPolicy,
	}).Do()

	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error setting IAM policy for %s.", u.DescribeResource()), err)
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

func canonicalFolderId(folder string) string {
	if strings.HasPrefix(folder, "folders/") {
		return folder
	}

	return "folders/" + folder
}

// v1 and v2beta policy are identical
func v1PolicyToV2Beta(in *cloudresourcemanager.Policy) (*resourceManagerV2Beta1.Policy, error) {
	out := &resourceManagerV2Beta1.Policy{}
	err := Convert(in, out)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert a v1 policy to a v2beta policy: %s", err)
	}
	return out, nil
}

func v2BetaPolicyToV1(in *resourceManagerV2Beta1.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(in, out)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert a v2beta policy to a v1 policy: %s", err)
	}
	return out, nil
}

// Retrieve the existing IAM Policy for a folder
func getFolderIamPolicyByFolderName(folderName string, config *Config) (*cloudresourcemanager.Policy, error) {
	p, err := config.clientResourceManagerV2Beta1.Folders.GetIamPolicy(folderName,
		&resourceManagerV2Beta1.GetIamPolicyRequest{}).Do()

	if err != nil {
		return nil, fmt.Errorf("Error retrieving IAM policy for folder %q: %s", folderName, err)
	}

	v1Policy, err := v2BetaPolicyToV1(p)
	if err != nil {
		return nil, err
	}

	return v1Policy, nil
}

func getFolderIamPolicyByParentAndDisplayName(parent, displayName string, config *Config) (*cloudresourcemanager.Policy, error) {
	queryString := fmt.Sprintf("lifecycleState=ACTIVE AND parent=%s AND displayName=%s", parent, displayName)
	searchRequest := &resourceManagerV2Beta1.SearchFoldersRequest{
		Query: queryString,
	}
	searchResponse, err := config.clientResourceManagerV2Beta1.Folders.Search(searchRequest).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return nil, fmt.Errorf("Folder Not Found: %s,%s", parent, displayName)
		}

		return nil, fmt.Errorf("Error reading folders: %s", err)
	}

	folders := searchResponse.Folders
	if len(folders) != 1 {
		return nil, fmt.Errorf("More than one folder found")
	}

	return getFolderIamPolicyByFolderName(folders[0].Name, config)
}
