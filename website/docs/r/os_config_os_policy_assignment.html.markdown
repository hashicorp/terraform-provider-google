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
subcategory: "OS Config"
page_title: "Google: google_os_config_os_policy_assignment"
description: |-
  Represents an OSPolicyAssignment resource.
---

# google_os_config_os_policy_assignment

Represents an OSPolicyAssignment resource.

## Example Usage - fixed_os_policy_assignment
An example of an osconfig os policy assignment with fixed rollout disruption budget
```hcl
resource "google_os_config_os_policy_assignment" "primary" {
  instance_filter {
    all = false

    exclusion_labels {
      labels = {
        label-two = "value-two"
      }
    }

    inclusion_labels {
      labels = {
        label-one = "value-one"
      }
    }

    inventories {
      os_short_name = "centos"
      os_version    = "8.*"
    }
  }

  location = "us-west1-a"
  name     = "assignment"

  os_policies {
    id   = "policy"
    mode = "VALIDATION"

    resource_groups {
      resources {
        id = "apt"

        pkg {
          desired_state = "INSTALLED"

          apt {
            name = "bazel"
          }
        }
      }

      inventory_filters {
        os_short_name = "centos"
        os_version    = "8.*"
      }
    }

    allow_no_resource_group_match = false
    description                   = "A test os policy"
  }

  rollout {
    disruption_budget {
      fixed = 1
    }

    min_wait_duration = "3.5s"
  }

  description = "A test os policy assignment"
  project     = "my-project-name"
}


```

## Argument Reference

The following arguments are supported:

* `instance_filter` -
  (Required)
  Required. Filter to select VMs.
  
* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  Resource name.
  
* `os_policies` -
  (Required)
  Required. List of OS policies to be applied to the VMs.
  
* `rollout` -
  (Required)
  Required. Rollout to deploy the OS policy assignment. A rollout is triggered in the following situations: 1) OSPolicyAssignment is created. 2) OSPolicyAssignment is updated and the update contains changes to one of the following fields: - instance_filter - os_policies 3) OSPolicyAssignment is deleted.
  


The `instance_filter` block supports:
    
* `all` -
  (Optional)
  Target all VMs in the project. If true, no other criteria is permitted.
    
* `exclusion_labels` -
  (Optional)
  List of label sets used for VM exclusion. If the list has more than one label set, the VM is excluded if any of the label sets are applicable for the VM.
    
* `inclusion_labels` -
  (Optional)
  List of label sets used for VM inclusion. If the list has more than one `LabelSet`, the VM is included if any of the label sets are applicable for the VM.
    
* `inventories` -
  (Optional)
  List of inventories to select VMs. A VM is selected if its inventory data matches at least one of the following inventories.
    
The `os_policies` block supports:
    
* `allow_no_resource_group_match` -
  (Optional)
  This flag determines the OS policy compliance status when none of the resource groups within the policy are applicable for a VM. Set this value to `true` if the policy needs to be reported as compliant even if the policy has nothing to validate or enforce.
    
* `description` -
  (Optional)
  Policy description. Length of the description is limited to 1024 characters.
    
* `id` -
  (Required)
  Required. The id of the OS policy with the following restrictions: * Must contain only lowercase letters, numbers, and hyphens. * Must start with a letter. * Must be between 1-63 characters. * Must end with a number or a letter. * Must be unique within the assignment.
    
* `mode` -
  (Required)
  Required. Policy mode Possible values: MODE_UNSPECIFIED, VALIDATION, ENFORCEMENT
    
* `resource_groups` -
  (Required)
  Required. List of resource groups for the policy. For a particular VM, resource groups are evaluated in the order specified and the first resource group that is applicable is selected and the rest are ignored. If none of the resource groups are applicable for a VM, the VM is considered to be non-compliant w.r.t this policy. This behavior can be toggled by the flag `allow_no_resource_group_match`
    
The `resource_groups` block supports:
    
* `inventory_filters` -
  (Optional)
  List of inventory filters for the resource group. The resources in this resource group are applied to the target VM if it satisfies at least one of the following inventory filters. For example, to apply this resource group to VMs running either `RHEL` or `CentOS` operating systems, specify 2 items for the list with following values: inventory_filters[0].os_short_name='rhel' and inventory_filters[1].os_short_name='centos' If the list is empty, this resource group will be applied to the target VM unconditionally.
    
* `resources` -
  (Required)
  Required. List of resources configured for this resource group. The resources are executed in the exact order specified here.
    
The `resources` block supports:
    
* `exec` -
  (Optional)
  Exec resource
    
* `file` -
  (Optional)
  File resource
    
