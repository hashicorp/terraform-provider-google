---
layout: "google"
page_title: "Google: google_admin_group"
sidebar_current: "docs-google-admin-group"
description: |-
  Creates a G Suite's Google Group.
---

# google_admin_group

Creates a G Suite's Google Group. For more information see
[the official help](https://support.google.com/a/topic/25838) and
[API](https://developers.google.com/admin-sdk/directory/v1/reference/groups).

Note: to use this resource with service account, you must
[delegate domain-wide authoritiy](https://developers.google.com/identity/protocols/OAuth2ServiceAccount?hl=ja#delegatingauthority)
in advance.

## Example Usage

```hcl
resource "google_admin_group" "devteam" {
  email       = "devteam@mycompany.com"
  name        = "devteam@mycompany.com"
  description = "Developer team"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The group's unique email address. If your account has multiple domains, select the appropriate domain.

* `name` - (Optional) The group's display name.

* `description` - (Optional) An extended textual description to help users know the purpose of a group. Maximum length is 4,096 characters.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exposed:

* `id` - The group's unique ID.

* `direct_members_count` - The number of users that are direct members of the group. If a group is a member (child) of this group (the parent), members of the child group are not counted in the `direct_members_count` property of the parent group.

* `admin_created` - Whether this group was created by an administrator (`true`) or by a user (`false`).

* `aliases` - List of a group's alias email addresses.  

* `non_editable_aliases` - List of the group's non-editable alias email addresses that are outside of the account's primary domain or subdomains. These are functioning email addresses used by the group.

## Import

The `import` command is currently not implemented yet for this resource.
