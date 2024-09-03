// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package certificatemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/certificatemanager/v1"
)

func DataSourceGoogleCertificateManagerCertificates() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCertificateManagerCertificate().Schema)
	tpgresource.DeleteFieldsFromSchema(dsSchema, "self_managed")

	return &schema.Resource{
		Read: dataSourceGoogleCertificateManagerCertificatesRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "global",
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
		},
	}
}

func dataSourceGoogleCertificateManagerCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("error fetching project for certificate: %s", err)
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return fmt.Errorf("error fetching region for certificate: %s", err)
	}

	filter := d.Get("filter").(string)

	certificates := make([]map[string]interface{}, 0)
	certificatesList, err := config.NewCertificateManagerClient(userAgent).Projects.Locations.Certificates.List(fmt.Sprintf("projects/%s/locations/%s", project, region)).Filter(filter).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Certificates : %s %s", project, region))
	}

	for _, certificate := range certificatesList.Certificates {
		if certificate != nil {
			certificates = append(certificates, map[string]interface{}{
				"name":         certificate.Name,
				"description":  certificate.Description,
				"labels":       certificate.Labels,
				"location":     region,
				"managed":      flattenCertificateManaged(certificate.Managed),
				"san_dnsnames": certificate.SanDnsnames,
				"scope":        certificate.Scope,
			})
		}
	}

	if err := d.Set("certificates", certificates); err != nil {
		return fmt.Errorf("error setting certificates: %s", err)
	}

	d.SetId(fmt.Sprintf(
		"projects/%s/locations/%s/certificates",
		project,
		region,
	))

	return nil
}

func flattenCertificateManaged(v *certificatemanager.ManagedCertificate) interface{} {
	if v == nil {
		return nil
	}

	output := make(map[string]interface{})

	output["authorization_attempt_info"] = flattenCertificateManagedAuthorizationAttemptInfo(v.AuthorizationAttemptInfo)
	output["dns_authorizations"] = v.DnsAuthorizations
	output["domains"] = v.Domains
	output["issuance_config"] = v.IssuanceConfig
	output["state"] = v.State
	output["provisioning_issue"] = flattenCertificateManagedProvisioningIssue(v.ProvisioningIssue)

	return []interface{}{output}
}

func flattenCertificateManagedAuthorizationAttemptInfo(v []*certificatemanager.AuthorizationAttemptInfo) interface{} {
	if v == nil {
		return nil
	}

	output := make([]interface{}, 0, len(v))

	for _, authorizationAttemptInfo := range v {
		output = append(output, map[string]interface{}{
			"details":        authorizationAttemptInfo.Details,
			"domain":         authorizationAttemptInfo.Domain,
			"failure_reason": authorizationAttemptInfo.FailureReason,
			"state":          authorizationAttemptInfo.State,
		})
	}

	return output
}

func flattenCertificateManagedProvisioningIssue(v *certificatemanager.ProvisioningIssue) interface{} {
	if v == nil {
		return nil
	}

	output := make(map[string]interface{})

	output["details"] = v.Details
	output["reason"] = v.Reason

	return []interface{}{output}
}