* `id` -
  (Required)
  Required. The id of the resource with the following restrictions: * Must contain only lowercase letters, numbers, and hyphens. * Must start with a letter. * Must be between 1-63 characters. * Must end with a number or a letter. * Must be unique within the OS policy.
    
* `pkg` -
  (Optional)
  Package resource
    
* `repository` -
  (Optional)
  Package repository resource
    
The `validate` block supports:
    
* `args` -
  (Optional)
  Optional arguments to pass to the source during execution.
    
* `file` -
  (Optional)
  A remote or local file.
    
* `interpreter` -
  (Required)
  Required. The script interpreter to use. Possible values: INTERPRETER_UNSPECIFIED, NONE, SHELL, POWERSHELL
    
* `output_file_path` -
  (Optional)
  Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.
    
* `script` -
  (Optional)
  An inline script. The size of the script is limited to 1024 characters.
    
The `source` block supports:
    
* `allow_insecure` -
  (Optional)
  Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.
    
* `gcs` -
  (Optional)
  A Cloud Storage object.
    
* `local_path` -
  (Optional)
  A local path within the VM to use.
    
* `remote` -
  (Optional)
  A generic remote file.
    
The `source` block supports:
    
* `allow_insecure` -
  (Optional)
  Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.
    
* `gcs` -
  (Optional)
  A Cloud Storage object.
    
* `local_path` -
  (Optional)
  A local path within the VM to use.
    
* `remote` -
  (Optional)
  A generic remote file.
    
The `source` block supports:
    
* `allow_insecure` -
  (Optional)
  Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.
    
* `gcs` -
  (Optional)
  A Cloud Storage object.
    
* `local_path` -
  (Optional)
  A local path within the VM to use.
    
* `remote` -
  (Optional)
  A generic remote file.
    
The `rollout` block supports:
    
* `disruption_budget` -
  (Required)
  Required. The maximum number (or percentage) of VMs per zone to disrupt at any given moment.
    
* `min_wait_duration` -
  (Required)
  Required. This determines the minimum duration of time to wait after the configuration changes are applied through the current rollout. A VM continues to count towards the `disruption_budget` at least until this duration of time has passed after configuration changes are applied.
    
The `disruption_budget` block supports:
    
* `fixed` -
  (Optional)
  Specifies a fixed value.
    
* `percent` -
  (Optional)
  Specifies the relative value defined as a percentage, which will be multiplied by a reference value.
    
- - -

* `description` -
  (Optional)
  OS policy assignment description. Length of the description is limited to 1024 characters.
  
* `project` -
  (Optional)
  The project for the resource
  
* `skip_await_rollout` -
  (Optional)
  Set to true to skip awaiting rollout during resource creation and update.
  


The `exclusion_labels` block supports:
    
* `labels` -
  (Optional)
  Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.
    
The `inclusion_labels` block supports:
    
* `labels` -
  (Optional)
  Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.
    
The `inventories` block supports:
    
* `os_short_name` -
  (Required)
  Required. The OS short name
    
* `os_version` -
  (Optional)
  The OS version Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of `7`, specify the following value for this field `7.*` An empty string matches all OS versions.
    
The `inventory_filters` block supports:
    
* `os_short_name` -
  (Required)
  Required. The OS short name
    
* `os_version` -
  (Optional)
  The OS version Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of `7`, specify the following value for this field `7.*` An empty string matches all OS versions.
    
The `exec` block supports:
    
* `enforce` -
  (Optional)
  What to run to bring this resource into the desired state. An exit code of 100 indicates "success", any other exit code indicates a failure running enforce.
    
* `validate` -
  (Required)
  Required. What to run to validate this resource is in the desired state. An exit code of 100 indicates "in desired state", and exit code of 101 indicates "not in desired state". Any other exit code indicates a failure running validate.
    
The `enforce` block supports:
    
* `args` -
  (Optional)
  Optional arguments to pass to the source during execution.
    
* `file` -
  (Optional)
  A remote or local file.
    
* `interpreter` -
  (Required)
  Required. The script interpreter to use. Possible values: INTERPRETER_UNSPECIFIED, NONE, SHELL, POWERSHELL
    
* `output_file_path` -
  (Optional)
  Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.
    
* `script` -
  (Optional)
  An inline script. The size of the script is limited to 1024 characters.
    
The `file` block supports:
    
* `allow_insecure` -
  (Optional)
  Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.
    
* `gcs` -
  (Optional)
  A Cloud Storage object.
    
* `local_path` -
  (Optional)
  A local path within the VM to use.
    
* `remote` -
  (Optional)
  A generic remote file.
    
The `gcs` block supports:
    
* `bucket` -
  (Required)
  Required. Bucket of the Cloud Storage object.
    
* `generation` -
  (Optional)
  Generation number of the Cloud Storage object.
    
* `object` -
  (Required)
  Required. Name of the Cloud Storage object.
    
