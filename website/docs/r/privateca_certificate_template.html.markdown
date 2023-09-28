---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     and is based on the DCL (https:#github.com/GoogleCloudPlatform/declarative-resource-client-library).
#     Changes will need to be made to the DCL or Magic Modules instead of here.
#
#     We are not currently able to accept contributions to this file. If changes
#     are required, please file an issue at https:#github.com/hashicorp/terraform-provider-google/issues/new/choose
#
# ----------------------------------------------------------------------------
subcategory: "Certificate Authority Service"
description: |-
  Certificate Authority Service provides reusable and parameterized templates that you can use for common certificate issuance scenarios. A certificate template represents a relatively static and well-defined certificate issuance schema within an organization.  A certificate template can essentially become a full-fledged vertical certificate issuance framework.
---

# google_privateca_certificate_template

Certificate Authority Service provides reusable and parameterized templates that you can use for common certificate issuance scenarios. A certificate template represents a relatively static and well-defined certificate issuance schema within an organization.  A certificate template can essentially become a full-fledged vertical certificate issuance framework.

For more information, see:
* [Understanding Certificate Templates](https://cloud.google.com/certificate-authority-service/docs/certificate-template)
* [Common configurations and Certificate Profiles](https://cloud.google.com/certificate-authority-service/docs/certificate-profile)
## Example Usage - basic_certificate_template
An example of a basic privateca certificate template
```hcl
resource "google_privateca_certificate_template" "primary" {
  location    = "us-west1"
  name        = "template"
  description = "An updated sample certificate template"

  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true

    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }

    known_extensions = ["EXTENDED_KEY_USAGE"]
  }

  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }

      value    = "c3RyaW5nCg=="
      critical = true
    }

    aia_ocsp_servers = ["string"]

    ca_options {
      is_ca                  = false
      max_issuer_path_length = 6
    }

    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }

      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }

      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }

    policy_ids {
      object_id_path = [1, 6]
    }
  }

  project = "my-project-name"

  labels = {
    label-two = "value-two"
  }
}


```

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  The resource name for this CertificateTemplate in the format `projects/*/locations/*/certificateTemplates/*`.
  


The `object_id` block supports:
    
* `object_id_path` -
  (Required)
  Required. The parts of an OID path. The most significant parts of the path come first.
    
- - -

* `description` -
  (Optional)
  Optional. A human-readable description of scenarios this template is intended for.
  
* `identity_constraints` -
  (Optional)
  Optional. Describes constraints on identities that may be appear in Certificates issued using this template. If this is omitted, then this template will not add restrictions on a certificate's identity.
  
* `labels` -
  (Optional)
  Optional. Labels with user-defined metadata.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration. Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `passthrough_extensions` -
  (Optional)
  Optional. Describes the set of X.509 extensions that may appear in a Certificate issued using this CertificateTemplate. If a certificate request sets extensions that don't appear in the passthrough_extensions, those extensions will be dropped. If the issuing CaPool's IssuancePolicy defines baseline_values that don't appear here, the certificate issuance request will fail. If this is omitted, then this template will not add restrictions on a certificate's X.509 extensions. These constraints do not apply to X.509 extensions set in this CertificateTemplate's predefined_values.
  
* `predefined_values` -
  (Optional)
  Optional. A set of X.509 values that will be applied to all issued certificates that use this template. If the certificate request includes conflicting values for the same properties, they will be overwritten by the values defined here. If the issuing CaPool's IssuancePolicy defines conflicting baseline_values for the same properties, the certificate issuance request will fail.
  
* `project` -
  (Optional)
  The project for the resource
  


The `identity_constraints` block supports:
    
* `allow_subject_alt_names_passthrough` -
  (Required)
  Required. If this is true, the SubjectAltNames extension may be copied from a certificate request into the signed certificate. Otherwise, the requested SubjectAltNames will be discarded.
    
* `allow_subject_passthrough` -
  (Required)
  Required. If this is true, the Subject field may be copied from a certificate request into the signed certificate. Otherwise, the requested Subject will be discarded.
    
* `cel_expression` -
  (Optional)
  Optional. A CEL expression that may be used to validate the resolved X.509 Subject and/or Subject Alternative Name before a certificate is signed. To see the full allowed syntax and some examples, see https://cloud.google.com/certificate-authority-service/docs/using-cel
    
The `cel_expression` block supports:
    
* `description` -
  (Optional)
  Optional. Description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.
    
* `expression` -
  (Optional)
  Textual representation of an expression in Common Expression Language syntax.
    
* `location` -
  (Optional)
  Optional. String indicating the location of the expression for error reporting, e.g. a file name and a position in the file.
    
* `title` -
  (Optional)
  Optional. Title for the expression, i.e. a short string describing its purpose. This can be used e.g. in UIs which allow to enter the expression.
    
The `passthrough_extensions` block supports:
    
* `additional_extensions` -
  (Optional)
  Optional. A set of ObjectIds identifying custom X.509 extensions. Will be combined with known_extensions to determine the full set of X.509 extensions.
    
* `known_extensions` -
  (Optional)
  Optional. A set of named X.509 extensions. Will be combined with additional_extensions to determine the full set of X.509 extensions.
    
The `additional_extensions` block supports:
    
* `object_id_path` -
  (Required)
  Required. The parts of an OID path. The most significant parts of the path come first.
    
The `predefined_values` block supports:
    
* `additional_extensions` -
  (Optional)
  Optional. Describes custom X.509 extensions.
    
* `aia_ocsp_servers` -
  (Optional)
  Optional. Describes Online Certificate Status Protocol (OCSP) endpoint addresses that appear in the "Authority Information Access" extension in the certificate.
    
* `ca_options` -
  (Optional)
  Optional. Describes options in this X509Parameters that are relevant in a CA certificate.
    
* `key_usage` -
  (Optional)
  Optional. Indicates the intended use for keys that correspond to a certificate.
    
* `policy_ids` -
  (Optional)
  Optional. Describes the X.509 certificate policy object identifiers, per https://tools.ietf.org/html/rfc5280#section-4.2.1.4.
    
The `additional_extensions` block supports:
    
* `critical` -
  (Optional)
  Optional. Indicates whether or not this extension is critical (i.e., if the client does not know how to handle this extension, the client should consider this to be an error).
    
* `object_id` -
  (Required)
  Required. The OID for this X.509 extension.
    
* `value` -
  (Required)
  Required. The value of this X.509 extension.
    
The `ca_options` block supports:
    
* `is_ca` -
  (Optional)
  Optional. Refers to the "CA" X.509 extension, which is a boolean value. When this value is missing, the extension will be omitted from the CA certificate.
    
* `max_issuer_path_length` -
  (Optional)
  Optional. Refers to the path length restriction X.509 extension. For a CA certificate, this value describes the depth of subordinate CA certificates that are allowed. If this value is less than 0, the request will fail. If this value is missing, the max path length will be omitted from the CA certificate.
    
The `key_usage` block supports:
    
* `base_key_usage` -
  (Optional)
  Describes high-level ways in which a key may be used.
    
* `extended_key_usage` -
  (Optional)
  Detailed scenarios in which a key may be used.
    
* `unknown_extended_key_usages` -
  (Optional)
  Used to describe extended key usages that are not listed in the KeyUsage.ExtendedKeyUsageOptions message.
    
The `base_key_usage` block supports:
    
* `cert_sign` -
  (Optional)
  The key may be used to sign certificates.
    
* `content_commitment` -
  (Optional)
  The key may be used for cryptographic commitments. Note that this may also be referred to as "non-repudiation".
    
* `crl_sign` -
  (Optional)
  The key may be used sign certificate revocation lists.
    
* `data_encipherment` -
  (Optional)
  The key may be used to encipher data.
    
* `decipher_only` -
  (Optional)
  The key may be used to decipher only.
    
* `digital_signature` -
  (Optional)
  The key may be used for digital signatures.
    
* `encipher_only` -
  (Optional)
  The key may be used to encipher only.
    
* `key_agreement` -
  (Optional)
  The key may be used in a key agreement protocol.
    
* `key_encipherment` -
  (Optional)
  The key may be used to encipher other keys.
    
The `extended_key_usage` block supports:
    
* `client_auth` -
  (Optional)
  Corresponds to OID 1.3.6.1.5.5.7.3.2. Officially described as "TLS WWW client authentication", though regularly used for non-WWW TLS.
    
* `code_signing` -
  (Optional)
  Corresponds to OID 1.3.6.1.5.5.7.3.3. Officially described as "Signing of downloadable executable code client authentication".
    
* `email_protection` -
  (Optional)
  Corresponds to OID 1.3.6.1.5.5.7.3.4. Officially described as "Email protection".
    
* `ocsp_signing` -
  (Optional)
  Corresponds to OID 1.3.6.1.5.5.7.3.9. Officially described as "Signing OCSP responses".
    
* `server_auth` -
  (Optional)
  Corresponds to OID 1.3.6.1.5.5.7.3.1. Officially described as "TLS WWW server authentication", though regularly used for non-WWW TLS.
    
* `time_stamping` -
  (Optional)
  Corresponds to OID 1.3.6.1.5.5.7.3.8. Officially described as "Binding the hash of an object to a time".
    
The `unknown_extended_key_usages` block supports:
    
* `object_id_path` -
  (Required)
  Required. The parts of an OID path. The most significant parts of the path come first.
    
The `policy_ids` block supports:
    
* `object_id_path` -
  (Required)
  Required. The parts of an OID path. The most significant parts of the path come first.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}`

* `create_time` -
  Output only. The time at which this CertificateTemplate was created.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
  
* `update_time` -
  Output only. The time at which this CertificateTemplate was updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

CertificateTemplate can be imported using any of these accepted formats:

```
$ terraform import google_privateca_certificate_template.default projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}
$ terraform import google_privateca_certificate_template.default {{project}}/{{location}}/{{name}}
$ terraform import google_privateca_certificate_template.default {{location}}/{{name}}
```



