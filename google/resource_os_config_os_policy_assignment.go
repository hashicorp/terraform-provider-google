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

func resourceOSConfigOSPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOSConfigOSPolicyAssignmentCreate,
		Read:   resourceOSConfigOSPolicyAssignmentRead,
		Update: resourceOSConfigOSPolicyAssignmentUpdate,
		Delete: resourceOSConfigOSPolicyAssignmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceOSConfigOSPolicyAssignmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_filter": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Filter to select VMs.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentInstanceFilterSchema(),
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
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesSchema(),
			},

			"rollout": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Rollout to deploy the OS policy assignment. A rollout is triggered in the following situations: 1) OSPolicyAssignment is created. 2) OSPolicyAssignment is updated and the update contains changes to one of the following fields: - instance_filter - os_policies 3) OSPolicyAssignment is deleted.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentRolloutSchema(),
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

func OSConfigOSPolicyAssignmentInstanceFilterSchema() *schema.Resource {
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
				Elem:        OSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsSchema(),
			},

			"inclusion_labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of label sets used for VM inclusion. If the list has more than one `LabelSet`, the VM is included if any of the label sets are applicable for the VM.",
				Elem:        OSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsSchema(),
			},

			"inventories": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of inventories to select VMs. A VM is selected if its inventory data matches at least one of the following inventories.",
				Elem:        OSConfigOSPolicyAssignmentInstanceFilterInventoriesSchema(),
			},
		},
	}
}

func OSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentInstanceFilterInventoriesSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesSchema() *schema.Resource {
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
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsSchema(),
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. List of resources configured for this resource group. The resources are executed in the exact order specified here.",
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSchema(),
			},

			"inventory_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of inventory filters for the resource group. The resources in this resource group are applied to the target VM if it satisfies at least one of the following inventory filters. For example, to apply this resource group to VMs running either `RHEL` or `CentOS` operating systems, specify 2 items for the list with following values: inventory_filters[0].os_short_name='rhel' and inventory_filters[1].os_short_name='centos' If the list is empty, this resource group will be applied to the target VM unconditionally.",
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSchema(),
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSchema() *schema.Resource {
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
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSchema(),
			},

			"file": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "File resource",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSchema(),
			},

			"pkg": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Package resource",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSchema(),
			},

			"repository": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Package repository resource",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySchema(),
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"validate": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. What to run to validate this resource is in the desired state. An exit code of 100 indicates \"in desired state\", and exit code of 101 indicates \"not in desired state\". Any other exit code indicates a failure running validate.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentExecSchema(),
			},

			"enforce": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Required. What to run to validate this resource is in the desired state. An exit code of 100 indicates \"in desired state\", and exit code of 101 indicates \"not in desired state\". Any other exit code indicates a failure running validate.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentExecSchema(),
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSchema() *schema.Resource {
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
				Description: "Required. A deb package.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentFileSchema(),
			},

			"permissions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consists of three octal digits which represent, in order, the permissions of the owner, group, and other users for the file (similarly to the numeric mode used in the linux chmod utility). Each digit represents a three bit number with the 4 bit corresponding to the read permissions, the 2 bit corresponds to the write bit, and the one bit corresponds to the execute permission. Default behavior is 755. Below are some examples of permissions and their associated values: read, write, and execute: 7 read and execute: 5 read and write: 6 read only: 4",
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSchema() *schema.Resource {
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
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSchema(),
			},

			"deb": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A deb package file.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSchema(),
			},

			"googet": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A package managed by GooGet.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSchema(),
			},

			"msi": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An MSI package.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSchema(),
			},

			"rpm": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An rpm package file.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSchema(),
			},

			"yum": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A package managed by YUM.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSchema(),
			},

			"zypper": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A package managed by Zypper.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSchema(),
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. A deb package.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentFileSchema(),
			},

			"pull_deps": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether dependencies should also be installed. - install when false: `dpkg -i package` - install when true: `apt-get update && apt-get -y install package.deb`",
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. A deb package.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentFileSchema(),
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. A deb package.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentFileSchema(),
			},

			"pull_deps": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether dependencies should also be installed. - install when false: `rpm --upgrade --replacepkgs package.rpm` - install when true: `yum -y install package.rpm` or `zypper -y install package.rpm`",
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"apt": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An Apt Repository.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSchema(),
			},

			"goo": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Goo Repository.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSchema(),
			},

			"yum": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Yum Repository.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSchema(),
			},

			"zypper": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A Zypper Repository.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSchema(),
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentRolloutSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disruption_budget": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The maximum number (or percentage) of VMs per zone to disrupt at any given moment.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentRolloutDisruptionBudgetSchema(),
			},

			"min_wait_duration": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. This determines the minimum duration of time to wait after the configuration changes are applied through the current rollout. A VM continues to count towards the `disruption_budget` at least until this duration of time has passed after configuration changes are applied.",
			},
		},
	}
}

