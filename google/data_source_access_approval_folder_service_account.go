package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAccessApprovalFolderServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccessApprovalFolderServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"folder_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccessApprovalFolderServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{AccessApprovalBasePath}}folders/{{folder_id}}/serviceAccount")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("AccessApprovalFolderServiceAccount %q", d.Id()))
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("account_email", res["accountEmail"]); err != nil {
		return fmt.Errorf("Error setting account_email: %s", err)
	}
	d.SetId(res["name"].(string))

	return nil
}
