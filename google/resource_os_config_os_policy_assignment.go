// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	osconfig "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/osconfig"
)

func resourceOsConfigOsPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOsConfigOsPolicyAssignmentCreate,
		Read:   resourceOsConfigOsPolicyAssignmentRead,
		Update: resourceOsConfigOsPolicyAssignmentUpdate,
		Delete: resourceOsConfigOsPolicyAssignmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceOsConfigOsPolicyAssignmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_filter": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Filter to select VMs.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentInstanceFilterSchema(),
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource name.",
			},

			"os_policies": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. List of OS policies to be applied to the VMs.",
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesSchema(),
			},

			"rollout": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Rollout to deploy the OS policy assignment. A rollout is triggered in the following situations: 1) OSPolicyAssignment is created. 2) OSPolicyAssignment is updated and the update contains changes to one of the following fields: - instance_filter - os_policies 3) OSPolicyAssignment is deleted.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentRolloutSchema(),
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OS policy assignment description. Length of the description is limited to 1024 characters.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"skip_await_rollout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to true to skip awaiting rollout during resource creation and update.",
			},

			"baseline": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. Indicates that this revision has been successfully rolled out in this zone and new VMs will be assigned OS policies from this revision. For a given OS policy assignment, there is only one revision with a value of `true` for this field.",
			},

			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. Indicates that this revision deletes the OS policy assignment.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The etag for this OS policy assignment. If this is provided on update, it must match the server's etag.",
			},

			"reconciling": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Output only. Indicates that reconciliation is in progress for the revision. This value is `true` when the `rollout_state` is one of: * IN_PROGRESS * CANCELLING",
			},

			"revision_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The timestamp that the revision was created.",
			},

			"revision_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. The assignment revision ID A new revision is committed whenever a rollout is triggered for a OS policy assignment",
			},

			"rollout_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. OS policy assignment rollout state Possible values: ROLLOUT_STATE_UNSPECIFIED, IN_PROGRESS, CANCELLING, CANCELLED, SUCCEEDED",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Server generated unique id for the OS policy assignment resource.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentInstanceFilterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Target all VMs in the project. If true, no other criteria is permitted.",
			},

			"exclusion_labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of label sets used for VM exclusion. If the list has more than one label set, the VM is excluded if any of the label sets are applicable for the VM.",
				Elem:        OsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsSchema(),
			},

			"inclusion_labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of label sets used for VM inclusion. If the list has more than one `LabelSet`, the VM is included if any of the label sets are applicable for the VM.",
				Elem:        OsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsSchema(),
			},

			"inventories": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of inventories to select VMs. A VM is selected if its inventory data matches at least one of the following inventories.",
				Elem:        OsConfigOsPolicyAssignmentInstanceFilterInventoriesSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func OsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func OsConfigOsPolicyAssignmentInstanceFilterInventoriesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"os_short_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The OS short name",
			},

			"os_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OS version Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of `7`, specify the following value for this field `7.*` An empty string matches all OS versions.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The id of the OS policy with the following restrictions: * Must contain only lowercase letters, numbers, and hyphens. * Must start with a letter. * Must be between 1-63 characters. * Must end with a number or a letter. * Must be unique within the assignment.",
			},

			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Policy mode Possible values: MODE_UNSPECIFIED, VALIDATION, ENFORCEMENT",
			},

			"resource_groups": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. List of resource groups for the policy. For a particular VM, resource groups are evaluated in the order specified and the first resource group that is applicable is selected and the rest are ignored. If none of the resource groups are applicable for a VM, the VM is considered to be non-compliant w.r.t this policy. This behavior can be toggled by the flag `allow_no_resource_group_match`",
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsSchema(),
			},

			"allow_no_resource_group_match": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This flag determines the OS policy compliance status when none of the resource groups within the policy are applicable for a VM. Set this value to `true` if the policy needs to be reported as compliant even if the policy has nothing to validate or enforce.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy description. Length of the description is limited to 1024 characters.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. List of resources configured for this resource group. The resources are executed in the exact order specified here.",
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesSchema(),
			},

			"inventory_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of inventory filters for the resource group. The resources in this resource group are applied to the target VM if it satisfies at least one of the following inventory filters. For example, to apply this resource group to VMs running either `RHEL` or `CentOS` operating systems, specify 2 items for the list with following values: inventory_filters[0].os_short_name='rhel' and inventory_filters[1].os_short_name='centos' If the list is empty, this resource group will be applied to the target VM unconditionally.",
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The id of the resource with the following restrictions: * Must contain only lowercase letters, numbers, and hyphens. * Must start with a letter. * Must be between 1-63 characters. * Must end with a number or a letter. * Must be unique within the OS policy.",
			},

			"exec": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Exec resource",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSchema(),
			},

			"file": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "File resource",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSchema(),
			},

			"pkg": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Package resource",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSchema(),
			},

			"repository": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Package repository resource",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"validate": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. What to run to validate this resource is in the desired state. An exit code of 100 indicates \"in desired state\", and exit code of 101 indicates \"not in desired state\". Any other exit code indicates a failure running validate.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSchema(),
			},

			"enforce": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "What to run to bring this resource into the desired state. An exit code of 100 indicates \"success\", any other exit code indicates a failure running enforce.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interpreter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The script interpreter to use. Possible values: INTERPRETER_UNSPECIFIED, NONE, SHELL, POWERSHELL",
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional arguments to pass to the source during execution.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A remote or local file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSchema(),
			},

			"output_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.",
			},

			"script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An inline script. The size of the script is limited to 1024 characters.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
			},

			"gcs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Cloud Storage object.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSchema(),
			},

			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A local path within the VM to use.",
			},

			"remote": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A generic remote file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Bucket of the Cloud Storage object.",
			},

			"object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Cloud Storage object.",
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Generation number of the Cloud Storage object.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
			},

			"sha256_checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 checksum of the remote file.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interpreter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The script interpreter to use. Possible values: INTERPRETER_UNSPECIFIED, NONE, SHELL, POWERSHELL",
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional arguments to pass to the source during execution.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A remote or local file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSchema(),
			},

			"output_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.",
			},

			"script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An inline script. The size of the script is limited to 1024 characters.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
			},

			"gcs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Cloud Storage object.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSchema(),
			},

			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A local path within the VM to use.",
			},

			"remote": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A generic remote file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Bucket of the Cloud Storage object.",
			},

			"object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Cloud Storage object.",
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Generation number of the Cloud Storage object.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
			},

			"sha256_checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 checksum of the remote file.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The absolute path of the file within the VM.",
			},

			"state": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Desired state of the file. Possible values: OS_POLICY_COMPLIANCE_STATE_UNSPECIFIED, COMPLIANT, NON_COMPLIANT, UNKNOWN, NO_OS_POLICIES_APPLICABLE",
			},

			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A a file with this content. The size of the content is limited to 1024 characters.",
			},

			"file": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A remote or local source.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSchema(),
			},

			"permissions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consists of three octal digits which represent, in order, the permissions of the owner, group, and other users for the file (similarly to the numeric mode used in the linux chmod utility). Each digit represents a three bit number with the 4 bit corresponding to the read permissions, the 2 bit corresponds to the write bit, and the one bit corresponds to the execute permission. Default behavior is 755. Below are some examples of permissions and their associated values: read, write, and execute: 7 read and execute: 5 read and write: 6 read only: 4",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
			},

			"gcs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Cloud Storage object.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSchema(),
			},

			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A local path within the VM to use.",
			},

			"remote": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A generic remote file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Bucket of the Cloud Storage object.",
			},

			"object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Cloud Storage object.",
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Generation number of the Cloud Storage object.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
			},

			"sha256_checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 checksum of the remote file.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"desired_state": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The desired state the agent should maintain for this package. Possible values: DESIRED_STATE_UNSPECIFIED, INSTALLED, REMOVED",
			},

			"apt": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A package managed by Apt.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSchema(),
			},

			"deb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A deb package file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSchema(),
			},

			"googet": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A package managed by GooGet.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSchema(),
			},

			"msi": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An MSI package.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSchema(),
			},

			"rpm": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An rpm package file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSchema(),
			},

			"yum": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A package managed by YUM.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSchema(),
			},

			"zypper": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A package managed by Zypper.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Package name.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. A deb package.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSchema(),
			},

			"pull_deps": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether dependencies should also be installed. - install when false: `dpkg -i package` - install when true: `apt-get update && apt-get -y install package.deb`",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
			},

			"gcs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Cloud Storage object.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSchema(),
			},

			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A local path within the VM to use.",
			},

			"remote": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A generic remote file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Bucket of the Cloud Storage object.",
			},

			"object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Cloud Storage object.",
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Generation number of the Cloud Storage object.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
			},

			"sha256_checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 checksum of the remote file.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Package name.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The MSI package.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSchema(),
			},

			"properties": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Additional properties to use during installation. This should be in the format of Property=Setting. Appended to the defaults of `ACTION=INSTALL REBOOT=ReallySuppress`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
			},

			"gcs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Cloud Storage object.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSchema(),
			},

			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A local path within the VM to use.",
			},

			"remote": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A generic remote file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Bucket of the Cloud Storage object.",
			},

			"object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Cloud Storage object.",
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Generation number of the Cloud Storage object.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
			},

			"sha256_checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 checksum of the remote file.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. An rpm package.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSchema(),
			},

			"pull_deps": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether dependencies should also be installed. - install when false: `rpm --upgrade --replacepkgs package.rpm` - install when true: `yum -y install package.rpm` or `zypper -y install package.rpm`",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
			},

			"gcs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Cloud Storage object.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSchema(),
			},

			"local_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A local path within the VM to use.",
			},

			"remote": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A generic remote file.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Bucket of the Cloud Storage object.",
			},

			"object": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Cloud Storage object.",
			},

			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Generation number of the Cloud Storage object.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
			},

			"sha256_checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 checksum of the remote file.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Package name.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Package name.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"apt": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An Apt Repository.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSchema(),
			},

			"goo": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Goo Repository.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSchema(),
			},

			"yum": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Yum Repository.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSchema(),
			},

			"zypper": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Zypper Repository.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSchema(),
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"archive_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Type of archive files in this repository. Possible values: ARCHIVE_TYPE_UNSPECIFIED, DEB, DEB_SRC",
			},

			"components": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. List of components for this repository. Must contain at least one item.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"distribution": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Distribution of this repository.",
			},

			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. URI for this repository.",
			},

			"gpg_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI of the key file for this repository. The agent maintains a keyring at `/etc/apt/trusted.gpg.d/osconfig_agent_managed.gpg`.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The name of the repository.",
			},

			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The url of the repository.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The location of the repository directory.",
			},

			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. A one word, unique name for this repository. This is the `repo id` in the yum config file and also the `display_name` if `display_name` is omitted. This id is also used as the unique identifier when checking for resource conflicts.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the repository.",
			},

			"gpg_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "URIs of GPG keys.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The location of the repository directory.",
			},

			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. A one word, unique name for this repository. This is the `repo id` in the zypper config file and also the `display_name` if `display_name` is omitted. This id is also used as the unique identifier when checking for GuestPolicy conflicts.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the repository.",
			},

			"gpg_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "URIs of GPG keys.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func OsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"os_short_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The OS short name",
			},

			"os_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OS version Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of `7`, specify the following value for this field `7.*` An empty string matches all OS versions.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentRolloutSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disruption_budget": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The maximum number (or percentage) of VMs per zone to disrupt at any given moment.",
				MaxItems:    1,
				Elem:        OsConfigOsPolicyAssignmentRolloutDisruptionBudgetSchema(),
			},

			"min_wait_duration": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. This determines the minimum duration of time to wait after the configuration changes are applied through the current rollout. A VM continues to count towards the `disruption_budget` at least until this duration of time has passed after configuration changes are applied.",
			},
		},
	}
}

