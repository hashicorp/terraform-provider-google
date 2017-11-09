---
layout: "google"
page_title: "Google: google_project"
sidebar_current: "docs-google-datasource-project"
description: |-
  Provides the Google Project details based on a name
---

# google\_project

Provides access to the latest available Google Project details based a given name.
See more about [project details](https://cloud.google.com/resource-manager/docs/cloud-platform-resource-hierarchy#projects) in the upstream docs.

```
data "google_project" "foo" {
  name = "foobar"
}

resource "google_project_services" "project" {
  project  = "${data.google_project.foo.project_id}"
  services = ["iam.googleapis.com", "cloudresourcemanager.googleapis.com"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The display name of the project.

## Attributes Reference

The following attribute is exported:

In addition to the arguments listed above, the following computed attributes are
exported:

* `project_id` - The project ID.

* `number` - The numeric identifier of the project.
