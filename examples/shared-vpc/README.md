# Shared Virtual Private Cloud Networking in Google Cloud

This is a template showcasing the shared VPC feature in Google Cloud.  It features
four projects:
- A host project, which owns a VPC
- Two service projects, each of which owns a VM connected to the VPC
- A fourth project, which owns a VM not connected to the VPC.

It is based on the diagram in the overview at [https://cloud.google.com/vpc/docs/shared-vpc](https://cloud.google.com/vpc/docs/shared-vpc).

Begin by [downloading your credentials from Google Cloud Console](https://www.terraform.io/docs/providers/google/#credentials); the default path for the downloaded file is `~/.gcloud/Terraform.json`.  If you use another path, update the `credentials_file_path` variable.  Ensure that these credentials have Organization-level permissions - this example will create and administer projects.

This example creates projects within an organization - to run it, you will need to have an Organization ID.  To get started using Organizations, read the quickstart [here](https://cloud.google.com/resource-manager/docs/quickstart-organizations).  Since it uses organizations, project-specific credentials won't work, and consequently this example is configured to use [application default credentials](https://developers.google.com/identity/protocols/application-default-credentials).  Ensure that the application default credentials have permission to create and manage projects and Shared VPCs (sometimes called 'XPN').  The example also requires you to specify a billing account, since it does start up a few VMs.

After you run `terraform apply` on this configuration, it will output the IP address of the second service project's VM, which (after it's done starting up) displays a page checking network connectivity to the other two VMs.

Run with a command like:
```
terraform apply \
        -var="region=us-central1" \
        -var="region_zone=us-central1-f" \
        -var="org_id=1234567" \
        -var="billing_account_id=XXXXXXXXXXXX"
```
