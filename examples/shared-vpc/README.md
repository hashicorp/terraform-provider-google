# Shared Virtual Private Cloud Networking in Google Cloud

This is a template showcasing the shared VPC feature in Google Cloud.  It features
four projects:
- A host project, which owns a VPC
- Two service projects, each of which owns a VM connected to the VPC
- A fourth project, which owns a VM not connected to the VPC.

It is based on the diagram in the overview at [https://cloud.google.com/vpc/docs/shared-vpc](https://cloud.google.com/vpc/docs/shared-vpc).

This example creates projects within an organization - to run it, you will need to have an Organization ID.  To get started using Organizations, read the quickstart [here](https://cloud.google.com/resource-manager/docs/quickstart-organizations).  Since it uses organizations, project-specific credentials won't work, and consequently this example is configured to use [application default credentials](https://developers.google.com/identity/protocols/application-default-credentials).  Ensure that the application default credentials have permission to create and manage projects and Shared VPCs (sometimes called 'XPN').

Since projects require globally unique names, you will need to provide a `project_base_id` variable.  Since project names can't be too long, make sure it's short - if it is too long, Terraform won't be able to create the project.

The example allows you to specify a billing account, if you want to ensure that the charges (which are minor) are grouped together.

After you run `terraform apply` on this configuration, it will output the IP address of the second service project's VM, which displays a page checking network connectivity to the other two VMs.

Run with a command like:
```
terraform apply \
        -var="region=us-central1" \
        -var="region_zone=us-central1-f" \
        -var="org_id=1234567" \
        -var="project_base_id=unique_string"
```

When you run `terraform destroy` (with the same variables), the projects created by `terraform apply` will not be deleted.  This is for convenience - because deleted projects cannot be recreated with the same name, you would need to create new projects with new globally-unique names for each `terraform apply` / `terraform destroy` pair.  However, you will need to re-import the projects if you wish to rerun `apply` after `destroy` - use commands like these, replacing `$PROJECT_BASE_ID` with the `project_base_id` you used above:

```
terraform import google_project.host_project host-project-$PROJECT_BASE_ID
terraform import google_project.service_project_1 service-project-$PROJECT_BASE_ID-1
terraform import google_project.service_project_2 service-project-$PROJECT_BASE_ID-2
terraform import google_project.standalone_project standalone-$PROJECT_BASE_ID
```