func OsConfigOsPolicyAssignmentRolloutDisruptionBudgetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fixed": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies a fixed value.",
			},

			"percent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the relative value defined as a percentage, which will be multiplied by a reference value.",
			},
		},
	}
}

func resourceOsConfigOsPolicyAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter:   expandOsConfigOsPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		OSPolicies:       expandOsConfigOsPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:          expandOsConfigOsPolicyAssignmentRollout(d.Get("rollout")),
		Description:      dcl.String(d.Get("description").(string)),
		Project:          dcl.String(project),
		SkipAwaitRollout: dcl.Bool(d.Get("skip_await_rollout").(bool)),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := CreateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLOsConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyOSPolicyAssignment(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating OSPolicyAssignment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating OSPolicyAssignment %q: %#v", d.Id(), res)

	return resourceOsConfigOsPolicyAssignmentRead(d, meta)
}

func resourceOsConfigOsPolicyAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter:   expandOsConfigOsPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		OSPolicies:       expandOsConfigOsPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:          expandOsConfigOsPolicyAssignmentRollout(d.Get("rollout")),
		Description:      dcl.String(d.Get("description").(string)),
		Project:          dcl.String(project),
		SkipAwaitRollout: dcl.Bool(d.Get("skip_await_rollout").(bool)),
	}

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLOsConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetOSPolicyAssignment(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("OsConfigOsPolicyAssignment %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("instance_filter", flattenOsConfigOsPolicyAssignmentInstanceFilter(res.InstanceFilter)); err != nil {
		return fmt.Errorf("error setting instance_filter in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("os_policies", flattenOsConfigOsPolicyAssignmentOSPoliciesArray(res.OSPolicies)); err != nil {
		return fmt.Errorf("error setting os_policies in state: %s", err)
	}
	if err = d.Set("rollout", flattenOsConfigOsPolicyAssignmentRollout(res.Rollout)); err != nil {
		return fmt.Errorf("error setting rollout in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("skip_await_rollout", res.SkipAwaitRollout); err != nil {
		return fmt.Errorf("error setting skip_await_rollout in state: %s", err)
	}
	if err = d.Set("baseline", res.Baseline); err != nil {
		return fmt.Errorf("error setting baseline in state: %s", err)
	}
	if err = d.Set("deleted", res.Deleted); err != nil {
		return fmt.Errorf("error setting deleted in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("reconciling", res.Reconciling); err != nil {
		return fmt.Errorf("error setting reconciling in state: %s", err)
	}
	if err = d.Set("revision_create_time", res.RevisionCreateTime); err != nil {
		return fmt.Errorf("error setting revision_create_time in state: %s", err)
	}
	if err = d.Set("revision_id", res.RevisionId); err != nil {
		return fmt.Errorf("error setting revision_id in state: %s", err)
	}
	if err = d.Set("rollout_state", res.RolloutState); err != nil {
		return fmt.Errorf("error setting rollout_state in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}

	return nil
}
func resourceOsConfigOsPolicyAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter:   expandOsConfigOsPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		OSPolicies:       expandOsConfigOsPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:          expandOsConfigOsPolicyAssignmentRollout(d.Get("rollout")),
		Description:      dcl.String(d.Get("description").(string)),
		Project:          dcl.String(project),
		SkipAwaitRollout: dcl.Bool(d.Get("skip_await_rollout").(bool)),
	}
	// Construct state hint from old values
	old := &osconfig.OSPolicyAssignment{
		InstanceFilter:   expandOsConfigOsPolicyAssignmentInstanceFilter(oldValue(d.GetChange("instance_filter"))),
		Location:         dcl.String(oldValue(d.GetChange("location")).(string)),
		Name:             dcl.String(oldValue(d.GetChange("name")).(string)),
		OSPolicies:       expandOsConfigOsPolicyAssignmentOSPoliciesArray(oldValue(d.GetChange("os_policies"))),
		Rollout:          expandOsConfigOsPolicyAssignmentRollout(oldValue(d.GetChange("rollout"))),
		Description:      dcl.String(oldValue(d.GetChange("description")).(string)),
		Project:          dcl.StringOrNil(oldValue(d.GetChange("project")).(string)),
		SkipAwaitRollout: dcl.Bool(oldValue(d.GetChange("skip_await_rollout")).(bool)),
	}
	directive := UpdateDirective
	directive = append(directive, dcl.WithStateHint(old))
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLOsConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyOSPolicyAssignment(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating OSPolicyAssignment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating OSPolicyAssignment %q: %#v", d.Id(), res)

	return resourceOsConfigOsPolicyAssignmentRead(d, meta)
}

func resourceOsConfigOsPolicyAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter:   expandOsConfigOsPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:         dcl.String(d.Get("location").(string)),
		Name:             dcl.String(d.Get("name").(string)),
		OSPolicies:       expandOsConfigOsPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:          expandOsConfigOsPolicyAssignmentRollout(d.Get("rollout")),
		Description:      dcl.String(d.Get("description").(string)),
		Project:          dcl.String(project),
		SkipAwaitRollout: dcl.Bool(d.Get("skip_await_rollout").(bool)),
	}

	log.Printf("[DEBUG] Deleting OSPolicyAssignment %q", d.Id())
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLOsConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteOSPolicyAssignment(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting OSPolicyAssignment: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting OSPolicyAssignment %q", d.Id())
	return nil
}

func resourceOsConfigOsPolicyAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/osPolicyAssignments/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandOsConfigOsPolicyAssignmentInstanceFilter(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilter {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilter
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilter
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilter{
		All:             dcl.Bool(obj["all"].(bool)),
		ExclusionLabels: expandOsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsArray(obj["exclusion_labels"]),
		InclusionLabels: expandOsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsArray(obj["inclusion_labels"]),
		Inventories:     expandOsConfigOsPolicyAssignmentInstanceFilterInventoriesArray(obj["inventories"]),
	}
}

func flattenOsConfigOsPolicyAssignmentInstanceFilter(obj *osconfig.OSPolicyAssignmentInstanceFilter) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"all":              obj.All,
		"exclusion_labels": flattenOsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsArray(obj.ExclusionLabels),
		"inclusion_labels": flattenOsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsArray(obj.InclusionLabels),
		"inventories":      flattenOsConfigOsPolicyAssignmentInstanceFilterInventoriesArray(obj.Inventories),
	}

	return []interface{}{transformed}

}
func expandOsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsArray(o interface{}) []osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels, 0, len(objs))
	for _, item := range objs {
		i := expandOsConfigOsPolicyAssignmentInstanceFilterExclusionLabels(item)
		items = append(items, *i)
	}

	return items
}

