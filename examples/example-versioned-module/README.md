# Terraform Google Cloud Platform Provider - Example Versioned Module

The `google` and `google-beta` split requires users to explicitly set
the version of the Google provider for Terraform that they are using;
see the [Google Provider Versions](https://www.terraform.io/docs/providers/google/provider_versions.html)
page for more details.

This has complicated module creation as the schema between `google`
and `google-beta` often differs; specifying a Beta feature with
the `google` provider will give an error. This example module
demonstrates how to create a "versioned" module that detects the
necessary version for a resource based on the fields specified.

This example only solves the simple case of a single beta field
in a single resource, but should give module developers the right
ideas on how to develop more complex modules intermixing `google`
and `google-beta`.