The `remote` block supports:
    
* `sha256_checksum` -
  (Optional)
  SHA256 checksum of the remote file.
    
* `uri` -
  (Required)
  Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.
    
The `file` block supports:
    
* `allow_insecure` -
  (Optional)
  Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.
    
* `gcs` -
  (Optional)
  A Cloud Storage object.
    
* `local_path` -
  (Optional)
  A local path within the VM to use.
    
* `remote` -
  (Optional)
  A generic remote file.
    
The `gcs` block supports:
    
* `bucket` -
  (Required)
  Required. Bucket of the Cloud Storage object.
    
* `generation` -
  (Optional)
  Generation number of the Cloud Storage object.
    
* `object` -
  (Required)
  Required. Name of the Cloud Storage object.
    
The `remote` block supports:
    
* `sha256_checksum` -
  (Optional)
  SHA256 checksum of the remote file.
    
* `uri` -
  (Required)
  Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.
    
The `file` block supports:
    
* `content` -
  (Optional)
  A a file with this content. The size of the content is limited to 1024 characters.
    
* `file` -
  (Optional)
  A remote or local source.
    
* `path` -
  (Required)
  Required. The absolute path of the file within the VM.
    
* `permissions` -
  Consists of three octal digits which represent, in order, the permissions of the owner, group, and other users for the file (similarly to the numeric mode used in the linux chmod utility). Each digit represents a three bit number with the 4 bit corresponding to the read permissions, the 2 bit corresponds to the write bit, and the one bit corresponds to the execute permission. Default behavior is 755. Below are some examples of permissions and their associated values: read, write, and execute: 7 read and execute: 5 read and write: 6 read only: 4
    
* `state` -
  (Required)
  Required. Desired state of the file. Possible values: OS_POLICY_COMPLIANCE_STATE_UNSPECIFIED, COMPLIANT, NON_COMPLIANT, UNKNOWN, NO_OS_POLICIES_APPLICABLE
    
The `file` block supports:
    
* `allow_insecure` -
  (Optional)
  Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.
    
* `gcs` -
  (Optional)
  A Cloud Storage object.
    
* `local_path` -
  (Optional)
  A local path within the VM to use.
    
* `remote` -
  (Optional)
  A generic remote file.
    
The `gcs` block supports:
    
* `bucket` -
  (Required)
  Required. Bucket of the Cloud Storage object.
    
* `generation` -
  (Optional)
  Generation number of the Cloud Storage object.
    
* `object` -
  (Required)
  Required. Name of the Cloud Storage object.
    
The `remote` block supports:
    
* `sha256_checksum` -
  (Optional)
  SHA256 checksum of the remote file.
    
* `uri` -
  (Required)
  Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.
    
The `pkg` block supports:
    
* `apt` -
  (Optional)
  A package managed by Apt.
    
* `deb` -
  (Optional)
  A deb package file.
    
* `desired_state` -
  (Required)
  Required. The desired state the agent should maintain for this package. Possible values: DESIRED_STATE_UNSPECIFIED, INSTALLED, REMOVED
    
* `googet` -
  (Optional)
  A package managed by GooGet.
    
* `msi` -
  (Optional)
  An MSI package.
    
* `rpm` -
  (Optional)
  An rpm package file.
    
* `yum` -
  (Optional)
  A package managed by YUM.
    
* `zypper` -
  (Optional)
  A package managed by Zypper.
    
The `apt` block supports:
    
* `name` -
  (Required)
  Required. Package name.
    
The `deb` block supports:
    
* `pull_deps` -
  (Optional)
  Whether dependencies should also be installed. - install when false: `dpkg -i package` - install when true: `apt-get update && apt-get -y install package.deb`
    
* `source` -
  (Required)
  Required. A deb package.
    
The `gcs` block supports:
    
* `bucket` -
  (Required)
  Required. Bucket of the Cloud Storage object.
    
* `generation` -
  (Optional)
  Generation number of the Cloud Storage object.
    
* `object` -
  (Required)
  Required. Name of the Cloud Storage object.
    
The `remote` block supports:
    
* `sha256_checksum` -
  (Optional)
  SHA256 checksum of the remote file.
    
* `uri` -
  (Required)
  Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.
    
The `googet` block supports:
    
* `name` -
  (Required)
  Required. Package name.
    
The `msi` block supports:
    
* `properties` -
  (Optional)
  Additional properties to use during installation. This should be in the format of Property=Setting. Appended to the defaults of `ACTION=INSTALL REBOOT=ReallySuppress`.
    
* `source` -
  (Required)
  Required. The MSI package.
    
The `gcs` block supports:
    
* `bucket` -
  (Required)
  Required. Bucket of the Cloud Storage object.
    
* `generation` -
  (Optional)
  Generation number of the Cloud Storage object.
    
