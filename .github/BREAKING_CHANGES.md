---
# ----------------------------------------------------------------------------
#
#   ***    AUTO GENERATED CODE     ***    Type: breaking-change-detector  ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     Changes will need to be made to the breaking-change-detector within Magic Modules instead of here.
#
# ----------------------------------------------------------------------------
---

# Breaking Changes and Provider Development

## Provider Versioning
As a provider is developed; resources are added, old resources are updated, and bugs are fixed.
These changes are [bundled together as a release](https://github.com/hashicorp/terraform-provider-google/releases/tag/v4.32.0).
Releases are numerically defined with a version number in the form of `MAJOR.MINOR.PATCH`.
Patch indicates bug fixes, minor represents new features, and major represents significant changes
which would be breaking to the customer if committed. Once a release is published the provider binary is copied to
[Hashicorp's provider registry](https://registry.terraform.io/browse/providers).

## Customer Trust
Terraform authors can write modular configurations, aptly named modules. These are shared within organizations and
[online](https://registry.terraform.io/browse/modules). Terraform configurations can specify [provider requirements](https://www.terraform.io/language/providers/requirements)
including a [version constraint field](https://www.terraform.io/language/providers/requirements#version-constraints).
The configuration will then [tie these version constraints](https://www.terraform.io/language/expressions/version-constraints)
to an approximate minor or exact full version. Maintaining trust and consistency on every `MINOR` or `MAJOR` version upgrade is critical.

If breaking changes are allowed within `MINOR` versions, trust in the provider will be eroded and module creators will
not have confidence in provider stability. This diminished trust will eventually lead to customers investing or deploying less to GCP.

## Breaking Changes

Now that we understand what defines a breaking change and that we don't want them.
What exactly constitutes a breaking change? Bellow we'll
go into the four categories and rules therein.


### Resource Inventory Level Breakages

* Resource/datasource naming conventions and entry differences.

<h4 id="resource-map-resource-removal-or-rename"> Removing or Renaming an Resource </h4>
In terraform resources should be retained whenever possible. A removable of an resource will result in a configuration breakage wherever a dependency on that resource exists. Renaming or Removing a resources are functionally equivalent in terms of configuration breakages.

### Resource Level Breakages

* Individual resource breakages like field entry removals or behavior within a resource.

<h4 id="resource-schema-field-removal-or-rename"> Removing or Renaming an field </h4>
In terraform fields should be retained whenever possible. A removable of an field will result in a configuration breakage wherever a dependency on that field exists. Renaming or Removing a field are functionally equivalent in terms of configuration breakages.

### Field Level Breakages

* Field level conventions like attribute changes and naming conventions.

<h4 id="field-optional-to-required"> Field becoming Required Field </h4>
A field cannot become required as existing terraform modules may not have this field defined. Thus breaking their modules in sequential plan or applies.

