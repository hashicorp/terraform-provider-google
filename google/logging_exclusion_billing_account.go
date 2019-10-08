package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/logging/v2"
)

var BillingAccountLoggingExclusionSchema = map[string]*schema.Schema{
	"billing_account": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

type BillingAccountLoggingExclusionUpdater struct {
	resourceType string
	resourceId   string
	Config       *Config
}

func NewBillingAccountLoggingExclusionUpdater(d *schema.ResourceData, config *Config) (ResourceLoggingExclusionUpdater, error) {
	billingAccount := d.Get("billing_account").(string)

	return &BillingAccountLoggingExclusionUpdater{
		resourceType: "billingAccounts",
		resourceId:   billingAccount,
		Config:       config,
	}, nil
}

func billingAccountLoggingExclusionIdParseFunc(d *schema.ResourceData, _ *Config) error {
	loggingExclusionId, err := parseLoggingExclusionId(d.Id())
	if err != nil {
		return err
	}

	if "billingAccounts" != loggingExclusionId.resourceType {
		return fmt.Errorf("Error importing logging exclusion, invalid resourceType %#v", loggingExclusionId.resourceType)
	}

	d.Set("billing_account", loggingExclusionId.resourceId)
	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) CreateLoggingExclusion(parent string, exclusion *logging.LogExclusion) error {
	_, err := u.Config.clientLogging.BillingAccounts.Exclusions.Create(parent, exclusion).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error creating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) ReadLoggingExclusion(id string) (*logging.LogExclusion, error) {
	exclusion, err := u.Config.clientLogging.BillingAccounts.Exclusions.Get(id).Do()

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Error retrieving logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return exclusion, nil
}

func (u *BillingAccountLoggingExclusionUpdater) UpdateLoggingExclusion(id string, exclusion *logging.LogExclusion, updateMask string) error {
	_, err := u.Config.clientLogging.BillingAccounts.Exclusions.Patch(id, exclusion).UpdateMask(updateMask).Do()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Error updating logging exclusion for %s: {{err}}", u.DescribeResource()), err)
	}

	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) DeleteLoggingExclusion(id string) error {
	_, err := u.Config.clientLogging.BillingAccounts.Exclusions.Delete(id).Do()
	if err != nil {
		return errwrap.Wrap(fmt.Errorf("Error deleting logging exclusion for %s.", u.DescribeResource()), err)
	}

	return nil
}

func (u *BillingAccountLoggingExclusionUpdater) GetResourceType() string {
	return u.resourceType
}

func (u *BillingAccountLoggingExclusionUpdater) GetResourceId() string {
	return u.resourceId
}

func (u *BillingAccountLoggingExclusionUpdater) DescribeResource() string {
	return fmt.Sprintf("%q %q", u.resourceType, u.resourceId)
}
