subcategory: "Site Verification"
description: |-
  A verification token is used to demonstrate ownership of a website or domain.
---

# google_site_verification_token

A verification token is used to demonstrate ownership of a website or domain.


To get more information about Token, see:

* [API documentation](https://developers.google.com/site-verification/v1)
* How-to Guides
    * [Getting Started](https://developers.google.com/site-verification/v1/getting_started)

<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=siteverification_token_site&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

## Example Usage - Site Verification via Site META Tag

```hcl
data "google_site_verification_token" "example" {
  type                = "SITE"
  identifier          = "https://www.example.com"
  verification_method = "META"
}
```

## Example Usage - Site Verification via DNS TXT Record

```hcl
data "google_site_verification_token" "example" {
  type                = "INET_DOMAIN"
  identifier          = "www.example.com"
  verification_method = "DNS_TXT"
}
```

## Argument Reference

The following arguments are supported:


* `type` -
  (Required)
  The type of resource to be verified, either a domain or a web site.
  Possible values are: `INET_DOMAIN`, `SITE`.

* `identifier` -
  (Required)
  The site identifier. If the type is set to SITE, the identifier is a URL. If the type is
  set to INET_DOMAIN, the identifier is a domain name.

* `verification_method` -
  (Required)
  The verification method for the Site Verification system to use to verify
  this site or domain.
  Possible values are: `ANALYTICS`, `DNS_CNAME`, `DNS_TXT`, `FILE`, `META`, `TAG_MANAGER`.


- - -

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `token` -
  The generated token for use in subsequent verification steps.


## Timeouts

This data source provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `read` - Default is 5 minutes.

## User Project Overrides

This data source supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
