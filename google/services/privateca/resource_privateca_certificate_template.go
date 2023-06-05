// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package privateca

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	privateca "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourcePrivatecaCertificateTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivatecaCertificateTemplateCreate,
		Read:   resourcePrivatecaCertificateTemplateRead,
		Update: resourcePrivatecaCertificateTemplateUpdate,
		Delete: resourcePrivatecaCertificateTemplateDelete,

		Importer: &schema.ResourceImporter{
			State: resourcePrivatecaCertificateTemplateImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource name for this CertificateTemplate in the format `projects/*/locations/*/certificateTemplates/*`.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. A human-readable description of scenarios this template is intended for.",
			},

			"identity_constraints": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Describes constraints on identities that may be appear in Certificates issued using this template. If this is omitted, then this template will not add restrictions on a certificate's identity.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplateIdentityConstraintsSchema(),
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Labels with user-defined metadata.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"passthrough_extensions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Describes the set of X.509 extensions that may appear in a Certificate issued using this CertificateTemplate. If a certificate request sets extensions that don't appear in the passthrough_extensions, those extensions will be dropped. If the issuing CaPool's IssuancePolicy defines baseline_values that don't appear here, the certificate issuance request will fail. If this is omitted, then this template will not add restrictions on a certificate's X.509 extensions. These constraints do not apply to X.509 extensions set in this CertificateTemplate's predefined_values.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplatePassthroughExtensionsSchema(),
			},

			"predefined_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A set of X.509 values that will be applied to all issued certificates that use this template. If the certificate request includes conflicting values for the same properties, they will be overwritten by the values defined here. If the issuing CaPool's IssuancePolicy defines conflicting baseline_values for the same properties, the certificate issuance request will fail.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplatePredefinedValuesSchema(),
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time at which this CertificateTemplate was created.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The time at which this CertificateTemplate was updated.",
			},
		},
	}
}

func PrivatecaCertificateTemplateIdentityConstraintsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_subject_alt_names_passthrough": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Required. If this is true, the SubjectAltNames extension may be copied from a certificate request into the signed certificate. Otherwise, the requested SubjectAltNames will be discarded.",
			},

			"allow_subject_passthrough": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Required. If this is true, the Subject field may be copied from a certificate request into the signed certificate. Otherwise, the requested Subject will be discarded.",
			},

			"cel_expression": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A CEL expression that may be used to validate the resolved X.509 Subject and/or Subject Alternative Name before a certificate is signed. To see the full allowed syntax and some examples, see https://cloud.google.com/certificate-authority-service/docs/using-cel",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplateIdentityConstraintsCelExpressionSchema(),
			},
		},
	}
}

func PrivatecaCertificateTemplateIdentityConstraintsCelExpressionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.",
			},

			"expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Textual representation of an expression in Common Expression Language syntax.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. String indicating the location of the expression for error reporting, e.g. a file name and a position in the file.",
			},

			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Title for the expression, i.e. a short string describing its purpose. This can be used e.g. in UIs which allow to enter the expression.",
			},
		},
	}
}