* `object` -
  (Required)
  Required. Name of the Cloud Storage object.
    
The `remote` block supports:
    
* `sha256_checksum` -
  (Optional)
  SHA256 checksum of the remote file.
    
* `uri` -
  (Required)
  Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.
    
The `rpm` block supports:
    
* `pull_deps` -
  (Optional)
  Whether dependencies should also be installed. - install when false: `rpm --upgrade --replacepkgs package.rpm` - install when true: `yum -y install package.rpm` or `zypper -y install package.rpm`
    
* `source` -
  (Required)
  Required. An rpm package.
    
The `gcs` block supports:
    
* `bucket` -
  (Required)
  Required. Bucket of the Cloud Storage object.
    
* `generation` -
  (Optional)
  Generation number of the Cloud Storage object.
    
* `object` -
  (Required)
  Required. Name of the Cloud Storage object.
    
The `remote` block supports:
    
* `sha256_checksum` -
  (Optional)
  SHA256 checksum of the remote file.
    
* `uri` -
  (Required)
  Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.
    
The `yum` block supports:
    
* `name` -
  (Required)
  Required. Package name.
    
The `zypper` block supports:
    
* `name` -
  (Required)
  Required. Package name.
    
The `repository` block supports:
    
* `apt` -
  (Optional)
  An Apt Repository.
    
* `goo` -
  (Optional)
  A Goo Repository.
    
* `yum` -
  (Optional)
  A Yum Repository.
    
* `zypper` -
  (Optional)
  A Zypper Repository.
    
The `apt` block supports:
    
* `archive_type` -
  (Required)
  Required. Type of archive files in this repository. Possible values: ARCHIVE_TYPE_UNSPECIFIED, DEB, DEB_SRC
    
* `components` -
  (Required)
  Required. List of components for this repository. Must contain at least one item.
    
* `distribution` -
  (Required)
  Required. Distribution of this repository.
    
* `gpg_key` -
  (Optional)
  URI of the key file for this repository. The agent maintains a keyring at `/etc/apt/trusted.gpg.d/osconfig_agent_managed.gpg`.
    
* `uri` -
  (Required)
  Required. URI for this repository.
    
The `goo` block supports:
    
* `name` -
  (Required)
  Required. The name of the repository.
    
* `url` -
  (Required)
  Required. The url of the repository.
    
The `yum` block supports:
    
* `base_url` -
  (Required)
  Required. The location of the repository directory.
    
* `display_name` -
  (Optional)
  The display name of the repository.
    
* `gpg_keys` -
  (Optional)
  URIs of GPG keys.
    
* `id` -
  (Required)
  Required. A one word, unique name for this repository. This is the `repo id` in the yum config file and also the `display_name` if `display_name` is omitted. This id is also used as the unique identifier when checking for resource conflicts.
    
The `zypper` block supports:
    
* `base_url` -
  (Required)
  Required. The location of the repository directory.
    
* `display_name` -
  (Optional)
  The display name of the repository.
    
* `gpg_keys` -
  (Optional)
  URIs of GPG keys.
    
* `id` -
  (Required)
  Required. A one word, unique name for this repository. This is the `repo id` in the zypper config file and also the `display_name` if `display_name` is omitted. This id is also used as the unique identifier when checking for GuestPolicy conflicts.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}`

* `baseline` -
  Output only. Indicates that this revision has been successfully rolled out in this zone and new VMs will be assigned OS policies from this revision. For a given OS policy assignment, there is only one revision with a value of `true` for this field.
  
* `deleted` -
  Output only. Indicates that this revision deletes the OS policy assignment.
  
* `etag` -
  The etag for this OS policy assignment. If this is provided on update, it must match the server's etag.
  
* `reconciling` -
  Output only. Indicates that reconciliation is in progress for the revision. This value is `true` when the `rollout_state` is one of: * IN_PROGRESS * CANCELLING
  
* `revision_create_time` -
  Output only. The timestamp that the revision was created.
  
* `revision_id` -
  Output only. The assignment revision ID A new revision is committed whenever a rollout is triggered for a OS policy assignment
  
* `rollout_state` -
  Output only. OS policy assignment rollout state Possible values: ROLLOUT_STATE_UNSPECIFIED, IN_PROGRESS, CANCELLING, CANCELLED, SUCCEEDED
  
* `uid` -
  Output only. Server generated unique id for the OS policy assignment resource.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

OSPolicyAssignment can be imported using any of these accepted formats:

```
$ terraform import google_os_config_os_policy_assignment.default projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}
$ terraform import google_os_config_os_policy_assignment.default {{project}}/{{location}}/{{name}}
$ terraform import google_os_config_os_policy_assignment.default {{location}}/{{name}}
```