func expandOsConfigOsPolicyAssignmentInstanceFilterExclusionLabels(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilterExclusionLabels
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels{
		Labels: checkStringMap(obj["labels"]),
	}
}

func flattenOsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsArray(objs []osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOsConfigOsPolicyAssignmentInstanceFilterExclusionLabels(&item)
		items = append(items, i)
	}

	return items
}

func flattenOsConfigOsPolicyAssignmentInstanceFilterExclusionLabels(obj *osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"labels": obj.Labels,
	}

	return transformed

}
func expandOsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsArray(o interface{}) []osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels, 0, len(objs))
	for _, item := range objs {
		i := expandOsConfigOsPolicyAssignmentInstanceFilterInclusionLabels(item)
		items = append(items, *i)
	}

	return items
}

func expandOsConfigOsPolicyAssignmentInstanceFilterInclusionLabels(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilterInclusionLabels
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels{
		Labels: checkStringMap(obj["labels"]),
	}
}

func flattenOsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsArray(objs []osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOsConfigOsPolicyAssignmentInstanceFilterInclusionLabels(&item)
		items = append(items, i)
	}

	return items
}

func flattenOsConfigOsPolicyAssignmentInstanceFilterInclusionLabels(obj *osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"labels": obj.Labels,
	}

	return transformed

}
func expandOsConfigOsPolicyAssignmentInstanceFilterInventoriesArray(o interface{}) []osconfig.OSPolicyAssignmentInstanceFilterInventories {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInventories, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInventories, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentInstanceFilterInventories, 0, len(objs))
	for _, item := range objs {
		i := expandOsConfigOsPolicyAssignmentInstanceFilterInventories(item)
		items = append(items, *i)
	}

	return items
}

