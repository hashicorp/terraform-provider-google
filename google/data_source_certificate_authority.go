package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivatecaCertificateAuthority() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourcePrivatecaCertificateAuthority().Schema)
	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "location")
	addOptionalFieldsToSchema(dsSchema, "pool")
	addOptionalFieldsToSchema(dsSchema, "certificate_authority_id")

	dsSchema["pem_csr"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Read:   dataSourcePrivatecaCertificateAuthorityRead,
		Schema: dsSchema,
	}
}

func dataSourcePrivatecaCertificateAuthorityRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return fmt.Errorf("Error generating user agent: %s", err)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	err = resourcePrivatecaCertificateAuthorityRead(d, meta)
	if err != nil {
		return err
	}

	// pem_csr is only applicable for SUBORDINATE CertificateAuthorities
	if d.Get("type") == "SUBORDINATE" {
		url, err := replaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}:fetch")
		if err != nil {
			return err
		}

		billingProject := ""

		project, err := getProject(d, config)
		if err != nil {
			return fmt.Errorf("Error fetching project for CertificateAuthority: %s", err)
		}
		billingProject = project

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("PrivatecaCertificateAuthority %q", d.Id()))
		}
		if err := d.Set("pem_csr", res["pemCsr"]); err != nil {
			return fmt.Errorf("Error fetching CertificateAuthority: %s", err)
		}
	}

	return nil
}