func PrivatecaCertificateTemplatePassthroughExtensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"additional_extensions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A set of ObjectIds identifying custom X.509 extensions. Will be combined with known_extensions to determine the full set of X.509 extensions.",
				Elem:        PrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsSchema(),
			},

			"known_extensions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A set of named X.509 extensions. Will be combined with additional_extensions to determine the full set of X.509 extensions.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func PrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_id_path": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"additional_extensions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Describes custom X.509 extensions.",
				Elem:        PrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsSchema(),
			},

			"aia_ocsp_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Describes Online Certificate Status Protocol (OCSP) endpoint addresses that appear in the \"Authority Information Access\" extension in the certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"ca_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Describes options in this X509Parameters that are relevant in a CA certificate.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplatePredefinedValuesCaOptionsSchema(),
			},

			"key_usage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Indicates the intended use for keys that correspond to a certificate.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplatePredefinedValuesKeyUsageSchema(),
			},

			"policy_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Describes the X.509 certificate policy object identifiers, per https://tools.ietf.org/html/rfc5280#section-4.2.1.4.",
				Elem:        PrivatecaCertificateTemplatePredefinedValuesPolicyIdsSchema(),
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The OID for this X.509 extension.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSchema(),
			},

			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The value of this X.509 extension.",
			},

			"critical": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Indicates whether or not this extension is critical (i.e., if the client does not know how to handle this extension, the client should consider this to be an error).",
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_id_path": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesCaOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_ca": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Refers to the \"CA\" X.509 extension, which is a boolean value. When this value is missing, the extension will be omitted from the CA certificate.",
			},

			"max_issuer_path_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Optional. Refers to the path length restriction X.509 extension. For a CA certificate, this value describes the depth of subordinate CA certificates that are allowed. If this value is less than 0, the request will fail. If this value is missing, the max path length will be omitted from the CA certificate.",
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesKeyUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"base_key_usage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Describes high-level ways in which a key may be used.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSchema(),
			},

			"extended_key_usage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Detailed scenarios in which a key may be used.",
				MaxItems:    1,
				Elem:        PrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSchema(),
			},

			"unknown_extended_key_usages": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Used to describe extended key usages that are not listed in the KeyUsage.ExtendedKeyUsageOptions message.",
				Elem:        PrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSchema(),
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cert_sign": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used to sign certificates.",
			},

			"content_commitment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used for cryptographic commitments. Note that this may also be referred to as \"non-repudiation\".",
			},

			"crl_sign": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used sign certificate revocation lists.",
			},

			"data_encipherment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used to encipher data.",
			},

			"decipher_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used to decipher only.",
			},

			"digital_signature": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used for digital signatures.",
			},

			"encipher_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used to encipher only.",
			},

			"key_agreement": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used in a key agreement protocol.",
			},

			"key_encipherment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The key may be used to encipher other keys.",
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Corresponds to OID 1.3.6.1.5.5.7.3.2. Officially described as \"TLS WWW client authentication\", though regularly used for non-WWW TLS.",
			},

			"code_signing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Corresponds to OID 1.3.6.1.5.5.7.3.3. Officially described as \"Signing of downloadable executable code client authentication\".",
			},

			"email_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Corresponds to OID 1.3.6.1.5.5.7.3.4. Officially described as \"Email protection\".",
			},

			"ocsp_signing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Corresponds to OID 1.3.6.1.5.5.7.3.9. Officially described as \"Signing OCSP responses\".",
			},

			"server_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Corresponds to OID 1.3.6.1.5.5.7.3.1. Officially described as \"TLS WWW server authentication\", though regularly used for non-WWW TLS.",
			},

			"time_stamping": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Corresponds to OID 1.3.6.1.5.5.7.3.8. Officially described as \"Binding the hash of an object to a time\".",
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_id_path": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func PrivatecaCertificateTemplatePredefinedValuesPolicyIdsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_id_path": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourcePrivatecaCertificateTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &privateca.CertificateTemplate{
		Location:              dcl.String(d.Get("location").(string)),
		Name:                  dcl.String(d.Get("name").(string)),
		Description:           dcl.String(d.Get("description").(string)),
		IdentityConstraints:   expandPrivatecaCertificateTemplateIdentityConstraints(d.Get("identity_constraints")),
		Labels:                tpgresource.CheckStringMap(d.Get("labels")),
		PassthroughExtensions: expandPrivatecaCertificateTemplatePassthroughExtensions(d.Get("passthrough_extensions")),
		PredefinedValues:      expandPrivatecaCertificateTemplatePredefinedValues(d.Get("predefined_values")),
		Project:               dcl.String(project),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLPrivatecaClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyCertificateTemplate(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating CertificateTemplate: %s", err)
	}

	log.Printf("[DEBUG] Finished creating CertificateTemplate %q: %#v", d.Id(), res)

	return resourcePrivatecaCertificateTemplateRead(d, meta)
}

func resourcePrivatecaCertificateTemplateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &privateca.CertificateTemplate{
		Location:              dcl.String(d.Get("location").(string)),
		Name:                  dcl.String(d.Get("name").(string)),
		Description:           dcl.String(d.Get("description").(string)),
		IdentityConstraints:   expandPrivatecaCertificateTemplateIdentityConstraints(d.Get("identity_constraints")),
		Labels:                tpgresource.CheckStringMap(d.Get("labels")),
		PassthroughExtensions: expandPrivatecaCertificateTemplatePassthroughExtensions(d.Get("passthrough_extensions")),
		PredefinedValues:      expandPrivatecaCertificateTemplatePredefinedValues(d.Get("predefined_values")),
		Project:               dcl.String(project),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLPrivatecaClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetCertificateTemplate(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("PrivatecaCertificateTemplate %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("identity_constraints", flattenPrivatecaCertificateTemplateIdentityConstraints(res.IdentityConstraints)); err != nil {
		return fmt.Errorf("error setting identity_constraints in state: %s", err)
	}
	if err = d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("passthrough_extensions", flattenPrivatecaCertificateTemplatePassthroughExtensions(res.PassthroughExtensions)); err != nil {
		return fmt.Errorf("error setting passthrough_extensions in state: %s", err)
	}
	if err = d.Set("predefined_values", flattenPrivatecaCertificateTemplatePredefinedValues(res.PredefinedValues)); err != nil {
		return fmt.Errorf("error setting predefined_values in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourcePrivatecaCertificateTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &privateca.CertificateTemplate{
		Location:              dcl.String(d.Get("location").(string)),
		Name:                  dcl.String(d.Get("name").(string)),
		Description:           dcl.String(d.Get("description").(string)),
		IdentityConstraints:   expandPrivatecaCertificateTemplateIdentityConstraints(d.Get("identity_constraints")),
		Labels:                tpgresource.CheckStringMap(d.Get("labels")),
		PassthroughExtensions: expandPrivatecaCertificateTemplatePassthroughExtensions(d.Get("passthrough_extensions")),
		PredefinedValues:      expandPrivatecaCertificateTemplatePredefinedValues(d.Get("predefined_values")),
		Project:               dcl.String(project),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLPrivatecaClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyCertificateTemplate(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating CertificateTemplate: %s", err)
	}

	log.Printf("[DEBUG] Finished creating CertificateTemplate %q: %#v", d.Id(), res)

	return resourcePrivatecaCertificateTemplateRead(d, meta)
}

func resourcePrivatecaCertificateTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &privateca.CertificateTemplate{
		Location:              dcl.String(d.Get("location").(string)),
		Name:                  dcl.String(d.Get("name").(string)),
		Description:           dcl.String(d.Get("description").(string)),
		IdentityConstraints:   expandPrivatecaCertificateTemplateIdentityConstraints(d.Get("identity_constraints")),
		Labels:                tpgresource.CheckStringMap(d.Get("labels")),
		PassthroughExtensions: expandPrivatecaCertificateTemplatePassthroughExtensions(d.Get("passthrough_extensions")),
		PredefinedValues:      expandPrivatecaCertificateTemplatePredefinedValues(d.Get("predefined_values")),
		Project:               dcl.String(project),
	}

	log.Printf("[DEBUG] Deleting CertificateTemplate %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLPrivatecaClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteCertificateTemplate(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting CertificateTemplate: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting CertificateTemplate %q", d.Id())
	return nil
}

func resourcePrivatecaCertificateTemplateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/certificateTemplates/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandPrivatecaCertificateTemplateIdentityConstraints(o interface{}) *privateca.CertificateTemplateIdentityConstraints {
	if o == nil {
		return privateca.EmptyCertificateTemplateIdentityConstraints
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplateIdentityConstraints
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplateIdentityConstraints{
		AllowSubjectAltNamesPassthrough: dcl.Bool(obj["allow_subject_alt_names_passthrough"].(bool)),
		AllowSubjectPassthrough:         dcl.Bool(obj["allow_subject_passthrough"].(bool)),
		CelExpression:                   expandPrivatecaCertificateTemplateIdentityConstraintsCelExpression(obj["cel_expression"]),
	}
}

func flattenPrivatecaCertificateTemplateIdentityConstraints(obj *privateca.CertificateTemplateIdentityConstraints) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_subject_alt_names_passthrough": obj.AllowSubjectAltNamesPassthrough,
		"allow_subject_passthrough":           obj.AllowSubjectPassthrough,
		"cel_expression":                      flattenPrivatecaCertificateTemplateIdentityConstraintsCelExpression(obj.CelExpression),
	}

	return []interface{}{transformed}

}

func expandPrivatecaCertificateTemplateIdentityConstraintsCelExpression(o interface{}) *privateca.CertificateTemplateIdentityConstraintsCelExpression {
	if o == nil {
		return privateca.EmptyCertificateTemplateIdentityConstraintsCelExpression
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplateIdentityConstraintsCelExpression
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplateIdentityConstraintsCelExpression{
		Description: dcl.String(obj["description"].(string)),
		Expression:  dcl.String(obj["expression"].(string)),
		Location:    dcl.String(obj["location"].(string)),
		Title:       dcl.String(obj["title"].(string)),
	}
}

func flattenPrivatecaCertificateTemplateIdentityConstraintsCelExpression(obj *privateca.CertificateTemplateIdentityConstraintsCelExpression) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"description": obj.Description,
		"expression":  obj.Expression,
		"location":    obj.Location,
		"title":       obj.Title,
	}

	return []interface{}{transformed}

}

