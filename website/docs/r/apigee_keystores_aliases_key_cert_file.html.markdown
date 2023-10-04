---
subcategory: "Apigee"
description: |-
  An alias from a key/certificate pair.
---

# google\_apigee\_keystores\_aliases\_key\_cert\_file

An alias from a key/certificate pair.

To get more information about KeystoresAliasesKeyCertFile, see:

* [API documentation](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.environments.keystores.aliases)
* How-to Guides
    * [Keystores Aliases](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.environments.keystores.aliases)

## Argument Reference

The following arguments are supported:


* `org_id` -
  (Required)
  Organization ID associated with the alias, without organization/ prefix

* `environment` -
  (Required)
  Environment associated with the alias

* `keystore` -
  (Required)
  Keystore Name

* `alias` -
  (Required)
  Alias Name

* `cert` -
  (Required)
  Cert content


- - -


* `key` -
  (Optional)
  Private Key content, omit if uploading to truststore

* `password` -
  (Optional)
  Password for the Private Key if it's encrypted


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}`

* `certs_info` -
  Chain of certificates under this alias.
  Structure is [documented below](#nested_certs_info).

* `type` -
  Optional.Type of Alias


<a name="nested_certs_info"></a>The `certs_info` block contains:

* `cert_info` -
  (Output)
  List of all properties in the object.
  Structure is [documented below](#nested_cert_info).


<a name="nested_cert_info"></a>The `cert_info` block contains:

* `version` -
  (Output)
  X.509 version.

* `subject` -
  (Output)
  X.509 subject.

* `issuer` -
  (Output)
  X.509 issuer.

* `expiry_date` -
  (Output)
  X.509 notAfter validity period in milliseconds since epoch.

* `valid_from` -
  (Output)
  X.509 notBefore validity period in milliseconds since epoch.

* `is_valid` -
  (Output)
  Flag that specifies whether the certificate is valid. 
  Flag is set to Yes if the certificate is valid, No if expired, or Not yet if not yet valid.

* `subject_alternative_names` -
  (Output)
  X.509 subject alternative names (SANs) extension.

* `sig_alg_name` -
  (Output)
  X.509 signatureAlgorithm.

* `public_key` -
  (Output)
  Public key component of the X.509 subject public key info.

* `basic_constraints` -
  (Output)
  X.509 basic constraints extension.

* `serial_number` -
  (Output)
  X.509 serial number.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


KeystoresAliasesKeyCertFile can be imported using any of these accepted formats:

* `organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}`
* `{{org_id}}/{{environment}}/{{keystore}}/{{alias}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import KeystoresAliasesKeyCertFile using one of the formats above. For example:

```tf
import {
  id = "organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}"
  to = google_apigee_keystores_aliases_key_cert_file.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), KeystoresAliasesKeyCertFile can be imported using one of the formats above. For example:

```
$ terraform import google_apigee_keystores_aliases_key_cert_file.default organizations/{{org_id}}/environments/{{environment}}/keystores/{{keystore}}/aliases/{{alias}}
$ terraform import google_apigee_keystores_aliases_key_cert_file.default {{org_id}}/{{environment}}/{{keystore}}/{{alias}}
```