func expandOsConfigOsPolicyAssignmentInstanceFilterInventories(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilterInventories {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilterInventories
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilterInventories{
		OSShortName: dcl.String(obj["os_short_name"].(string)),
		OSVersion:   dcl.String(obj["os_version"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentInstanceFilterInventoriesArray(objs []osconfig.OSPolicyAssignmentInstanceFilterInventories) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOsConfigOsPolicyAssignmentInstanceFilterInventories(&item)
		items = append(items, i)
	}

	return items
}

func flattenOsConfigOsPolicyAssignmentInstanceFilterInventories(obj *osconfig.OSPolicyAssignmentInstanceFilterInventories) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"os_short_name": obj.OSShortName,
		"os_version":    obj.OSVersion,
	}

	return transformed

}
func expandOsConfigOsPolicyAssignmentOSPoliciesArray(o interface{}) []osconfig.OSPolicyAssignmentOSPolicies {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPolicies, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]osconfig.OSPolicyAssignmentOSPolicies, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPolicies, 0, len(objs))
	for _, item := range objs {
		i := expandOsConfigOsPolicyAssignmentOSPolicies(item)
		items = append(items, *i)
	}

	return items
}

func expandOsConfigOsPolicyAssignmentOSPolicies(o interface{}) *osconfig.OSPolicyAssignmentOSPolicies {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPolicies
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPolicies{
		Id:                        dcl.String(obj["id"].(string)),
		Mode:                      osconfig.OSPolicyAssignmentOSPoliciesModeEnumRef(obj["mode"].(string)),
		ResourceGroups:            expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsArray(obj["resource_groups"]),
		AllowNoResourceGroupMatch: dcl.Bool(obj["allow_no_resource_group_match"].(bool)),
		Description:               dcl.String(obj["description"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesArray(objs []osconfig.OSPolicyAssignmentOSPolicies) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOsConfigOsPolicyAssignmentOSPolicies(&item)
		items = append(items, i)
	}

	return items
}

func flattenOsConfigOsPolicyAssignmentOSPolicies(obj *osconfig.OSPolicyAssignmentOSPolicies) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"id":                            obj.Id,
		"mode":                          obj.Mode,
		"resource_groups":               flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsArray(obj.ResourceGroups),
		"allow_no_resource_group_match": obj.AllowNoResourceGroupMatch,
		"description":                   obj.Description,
	}

	return transformed

}
func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsArray(o interface{}) []osconfig.OSPolicyAssignmentOSPoliciesResourceGroups {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroups, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroups, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroups, 0, len(objs))
	for _, item := range objs {
		i := expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroups(item)
		items = append(items, *i)
	}

	return items
}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroups(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroups {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroups
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroups{
		Resources:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(obj["resources"]),
		InventoryFilters: expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(obj["inventory_filters"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsArray(objs []osconfig.OSPolicyAssignmentOSPoliciesResourceGroups) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroups(&item)
		items = append(items, i)
	}

	return items
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroups(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroups) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"resources":         flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(obj.Resources),
		"inventory_filters": flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(obj.InventoryFilters),
	}

	return transformed

}
func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(o interface{}) []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0, len(objs))
	for _, item := range objs {
		i := expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResources(item)
		items = append(items, *i)
	}

	return items
}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResources(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources{
		Id:         dcl.String(obj["id"].(string)),
		Exec:       expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(obj["exec"]),
		File:       expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(obj["file"]),
		Pkg:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(obj["pkg"]),
		Repository: expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(obj["repository"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(objs []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResources(&item)
		items = append(items, i)
	}

	return items
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResources(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"id":         obj.Id,
		"exec":       flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(obj.Exec),
		"file":       flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(obj.File),
		"pkg":        flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(obj.Pkg),
		"repository": flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(obj.Repository),
	}

	return transformed

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{
		Validate: expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(obj["validate"]),
		Enforce:  expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(obj["enforce"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"validate": flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(obj.Validate),
		"enforce":  flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(obj.Enforce),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{
		Interpreter:    osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumRef(obj["interpreter"].(string)),
		Args:           expandStringArray(obj["args"]),
		File:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(obj["file"]),
		OutputFilePath: dcl.String(obj["output_file_path"].(string)),
		Script:         dcl.String(obj["script"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"interpreter":      obj.Interpreter,
		"args":             obj.Args,
		"file":             flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(obj.File),
		"output_file_path": obj.OutputFilePath,
		"script":           obj.Script,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{
		AllowInsecure: dcl.Bool(obj["allow_insecure"].(bool)),
		Gcs:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(obj["gcs"]),
		LocalPath:     dcl.String(obj["local_path"].(string)),
		Remote:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(obj["remote"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_insecure": obj.AllowInsecure,
		"gcs":            flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(obj.Gcs),
		"local_path":     obj.LocalPath,
		"remote":         flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(obj.Remote),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{
		Bucket:     dcl.String(obj["bucket"].(string)),
		Object:     dcl.String(obj["object"].(string)),
		Generation: dcl.Int64(int64(obj["generation"].(int))),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"bucket":     obj.Bucket,
		"object":     obj.Object,
		"generation": obj.Generation,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{
		Uri:            dcl.String(obj["uri"].(string)),
		Sha256Checksum: dcl.String(obj["sha256_checksum"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri":             obj.Uri,
		"sha256_checksum": obj.Sha256Checksum,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{
		Interpreter:    osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumRef(obj["interpreter"].(string)),
		Args:           expandStringArray(obj["args"]),
		File:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(obj["file"]),
		OutputFilePath: dcl.String(obj["output_file_path"].(string)),
		Script:         dcl.String(obj["script"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"interpreter":      obj.Interpreter,
		"args":             obj.Args,
		"file":             flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(obj.File),
		"output_file_path": obj.OutputFilePath,
		"script":           obj.Script,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{
		AllowInsecure: dcl.Bool(obj["allow_insecure"].(bool)),
		Gcs:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(obj["gcs"]),
		LocalPath:     dcl.String(obj["local_path"].(string)),
		Remote:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(obj["remote"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_insecure": obj.AllowInsecure,
		"gcs":            flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(obj.Gcs),
		"local_path":     obj.LocalPath,
		"remote":         flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(obj.Remote),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{
		Bucket:     dcl.String(obj["bucket"].(string)),
		Object:     dcl.String(obj["object"].(string)),
		Generation: dcl.Int64(int64(obj["generation"].(int))),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"bucket":     obj.Bucket,
		"object":     obj.Object,
		"generation": obj.Generation,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{
		Uri:            dcl.String(obj["uri"].(string)),
		Sha256Checksum: dcl.String(obj["sha256_checksum"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri":             obj.Uri,
		"sha256_checksum": obj.Sha256Checksum,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{
		Path:    dcl.String(obj["path"].(string)),
		State:   osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumRef(obj["state"].(string)),
		Content: dcl.String(obj["content"].(string)),
		File:    expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(obj["file"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"path":        obj.Path,
		"state":       obj.State,
		"content":     obj.Content,
		"file":        flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(obj.File),
		"permissions": obj.Permissions,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{
		AllowInsecure: dcl.Bool(obj["allow_insecure"].(bool)),
		Gcs:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(obj["gcs"]),
		LocalPath:     dcl.String(obj["local_path"].(string)),
		Remote:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(obj["remote"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_insecure": obj.AllowInsecure,
		"gcs":            flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(obj.Gcs),
		"local_path":     obj.LocalPath,
		"remote":         flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(obj.Remote),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{
		Bucket:     dcl.String(obj["bucket"].(string)),
		Object:     dcl.String(obj["object"].(string)),
		Generation: dcl.Int64(int64(obj["generation"].(int))),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"bucket":     obj.Bucket,
		"object":     obj.Object,
		"generation": obj.Generation,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{
		Uri:            dcl.String(obj["uri"].(string)),
		Sha256Checksum: dcl.String(obj["sha256_checksum"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri":             obj.Uri,
		"sha256_checksum": obj.Sha256Checksum,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{
		DesiredState: osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumRef(obj["desired_state"].(string)),
		Apt:          expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(obj["apt"]),
		Deb:          expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(obj["deb"]),
		Googet:       expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(obj["googet"]),
		Msi:          expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(obj["msi"]),
		Rpm:          expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(obj["rpm"]),
		Yum:          expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(obj["yum"]),
		Zypper:       expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(obj["zypper"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"desired_state": obj.DesiredState,
		"apt":           flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(obj.Apt),
		"deb":           flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(obj.Deb),
		"googet":        flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(obj.Googet),
		"msi":           flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(obj.Msi),
		"rpm":           flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(obj.Rpm),
		"yum":           flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(obj.Yum),
		"zypper":        flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(obj.Zypper),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{
		Source:   expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(obj["source"]),
		PullDeps: dcl.Bool(obj["pull_deps"].(bool)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"source":    flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(obj.Source),
		"pull_deps": obj.PullDeps,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{
		AllowInsecure: dcl.Bool(obj["allow_insecure"].(bool)),
		Gcs:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(obj["gcs"]),
		LocalPath:     dcl.String(obj["local_path"].(string)),
		Remote:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(obj["remote"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_insecure": obj.AllowInsecure,
		"gcs":            flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(obj.Gcs),
		"local_path":     obj.LocalPath,
		"remote":         flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(obj.Remote),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{
		Bucket:     dcl.String(obj["bucket"].(string)),
		Object:     dcl.String(obj["object"].(string)),
		Generation: dcl.Int64(int64(obj["generation"].(int))),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"bucket":     obj.Bucket,
		"object":     obj.Object,
		"generation": obj.Generation,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{
		Uri:            dcl.String(obj["uri"].(string)),
		Sha256Checksum: dcl.String(obj["sha256_checksum"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri":             obj.Uri,
		"sha256_checksum": obj.Sha256Checksum,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{
		Source:     expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(obj["source"]),
		Properties: expandStringArray(obj["properties"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"source":     flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(obj.Source),
		"properties": obj.Properties,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{
		AllowInsecure: dcl.Bool(obj["allow_insecure"].(bool)),
		Gcs:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(obj["gcs"]),
		LocalPath:     dcl.String(obj["local_path"].(string)),
		Remote:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(obj["remote"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_insecure": obj.AllowInsecure,
		"gcs":            flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(obj.Gcs),
		"local_path":     obj.LocalPath,
		"remote":         flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(obj.Remote),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{
		Bucket:     dcl.String(obj["bucket"].(string)),
		Object:     dcl.String(obj["object"].(string)),
		Generation: dcl.Int64(int64(obj["generation"].(int))),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"bucket":     obj.Bucket,
		"object":     obj.Object,
		"generation": obj.Generation,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{
		Uri:            dcl.String(obj["uri"].(string)),
		Sha256Checksum: dcl.String(obj["sha256_checksum"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri":             obj.Uri,
		"sha256_checksum": obj.Sha256Checksum,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{
		Source:   expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(obj["source"]),
		PullDeps: dcl.Bool(obj["pull_deps"].(bool)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"source":    flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(obj.Source),
		"pull_deps": obj.PullDeps,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{
		AllowInsecure: dcl.Bool(obj["allow_insecure"].(bool)),
		Gcs:           expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(obj["gcs"]),
		LocalPath:     dcl.String(obj["local_path"].(string)),
		Remote:        expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(obj["remote"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_insecure": obj.AllowInsecure,
		"gcs":            flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(obj.Gcs),
		"local_path":     obj.LocalPath,
		"remote":         flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(obj.Remote),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{
		Bucket:     dcl.String(obj["bucket"].(string)),
		Object:     dcl.String(obj["object"].(string)),
		Generation: dcl.Int64(int64(obj["generation"].(int))),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"bucket":     obj.Bucket,
		"object":     obj.Object,
		"generation": obj.Generation,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{
		Uri:            dcl.String(obj["uri"].(string)),
		Sha256Checksum: dcl.String(obj["sha256_checksum"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri":             obj.Uri,
		"sha256_checksum": obj.Sha256Checksum,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{
		Apt:    expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(obj["apt"]),
		Goo:    expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(obj["goo"]),
		Yum:    expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(obj["yum"]),
		Zypper: expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(obj["zypper"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"apt":    flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(obj.Apt),
		"goo":    flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(obj.Goo),
		"yum":    flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(obj.Yum),
		"zypper": flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(obj.Zypper),
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{
		ArchiveType:  osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumRef(obj["archive_type"].(string)),
		Components:   expandStringArray(obj["components"]),
		Distribution: dcl.String(obj["distribution"].(string)),
		Uri:          dcl.String(obj["uri"].(string)),
		GpgKey:       dcl.String(obj["gpg_key"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"archive_type": obj.ArchiveType,
		"components":   obj.Components,
		"distribution": obj.Distribution,
		"uri":          obj.Uri,
		"gpg_key":      obj.GpgKey,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{
		Name: dcl.String(obj["name"].(string)),
		Url:  dcl.String(obj["url"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
		"url":  obj.Url,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{
		BaseUrl:     dcl.String(obj["base_url"].(string)),
		Id:          dcl.String(obj["id"].(string)),
		DisplayName: dcl.String(obj["display_name"].(string)),
		GpgKeys:     expandStringArray(obj["gpg_keys"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"base_url":     obj.BaseUrl,
		"id":           obj.Id,
		"display_name": obj.DisplayName,
		"gpg_keys":     obj.GpgKeys,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{
		BaseUrl:     dcl.String(obj["base_url"].(string)),
		Id:          dcl.String(obj["id"].(string)),
		DisplayName: dcl.String(obj["display_name"].(string)),
		GpgKeys:     expandStringArray(obj["gpg_keys"]),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"base_url":     obj.BaseUrl,
		"id":           obj.Id,
		"display_name": obj.DisplayName,
		"gpg_keys":     obj.GpgKeys,
	}

	return []interface{}{transformed}

}
func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(o interface{}) []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0, len(objs))
	for _, item := range objs {
		i := expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(item)
		items = append(items, *i)
	}

	return items
}

func expandOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{
		OSShortName: dcl.String(obj["os_short_name"].(string)),
		OSVersion:   dcl.String(obj["os_version"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(objs []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(&item)
		items = append(items, i)
	}

	return items
}

func flattenOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"os_short_name": obj.OSShortName,
		"os_version":    obj.OSVersion,
	}

	return transformed

}

func expandOsConfigOsPolicyAssignmentRollout(o interface{}) *osconfig.OSPolicyAssignmentRollout {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentRollout
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentRollout
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentRollout{
		DisruptionBudget: expandOsConfigOsPolicyAssignmentRolloutDisruptionBudget(obj["disruption_budget"]),
		MinWaitDuration:  dcl.String(obj["min_wait_duration"].(string)),
	}
}

func flattenOsConfigOsPolicyAssignmentRollout(obj *osconfig.OSPolicyAssignmentRollout) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"disruption_budget": flattenOsConfigOsPolicyAssignmentRolloutDisruptionBudget(obj.DisruptionBudget),
		"min_wait_duration": obj.MinWaitDuration,
	}

	return []interface{}{transformed}

}

func expandOsConfigOsPolicyAssignmentRolloutDisruptionBudget(o interface{}) *osconfig.OSPolicyAssignmentRolloutDisruptionBudget {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentRolloutDisruptionBudget
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return osconfig.EmptyOSPolicyAssignmentRolloutDisruptionBudget
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentRolloutDisruptionBudget{
		Fixed:   dcl.Int64(int64(obj["fixed"].(int))),
		Percent: dcl.Int64(int64(obj["percent"].(int))),
	}
}

func flattenOsConfigOsPolicyAssignmentRolloutDisruptionBudget(obj *osconfig.OSPolicyAssignmentRolloutDisruptionBudget) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"fixed":   obj.Fixed,
		"percent": obj.Percent,
	}

	return []interface{}{transformed}

}
