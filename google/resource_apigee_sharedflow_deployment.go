package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceApigeeSharedFlowDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigeeSharedflowDeploymentCreate,
		Read:   resourceApigeeSharedflowDeploymentRead,
		Delete: resourceApigeeSharedflowDeploymentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApigeeSharedflowDeploymentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource ID of the environment.`,
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The Apigee Organization associated with the Apigee instance`,
			},
			"revision": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Revision of the Sharedflow to be deployed.`,
			},
			"service_account": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: `The service account represents the identity of the deployed proxy, and determines what permissions it has. The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.`,
			},
			"sharedflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Id of the Sharedflow to be deployed.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceApigeeSharedflowDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments?override=true&serviceAccount={{service_account}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new SharedflowDeployment at %s", url)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, nil, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating SharedflowDeployment: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating SharedflowDeployment %q: %#v", d.Id(), res)

	return resourceApigeeSharedflowDeploymentRead(d, meta)
}

func resourceApigeeSharedflowDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	log.Printf("[DEBUG] Reading SharedflowDeployment at %s", url)

	res, err := SendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ApigeeSharedflowDeployment %q", d.Id()))
	}
	log.Printf("[DEBUG] ApigeeSharedflowDeployment deployStartTime %s", res["deployStartTime"])

	return nil
}

func resourceApigeeSharedflowDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := replaceVars(d, config, "{{ApigeeBasePath}}organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting SharedflowDeployment %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "SharedflowDeployment")
	}

	log.Printf("[DEBUG] Finished deleting SharedflowDeployment %q: %#v", d.Id(), res)
	return nil
}

func resourceApigeeSharedflowDeploymentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"organizations/(?P<org_id>[^/]+)/environments/(?P<environment>[^/]+)/sharedflows/(?P<sharedflow_id>[^/]+)/revisions/(?P<revision>[^/]+)",
		"(?P<org_id>[^/]+)/(?P<environment>[^/]+)/(?P<sharedflow_id>[^/]+)/(?P<revision>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenApigeeSharedflowDeploymentOrgId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeSharedflowDeploymentEnvironment(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeSharedflowDeploymentSharedflowId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeSharedflowDeploymentRevision(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeSharedflowDeploymentServiceAccount(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}