func expandPrivatecaCertificateTemplatePassthroughExtensions(o interface{}) *privateca.CertificateTemplatePassthroughExtensions {
	if o == nil {
		return privateca.EmptyCertificateTemplatePassthroughExtensions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplatePassthroughExtensions
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplatePassthroughExtensions{
		AdditionalExtensions: expandPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsArray(obj["additional_extensions"]),
		KnownExtensions:      expandPrivatecaCertificateTemplatePassthroughExtensionsKnownExtensionsArray(obj["known_extensions"]),
	}
}

func flattenPrivatecaCertificateTemplatePassthroughExtensions(obj *privateca.CertificateTemplatePassthroughExtensions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"additional_extensions": flattenPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsArray(obj.AdditionalExtensions),
		"known_extensions":      flattenPrivatecaCertificateTemplatePassthroughExtensionsKnownExtensionsArray(obj.KnownExtensions),
	}

	return []interface{}{transformed}

}
func expandPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsArray(o interface{}) []privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	if o == nil {
		return make([]privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions, 0)
	}

	items := make([]privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions, 0, len(objs))
	for _, item := range objs {
		i := expandPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensions(item)
		items = append(items, *i)
	}

	return items
}

func expandPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensions(o interface{}) *privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions {
	if o == nil {
		return privateca.EmptyCertificateTemplatePassthroughExtensionsAdditionalExtensions
	}

	obj := o.(map[string]interface{})
	return &privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions{
		ObjectIdPath: tpgdclresource.ExpandIntegerArray(obj["object_id_path"]),
	}
}

func flattenPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsArray(objs []privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensions(&item)
		items = append(items, i)
	}

	return items
}

func flattenPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensions(obj *privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"object_id_path": obj.ObjectIdPath,
	}

	return transformed

}

func expandPrivatecaCertificateTemplatePredefinedValues(o interface{}) *privateca.CertificateTemplatePredefinedValues {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValues
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplatePredefinedValues
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValues{
		AdditionalExtensions: expandPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsArray(obj["additional_extensions"]),
		AiaOcspServers:       tpgdclresource.ExpandStringArray(obj["aia_ocsp_servers"]),
		CaOptions:            expandPrivatecaCertificateTemplatePredefinedValuesCaOptions(obj["ca_options"]),
		KeyUsage:             expandPrivatecaCertificateTemplatePredefinedValuesKeyUsage(obj["key_usage"]),
		PolicyIds:            expandPrivatecaCertificateTemplatePredefinedValuesPolicyIdsArray(obj["policy_ids"]),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValues(obj *privateca.CertificateTemplatePredefinedValues) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"additional_extensions": flattenPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsArray(obj.AdditionalExtensions),
		"aia_ocsp_servers":      obj.AiaOcspServers,
		"ca_options":            flattenPrivatecaCertificateTemplatePredefinedValuesCaOptions(obj.CaOptions),
		"key_usage":             flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsage(obj.KeyUsage),
		"policy_ids":            flattenPrivatecaCertificateTemplatePredefinedValuesPolicyIdsArray(obj.PolicyIds),
	}

	return []interface{}{transformed}

}
func expandPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsArray(o interface{}) []privateca.CertificateTemplatePredefinedValuesAdditionalExtensions {
	if o == nil {
		return make([]privateca.CertificateTemplatePredefinedValuesAdditionalExtensions, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]privateca.CertificateTemplatePredefinedValuesAdditionalExtensions, 0)
	}

	items := make([]privateca.CertificateTemplatePredefinedValuesAdditionalExtensions, 0, len(objs))
	for _, item := range objs {
		i := expandPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensions(item)
		items = append(items, *i)
	}

	return items
}

func expandPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensions(o interface{}) *privateca.CertificateTemplatePredefinedValuesAdditionalExtensions {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesAdditionalExtensions
	}

	obj := o.(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesAdditionalExtensions{
		ObjectId: expandPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(obj["object_id"]),
		Value:    dcl.String(obj["value"].(string)),
		Critical: dcl.Bool(obj["critical"].(bool)),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsArray(objs []privateca.CertificateTemplatePredefinedValuesAdditionalExtensions) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensions(&item)
		items = append(items, i)
	}

	return items
}

func flattenPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensions(obj *privateca.CertificateTemplatePredefinedValuesAdditionalExtensions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"object_id": flattenPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(obj.ObjectId),
		"value":     obj.Value,
		"critical":  obj.Critical,
	}

	return transformed

}

func expandPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(o interface{}) *privateca.CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{
		ObjectIdPath: tpgdclresource.ExpandIntegerArray(obj["object_id_path"]),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(obj *privateca.CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"object_id_path": obj.ObjectIdPath,
	}

	return []interface{}{transformed}

}

func expandPrivatecaCertificateTemplatePredefinedValuesCaOptions(o interface{}) *privateca.CertificateTemplatePredefinedValuesCaOptions {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesCaOptions
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesCaOptions
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesCaOptions{
		IsCa:                dcl.Bool(obj["is_ca"].(bool)),
		MaxIssuerPathLength: dcl.Int64(int64(obj["max_issuer_path_length"].(int))),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesCaOptions(obj *privateca.CertificateTemplatePredefinedValuesCaOptions) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"is_ca":                  obj.IsCa,
		"max_issuer_path_length": obj.MaxIssuerPathLength,
	}

	return []interface{}{transformed}

}

func expandPrivatecaCertificateTemplatePredefinedValuesKeyUsage(o interface{}) *privateca.CertificateTemplatePredefinedValuesKeyUsage {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesKeyUsage
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesKeyUsage
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesKeyUsage{
		BaseKeyUsage:             expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(obj["base_key_usage"]),
		ExtendedKeyUsage:         expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(obj["extended_key_usage"]),
		UnknownExtendedKeyUsages: expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesArray(obj["unknown_extended_key_usages"]),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsage(obj *privateca.CertificateTemplatePredefinedValuesKeyUsage) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"base_key_usage":              flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(obj.BaseKeyUsage),
		"extended_key_usage":          flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(obj.ExtendedKeyUsage),
		"unknown_extended_key_usages": flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesArray(obj.UnknownExtendedKeyUsages),
	}

	return []interface{}{transformed}

}

func expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(o interface{}) *privateca.CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{
		CertSign:          dcl.Bool(obj["cert_sign"].(bool)),
		ContentCommitment: dcl.Bool(obj["content_commitment"].(bool)),
		CrlSign:           dcl.Bool(obj["crl_sign"].(bool)),
		DataEncipherment:  dcl.Bool(obj["data_encipherment"].(bool)),
		DecipherOnly:      dcl.Bool(obj["decipher_only"].(bool)),
		DigitalSignature:  dcl.Bool(obj["digital_signature"].(bool)),
		EncipherOnly:      dcl.Bool(obj["encipher_only"].(bool)),
		KeyAgreement:      dcl.Bool(obj["key_agreement"].(bool)),
		KeyEncipherment:   dcl.Bool(obj["key_encipherment"].(bool)),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(obj *privateca.CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cert_sign":          obj.CertSign,
		"content_commitment": obj.ContentCommitment,
		"crl_sign":           obj.CrlSign,
		"data_encipherment":  obj.DataEncipherment,
		"decipher_only":      obj.DecipherOnly,
		"digital_signature":  obj.DigitalSignature,
		"encipher_only":      obj.EncipherOnly,
		"key_agreement":      obj.KeyAgreement,
		"key_encipherment":   obj.KeyEncipherment,
	}

	return []interface{}{transformed}

}

func expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(o interface{}) *privateca.CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
	}
	obj := objArr[0].(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{
		ClientAuth:      dcl.Bool(obj["client_auth"].(bool)),
		CodeSigning:     dcl.Bool(obj["code_signing"].(bool)),
		EmailProtection: dcl.Bool(obj["email_protection"].(bool)),
		OcspSigning:     dcl.Bool(obj["ocsp_signing"].(bool)),
		ServerAuth:      dcl.Bool(obj["server_auth"].(bool)),
		TimeStamping:    dcl.Bool(obj["time_stamping"].(bool)),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(obj *privateca.CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"client_auth":      obj.ClientAuth,
		"code_signing":     obj.CodeSigning,
		"email_protection": obj.EmailProtection,
		"ocsp_signing":     obj.OcspSigning,
		"server_auth":      obj.ServerAuth,
		"time_stamping":    obj.TimeStamping,
	}

	return []interface{}{transformed}

}
func expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesArray(o interface{}) []privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	if o == nil {
		return make([]privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, 0)
	}

	items := make([]privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages, 0, len(objs))
	for _, item := range objs {
		i := expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(item)
		items = append(items, *i)
	}

	return items
}

func expandPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(o interface{}) *privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages
	}

	obj := o.(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{
		ObjectIdPath: tpgdclresource.ExpandIntegerArray(obj["object_id_path"]),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesArray(objs []privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(&item)
		items = append(items, i)
	}

	return items
}

func flattenPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(obj *privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"object_id_path": obj.ObjectIdPath,
	}

	return transformed

}
func expandPrivatecaCertificateTemplatePredefinedValuesPolicyIdsArray(o interface{}) []privateca.CertificateTemplatePredefinedValuesPolicyIds {
	if o == nil {
		return make([]privateca.CertificateTemplatePredefinedValuesPolicyIds, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]privateca.CertificateTemplatePredefinedValuesPolicyIds, 0)
	}

	items := make([]privateca.CertificateTemplatePredefinedValuesPolicyIds, 0, len(objs))
	for _, item := range objs {
		i := expandPrivatecaCertificateTemplatePredefinedValuesPolicyIds(item)
		items = append(items, *i)
	}

	return items
}

func expandPrivatecaCertificateTemplatePredefinedValuesPolicyIds(o interface{}) *privateca.CertificateTemplatePredefinedValuesPolicyIds {
	if o == nil {
		return privateca.EmptyCertificateTemplatePredefinedValuesPolicyIds
	}

	obj := o.(map[string]interface{})
	return &privateca.CertificateTemplatePredefinedValuesPolicyIds{
		ObjectIdPath: tpgdclresource.ExpandIntegerArray(obj["object_id_path"]),
	}
}

func flattenPrivatecaCertificateTemplatePredefinedValuesPolicyIdsArray(objs []privateca.CertificateTemplatePredefinedValuesPolicyIds) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenPrivatecaCertificateTemplatePredefinedValuesPolicyIds(&item)
		items = append(items, i)
	}

	return items
}

func flattenPrivatecaCertificateTemplatePredefinedValuesPolicyIds(obj *privateca.CertificateTemplatePredefinedValuesPolicyIds) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"object_id_path": obj.ObjectIdPath,
	}

	return transformed

}
func flattenPrivatecaCertificateTemplatePassthroughExtensionsKnownExtensionsArray(obj []privateca.CertificateTemplatePassthroughExtensionsKnownExtensionsEnum) interface{} {
	if obj == nil {
		return nil
	}
	items := []string{}
	for _, item := range obj {
		items = append(items, string(item))
	}
	return items
}
func expandPrivatecaCertificateTemplatePassthroughExtensionsKnownExtensionsArray(o interface{}) []privateca.CertificateTemplatePassthroughExtensionsKnownExtensionsEnum {
	objs := o.([]interface{})
	items := make([]privateca.CertificateTemplatePassthroughExtensionsKnownExtensionsEnum, 0, len(objs))
	for _, item := range objs {
		i := privateca.CertificateTemplatePassthroughExtensionsKnownExtensionsEnumRef(item.(string))
		items = append(items, *i)
	}
	return items
}
