subcategory: "OS Config"

description: |-
    OS policy assignment is an API resource that is used to apply a set of OS policies to a dynamically targeted group of Compute Engine VM instances.
---

# google_os_config_os_policy_assignment

OS policy assignment is an API resource that is used to apply a set of OS
policies to a dynamically targeted group of Compute Engine VM instances. An OS
policy is used to define the desired state configuration for a Compute Engine VM
instance through a set of configuration resources that provide capabilities such
as installing or removing software packages, or executing a script. For more
information about the OS policy resource definitions and examples, see
[OS policy and OS policy assignment](https://cloud.google.com/compute/docs/os-configuration-management/working-with-os-policies).

To get more information about OSPolicyAssignment, see:

*   [API documentation](https://cloud.google.com/compute/docs/osconfig/rest/v1/projects.locations.osPolicyAssignments)
*   How-to Guides
    *   [Official Documentation](https://cloud.google.com/compute/docs/os-configuration-management/create-os-policy-assignment)

<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_working_dir=os_config_os_policy_assignment_basic&cloudshell_image=gcr.io%2Fgraphite-cloud-shell-images%2Fterraform%3Alatest&open_in_editor=main.tf&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

## Example Usage - Os Config Os Policy Assignment Basic

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

  location = "us-central1-a"
  name     = "policy-assignment"

  os_policies {
    id   = "policy"
    mode = "VALIDATION"

    resource_groups {
      resources {
        id = "apt-to-yum"

        repository {
          apt {
            archive_type = "DEB"
            components   = ["doc"]
            distribution = "debian"
            uri          = "https://atl.mirrors.clouvider.net/debian"
            gpg_key      = ".gnupg/pubring.kbx"
          }
        }
      }
      inventory_filters {
        os_short_name = "centos"
        os_version    = "8.*"
      }

      resources {
        id = "exec1"
        exec {
          validate {
            interpreter = "SHELL"
            args        = ["arg1"]
            file {
              local_path = "$HOME/script.sh"
            }
            output_file_path = "$HOME/out"
          }
          enforce {
            interpreter = "SHELL"
            args        = ["arg1"]
            file {
              allow_insecure = true
              remote {
                uri             = "https://www.example.com/script.sh"
                sha256_checksum = "c7938fed83afdccbb0e86a2a2e4cad7d5035012ca3214b4a61268393635c3063"
              }
            }
            output_file_path = "$HOME/out"
          }
        }
      }
    }
    allow_no_resource_group_match = false
    description                   = "A test os policy"
  }

  rollout {
    disruption_budget {
      percent = 100
    }

    min_wait_duration = "3s"
  }

  description = "A test os policy assignment"
}
```

## Argument Reference

The following arguments are supported:

*   `name` - (Required) Resource name.

*   `os_policies` - (Required) List of OS policies to be applied to the VMs.
    Structure is [documented below](#nested_os_policies).

*   `instance_filter` - (Required) Filter to select VMs. Structure is
    [documented below](#nested_instance_filter).

*   `rollout` - (Required) Rollout to deploy the OS policy assignment. A rollout
    is triggered in the following situations: 1) OSPolicyAssignment is created.
    2) OSPolicyAssignment is updated and the update contains changes to one of
    the following fields: - instance_filter - os_policies 3) OSPolicyAssignment
    is deleted. Structure is [documented below](#nested_rollout).

*   `location` - (Required) The location for the resource

<a name="nested_os_policies"></a>The `os_policies` block supports:

*   `id` - (Required) The id of the OS policy with the following restrictions:

    *   Must contain only lowercase letters, numbers, and hyphens.
    *   Must start with a letter.
    *   Must be between 1-63 characters.
    *   Must end with a number or a letter.
    *   Must be unique within the assignment.

*   `description` - (Optional) Policy description. Length of the description is
    limited to 1024 characters.

*   `mode` - (Required) Policy mode Possible values are: `MODE_UNSPECIFIED`,
    `VALIDATION`, `ENFORCEMENT`.

*   `resource_groups` - (Required) List of resource groups for the policy. For a
    particular VM, resource groups are evaluated in the order specified and the
    first resource group that is applicable is selected and the rest are
    ignored. If none of the resource groups are applicable for a VM, the VM is
    considered to be non-compliant w.r.t this policy. This behavior can be
    toggled by the flag `allow_no_resource_group_match` Structure is
    [documented below](#nested_resource_groups).

*   `allow_no_resource_group_match` - (Optional) This flag determines the OS
    policy compliance status when none of the resource groups within the policy
    are applicable for a VM. Set this value to `true` if the policy needs to be
    reported as compliant even if the policy has nothing to validate or enforce.

<a name="nested_resource_groups"></a>The `resource_groups` block supports:

*   `inventory_filters` - (Optional) List of inventory filters for the resource
    group. The resources in this resource group are applied to the target VM if
    it satisfies at least one of the following inventory filters. For example,
    to apply this resource group to VMs running either `RHEL` or `CentOS`
    operating systems, specify 2 items for the list with following values:
    inventory_filters[0].os_short_name='rhel' and
    inventory_filters[1].os_short_name='centos' If the list is empty, this
    resource group will be applied to the target VM unconditionally. Structure
    is [documented below](#nested_inventory_filters).

*   `resources` - (Required) List of resources configured for this resource
    group. The resources are executed in the exact order specified here.
    Structure is [documented below](#nested_resources).

<a name="nested_inventory_filters"></a>The `inventory_filters` block supports:

*   `os_short_name` - (Required) The OS short name

*   `os_version` - (Optional) The OS version Prefix matches are supported if
    asterisk(*) is provided as the last character. For example, to match all
    versions with a major version of `7`, specify the following value for this
    field `7.*` An empty string matches all OS versions.

<a name="nested_resources"></a>The `resources` block supports:

*   `id` - (Required) The id of the resource with the following restrictions:

    *   Must contain only lowercase letters, numbers, and hyphens.
    *   Must start with a letter.
    *   Must be between 1-63 characters.
    *   Must end with a number or a letter.
    *   Must be unique within the OS policy.

*   `pkg` - (Optional) Package resource Structure is
    [documented below](#nested_pkg).

*   `repository` - (Optional) Package repository resource Structure is
    [documented below](#nested_repository).

*   `exec` - (Optional) Exec resource Structure is
    [documented below](#nested_exec).

*   `file` - (Optional) File resource Structure is
    [documented below](#nested_file).

<a name="nested_pkg"></a>The `pkg` block supports:

*   `desired_state` - (Required) The desired state the agent should maintain for
    this package. Possible values are: `DESIRED_STATE_UNSPECIFIED`, `INSTALLED`,
    `REMOVED`.

*   `apt` - (Optional) A package managed by Apt. Structure is
    [documented below](#nested_apt).

*   `deb` - (Optional) A deb package file. Structure is
    [documented below](#nested_deb).

*   `yum` - (Optional) A package managed by YUM. Structure is
    [documented below](#nested_yum).

*   `zypper` - (Optional) A package managed by Zypper. Structure is
    [documented below](#nested_zypper).

*   `rpm` - (Optional) An rpm package file. Structure is
    [documented below](#nested_rpm).

*   `googet` - (Optional) A package managed by GooGet. Structure is
    [documented below](#nested_googet).

*   `msi` - (Optional) An MSI package. Structure is
    [documented below](#nested_msi).

<a name="nested_apt"></a>The `apt` block supports:

*   `name` - (Required) Package name.

<a name="nested_deb"></a>The `deb` block supports:

*   `source` - (Required) A deb package. Structure is
    [documented below](#nested_source).

*   `pull_deps` - (Optional) Whether dependencies should also be installed. -
    install when false: `dpkg -i package` - install when true: `apt-get update
    && apt-get -y install package.deb`

<a name="nested_source"></a>The `source` block supports:

*   `remote` - (Optional) A generic remote file. Structure is
    [documented below](#nested_remote).

*   `gcs` - (Optional) A Cloud Storage object. Structure is
    [documented below](#nested_gcs).

*   `local_path` - (Optional) A local path within the VM to use.

*   `allow_insecure` - (Optional) Defaults to false. When false, files are
    subject to validations based on the file type: Remote: A checksum must be
    specified. Cloud Storage: An object generation number must be specified.

<a name="nested_remote"></a>The `remote` block supports:

*   `uri` - (Required) URI from which to fetch the object. It should contain
    both the protocol and path following the format `{protocol}://{location}`.

*   `sha256_checksum` - (Optional) SHA256 checksum of the remote file.

<a name="nested_gcs"></a>The `gcs` block supports:

*   `bucket` - (Required) Bucket of the Cloud Storage object.

*   `object` - (Required) Name of the Cloud Storage object.

*   `generation` - (Optional) Generation number of the Cloud Storage object.

<a name="nested_yum"></a>The `yum` block supports:

*   `name` - (Required) Package name.

<a name="nested_zypper"></a>The `zypper` block supports:

*   `name` - (Required) Package name.

<a name="nested_rpm"></a>The `rpm` block supports:

*   `source` - (Required) An rpm package. Structure is
    [documented below](#nested_source).

*   `pull_deps` - (Optional) Whether dependencies should also be installed. -
    install when false: `rpm --upgrade --replacepkgs package.rpm` - install when
    true: `yum -y install package.rpm` or `zypper -y install package.rpm`

<a name="nested_source"></a>The `source` block supports:

*   `remote` - (Optional) A generic remote file. Structure is
    [documented below](#nested_remote).

*   `gcs` - (Optional) A Cloud Storage object. Structure is
    [documented below](#nested_gcs).

*   `local_path` - (Optional) A local path within the VM to use.

*   `allow_insecure` - (Optional) Defaults to false. When false, files are
    subject to validations based on the file type: Remote: A checksum must be
    specified. Cloud Storage: An object generation number must be specified.

<a name="nested_remote"></a>The `remote` block supports:

*   `uri` - (Required) URI from which to fetch the object. It should contain
    both the protocol and path following the format `{protocol}://{location}`.

*   `sha256_checksum` - (Optional) SHA256 checksum of the remote file.

<a name="nested_gcs"></a>The `gcs` block supports:

*   `bucket` - (Required) Bucket of the Cloud Storage object.

*   `object` - (Required) Name of the Cloud Storage object.

*   `generation` - (Optional) Generation number of the Cloud Storage object.

<a name="nested_googet"></a>The `googet` block supports:

*   `name` - (Required) Package name.

<a name="nested_msi"></a>The `msi` block supports:

*   `source` - (Required) The MSI package. Structure is
    [documented below](#nested_source).

*   `properties` - (Optional) Additional properties to use during installation.
    This should be in the format of Property=Setting. Appended to the defaults
    of `ACTION=INSTALL REBOOT=ReallySuppress`.

<a name="nested_source"></a>The `source` block supports:

*   `remote` - (Optional) A generic remote file. Structure is
    [documented below](#nested_remote).

*   `gcs` - (Optional) A Cloud Storage object. Structure is
    [documented below](#nested_gcs).

*   `local_path` - (Optional) A local path within the VM to use.

*   `allow_insecure` - (Optional) Defaults to false. When false, files are
    subject to validations based on the file type: Remote: A checksum must be
    specified. Cloud Storage: An object generation number must be specified.

<a name="nested_remote"></a>The `remote` block supports:

*   `uri` - (Required) URI from which to fetch the object. It should contain
    both the protocol and path following the format `{protocol}://{location}`.

*   `sha256_checksum` - (Optional) SHA256 checksum of the remote file.

<a name="nested_gcs"></a>The `gcs` block supports:

*   `bucket` - (Required) Bucket of the Cloud Storage object.

*   `object` - (Required) Name of the Cloud Storage object.

*   `generation` - (Optional) Generation number of the Cloud Storage object.

<a name="nested_repository"></a>The `repository` block supports:

*   `apt` - (Optional) An Apt Repository. Structure is
    [documented below](#nested_apt).

*   `yum` - (Optional) A Yum Repository. Structure is
    [documented below](#nested_yum).

*   `zypper` - (Optional) A Zypper Repository. Structure is
    [documented below](#nested_zypper).

*   `goo` - (Optional) A Goo Repository. Structure is
    [documented below](#nested_goo).

<a name="nested_apt"></a>The `apt` block supports:

*   `archive_type` - (Required) Type of archive files in this repository.
    Possible values are: `ARCHIVE_TYPE_UNSPECIFIED`, `DEB`, `DEB_SRC`.

*   `uri` - (Required) URI for this repository.

*   `distribution` - (Required) Distribution of this repository.

*   `components` - (Required) List of components for this repository. Must
    contain at least one item.

*   `gpg_key` - (Optional) URI of the key file for this repository. The agent
    maintains a keyring at `/etc/apt/trusted.gpg.d/osconfig_agent_managed.gpg`.

<a name="nested_yum"></a>The `yum` block supports:

*   `id` - (Required) A one word, unique name for this repository. This is the
    `repo id` in the yum config file and also the `display_name` if
    `display_name` is omitted. This id is also used as the unique identifier
    when checking for resource conflicts.

*   `display_name` - (Optional) The display name of the repository.

*   `base_url` - (Required) The location of the repository directory.

*   `gpg_keys` - (Optional) URIs of GPG keys.

<a name="nested_zypper"></a>The `zypper` block supports:

*   `id` - (Required) A one word, unique name for this repository. This is the
    `repo id` in the zypper config file and also the `display_name` if
    `display_name` is omitted. This id is also used as the unique identifier
    when checking for GuestPolicy conflicts.

*   `display_name` - (Optional) The display name of the repository.

*   `base_url` - (Required) The location of the repository directory.

*   `gpg_keys` - (Optional) URIs of GPG keys.

<a name="nested_goo"></a>The `goo` block supports:

*   `name` - (Required) The name of the repository.

*   `url` - (Required) The url of the repository.

<a name="nested_exec"></a>The `exec` block supports:

*   `validate` - (Required) What to run to validate this resource is in the
    desired state. An exit code of 100 indicates "in desired state", and exit
    code of 101 indicates "not in desired state". Any other exit code indicates
    a failure running validate. Structure is
    [documented below](#nested_validate).

*   `enforce` - (Optional) What to run to bring this resource into the desired
    state. An exit code of 100 indicates "success", any other exit code
    indicates a failure running enforce. Structure is
    [documented below](#nested_enforce).

<a name="nested_validate"></a>The `validate` block supports:

*   `file` - (Optional) A remote or local file. Structure is
    [documented below](#nested_file).

*   `script` - (Optional) An inline script. The size of the script is limited to
    1024 characters.

*   `args` - (Optional) Optional arguments to pass to the source during
    execution.

*   `interpreter` - (Required) The script interpreter to use. Possible values
    are: `INTERPRETER_UNSPECIFIED`, `NONE`, `SHELL`, `POWERSHELL`.

*   `output_file_path` - (Optional) Only recorded for enforce Exec. Path to an
    output file (that is created by this Exec) whose content will be recorded in
    OSPolicyResourceCompliance after a successful run. Absence or failure to
    read this file will result in this ExecResource being non-compliant. Output
    file size is limited to 100K bytes.

<a name="nested_file"></a>The `file` block supports:

*   `remote` - (Optional) A generic remote file. Structure is
    [documented below](#nested_remote).

*   `gcs` - (Optional) A Cloud Storage object. Structure is
    [documented below](#nested_gcs).

*   `local_path` - (Optional) A local path within the VM to use.

*   `allow_insecure` - (Optional) Defaults to false. When false, files are
    subject to validations based on the file type: Remote: A checksum must be
    specified. Cloud Storage: An object generation number must be specified.

<a name="nested_remote"></a>The `remote` block supports:

*   `uri` - (Required) URI from which to fetch the object. It should contain
    both the protocol and path following the format `{protocol}://{location}`.

*   `sha256_checksum` - (Optional) SHA256 checksum of the remote file.

<a name="nested_gcs"></a>The `gcs` block supports:

*   `bucket` - (Required) Bucket of the Cloud Storage object.

*   `object` - (Required) Name of the Cloud Storage object.

*   `generation` - (Optional) Generation number of the Cloud Storage object.

<a name="nested_enforce"></a>The `enforce` block supports:

*   `file` - (Optional) A remote or local file. Structure is
    [documented below](#nested_file).

*   `script` - (Optional) An inline script. The size of the script is limited to
    1024 characters.

*   `args` - (Optional) Optional arguments to pass to the source during
    execution.

*   `interpreter` - (Required) The script interpreter to use. Possible values
    are: `INTERPRETER_UNSPECIFIED`, `NONE`, `SHELL`, `POWERSHELL`.

*   `output_file_path` - (Optional) Only recorded for enforce Exec. Path to an
    output file (that is created by this Exec) whose content will be recorded in
    OSPolicyResourceCompliance after a successful run. Absence or failure to
    read this file will result in this ExecResource being non-compliant. Output
    file size is limited to 100K bytes.

<a name="nested_file"></a>The `file` block supports:

*   `remote` - (Optional) A generic remote file. Structure is
    [documented below](#nested_remote).

*   `gcs` - (Optional) A Cloud Storage object. Structure is
    [documented below](#nested_gcs).

*   `local_path` - (Optional) A local path within the VM to use.

*   `allow_insecure` - (Optional) Defaults to false. When false, files are
    subject to validations based on the file type: Remote: A checksum must be
    specified. Cloud Storage: An object generation number must be specified.

<a name="nested_remote"></a>The `remote` block supports:

*   `uri` - (Required) URI from which to fetch the object. It should contain
    both the protocol and path following the format `{protocol}://{location}`.

*   `sha256_checksum` - (Optional) SHA256 checksum of the remote file.

<a name="nested_gcs"></a>The `gcs` block supports:

*   `bucket` - (Required) Bucket of the Cloud Storage object.

*   `object` - (Required) Name of the Cloud Storage object.

*   `generation` - (Optional) Generation number of the Cloud Storage object.

<a name="nested_file"></a>The `file` block supports:

*   `file` - (Optional) A remote or local source. Structure is
    [documented below](#nested_file).

*   `content` - (Optional) A a file with this content. The size of the content
    is limited to 1024 characters.

*   `path` - (Required) The absolute path of the file within the VM.

*   `state` - (Required) Desired state of the file. Possible values are:
    `DESIRED_STATE_UNSPECIFIED`, `PRESENT`, `ABSENT`, `CONTENTS_MATCH`.

*   `permissions` - (Output) Consists of three octal digits which represent, in
    order, the permissions of the owner, group, and other users for the file
    (similarly to the numeric mode used in the linux chmod utility). Each digit
    represents a three bit number with the 4 bit corresponding to the read
    permissions, the 2 bit corresponds to the write bit, and the one bit
    corresponds to the execute permission. Default behavior is 755. Below are
    some examples of permissions and their associated values: read, write, and
    execute: 7 read and execute: 5 read and write: 6 read only: 4

<a name="nested_file"></a>The `file` block supports:

*   `remote` - (Optional) A generic remote file. Structure is
    [documented below](#nested_remote).

*   `gcs` - (Optional) A Cloud Storage object. Structure is
    [documented below](#nested_gcs).

*   `local_path` - (Optional) A local path within the VM to use.

*   `allow_insecure` - (Optional) Defaults to false. When false, files are
    subject to validations based on the file type: Remote: A checksum must be
    specified. Cloud Storage: An object generation number must be specified.

<a name="nested_remote"></a>The `remote` block supports:

*   `uri` - (Required) URI from which to fetch the object. It should contain
    both the protocol and path following the format `{protocol}://{location}`.

*   `sha256_checksum` - (Optional) SHA256 checksum of the remote file.

<a name="nested_gcs"></a>The `gcs` block supports:

*   `bucket` - (Required) Bucket of the Cloud Storage object.

*   `object` - (Required) Name of the Cloud Storage object.

*   `generation` - (Optional) Generation number of the Cloud Storage object.

<a name="nested_instance_filter"></a>The `instance_filter` block supports:

*   `all` - (Optional) Target all VMs in the project. If true, no other criteria
    is permitted.

*   `inclusion_labels` - (Optional) List of label sets used for VM inclusion. If
    the list has more than one `LabelSet`, the VM is included if any of the
    label sets are applicable for the VM. Structure is
    [documented below](#nested_inclusion_labels).

*   `exclusion_labels` - (Optional) List of label sets used for VM exclusion. If
    the list has more than one label set, the VM is excluded if any of the label
    sets are applicable for the VM. Structure is
    [documented below](#nested_exclusion_labels).

*   `inventories` - (Optional) List of inventories to select VMs. A VM is
    selected if its inventory data matches at least one of the following
    inventories. Structure is [documented below](#nested_inventories).

<a name="nested_inclusion_labels"></a>The `inclusion_labels` block supports:

*   `labels` - (Optional) Labels are identified by key/value pairs in this map.
    A VM should contain all the key/value pairs specified in this map to be
    selected.

<a name="nested_exclusion_labels"></a>The `exclusion_labels` block supports:

*   `labels` - (Optional) Labels are identified by key/value pairs in this map.
    A VM should contain all the key/value pairs specified in this map to be
    selected.

<a name="nested_inventories"></a>The `inventories` block supports:

*   `os_short_name` - (Required) The OS short name

*   `os_version` - (Optional) The OS version Prefix matches are supported if
    asterisk(*) is provided as the last character. For example, to match all
    versions with a major version of `7`, specify the following value for this
    field `7.*` An empty string matches all OS versions.

<a name="nested_rollout"></a>The `rollout` block supports:

*   `disruption_budget` - (Required) The maximum number (or percentage) of VMs
    per zone to disrupt at any given moment. Structure is
    [documented below](#nested_disruption_budget).

*   `min_wait_duration` - (Required) This determines the minimum duration of
    time to wait after the configuration changes are applied through the current
    rollout. A VM continues to count towards the `disruption_budget` at least
    until this duration of time has passed after configuration changes are
    applied.

<a name="nested_disruption_budget"></a>The `disruption_budget` block supports:

*   `fixed` - (Optional) Specifies a fixed value.

*   `percent` - (Optional) Specifies the relative value defined as a percentage,
    which will be multiplied by a reference value.

--------------------------------------------------------------------------------

*   `description` - (Optional) OS policy assignment description. Length of the
    description is limited to 1024 characters.

*   `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

*   `skip_await_rollout` - (Optional) Set to true to skip awaiting rollout
    during resource creation and update.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

*   `id` - an identifier for the resource with format
    `projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}`

*   `revision_id` - Output only. The assignment revision ID A new revision is
    committed whenever a rollout is triggered for a OS policy assignment

*   `revision_create_time` - Output only. The timestamp that the revision was
    created.

*   `etag` - The etag for this OS policy assignment. If this is provided on
    update, it must match the server's etag.

*   `rollout_state` - Output only. OS policy assignment rollout state

*   `baseline` - Output only. Indicates that this revision has been successfully
    rolled out in this zone and new VMs will be assigned OS policies from this
    revision. For a given OS policy assignment, there is only one revision with
    a value of `true` for this field.

*   `deleted` - Output only. Indicates that this revision deletes the OS policy
    assignment.

*   `reconciling` - Output only. Indicates that reconciliation is in progress
    for the revision. This value is `true` when the `rollout_state` is one of:

    *   IN_PROGRESS
    *   CANCELLING

*   `uid` - Output only. Server generated unique id for the OS policy assignment
    resource.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts)
configuration options:

-   `create` - Default is 20 minutes.
-   `update` - Default is 20 minutes.
-   `delete` - Default is 20 minutes.

## Import

OSPolicyAssignment can be imported using any of these accepted formats:

```
$ terraform import google_os_config_os_policy_assignment.default projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}
$ terraform import google_os_config_os_policy_assignment.default {{project}}/{{location}}/{{name}}
$ terraform import google_os_config_os_policy_assignment.default {{location}}/{{name}}
```

## User Project Overrides

This resource supports
[User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