func OSConfigOSPolicyAssignmentRolloutDisruptionBudgetSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPolicyAssignmentFileSchema() *schema.Resource {
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
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentFileGcsSchema(),
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
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentFileRemoteSchema(),
			},
		},
	}
}

func OSConfigOSPolicyAssignmentOSPolicyAssignmentFileGcsSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPolicyAssignmentFileRemoteSchema() *schema.Resource {
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

func OSConfigOSPolicyAssignmentOSPolicyAssignmentExecSchema() *schema.Resource {
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
				Description: "Required. A deb package.",
				MaxItems:    1,
				Elem:        OSConfigOSPolicyAssignmentOSPolicyAssignmentFileSchema(),
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

func resourceOSConfigOSPolicyAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter: expandOSConfigOSPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		OSPolicies:     expandOSConfigOSPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:        expandOSConfigOSPolicyAssignmentRollout(d.Get("rollout")),
		Description:    dcl.String(d.Get("description").(string)),
		Project:        dcl.String(project),
	}

	id, err := replaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}")
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	createDirective := CreateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLOSConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyOSPolicyAssignment(context.Background(), obj, createDirective...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating OSPolicyAssignment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating OSPolicyAssignment %q: %#v", d.Id(), res)

	return resourceOSConfigOSPolicyAssignmentRead(d, meta)
}

func resourceOSConfigOSPolicyAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter: expandOSConfigOSPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		OSPolicies:     expandOSConfigOSPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:        expandOSConfigOSPolicyAssignmentRollout(d.Get("rollout")),
		Description:    dcl.String(d.Get("description").(string)),
		Project:        dcl.String(project),
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
	client := NewDCLOSConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := replaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetOSPolicyAssignment(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("OSConfigOSPolicyAssignment %q", d.Id())
		return handleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("instance_filter", flattenOSConfigOSPolicyAssignmentInstanceFilter(res.InstanceFilter)); err != nil {
		return fmt.Errorf("error setting instance_filter in state: %s", err)
	}
	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("os_policies", flattenOSConfigOSPolicyAssignmentOSPoliciesArray(res.OSPolicies)); err != nil {
		return fmt.Errorf("error setting os_policies in state: %s", err)
	}
	if err = d.Set("rollout", flattenOSConfigOSPolicyAssignmentRollout(res.Rollout)); err != nil {
		return fmt.Errorf("error setting rollout in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
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
func resourceOSConfigOSPolicyAssignmentUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter: expandOSConfigOSPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		OSPolicies:     expandOSConfigOSPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:        expandOSConfigOSPolicyAssignmentRollout(d.Get("rollout")),
		Description:    dcl.String(d.Get("description").(string)),
		Project:        dcl.String(project),
	}
	directive := UpdateDirective
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := NewDCLOSConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
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

	return resourceOSConfigOSPolicyAssignmentRead(d, meta)
}

func resourceOSConfigOSPolicyAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := &osconfig.OSPolicyAssignment{
		InstanceFilter: expandOSConfigOSPolicyAssignmentInstanceFilter(d.Get("instance_filter")),
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		OSPolicies:     expandOSConfigOSPolicyAssignmentOSPoliciesArray(d.Get("os_policies")),
		Rollout:        expandOSConfigOSPolicyAssignmentRollout(d.Get("rollout")),
		Description:    dcl.String(d.Get("description").(string)),
		Project:        dcl.String(project),
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
	client := NewDCLOSConfigClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
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

func resourceOSConfigOSPolicyAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func expandOSConfigOSPolicyAssignmentInstanceFilter(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilter {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilter
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilter
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilter{
		All:             dcl.Bool(obj["all"].(bool)),
		ExclusionLabels: expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsArray(obj["exclusion_labels"]),
		InclusionLabels: expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsArray(obj["inclusion_labels"]),
		Inventories:     expandOSConfigOSPolicyAssignmentInstanceFilterInventoriesArray(obj["inventories"]),
	}
}

func flattenOSConfigOSPolicyAssignmentInstanceFilter(obj *osconfig.OSPolicyAssignmentInstanceFilter) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"all":              obj.All,
		"exclusion_labels": flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsArray(obj.ExclusionLabels),
		"inclusion_labels": flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsArray(obj.InclusionLabels),
		"inventories":      flattenOSConfigOSPolicyAssignmentInstanceFilterInventoriesArray(obj.Inventories),
	}

	return []interface{}{transformed}

}
func expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsArray(o interface{}) []osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels, 0, len(objs))
	for _, item := range objs {
		i := expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(item)
		items = append(items, *i)
	}

	return items
}

func expandOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilterExclusionLabels
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels{
		Labels: checkStringMap(obj["labels"]),
	}
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabelsArray(objs []osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(&item)
		items = append(items, i)
	}

	return items
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterExclusionLabels(obj *osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"labels": obj.Labels,
	}

	return transformed

}
func expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsArray(o interface{}) []osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels, 0, len(objs))
	for _, item := range objs {
		i := expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(item)
		items = append(items, *i)
	}

	return items
}

func expandOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilterInclusionLabels
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels{
		Labels: checkStringMap(obj["labels"]),
	}
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabelsArray(objs []osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(&item)
		items = append(items, i)
	}

	return items
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterInclusionLabels(obj *osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"labels": obj.Labels,
	}

	return transformed

}
func expandOSConfigOSPolicyAssignmentInstanceFilterInventoriesArray(o interface{}) []osconfig.OSPolicyAssignmentInstanceFilterInventories {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInventories, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]osconfig.OSPolicyAssignmentInstanceFilterInventories, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentInstanceFilterInventories, 0, len(objs))
	for _, item := range objs {
		i := expandOSConfigOSPolicyAssignmentInstanceFilterInventories(item)
		items = append(items, *i)
	}

	return items
}

func expandOSConfigOSPolicyAssignmentInstanceFilterInventories(o interface{}) *osconfig.OSPolicyAssignmentInstanceFilterInventories {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentInstanceFilterInventories
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentInstanceFilterInventories{
		OSShortName: dcl.String(obj["os_short_name"].(string)),
		OSVersion:   dcl.String(obj["os_version"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterInventoriesArray(objs []osconfig.OSPolicyAssignmentInstanceFilterInventories) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOSConfigOSPolicyAssignmentInstanceFilterInventories(&item)
		items = append(items, i)
	}

	return items
}

func flattenOSConfigOSPolicyAssignmentInstanceFilterInventories(obj *osconfig.OSPolicyAssignmentInstanceFilterInventories) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"os_short_name": obj.OSShortName,
		"os_version":    obj.OSVersion,
	}

	return transformed

}
func expandOSConfigOSPolicyAssignmentOSPoliciesArray(o interface{}) []osconfig.OSPolicyAssignmentOSPolicies {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPolicies, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]osconfig.OSPolicyAssignmentOSPolicies, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPolicies, 0, len(objs))
	for _, item := range objs {
		i := expandOSConfigOSPolicyAssignmentOSPolicies(item)
		items = append(items, *i)
	}

	return items
}

