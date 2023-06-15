// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package privateca

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourcePrivatecaCertificateAuthority() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourcePrivatecaCertificateAuthority().Schema)
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "location")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "pool")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "certificate_authority_id")

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
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return fmt.Errorf("Error generating user agent: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	err = resourcePrivatecaCertificateAuthorityRead(d, meta)
	if err != nil {
		return err
	}

	// pem_csr is only applicable for SUBORDINATE CertificateAuthorities when their state is AWAITING_USER_ACTIVATION
	if d.Get("type") == "SUBORDINATE" && d.Get("state") == "AWAITING_USER_ACTIVATION" {
		url, err := tpgresource.ReplaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}:fetch")
		if err != nil {
			return err
		}

		billingProject := ""

		project, err := tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("Error fetching project for CertificateAuthority: %s", err)
		}
		billingProject = project

		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("PrivatecaCertificateAuthority %q", d.Id()))
		}
		if err := d.Set("pem_csr", res["pemCsr"]); err != nil {
			return fmt.Errorf("Error fetching CertificateAuthority: %s", err)
		}
	}

	return nil
}
