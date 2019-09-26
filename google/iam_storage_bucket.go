package google

import (
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/storage/v1"
)

var IamStorageBucketSchema = map[string]*schema.Schema{
	"bucket": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

func StorageBucketIdParseFunc(d *schema.ResourceData, _ *Config) error {
	d.Set("bucket", d.Id())
	return nil
}

type StorageBucketIamUpdater struct {
	bucket string
	Config *Config
}

func NewStorageBucketIamUpdater(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error) {
	bucket := d.Get("bucket").(string)

	return &StorageBucketIamUpdater{
		bucket: bucket,
		Config: config,
	}, nil
}

func (u *StorageBucketIamUpdater) GetResourceIamPolicy() (*cloudresourcemanager.Policy, error) {
	p, err := u.Config.clientStorage.Buckets.GetIamPolicy(u.bucket).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	cloudResourcePolicy, err := storageToResourceManagerPolicy(p)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return cloudResourcePolicy, nil
}

func (u *StorageBucketIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	storagePolicy, err := resourceManagerToStoragePolicy(policy)

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Invalid IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	ppolicy, err := u.Config.clientStorage.Buckets.GetIamPolicy(u.bucket).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}
	storagePolicy.Etag = ppolicy.Etag
	_, err = u.Config.clientStorage.Buckets.SetIamPolicy(u.bucket, storagePolicy).Do()

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error setting IAM policy for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *StorageBucketIamUpdater) GetResourceId() string {
	return u.bucket
}

func (u *StorageBucketIamUpdater) GetMutexKey() string {
	return fmt.Sprintf("iam-storage-bucket-%s", u.bucket)
}

func (u *StorageBucketIamUpdater) DescribeResource() string {
	return fmt.Sprintf("Storage Bucket %q", u.bucket)
}

func resourceManagerToStoragePolicy(p *cloudresourcemanager.Policy) (*storage.Policy, error) {
	out := &storage.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a v1 policy to a storage policy: {{err}}", err)
	}
	return out, nil
}

func storageToResourceManagerPolicy(p *storage.Policy) (*cloudresourcemanager.Policy, error) {
	out := &cloudresourcemanager.Policy{}
	err := Convert(p, out)
	if err != nil {
		return nil, errwrap.Wrapf("Cannot convert a storage policy to a v1 policy: {{err}}", err)
	}
	return out, nil
}