func expandOSConfigOSPolicyAssignmentOSPolicies(o interface{}) *osconfig.OSPolicyAssignmentOSPolicies {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPolicies
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPolicies{
		Id:                        dcl.String(obj["id"].(string)),
		Mode:                      osconfig.OSPolicyAssignmentOSPoliciesModeEnumRef(obj["mode"].(string)),
		ResourceGroups:            expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsArray(obj["resource_groups"]),
		AllowNoResourceGroupMatch: dcl.Bool(obj["allow_no_resource_group_match"].(bool)),
		Description:               dcl.String(obj["description"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesArray(objs []osconfig.OSPolicyAssignmentOSPolicies) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOSConfigOSPolicyAssignmentOSPolicies(&item)
		items = append(items, i)
	}

	return items
}

func flattenOSConfigOSPolicyAssignmentOSPolicies(obj *osconfig.OSPolicyAssignmentOSPolicies) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"id":                            obj.Id,
		"mode":                          obj.Mode,
		"resource_groups":               flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsArray(obj.ResourceGroups),
		"allow_no_resource_group_match": obj.AllowNoResourceGroupMatch,
		"description":                   obj.Description,
	}

	return transformed

}
func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsArray(o interface{}) []osconfig.OSPolicyAssignmentOSPoliciesResourceGroups {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroups, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroups, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroups, 0, len(objs))
	for _, item := range objs {
		i := expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroups(item)
		items = append(items, *i)
	}

	return items
}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroups(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroups {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroups
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroups{
		Resources:        expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(obj["resources"]),
		InventoryFilters: expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(obj["inventory_filters"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsArray(objs []osconfig.OSPolicyAssignmentOSPoliciesResourceGroups) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroups(&item)
		items = append(items, i)
	}

	return items
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroups(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroups) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"resources":         flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(obj.Resources),
		"inventory_filters": flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(obj.InventoryFilters),
	}

	return transformed

}
func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(o interface{}) []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0, len(objs))
	for _, item := range objs {
		i := expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResources(item)
		items = append(items, *i)
	}

	return items
}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResources(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources{
		Id:         dcl.String(obj["id"].(string)),
		Exec:       expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(obj["exec"]),
		File:       expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(obj["file"]),
		Pkg:        expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(obj["pkg"]),
		Repository: expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(obj["repository"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesArray(objs []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResources(&item)
		items = append(items, i)
	}

	return items
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResources(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"id":         obj.Id,
		"exec":       flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(obj.Exec),
		"file":       flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(obj.File),
		"pkg":        flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(obj.Pkg),
		"repository": flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(obj.Repository),
	}

	return transformed

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{
		Validate: expandOSConfigOSPolicyAssignmentOSPolicyAssignmentExec(obj["validate"]),
		Enforce:  expandOSConfigOSPolicyAssignmentOSPolicyAssignmentExec(obj["enforce"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"validate": flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentExec(obj.Validate),
		"enforce":  flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentExec(obj.Enforce),
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{
		Path:    dcl.String(obj["path"].(string)),
		State:   osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumRef(obj["state"].(string)),
		Content: dcl.String(obj["content"].(string)),
		File:    expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj["file"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"path":        obj.Path,
		"state":       obj.State,
		"content":     obj.Content,
		"file":        flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj.File),
		"permissions": obj.Permissions,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{
		DesiredState: osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumRef(obj["desired_state"].(string)),
		Apt:          expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(obj["apt"]),
		Deb:          expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(obj["deb"]),
		Googet:       expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(obj["googet"]),
		Msi:          expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(obj["msi"]),
		Rpm:          expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(obj["rpm"]),
		Yum:          expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(obj["yum"]),
		Zypper:       expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(obj["zypper"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"desired_state": obj.DesiredState,
		"apt":           flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(obj.Apt),
		"deb":           flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(obj.Deb),
		"googet":        flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(obj.Googet),
		"msi":           flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(obj.Msi),
		"rpm":           flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(obj.Rpm),
		"yum":           flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(obj.Yum),
		"zypper":        flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(obj.Zypper),
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{
		Source:   expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj["source"]),
		PullDeps: dcl.Bool(obj["pull_deps"].(bool)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"source":    flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj.Source),
		"pull_deps": obj.PullDeps,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{
		Source:     expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj["source"]),
		Properties: expandStringArray(obj["properties"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"source":     flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj.Source),
		"properties": obj.Properties,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{
		Source:   expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj["source"]),
		PullDeps: dcl.Bool(obj["pull_deps"].(bool)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"source":    flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj.Source),
		"pull_deps": obj.PullDeps,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{
		Name: dcl.String(obj["name"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{
		Apt:    expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(obj["apt"]),
		Goo:    expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(obj["goo"]),
		Yum:    expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(obj["yum"]),
		Zypper: expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(obj["zypper"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"apt":    flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(obj.Apt),
		"goo":    flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(obj.Goo),
		"yum":    flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(obj.Yum),
		"zypper": flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(obj.Zypper),
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
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

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) interface{} {
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

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{
		Name: dcl.String(obj["name"].(string)),
		Url:  dcl.String(obj["url"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"name": obj.Name,
		"url":  obj.Url,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
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

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) interface{} {
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

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
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

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) interface{} {
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
func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(o interface{}) []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if o == nil {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 {
		return make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0)
	}

	items := make([]osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0, len(objs))
	for _, item := range objs {
		i := expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(item)
		items = append(items, *i)
	}

	return items
}

func expandOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(o interface{}) *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	}

	obj := o.(map[string]interface{})
	return &osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{
		OSShortName: dcl.String(obj["os_short_name"].(string)),
		OSVersion:   dcl.String(obj["os_version"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersArray(objs []osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(&item)
		items = append(items, i)
	}

	return items
}

func flattenOSConfigOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(obj *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"os_short_name": obj.OSShortName,
		"os_version":    obj.OSVersion,
	}

	return transformed

}

func expandOSConfigOSPolicyAssignmentRollout(o interface{}) *osconfig.OSPolicyAssignmentRollout {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentRollout
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentRollout
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentRollout{
		DisruptionBudget: expandOSConfigOSPolicyAssignmentRolloutDisruptionBudget(obj["disruption_budget"]),
		MinWaitDuration:  dcl.String(obj["min_wait_duration"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentRollout(obj *osconfig.OSPolicyAssignmentRollout) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"disruption_budget": flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudget(obj.DisruptionBudget),
		"min_wait_duration": obj.MinWaitDuration,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentRolloutDisruptionBudget(o interface{}) *osconfig.OSPolicyAssignmentRolloutDisruptionBudget {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentRolloutDisruptionBudget
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentRolloutDisruptionBudget
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentRolloutDisruptionBudget{
		Fixed:   dcl.Int64(int64(obj["fixed"].(int))),
		Percent: dcl.Int64(int64(obj["percent"].(int))),
	}
}

func flattenOSConfigOSPolicyAssignmentRolloutDisruptionBudget(obj *osconfig.OSPolicyAssignmentRolloutDisruptionBudget) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"fixed":   obj.Fixed,
		"percent": obj.Percent,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(o interface{}) *osconfig.OSPolicyAssignmentFile {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentFile
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentFile
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentFile{
		AllowInsecure: dcl.Bool(obj["allow_insecure"].(bool)),
		Gcs:           expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFileGcs(obj["gcs"]),
		LocalPath:     dcl.String(obj["local_path"].(string)),
		Remote:        expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFileRemote(obj["remote"]),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj *osconfig.OSPolicyAssignmentFile) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"allow_insecure": obj.AllowInsecure,
		"gcs":            flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFileGcs(obj.Gcs),
		"local_path":     obj.LocalPath,
		"remote":         flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFileRemote(obj.Remote),
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFileGcs(o interface{}) *osconfig.OSPolicyAssignmentFileGcs {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentFileGcs
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentFileGcs
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentFileGcs{
		Bucket:     dcl.String(obj["bucket"].(string)),
		Object:     dcl.String(obj["object"].(string)),
		Generation: dcl.Int64(int64(obj["generation"].(int))),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFileGcs(obj *osconfig.OSPolicyAssignmentFileGcs) interface{} {
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

func expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFileRemote(o interface{}) *osconfig.OSPolicyAssignmentFileRemote {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentFileRemote
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentFileRemote
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentFileRemote{
		Uri:            dcl.String(obj["uri"].(string)),
		Sha256Checksum: dcl.String(obj["sha256_checksum"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFileRemote(obj *osconfig.OSPolicyAssignmentFileRemote) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"uri":             obj.Uri,
		"sha256_checksum": obj.Sha256Checksum,
	}

	return []interface{}{transformed}

}

func expandOSConfigOSPolicyAssignmentOSPolicyAssignmentExec(o interface{}) *osconfig.OSPolicyAssignmentExec {
	if o == nil {
		return osconfig.EmptyOSPolicyAssignmentExec
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 {
		return osconfig.EmptyOSPolicyAssignmentExec
	}
	obj := objArr[0].(map[string]interface{})
	return &osconfig.OSPolicyAssignmentExec{
		Interpreter:    osconfig.OSPolicyAssignmentExecInterpreterEnumRef(obj["interpreter"].(string)),
		Args:           expandStringArray(obj["args"]),
		File:           expandOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj["file"]),
		OutputFilePath: dcl.String(obj["output_file_path"].(string)),
		Script:         dcl.String(obj["script"].(string)),
	}
}

func flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentExec(obj *osconfig.OSPolicyAssignmentExec) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"interpreter":      obj.Interpreter,
		"args":             obj.Args,
		"file":             flattenOSConfigOSPolicyAssignmentOSPolicyAssignmentFile(obj.File),
		"output_file_path": obj.OutputFilePath,
		"script":           obj.Script,
	}

	return []interface{}{transformed}

}
