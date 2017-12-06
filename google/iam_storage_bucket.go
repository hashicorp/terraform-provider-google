package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/storage/v1"
)

var IamStorageBucketSchema = map[string]*schema.Schema{
	"bucket": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
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
		return nil, fmt.Errorf("Error retrieving IAM policy for %s: %s", u.DescribeResource(), err)
	}

	cloudResourcePolicy, err := storageToResourceManagerPolicy(p)
	if err != nil {
		return nil, fmt.Errorf("Invalid IAM policy for %s: %s", u.DescribeResource(), err)
	}

	return cloudResourcePolicy, nil
}

func (u *StorageBucketIamUpdater) SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error {
	storagePolicy, err := resourceManagerToStoragePolicy(policy)

	if err != nil {
		return fmt.Errorf("Invalid IAM policy for %s: %s", u.DescribeResource(), err)
	}

	_, err = u.Config.clientStorage.Buckets.SetIamPolicy(u.bucket, storagePolicy).Do()

	if err != nil {
		return fmt.Errorf("Error setting IAM policy for %s: %s", u.DescribeResource(), err)
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

func resourceManagerToStoragePolicy(p *cloudresourcemanager.Policy) (policy *storage.Policy, err error) {
	policy = &storage.Policy{}
	err = Convert(p, policy)
	return
}

func storageToResourceManagerPolicy(p *storage.Policy) (policy *cloudresourcemanager.Policy, err error) {
	policy = &cloudresourcemanager.Policy{}
	err = Convert(p, policy)
	return
}
